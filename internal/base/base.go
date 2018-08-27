package base

import (
	"fmt"
	"log"
	"os"
	"strings"

	"flag"
)

var (
	// Dir is the directory where local package files are kept
	Dir string

	// verbose indicates that we must print
	verbose bool
)

// A Command is an implementation of a go command
// like go build or go fix.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)
	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string
	// Short is the short description shown in the 'go help' output.
	Short string
	// Long is the long message shown in the 'go help <this-command>' output.
	Long string
	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
	// OmitBaseFlags indicates that the command will not support the
	// base flags.
	OmitBaseFlags bool
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'go help'.
var Commands []*Command

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	if !c.OmitBaseFlags {
		fmt.Fprint(os.Stderr, `
goproxy also supports the following flags for all its commands:

    -dir       directory containing package data

`)
	}
	os.Exit(2)
}

// Runnable reports whether the command can be run; otherwise
// it is a documentation pseudo-command such as importpath.
func (c *Command) Runnable() bool {
	return c.Run != nil
}

// AddBaseFlags adds base flags that are common to all commands
func AddBaseFlags(flags *flag.FlagSet) {
	flags.StringVar(&Dir, "dir", "./repo", "Directory containing package data")
	flags.BoolVar(&verbose, "v", false, "print detailed output")
}

// Log will print a message to stdout if and only if verbose is enabled
func Log(msg ...interface{}) {
	if verbose {
		log.Println(msg...)
	}
}

// Logf will print a message to stdout if and only if verbose is enabled
func Logf(format string, a ...interface{}) {
	if verbose {
		msg := format
		if msg == "" && len(a) > 0 {
			msg = fmt.Sprint(a...)
		} else if msg != "" && len(a) > 0 {
			msg = fmt.Sprintf(format, a...)
		}
		log.Println(msg)
	}
}
