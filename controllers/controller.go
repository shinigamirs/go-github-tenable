package controllers

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"os"
)

func InitEcho() *echo.Echo {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	return e
}

func AddRoutes(e *echo.Echo) {
	AddGithubRoute(e)
}
