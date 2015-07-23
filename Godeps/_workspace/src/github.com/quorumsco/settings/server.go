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

func (config Config) CustomServer(name, defaultHost string, defaultPort int) (Server, error) {
	var server Server

	if err := json.Unmarshal(config.Components[name], &server); err != nil {
		logs.Warning("%s: %s", err.Error(), "missing or wrong '"+name+"' configuration, ignoring")
	}

	if server.Host == "" {
		server.Host = defaultHost
		logs.Warning("missing " + name + " 'host' configuration, assuming default value: '" + defaultHost + "'")
	}

	if server.Port == 0 {
		server.Port = defaultPort
		logs.Warning("missing "+name+" 'port' configuration, assuming default value: %d", defaultPort)
	}

	return server, nil
}

func (config Config) Server() (Server, error) {
	return config.CustomServer("server", "0.0.0.0", 8080)
}

func (s Server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
