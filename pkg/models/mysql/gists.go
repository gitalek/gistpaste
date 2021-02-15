package mysql

import (
	"database/sql"
	"github.com/gitalek/gistpaste/pkg/models"
)

type GistModel struct {
	DB *sql.DB
}

func (m GistModel) Insert(title, content, expires string) (int, error) {
	return 0, nil
}

func (m GistModel) Get(id string) (*models.Gist, error) {
	return nil, nil
}

// latest method returns the 10 most recently created gits.
func (m GistModel) latest() ([]*models.Gist, error) {
	return nil, nil
}
