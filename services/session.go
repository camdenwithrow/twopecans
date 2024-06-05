package services

import (
	"github.com/camdenwithrow/twopecans/config"
	"github.com/gorilla/sessions"
)

const (
	SessionName = "session"
)

type SessionOptions struct {
	CookiesKey string
	MaxAge     int
	HttpOnly   bool
	Secure     bool // Should be true if the site is served over HTTPS (prod)
}

var CookieConfig = SessionOptions{
	CookiesKey: config.CookiesAuthSecret,
	MaxAge:     60 * 60 * 24 * 30,
	HttpOnly:   true,
	Secure:     config.Environment == "production",
}

func NewCookieStore(opts SessionOptions) *sessions.CookieStore {
	store := sessions.NewCookieStore([]byte(opts.CookiesKey))

	store.MaxAge(opts.MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = opts.HttpOnly
	store.Options.Secure = opts.Secure

	return store
}
