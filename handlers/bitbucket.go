package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/docker"
)

// bitbucketPayload contains the pull request object.
type bitbucketPayload struct {
	PullRequest struct {
		State  string `json:"state"`
		Source struct {
			Branch struct {
				Name string `json:"name"`
			} `json:"branch"`
		} `json:"source"`
	} `json:"pullrequest"`
}

// BitbucketHandler will handle the payload from Bitbucket and remove containers
// based on a domain if they exists where the issue key is the prefix in the domain.
func BitbucketHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p bitbucketPayload
	if err := decoder.Decode(&p); err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	// Get the domain.
	domain, err := getDomain()
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	// Check pull request state.
	if strings.ToLower(p.PullRequest.State) != "merged" {
		fmt.Fprintf(w, "Only pull requested with state `merged` can be handle.")
		return
	}

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", strings.ToLower(p.PullRequest.Source.Branch.Name), domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Removed %d containers with name suffix %s", count, suffix)
}
