package device

import (
	"fmt"
	"io"
	"log"
	"os"
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

// NewIoTScript creates a new *Script instance.
func NewIoTScript(name string) *Script {
	return &Script{
		name: name,
	}
}

// loadScript from an I/O reader and read the data into an *Script instance.
func loadScript(r io.Reader, src *Script) (err error) {
	src.code, err = io.ReadAll(r)
	if err != nil {
		return err
	}

	if src.Length() == 0 {
		return ErrScriptEmpty
	}

	return nil
}

// LoadScriptFromPath reads and loads a script from the specified file path.
// It opens the file, ensures it's closed after reading, and processes the script.
// Returns an error if the file cannot be opened or the script cannot be loaded.
func LoadScriptFromPath(fp string, src *Script) error {
	if fp == "" {
		return fmt.Errorf("the script file path cannot be empty")
	}

	f, err := os.Open(fp)
	if err != nil {
		return err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.Printf("Script close error: %v", err)
		}
	}()

	return loadScript(f, src)
}
