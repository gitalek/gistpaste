package main

import "github.com/gitalek/gistpaste/pkg/models"

type templateData struct {
	Gist  *models.Gist
	Gists []*models.Gist
}
