package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

func YamlUnmarshal(data []byte) (*Config, error) {
	var cfg Config

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func FromYamlFile(filename string) (*Config, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return YamlUnmarshal(bytes)
}
