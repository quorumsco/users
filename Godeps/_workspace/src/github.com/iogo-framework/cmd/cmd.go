/*
Package cmd is a simple wrapper around codegangsta.cli to provide a
basic command line interface with no subcommands.
*/
package cmd

import "github.com/codegangsta/cli"

const helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]
VERSION:
   {{.Version}}{{if len .Authors}}
AUTHOR(S): 
   {{range .Authors}}{{ . }}{{end}}{{end}}
{{if .Commands}}
COMMANDS:
   {{range .Commands}}{{join .Names ", "}}{{ "\t" }}{{.Usage}}
   {{end}}{{end}}{{if .Flags}}OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

// New creates a new cmd.
func New() *cli.App {
	app := cli.NewApp()
	app.Action = func(ctx *cli.Context) {}
	app.HideHelp = true
	cli.AppHelpTemplate = helpTemplate

	return app
}
