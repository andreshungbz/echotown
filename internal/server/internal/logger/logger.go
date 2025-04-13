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

// NewServer returns a [log.Logger] that writes to both standard output and to log/echotown.log.
// Its second return value is a cleanup function for closing the file.
func NewServer() (*log.Logger, func(), error) {
	// create log directory if it doesn't exist
	// provides read and execution permissions to group and everyone, and additional write permission to owner
	logPath := "log/echotown.log"
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		return nil, nil, err
	}

	// creates file if it doesn't exist, opens in write-only mode, and appends if it exists already
	// provides read and write permissions to owner, group and everyone
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// write both to the console and to file
	mw := io.MultiWriter(os.Stdout, file)
	logger := log.New(mw, "", log.Ldate|log.Ltime)

	// provide a cleanup function to close the file properly
	return logger, func() { file.Close() }, nil
}

// NewClient returns a [log.Logger] that writes to log/{addr}.log where addr is a
// passed in [net.Addr] whose periods and colons are replaced with underscores.
// Its second return value is a cleanup function for closing the file.
func NewClient(addr net.Addr) (*log.Logger, func(), error) {
	// replace periods and colons of address to underscores and append .log
	// e.g. 192.168.18.125:50000 becomes 192_168_18_125_p50000.log
	addrStr := strings.ReplaceAll(strings.ReplaceAll(addr.String(), ".", "_"), ":", "_p") + ".log"

	// create log directory if it doesn't exist
	// provides read and execution permissions to group and everyone, and additional write permission to owner
	logPath := fmt.Sprintf("log/%s", addrStr)
	err := os.MkdirAll(filepath.Dir(logPath), 0755)
	if err != nil {
		return nil, nil, err
	}

	// creates file if it doesn't exist, opens in write-only mode, and appends if it exists already
	// provides read and write permissions to owner, group and everyone
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, nil, err
	}

	// write to file only
	logger := log.New(file, "", log.Ldate|log.Ltime)

	// provide a cleanup function to close the file properly
	return logger, func() { file.Close() }, nil
}
