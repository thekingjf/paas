package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

type server struct {
	db *sql.DB
}

func main() {
	r := chi.NewRouter()
	db, err := sql.Open("sqlite", "paas.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err2 := db.Exec(`CREATE TABLE IF NOT EXISTS apps (
		name TEXT PRIMARY KEY,
		container_id TEXT,
		port INTEGER,
		status TEXT
	)`)

	if err2 != nil {
		log.Fatal(err2)
	}

	s := &server{db: db}

	r.Post("/apps", s.ok)
	r.Get("/apps/{id}", s.getAppId)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func (s *server) ok(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("ok"))
}

func (s *server) getAppId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	//result, err := s.db.Exec(fmt.Sprintf(`SELECT name, port, status WHERE container_id = %s`, id))

	// if err != nil {
	// 	log.Fatal(err)
	// }

	w.Write([]byte(id))
}
