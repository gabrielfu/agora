package internal

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type RequestDatabase struct {
	db *sql.DB
}

func NewRquestDatabase(path string) (*RequestDatabase, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &RequestDatabase{db: db}, nil
}

func (r *RequestDatabase) Close() error {
	return r.db.Close()
}

func (r *RequestDatabase) CreateRequest(req Request) error {
	_, err := r.db.Exec(
		"INSERT INTO requests (id, method, url, body, params, headers, auth) VALUES (?, ?, ?, ?, ?, ?, ?)",
		req.ID(), req.Method, req.URL, req.Body, req.Params, req.Headers, req.Auth,
	)
	return err
}

func (r *RequestDatabase) GetRequest(id string) (Request, error) {
	var req Request
	err := r.db.QueryRow(
		"SELECT id, method, url, body, params, headers, auth FROM requests WHERE id = ?",
		id,
	).Scan(&req.id, &req.Method, &req.URL, &req.Body, &req.Params, &req.Headers, &req.Auth)
	return req, err
}

func (r *RequestDatabase) ListRequests() ([]Request, error) {
	rows, err := r.db.Query("SELECT id, method, url, body, params, headers, auth FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []Request
	for rows.Next() {
		var req Request
		if err := rows.Scan(&req.id, &req.Method, &req.URL, &req.Body, &req.Params, &req.Headers, &req.Auth); err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *RequestDatabase) UpdateRequest(req Request) error {
	_, err := r.db.Exec(
		"UPDATE requests SET method = ?, url = ?, body = ?, params = ?, headers = ?, auth = ? WHERE id = ?",
		req.Method, req.URL, req.Body, req.Params, req.Headers, req.Auth, req.ID,
	)
	return err
}

func (r *RequestDatabase) DeleteRequest(id string) error {
	_, err := r.db.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}
