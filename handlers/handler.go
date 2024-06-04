package handlers

import (
	"github.com/camdenwithrow/twopecans/db"
)

type Handler struct {
	store db.Store
	env   string
}

func New(env string, store db.Store) *Handler {
	return &Handler{
		store: store,
		env:   env,
	}
}
