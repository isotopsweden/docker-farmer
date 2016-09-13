package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/config"
	"github.com/isotopsweden/docker-farmer/docker"
)

// gitlabPayload contains the object attributes object.
type gitlabPayload struct {
	After            string `json:"after"`
	Before           string `json:"before"`
	ObjectAttributes struct {
		SourceBranch string `json:"source_branch"`
		State        string `json:"state"`
	} `json:"object_attributes"`
	ObjectKind string `json:"object_kind"`
	Ref        string `json:"ref"`
}

// GitlabHandler will handle the payload from GitLab and remove containers
// based on a domain if they exists where the issue key is the prefix in the domain.
func GitlabHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p gitlabPayload
	if err := decoder.Decode(&p); err != nil {
		write(w, fmt.Sprintf("GitLab: %s", err.Error()))
		return
	}

	// Bail if not right webhook object kind.
	if p.ObjectKind != "push" && p.ObjectKind != "merge_request" {
		write(w, "GitLab: Only `merge_request` and `push` webhooks works.")
		return
	}

	var branch []string

	// Test object kind.
	if p.ObjectKind == "push" {
		// A branch that is deleted don't have any after commit, but a before commit.
		if p.Before != "0000000000000000000000000000000000000000" && p.After == "0000000000000000000000000000000000000000" {
			// Split branch on `/`.
			branch = strings.Split(strings.ToLower(p.Ref), "/")
		} else {
			write(w, "GitLab: Only push events with empty after commit and existing before commit works.")
			return
		}
	} else {
		// Check merge request state.
		if strings.ToLower(p.ObjectAttributes.State) != "merged" {
			write(w, "GitLab: Only merge requested with state `merged` can be handle.")
			return
		}

		// Split branch on `/`.
		branch = strings.Split(strings.ToLower(p.ObjectAttributes.SourceBranch), "/")
	}

	// Bail if branch is empty.
	if len(branch) == 0 {
		write(w, "GitLab: branch is empty, stopping.")
		return
	}

	// Get the domain.
	domain := config.Get().Domain

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", branch[len(branch)-1], domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		write(w, fmt.Sprintf("GitLab: %s", err.Error()))
		return
	}

	write(w, fmt.Sprintf("GitLab: Removed %d containers with name suffix %s", count, suffix))
}
