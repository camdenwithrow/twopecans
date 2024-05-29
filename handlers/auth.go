package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"

	"github.com/camdenwithrow/twopecans/views"
)

type ProviderContextKey struct{}

func (h *Handler) HandleLogin(c echo.Context) error {
	return Render(c, http.StatusOK, views.Login(h.env))
}

func (h *Handler) HandleProviderLogin(c echo.Context) error {
	provider := c.Param("provider")
	r := c.Request().WithContext(context.WithValue(context.Background(), gothic.ProviderParamKey, provider))
	// if p, ok := req.Context().Value(ProviderParamKey).(string); ok {

	// try to get the user without re-authenticating
	if u, err := gothic.CompleteUserAuth(c.Response().Writer, r); err == nil {
		log.Printf("User already authenticated! %v", u)

		views.Login(h.env)
	} else {
		gothic.BeginAuthHandler(c.Response().Writer, c.Request())
	}
	return nil
}

func (h *Handler) HandleAuthCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response().Writer, c.Request())
	if err != nil {
		fmt.Fprintln(c.Response().Writer, err)
		return err
	}

	err = h.auth.StoreUserSession(c, user)
	if err != nil {
		log.Println(err)
		return err
	}

	c.Response().Writer.Header().Set("Location", "/")
	c.Response().Writer.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (h *Handler) HandleLogout(c echo.Context) error {
	log.Println("Logging out...")

	err := gothic.Logout(c.Response().Writer, c.Request())
	if err != nil {
		log.Println(err)
		return err
	}

	h.auth.RemoveUserSession(c.Response().Writer, c.Request())

	c.Response().Writer.Header().Set("Location", "/")
	c.Response().Writer.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}
