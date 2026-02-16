package fileio

import (
	"fmt"
	"os"
)

// Reader handles file reading operations
type Reader interface {
	ReadFile(path string) (string, error)
}

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (r *FileReader) ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("input not found: %w", err)
	}
	return string(content), nil
}

// Writer handles file writing operations
type Writer interface {
	WriteFile(path string, content string) error
}

type FileWriter struct{}

func NewFileWriter() *FileWriter {
	return &FileWriter{}
}

func (w *FileWriter) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

// Checker handles file existence checks
type Checker interface {
	Exists(path string) bool
}

type FileChecker struct{}

func NewFileChecker() *FileChecker {
	return &FileChecker{}
}

func (c *FileChecker) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
