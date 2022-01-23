package controllers

import (
	"github.com/labstack/echo"
	infrastructure "go-git-tenable/githubAPI/infrastructure"
)

func AddGithubRoute(e *echo.Echo) {
	e.GET("/login", infrastructure.GithubLoginController)
	e.GET("/callback/handler", infrastructure.GithubCallbackHandlerController)
	e.GET("/github/listRepo", infrastructure.GitListRepositoryController)
	e.POST("/github/createRepo", infrastructure.GithubCreateRepositoryController)
	e.GET("/github/getRepo", infrastructure.GithubGetRepositoryController)
	e.POST("/github/createBranch", infrastructure.GithubCreateBranchController)
}
