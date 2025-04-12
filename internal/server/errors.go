package server

import (
	"errors"
	"fmt"
	"net"
)

var (
	ERROR_LONG_MSG = errors.New("Message cannot be longer than 1024 bytes!")
	ERROR_NON_PRNT = errors.New("Message contains non-printable characters!")
	ERROR_BAD_UTF8 = errors.New("Message contains invalid UTF-8 characters!")
)

func createError(message string, address net.Addr, err error) string {
	return fmt.Sprintf("[ERROR] [%v] %v: %v\n", address, message, err)
}
