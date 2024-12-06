package device

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

// Script holds the name and the contents of an IoT device script.
type Script struct {
	name string
	code []byte
}

// Name of the file the script was loaded from.
func (s *Script) Name() string {
	return s.name
}

// Code returns the content of the script.
func (s *Script) Code() []byte {
	return s.code
}

// Length returns the length of the script content.
func (s *Script) Length() int {
	return len(s.code)
}

// loadScript from an I/O reader and read the data into an *Script instance.
func loadScript(r io.ReadCloser, src *Script) (err error) {
	defer func() {
		err = r.Close()
		if err != nil {
			log.Printf("Script close error: %v", err)
		}
	}()

	src.code, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	if src.Length() == 0 {
		return ErrScriptEmpty
	}

	return nil
}

// LoadScriptsFromPath reads and loads one or more script file paths.
// It opens each file and processes the script.
// An error is returned if a file cannot be opened or a script cannot be loaded.
func LoadScriptsFromPath(fps []string) ([]*Script, error) {
	scripts := make([]*Script, len(fps))

	for i, fp := range fps {
		if fp == "" {
			return nil, fmt.Errorf("the script file path cannot be empty")
		}

		f, err := os.Open(fp)
		if err != nil {
			return nil, err
		}

		scripts[i] = &Script{
			name: path.Base(fp),
		}

		err = loadScript(f, scripts[i])
		if err != nil {
			return nil, err
		}
	}

	return scripts, nil
}
