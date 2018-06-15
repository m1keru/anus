package tools

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

//Config -- Config for startup
var Config ServerConfig

//ServerConfig -- struct to handle config
type ServerConfig struct {
	Network struct {
		Listen string
	}
	Log struct {
		LogType  string `yaml:"type"`
		LogPath  string `yaml:"path"`
		LogLevel string `yaml:"level"`
	}
	SSH struct {
		Host     string
		Port     string
		Login    string
		Password string
		Keypath  string
	}
	Db struct {
		ProdPath string `yaml:"prod_path"`
		TestPath string `yaml:"test_path"`
	}
	Debug bool
}

//SetupConfig -- init vars
func SetupConfig(configFile string) {
	configfile, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Не могу прочесть конфиг-файл")
	}
	err = yaml.Unmarshal(configfile, &Config)
	if err != nil {
		log.Fatal("Не могу распарсить конфиг-файл")
	}
}
