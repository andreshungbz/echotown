package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// monitorTermSig creates a channel to monitor OS termination signals
// and launches a single goroutine to print and log a server end message when
// the signal is received.
func monitorTermSig(logger *log.Logger) {
	signalChan := make(chan os.Signal, 1)

	// register channel to receive interrupt and termination OS signals
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// when a proper signal is received, log message and exit program
	go func() {
		<-signalChan
		fmt.Println() // print a new line just for the console
		logger.Print("[INFO] Echo Town server stopped.\n\n")
		os.Exit(0)
	}()
}
