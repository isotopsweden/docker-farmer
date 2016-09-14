package main

import (
	"flag"
	"fmt"
	"html/template"
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
func realpath(file string) string {
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
	config.Init(realpath(*configFlag))
	c := config.Get()

	if c.Listen[0] == ':' {
		c.Listen = "0.0.0.0" + c.Listen
	}

	// Index route.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates := template.Must(template.ParseFiles(realpath(*publicFlag) + "/index.html"))
		err := templates.ExecuteTemplate(w, "index.html", map[string]interface{}{
			"Config": c,
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Assets route.
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(realpath(*publicFlag)+"/assets"))))

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

	// Test service route.
	http.HandleFunc("/services/test", handlers.TestHandler)

	fmt.Printf("Listening to http://%s\n", c.Listen)

	// Listen to port.
	log.Fatal(http.ListenAndServe(c.Listen, nil))
}
