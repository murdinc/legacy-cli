package cli

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

// Context is a type that is passed through to
// each Handler action in a cli application. Context
// can be used to retrieve context-specific Args and
// parsed command-line options.
type Context struct {
	App       *App
	flagSet   *flag.FlagSet
	globalSet *flag.FlagSet
	setFlags  map[string]bool
	SetArgs	  map[string]string
}

// Creates a new context. For use in when invoking an App or Command action.
func NewContext(app *App, set *flag.FlagSet, globalSet *flag.FlagSet, setArgs map[string]string) *Context {
	return &Context{app, set, globalSet, nil, setArgs}
}

// Looks up the value of a local int flag, returns 0 if no int flag exists
func (c *Context) Int(name string) int {
	return lookupInt(name, c.flagSet)
}

// Looks up the value of a local float64 flag, returns 0 if no float64 flag exists
func (c *Context) Float64(name string) float64 {
	return lookupFloat64(name, c.flagSet)
}

// Looks up the value of a local bool flag, returns false if no bool flag exists
func (c *Context) Bool(name string) bool {
	return lookupBool(name, c.flagSet)
}

// Looks up the value of a local string flag, returns "" if no string flag exists
func (c *Context) String(name string) string {
	return lookupString(name, c.flagSet)
}

// Looks up the value of a local string slice flag, returns nil if no string slice flag exists
func (c *Context) StringSlice(name string) []string {
	return lookupStringSlice(name, c.flagSet)
}

// Looks up the value of a local int slice flag, returns nil if no int slice flag exists
func (c *Context) IntSlice(name string) []int {
	return lookupIntSlice(name, c.flagSet)
}

// Looks up the value of a global int flag, returns 0 if no int flag exists
func (c *Context) GlobalInt(name string) int {
	return lookupInt(name, c.globalSet)
}

// Looks up the value of a global bool flag, returns false if no bool flag exists
func (c *Context) GlobalBool(name string) bool {
	return lookupBool(name, c.globalSet)
}

// Looks up the value of a global string flag, returns "" if no string flag exists
func (c *Context) GlobalString(name string) string {
	return lookupString(name, c.globalSet)
}

// Looks up the value of a global string slice flag, returns nil if no string slice flag exists
func (c *Context) GlobalStringSlice(name string) []string {
	return lookupStringSlice(name, c.globalSet)
}

// Looks up the value of a global int slice flag, returns nil if no int slice flag exists
func (c *Context) GlobalIntSlice(name string) []int {
	return lookupIntSlice(name, c.globalSet)
}

// Determines if the flag was actually set exists
func (c *Context) IsSet(name string) bool {
	if c.setFlags == nil {
		c.setFlags = make(map[string]bool)
		c.flagSet.Visit(func(f *flag.Flag) {
			c.setFlags[f.Name] = true
		})
	}
	return c.setFlags[name] == true
}

// Return an Argument by key
func (c *Context) NamedArg(k string) string {
	arg := c.SetArgs[k]
	return arg
}

type Args []string

// Returns the command line arguments associated with the context.
func (c *Context) Args() Args {
	args := Args(c.flagSet.Args())
	return args
}

// Returns the nth argument, or else a blank string
func (a Args) Get(n int) string {
	if len(a) > n {
		return a[n]
	}
	return ""
}

// Returns the first argument, or else a blank string
func (a Args) First() string {
	return a.Get(0)
}

// Return the rest of the arguments (not the first one)
// or else an empty string slice
func (a Args) Tail() []string {
	if len(a) >= 2 {
		return []string(a)[1:]
	}
	return []string{}
}

// Checks if there are any arguments present
func (a Args) Present() bool {
	return len(a) != 0
}


func lookupInt(name string, set *flag.FlagSet) int {
	f := set.Lookup(name)
	if f != nil {
		val, err := strconv.Atoi(f.Value.String())
		if err != nil {
			return 0
		}
		return val
	}

	return 0
}

func lookupFloat64(name string, set *flag.FlagSet) float64 {
	f := set.Lookup(name)
	if f != nil {
		val, err := strconv.ParseFloat(f.Value.String(), 64)
		if err != nil {
			return 0
		}
		return val
	}

	return 0
}

func lookupString(name string, set *flag.FlagSet) string {
	f := set.Lookup(name)
	if f != nil {
		return f.Value.String()
	}

	return ""
}

func lookupStringSlice(name string, set *flag.FlagSet) []string {
	f := set.Lookup(name)
	if f != nil {
		return (f.Value.(*StringSlice)).Value()

	}

	return nil
}

func lookupIntSlice(name string, set *flag.FlagSet) []int {
	f := set.Lookup(name)
	if f != nil {
		return (f.Value.(*IntSlice)).Value()

	}

	return nil
}

func lookupBool(name string, set *flag.FlagSet) bool {
	f := set.Lookup(name)
	if f != nil {
		val, err := strconv.ParseBool(f.Value.String())
		if err != nil {
			return false
		}
		return val
	}

	return false
}

func normalizeFlags(flags []Flag, set *flag.FlagSet) error {
	visited := make(map[string]bool)
	set.Visit(func(f *flag.Flag) {
		visited[f.Name] = true
	})
	for _, f := range flags {
		parts := strings.Split(f.getName(), ",")
		if len(parts) == 1 {
			continue
		}
		var ff *flag.Flag
		for _, name := range parts {
			name = strings.Trim(name, " ")
			if visited[name] {
				if ff != nil {
					return errors.New("Cannot use two forms of the same flag: " + name + " " + ff.Name)
				}
				ff = set.Lookup(name)
			}
		}
		if ff == nil {
			continue
		}
		for _, name := range parts {
			name = strings.Trim(name, " ")
			set.Set(name, ff.Value.String())
		}
	}
	return nil
}
