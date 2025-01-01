package command

import "errors"

var (
	ErrNotFound      = errors.New("command not found")
	ErrInvalid       = errors.New("invalid command")
	ErrArgumentParse = errors.New("error parsing argument")
	ErrFlagConflict  = errors.New("conflicting command flags")
)
