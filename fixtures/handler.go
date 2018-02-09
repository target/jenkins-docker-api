package fixtures

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// FakeAPIHandler returns an http.Handler that is capable of handling a variety
// of mock Docker API requests and returning mock responses.
func FakeAPIHandler() http.Handler {
	gin.SetMode(gin.TestMode)

	e := gin.New()

	// Faked Docker API endpoints
	e.GET("/v1.36/services", DockerServiceList)
	e.GET("/v1.36/services/:serviceID", DockerService)
	e.POST("/v1.36/services/:serviceID/update", DockerServiceUpdate)
	e.POST("/v1.36./services/:serviceID/remove", DockerServiceRemove)
	e.DELETE("/v1.36/services/:serviceID", DockerService)

	// Faked GitHub API endpoints
	e.GET("/api/v3/user/teams", GitHubTeams)
	e.GET("/api/v3/user", GitHubUser)

	return e
}

func DockerServiceList(c *gin.Context) {
	c.String(http.StatusOK, DockerSwarmServicesResponse)
}

func DockerService(c *gin.Context) {
	c.String(http.StatusOK, DockerSwarmServiceResponse)
}

func DockerServiceUpdate(c *gin.Context) {
	c.String(http.StatusOK, "{}")
}

func DockerServiceRemove(c *gin.Context) {
	c.String(http.StatusOK, "{}")
}

func GitHubTeams(c *gin.Context) {
	c.String(http.StatusOK, GitHubTeamsResponse)
}

func GitHubUser(c *gin.Context) {
	c.String(http.StatusOK, GitHubUserResponse)
}
