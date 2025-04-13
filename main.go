// Echo Town is a TCP echo server program that handles multiple clients, writes logs,
// has custom personality responses, and has a simple command protocol.
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

	if !validatePort(*port) {
		fmt.Fprintf(os.Stderr, "[ERROR] Port %d is invalid. Port must be between 1 and 65535.\n", *port)
		os.Exit(1)
	}

	server.Start(*port)
}

// validatePort returns true if passed in port is valid and false otherwise.
func validatePort(port int) bool {
	return port >= 1 && port <= 65535
}
