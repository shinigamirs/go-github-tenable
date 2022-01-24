package infrastructure

import (
	"github.com/labstack/echo/v4"
	"go-github-tenable/githubAPI/application"
)

// GithubLoginController is the controller for github login
func GithubLoginController(c echo.Context) error {
	return application.GithubLogin(c)
}

// GitListRepositoryController is the controller for github list repo
func GitListRepositoryController(c echo.Context) error {
	return application.GithubListRepository(c)
}

// GithubCallbackHandlerController is the controller for call back handler
func GithubCallbackHandlerController(c echo.Context) error {
	return application.GithubLoginCallbackHandler(c)
}

// GithubCreateRepositoryController is the controller to create repository from Github
func GithubCreateRepositoryController(c echo.Context) error {
	return application.GithubCreateRepo(c)
}

// GithubGetRepositoryController is the controller to get repository from Github
func GithubGetRepositoryController(c echo.Context) error {
	return application.GithubGetRepo(c)
}

// GithubCreateBranchController is the controller to create branch in Github repository
func GithubCreateBranchController(c echo.Context) error {
	return application.GithubCreateBranch(c)
}

// GithubCreatePullRequestController is the controller to create pull request in Github
func GithubCreatePullRequestController(c echo.Context) error {
	return application.CreateGithubPullRequest(c)
}
