package handlers

import (
	"net/http"

	"github.com/camdenwithrow/twopecans/views"
	"github.com/labstack/echo/v4"
)

func (h *Handler) HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, views.Home(h.env))
}
