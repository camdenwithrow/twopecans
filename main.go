package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/camdenwithrow/twopecans/config"
	"github.com/camdenwithrow/twopecans/db"
	"github.com/camdenwithrow/twopecans/handlers"
	"github.com/camdenwithrow/twopecans/services"
)

func main() {
	godotenv.Load()

	env := os.Getenv("ENVIRONMENT")

	if env != "dev" && env != "prod" {
		panic("NO ENVIRONMENT SET")
	}

	dbConfig := db.SqliteConfig{
		BaseUrl: os.Getenv("TURSO_DATABASE_URL"),
		Token:   os.Getenv("TURSO_AUTH_TOKEN"),
	}

	sqlDB, err := db.NewSqliteDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	sqlStore := db.NewSQLStore(sqlDB)

	sessionStore := services.NewCookieStore(services.SessionOptions{
		CookiesKey: config.Envs.CookiesAuthSecret,
		MaxAge:     config.Envs.CookiesAuthAgeInSeconds,
		Secure:     config.Envs.CookiesAuthIsSecure,
		HttpOnly:   config.Envs.CookiesAuthIsHttpOnly,
	})

	authService := services.NewAuthService(sessionStore)
	handler := handlers.New(env, sqlStore, authService)

	e := echo.New()

	e.Static("/js", "static/js")
	e.Static("/css", "static/css")
	e.Static("/img", "static/img")

	e.GET("/", handler.HomeHandler)
	e.GET("/recipes/:id", handler.GetRecipeHandler)

	// e.GET("/auth/:provider", handler.HandlerGoogleLogin)
	// e.GET("/auth/:provider/callback", handler.HandleGoogleCallback)
	// e.GET("/auth/logout/:provider", handler.HandleLogout)
	e.GET("/login", handler.HandleLogin) 

	e.Logger.Fatal(e.Start(":4321"))
}
