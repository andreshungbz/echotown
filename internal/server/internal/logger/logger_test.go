package logger

import (
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	logger, cleanup, err := NewServer()

	if err != nil {
		t.Fatalf("NewServer() error = %v; want nil", err)
	}

	if logger == nil {
		t.Fatal("NewServer() returned nil logger; want non-nil")
	}

	defer cleanup()
}

func TestNewClient(t *testing.T) {
	addr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: 12345,
	}

	logger, cleanup, err := NewClient(addr)

	if err != nil {
		t.Fatalf("NewClient() error = %v; want nil", err)
	}

	if logger == nil {
		t.Fatal("NewClient() returned nil logger; want non-nil")
	}

	defer cleanup()
}
