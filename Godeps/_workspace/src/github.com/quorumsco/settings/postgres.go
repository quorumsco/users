package settings

import (
	"encoding/json"
	"fmt"

	"github.com/iogo-framework/logs"
)

type Postgres struct {
	DB       string
	Host     string
	Port     int
	User     string
	Password string
}

func (config Config) Postgres() (Postgres, error) {
	var postgres Postgres

	err := json.Unmarshal(config.Components["postgres"], &postgres)
	if err != nil {
		logs.Warning("%s: %s", err.Error(), "missing or wrong 'postgres' configuration, ignoring")
	}

	if postgres.DB == "" {
		postgres.DB = "postgres"
		logs.Warning("missing postgres 'db' configuration, assuming default value: 'postgres'")
	}
	if postgres.Host == "" {
		postgres.Host = "postgres"
		logs.Warning("missing postgres 'host' configuration, assuming default value: 'postgres'")
	}
	if postgres.Port == 0 {
		postgres.Port = 5432
		logs.Warning("missing postgres 'port' configuration, assuming default value: 5432")
	}
	if postgres.User == "" {
		postgres.User = "postgres"
		logs.Warning("missing postgres 'user' configuration, assuming default value: 'postgres'")
	}

	if postgres.Password == "" {
		logs.Warning("missing postgres 'password' configuration, assuming no password is needed")
	}

	return postgres, nil
}

func (p Postgres) String() string {
	if p.Password == "" {
		return fmt.Sprintf("postgres://%s@%s:%d/%s?sslmode=disable",
			p.User, p.Host, p.Port, p.DB)
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DB)
}
