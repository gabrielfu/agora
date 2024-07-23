package internal

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type RequestFileStore struct {
	root string
}

func NewRquestFileStore(root string) (*RequestFileStore, error) {
	if err := os.Mkdir(root, 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}
	return &RequestFileStore{root: root}, nil
}

func (r *RequestFileStore) calcFilename(id string) string {
	return filepath.Join(r.root, id)
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

func (r *RequestFileStore) CreateRequest(req Request) error {
	data, err := yaml.Marshal(req)
	if err != nil {
		return err
	}
	filename := r.calcFilename(req.id)
	return os.WriteFile(filename, data, 0755)
}

func (r *RequestFileStore) GetRequest(id string) (Request, error) {
	filename := r.calcFilename(id)
	return readFile(filename)
}

func (r *RequestFileStore) ListRequests() ([]Request, error) {
	entries, err := os.ReadDir(r.root)
	if err != nil {
		return nil, err
	}

	read := func(filename string, wg *sync.WaitGroup, ch chan<- Request, ech chan<- error) {
		defer wg.Done()
		req, err := readFile(filename)
		if err != nil {
			ech <- err
		} else {
			ch <- req
		}
	}

	var wg sync.WaitGroup
	ch := make(chan Request)
	ech := make(chan error)
	for _, e := range entries {
		wg.Add(1)
		filename := filepath.Join(r.root, e.Name())
		go read(filename, &wg, ch, ech)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(ech)
	}()

	for err := range ech {
		return nil, err
	}

	var reqs []Request
	for req := range ch {
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (r *RequestFileStore) UpdateRequest(req Request) error {
	return r.CreateRequest(req)
}

func (r *RequestFileStore) DeleteRequest(id string) error {
	filename := r.calcFilename(id)
	return os.Remove(filename)
}
