package server

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/target/jenkins-docker-api/docker"
)

// Update represents a server function we use for handling PUT requests
// to the API. It will update a Jenkins master by updating its
// corresponding Docker service using the Docker API.
//
// swagger:operation PUT /jenkins/update/:master update putUpdate
//
// Update Jenkins Master
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
//     description: Successful update of the Jenkins master
//     schema:
//       type: string
//   '400':
//     description: There is something wrong about the request
//     schema:
//       type: string
//   '401':
//     description: Request to update is unauthorized
//     schema:
//       type: string
//   '404':
//     description: Jenkins master does not exist as a running Docker service in docker swarm
//     schema:
//       type: string
//   '500':
//     description: There was a problem within the API
//     schema:
//       type: string
func Update(c *gin.Context) {
	log.Info("Updating service")

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

	log.Infof("Attempting upgrade of %s:%s", serviceName, jenkinsImage)

	adminConfig := os.Getenv("JENKINS_ADMIN_CONFIG_PATH")
	imageUpdated, err := updateMasterJSONFile(serviceName, jenkinsImage, adminConfig)
	if err != nil {
		log.Errorf("Failed to update %s.json. Error: %s", serviceName, err.Error())
		msg := fmt.Sprintf("Oops, there was an internal problem updating '%s'", serviceName)
		c.String(http.StatusInternalServerError, msg)
		return
	}
	if !imageUpdated {
		msg := fmt.Sprintf("No updated needed. %s is already up-to-date", serviceName)
		log.Info(msg)
		c.String(http.StatusOK, msg)
		return
	}

	err = docker.UpdateService(serviceName, jenkinsImage, false)
	if err != nil {
		log.Error(err)
		msg := "Oops, there was a problem updating with the Docker API"
		c.String(http.StatusInternalServerError, msg)
		return
	}

	msg := fmt.Sprintf("Successfully updated %s", serviceName)
	log.Info(msg)
	c.String(http.StatusOK, msg)
}

// UpdateAll represents a server function we use for handling PUT requests
// to the API. It will update all Jenkins masters by updating all
// actively running Docker services using the Docker API.
//
// swagger:operation PUT /jenkins/admin/update_all updateAll putUpdateAll
//
// Update All Jenkins Masters (admins only)
//
// ---
// x-success_http_code: '200'
// x-response_time_sla: 5000ms
// x-expected_tps: '5'
// produces:
// - application/json
// parameters:
// - in: header
//   name: Authorization
//   description: GitHub personal access token
//   required: true
//   type: string
//   default: 'Personal Access Token '
// responses:
//   '200':
//     description: Successful update of all Jenkins masters
//     schema:
//       type: string
//   '401':
//     description: Request to update all is unauthorized
//     schema:
//       type: string
//   '500':
//     description: There was a problem within the API
//     schema:
//       type: string
func UpdateAll(c *gin.Context) {
	log.Info("Updating all services")

	// Get all running services
	services, err := docker.GetServices()
	if err != nil {
		log.Error(err)
		msg := "Oops, there was a problem finding all active services with the Docker API"
		c.String(http.StatusInternalServerError, msg)
		return
	}

	masters := []string{}
	for _, service := range services {
		if docker.ServiceExists(services, service.Spec.Name) {
			masters = append(masters, service.Spec.Name)
		}
	}

	// Authenticate user for updating
	adminOrg = os.Getenv("GITHUB_ADMIN_ORG")
	adminTeam = os.Getenv("GITHUB_ADMIN_TEAM")
	adminMap := make(map[string][]string)
	adminMap[adminOrg] = []string{adminTeam}

	err = authenticateUser(c, "", adminMap)
	if err != nil {
		log.Error(err)
		msg := fmt.Sprintf("Unable to authenticate your request to update because %s", err.Error())
		c.String(http.StatusUnauthorized, msg)
		return
	}

	jenkinsImage := os.Getenv("JENKINS_IMAGE")
	adminConfig := os.Getenv("JENKINS_ADMIN_CONFIG_PATH")
	updateCount := 0
	failures := []string{}
	for _, master := range masters {
		// Update image
		imageUpdated, err := updateMasterJSONFile(master, jenkinsImage, adminConfig)
		if err != nil {
			log.Errorf("Failed to read %s.json. Error: %s", master, err.Error())
			failures = append(failures, master)
			continue
		}
		if imageUpdated {
			err := docker.UpdateService(master, jenkinsImage, false)
			if err != nil {
				log.Errorf("Failed to update '%s'.", master)
				failures = append(failures, master)
				continue
			}
			updateCount++
			log.Info(fmt.Sprintf("%s has been updated", master))
		} else {
			log.Info(fmt.Sprintf("%s is already up-to-date", master))
		}
	}

	// Either all services were updated or some updated and some were already up-to-date
	if len(failures) == 0 {
		msg := fmt.Sprintf("Successfully updated %d masters", updateCount)
		log.Info(msg)
		c.String(http.StatusOK, msg)
		return
	}

	// There were one ore more docker service update failures
	msg := fmt.Sprintf("Updated %d masters successfully, but failed to update these masters:", updateCount)
	for i, fail := range failures {
		if i == (len(failures) - 1) {
			msg = fmt.Sprintf(msg+" %s.", fail)
			break
		}
		msg = fmt.Sprintf(msg+" %s,", fail)
	}
	log.Info(msg)
	c.String(http.StatusInternalServerError, msg)
}
