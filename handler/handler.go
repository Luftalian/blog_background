package handler

import (
	"blog-backend/model"
)

type Handler struct {
	Repo   *model.Repository
	Config *model.Configuration
}

func New(repo *model.Repository, config *model.Configuration) *Handler {
	return &Handler{
		Repo:   repo,
		Config: config,
	}
}
