package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/isotopsweden/docker-farmer/config"
	"github.com/isotopsweden/docker-farmer/docker"
)

// jiraPayload contains the payload from Jira.
type jiraPayload struct {
	Key string `json:"key"`
}

// JiraHandler will handle the payload from JIRA and remove containers
// based on a domain if they exists where the issue key is the prefix in the domain.
func JiraHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var p jiraPayload
	if err := decoder.Decode(&p); err != nil {
		write(w, fmt.Sprintf("Jira: %s", err.Error()))
		return
	}

	// Get the domain.
	domain := config.Get().Domain

	// Remove containers for the suffix.
	suffix := fmt.Sprintf("%s.%s", strings.ToLower(p.Key), domain)
	count, err := docker.RemoveContainers(suffix)
	if err != nil {
		write(w, fmt.Sprintf("Jira: %s", err.Error()))
		return
	}

	write(w, fmt.Sprintf("Jira: Removed %d containers with name suffix %s", count, suffix))
}
