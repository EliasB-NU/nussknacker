package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type CFG struct {
	BotToken  string `yaml:"Bot-Token"`
	ChannelID string `yaml:"ChannelID"`
	SendTime  string `yaml:"SendTime"`
}

func GetConfig() *CFG {
	const file = "config.yaml"
	var config CFG

	cfgFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer func(cfgFile *os.File) {
		err := cfgFile.Close()
		if err != nil {
			panic(err)
		}
	}(cfgFile)

	yamlParser := yaml.NewDecoder(cfgFile)
	err = yamlParser.Decode(&config)
	if err != nil {
		panic(err)
	}

	return &config
}
