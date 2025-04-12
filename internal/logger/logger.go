// Package logger contains functions that handle creating appropriate [net.Logger] instances.
package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// NewServer returns a [log.Logger] that writes to both standard output and to /log/echotown.log.
// Its second return value is a cleanup function for closing the file.
func NewServer() (*log.Logger, func(), error) {
	// create log directory if it doesn't exist
	// privides read and execution permissions to group and everyone, and additional write permission to owner
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
