package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/isotopsweden/docker-farmer/config"
	"github.com/isotopsweden/docker-farmer/docker"
)

// ConfigHandler will handle the config api route.
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(config.Get())

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		fmt.Fprintf(w, string(j))
	}
}

// ContainersHandler will handle the containers api route.
func ContainersHandler(w http.ResponseWriter, r *http.Request) {
	c := config.Get()
	containers, err := docker.GetContainers(c.Domain)
	all := r.URL.Query().Get("all")

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		exclude := config.Get().Sites.Exclude
		list := []types.Container{}

		for _, c := range containers {
			if all != "true" && stringInSlice(c.Image, exclude) {
				continue
			}

			list = append(list, c)
		}

		j, err := json.Marshal(list)

		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, string(j))
		}
	}
}

// DatabaseHandler will handle the database api route.
func DatabaseHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	typ := r.URL.Query().Get("type")

	if name == "" {
		fmt.Fprintf(w, "No name query string")
		return
	}

	if typ == "" {
		typ = "mysql"
	}

	conf := config.Get()

	ok := false
	var err error

	switch typ {
	case "mysql":
		ok, err = docker.DeleteMySQLDatabase(conf.Database.User, conf.Database.Password, conf.Database.Prefix, name, conf.Database.Container)
		break
	default:
		break
	}

	if err != nil {
		fmt.Fprintf(w, err.Error())
	} else {
		j, err := json.Marshal(map[string]bool{
			"success": ok,
		})

		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, string(j))
		}
	}
}
