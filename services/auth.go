package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"

	"github.com/camdenwithrow/twopecans/config"
)

type AuthService struct{}

var cfg = config.GetConfig()

func NewAuthService(store sessions.Store) *AuthService {
	gothic.Store = store

	goth.UseProviders(
		google.New(
			cfg.GoogleClientID,
			cfg.GoogleClientSecret,
			buildCallbackURL("google"),
		),
		github.New(
			cfg.GithubClientID,
			cfg.GoogleClientSecret,
			buildCallbackURL("github"),
		),
	)

	return &AuthService{}
}

func (s *AuthService) GetSessionUser(r *http.Request) (goth.User, error) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		return goth.User{}, err
	}

	u := session.Values["user"]
	if u == nil {
		return goth.User{}, fmt.Errorf("user is not authenticated! %v", u)
	}

	return u.(goth.User), nil
}

func (s *AuthService) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, _ := gothic.Store.Get(r, SessionName)

	session.Values["user"] = user

	err := session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RemoveUserSession(w http.ResponseWriter, r *http.Request) {
	session, err := gothic.Store.Get(r, SessionName)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["user"] = goth.User{}
	// delete the cookie immediately
	session.Options.MaxAge = -1

	session.Save(r, w)
}

func RequireAuth(handlerFunc echo.HandlerFunc, auth *AuthService) echo.HandlerFunc {
	return func(c echo.Context) error {
		w, r := c.Response().Writer, c.Request()
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return nil
		}

		log.Printf("user is authenticated! user: %v!", session.FirstName)

		handlerFunc(c)
		return nil
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s:%s/auth/%s/callback", cfg.PublicHost, cfg.Port, provider)
}
