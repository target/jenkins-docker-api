package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/target/gelvedere/model"
	ghClient "github.com/target/jenkins-docker-api/github"
	"github.com/target/jenkins-docker-api/util"
)

var (
	adminOrg  string
	adminTeam string
)

// authenticateUser calls the github client to authenticate user
func authenticateUser(c *gin.Context, serviceName string, adminMap map[string][]string) error {
	// Verify if user is authorized
	token := ghClient.GetToken(c)
	if token == "" {
		return fmt.Errorf("no authentication token was found in the request")
	}

	// Get user's teams
	teams, _ := ghClient.GetUserTeams(token)

	if adminMap == nil {
		// Get map of admins of Jenkins master
		usrConfig := os.Getenv("JENKINS_USER_CONFIG_PATH")
		adminMap, _ = getAdminMap(serviceName, usrConfig)
	}

	// Validate if teams exist in adminMap
	if isAuthorizedUser(teams, adminMap) {
		user := ghClient.GetUserName(token)
		log.Info(fmt.Sprintf("Successfully authorized %s", user))
	} else {
		err := fmt.Errorf("the provided credentials are unauthorized")
		log.Error(err)
		return err
	}

	return nil
}

// isAuthorizedUser verifies if the Jenkins master name is part of the adminMap
func isAuthorizedUser(teams []*github.Team, adminMap map[string][]string) bool {
	if adminMap == nil {
		log.Error("Admin map is nil")
		return false
	}
	// Verify if user is in team
	for _, team := range teams {
		orgName := team.Organization.GetLogin()
		teamName := team.GetName()
		if orgTeams, ok := adminMap[orgName]; ok {
			if util.SliceContains(teamName, orgTeams) {
				return true
			}
		}
	}
	return false
}

// getAdminMap parses Jenkins master info from the user config directory
func getAdminMap(master string, configPath string) (map[string][]string, error) {
	// Get JSON contents
	configPath = util.TrimSuffix(configPath, "/") // remove trailing "/"
	env := os.Getenv("JENKINS_ENV")               // env = test or prod
	filename := fmt.Sprintf("%s/%s/%s.json", configPath, env, master)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// Unmarshal JSON contents to UserConfig
	var uc model.UserConfig
	json.Unmarshal(bytes, &uc)

	// Create map from user config
	adminMap := createAdminMap(uc.Admins)

	return adminMap, nil
}

// createAdminMap creates a mapping of admins with their git orgs/teams
func createAdminMap(admins string) map[string][]string {
	adminMap := make(map[string][]string)
	if strings.Contains(admins, ",") { // Multiple Organizations
		adminSplit := strings.Split(admins, ",")
		for _, admin := range adminSplit {
			splitAdmin := strings.Split(admin, "*")
			if teams, ok := adminMap[splitAdmin[0]]; ok {
				teams = append(teams, splitAdmin[1])
				adminMap[splitAdmin[0]] = teams
			} else { // TODO - make else if
				teams := []string{}
				teams = append(teams, splitAdmin[1])
				adminMap[splitAdmin[0]] = teams
			}
		}
	} else if strings.Contains(admins, "*") { // Single Organization
		adminSplit := strings.Split(admins, "*")
		teams := []string{adminSplit[1]}
		adminMap[adminSplit[0]] = teams
	} else {
		// String input does not follow format: "<org_name>*<team_name>""
		return nil
	}
	// Give admins access to master
	adminOrg = os.Getenv("GITHUB_ADMIN_ORG")
	adminTeam = os.Getenv("GITHUB_ADMIN_TEAM")
	adminMap[adminOrg] = []string{adminTeam}
	return adminMap
}

// updateMasterJSONFile updates the master's version via the master json
func updateMasterJSONFile(master, image, configPath string) (bool, error) {
	// Get masters json file
	configPath = util.TrimSuffix(configPath, "/") // remove trailing "/"
	filename := fmt.Sprintf("%s/%s.json", configPath, master)
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, err
	}

	// Unmarshal json data to a struct and do a image check
	var ac model.AdminConfig
	json.Unmarshal(bytes, &ac)

	if image != ac.Image {
		// Marshal json back to bytes and rewrite to file
		ac.Image = image
		bytes, _ := json.MarshalIndent(&ac, "", "	")
		ioutil.WriteFile(filename, bytes, 0644)
		return true, nil
	}

	// Image is already updated
	return false, nil
}
