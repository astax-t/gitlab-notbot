package main

import (
	"net/http"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
	"github.com/joho/godotenv"
)

var config Config
var git *gitlab.Client

func main() {
	PrepareConfig()

	PrepareGitLabClient()

	StartServer()
}

func PrepareConfig()  {
	config.loadDefault()
	err := godotenv.Load()
	if err != nil {
		log(LOG_MESSAGE, "Can't load .env file", err)
	}

	err = config.populate()
	if err != nil {
		log(LOG_FATAL, "Configuration error: ", err)
	}
}

func StartServer()  {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		data := whData{}
		err := data.Prepare(request.Body, request.ContentLength)
		if err != nil {
			return
		}
		go ProcessRequest(&data)
	})

	host := config.ListenHost
	if host == "*" {
		host = ""
	}
	address := fmt.Sprintf("%v:%d", host, config.ListenPort)

	log(LOG_INFO, fmt.Sprintf("Starting server at %v", address), nil)

	log(LOG_FATAL, "", http.ListenAndServe(address, nil))
}

func ProcessRequest(data *whData) {
	kind := data.GetKind()
	if kind != "issue" {
		log(LOG_INFO, "Unsupported object in webhook - " + kind, nil)
		return
	}
	log(LOG_DEBUG, "Object kind - " + kind, nil)

	/*if !data.LabelsChanged() {
		log(LOG_INFO, "Labels are unchanged, no need to do anything", nil)
		return
	}
	log(LOG_DEBUG, "Labels have been changed, processing the list", nil)/**/

	labels, err := data.GetLabels()
	if err != nil {
		return
	}
	log(LOG_DEBUG, fmt.Sprintf("Current labels %v", labels), nil)

	issue, err := data.GetIssue()
	if err != nil {
		return
	}
	log(LOG_DEBUG, fmt.Sprintf("Issue: %v", issue), nil)

	project, err := data.GetProject()
	if err != nil {
		return
	}
	log(LOG_DEBUG, fmt.Sprintf("Project: %v", project), nil)

	projectLabels, err := GetAvailableLabels(project)
	if err != nil {
		return
	}
	log(LOG_DEBUG, fmt.Sprintf("Project labels %v", projectLabels), nil)

	labelsToSet := ComputeLabelsToSet(labels, projectLabels)
	log(LOG_DEBUG, fmt.Sprintf("Computed labels to set: %v", labelsToSet), nil)

	CreateMissingLabels(project, labelsToSet, projectLabels)

	SetIssueLabels(project, issue, labelsToSet)

	log(LOG_INFO, fmt.Sprintf("Successfully set labels for issue '%v'", issue.Title), nil)
}

func ComputeLabelsToSet(current, all []label) []label {
	result := make([]label, 0, len(all)/2)

	var exists bool

	currentMap := make(map[string]label)
	for _, l := range current {
		if IsNegativeLabel(&l.Title) {
			continue
		}
		currentMap[l.Title] = l
		// keep all existing non-negative labels
		result = append(result, l)
	}

	for _, al := range all {
		if IsNegativeLabel(&al.Title) {
			continue
		}
		if _, exists = currentMap[al.Title]; !exists {
			result = append(result, label{Title: config.LabelPrefix + al.Title})
		}
	}

	return result
}

func IsNegativeLabel(name *string) bool {
	return strings.HasPrefix(*name, config.LabelPrefix)
}

