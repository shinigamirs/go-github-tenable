package domain

import (
	"errors"
	"log"
	"os"
)

// GetGithubClientSecret is to fetch GITHUB_CLIENT_SECRET env variable
func GetGithubClientSecret() (string, error) {

	githubClientSecret, exists := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !exists {
		log.Fatal("Github ClientSecret not defined in .env file")
		return "", errors.New("GITHUB_CLIENT_SECRET is not define")
	}

	return githubClientSecret, nil
}

// GetGithubClientID is to fetch GITHUB_CLIENT_ID env variable
func GetGithubClientID() (string, error) {

	githubClientID, exists := os.LookupEnv("GITHUB_CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
		return "", errors.New("GITHUB_CLIENT_ID is not define")
	}

	return githubClientID, nil
}
