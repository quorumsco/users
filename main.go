package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/iogo-framework/application"
	"github.com/iogo-framework/cmd"
	"github.com/iogo-framework/databases"
	"github.com/iogo-framework/logs"
	"github.com/iogo-framework/router"
	"github.com/iogo-framework/settings"
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
		cli.BoolFlag{Name: "debug, d", Usage: "debug log level"},
		cli.HelpFlag,
	}...)
	cmd.RunAndExitOnError()
}

func serve(ctx *cli.Context) error {
	var app *application.Application
	var err error

	config, err := settings.Parse(ctx.String("config"))
	if err != nil && ctx.String("config") != "" {
		logs.Error(err)
	}

	if ctx.Bool("debug") || config.Debug() {
		logs.Level(logs.DebugLevel)
	}

	dialect, args, err := config.SqlDB()
	if err != nil {
		logs.Critical(err)
		os.Exit(1)
	}
	logs.Debug("database type: %s", dialect)

	app = application.New()
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
	app.Components["Mux"] = router.New()

	if ctx.Bool("debug") || config.Debug() {
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
	db, err := databases.InitGORM(dialect, args)
	if err != nil {
		return err
	}

	db.LogMode(true)

	db.AutoMigrate(models.Models()...)

	return nil
}
