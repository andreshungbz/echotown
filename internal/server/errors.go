package server

import "errors"

var (
	ERROR_LONG_MSG = errors.New("[ERROR] Message cannot be longer than 1024 bytes!")
	ERROR_NON_PRNT = errors.New("[ERROR] Message contains non-printable characters!")
	ERROR_BAD_UTF8 = errors.New("[ERROR] Message contains invalid UTF-8 characters!")
)
