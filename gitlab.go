package main

import (
	"github.com/xanzy/go-gitlab"
	"fmt"
	"net/http"
)

func PrepareGitLabClient()  {
	git = gitlab.NewClient(nil, config.GitLabToken)
	git.SetBaseURL(config.GitLabUrl + "/api/v4")
}

func GetAvailableLabels(project project) ([]label, error) {
	// List all labels
	labels, _, err := git.Labels.ListLabels(project.Id, func(req *http.Request) error {
		if req.URL.RawQuery != "" {
			req.URL.RawQuery = req.URL.RawQuery + "&per_page=100"
		} else {
			req.URL.RawQuery = "per_page=100"
		}
		return nil
	})
	if err != nil {
		log(LOG_ERROR, "error reading project labels", err)
		return nil, err
	}

	// convert labels to internal simpler format
	result := make([]label, 0, len(labels))
	for _, l := range labels {
		result = append(result, label{Title:l.Name})
	}

	return result, nil
}

func CreateMissingLabels(project project, required, current []label) {
	for i := range required {
		exists := false
		for j := range current {
			if current[j].Title == required[i].Title {
				exists = true
				break
			}
		}

		if !exists {
			log(LOG_INFO, "Creating new project label: " + required[i].Title, nil)
		} else
		{
			log(LOG_DEBUG, "Project label already exists: " + required[i].Title, nil)
			continue
		}

		labelInfo := &gitlab.CreateLabelOptions{
			Name:  gitlab.String(required[i].Title),
			Color: gitlab.String(config.LabelColor),
			Description: gitlab.String("Auto-created by NotBot"),
		}
		newLabel, _, err := git.Labels.CreateLabel(project.Id, labelInfo)
		if err != nil {
			log(LOG_ERROR, "Error while creating label", err)
		} else {
			log(LOG_DEBUG, "New project label successfully created: " + newLabel.Name, nil)
		}
	}
}

func SetIssueLabels(project project, issue issue, labels []label) {
	log(LOG_DEBUG, fmt.Sprintf("Setting labels for issue %v : %v ", issue.Iid, labels), nil)

	labelsStr := make([]string, 0, len(labels))
	for i := range labels {
		labelsStr = append(labelsStr, labels[i].Title)
	}

	issueInfo := &gitlab.UpdateIssueOptions{
		Labels: labelsStr,
	}
	_, _, err := git.Issues.UpdateIssue(project.Id, issue.Iid, issueInfo)
	if err != nil {
		log(LOG_ERROR, "Error while setting issue labels", err)
	} else {
		log(LOG_DEBUG, "Issue labels successfully updated", nil)
	}
}
