package command

import (
	"testing"
	"time"
)

func TestCommandResponses(t *testing.T) {
	t.Run("Matching Commands", func(t *testing.T) {
		for input, expected := range commands {
			gotMsg, gotClose := Parse(input)

			if input == "/time" {
				expectedMsg := time.Now().String()
				if gotClose != expected.Close {
					t.Errorf("Parse(%q) = (%q, %v); want (%q, %v)", input, gotMsg, gotClose, expectedMsg, expected.Close)
				}
			} else {
				if gotMsg != expected.Message(input) || gotClose != expected.Close {
					t.Errorf("Parse(%q) = (%q, %v); want (%q, %v)",
						input, gotMsg, gotClose, expected.Message(input), expected.Close)
				}
			}
		}
	})
}
