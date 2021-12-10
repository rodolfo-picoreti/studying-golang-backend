package config

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type rootConfig struct {
	Service serviceConfig `yaml:"service"`
	DB      dbConfig      `yaml:"db"`
	Cache   cacheConfig   `yaml:"cache"`
	Tracing tracingConfig `yaml:"tracing"`
}

type serviceConfig struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type dbConfig struct {
	URI    string `yaml:"uri"`
	Prefix string `yaml:"prefix"`
}

type cacheConfig struct {
	Hostname          string        `yaml:"hostname"`
	Password          string        `yaml:"password"`
	DefaultExpiration time.Duration `yaml:"defaultExpiration"`
}

type tracingConfig struct {
	Target string `yaml:"target"`
}

func getConfigFilePath() string {
	if value, ok := os.LookupEnv("CONFIG_PATH"); ok {
		return value
	}
	return "./config.yaml"
}

func ReadConfig() *rootConfig {
	configPath := getConfigFilePath()

	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := &rootConfig{}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		panic(err)
	}

	return config
}
