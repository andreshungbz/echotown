// Package personality contains custom server responses to specific client inputs.
package personality

type personalityResponse struct {
	// message to be printed to client
	Message string

	// determines whether to close the connection
	Close bool
}

// extendible map of specific input to custom response
var responses = make(map[string]personalityResponse)

// initiate map
func init() {
	responses["hello"] = personalityResponse{
		Message: "Hi there!",
		Close:   false,
	}

	responses[""] = personalityResponse{
		Message: "Say something...",
		Close:   false,
	}

	responses["bye"] = personalityResponse{
		Message: "Goodbye!",
		Close:   true,
	}
}

// Parse compares the input to pre-determined keys, and if they match,
// the custom server response is returned. Otherwise return the same input.
// The returned boolean determines whether the connection should be closed.
func Parse(input string) (string, bool) {
	for key, response := range responses {
		if input == key {
			return response.Message, response.Close
		}
	}

	return input, false
}
