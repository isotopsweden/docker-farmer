package docker

import (
	"crypto/md5"
	"fmt"
	"log"

	"context"

	"github.com/docker/docker/api/types"
)

// DeleteMySQLDatabase will try to delete a database based on prefix and domain,
// the domain will be converted to a md5 hash.
func DeleteMySQLDatabase(user, password, prefix, name, container string) (bool, error) {
	client, err := getDockerClient()
	ctx := context.Background()

	if err != nil {
		return false, err
	}

	if user == "" {
		user = "root"
	}

	if password == "" {
		password = "root"
	}

	// Remove trailing slash from name.
	name = name[1:]

	// Create domain based on prefix and md5.
	dbname := fmt.Sprintf("%s%s", prefix, fmt.Sprintf("%x", md5.Sum([]byte(name))))

	log.Println(fmt.Sprintf("Trying to run drop database on %s", name))

	// Exec create on container.
	res, err := client.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Cmd: []string{
			"mysql",
			fmt.Sprintf("-u%s", user),
			fmt.Sprintf("-p%s", password),
			fmt.Sprintf("-edrop database %s;", dbname),
		},
	})

	log.Println(fmt.Sprintf("Exec create response from container %s: %s", container, res.ID))

	if err != nil {
		return false, err
	}

	// Exec start on container.
	err = client.ContainerExecStart(ctx, res.ID, types.ExecStartCheck{
		Detach: false,
	})

	if err != nil {
		return false, err
	}

	inspect, err := client.ContainerExecInspect(ctx, res.ID)

	if err != nil {
		return false, err
	}

	return inspect.ExitCode == 0, nil
}
