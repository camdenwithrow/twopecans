package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/camdenwithrow/twopecans/db"
	"github.com/camdenwithrow/twopecans/handlers"
)

func main() {
	godotenv.Load()

	env := os.Getenv("ENVIRONMENT")

	if env != "dev" && env != "prod" {
		panic("NO ENVIRONMENT SET")
	}

	config := db.SqliteConfig{
		BaseUrl: os.Getenv("TURSO_DATABASE_URL"),
		Token:   os.Getenv("TURSO_AUTH_TOKEN"),
	}

	sqlDB, err := db.NewSqliteDB(config)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	sqlStore := db.NewSQLStore(sqlDB)

	handler := handlers.New(env, sqlStore)

	e := echo.New()

	e.Static("/js", "static/js")
	e.Static("/css", "static/css")
	e.Static("/img", "static/img")

	e.GET("/", handler.HomeHandler)
	e.GET("/recipes/:id", handler.GetRecipeHandler)

	e.Logger.Fatal(e.Start(":4321"))
}
