package server

import (
	"bufio"
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

// createResponse reads the client input from the passed [bufio.Reader] and constructs
// the server response string to write back to the client. It reads until a newline is
// encountered, performs validation, and clears the buffer at the end. Invalid input
// is processed as messages written back to the client.
func createResponse(reader *bufio.Reader) (string, error) {
	defer reader.Discard(reader.Buffered()) // clear buffer at the end

	// read until a newline character
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return "", err
	}

	// override input with a reject message if input is larger than 1024 bytes
	if len(input) > 1024 {
		input = "[ERROR] Message cannot be longer than 1024 bytes!"
	}

	// prepend server responses
	response := fmt.Sprintf("[Echo Town]: %s", input)

	return response, nil
}
