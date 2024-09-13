package parser

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func Parse() (map[string]Server, error) {
	ServerMap := map[string]Server{}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/var/lib/ngonx/config.toml"
	}

	_, err := toml.DecodeFile(configPath, &ServerMap)
	if err != nil {
		return nil, err
	}

	return ServerMap, nil
}
