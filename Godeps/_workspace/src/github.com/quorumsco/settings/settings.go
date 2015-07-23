package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/iogo-framework/logs"
)

type Config struct {
	Components map[string]json.RawMessage
	Settings   map[string]interface{}
}

func (config Config) Debug() bool {
	debug, ok := config.Settings["debug"].(bool)
	if !ok {
		return false
	}
	return debug
}

func (config Config) Migrate() bool {
	debug, ok := config.Settings["migrate"].(bool)
	if !ok {
		return false
	}
	return debug
}

func (config Config) Dialect() (string, error) {
	client, ok := config.Settings["database"].(string)
	if ok && client != "" {
		return client, nil
	}
	logs.Warning("missing 'database' configuration, assuming default value: 'sqlite3'")
	return "sqlite3", nil
}

func (config Config) SqlDB() (dialect, args string, err error) {
	dialect, err = config.Dialect()
	if err != nil {
		return
	}

	switch dialect {
	case "postgres":
		var postgres Postgres
		postgres, err = config.Postgres()
		if err != nil {
			return
		}
		args = postgres.String()
	case "sqlite3":
		var sqlite Sqlite3
		sqlite, err = config.Sqlite3()
		if err != nil {
			return
		}
		args = sqlite.String()
	default:
		err = fmt.Errorf("unknown sql dialect '%s'", dialect)
	}
	return
}

func Parse(file string) (Config, error) {
	var config Config

	f, err := os.Open(file)
	if err != nil {
		return config, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(buf, &config); err != nil {
		return config, err
	}

	return config, nil
}
