package storage

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Filesystem struct {
	Root string
	mu   sync.RWMutex
}

func NewFilesystem(root string) (*Filesystem, error) {
	s := &Filesystem{
		Root: root,
	}
	return s, os.MkdirAll(s.Root, 0744)
}

// CleanPath removes any path transversal nonsense
func CleanPath(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return filepath.Clean(path)
}

// Takes a possibly multilevel path and flattens it by dropping any slashes
func FlattenPath(path string, separator string) string {
	return strings.Replace(path, string(os.PathSeparator), separator, -1)
}

func (s *Filesystem) SaveName(ctx context.Context, rawShort, url string) error {
	short, err := sanitizeShort(rawShort)
	if err != nil {
		return err
	}
	if _, err := validateURL(url); err != nil {
		return err
	}

	short = FlattenPath(CleanPath(short), "_")

	s.mu.Lock()
	err = ioutil.WriteFile(filepath.Join(s.Root, short), []byte(url), 0744)
	s.mu.Unlock()

	return err
}

func (s *Filesystem) Load(ctx context.Context, rawShort string) (string, error) {
	short, err := sanitizeShort(rawShort)
	if err != nil {
		return "", err
	}

	short = FlattenPath(CleanPath(short), "_")

	s.mu.Lock()
	urlBytes, err := ioutil.ReadFile(filepath.Join(s.Root, short))
	s.mu.Unlock()

	if _, ok := err.(*os.PathError); ok {
		return "", ErrShortNotSet
	}

	return string(urlBytes), err
}
