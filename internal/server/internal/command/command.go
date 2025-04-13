// Package command contains custom server responses to specific command protocol strings
package command

import (
	"strings"
	"time"
)

type CommandResponse struct {
	// message function that returns a string to be printed to client
	Message func(string) string

	// determines whether to close the connection
	Close bool
}

// extendible map of specific command protocols to custom response
var commands = map[string]CommandResponse{
	"/time": {func(string) string { return time.Now().String() }, false},
	"/quit": {func(string) string { return "" }, true},
	"/echo": {func(input string) string {
		if len(input) > len("/echo") {
			return input[len("/echo "):]
		}
		return ""
	}, false},
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
