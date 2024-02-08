package main

import (
	"PracticalGoStudy/practice/cmd"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	err := handleCommand(os.Stdout, os.Args[1:])
	if err != nil {
		os.Exit(1)
	}
}

func handleCommand(w io.Writer, args []string) error {
	var err error
	if len(os.Args) < 2 {
		err = cmd.ErrInvalidSubCommand
	}
	if err == nil {
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
			err = cmd.ErrInvalidSubCommand
		}
	}

	if errors.Is(err, cmd.ErrNoServerSpecified) || errors.Is(err, cmd.ErrInvalidSubCommand) {
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
