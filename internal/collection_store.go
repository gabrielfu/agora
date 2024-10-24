package internal

import (
	"os"
	"path/filepath"
)

const DEFAULT_COLLECTION_NAME = "default"

// default: ~/.agora
func defaultRootDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".agora"), nil
}

type CollectionStore struct {
	root       string
	collection string
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

func (c *CollectionStore) Collection() string {
	return c.collection
}

func (c *CollectionStore) SetCollection(collection string) {
	c.collection = collection
}

func (c *CollectionStore) CollectionDir() string {
	return filepath.Join(c.Root(), "collections", c.collection)
}

func (c *CollectionStore) CollectionRequestDir() string {
	return filepath.Join(c.CollectionDir(), "requests")
}

var DefaultCollectionStore = &CollectionStore{
	collection: DEFAULT_COLLECTION_NAME,
}
