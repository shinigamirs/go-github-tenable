package controllers

import (
	"github.com/labstack/echo/v4"
	"go-github-tenable/githubAPI/infrastructure"
)

func AddLoginRoute(e *echo.Echo) {
	t := e.Group("")
	t.GET("login", infrastructure.GithubLoginController)
	t.GET("callback/handler", infrastructure.GithubCallbackHandlerController)
}
