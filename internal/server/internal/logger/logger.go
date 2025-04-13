// Package logger contains functions that handle creating appropriate [net.Logger] instances.
package logger

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// NewServer returns a [log.Logger] that writes to both standard output and to log/echotown_{addr}.log
// where addr is a formatted filename based on the passed in [net.Addr].
// Its second return value is a cleanup function for closing the file.
func NewServer(addr net.Addr) (*log.Logger, func(), error) {
	fileName := "echotown_" + formatFileString(addr.String()) + ".log"

	// create log directory if it doesn't exist
	logPath := fmt.Sprintf("log/%s", fileName)
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		return nil, nil, err
	}

	// creates file if it doesn't exist
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// write to both the console and file
	mw := io.MultiWriter(os.Stdout, file)
	logger := log.New(mw, "", log.Ldate|log.Ltime)

	return logger, func() { file.Close() }, nil
}

// NewClient returns a [log.Logger] that writes to log/{addr}.log where addr is a
// formatted filename based on the passed in [net.Addr].
// Its second return value is a cleanup function for closing the file.
func NewClient(addr net.Addr) (*log.Logger, func(), error) {
	fileName := "client_" + formatFileString(addr.String()) + ".log"

	// create log directory if it doesn't exist
	logPath := fmt.Sprintf("log/%s", fileName)
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		return nil, nil, err
	}

	// creates file if it doesn't exist
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// write to file only
	logger := log.New(file, "", log.Ldate|log.Ltime)

	return logger, func() { file.Close() }, nil
}

// formatFileString replaces the IPv4 and IPv6 loopback address to the string "localhost",
// replaces periods with underscores, and then replaces colons with _p.
// e.g. 192.168.18.125:50000 becomes 192_168_18_125_p50000.
func formatFileString(addrStr string) string {
	// replace IPv6 loopback address with "localhost"
	addrStr = strings.Replace(addrStr, "[::1]", "localhost", 1)

	// replace IPv4 loopback address with "localhost"
	addrStr = strings.Replace(addrStr, "127.0.0.1", "localhost", 1)

	// replace periods with underscores
	addrStr = strings.ReplaceAll(addrStr, ".", "_")

	// replace port extention with _p
	addrStr = strings.ReplaceAll(addrStr, ":", "_p")

	return addrStr
}
