package command

import "errors"

// Predefined errors for command handling, covering several common cases.
var (
	ErrNotFound      = errors.New("command not found")
	ErrInvalid       = errors.New("invalid command")
	ErrArgumentParse = errors.New("error parsing argument")
	ErrFlagConflict  = errors.New("conflicting command flags")
)
