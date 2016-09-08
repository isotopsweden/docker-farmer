package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/docker"
)

// gitlabPayload contains the object attributes object.
type gitlabPayload struct {
	ObjectAttributes struct {
		SourceBranch string `json:"source_branch"`
		State        string `json:"state"`
	} `json:"object_attributes"`
}

// GitlabHandler will handle the payload from GitLab and remove containers
// based on a domain if they exists where the issue key is the prefix in the domain.
func GitlabHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p gitlabPayload
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

	// Check merge request state.
	if strings.ToLower(p.ObjectAttributes.State) != "merged" {
		fmt.Fprintf(w, "Only merge requested with state `merged` can be handle.")
		return
	}

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", strings.ToLower(p.ObjectAttributes.SourceBranch), domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprintf(w, "Removed %d containers with name suffix %s", count, suffix)
}
