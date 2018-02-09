package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/target/jenkins-docker-api/docker"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// Shutdown represents a server function we use for handling PUT requests
// to the API. It will shutdown a Jenkins master by by forcing a Docker Swarm
// service remove using the Docker API.
//
// swagger:operation PUT /jenkins/shutdown/:master shutdown putShutdown
//
// Shutdown Jenkins Master
//
// ---
// x-success_http_code: '200'
// x-response_time_sla: 5000ms
// x-expected_tps: '5'
// produces:
// - text/plain
// parameters:
// - in: header
//   name: Authorization
//   description: GitHub personal access token
//   required: true
//   type: string
//   default: 'Personal Access Token '
// responses:
//   '200':
//     description: Successful shutdown of the Jenkins master
//     schema:
//       type: string
//   '400':
//     description: There is something wrong about the request
//     schema:
//       type: string
//   '401':
//     description: Request to shutdown is unauthorized
//     schema:
//      type: string
//   '404':
//     description: Jenkins master does not exist as a running Docker service in docker swarm
//     schema:
//       type: string
//   '500':
//     description: There was a problem within the API
//     schema:
//       type: string
func Shutdown(c *gin.Context) {
	log.Info("Shutting down service")

	// Get Jenkins master name
	serviceName := c.Param("master")
	if serviceName == "" {
		msg := "Master name is required"
		log.Error(msg)
		c.String(http.StatusBadRequest, msg)
		return
	}

	// Get running services
	services, err := docker.GetServices()
	if err != nil {
		log.Error(err)
		msg := "Oops, there was a problem finding all active services with the Docker API"
		c.String(http.StatusInternalServerError, msg)
		return
	}

	// Check if service exists
	if !docker.ServiceExists(services, serviceName) {
		log.Errorf("Did not find Docker service '%s'", serviceName)
		msg := fmt.Sprintf("Did not find master '%s'", serviceName)
		c.String(http.StatusNotFound, msg)
		return
	}

	// Authenticating user for update
	err = authenticateUser(c, serviceName, nil)
	if err != nil {
		log.Error(err)
		msg := fmt.Sprintf("Unable to authenticate your request to update because %s", err.Error())
		c.String(http.StatusUnauthorized, msg)
		return
	}

	jenkinsImage := os.Getenv("JENKINS_IMAGE")

	log.Infof("Attempting shutdown of %s:%s", serviceName, jenkinsImage)

	err = docker.RemoveService(serviceName)
	if err != nil {
		log.Error(err)
		msg := fmt.Sprintf("Unable to shutdown master because %s", err.Error())
		c.String(http.StatusInternalServerError, msg)
		return
	}

	msg := fmt.Sprintf("Successfully shutdown %s", serviceName)
	log.Info(msg)
	c.String(http.StatusOK, msg)
}
