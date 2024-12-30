package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig `yaml:"server"`
	Database DBConfig     `yaml:"database"`
	Cors     CORSConfig   `yaml:"cors"`
	Search   SearchEngine `yaml:"searchEngine"`
}

type ServerConfig struct {
	Port            int           `yaml:"port"`
	MaxReadTimeout  time.Duration `yaml:"maxReadTimeout"`
	MaxWriteTimeout time.Duration `yaml:"maxWriteTimeout"`
	GracefulTimeout time.Duration `yaml:"gracefulTimeout"`
}

type DBConfig struct {
	Url string `yaml:"url"`
}

type CORSConfig struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`
	AllowedMethods   []string `yaml:"allowed_methods"`
	AllowedHeaders   []string `yaml:"allowed_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

type SearchEngine struct {
	Url    string `yaml:"url"`
	ApiKey string `yaml:"api_key"`
}

func Load(path string) (*Config, error) {
	var data, err = os.ReadFile(path)

	if err != nil {
		return nil, err
	}
	var configFile Config
	if err := yaml.Unmarshal(data, &configFile); err != nil {
		return nil, err
	}

	return &configFile, nil
}
