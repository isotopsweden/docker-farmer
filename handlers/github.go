package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/config"
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
		write(w, fmt.Sprintf("GitHub: %s", err.Error()))
		return
	}

	// Get the domain.
	domain := config.Get().Domain

	// Check if the pull request is merged.
	if !p.PullRequest.Merged {
		write(w, "GitHub: Only merged pull requests can be handle.")
		return
	}

	// Split branch on `/`.
	branch := strings.Split(strings.ToLower(p.PullRequest.Head.Ref), "/")

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", branch[len(branch)-1], domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		write(w, fmt.Sprintf("Bitbucket: %s", err.Error()))
		return
	}

	write(w, fmt.Sprintf("GitHub: Removed %d containers with name suffix %s", count, suffix))
}
