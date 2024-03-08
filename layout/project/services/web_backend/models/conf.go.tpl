package models

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Db struct {
	Type string `yaml:"type"`
	Dsn  string `yaml:"dsn"`
}

type Net struct {
	Address string `yaml:"address"`
}

type RunConf struct {
	Db  Db  `yaml:"db"`
	Net Net `yaml:"net"`
}

var Conf *RunConf

func LoadRunConf(file string) error {
	Conf = &RunConf{}
	buff, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(buff, Conf)
	if err != nil {
		return err
	}

	return nil
}
