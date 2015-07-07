package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
	"github.com/iogo-framework/application"
	"github.com/iogo-framework/cmd"
	"github.com/iogo-framework/databases"
	"github.com/iogo-framework/logs"
	"github.com/iogo-framework/router"
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
		cli.StringFlag{Name: "listen-host", Value: "0.0.0.0", Usage: "server listening host", EnvVar: "LISTEN_HOST"},
		cli.IntFlag{Name: "listen-port", Value: 8080, Usage: "server listening port", EnvVar: "LISTEN_PORT"},

		cli.StringFlag{Name: "sql-dialect", Value: "sqlite3", Usage: "database dialect ('sqlite3' or 'postgres')", EnvVar: "SQL_DIALECT"},

		cli.StringFlag{Name: "postgres-host", Value: "postgres", Usage: "postgresql host", EnvVar: "POSTGRES_HOST"},
		cli.IntFlag{Name: "postgres-port", Value: 5432, Usage: "postgresql port", EnvVar: "POSTGRES_PORT"},
		cli.StringFlag{Name: "postgres-user", Value: "postgres", Usage: "postgresql user", EnvVar: "POSTGRES_USER"},
		cli.StringFlag{Name: "postgres-password", Value: "postgres", Usage: "postgresql password", EnvVar: "POSTGRES_PASSWORD"},
		cli.StringFlag{Name: "postgres-db", Value: "postgres", Usage: "postgresql database", EnvVar: "POSTGRES_DB"},

		cli.StringFlag{Name: "sqlite-path", Value: "/tmp/db.sqlite", Usage: "sqlite path", EnvVar: "SQLITE_PATH"},

		cli.BoolFlag{Name: "migrate, m", Usage: "migrate the database", EnvVar: "MIGRATE"},
		cli.BoolFlag{Name: "debug, d", Usage: "print debug information", EnvVar: "DEBUG"},
		cli.HelpFlag,
	}...)
	cmd.RunAndExitOnError()
}

func serve(ctx *cli.Context) error {
	var app *application.Application
	var err error

	if ctx.Bool("debug") {
		logs.Level(logs.DebugLevel)
	}

	var dialect, args string

	switch ctx.String("sql-dialect") {
	case "postgres":
		dialect = "postgres"
		args = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			ctx.String("postgres-user"),
			ctx.String("postgres-password"),
			ctx.String("postgres-host"),
			ctx.Int("postgres-port"),
			ctx.String("postgres-db"),
		)
		logs.Debug("Loading database %s at %s", dialect, ctx.String("postgres-host"))
	case "sqlite3":
		fallthrough
	default:
		dialect = "sqlite3"
		args = ctx.String("sqlite-path")
		logs.Debug("Loading database %s at %s", dialect, args)
	}

	if ctx.Bool("migrate") {
		if err := migrate(dialect, args); err != nil {
			logs.Critical(err)
			os.Exit(1)
		}
		logs.Debug("Database migrated successfully")
	}

	app = application.New()
	if app.Components["DB"], err = databases.InitSQLX(dialect, args); err != nil {
		logs.Critical(err)
		os.Exit(1)
	}
	app.Components["Templates"] = views.Templates()
	app.Components["Mux"] = router.New()

	app.Use(router.Logger)
	app.Use(app.Apply)

	app.Get("/users/register", controllers.Register)
	app.Post("/users/register", controllers.Register)

	return app.Serve(fmt.Sprintf("%s:%d", ctx.String("listen-host"), ctx.Int("listen-port")))
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
