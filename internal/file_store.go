package internal

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const REQUEST_FOLDER_NAME = "requests"

type CollectionRequest struct {
	Collection string
	Request    Request
}

type RequestFileStore struct {
	rootStore *RootStore
}

func NewRequestFileStore(rootStore *RootStore) *RequestFileStore {
	return &RequestFileStore{rootStore: rootStore}
}

func (r *RequestFileStore) calcFilename(collection, id string) string {
	return filepath.Join(r.rootStore.CollectionDir(collection), REQUEST_FOLDER_NAME, id)
}

func readFile(filename string) (Request, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return Request{}, err
	}

	var req Request
	err = yaml.Unmarshal(data, &req)
	if err != nil {
		return Request{}, err
	}
	return req, nil
}

func (r *RequestFileStore) CreateRequest(collection string, req Request) error {
	data, err := yaml.Marshal(req)
	if err != nil {
		return err
	}
	filename := r.calcFilename(collection, req.ID)
	return os.WriteFile(filename, data, 0755)
}

func (r *RequestFileStore) GetRequest(collection, id string) (Request, error) {
	filename := r.calcFilename(collection, id)
	return readFile(filename)
}

func (r *RequestFileStore) ListRequests() ([]CollectionRequest, error) {
	collections, err := r.rootStore.ListCollections()
	if err != nil {
		return nil, err
	}

	var collectionEntries []CollectionRequest
	for _, collection := range collections {
		collectionDir := r.rootStore.CollectionDir(collection)
		entries, err := os.ReadDir(collectionDir)
		if err != nil {
			return nil, err
		}

		var wg sync.WaitGroup
		reqs := make([]CollectionRequest, len(entries))
		errs := make([]error, len(entries))
		for i, e := range entries {
			wg.Add(1)
			go func(i int, e os.DirEntry) {
				defer wg.Done()
				filename := filepath.Join(collectionDir, e.Name())
				req, err := readFile(filename)
				if err != nil {
					errs[i] = err
					return
				}
				reqs[i] = CollectionRequest{Collection: collection, Request: req}
			}(i, e)
		}
		wg.Wait()

		for _, err := range errs {
			if err != nil {
				return nil, err
			}
		}
		collectionEntries = append(collectionEntries, reqs...)
	}
	return collectionEntries, nil
}

func (r *RequestFileStore) UpdateRequest(collection string, req Request) error {
	return r.CreateRequest(collection, req)
}

func (r *RequestFileStore) DeleteRequest(collection, id string) error {
	filename := r.calcFilename(collection, id)
	return os.Remove(filename)
}
