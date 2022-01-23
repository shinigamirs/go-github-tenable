package controllers

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"os"
)

func InitEcho() *echo.Echo {
	e := echo.New()
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
	return e
}

func AddRoutes(e *echo.Echo) {
	AddGithubRoute(e)
	AddLoginRoute(e)
}
