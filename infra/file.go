package infra

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

type (
	// FileWriter implements the domain.FileWriter interface.
	FileWriter struct{}
)

// Append implements the domain.FileWriter interface.
func (writer FileWriter) Append(dest string, data []byte) error {
	f, err := os.OpenFile(dest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to open file: %s", dest))
	}
	defer f.Close()
	fmt.Printf("add a test skelton code to a test file: %s\n", dest)
	if _, err := f.Write(data); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write file: %s", dest))
	}
	return nil
}

// Create implements the domain.FileWriter interface.
func (writer FileWriter) Create(dest string, data []byte) error {
	if err := ioutil.WriteFile(dest, data, 0644); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create test file: %s", dest))
	}
	return nil
}
