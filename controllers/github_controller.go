package controllers

import (
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"go-github-tenable/githubAPI/infrastructure"
	"net/http"
)

func SessionAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, _ := session.Get("session", c)
			if sess.IsNew == true {
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			if sess.Values["access_token"].(string) != "" {
				return next(c)
			}
			return c.Redirect(http.StatusTemporaryRedirect, "/login")

		}
	}
}

func AddGithubRoute(e *echo.Echo) {
	t := e.Group("/")
	t.Use(SessionAuth())
	t.GET("github/listRepo", infrastructure.GitListRepositoryController)
	t.POST("github/createRepo", infrastructure.GithubCreateRepositoryController)
	t.GET("github/getRepo", infrastructure.GithubGetRepositoryController)
	t.POST("github/createBranch", infrastructure.GithubCreateBranchController)
	t.POST("github/createPullRequest", infrastructure.GithubCreatePullRequestController)
}
