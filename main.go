package main

import (
	"database/sql"
	"embed"
	"encoding/gob"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var (
	db             *sql.DB
	sessionManager *scs.SessionManager
)

//go:embed static
var static embed.FS

//go:embed tmpls
var tmpls embed.FS

func main() {
	gob.Register(User{})

	err := InitDB()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = InitCaddy()

	if err != nil {
		log.Fatal(err)
	}

	err = InitXray()

	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(static))

	mux.HandleFunc("GET /", HomeHandler)
	mux.HandleFunc("POST /", HomeHandler)
	mux.HandleFunc("GET /config", ConfigHandler)
	mux.HandleFunc("POST /config", ConfigHandler)
	mux.HandleFunc("GET /login", LoginHandler)
	mux.HandleFunc("POST /login", LoginHandler)
	mux.HandleFunc("POST /logout", LogoutHandler)
	mux.HandleFunc("GET /register", RegisterHandler)
	mux.HandleFunc("POST /register", RegisterHandler)

	sessionManager = scs.New()
	sessionManager.Store = sqlite3store.New(db)

	err = http.ListenAndServe(":8080", sessionManager.LoadAndSave(mux))

	log.Fatal(err)
}
