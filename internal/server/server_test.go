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

	t.Run("Client Message Too Long", func(t *testing.T) {
		conn, reader := initiateConn(port, t)
		defer conn.Close()

		longMessage := strings.Repeat("A", 2048) + "\n"

		_, err := conn.Write([]byte(longMessage))
		if err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}

		assertResponse(ERROR_LONG_MSG.Error(), &reader, t)
	})

	t.Run("Client Message Contains Non-Printable Characters", func(t *testing.T) {
		conn, reader := initiateConn(port, t)
		defer conn.Close()

		badChars := "Hello\x00World\n" // null byte is non-printable
		_, err := conn.Write([]byte(badChars))
		if err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}

		assertResponse(ERROR_NON_PRNT.Error(), &reader, t)
	})

	t.Run("Client Message Contains Invalid UTF-8", func(t *testing.T) {
		conn, reader := initiateConn(port, t)
		defer conn.Close()

		// Invalid UTF-8 byte sequence
		invalidUTF8 := []byte{0xff, 0xfe, 0xfd, '\n'}
		_, err := conn.Write(invalidUTF8)
		if err != nil {
			t.Fatalf("Failed to send message: %v", err)
		}

		assertResponse(ERROR_BAD_UTF8.Error(), &reader, t)
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

	// clear server prepend
	_, _ = reader.ReadString(' ')
	_, _ = reader.ReadString(' ')

	// read response
	response, err := reader.ReadString('\n')
	if err != nil {
		t.Fatalf("Failed to read echo response: %v", err)
	}

	// clear whitespace
	response = strings.TrimSpace(response)

	// validate
	if response != want {
		t.Errorf("Expected %q, got: %q", want, response)
	}
}
