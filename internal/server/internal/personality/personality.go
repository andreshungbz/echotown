// Package personality contains custom server responses to specific client inputs.
package personality

type PersonalityResponse struct {
	// message to be printed to client
	Message string

	// determines whether to close the connection
	Close bool
}

// extendible map of specific input to custom response
var responses = make(map[string]PersonalityResponse)

// initiate map
func init() {
	responses["hello"] = PersonalityResponse{
		Message: "Hi there!",
		Close:   false,
	}

	responses[""] = PersonalityResponse{
		Message: "Say something...",
		Close:   false,
	}

	responses["bye"] = PersonalityResponse{
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
