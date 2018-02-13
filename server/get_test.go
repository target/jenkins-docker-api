package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/target/jenkins-docker-api/fixtures"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func Test_GetAll(t *testing.T) {
	log.SetLevel(log.PanicLevel)

	gin.SetMode(gin.TestMode)

	var path string
	var request *http.Request

	expected := new(bytes.Buffer)
	actual := new(bytes.Buffer)

	rr := httptest.NewRecorder()

	s := httptest.NewServer(fixtures.FakeAPIHandler())

	r := gin.New()

	g := goblin.Goblin(t)
	g.Describe("GetAll", func() {
		r.GET("/jenkins", GetAll)

		g.Before(func() {
			rr = httptest.NewRecorder()

			path = "/jenkins"
			request = fixtures.CreateFakeRequest("GET", path, "", nil)

			expected = new(bytes.Buffer)
			actual = new(bytes.Buffer)

			os.Setenv("DOCKER_HOST", s.URL)
			os.Setenv("DOCKER_API_VERSION", "v1.36")
		})

		g.After(func() {
			s.Close()
		})

		g.It("- should get all test Jenkins masters", func() {
			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			r.ServeHTTP(rr, request)

			// Make sure we get the proper response we expect
			g.Assert(rr.Code).Equal(http.StatusOK)

			// Trim all extra whitespaces and newline characters
			_ = json.Compact(expected, []byte(fixtures.JenkinsGetAllResponse))
			_ = json.Compact(actual, rr.Body.Bytes())

			g.Assert(expected).Equal(actual)
		})
	})
}
