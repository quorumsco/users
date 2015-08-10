package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/quorumsco/application"
	"github.com/quorumsco/cmd"
	"github.com/quorumsco/databases"
	"github.com/quorumsco/gojimux"
	"github.com/quorumsco/logs"
	"github.com/quorumsco/router"
	"github.com/quorumsco/settings"
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
		cli.StringFlag{Name: "config, c", Usage: "configuration file", EnvVar: "CONFIG"},
		cli.HelpFlag,
	}...)
	cmd.RunAndExitOnError()
}

func serve(ctx *cli.Context) error {
	var err error

	var config settings.Config
	if ctx.String("config") != "" {
		config, err = settings.Parse(ctx.String("config"))
		if err != nil {
			logs.Error(err)
		}
	}

	if config.Debug() {
		logs.Level(logs.DebugLevel)
	}

	dialect, args, err := config.SqlDB()
	if err != nil {
		logs.Critical(err)
		os.Exit(1)
	}
	logs.Debug("database type: %s", dialect)

	var app = application.New()
	if app.Components["DB"], err = databases.InitSQLX(dialect, args); err != nil {
		logs.Critical(err)
		os.Exit(1)
	}
	logs.Debug("connected to %s", args)

	if config.Migrate() {
		if err := migrate(dialect, args); err != nil {
			logs.Critical(err)
			os.Exit(1)
		}
		logs.Debug("database migrated successfully")
	}

	app.Components["Templates"] = views.Templates()
	// app.Components["Mux"] = router.New()
	app.Components["Mux"] = gojimux.New()

	if config.Debug() {
		app.Use(router.Logger)
	}

	app.Use(app.Apply)

	app.Get("/users/register", controllers.Register)
	app.Post("/users/register", controllers.Register)

	server, err := config.Server()
	if err != nil {
		logs.Critical(err)
		os.Exit(1)
	}
	return app.Serve(server.String())
}

func migrate(dialect string, args string) error {
	var db, err = databases.InitGORM(dialect, args)
	if err != nil {
		return err
	}

	db.LogMode(true)

	db.AutoMigrate(models.Models()...)

	return nil
}
