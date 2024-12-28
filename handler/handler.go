package handler

import (
	"blog-backend/model"

	"google.golang.org/api/drive/v3"
)

type Handler struct {
	Repo         *model.Repository
	Config       *model.Configuration
	DriveService *drive.Service
}

func New(repo *model.Repository, config *model.Configuration, srv *drive.Service) *Handler {
	return &Handler{
		Repo:         repo,
		Config:       config,
		DriveService: srv,
	}
}
