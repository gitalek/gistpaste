package mysql

import (
	"database/sql"
	"github.com/gitalek/gistpaste/pkg/models"
)

type GistModel struct {
	DB *sql.DB
}

func (m GistModel) Insert(title, content, expires string) (int, error) {
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

func (m GistModel) Get(id string) (*models.Gist, error) {
	return nil, nil
}

// latest method returns the 10 most recently created gits.
func (m GistModel) latest() ([]*models.Gist, error) {
	return nil, nil
}
