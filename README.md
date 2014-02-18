
# cli.go
cli.go is simple, fast, and fun package for building command line apps in Go. The goal is to enable developers to write fast and distributable command line applications in an expressive way.

The majority of the code is the [cli library written by codegansta] (https://github.com/codegangsta/cli)
I have also included the [termtable library written by stevedomin] (https://github.com/stevedomin/termtable)

This library is *NOT* backwards compatible with the version written by codegangsta, as I have changed some terminologies and types to accomodate named arguments. This was done to allow for mandatory and optional arguments to commands. It should not be too time consuming to adapt an existing project to use this library though, all of the changes can be made in the place where NewApp() is defined. 

I also added ansi color  support via text/template and a function to display errors. 


## Overview
Command line apps are usually so tiny that there is absolutely no reason why your code should *not* be self-documenting. Things like generating help text and parsing command flags/options should not hinder productivity when writing a command line app.

This is where cli.go comes into play. cli.go makes command line programming fun, organized, and expressive!

## Installation
Make sure you have the a working Go environment (go 1.1 is *required*). [See the install instructions](http://golang.org/doc/install.html).

To install cli.go, simply run:
```
$ go get github.com/murdinc/cli
```

Make sure your PATH includes to the `$GOPATH/bin` directory so your commands can be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

## Getting Started
One of the philosophies behind cli.go is that an API should be playful and full of discovery. So a cli.go app can be as little as one line of code in `main()`. 

``` go
package main

import (
  "os"
  "github.com/murdinc/cli"
)

func main() {
  cli.NewApp().Run(os.Args)
}
```

This app will run and show help text, but is not very useful. Let's give an action to execute and some help documentation:

``` go
package main

import (
  "os"
  "github.com/murdinc/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "boom"
  app.Usage = "make an explosive entrance"
  app.Action = func(c *cli.Context) {
    println("boom! I say!")
  }
  
  app.Run(os.Args)
}
```

Running this already gives you a ton of functionality, plus support for things like subcommands and flags, which are covered below.

## Example

Being a programmer can be a lonely job. Thankfully by the power of automation that is not the case! Let's create a greeter app to fend off our demons of loneliness!

``` go
/* greet.go */
package main

import (
  "os"
  "github.com/murdinc/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "greet"
  app.Usage = "fight the loneliness!"
  app.Action = func(c *cli.Context) {
    println("Hello friend!")
  }
  
  app.Run(os.Args)
}
```

Install our command to the `$GOPATH/bin` directory:

```
$ go install
```

Finally run our new command:

```
$ greet
Hello friend!
```

cli.go also generates some bitchass help text:
```
$ greet help
NAME:
    greet - fight the loneliness!

USAGE:
    greet [global options] command [command options] [arguments...]

VERSION:
    0.0.0

COMMANDS:
    help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS
    --version	Shows version information
```

### Arguments
You can lookup arguments by calling the `Args` function on cli.Context.

``` go
...
app.Action = func(c *cli.Context) {
  println("Hello", c.Args()[0])
}
...
```

### Flags
Setting and querying flags is simple.
``` go
...
app.Flags = []cli.Flag {
  cli.StringFlag{"lang", "english", "language for the greeting"},
}
app.Action = func(c *cli.Context) {
  name := "someone"
  if len(c.Args()) > 0 {
    name = c.Args()[0]
  }
  if c.String("lang") == "spanish" {
    println("Hola", name)
  } else {
    println("Hello", name)
  }
}
...
```

#### Alternate Names

You can set alternate (or short) names for flags by providing a comma-delimited list for the Name. e.g.

``` go
app.Flags = []cli.Flag {
  cli.StringFlag{"lang, l", "english", "language for the greeting"},
}
```

That flag can then be set with `--lang spanish` or `-l spanish`. Note that giving two different forms of the same flag in the same command invocation is an error.

### Subcommands

Subcommands can be defined for a more git-like command line app.
```go
...
app.Commands = []cli.Command{
  {
    Name:      "add",
    ShortName: "a",
    Usage:     "add a task to the list",
    Action: func(c *cli.Context) {
      println("added task: ", c.FirstArg())
    },
  },
  {
    Name:      "complete",
    ShortName: "c",
    Usage:     "complete a task on the list",
    Action: func(c *cli.Context) {
      println("completed task: ", c.FirstArg())
    },
  },
}
...
```

## About
cli.go is written by none other than the [Code Gangsta](http://codegangsta.io) - This version of the library was modified by murdinc to add additional functionality
