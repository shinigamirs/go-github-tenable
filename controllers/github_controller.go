package controllers

import (
	"github.com/labstack/echo/v4"
	"go-github-tenable/githubAPI/infrastructure"
)

func AddGithubRoute(e *echo.Echo) {
	t := e.Group("/")
	t.GET("github/listRepo", infrastructure.GitListRepositoryController)
	t.POST("github/createRepo", infrastructure.GithubCreateRepositoryController)
	t.GET("github/getRepo", infrastructure.GithubGetRepositoryController)
	t.POST("github/createBranch", infrastructure.GithubCreateBranchController)
	t.POST("github/createPullRequest", infrastructure.GithubCreatePullRequestController)
}
