package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Command is a subcommand for a cli.App.
type Command struct {
	// The name of the command
	Name string
	// short name of the command. Typically one character
	ShortName string
	// A short example of a formed command
	Example string
	// A description of what this command does
	Description string
	// The function to call when this command is invoked
	Action func(context *Context)
	// The list of required and optional arguments
	Arguments       []Argument
	// List of flags to parse
	Flags []Flag
}

// Argument can be required or optional arguments for a cli.App commands.
type Argument struct {
	// The name of the command
	Name string
	// short name of the command. Typically one character
	Usage string
	// A longer explaination of how the command works
	Description string
	// Required defines if this argument is required for the command or not
	Optional bool
}

// Invokes the command given the context, parses ctx.Args() to generate command-specific flags
func (c Command) Run(ctx *Context) error {
	// append help to flags
	c.Flags = append(
		c.Flags,
		BoolFlag{"help, h", "show help"},
	)

	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)

	firstFlagIndex := -1
	for index, arg := range ctx.Args() {
		if strings.HasPrefix(arg, "-") {
			firstFlagIndex = index
			break
		}
	}

	var err error
	if firstFlagIndex > -1 {
		args := ctx.Args()
		regularArgs := args[1:firstFlagIndex]
		flagArgs := args[firstFlagIndex:]
		c.BuildCustomArgs(ctx)
		err = set.Parse(append(flagArgs, regularArgs...))
	} else {
		c.BuildCustomArgs(ctx)
		err = set.Parse(ctx.Args().Tail())
	}

	if err != nil {
		fmt.Printf("Incorrect Usage.\n\n")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return err
	}

	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		fmt.Println("")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return nerr
	}
	context := NewContext(ctx.App, set, ctx.globalSet, ctx.SetArgs)
	if checkCommandHelp(context, c.Name) {
		return nil
	}

	c.Action(context)
	return nil
}

func (c Command) BuildCustomArgs(ctx *Context) error {
	// Combine the argument names with their entered values
	// The order in the command input are known, their order in the context is unknown, so we do some crazy stuff to get them sorted out
	i := 0
	m := make(map[string]string)

	for _,val := range ctx.Args()[1:]  {
		if !strings.HasPrefix(val, "-") {
			name := c.Arguments[i].Name
			value := val
			m[name] = value
			i++
		}
	}

	ctx.SetArgs = m
	return nil
}


// Returns true if Command.Name or Command.ShortName matches given name
func (c Command) HasName(name string) bool {
	return c.Name == name || c.ShortName == name
}
