package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `insert into snippets(title, content, created, expires)
		values(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `select id, title, content, created, expires from snippets where id = ? and expires > utc_timestamp()`
	snippet := &Snippet{}

	err := m.DB.QueryRow(stmt, id).Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// multi record sql query
	stmt := `select id, title, content, created, expires from snippets where expires > utc_timestamp() order by id desc limit 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	s := []*Snippet{}
	for rows.Next() {
		snippet := &Snippet{}
		err := rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		s = append(s, snippet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return s, nil
}

func (m *SnippetModel) ExampleTransaction() error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("insert into ..")
	if err != nil {
		return err
	}

	_, err = tx.Exec("update ...")
	if err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

// example Pre statement
type PreStmtModel struct {
	DB         *sql.DB
	InsertStmt *sql.Stmt
}

func NewPreStmtModel(db *sql.DB) (*PreStmtModel, error) {
	insertStmt, err := db.Prepare("insert into ...")
	if err != nil {
		return nil, err
	}

	return &PreStmtModel{
		DB:         db,
		InsertStmt: insertStmt,
	}, nil

}

func (m *PreStmtModel) insert(args ...interface{}) error {
	_, err := m.InsertStmt.Exec(args...)
	return err
}
