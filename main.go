package main

import (
	"archive/tar"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/moby/moby/client"
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

	dockerCli, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer dockerCli.Close()

	s := &server{db: db, docker: dockerCli}

	r.Post("/apps", s.create)
	r.Get("/apps/{name}", s.getAppId)
	r.Post("/apps/{name}/deploy", s.deploy)

	log.Fatal(http.ListenAndServe(":8080", r))

}

func (s *server) deploy(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	options := client.ImageBuildOptions{Tags: []string{name + ":latest"},
		Dockerfile: "Dockerfile"}

	result, err := s.docker.ImageBuild(r.Context(), r.Body, options)

	if err != nil {
		log.Println(err)
		http.Error(w, "build failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	defer result.Body.Close()
	io.Copy(os.Stdout, result.Body)
}

func (s *server) create(w http.ResponseWriter, r *http.Request) {
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

func streamIn(folderLocation string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	folder, err := filepath.Abs(folderLocation)
	writer := tar.NewWriter(&buf)

	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	for _, item := range files {
		err = streamInHelper(folderLocation, folder, "", item, writer)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()
	return &buf, nil
}

func streamInDive(folderLocation, src string, writer *tar.Writer) error {

	folder := filepath.Join(folderLocation, src)

	files, err := os.ReadDir(folder)

	if err != nil {
		return err
	}

	for _, item := range files {
		err = streamInHelper(folderLocation, folder, src, item, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func streamInHelper(folderLocation, folder, src string, item os.DirEntry, writer *tar.Writer) error {
	if item.IsDir() {
		err := streamInDive(folderLocation, filepath.Join(src, item.Name()), writer)

		if err != nil {
			return err
		}
	} else {
		file, err := item.Info()

		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(file, "")

		if err != nil {
			return err
		}

		header.Name = filepath.Join(src, file.Name())
		err = writer.WriteHeader(header)

		if err != nil {
			return err
		}

		book, err := os.ReadFile(filepath.Join(folder, file.Name()))

		if err != nil {
			return err
		}

		_, err = writer.Write(book)

		if err != nil {
			return err
		}
	}

	return nil
}
