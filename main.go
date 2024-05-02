package main

import (
	"net/http"
	"os"

	"github.com/a-h/templ"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"github.com/camdenwithrow/twopecans/db"
	"github.com/camdenwithrow/twopecans/views"
)

func main() {
	godotenv.Load()

	env := os.Getenv("ENVIRONMENT")

	if env != "dev" && env != "prod" {
		panic("NO ENVIRONMENT SET")
	}

	database := db.OpenDatabase()
	defer database.Close()

	e := echo.New()

	e.Static("/js", "static/js")
	e.Static("/css", "css")

	e.GET("/", HomeHandler(env))

	e.Logger.Fatal(e.Start(":4321"))
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	ctx.Response().Writer.WriteHeader(statusCode)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return t.Render(ctx.Request().Context(), ctx.Response().Writer)
}

func HomeHandler(env string) echo.HandlerFunc {
	return func(c echo.Context) error {
		return Render(c, http.StatusOK, views.Home(env))
	}
}
