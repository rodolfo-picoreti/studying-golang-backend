package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type rootConfig struct {
	DB dbConfig `yaml:"db"`
}

type dbConfig struct {
	URI    string `yaml:"uri"`
	Prefix string `yaml:"prefix"`
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
