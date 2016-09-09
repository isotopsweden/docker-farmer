package docker

import (
	"crypto/md5"
	"fmt"
	"log"
	"strings"
	"unicode"

	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
)

// DeleteDatabase will try to delete a database based on prefix and domain,
// the domain will be converted to a md5 hash.
func DeleteDatabase(user, password, prefix, name, container string) (bool, error) {
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

	// Trim whitespace in container name.
	name = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, name)

	// Create domain based on prefix and md5.
	dbname := fmt.Sprintf("%s%s", prefix, fmt.Sprintf("%x", md5.Sum([]byte(name))))

	// Create mysql command.
	cmd := fmt.Sprintf("mysql -u%s -p%s -e\"drop database %s\"", user, password, dbname)

	log.Println(fmt.Sprintf("Trying to run: %s on container %s with name %s", cmd, container, name))

	// Exec create on container.
	res, err := client.ContainerExecCreate(ctx, container, types.ExecConfig{
		AttachStdin:  false,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
		Cmd: []string{
			cmd,
		},
	})

	log.Println(fmt.Sprintf("Exec create response from container %s: %s", container, res.ID))

	if err != nil {
		return false, err
	}

	// Exec start on container.
	err = client.ContainerExecStart(ctx, res.ID, types.ExecStartCheck{})

	if err != nil {
		return false, err
	}

	return true, nil
}
