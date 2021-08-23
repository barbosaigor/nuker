package bufwriter

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		w        io.Writer
		location string
	}

	type test struct {
		name    string
		args    args
		wantErr bool
		wantNil bool
	}

	tests := []test{
		{
			name: "when location and writer is provided, should build successfully",
			args: args{
				w:        io.Discard,
				location: "location.txt",
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "when location is provided but writer not, should return error",
			args: args{
				location: "location.txt",
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "when writer is provided but location not, should return error",
			args: args{
				w: io.Discard,
			},
			wantErr: true,
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := New(tc.args.w, tc.args.location)

			assert.Equal(t, tc.wantNil, got == nil)
			assert.Equal(t, tc.wantErr, err != nil)
		})
	}
}
