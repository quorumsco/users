package controllers

import (
	"net/http"

	"github.com/iogo-framework/application"
	"github.com/iogo-framework/router"
	"github.com/jmoiron/sqlx"
)

func getDB(r *http.Request) *sqlx.DB {
	return router.GetContext(r).Env["Application"].(*application.Application).Components["DB"].(*sqlx.DB)
}
