package models

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func Models() []interface{} {
	return []interface{}{
		&User{},
	}
}

func Stores() []interface{} {
	return []interface{}{
		UserSQL{},
	}
}
