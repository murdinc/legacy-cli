package cli

import (
	"fmt"
)

// The text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var AppHelpTemplate = `{{ printf "\033[%qm%v\033[0m" "33" .Name }} - {{ printf "\033[%qm%v\033[0m" "33" .Version }}

{{ ansi "fgyellow"}}Usage:{{ ansi ""}}
   {{.Name}} [global options] command [command options] [arguments...]

{{ ansi "fgyellow"}}Commands:{{ ansi ""}}
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Description}}
   {{end}}
{{ ansi "fgyellow"}}Global Options:{{ ansi ""}}
   {{range .Flags}}{{.}}
   {{end}}
`

// The text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = `{{ ansi "fgyellow"}}Usage:{{ ansi ""}}
   {{.Name}} {{range .Arguments}}{{if .Optional}}[{{.Name}}] {{else}}{{.Name}} {{end}}{{end}}[--flags]

{{ ansi "fgyellow"}}Arguments:{{ ansi ""}}{{range .Arguments}}
   {{if .Optional}}[{{.Name}}]{{else}}{{.Name}}{{end}}{{ "\t" }}{{.Description}}{{end}}

{{ ansi "fgyellow"}}Flags:{{ ansi ""}}
   {{range .Flags}}{{.}}{{end}}

{{ ansi "fgyellow"}}Example:{{ ansi ""}}
   {{.Example}}
`

var helpCommand = Command{
	Name:      "help",
	ShortName: "h",
	Description:     "Shows a list of commands or help for one command",
	Action: func(c *Context) {
		args := c.Args()
		if args.Present() {
			ShowCommandHelp(c, args.First())
		} else {
			ShowAppHelp(c)
		}
	},
}

// Prints help for the App
func ShowAppHelp(c *Context) {
	printAnsi(AppHelpTemplate, c.App)
}

// Prints an error message
func PrintError(c *Context) {
        printAnsi(AppHelpTemplate, c.App)
}


// Prints help for the given command
func ShowCommandHelp(c *Context, command string) {
	for _, c := range c.App.Commands {
		if c.HasName(command) {
			printAnsi(CommandHelpTemplate, c)
			return
		}
	}
	fmt.Printf("No help topic for '%v'\n", command)
}

// Prints the version number of the App
func ShowVersion(c *Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
}

func checkVersion(c *Context) bool {
	if c.GlobalBool("version") {
		ShowVersion(c)
		return true
	}

	return false
}

func checkHelp(c *Context) bool {
	if c.GlobalBool("h") || c.GlobalBool("help") {
		ShowAppHelp(c)
		return true
	}

	return false
}

func checkCommandHelp(c *Context, name string) bool {
	if c.Bool("h") || c.Bool("help") {
		ShowCommandHelp(c, name)
		return true
	}

	return false
}
