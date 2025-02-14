package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	appFolderName = "ivy_lee_todo"
	dbFileName    = "data.db"
)

type Db struct {
	internal *sql.DB
}

// NewSqlite3Db creates a new SQLite3 database.
//
// When `initDb` variable is true, then the database file will be always deleted and recreated.
func NewSqlite3Db(initialize bool) *Db {
	log.Println("Initializing SQLite3 database")
	p, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	p = filepath.Join(p, appFolderName)
	if _, err := os.Stat(p); err != nil {
		log.Println("Creating app dir in a cache folder")
		if err = os.Mkdir(p, 0755); err != nil {
			log.Fatal(err)
		}
	}
	p = filepath.Join(p, dbFileName)
	if initialize {
		log.Println("Deleting a database file")
		os.Remove(p)
	}
	log.Println("Loading a database file in the app folder")
	db, err := sql.Open("sqlite3", p)
	if err != nil {
		log.Fatal(err)
	}
	initSchema(db)
	return &Db{internal: db}
}

func initSchema(db *sql.DB) {
	log.Println("Creating tables if not exists")
	stmt := `
CREATE TABLE IF NOT EXISTS task (
	id INTEGER PRIMARY KEY AUTOINCREMENT, -- ID
	title TEXT NOT NULL, -- Title
	description TEXT NOT NULL DEFAULT '', -- Description
	actionable INTEGER NOT NULL DEFAULT 0 CHECK (actionable IN (0 ,1)), -- A boolean type indicating whether the task is actionable
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- Store as DATETIME
	deleted_at DATETIME -- This has a value if the task is deleted
)`
	if _, err := db.Exec(stmt); err != nil {
		log.Fatal(err)
	}
}

func (db *Db) Begin() error {
	return db.Begin()
}

func (db *Db) Commit() error {
	return db.Commit()
}

func (db *Db) Exec(query string, args ...any) (sql.Result, error) {
	return db.internal.Exec(query, args...)
}

func (db *Db) Query(query string, args ...any) (*sql.Rows, error) {
	return db.internal.Query(query, args...)
}

func (db *Db) QueryRow(query string, args ...any) *sql.Row {
	return db.internal.QueryRow(query, args...)
}
