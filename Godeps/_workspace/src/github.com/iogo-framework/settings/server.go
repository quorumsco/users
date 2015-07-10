package settings

import (
	"encoding/json"
	"fmt"

	"github.com/iogo-framework/logs"
)

type Server struct {
	Host string
	Port int
}

func (config Config) Server() (Server, error) {
	var server Server

	if err := json.Unmarshal(config.Components["server"], &server); err != nil {
		logs.Warning("%s: %s", err.Error(), "missing or wrong 'server' configuration, ignoring")
	}

	if server.Host == "" {
		server.Host = "0.0.0.0"
		logs.Warning("missing server 'host' configuration, assuming default value: '0.0.0.0'")
	}

	if server.Port == 0 {
		server.Port = 8080
		logs.Warning("missing server 'port' configuration, assuming default value: 8080")
	}

	return server, nil
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
