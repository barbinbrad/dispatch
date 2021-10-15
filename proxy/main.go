package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	var (
		port = flags.Int("port", 4720, "port to listen on")
	)
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	addr := fmt.Sprintf("localhost:%d", *port)
	server, err := NewServer()
	if err != nil {
		return err
	}
	fmt.Printf("listening on: %s\n", addr)
	return http.ListenAndServe(addr, server.router)
}
