package command

import (
	"strings"
	"testing"
	"time"
)

func TestCommandResponses(t *testing.T) {
	t.Run("Matching Commands", func(t *testing.T) {
		for input, expected := range commands {
			gotMsg, gotClose := Parse(input)

			switch input {
			case "/time":
				expectedMsg := time.Now().String()
				if gotClose != expected.Close {
					t.Errorf("Parse(%q) = (%q, %v); want (%q, %v)", input, gotMsg, gotClose, expectedMsg, expected.Close)
				}

			case "/help":
				expectedLines := []string{
					"/time - Displays the server time.",
					"/quit - Closes the connection to the server.",
					"/echo - Returns only the message. Same result if not used.",
					"/help - Lists all available commands and their descriptions.",
				}

				for _, line := range expectedLines {
					if !strings.Contains(gotMsg, line) {
						t.Errorf("Parse(%q) output missing line %q\nGot:\n%s", input, line, gotMsg)
					}
				}

				if gotClose != expected.Close {
					t.Errorf("Parse(%q) = close: %v; expected %v", input, gotClose, expected.Close)
				}

			default:
				wantMsg := expected.Message(input)
				if gotMsg != wantMsg || gotClose != expected.Close {
					t.Errorf("Parse(%q) = (%q, %v); want (%q, %v)",
						input, gotMsg, gotClose, wantMsg, expected.Close)
				}
			}
		}
	})
}
