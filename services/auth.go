package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/camdenwithrow/twopecans/config"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

type AuthService struct{}

func NewAuthService(store sessions.Store) *AuthService {
	gothic.Store = store

	goth.UseProviders(
		github.New(
			config.Envs.GithubClientID,
			config.Envs.GithubClientSecret,
			buildCallbackURL("github"),
		),
		google.New(
			config.Envs.GoogleClientID,
			config.Envs.GoogleClientSecret,
			buildCallbackURL("google"),
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

func RequireAuth(handlerFunc http.HandlerFunc, auth *AuthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := auth.GetSessionUser(r)
		if err != nil {
			log.Println("User is not authenticated!")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}

		log.Printf("user is authenticated! user: %v!", session.FirstName)

		handlerFunc(w, r)
	}
}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s:%s/auth/%s/callback", config.Envs.PublicHost, config.Envs.Port, provider)
}
