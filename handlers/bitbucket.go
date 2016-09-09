package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/config"
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
		write(w, fmt.Sprintf("Bitbucket: %s", err.Error()))
		return
	}

	// Get the domain.
	domain := config.Get().Domain

	// Check pull request state.
	if strings.ToLower(p.PullRequest.State) != "merged" {
		write(w, "Bitbucket: Only pull requested with state `merged` can be handle.")
		return
	}

	// Split branch on `/`.
	branch := strings.Split(strings.ToLower(p.PullRequest.Source.Branch.Name), "/")

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", branch[len(branch)-1], domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		write(w, fmt.Sprintf("Bitbucket: %s", err.Error()))
		return
	}

	write(w, fmt.Sprintf("Bitbucket: Removed %d containers with name suffix %s", count, suffix))
}
