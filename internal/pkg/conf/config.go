package conf

import (
	"log"
)

func NewConfig(path string, config interface{}) interface{} {
	conf, err := Load(path, config)
	if err != nil {
		log.Fatal("Error load config", err)
	}

	return conf
}
