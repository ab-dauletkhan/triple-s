package core

import "errors"

var (
	ErrIncorrectPort = errors.New("incorrect port number, range must be between 1-65535")
	ErrEmptyDir      = errors.New("empty directory path")
)
