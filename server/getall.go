package server

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/target/jenkins-docker-api/docker"
)

// GetAll represents a server function we use for handling GET requests
// to the API. It will grab a list of all active Jenkins masters by retrieving
// corresponding running Docker services using the Docker API.
//
// swagger:operation GET /jenkins get getAll
//
// Get All Jenkins Masters
//
// ---
// x-success_http_code: '200'
// x-response_time_sla: 5000ms
// x-expected_tps: '5'
// produces:
// - text/plain
// responses:
//   '200':
//     description: Successful retrieval of all active Jenkins masters
//     schema:
//       type: string
//   '500':
//     description: There was a problem within the API
//     schema:
//       type: string
func GetAll(c *gin.Context) {
	log.Info("Attempting retrieval of services")

	// Get running services
	services, err := docker.GetServices()
	if err != nil {
		log.Error(err)
		msg := "Oops, there was a problem finding all active services with the Docker API"
		c.String(http.StatusInternalServerError, msg)
		return
	}

	// Ensure that we only return jenkins master
	res := []string{}
	for _, service := range services {
		if docker.ServiceExists(services, service.Spec.Name) {
			res = append(res, service.Spec.Name)
		}
	}

	if len(res) > 0 {
		log.Info("Successfully retrieved services")
	} else {
		log.Info("No services were retrieved")
	}

	c.JSON(http.StatusOK, res)
}
