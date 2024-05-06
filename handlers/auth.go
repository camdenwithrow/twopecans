package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/camdenwithrow/twopecans/views"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func (h *Handler) HandleLogin(c echo.Context) error {
	return Render(c, http.StatusOK, views.Login(h.env))
}

func (h *Handler) HandleProviderLogin(c echo.Context) {
	// try to get the user without re-authenticating
	if u, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request()); err == nil {
		log.Printf("User already authenticated! %v", u)

		views.Login(h.env)
	} else {
		gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	}
}

func (h *Handler) HandleAuthCallbackFunction(c echo.Context) {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		fmt.Fprintln(c.Response().Writer, err)
		return
	}

	err = h.auth.StoreUserSession(c.Response().Writer, c.Request(), user)
	if err != nil {
		log.Println(err)
		return
	}

	c.Response().Writer.Header().Set("Location", "/")
	c.Response().Writer.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) HandleLogout(c echo.Context) {
	log.Println("Logging out...")

	err := gothic.Logout(c.Response().Writer, c.Request())
	if err != nil {
		log.Println(err)
		return
	}

	h.auth.RemoveUserSession(c.Response().Writer, c.Request())

	c.Response().Writer.Header().Set("Location", "/")
	c.Response().Writer.WriteHeader(http.StatusTemporaryRedirect)
}
