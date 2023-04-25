package models

import (
	"database/sql"
	"errors"
	"time"
)

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
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId() // used to get value of "auto increment"
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	//*sql.Row type (always returns a sql.Row, errors are deflected untill row.Scan() is run)
	// in that case ErrNoRows is returned
	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{} // initialize a zeroed-struct reference (to be used to store the query)

	// row.Scan takes references to the locations where the returned row from the dbms
	// is to be written
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err // any other type of error
		}
	}
	return s, nil

}

// returns the latest 10 rows from the database
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// returns a resultset which can be traversed
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// (important to call), as rows do not get closed unless closed
	// can lead to resource hogging
	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		s := &Snippet{}

		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	// (important)have to call this to verify if the iteration is successfull
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

func (m *SnippetModel) ExampleTransaction() error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}
	//(important step) Rollback returns to previous commit even if any one
	// of the database operations are not executed correctly
	defer tx.Rollback()

	_, err = tx.Exec("INSERT....")
	if err != nil {
		return err
	}

	_, err = tx.Exec("Update...")
	if err != nil {
		return nil
	}

	err = tx.Commit()
	return err
}
