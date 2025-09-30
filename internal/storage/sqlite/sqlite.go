package sqlite

import (
	"database/sql"
	"errors"
	
	"github.com/Magic-B/url-shortener/internal/storage"
	"github.com/Magic-B/url-shortener/pkg/apperr"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, apperr.ErrWrapper(op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`)

	if err != nil {
		return nil, apperr.ErrWrapper(op, err)
	}

	if _, err = stmt.Exec(); err != nil {
		return nil, apperr.ErrWrapper(op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlite.SaveURL"

	smtp, err := s.db.Prepare("INSERT INTO url (url, alias) VALUES (?, ?)")

	if err != nil {
		return 0, apperr.ErrWrapper(op, err)
	}
	
	res, err := smtp.Exec(urlToSave, alias)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, apperr.ErrWrapper(op, storage.ErrURLExist)
		}

		return 0, apperr.ErrWrapper(op, err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, apperr.ErrWrapper(op, err, "failed to get last insert id")
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlite.GetUrl"
	
	smtp, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")

	if err != nil {
		return "", apperr.ErrWrapper(op, err)
	}

	var res string

	err = smtp.QueryRow(alias).Scan(&res)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperr.ErrWrapper(op, storage.ErrURLNotFound)
		}

		return "", apperr.ErrWrapper(op, err)
	}

	return res, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlite.DeleteURL"
	
	smtp, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")

	if err != nil {
		return apperr.ErrWrapper(op, err)
	}
	
	if _, err = smtp.Exec(alias); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return apperr.ErrWrapper(op, storage.ErrURLNotFound)
		}

		return apperr.ErrWrapper(op, err)
	}

	return nil
}
