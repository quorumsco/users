package controllers

import (
	"net/http"
	"text/template"

	"github.com/iogo-framework/application"
	"github.com/iogo-framework/router"
	"github.com/jmoiron/sqlx"
)

func getDB(r *http.Request) *sqlx.DB {
	return router.Context(r).Env["Application"].(*application.Application).Components["DB"].(*sqlx.DB)
}

func getTemplates(r *http.Request) map[string]*template.Template {
	return router.Context(r).Env["Application"].(*application.Application).Components["Templates"].(map[string]*template.Template)
}
