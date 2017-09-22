package main

import "gopkg.in/mgo.v2/bson"

type ProjectResp struct {
	Errno int `json:"errno,omitempty" bson:"errno,omitempty"`
	Msg string `json:"msg,omitempty" bson:"msg,omitempty"`
	Project *ProjectModel `json:"project,omitempty" bson:"project,omitempty"`
	Projects []*ProjectModel `json:"projects,omitempty" bson:"projects,omitempty"`
}

type ProjectModel struct {
	Id          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Team        string        `json:"team,omitempty" bson:"team,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	Git         string        `json:"git,omitempty" bson:"git,omitempty"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`
}

