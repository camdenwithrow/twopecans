package handlers

import (
	"net/http"

	"github.com/camdenwithrow/twopecans/views"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetRecipeHandler(c echo.Context) error {
	id := c.Param("id")
	return Render(c, http.StatusOK, views.Recipe(h.env, id))
}
