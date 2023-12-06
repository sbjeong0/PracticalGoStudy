package main

import (
	"errors"
	"fmt"
	"github.com/sbjeong0/sub-cmd-arch/cmd"
	"io"
	"os"
)

var errInvalidSubCommand = errors.New("Invalid sub-command specified")

func main() {
	err := handleCommand(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}

func handleCommand(w io.Writer, args []string) error {
	var err error
	if len(os.Args) < 1 {
		err = errInvalidSubCommand
	}
	switch os.Args[1] {
	case "http":
		err = cmd.HandleHttp(os.Stdout, os.Args[2:])
	case "grpc":
		err = cmd.HandleGrpc(os.Stdout, os.Args[2:])
	case "-h":
		printUsage(os.Stdout)
	case "--help":
		printUsage(os.Stdout)
	default:
		err = errInvalidSubCommand
	}

	if errors.Is(err, cmd.ErrNoServerSpecified) || errors.Is(err, errInvalidSubCommand) {
		fmt.Fprintln(w, err)
		printUsage(w)
	}
	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, "Usage: [http|grpc] -h\n")
	cmd.HandleGrpc(w, []string{"-h"})
	cmd.HandleHttp(w, []string{"-h"})
}
