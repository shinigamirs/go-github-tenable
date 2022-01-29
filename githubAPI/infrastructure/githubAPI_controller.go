package infrastructure

import (
	"github.com/labstack/echo/v4"
	"go-github-tenable/githubAPI/application"
)

// GithubLoginController is the controller for GitHub login
func GithubLoginController(c echo.Context) error {
	return application.GithubLogin(c)
}

// GitListRepositoryController is the controller for GitHub list repo
func GitListRepositoryController(c echo.Context) error {
	return application.GithubListRepository(c)
}

// GithubCallbackHandlerController is the controller for call back handler
func GithubCallbackHandlerController(c echo.Context) error {
	return application.GithubLoginCallbackHandler(c)
}

// GithubCreateRepositoryController is the controller to create repository from GitHub
func GithubCreateRepositoryController(c echo.Context) error {
	return application.GithubCreateRepo(c)
}

// GithubGetRepositoryController is the controller to get repository from GitHub
func GithubGetRepositoryController(c echo.Context) error {
	return application.GithubGetRepo(c)
}

// GithubCreateBranchController is the controller to create branch in GitHub repository
func GithubCreateBranchController(c echo.Context) error {
	return application.GithubCreateBranch(c)
}

// GithubCreatePullRequestController is the controller to create pull request in GitHub
func GithubCreatePullRequestController(c echo.Context) error {
	return application.CreateGithubPullRequest(c)
}

// GithubCreateContentRequestController is the controller to create file in GitHub repository
func GithubCreateContentRequestController(c echo.Context) error {
	return application.CreateRepositoryContent(c)
}
