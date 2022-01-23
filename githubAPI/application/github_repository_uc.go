package application

import (
	"context"
	"fmt"
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
	repoName              string `json:"repoName" validate:"required"`
	sourceBranchName      string `json:"sourceBranchName" default:"main"`
	destinationBranchName string `json:"destinationBranchName" validate:"required"`
	private               bool   `json:"private"`
}

type createPullRequestParam struct {
	repoName  string  `json:"repoName" validate:"required"`
	prSubject string  `json:"prSubject" validate:"required"`
	head      *string `json:"head" validate:"required"`
	base      *string `json:"base" validate:"required"`
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
	// TODO get branch and add to output list
	//repo, resp, err := client.Repositories.GetBranch(ctx, userName, param.Name)
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
	if repoName == "" {
		return echo.NewHTTPError(http.StatusNotFound, "Please specify a repo")
	}
	userName, _, _ := client.Users.Get(ctx, "")
	repo, resp, err := client.Repositories.Get(ctx, *userName.Login, repoName)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
		} else {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	}
	return c.JSON(http.StatusOK, repo)
}

func GithubCreateBranch(c echo.Context) error {
	var param createBranchParam
	client, err, ctx := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&param)
	if err != nil {
		return err
	}
	if param.sourceBranchName == "" {
		param.sourceBranchName = "master"
	}
	userName, _, _ := client.Users.Get(ctx, "")
	sourceBranchRefString := fmt.Sprintf("refs/heads/%s", param.sourceBranchName)
	sourceBranchRef, resp, err := client.Git.GetRef(ctx, *userName.Login, param.repoName, sourceBranchRefString)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
		} else {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	}
	newBranchRefString := fmt.Sprintf("refs/heads/%s", param.destinationBranchName)
	*sourceBranchRef.Ref = newBranchRefString
	newBranchRef, resp, err := client.Git.CreateRef(ctx, *userName.Login, param.repoName, sourceBranchRef)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
		} else {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	}
	log.Info(fmt.Sprintf("Branch %s Created Successfully", param.destinationBranchName))
	return c.JSON(http.StatusOK, newBranchRef)
}

func CreateGithubPullRequest(c echo.Context) error {
	var param createPullRequestParam
	client, err, ctx := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&param)
	if err != nil {
		return err
	}
	createPullRequest := &github.NewPullRequest{
		Head: param.head,
		Base: param.base,
	}
	userName, _, _ := client.Users.Get(ctx, "")
	repo, _, err := client.PullRequests.Create(ctx, *userName.Login, param.repoName, createPullRequest)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	log.Info("PullRequest Created Successfully")
	return c.JSON(http.StatusOK, repo)
}
