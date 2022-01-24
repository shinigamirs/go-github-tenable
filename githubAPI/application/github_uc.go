package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"go-github-tenable/githubAPI/domain"
	"io/ioutil"
	"net/http"
	"os"
)

// getGithubAccessToken is to get the github access token
func getGithubAccessToken(code string) string {

	clientID, _ := domain.GetGithubClientID()
	clientSecret, _ := domain.GetGithubClientSecret()

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	// request to github to get accessToken
	req, reqErr := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON))
	if reqErr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		log.Panic("Github Request failed")
	}

	respBody, _ := ioutil.ReadAll(resp.Body)

	// Represents the response received from Github
	type githubAccessTokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	var githubResponse githubAccessTokenResponse
	json.Unmarshal(respBody, &githubResponse)

	return githubResponse.AccessToken
}

// GithubLogin is to login in github using oauth
func GithubLogin(c echo.Context) error {
	githubClientID, err := domain.GetGithubClientID()
	if err != nil {
		return err
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create the dynamic redirect URL for login
	redirectURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&scope=repo&redirect_uri="+
			"http://%s:%s/callback/handler", githubClientID, host, port,
	)
	log.Info("Redirect URL", redirectURL)
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GithubLoginCallbackHandler to store access_token in session
func GithubLoginCallbackHandler(c echo.Context) error {
	// To get the code out of github login url
	code := c.QueryParam("code")
	githubAccessToken := getGithubAccessToken(code)
	sess, _ := session.Get("session", c)
	sess.Values["access_token"] = githubAccessToken
	sess.Save(c.Request(), c.Response())
	log.Info("Call back handler done")
	return c.Redirect(http.StatusTemporaryRedirect, "/github/listRepo")
}
