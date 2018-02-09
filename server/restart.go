package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/target/jenkins-docker-api/docker"
)

// Restart represents a server function we use for handling PUT requests
// to the API. It will restart a Jenkins master by by forcing a Docker Swarm
// service update using the Docker API.
//
// swagger:operation PUT /jenkins/restart/:master restart putRestart
//
// Restart Jenkins Master
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
//     description: Successful restart of the Jenkins master
//     schema:
//       type: string
//   '400':
//     description: There is something wrong about the request
//     schema:
//       type: string
//   '401':
//     description: Request to restart is unauthorized
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
func Restart(c *gin.Context) {
	log.Info("Restarting service")

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

	// Search for specified service and exec a service update
	if !docker.ServiceExists(services, serviceName) {
		log.Errorf("Did not find Docker service '%s'", serviceName)
		msg := fmt.Sprintf("Did not find master '%s'", serviceName)
		c.String(http.StatusNotFound, msg)
		return
	}

	// Authenticating user for restart
	err = authenticateUser(c, serviceName, nil)
	if err != nil {
		log.Error(err)
		msg := fmt.Sprintf("Unable to authenticate your request to restart because %s", err.Error())
		c.String(http.StatusUnauthorized, msg)
		return
	}

	log.Info(fmt.Sprintf("Attempting restart of %s", serviceName))

	err = docker.UpdateService(serviceName, "", true)
	if err != nil {
		log.Error(err)
		msg := "Oops, there was a problem restarting with the Docker API"
		c.String(http.StatusInternalServerError, msg)
		return
	}

	msg := fmt.Sprintf("Successfully restarted %s", serviceName)
	log.Info(msg)
	c.String(http.StatusOK, msg)
}
