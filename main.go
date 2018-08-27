package main

import (
	"fmt"
	"goproxy/internal/base"
	"goproxy/internal/servecmd"
	"os"

	"github.com/namsral/flag"
)

func init() {
	base.Commands = []*base.Command{servecmd.CmdServe}
	// add the flags common to all commands
	for _, cmd := range base.Commands {
		base.AddBaseFlags(&cmd.Flag)
	}
}

func main() {
	flag.Usage = printUsage

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		printUsage()
		os.Exit(-1)
	}

	if args[0] == "help" {
		printUsage()
		os.Exit(0)
	}

	for _, cmd := range base.Commands {
		if cmd.Name() == args[0] && cmd.Runnable() {
			cmd.Flag.Usage = func() { cmd.Usage() }
			cmd.Flag.Parse(args[1:])

			args = cmd.Flag.Args()
			cmd.Run(cmd, args)
			return
		}
	}

	fmt.Printf("Unknown command %s\n", args[0])
	printUsage()
	os.Exit(-2)
}

func printUsage() {
	fmt.Print(`
Usage: goproxy [command] [arguments]

Available commands:

`)
	for _, cmd := range base.Commands {
		fmt.Printf("\t%-15s %s\n", cmd.Name(), cmd.Short)
	}
	fmt.Printf("\nYou can run goproxy [command] -h for info on each command.\n\n")
}
