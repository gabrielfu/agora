package internal

import (
	"database/sql"
	"encoding/json"

	_ "github.com/mattn/go-sqlite3"
)

type RequestDAO struct {
	id      string
	Method  string
	URL     string
	Body    string
	Params  string
	Headers string
	Auth    string
}

func fromRequest(req Request) (RequestDAO, error) {
	body, err := json.Marshal(req.Body)
	if err != nil {
		return RequestDAO{}, err
	}
	params, err := json.Marshal(req.Params)
	if err != nil {
		return RequestDAO{}, err
	}
	headers, err := json.Marshal(req.Headers)
	if err != nil {
		return RequestDAO{}, err
	}
	return RequestDAO{
		id:      req.ID(),
		Method:  req.Method,
		URL:     req.URL,
		Body:    string(body),
		Params:  string(params),
		Headers: string(headers),
		Auth:    req.Auth,
	}, nil
}

func (r *RequestDAO) toRequest() (Request, error) {
	var body, params, headers map[string]string
	if err := json.Unmarshal([]byte(r.Body), &body); err != nil {
		return Request{}, err
	}
	if err := json.Unmarshal([]byte(r.Params), &params); err != nil {
		return Request{}, err
	}
	if err := json.Unmarshal([]byte(r.Headers), &headers); err != nil {
		return Request{}, err
	}
	return Request{
		id:      r.id,
		Method:  r.Method,
		URL:     r.URL,
		Body:    body,
		Params:  params,
		Headers: headers,
		Auth:    r.Auth,
	}, nil
}

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
	r := &RequestDatabase{db: db}
	if err = r.Init(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *RequestDatabase) Close() error {
	return r.db.Close()
}

func (r *RequestDatabase) Init() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS requests (
			id TEXT PRIMARY KEY,
			method TEXT,
			url TEXT,
			body TEXT,
			params TEXT,
			headers TEXT,
			auth TEXT
		)
	`)
	return err
}

func (r *RequestDatabase) CreateRequest(req Request) error {
	dao, err := fromRequest(req)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		"INSERT INTO requests (id, method, url, body, params, headers, auth) VALUES (?, ?, ?, ?, ?, ?, ?)",
		dao.id, dao.Method, dao.URL, dao.Body, dao.Params, dao.Headers, dao.Auth,
	)
	return err
}

func (r *RequestDatabase) GetRequest(id string) (Request, error) {
	var dao RequestDAO
	err := r.db.QueryRow(
		"SELECT id, method, url, body, params, headers, auth FROM requests WHERE id = ?",
		id,
	).Scan(&dao.id, &dao.Method, &dao.URL, &dao.Body, &dao.Params, &dao.Headers, &dao.Auth)
	if err != nil {
		return Request{}, err
	}
	return dao.toRequest()
}

func (r *RequestDatabase) ListRequests() ([]Request, error) {
	rows, err := r.db.Query("SELECT id, method, url, body, params, headers, auth FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []Request
	for rows.Next() {
		var dao RequestDAO
		if err := rows.Scan(&dao.id, &dao.Method, &dao.URL, &dao.Body, &dao.Params, &dao.Headers, &dao.Auth); err != nil {
			return nil, err
		}
		req, err := dao.toRequest()
		if err != nil {
			return nil, err
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func (r *RequestDatabase) UpdateRequest(req Request) error {
	dao, err := fromRequest(req)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(
		"UPDATE requests SET method = ?, url = ?, body = ?, params = ?, headers = ?, auth = ? WHERE id = ?",
		dao.Method, dao.URL, dao.Body, dao.Params, dao.Headers, dao.Auth, dao.id,
	)
	return err
}

func (r *RequestDatabase) DeleteRequest(id string) error {
	_, err := r.db.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}
