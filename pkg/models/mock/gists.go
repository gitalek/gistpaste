package mock

import (
"github.com/gitalek/gistpaste/pkg/models"
"time"
)

var mockGist = &models.Gist{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type GistModel struct {}

func (m *GistModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *GistModel) Get(id int) (*models.Gist, error) {
	switch id {
	case 1:
		return mockGist, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *GistModel) Latest() ([]*models.Gist, error) {
	return []*models.Gist{mockGist}, nil
}
