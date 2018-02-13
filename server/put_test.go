package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/target/jenkins-docker-api/fixtures"

	log "github.com/Sirupsen/logrus"
	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func Test_Put(t *testing.T) {
	log.SetLevel(log.PanicLevel)

	gin.SetMode(gin.TestMode)

	var path string
	var request *http.Request

	rr := httptest.NewRecorder()

	s := httptest.NewServer(fixtures.FakeAPIHandler())

	r := gin.New()

	g := goblin.Goblin(t)
	g.Describe("Tests for restart.go and update.go", func() {
		testMaster := "loadtest"

		// Init test endpoints
		r.PUT("/jenkins/restart/:master", Restart)
		r.PUT("/jenkins/update/:master", Update)
		r.PUT("/jenkins/admin/update_all", UpdateAll)
		r.PUT("/jenkins/shutdown/:master", Shutdown)

		// Set revelant env variables for mocking client API calls
		os.Setenv("DOCKER_HOST", s.URL)
		os.Setenv("DOCKER_API_VERSION", "v1.36")

		os.Setenv("GITHUB_API_URL", s.URL+"/api/v3/")

		// Set admin GitHub env variables
		os.Setenv("GITHUB_ADMIN_ORG", "Jenkins")
		os.Setenv("GITHUB_ADMIN_TEAM", "Admins")

		// Set revelent Jenkins env variables
		os.Setenv("JENKINS_IMAGE", "jenkins.docker.test.com/jenkins-test:v2")
		loadtestUserConfigPath, _ := filepath.Abs("../fixtures/test/user/")
		os.Setenv("JENKINS_USER_CONFIG_PATH", loadtestUserConfigPath)
		loadtestAdminConfigPath, _ := filepath.Abs("../fixtures/test/admin/")
		os.Setenv("JENKINS_ADMIN_CONFIG_PATH", loadtestAdminConfigPath)

		// Init a new mocker before each test
		g.Before(func() {
			rr = httptest.NewRecorder()
		})

		// Close server after each test
		g.After(func() {
			s.Close()
		})

		g.It("- should restart the test Jenkins master", func() {
			// Create an http request for the Restart endpoint
			path = "/jenkins/restart/" + testMaster
			request = fixtures.CreateFakeRequest("PUT", path, "token test_token", nil)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(rr, request)

			// Make sure we get the proper response we expect
			g.Assert(rr.Code).Equal(http.StatusOK)
		})

		g.It("- should update the test Jenkins master", func() {
			// Create an http request for the Update endpoint
			path = "/jenkins/update/" + testMaster
			request = fixtures.CreateFakeRequest("PUT", path, "token test_token", nil)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(rr, request)

			// Make sure we get the proper response we expect
			g.Assert(rr.Code).Equal(http.StatusOK)
		})

		g.It("- should update all test Jenkins masters", func() {
			// Create an http request for the UpdateAll endpoint
			path = "/jenkins/admin/update_all"
			request = fixtures.CreateFakeRequest("PUT", path, "token test_token", nil)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(rr, request)

			// Make sure we get the proper response we expect
			g.Assert(rr.Code).Equal(http.StatusOK)
		})

		g.It("- should shutdown test Jenkins masters", func() {
			// Create an http request for the shutdown endpoint
			path = "/jenkins/shutdown/" + testMaster
			request = fixtures.CreateFakeRequest("PUT", path, "token test_token", nil)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(rr, request)

			// Make sure we get the proper response we expect
			g.Assert(rr.Code).Equal(http.StatusOK)
		})
	})
}
