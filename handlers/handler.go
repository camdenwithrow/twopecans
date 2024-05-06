package handlers

import (
	"github.com/camdenwithrow/twopecans/db"
	"github.com/camdenwithrow/twopecans/services"
)

type Handler struct {
	store db.Store
	env   string
	auth  *services.AuthService
}

func New(env string, store db.Store, auth *services.AuthService) *Handler {
	return &Handler{
		store: store,
		env:   env,
		auth:  auth,
	}
}
