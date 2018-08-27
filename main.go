package main

import (
	"fmt"
	"goproxy/internal/commandline/base"
	"goproxy/internal/commandline/servecmd"
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

	/*	dataDir := flag.String("dir", "repo", "Directory containing the files blah blah")
		httpAddr := flag.String("http", ":8080", "address to bind the HTTP server to")
		tls := flag.Bool("tls", false, "Whether to use TLS")
		tlsCert := flag.String("cert", "", "TLS certificate file")
		tlsKey := flag.String("key", "", "TLS key file")
		flag.Parse()
		panic(proxy.Run(*dataDir, *httpAddr, proxy.TLS(*tls, *tlsCert, *tlsKey)))*/
}

func printUsage() {
	fmt.Print(`
Usage: goproxy [command] [arguments]

Available commands:

`)
	for _, cmd := range base.Commands {
		fmt.Printf("\t%-15s %s\n", cmd.Name(), cmd.Short)
	}
	fmt.Printf("You can run goproxy [command] -h for info on each command.\n\n")
}
