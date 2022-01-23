package infrastructure

import (
	"github.com/labstack/echo"
	"go-github-tenable/githubAPI/application"
)

func GithubLoginController(c echo.Context) error {
	return application.GithubLogin(c)
}

func GitListRepositoryController(c echo.Context) error {
	return application.GithubListRepository(c)
}

func GithubCallbackHandlerController(c echo.Context) error {
	return application.GithubLoginCallbackHandler(c)
}

func GithubCreateRepositoryController(c echo.Context) error {
	return application.GithubCreateRepo(c)
}

func GithubGetRepositoryController(c echo.Context) error {
	return application.GithubGetRepo(c)
}

func GithubCreateBranchController(c echo.Context) error {
	return application.GithubCreateBranch(c)
}

func GithubCreatePullRequestController(c echo.Context) error {
	return application.CreateGithubPullRequest(c)
}
