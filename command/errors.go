package command

import "errors"

var (
	ErrNotFound      = errors.New("command not found")
	ErrArgumentParse = errors.New("error parsing argument")
	ErrInvalid       = errors.New("invalid command")
)
