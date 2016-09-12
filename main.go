package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/isotopsweden/docker-farmer/config"
	"github.com/isotopsweden/docker-farmer/handlers"
	"github.com/kardianos/osext"
)

var (
	configFlag = flag.String("config", "config.json", "Path to config file")
	publicFlag = flag.String("public", "public", "Path to public directory")
)

// path will return right path for file, looks at the
// given file first and then looks in the executable folder.
func path(file string) string {
	if _, err := os.Stat(file); err == nil {
		return file
	}

	path, _ := osext.ExecutableFolder()

	if _, err := os.Stat(path + "/" + file); os.IsNotExist(err) {
		return file
	}

	return path + "/" + file
}

func main() {
	flag.Parse()

	// Init config.
	config.Init(path(*configFlag))
	c := config.Get()

	if c.Listen[0] == ':' {
		c.Listen = "0.0.0.0" + c.Listen
	}

	// Index route.
	http.Handle("/", http.FileServer(http.Dir(path(*publicFlag))))

	// Config api route.
	http.HandleFunc("/api/config", handlers.ConfigHandler)

	// Containers api route.
	http.HandleFunc("/api/containers", handlers.ContainersHandler)

	// Database api route.
	http.HandleFunc("/api/database", handlers.DatabaseHandler)

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
