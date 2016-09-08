package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/docker"
)

// githubPayload contains the pull request object.
type githubPayload struct {
	PullRequest struct {
		Head struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Merged bool `json:"merged"`
	} `json:"pull_request"`
}

// GithubHandler will handle the payload from GitHub and remove containers
// based on a domain if they exists where the issue key is the prefix in the domain.
func GithubHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p githubPayload
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

	// Check if the pull request is merged.
	if !p.PullRequest.Merged {
		fmt.Fprintf(w, "Only merged pull requests can be handle.")
		return
	}

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", strings.ToLower(p.PullRequest.Head.Ref), domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Removed %d containers with name suffix %s", count, suffix)
}
