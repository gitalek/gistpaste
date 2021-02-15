package mysql

import (
	"database/sql"
	"errors"
	"github.com/gitalek/gistpaste/pkg/models"
)

type GistModel struct {
	DB *sql.DB
}

func (m *GistModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO gists (title, content, created, expires)
				VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

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

func (m *GistModel) Get(id string) (*models.Gist, error) {
	stmt := `SELECT id, title, content, created, expires FROM gists
				WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	g := new(models.Gist)
	err := row.Scan(&g.ID, &g.Title, &g.Content, &g.Created, &g.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return g, nil
}

// latest method returns the 10 most recently created gits.
func (m *GistModel) latest() ([]*models.Gist, error) {
	return nil, nil
}
