package parser

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Server struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	EndPoint string `toml:"endpoint"`
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

type URLAndEndPoint struct {
	URL      string
	EndPoint string
}

type URLAndEndPointList []URLAndEndPoint

func (s *ServerMap) GetURLList() URLAndEndPointList {
	list := make(URLAndEndPointList, 0, len(*s))
	for _, server := range *s {
		url := fmt.Sprintf("http://%s:%d", server.Host, server.Port)
		list = append(list, URLAndEndPoint{URL: url, EndPoint: server.EndPoint})
	}

	return list
}
