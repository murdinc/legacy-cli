package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"text/template"
)

type errorMessage struct {
	Title   string
	Message string
	Blank   string
}

func PrintAnsi(templ string, data interface{}) {
	ansiTemplate := template.FuncMap{
		"ansi": AnsiCode,
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Funcs(ansiTemplate).Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

// AnsiCode outputs the ansi codes for changing terminal colors and behaviors.
func AnsiCode(code string) string {
	var ansiCode int
	switch code {
	case "reset":
		ansiCode = 0
	case "bright":
		ansiCode = 1
	case "dim":
		ansiCode = 2
	case "underscore":
		ansiCode = 4
	case "blink":
		ansiCode = 5
	case "reverse":
		ansiCode = 7
	case "hidden":
		ansiCode = 8
	case "fgblack":
		ansiCode = 30
	case "fgred":
		ansiCode = 31
	case "fggreen":
		ansiCode = 32
	case "fgyellow":
		ansiCode = 33
	case "fgblue":
		ansiCode = 34
	case "fgmagenta":
		ansiCode = 35
	case "fgcyan":
		ansiCode = 36
	case "fgwhite":
		ansiCode = 37
	case "bgblack":
		ansiCode = 40
	case "bgred":
		ansiCode = 41
	case "bggreen":
		ansiCode = 42
	case "bgyellow":
		ansiCode = 43
	case "bgblue":
		ansiCode = 44
	case "bgmagenta":
		ansiCode = 45
	case "bgcyan":
		ansiCode = 46
	case "bgwhite":
		ansiCode = 47
	}
	output := fmt.Sprintf("\033[%dm", ansiCode)
	return output
}

var ErrorMessageTemplate = `
{{ ansi "fgwhite"}}{{ ansi "bgred"}}{{.Blank}}{{ ansi ""}}
{{ ansi "fgwhite"}}{{ ansi "bgred"}}{{.Title}}{{ ansi ""}}
{{ ansi "fgwhite"}}{{ ansi "bgred"}}{{.Blank}}{{ ansi ""}}
{{ ansi "fgwhite"}}{{ ansi "bgred"}}{{.Message}}{{ ansi ""}}
{{ ansi "fgwhite"}}{{ ansi "bgred"}}{{.Blank}}{{ ansi ""}}

`

func ShowErrorMessage(title string, message string) {

	title = fmt.Sprintf("[%s]", title)

	// Figure out how wide of a notification we are going to be building
	titleWidth := len(title)
	msgWidth := len(message)
	var totalWidth int
	// TODO make this use math.Max or something? idk
	if titleWidth > msgWidth {
		totalWidth = titleWidth + 4
	} else {
		totalWidth = msgWidth + 4
	}

	totalWidth += (totalWidth % 2)

	blank := strings.Repeat(" ", totalWidth)

	// Pad our strings until they are centered at the same width
	title = addPadding(title, totalWidth)
	message = addPadding(message, totalWidth)

	// Finally print the output
	PrintAnsi(ErrorMessageTemplate, errorMessage{Title: title, Message: message, Blank: blank})
}

func addPadding(s string, w int) string {
	if len(s)%2 != 0 {
		s += " "
	}

	padding := strings.Repeat(" ", (w-len(s)%w)/2)
	t := []string{padding, s, padding}
	padded := strings.Join(t, "")
	return padded
}
