package domain

import (
	"errors"
	"log"
	"os"
)

func GetGithubClientSecret() (string, error) {

	githubClientSecret, exists := os.LookupEnv("GITHUB_CLIENT_SECRET")
	if !exists {
		log.Fatal("Github ClientSecret not defined in .env file")
		return "", errors.New("Github ClientSecret is not define")
	}

	return githubClientSecret, nil
}

func GetGithubClientID() (string, error) {

	githubClientID, exists := os.LookupEnv("GITHUB_CLIENT_ID")
	if !exists {
		log.Fatal("Github Client ID not defined in .env file")
		return "", errors.New("Github ClientID is not define")
	}

	return githubClientID, nil
}
