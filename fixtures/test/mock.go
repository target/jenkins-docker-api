package test

import "github.com/google/go-github/github"

func MockOrg(loginName string) *github.Organization {
	login := sPtr(loginName)
	return &github.Organization{
		Login: login,
	}
}

func MockTeam(teamName string, loginName string) *github.Team {
	name := sPtr(teamName)
	return &github.Team{
		Name:         name,
		Organization: MockOrg(loginName),
	}
}

func sPtr(s string) *string { return &s }
