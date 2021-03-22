package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Store interface {
	Open() error
	Close() error

	GetMovieById(id int64) (*Movie, error)
	GetMovies() ([]*Movie, error)
	CreateMovie(m *Movie) error

	FindUser(username, password string) (bool, error)
}

type dbStore struct {
	db *sqlx.DB
}

var schema = `
CREATE TABLE IF NOT EXISTS movie
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT,
	release_date TEXT,
	duration INTEGER,
	trailer_url TEXT
);

CREATE TABLE IF NOT EXISTS user
(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user TEXT,
	password TEXT
);
`

func (store *dbStore) Open() error {
	db, err := sqlx.Connect("sqlite3", "goflix.db")
	if err != nil {
		return err
	}
	log.Println("Connected to DB!")
	db.MustExec(schema)
	store.db = db
	return nil
}

func (store *dbStore) Close() error {
	if err := store.db.Close(); err != nil {
		return err
	}
	return nil
}

func (store *dbStore) GetMovies() ([]*Movie, error) {
	var movies []*Movie
	if err := store.db.Select(&movies, "SELECT * FROM movie"); err != nil {
		return movies, err
	}
	return movies, nil
}

func (store *dbStore) GetMovieById(id int64) (*Movie, error) {
	var movie = &Movie{}
	fmt.Println(id)
	if err := store.db.Get(movie, "SELECT * FROM movie WHERE id=$1", id); err != nil {
		return movie, err
	}
	return movie, nil
}

func (store *dbStore) CreateMovie(m *Movie) error {
	res, err := store.db.Exec("INSERT INTO movie (title, release_date, duration, trailer_url) VALUES (?, ?, ?, ?)",
		m.Title, m.ReleaseDate, m.Duration, m.TrailerURL)

	if err != nil {
		return err
	}

	m.ID, err = res.LastInsertId()
	return err
}

func (store *dbStore) FindUser(username, password string) (bool, error) {
	var count int
	if err := store.db.Get(&count, "SELECT COUNT(id) FROM user WHERE user=$1 AND password=$2", username, password); err != nil {
		return false, err
	}

	return count == 1, nil
}
