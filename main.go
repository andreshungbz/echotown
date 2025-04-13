package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/andreshungbz/echotown/internal/server"
)

func main() {
	port := flag.Int("port", 4000, "Designated port to start the server on.")
	flag.Parse()

	if *port < 1 || *port > 65535 {
		fmt.Fprintf(os.Stderr, "[ERROR] Port %d is invalid. Port must be between 1 and 65535.\n", *port)
		os.Exit(1)
	}

	server.Start(*port)
}
