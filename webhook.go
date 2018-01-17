package main

import (
	"encoding/json"
	"io"
	"errors"
	"fmt"
)

type whData struct {
	data map[string]*json.RawMessage
}

type label struct {
	Id int `json:"id,int"`
	Title string `json:"title"`
//	ProjectId int `json:"project_id,int"`
//	Description string `json:"description"`
//	Type string `json:"type"`   // can be "GroupLabel" or "ProjectLabel"
//	GroupId int `json:"group_id,int"`
}

type project struct {
	Id int `json:"id,int"`
	Path string `json:"path_with_namespace"`
}

type issue struct {
	Id int `json:"id,int"`
	Iid int `json:"iid,int"`
	Title string `json:"title"`
}

type user struct {
	Username string `json:"username"`
}

var jsonError = errors.New("JSON decoding error")
var dataError = errors.New("webhook data error")

func (wh *whData) Prepare(body io.Reader, length int64) error {
	log(LOG_DEBUG, fmt.Sprintf("Request length - %v bytes", length), nil)
	bodyStr := make([]byte, length)
	readLen, err := io.ReadFull(body, bodyStr)
	if err != nil {
		log(LOG_ERROR, "Error while reading request data", err)
		log(LOG_ERROR, fmt.Sprintf("Mismatching data length: expected %v got %v", length, readLen), err)
		return dataError
	}
	//log(LOG_DEBUG, "Full request:\n" + string(bodyStr), nil)

	err = json.Unmarshal(bodyStr[:length], &wh.data)
	if err != nil {
		log(LOG_ERROR, "Top level JSON decoding error", err)
		return jsonError
	}

	return nil
}

func (wh whData) GetKind() (objectKind string) {
	err := json.Unmarshal(*wh.data["object_kind"], &objectKind)
	if err != nil {
		log(LOG_ERROR, "JSON error while getting object kind", err)
		return ""
	}

	return
}

func (wh whData) LabelsChanged() bool {
	var changes map[string]*json.RawMessage
	err := json.Unmarshal(*wh.data["changes"], &changes)
	if err != nil {
		log(LOG_ERROR, "JSON error while getting list of changes", err)
		return false
	}

	_, exists := changes["labels"]
	return exists
}

func (wh whData) GetLabels() (labels []label, err error) {
	err = json.Unmarshal(*wh.data["labels"], &labels)
	if err != nil {
		log(LOG_ERROR, "JSON error while getting list of labels", err)
		return labels, jsonError
	}

	return labels, nil
}

func (wh whData) GetProject() (project project, err error) {
	err = json.Unmarshal(*wh.data["project"], &project)
	if err != nil {
		log(LOG_ERROR, "JSON error while getting project details", err)
		return project, jsonError
	}

	return project, nil
}

func (wh whData) GetIssue() (issue issue, err error) {
	err = json.Unmarshal(*wh.data["object_attributes"], &issue)
	if err != nil {
		log(LOG_ERROR, "JSON error while getting issue details", err)
		return issue, jsonError
	}

	return issue, nil
}

func (wh whData) GetUsername() string {
	var usr user
	err := json.Unmarshal(*wh.data["user"], &usr)
	if err != nil {
		log(LOG_ERROR, "JSON error while getting username of the changer", err)
		return ""
	}

	return usr.Username
}

