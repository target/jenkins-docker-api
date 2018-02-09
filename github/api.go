package client

import (
	"context"
	"net/url"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubAPIURL = os.Getenv("GITHUB_API_URL")

// GetToken retrieves the oauth token from request header
func GetToken(c *gin.Context) string {
	log.Info("Obtaining oauth token from header")

	header := c.Request.Header["Authorization"]
	if len(header) == 0 {
		return ""
	}

	headerVal := header[0]
	if !strings.Contains(headerVal, "token ") {
		return ""
	}

	token := strings.Split(headerVal, " ")[1]
	return token
}

// GetUserName gets the Git username of a user from git
func GetUserName(token string) string {
	log.Info("Grabbing GitHub username")

	if githubAPIURL == "" {
		githubAPIURL = os.Getenv("GITHUB_API_URL")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	client.BaseURL, _ = url.Parse(githubAPIURL)

	user, _, _ := client.Users.Get(ctx, "")

	return *user.Login
}

// GetUserTeams gets the users teams in Github
func GetUserTeams(token string) ([]*github.Team, error) {
	log.Info("Authorizing user")

	if githubAPIURL == "" {
		githubAPIURL = os.Getenv("GITHUB_API_URL")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	client.BaseURL, _ = url.Parse(githubAPIURL)

	// list all teams for the authenticated user
	teams, _, err := client.Organizations.ListUserTeams(ctx, &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	return teams, nil
}
