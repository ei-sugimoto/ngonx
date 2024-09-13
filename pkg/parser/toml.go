package parser

import "github.com/BurntSushi/toml"

type Server struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func Parse() (map[string]Server, error) {
	ServerMap := map[string]Server{}

	_, err := toml.DecodeFile("/var/lib/ngonx/config.toml", &ServerMap)
	if err != nil {
		return nil, err
	}

	return ServerMap, nil
}
