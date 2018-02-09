package docker

import (
	"context"
	"fmt"

	"github.com/target/jenkins-docker-api/util"

	log "github.com/Sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
)

// GetServices returns a list of running Docker swarm services
func GetServices() ([]swarm.Service, error) {
	log.Info("Grabbing list of running services")

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	services, err := cli.ServiceList(ctx, types.ServiceListOptions{})
	if err != nil {
		return nil, err
	}

	return services, nil
}

// ServiceExists verifies if a specified service is running in the Docker swarm
func ServiceExists(services []swarm.Service, serviceName string) bool {
	var exists = false
	for _, service := range services {
		if serviceName == service.Spec.Name {
			environment := service.Spec.TaskTemplate.ContainerSpec.Env
			if util.SliceContains("JENKINS_URL", environment) {
				log.Info(fmt.Sprintf("Found service %s", service.Spec.Name))
				exists = true
				break
			}
		}
	}
	return exists
}

func RemoveService(serviceNameOrID string) error {
	log.Infof("Removing service %s", serviceNameOrID)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	err = cli.ServiceRemove(ctx, serviceNameOrID)
	if err != nil {
		return err
	}

	return nil
}

// UpdateService updates a specified Docker swarm service
func UpdateService(serviceNameOrID, jenkinsImage string, forceUpdate bool) error {
	log.Infof("Updating service %s", serviceNameOrID)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	// get service
	swarmService, _, err := cli.ServiceInspectWithRaw(ctx, serviceNameOrID, types.ServiceInspectOptions{})
	if err != nil {
		return err
	}

	// update service image/version
	if jenkinsImage != "" {
		swarmService.Spec.TaskTemplate.ContainerSpec.Image = jenkinsImage
	}

	// force update?
	if forceUpdate {
		swarmService.Spec.TaskTemplate.ForceUpdate++
	}

	// get update specs
	swarmVersion := swarmService.Version
	serviceSpec := swarmService.Spec
	options := types.ServiceUpdateOptions{}

	// update the service
	response, err := cli.ServiceUpdate(ctx, serviceNameOrID, swarmVersion, serviceSpec, options)
	if err != nil {
		return err
	}

	// log warnings
	if len(response.Warnings) > 0 {
		for _, warning := range response.Warnings {
			log.Printf("WARNING: %s:", warning)
		}
	}
	return nil
}
