package parser

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type ServerMap map[string]Server

func NewServer() ServerMap {
	return make(ServerMap)
}

func (s *ServerMap) Parse() error {

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/var/lib/ngonx/config.toml"
	}

	_, err := toml.DecodeFile(configPath, &s)
	if err != nil {
		return err
	}

	return nil
}
