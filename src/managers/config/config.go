package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config ...
type Config struct {
	URL         string  `yaml:"url"`
	Method      string  `yaml:"method"`
	ContentType string  `yaml:"contenttype"`
	Body        string  `yaml:"body"`
	Clients     int     `yaml:"clients"`
	Timelife    int64   `yaml:"timelife"`
	Timeout     int64   `yaml:"timeout"`
	Logfile     string  `yaml:"logfile"`
	Proxy       []proxy `yaml:"proxy"`
}

type proxy struct {
	IP       string `yaml:"name"`
	Port     int    `yaml:"port"`
	Login    string `yaml:"login"`
	Password string `yaml:"password"`
}

var (
	yml []byte
	err error
)

// New ...
func New(file string) *Config {
	if yml, err = ioutil.ReadFile(file); err != nil {
		panic(err)
	}
	cfg := &Config{}
	e := yaml.Unmarshal(yml, cfg)

	if e != nil {
		panic(e.Error())
	}

	return cfg
}
