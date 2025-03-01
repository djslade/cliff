package main

import "github.com/djslade/cliff/internal/models"

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
