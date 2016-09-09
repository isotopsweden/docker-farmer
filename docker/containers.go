package docker

import (
	"log"
	"strings"

	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/isotopsweden/docker-farmer/config"
)

// getDockerClient will return the Docker client or error.
func getDockerClient() (*client.Client, error) {
	defaultHeaders := map[string]string{"User-Agent": "docker-farmer"}
	client, err := client.NewClient(GetHost(), GetVersion(), nil, defaultHeaders)

	if err != nil {
		return nil, err
	}

	return client, nil
}

// GetHost returns Docker host.
func GetHost() string {
	if config.Get().Docker.Host == "" {
		return "unix:///var/run/docker.sock"
	}

	return config.Get().Docker.Host
}

// GetVersion returns Docker version.
func GetVersion() string {
	return config.Get().Docker.Version
}

// GetContainers returns all containers for a domain suffix or a error.
func GetContainers(domain string) ([]types.Container, error) {
	client, err := getDockerClient()
	ctx := context.Background()

	if err != nil {
		return nil, err
	}

	options := types.ContainerListOptions{All: true}
	containers, err := client.ContainerList(ctx, options)

	if err != nil {
		return nil, err
	}

	result := []types.Container{}
	for _, c := range containers {
		// No name on the container.
		if len(c.Names) < 1 {
			continue
		}

		// Domain suffix does not exists.
		if !strings.HasSuffix(c.Names[0], domain) {
			continue
		}

		result = append(result, c)
	}

	return result, nil
}

// RemoveContainers will remove containers with the domain suffix
// and return a count of containers removed or a error.
func RemoveContainers(domain string) (int, error) {
	client, err := getDockerClient()
	ctx := context.Background()

	if err != nil {
		return 0, err
	}

	containers, err := GetContainers(domain)

	if err != nil {
		return 0, err
	}

	count := 0
	conf := config.Get()

	for _, c := range containers {
		// Try to force remove the container.
		if err := client.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{
			RemoveVolumes: true,
			Force:         true,
		}); err != nil {
			log.Println("Docker error: ", err.Error())
			continue
		}

		// Try to delete database
		_, err := DeleteDatabase(conf.Database.User, conf.Database.Password, conf.Database.Prefix, c.Names[0], conf.Database.Container)
		if err != nil {
			log.Println("Database error: ", err.Error())
		}

		count++
	}

	return count, nil
}
