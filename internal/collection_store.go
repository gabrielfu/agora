package internal

import (
	"fmt"
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
	root              string
	currentCollection string
}

func NewCollectionStore(root string) (*CollectionStore, error) {
	// handles case where default collection was renamed by user
	c := &CollectionStore{root: root, currentCollection: DEFAULT_COLLECTION_NAME}
	if err := c.CreateCollection(c.currentCollection); err != nil && !os.IsExist(err) {
		return nil, err
	}
	return c, nil
}

func NewDefaultCollectionStore() (*CollectionStore, error) {
	dir, err := defaultRootDir()
	if err != nil {
		return nil, err
	}
	return NewCollectionStore(dir)
}

func (c *CollectionStore) Root() string {
	return c.root
}

func (c *CollectionStore) SetRoot(dir string) {
	c.root = dir
}

func (c *CollectionStore) GetFirstCollection() (string, error) {
	collectionsDir := filepath.Join(c.Root(), "collections")
	entries, err := os.ReadDir(collectionsDir)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		return entry.Name(), nil
	}
	return "", fmt.Errorf("no collections found")
}

func (c *CollectionStore) ListCollections() ([]string, error) {
	collectionsDir := filepath.Join(c.Root(), "collections")
	entries, err := os.ReadDir(collectionsDir)
	if err != nil {
		return nil, err
	}
	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	return dirs, nil
}

func (c *CollectionStore) CollectionDir(collection string) string {
	return filepath.Join(c.Root(), "collections", collection)
}

func (c *CollectionStore) CollectionExists(collection string) bool {
	_, err := os.Stat(filepath.Join(c.Root(), "collections", collection))
	return err == nil
}

func (c *CollectionStore) CreateCollection(collection string) error {
	return os.MkdirAll(c.CollectionDir(collection), 0755)
}

func (c *CollectionStore) DeleteCollection(collection string) error {
	return os.RemoveAll(c.CollectionDir(collection))
}

func (c *CollectionStore) RenameCollection(oldName, newName string) error {
	return os.Rename(c.CollectionDir(oldName), c.CollectionDir(newName))
}

func (c *CollectionStore) CurrentCollection() string {
	return c.currentCollection
}

func (c *CollectionStore) SetCurrentCollection(collection string) {
	c.currentCollection = collection
}

func (c *CollectionStore) CurrentCollectionDir() string {
	return c.CollectionDir(c.currentCollection)
}

func (c *CollectionStore) CurrentCollectionRequestDir() string {
	return filepath.Join(c.CurrentCollectionDir(), "requests")
}
