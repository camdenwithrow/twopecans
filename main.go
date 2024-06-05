package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"github.com/camdenwithrow/twopecans/config"
	"github.com/camdenwithrow/twopecans/db"
	"github.com/camdenwithrow/twopecans/handlers"
	"github.com/camdenwithrow/twopecans/services"
)

const (
	DEV  = "development"
	TEST = "test"
	PROD = "production"
)

func main() {
	config.LoadConfig()

	sqlDB, err := db.NewSqliteDB(db.SqliteConfig{
		BaseUrl: config.DBUrl,
		Token:   config.DBAuthToken,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	store := db.NewSQLStore(sqlDB)

	sessionStore := services.NewCookieStore(services.CookieConfig)

	authService := services.NewAuthService(sessionStore)
	handler := handlers.New(config.Environment, store, authService)

	e := echo.New()

	e.Static("/js", "static/js")
	e.Static("/css", "static/css")
	e.Static("/img", "static/img")

	e.GET("/", handler.HomeHandler)
	e.GET("/recipes/:id", handler.GetRecipeHandler)

	e.GET("/auth/:provider", handler.HandleProviderLogin)
	e.GET("/auth/:provider/callback", handler.HandleAuthCallback)
	e.GET("/auth/logout/:provider", handler.HandleLogout)
	e.GET("/login", handler.HandleLogin)

	e.Logger.Fatal(e.Start(config.Port))
}
