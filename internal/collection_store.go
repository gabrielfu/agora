package internal

import (
	"os"
	"path/filepath"
)

// default: ~/.agora
func defaultRootDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agora"), nil
}

type CollectionStore struct {
	root string
}

func (c *CollectionStore) Root() string {
	return c.root
}

func (c *CollectionStore) SetRoot(dir string) {
	c.root = dir
}

func (c *CollectionStore) SetDefaultRoot() error {
	dir, err := defaultRootDir()
	if err != nil {
		return err
	}
	c.SetRoot(dir)
	return nil
}

func (c *CollectionStore) CollectionDir(collection string) string {
	return filepath.Join(c.Root(), "collections", collection)
}

func (c *CollectionStore) CollectionRequestDir(collection string) string {
	return filepath.Join(c.CollectionDir(collection), "requests")
}

var DefaultCollectionStore = &CollectionStore{}
