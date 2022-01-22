package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	domain "go-git-tenable/githubAPI/domain"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func getGithubAccessToken(code string) string {

	clientID, _ := domain.GetGithubClientID()
	clientSecret, _ := domain.GetGithubClientSecret()

	requestBodyMap := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code}
	requestJSON, _ := json.Marshal(requestBodyMap)

	req, reqErr := http.NewRequest("POST", "https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON))
	if reqErr != nil {
		log.Panic("Request creation failed")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		log.Panic("Request failed")
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

func GithubLogin(c echo.Context) error {
	// Get the environment variable
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
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri="+
			"http://%s:%s/callback/handler", githubClientID, host, port,
	)
	log.Println("Redirect URL", redirectURL)
	return c.Redirect(http.StatusPermanentRedirect, redirectURL)
}

func GithubLoginCallbackHandler(c echo.Context) error {
	code := c.QueryParam("code")
	githubAccessToken := getGithubAccessToken(code)
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	}
	log.Println("Access_Token", githubAccessToken)
	sess.Values["access_token"] = githubAccessToken
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusTemporaryRedirect, "/github/listRepo")
}
