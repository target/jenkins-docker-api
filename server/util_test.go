package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/franela/goblin"
	"github.com/google/go-github/github"
	"github.com/target/gelvedere/model"
	"github.com/target/jenkins-docker-api/fixtures/test"
)

func Test_utils(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Tests for server/util.go", func() {

		g.Describe("The isAuthorizedUser function", func() {
			g.It("should return true if a team is part of the adminMap", func() {
				team := test.MockTeam("mockTeam", "mockOrg")
				teams := []*github.Team{team}

				adminMap := make(map[string][]string)
				adminMap["mockOrg"] = []string{"mockTeam"}

				actual := isAuthorizedUser(teams, adminMap)

				g.Assert(actual).IsTrue()
			})

			g.It("should return true if multiple teams are part of the adminMap", func() {
				team1 := test.MockTeam("mockTeam1", "mockOrg1")
				team2 := test.MockTeam("mockTeam2", "mockOrg2")
				team3 := test.MockTeam("mockTeam3", "mockOrg3")
				teams := []*github.Team{team1, team2, team3}

				adminMap := make(map[string][]string)
				adminMap["mockOrg1"] = []string{"mockTeam1"}
				adminMap["mockOrg2"] = []string{"mockTeam2"}
				adminMap["mockOrg3"] = []string{"mockTeam3"}

				actual := isAuthorizedUser(teams, adminMap)

				g.Assert(actual).IsTrue()
			})

			g.It("should return false if a team is not part of the adminMap", func() {
				team := test.MockTeam("mockTeam", "mockOrg")
				teams := []*github.Team{team}

				adminMap := make(map[string][]string)
				adminMap["org"] = []string{"team"}

				actual := isAuthorizedUser(teams, adminMap)

				g.Assert(actual).IsFalse()
			})

			g.It("should return false if adminMap is nil", func() {
				team := test.MockTeam("mockTeam", "mockOrg")
				teams := []*github.Team{team}

				actual := isAuthorizedUser(teams, nil)

				g.Assert(actual).IsFalse()
			})
		})

		g.Describe("The getAdminMap function", func() {
			g.BeforeEach(func() {
				os.Setenv("JENKINS_ENV", "test")
				os.Setenv("GITHUB_ADMIN_ORG", "Jenkins")
				os.Setenv("GITHUB_ADMIN_TEAM", "Admins")
			})

			g.It("should return a map based on the mock JSON file", func() {
				master := "mock"
				mockConfigPath, _ := filepath.Abs("../fixtures/")

				exp := make(map[string][]string)
				exp["mockOrg"] = []string{"mockTeam"}
				exp["Jenkins"] = []string{"Admins"}

				actual, err := getAdminMap(master, mockConfigPath)
				if err != nil {
					g.Fail(err)
				}

				eq := reflect.DeepEqual(exp, actual)

				g.Assert(eq).IsTrue()
			})

			g.It("should return nil if provided an invalid file path", func() {
				master := "mock"
				mockConfigPath, _ := filepath.Abs("../bad/path/")

				actual, err := getAdminMap(master, mockConfigPath)
				if err != nil {
					g.Assert(actual == nil).IsTrue()
				} else {
					g.Fail(nil)
				}
			})
		})

		g.Describe("The createAdminMap function", func() {
			g.It("should return a map of a single organization with one team + admins", func() {
				testadmins := "foo*bar"

				exp := make(map[string][]string)
				exp["foo"] = []string{"bar"}
				exp["Jenkins"] = []string{"Admins"}

				actual := createAdminMap(testadmins)

				eq := reflect.DeepEqual(exp, actual)

				g.Assert(eq).IsTrue()
			})
			g.It("should return a map of a single organization with multiple teams + admins", func() {
				testadmins := "foo*bar,foo*bar1,foo*bar2"

				exp := make(map[string][]string)
				exp["foo"] = []string{"bar", "bar1", "bar2"}
				exp["Jenkins"] = []string{"Admins"}

				actual := createAdminMap(testadmins)

				eq := reflect.DeepEqual(exp, actual)

				g.Assert(eq).IsTrue()
			})
			g.It("should return a map of multiple organizations and teams + admins", func() {
				testadmins := "foo*bar,bar*foo"

				exp := make(map[string][]string)
				exp["foo"] = []string{"bar"}
				exp["bar"] = []string{"foo"}
				exp["Jenkins"] = []string{"Admins"}

				actual := createAdminMap(testadmins)

				eq := reflect.DeepEqual(exp, actual)

				g.Assert(eq).IsTrue()
			})
			g.It("should return a nil map if input does not follow format: <org_name>*<team_name>", func() {
				testadmins := "this is a bad string"

				actual := createAdminMap(testadmins)

				eq := reflect.DeepEqual(nil, actual)

				g.Assert(eq).IsFalse()
			})
		})

		g.Describe("The updateMasterJSONFile function", func() {
			// Set test images
			const oldImage = "target/jenkins-docker-master:2.73.3-1"
			const newImage = "target/jenkins-docker-master:2.73.3-2"
			const filename = "secrets"
			var adminConfigPath, _ = filepath.Abs("../fixtures/test/")

			// Set initial image in test/secrets.json before updating
			g.BeforeEach(func() {
				file := fmt.Sprintf("%s/%s.json", adminConfigPath, filename)
				bytes, err := ioutil.ReadFile(file)
				if err != nil {
					g.Fail(err)
				}

				var ac model.AdminConfig
				json.Unmarshal(bytes, &ac)

				ac.Image = oldImage
				bytes, _ = json.MarshalIndent(&ac, "", "	")

				err = ioutil.WriteFile(file, bytes, 0644)
				if err != nil {
					g.Fail(err)
				}
			})
			g.It("should return true when updating to a new image to indicate a successful update", func() {
				// Update image
				updated, err := updateMasterJSONFile(filename, newImage, adminConfigPath)

				g.Assert(updated && err == nil).IsTrue()

				// Get newly updated image
				file := fmt.Sprintf("%s/%s.json", adminConfigPath, filename)
				bytes, err := ioutil.ReadFile(file)
				if err != nil {
					g.Fail(err)
				}

				var ac model.AdminConfig
				json.Unmarshal(bytes, &ac)

				g.Assert(ac.Image != oldImage).IsTrue()
				g.Assert(ac.Image == newImage).IsTrue()
			})
			g.It("should return false when updating an already up-to-date image", func() {
				// Update image with old (same) image
				updated, err := updateMasterJSONFile(filename, oldImage, adminConfigPath)

				g.Assert(!updated && err == nil).IsTrue()

				// Get newly updated image
				file := fmt.Sprintf("%s/%s.json", adminConfigPath, filename)
				bytes, err := ioutil.ReadFile(file)
				if err != nil {
					g.Fail(err)
				}

				var ac model.AdminConfig
				json.Unmarshal(bytes, &ac)

				g.Assert(ac.Image == oldImage).IsTrue()
				g.Assert(ac.Image != newImage).IsTrue()
			})
			g.It("should return false to indicate there was an error when given a bad file path", func() {
				// Update image with old (same) image
				updated, err := updateMasterJSONFile(filename, oldImage, "/bad/path")
				g.Assert(!updated && err != nil).IsTrue()
			})
		})
	})
}
