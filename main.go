package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/docker/docker/client"
	"github.com/go-chi/chi/v5"
	sqlite "modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type server struct {
	db     *sql.DB
	docker *client.Client
}

type createAppReq struct {
	Name string `json:"name"`
}

type appInfo struct {
	Name         string `json:"name"`
	Container_id string `json:"container_id"`
	Port         int    `json:"port"`
	Status       string `json:"status"`
}

func main() {
	r := chi.NewRouter()
	db, err := sql.Open("sqlite", "paas.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS apps (
		name TEXT PRIMARY KEY,
		container_id TEXT,
		port INTEGER,
		status TEXT
	)`)

	if err != nil {
		log.Fatal(err)
	}

	dockerCli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	defer dockerCli.Close()

	s := &server{db: db, docker: dockerCli}

	r.Post("/apps", s.ok)
	r.Get("/apps/{name}", s.getAppId)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (s *server) ok(w http.ResponseWriter, r *http.Request) {
	var req createAppReq

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Empty Name", http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("INSERT INTO apps (name, container_id, port, status) VALUES (?, ?, ?, ?)",
		req.Name, "", 0, "created")

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_PRIMARYKEY {
		http.Error(w, "app name already taken", http.StatusConflict)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *server) getAppId(w http.ResponseWriter, r *http.Request) {
	var app appInfo
	name := chi.URLParam(r, "name")
	err := s.db.
		QueryRow(`SELECT name, container_id, port, status FROM apps WHERE name = ?`, name).
		Scan(&app.Name, &app.Container_id, &app.Port, &app.Status)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&app)
}
