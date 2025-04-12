package server

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"unicode/utf8"
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
// encountered, performs validation, and clears the buffer at the end.
func createResponse(reader *bufio.Reader) (string, error) {
	defer reader.Discard(reader.Buffered()) // clear buffer at the end

	// read until a newline character
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from client:", err)
		return "", err
	}

	validateInput(&input)

	// prepend server responses
	response := fmt.Sprintf("[Echo Town]: %s", input)

	return response, nil
}

// validateInput modifies the input by checking for size and bad characters.
// It also trims the input of any whitespace. Invalid input is processed
// as messages written back to the client.
func validateInput(input *string) {
	// validate input larger than 1024 bytes
	if len(*input) > 1024 {
		*input = "[ERROR] Message cannot be longer than 1024 bytes!"
	}

	// trim whitespace
	*input = strings.TrimSpace(*input)

	// validate input with non-printable characters (bad characters)
	for _, rune := range *input {
		if !strconv.IsPrint(rune) {
			*input = "[ERROR] Message contains non-printable characters!"
			break
		}
	}

	// validate UTF-8 string
	if !utf8.ValidString(*input) {
		*input = "[ERROR] Message contains invalid UTF-8 characters!"
	}
}
