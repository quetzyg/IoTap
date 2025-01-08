package device

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// Script holds the path and the contents of an IoT device script.
type Script struct {
	path string
	code []byte
}

// Name of the file the script was loaded from.
func (s *Script) Name() string {
	return path.Base(s.path)
}

// Code returns the content of the script.
func (s *Script) Code() []byte {
	return s.code
}

// Length returns the length of the script content.
func (s *Script) Length() int {
	return len(s.code)
}

// NewScript creates a new *Script instance by parsing data from the provided reader.
// It returns an error if the data is invalid or cannot be parsed.
func NewScript(r io.Reader) (*Script, error) {
	var (
		src = &Script{}
		err error
	)

	src.code, err = io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if src.Length() == 0 {
		return nil, fmt.Errorf("%w: %s", ErrScriptEmpty, src.path)
	}

	return src, nil
}

// LoadScript creates a new *Script instance from a file at the given path.
// It returns an error if the file cannot be opened or contains invalid data.
func LoadScript(fp string) (*Script, error) {
	if fp == "" {
		return nil, ErrFilePathEmpty
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Script close error: %v", err)
		}
	}()

	src, err := NewScript(f)
	if err != nil {
		return nil, err
	}

	src.path = fp

	return src, nil
}

// LoadScripts creates a slice of *Script instances from multiple file paths.
// It returns an error if any of the files cannot be opened or contain invalid data.
func LoadScripts(fps []string) ([]*Script, error) {
	if len(fps) == 0 {
		return nil, ErrFilePathEmpty
	}

	var (
		err     error
		scripts = make([]*Script, len(fps))
	)

	for i, fp := range fps {
		if fp == "" {
			return nil, ErrFilePathEmpty
		}

		scripts[i], err = LoadScript(fp)
		if err != nil {
			return nil, err
		}
	}

	return scripts, nil
}
