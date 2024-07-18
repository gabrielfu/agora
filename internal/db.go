package internal

import (
	"database/sql"
	"encoding/json"

	"github.com/gabrielfu/agora/tui/styles"
	_ "github.com/mattn/go-sqlite3"
)

type requestDAO struct {
	id      string
	Name    string
	Method  string
	URL     string
	Body    string
	Params  string
	Headers string
	Auth    string
}

func fromRequest(req Request) (requestDAO, error) {
	body, err := json.Marshal(req.Body)
	if err != nil {
		return requestDAO{}, err
	}
	params, err := json.Marshal(req.Params)
	if err != nil {
		return requestDAO{}, err
	}
	headers, err := json.Marshal(req.Headers)
	if err != nil {
		return requestDAO{}, err
	}
	return requestDAO{
		id:      req.ID(),
		Name:    req.Name,
		Method:  req.Method,
		URL:     req.URL,
		Body:    styles.PrettifyJsonIfValid(string(body)),
		Params:  string(params),
		Headers: string(headers),
		Auth:    req.Auth,
	}, nil
}

func (r *requestDAO) toRequest() (Request, error) {
	var body string
	var params, headers KVPairs
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
		Name:    r.Name,
		Method:  r.Method,
		URL:     r.URL,
		Body:    styles.PrettifyJsonIfValid(body),
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
			name TEXT,
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
		"INSERT INTO requests (id, name, method, url, body, params, headers, auth) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		dao.id, dao.Name, dao.Method, dao.URL, dao.Body, dao.Params, dao.Headers, dao.Auth,
	)
	return err
}

func (r *RequestDatabase) GetRequest(id string) (Request, error) {
	var dao requestDAO
	err := r.db.QueryRow(
		"SELECT id, name, method, url, body, params, headers, auth FROM requests WHERE id = ?",
		id,
	).Scan(&dao.id, &dao.Name, &dao.Method, &dao.URL, &dao.Body, &dao.Params, &dao.Headers, &dao.Auth)
	if err != nil {
		return Request{}, err
	}
	return dao.toRequest()
}

func (r *RequestDatabase) ListRequests() ([]Request, error) {
	rows, err := r.db.Query("SELECT id, name, method, url, body, params, headers, auth FROM requests")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var requests []Request
	for rows.Next() {
		var dao requestDAO
		if err := rows.Scan(&dao.id, &dao.Name, &dao.Method, &dao.URL, &dao.Body, &dao.Params, &dao.Headers, &dao.Auth); err != nil {
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
		"UPDATE requests SET name = ?, method = ?, url = ?, body = ?, params = ?, headers = ?, auth = ? WHERE id = ?",
		dao.Name, dao.Method, dao.URL, dao.Body, dao.Params, dao.Headers, dao.Auth, dao.id,
	)
	return err
}

func (r *RequestDatabase) DeleteRequest(id string) error {
	_, err := r.db.Exec("DELETE FROM requests WHERE id = ?", id)
	return err
}
