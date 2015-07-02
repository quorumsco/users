package main

import (
	"net/http"
	"runtime"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/iogo-framework/application"
	"github.com/iogo-framework/cmd"
	"github.com/iogo-framework/logs"
	"github.com/iogo-framework/router"
	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"
	"github.com/quorumsco/contacts/models"
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
		cli.StringFlag{Name: "cpu, cpuprofile", Usage: "cpu profiling"},
		cli.BoolFlag{Name: "m, migrate", Usage: "migrate the database"},
		cli.StringFlag{Name: "listen, l", Value: "localhost:8080", Usage: "server listening port"},
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

	if app, err = application.New(); err != nil {
		return err
	}

	app.Components["DB"], err = initSQLX()
	if err != nil {
		return err
	}

	app.Components["Templates"] = make(map[string]*template.Template)

	app.Mux = router.New()

	app.Use(router.Logger)
	app.Use(app.Apply)
	app.Use(cors)

	app.Serve(ctx.String("listen"))

	return nil
}

func cors(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "access-control-allow-origin,content-type")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
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

func initSQLX() (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	if db, err = sqlx.Connect("sqlite3", "/tmp/contacts.db"); err != nil {
		return nil, err
	}

	db.Ping()
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	return db, nil
}
