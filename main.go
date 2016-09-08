package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/isotopsweden/docker-farmer/docker"
	"github.com/isotopsweden/docker-farmer/handlers"
)

var (
	configFlag = flag.String("config", "", "Path to config file")
)

// Config represents a config struct.
type Config struct {
	Domain string
	Docker struct {
		Host    string
		Version string
	}
	Listen string
}

// getConfig will return a config struct.
func getConfig() Config {
	path := "config.json"

	if len(*configFlag) > 0 {
		path = *configFlag
	}

	file, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Printf("Config error: %v\n", err)
		return Config{}
	}

	var config Config

	json.Unmarshal(file, &config)

	return config
}

func main() {
	flag.Parse()

	config := getConfig()

	if config.Listen[0] == ':' {
		config.Listen = "0.0.0.0" + config.Listen
	}

	// Setup required docker host and version information.
	docker.SetHost(config.Docker.Host)
	docker.SetVersion(config.Docker.Version)

	// Index route.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Docker Farmer!")
	})

	// Containers route.
	http.HandleFunc("/containers", func(w http.ResponseWriter, r *http.Request) {
		containers, err := docker.GetContainers(config.Domain)

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

	fmt.Printf("Listening to http://%s\n", config.Listen)

	// Listen to port.
	log.Fatal(http.ListenAndServe(config.Listen, nil))
}
