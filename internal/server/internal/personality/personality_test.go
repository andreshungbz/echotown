package personality

import (
	"testing"
)

func TestPersonalityResponses(t *testing.T) {
	t.Run("Matching Inputs", func(t *testing.T) {
		for input, expected := range responses {
			gotMsg, gotClose := Parse(input)

			if gotMsg != expected.Message || gotClose != expected.Close {
				t.Errorf("Parse(%q) = (%q, %v); want (%q, %v)",
					input, gotMsg, gotClose, expected.Message, expected.Close)
			}
		}
	})

	t.Run("Non-matching Input", func(t *testing.T) {
		input := "test"
		var close bool

		input, close = Parse(input)

		if input != "test" || close != false {
			t.Errorf("Parse(%q) = (%q, %v); want (%q, %v)",
				input, "test", false, input, close)
		}
	})
}
