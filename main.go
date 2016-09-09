package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/isotopsweden/docker-farmer/config"
	"github.com/isotopsweden/docker-farmer/docker"
	"github.com/isotopsweden/docker-farmer/handlers"
)

var (
	configFlag = flag.String("config", "", "Path to config file")
)

func main() {
	flag.Parse()

	// Init config.
	config.Init(*configFlag)
	c := config.Get()

	if c.Listen[0] == ':' {
		c.Listen = "0.0.0.0" + c.Listen
	}

	// Setup required docker host and version information.
	docker.SetHost(c.Docker.Host)
	docker.SetVersion(c.Docker.Version)

	// Index route.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		containers, err := docker.GetContainers(c.Domain)

		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			sites := []string{}

			for _, c := range containers {
				name := c.Names[0][1:]
				sites = append(sites, name)
			}

			j, err := json.Marshal(map[string]interface{}{
				"sites": sites,
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, string(j))
			}
		}
	})

	// Containers route.
	http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		containers, err := docker.GetContainers(c.Domain)

		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			j, err := json.Marshal(containers)

			if err != nil {
				fmt.Fprintf(w, err.Error())
			} else {
				fmt.Fprintf(w, string(j))
			}
		}
	})

	// BitBucket service route.
	http.HandleFunc("/services/bitbucket", handlers.BitbucketHandler)

	// GitHub service route.
	http.HandleFunc("/services/github", handlers.GithubHandler)

	// GitLab service route.
	http.HandleFunc("/services/gitlab", handlers.GitlabHandler)

	// Jira service route.
	http.HandleFunc("/services/jira", handlers.JiraHandler)

	fmt.Printf("Listening to http://%s\n", c.Listen)

	// Listen to port.
	log.Fatal(http.ListenAndServe(c.Listen, nil))
}
