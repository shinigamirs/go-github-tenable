package application

import (
	"context"
	"github.com/google/go-github/github"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"net/http"
)

type githubRepositoryListOutput struct {
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
}

type createBranchParam struct {
	Name    string `json:"name" validate:"required"`
	Private bool   `json:"private"`
}

// createGithubClient return a github client
func createGithubClient(c echo.Context) (*github.Client, error, context.Context) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err, nil
	}
	// creates an empty context
	ctx := context.Background()
	accessToken := sess.Values["access_token"]
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken.(string)},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil, ctx
}

// GithubListRepository endpoint for listing repository
func GithubListRepository(c echo.Context) error {
	client, err, ctx := createGithubClient(c)
	if err != nil {
		return err
	}
	repos, _, err := client.Repositories.List(ctx, "", nil)
	repoList := make([]githubRepositoryListOutput, len(repos))
	i := 0
	for j := range repos {
		repoList[i].Name = repos[j].GetName()
		repoList[i].FullName = repos[j].GetFullName()
		repoList[i].Description = repos[j].GetDescription()
		i++
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, repoList)
}

// GithubCreateRepo create a repo under auth user
func GithubCreateRepo(c echo.Context) error {
	var repo github.Repository
	client, err, ctx := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&repo)
	if err != nil {
		return err
	}
	repos, _, err := client.Repositories.Create(ctx, "", &repo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, repos)
}

func GithubGetRepo(c echo.Context) error {
	client, err, ctx := createGithubClient(c)
	if err != nil {
		return err
	}
	repoName := c.QueryParam("name")
	userName, _, _ := client.Users.Get(ctx, "")
	repo, resp, err := client.Repositories.Get(ctx, *userName.Login, repoName)
	if err != nil && resp.StatusCode == 404 {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
	}
	return c.JSON(http.StatusOK, repo)
}

//func GithubCreateBranch(c echo.Context) error {
//	var param *createBranchParam
//	client, err, ctx := createGithubClient(c)
//	if err != nil {
//		return err
//	}
//	err = c.Bind(&param)
//	if err != nil {
//		return err
//	}
//	userName, _, _ := client.Users.Get(ctx, "")
//	repo, resp, err := client.Repositories.Create(ctx, *userName.Login, repoName)
//	if err != nil && resp.StatusCode == 404 {
//		log.Error(err)
//		return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
//	}
//	return c.JSON(http.StatusOK, repo)
//}
