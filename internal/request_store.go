package internal

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// File store of a single collection
type RequestFileStore struct {
	root string
}

func NewRequestFileStore(root string) (*RequestFileStore, error) {
	if err := os.MkdirAll(root, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &RequestFileStore{root: root}, nil
}

func (r *RequestFileStore) calcRequestFilename(id string) string {
	return filepath.Join(r.root, id)
}

func readRequestFile(filename string) (Request, error) {
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

// Catalog file is a list of request IDs
// and maintains the order of requests

func (r *RequestFileStore) calcCatalogFilename() string {
	return filepath.Join(r.root, ".catalog")
}

func (r *RequestFileStore) createCatalogIfNotExists() error {
	if _, err := os.Stat(r.calcCatalogFilename()); os.IsNotExist(err) {
		requests, err := r.listRequestsUnordered()
		if err != nil {
			return err
		}
		catalog := make([]string, len(requests))
		for i, req := range requests {
			catalog[i] = req.ID
		}
		return r.WriteCatalog(catalog)
	}
	return nil
}
func (r *RequestFileStore) addToCatalog(id string) error {
	catalog, err := r.ReadCatalog()
	if err != nil {
		return err
	}
	for _, v := range catalog {
		if v == id {
			return nil
		}
	}
	catalog = append(catalog, id)
	return r.WriteCatalog(catalog)
}

func (r *RequestFileStore) removeFromCatalog(id string) error {
	catalog, err := r.ReadCatalog()
	if err != nil {
		return err
	}
	for i, v := range catalog {
		if v == id {
			catalog = append(catalog[:i], catalog[i+1:]...)
			return r.WriteCatalog(catalog)
		}
	}
	return nil
}

func (r *RequestFileStore) WriteCatalog(catalog []string) error {
	data, err := yaml.Marshal(catalog)
	if err != nil {
		return err
	}
	filename := r.calcCatalogFilename()
	return os.WriteFile(filename, data, 0644)
}

func (r *RequestFileStore) ReadCatalog() ([]string, error) {
	err := r.createCatalogIfNotExists()
	if err != nil {
		return nil, err
	}
	filename := r.calcCatalogFilename()
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var catalog []string
	err = yaml.Unmarshal(data, &catalog)
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

func (r *RequestFileStore) CreateRequest(req Request) error {
	data, err := yaml.Marshal(req)
	if err != nil {
		return err
	}
	filename := r.calcRequestFilename(req.ID)
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return r.addToCatalog(req.ID)
}

func (r *RequestFileStore) GetRequest(id string) (Request, error) {
	filename := r.calcRequestFilename(id)
	return readRequestFile(filename)
}

func (r *RequestFileStore) listRequestsUnordered() ([]Request, error) {
	rawEntries, err := os.ReadDir(r.root)
	if err != nil {
		return nil, err
	}
	var entries []os.DirEntry
	for _, e := range rawEntries {
		if e.IsDir() {
			continue
		}
		if e.Name() == ".catalog" {
			continue
		}
		entries = append(entries, e)
	}

	var wg sync.WaitGroup
	reqs := make([]Request, len(entries))
	errs := make([]error, len(entries))
	for i, e := range entries {
		wg.Add(1)
		go func(i int, e os.DirEntry) {
			defer wg.Done()
			filename := filepath.Join(r.root, e.Name())
			req, err := readRequestFile(filename)
			if err != nil {
				errs[i] = err
				return
			}
			reqs[i] = req
		}(i, e)
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return reqs, nil
}

func (r *RequestFileStore) ListRequests() ([]Request, error) {
	requests, err := r.listRequestsUnordered()
	if err != nil {
		return nil, err
	}
	catalog, err := r.ReadCatalog()
	if err != nil {
		return nil, err
	}
	// sort requests by catalog order
	sortMap := make(map[string]int, len(catalog))
	for i, id := range catalog {
		sortMap[id] = i
	}
	sortedRequests := make([]Request, len(requests))
	for _, id := range catalog {
		for _, req := range requests {
			if req.ID == id {
				sortedRequests[sortMap[id]] = req
				break
			}
		}
	}
	return sortedRequests, nil
}

func (r *RequestFileStore) UpdateRequest(req Request) error {
	return r.CreateRequest(req)
}

func (r *RequestFileStore) DeleteRequest(id string) error {
	filename := r.calcRequestFilename(id)
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return r.removeFromCatalog(id)
}
