package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"url-shortener/internal/storage"

	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(host, port, user, pass, name string) (*Storage, error) {
	psqlConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, name)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", psqlConn, err)
	}
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS public.url(
		id INTEGER PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
		alias VARCHAR NOT NULL UNIQUE,
		url VARCHAR NOT NULL);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", psqlConn, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", psqlConn, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", psqlConn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	stmt := "INSERT INTO url(url, alias) VALUES($1, $2)"
	_, err := s.db.Exec(stmt, urlToSave, alias)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				fmt.Println(pqErr.Code)
				return fmt.Errorf("INSERT %s", storage.ErrUrlExists)
			}
		}
		return fmt.Errorf("INSERT %w", err)
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	stmt := "SELECT u.url FROM url u WHERE u.alias=$1"
	var queryRes string
	err := s.db.QueryRow(stmt, alias).Scan(&queryRes)
	if errors.Is(err, sql.ErrNoRows) {
		return "", storage.ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("SELECT executing %w", err)
	}
	return queryRes, nil
}

func (s *Storage) DeleteURL(alias string) error {
	stmt := "DELETE FROM url WHERE url.url=$1"
	res, err := s.db.Exec(stmt, alias)
	if err != nil {
		return fmt.Errorf("DELETE %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil || rows == 0 {
		return fmt.Errorf("DELETE %w", err)
	}
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
