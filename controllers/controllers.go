package controllers

import (
	"net/http"
	"text/template"

	"github.com/jmoiron/sqlx"
	"github.com/quorumsco/application"
	"github.com/quorumsco/router"
)

func getDB(r *http.Request) *sqlx.DB {
	return router.Context(r).Env["Application"].(*application.Application).Components["DB"].(*sqlx.DB)
}

func getTemplates(r *http.Request) map[string]*template.Template {
	return router.Context(r).Env["Application"].(*application.Application).Components["Templates"].(map[string]*template.Template)
}
