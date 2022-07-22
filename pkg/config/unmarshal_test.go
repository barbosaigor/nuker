package config

import (
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configData = []byte(`
stages:
- name: stage 1 
  steps:
  - name: step 1
    containers:
    - name: container 1
      duration: 10
      holdfor: 3
      min: 5
      max: 15
      network:
        protocol: http
        path: /resource
        method: POST
        host: "api.com"
        timeout: 22
        headers:
            "Content-Type": "application/json"
        body: '{"name":"bob"}'
`)

var invalidCfgData = []byte(`
stages
- name: stage 1 
  steps
  - name: step 1
    containers:
    - name: container 1
      duration: 10
      holdfor: 3
      min: 5
      max: 15
      network:
        protocol: http
        path: /resource
        method: POST
        host: "api.com"
        timeout: 22
        headers:
            "Content-Type": "application/json"
        body: '{"name":"bob"}'
`)

func TestYamlUnmarshal(t *testing.T) {
	type args struct {
		data []byte
	}
	type test struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}

	tests := []test{
		{
			name: "when input is valid, should unmarshal config successfully",
			args: args{
				data: configData,
			},
			want: &Config{
				Stages: []Stage{
					{
						Name: "stage 1",
						Steps: []Step{
							{
								Name: "step 1",
								Containers: []Container{
									{
										Name:     "container 1",
										Duration: 10,
										HoldFor:  3,
										Min:      5,
										Max:      15,
										Network: Network{
											Protocol: "http",
											Host:     "api.com",
											Path:     "/resource",
											Method:   "POST",
											Headers: map[string]string{
												"Content-Type": "application/json",
											},
											Timeout: 22,
											Body:    `{"name":"bob"}`,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "when input is not valid, should not build config and return error",
			args: args{
				data: invalidCfgData,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := YamlUnmarshal(tc.args.data)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}

func TestFromYamlFile(t *testing.T) {
	validFile, err := ioutil.TempFile("", "testyamlfile.yaml")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(validFile.Name())
	ioutil.WriteFile(validFile.Name(), configData, 0o644)

	invalidFile, err := ioutil.TempFile("", "testinvalidyamlfile.yaml")
	if err != nil {
		panic(err)
	}
	defer syscall.Unlink(invalidFile.Name())
	ioutil.WriteFile(invalidFile.Name(), invalidCfgData, 0o644)

	type args struct {
		fileName string
	}
	type test struct {
		name    string
		args    args
		want    *Config
		wantErr bool
	}

	tests := []test{
		{
			name: "when file is valid, should unmarshal config successfully",
			args: args{
				fileName: validFile.Name(),
			},
			want: &Config{
				Stages: []Stage{
					{
						Name: "stage 1",
						Steps: []Step{
							{
								Name: "step 1",
								Containers: []Container{
									{
										Name:     "container 1",
										Duration: 10,
										HoldFor:  3,
										Min:      5,
										Max:      15,
										Network: Network{
											Protocol: "http",
											Host:     "api.com",
											Path:     "/resource",
											Method:   "POST",
											Headers: map[string]string{
												"Content-Type": "application/json",
											},
											Timeout: 22,
											Body:    `{"name":"bob"}`,
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "when input is not valid, should not build config and return error",
			args: args{
				fileName: invalidFile.Name(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "when file doesn't exists, should not build config and return error",
			args: args{
				fileName: "unknow.yaml",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FromYamlFile(tc.args.fileName)
			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
