package main

import (
	"fmt"
	"io/ioutil"

	"github.com/go-yaml/yaml"
)

type Config struct {
	OpenApiKey string `yaml:"OpenApiKey"`
	BotToken   string `yaml:"BotToken"`
	Port       string `yaml:"Port"`
}

func ReadConfig(fileaname string) (*Config, error) {
	var config Config

	yamlFile, err := ioutil.ReadFile(fileaname)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML file: %v", err)
	}

	return &config, nil
}
