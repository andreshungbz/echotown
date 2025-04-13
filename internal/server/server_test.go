package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"
)

func TestEchoServerResponse(t *testing.T) {
	// start a test server on port 4001
	port := 4001
	go Start(port)
	time.Sleep(100 * time.Millisecond) // give server time to start up

	t.Run("Validation Responses", func(t *testing.T) {
		tests := []struct {
			name           string
			input          string
			expectedOutput string
		}{
			{"Client Message Too Long", strings.Repeat("A", 2048) + "\n", ERROR_LONG_MSG.Error()},
			{"Client Message Contains Non-Printable Characters", "Hello\x00World\n", ERROR_NON_PRNT.Error()},
			{"Client Message Contains Invalid UTF-8", string([]byte{0xff, 0xfe, 0xfd, '\n'}), ERROR_BAD_UTF8.Error()},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				conn, reader := initiateConn(port, t)
				defer conn.Close()

				_, err := conn.Write([]byte(tt.input))
				if err != nil {
					t.Fatalf("Failed to send message: %v", err)
				}

				assertResponse(tt.expectedOutput, &reader, t)
			})
		}
	})

	time.Sleep(100 * time.Millisecond) // give time to log last disconnect
}

// initiateConn creates and returns a [net.Conn] and associated [bufio.Reader]
// with the initial server messages already cleared.
func initiateConn(port int, t *testing.T) (net.Conn, bufio.Reader) {
	t.Helper()

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	reader := bufio.NewReader(conn)
	_, _ = reader.ReadString('\n') // welcome line
	_, _ = reader.ReadString('\n') // connection line
	_, _ = reader.ReadString(' ')  // prompt

	return conn, *reader
}

// assertResponse checks to see if the respone string matches properly.
func assertResponse(want string, reader *bufio.Reader, t *testing.T) {
	t.Helper()

	// read response
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read echo response: %v", err)
	}

	// clear whitespace
	response = strings.TrimSpace(response)

	// validate
	if !strings.Contains(response, want) {
		t.Errorf("Expected %q, got: %q", want, response)
	}
}
