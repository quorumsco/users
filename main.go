package main

import (
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/iogo-framework/application"
	"github.com/iogo-framework/cmd"
	"github.com/iogo-framework/logs"
	"github.com/iogo-framework/router"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/quorumsco/users/controllers"
	"github.com/quorumsco/users/models"
	"github.com/quorumsco/users/views"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	cmd := cmd.New()
	cmd.Name = "users"
	cmd.Usage = "quorums users backend"
	cmd.Version = "0.0.1"
	cmd.Before = serve
	cmd.Flags = append(cmd.Flags, []cli.Flag{
		cli.BoolFlag{Name: "migrate, m", Usage: "migrate the database"},
		cli.StringFlag{Name: "listen, l", Value: "0.0.0.0:8080", Usage: "server listening host:port"},
		cli.BoolFlag{Name: "debug, d", Usage: "print debug information"},
		cli.HelpFlag,
	}...)
	cmd.RunAndExitOnError()
}

func serve(ctx *cli.Context) error {
	var app *application.Application
	var err error

	if ctx.Bool("migrate") {
		migrate()
		logs.Debug("Database migrated")
	}

	app = application.New()
	if app.Components["DB"], err = sqlx.Connect("sqlite3", "/tmp/contacts.db"); err != nil {
		return err
	}
	app.Components["Templates"] = views.Templates()
	app.Components["Mux"] = router.New()

	app.Use(router.Logger)
	app.Use(app.Apply)

	app.Get("/users/register", controllers.Register)
	app.Post("/users/register", controllers.Register)

	app.Serve(ctx.String("listen"))

	return nil
}

func migrate() {
	db, err := gorm.Open("sqlite3", "/tmp/contacts.db")
	if err != nil {
		logs.Error(err)
		return
	}

	err = db.DB().Ping()
	if err != nil {
		logs.Error(err)
		return
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(false)

	db.AutoMigrate(models.Models()...)
}
