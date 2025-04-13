// Package command contains custom server responses to specific command protocol strings
package command

import (
	"strings"
	"time"
)

type commandResponse struct {
	// message function that returns a string to be printed to client
	Message func(string) string

	// determines whether to close the connection
	Close bool

	// what the command does
	Description string
}

// extendible map of specific command protocols to custom response
var commands = make(map[string]commandResponse)

// initiate map
func init() {
	commands["/time"] = commandResponse{
		Message: func(string) string {
			return time.Now().String()
		},
		Close:       false,
		Description: "Displays the server time.",
	}

	commands["/quit"] = commandResponse{
		Message: func(string) string {
			return ""
		},
		Close:       true,
		Description: "Closes the connection to the server.",
	}

	commands["/echo"] = commandResponse{
		Message: func(input string) string {
			if len(input) > len("/echo") {
				return input[len("/echo "):]
			}
			return ""
		},
		Close:       false,
		Description: "Returns only the message. Same result if not used.",
	}

	commands["/help"] = commandResponse{
		Message: func(string) string {
			var helpMessage strings.Builder
			for key, response := range commands {
				helpMessage.WriteString("\n" + key + " - " + response.Description)
			}
			return helpMessage.String()
		},
		Close:       false,
		Description: "Lists all available commands and their descriptions.",
	}
}

// Parse compares the input to pre-determined keys, and if they match,
// the custom server response is returned. Otherwise return the same input.
// The returned boolean determines whether the connection should be closed.
func Parse(input string) (string, bool) {
	for key, response := range commands {
		if strings.HasPrefix(input, key) {
			return response.Message(input), response.Close
		}
	}

	return input, false
}
