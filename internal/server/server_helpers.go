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
		return "", err
	}

	// validate and modify input if necessary
	err = validateInput(&input)
	if err != nil {
		return "", err
	}

	// prepend server responses
	response := fmt.Sprintf("[Echo Town]: %s", input)

	return response, nil
}

// validateInput modifies the input by checking for size and bad characters.
// It also trims the input of any whitespace.
func validateInput(input *string) error {
	// validate input larger than 1024 bytes
	if len(*input) > 1024 {
		return ERROR_LONG_MSG
	}

	// trim whitespace
	*input = strings.TrimSpace(*input)

	// validate input with non-printable characters (bad characters)
	for _, rune := range *input {
		if !strconv.IsPrint(rune) {
			return ERROR_NON_PRNT
		}
	}

	// validate UTF-8 string
	if !utf8.ValidString(*input) {
		return ERROR_BAD_UTF8
	}

	// append newline
	*input = *input + "\n"

	return nil
}
