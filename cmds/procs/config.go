package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ProcsConfig struct {
	Name        string `yaml:"name"`
	Cmd         string `yaml:"cmd"`
	Dir         string `yaml:"dir"`
	Log         string `yaml:"log"`
	StopRestart bool   `yaml:"stop_restart"`
	StopSignal  int    `yaml:"stop_signal"`
	ShowScreen  bool   `yaml:"show_screen"`
}

type NetConfig struct {
	Addr string `yaml:"addr"`
}

type RunConfig struct {
	Procs []ProcsConfig `yaml:"procs"`
	Net   []NetConfig   `yaml:"net"`
}

func LoadConf(configName string) (*RunConfig, error) {
	conf := &RunConfig{}
	buff, err := os.ReadFile(configName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buff, conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}
