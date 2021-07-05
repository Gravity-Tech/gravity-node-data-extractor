package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	MainConfigFile = "config.json"
)

const (
	MountedVolume = "/etc/extractor"
)

type MainConfig struct {
	SourceNodeURL, DestinationNodeURL   string
	SourceDecimals, DestinationDecimals int64
	IBPortAddress                       string
	LUPortAddress                       string
	Meta                                map[string]string
}

func ParseMainConfig(confName string) (*MainConfig, error) {
	configName := MainConfigFile
	if confName != "" {
		configName = confName
	}

  var err error
  var data []byte

	for _, path := range []string{".", MountedVolume} {
    data, err = ioutil.ReadFile(path + "/" + configName)

		if err == nil {
      break
		}
	}

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
