package controllers

import (
	"github.com/labstack/echo"
	infrastructure "go-github-tenable/githubAPI/infrastructure"
)

func AddLoginRoute(e *echo.Echo) {
	t := e.Group("")
	t.GET("login", infrastructure.GithubLoginController)
	t.GET("callback/handler", infrastructure.GithubCallbackHandlerController)
}
