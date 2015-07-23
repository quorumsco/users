package settings

import (
	"encoding/json"

	"github.com/iogo-framework/logs"
)

type Sqlite3 struct {
	Path string
}

func (config Config) Sqlite3() (Sqlite3, error) {
	var sqlite Sqlite3

	if err := json.Unmarshal(config.Components["sqlite3"], &sqlite); err != nil {
		logs.Warning("%s: %s", err.Error(), "missing or wrong 'sqlite3' configuration, ignoring")
	}

	if sqlite.Path == "" {
		logs.Warning("missing sqlite3 'path' configuration, assuming default value: 'db.sqlite'")
		sqlite.Path = "db.sqlite"
	}

	return sqlite, nil
}

func (s Sqlite3) String() string {
	return s.Path
}
