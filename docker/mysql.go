package docker

import (
	"crypto/md5"
	"fmt"
	"log"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// DeleteDatabase will try to delete a database based on prefix and domain,
// the domain will be converted to a md5 hash.
func DeleteDatabase(user, password, prefix, domain, container string) (bool, error) {
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

	dbname := fmt.Sprintf("%s%s", prefix, fmt.Sprintf("%x", md5.Sum([]byte(domain))))
	cmd := fmt.Sprintf("mysql -u%s -p%s -e\"drop database %s\"", user, password, dbname)

	log.Println(fmt.Sprintf("Trying to run: %s on container %s", cmd, container))

	res, err := client.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd: []string{
			cmd,
		},
	})

	log.Println(fmt.Sprintf("Response from container %s: %s", container, res.ID))

	if err != nil {
		return false, err
	}

	return true, nil
}
