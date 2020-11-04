package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	MainConfigFile = "config.json"
)

type MainConfig struct {
	SourceNodeURL, DestinationNodeURL   string
	SourceDecimals, DestinationDecimals int64
	IBPortAddress                       string
	LUPortAddress                       string
}

func ParseMainConfig(confName string) (*MainConfig, error) {
	configName := MainConfigFile
	if confName != "" {
		configName = confName
	}

	data, err := ioutil.ReadFile(configName)

	if err != nil {
		return nil, err
	}

	var config *MainConfig
	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}
