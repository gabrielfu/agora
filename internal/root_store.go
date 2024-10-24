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

type RootStore struct {
	root string
}

func NewRootStore(root string) *RootStore {
	return &RootStore{root: root}
}

func NewDefaultRootStore() (*RootStore, error) {
	dir, err := defaultRootDir()
	if err != nil {
		return nil, err
	}
	return NewRootStore(dir), nil
}

func (r *RootStore) Root() string {
	return r.root
}

func (r *RootStore) SetRoot(dir string) {
	r.root = dir
}

func (r *RootStore) ListCollections() ([]string, error) {
	collectionsDir := filepath.Join(r.Root(), "collections")
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

func (r *RootStore) CollectionDir(collection string) string {
	return filepath.Join(r.Root(), "collections", collection)
}

func (r *RootStore) CollectionExists(collection string) bool {
	_, err := os.Stat(r.CollectionDir(collection))
	return err == nil
}

func (r *RootStore) CreateCollection(collection string) error {
	return os.MkdirAll(r.CollectionDir(collection), 0755)
}

func (r *RootStore) DeleteCollection(collection string) error {
	return os.RemoveAll(r.CollectionDir(collection))
}

func (r *RootStore) RenameCollection(oldName, newName string) error {
	return os.Rename(r.CollectionDir(oldName), r.CollectionDir(newName))
}
