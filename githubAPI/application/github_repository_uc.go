package application

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"net/http"
)

//githubRepositoryListOutput struct for listing repository
type githubRepositoryListOutput struct {
	Name        string `json:"name"`
	FullName    string `json:"fullName"`
	Description string `json:"description"`
}

//createBranchParam struct for creating branch
type createBranchParam struct {
	RepoName              string `json:"repoName" validate:"required"`
	SourceBranchName      string `json:"sourceBranchName" default:"main"`
	DestinationBranchName string `json:"destinationBranchName" validate:"required"`
	Private               bool   `json:"private"`
}

//createPullRequestParam struct for creating pull request
type createPullRequestParam struct {
	RepoName  string `json:"repoName" validate:"required"`
	PrSubject string `json:"prSubject" validate:"required"`
	Head      string `json:"head" validate:"required"`
	Base      string `json:"base" validate:"required"`
}

//createFileContentParam struct for creating a file in github
type createFileContentParam struct {
	RepoName    string `json:"repoName" validate:"required"`
	BranchName  string `json:"branchName" validate:"required"`
	FileName    string `json:"fileName" validate:"required"`
	FileContent string `json:"fileContent" validate:"required"`
	Path        string `json:"path" validate:"required"`
	Message     string `json:"message" validate:"required"`
}

var ctx = context.Background()

// createGithubClient return a github client
func createGithubClient(c echo.Context) (*github.Client, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return nil, err
	}
	accessToken := sess.Values["access_token"]
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken.(string)},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc), nil
}

// GithubListRepository endpoint for listing repository
func GithubListRepository(c echo.Context) error {
	client, err := createGithubClient(c)
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
	var param createBranchParam
	client, err := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&param)
	repo := github.Repository{
		Name: &param.RepoName,
	}
	if err != nil {
		return err
	}
	repos, _, err := client.Repositories.Create(ctx, "", &repo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, repos)
}

// GithubGetRepo endpoint to get repo
func GithubGetRepo(c echo.Context) error {
	client, err := createGithubClient(c)
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

// GithubCreateBranch endpoint to create branch
func GithubCreateBranch(c echo.Context) error {
	var param createBranchParam
	client, err := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&param)
	if err != nil {
		return err
	}
	if param.SourceBranchName == "" {
		param.SourceBranchName = "master"
	}
	userName, _, _ := client.Users.Get(ctx, "")
	sourceBranchRefString := fmt.Sprintf("refs/heads/%s", param.SourceBranchName)
	sourceBranchRef, resp, err := client.Git.GetRef(ctx, *userName.Login, param.RepoName, sourceBranchRefString)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
		} else {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	}
	newBranchRefString := fmt.Sprintf("refs/heads/%s", param.DestinationBranchName)
	*sourceBranchRef.Ref = newBranchRefString
	newBranchRef, resp, err := client.Git.CreateRef(ctx, *userName.Login, param.RepoName, sourceBranchRef)
	if err != nil {
		log.Error(err)
		if resp.StatusCode == 404 {
			return echo.NewHTTPError(http.StatusNotFound, "Repo Not Found")
		} else {
			return echo.NewHTTPError(http.StatusOK, err.Error())
		}
	}
	log.Info(fmt.Sprintf("Branch %s Created Successfully", param.DestinationBranchName))
	return c.JSON(http.StatusOK, newBranchRef)
}

// CreateGithubPullRequest endpoint to create pull request
func CreateGithubPullRequest(c echo.Context) error {
	var param createPullRequestParam
	client, err := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&param)
	if err != nil {
		return err
	}
	createPullRequest := &github.NewPullRequest{
		Title: &param.PrSubject,
		Head:  &param.Head,
		Base:  &param.Base,
	}
	userName, _, _ := client.Users.Get(ctx, "")
	repo, _, err := client.PullRequests.Create(ctx, *userName.Login, param.RepoName, createPullRequest)
	if err != nil {
		log.Error(err.Error())
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	log.Info("PullRequest Created Successfully")
	return c.JSON(http.StatusOK, repo)
}

// CreateRepositoryContent creates a file with base64 content and commit using the message
func CreateRepositoryContent(c echo.Context) error {
	var createFileParam createFileContentParam
	client, err := createGithubClient(c)
	if err != nil {
		return err
	}
	err = c.Bind(&createFileParam)
	repositoryFileContent := &github.RepositoryContentFileOptions{
		Message: &createFileParam.Message,
		Content: []byte(createFileParam.FileContent),
		Branch:  &createFileParam.BranchName,
	}
	if createFileParam.Path == "" {
		createFileParam.Path = createFileParam.FileName
	}
	userName, _, _ := client.Users.Get(ctx, "")
	content, resp, err := client.Repositories.CreateFile(ctx, *userName.Login, createFileParam.RepoName, createFileParam.Path,
		repositoryFileContent)
	if err != nil {
		log.Error(err.Error())
		if resp.StatusCode == http.StatusNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else if resp.StatusCode == http.StatusConflict {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}
	}
	return c.JSON(http.StatusOK, content)
}
