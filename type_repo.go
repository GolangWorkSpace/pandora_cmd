package main

import "gopkg.in/mgo.v2/bson"

type RepoResp struct {
	Errno int `json:"errno,omitempty" bson:"errno,omitempty"`
	Msg string `json:"msg,omitempty" bson:"msg,omitempty"`
	TaskCount int `json:"task_count,omitempty" bson:"task_count,omitempty"`
	Repos []*RepoModel `json:"repos,omitempty" bson:"repos,omitempty"`
}

type RepoModel struct {
	Id          bson.ObjectId     `json:"_id,omitempty" bson:"_id,omitempty"`
	RepoName    string            `json:"repo_name,omitempty" bson:"repo_name,omitempty"`
	Git         string            `json:"git,omitempty" bson:"git,omitempty"`
	Cycle       int               `json:"cycle" bson:"cycle"`
	Exclude     map[string]string `json:"exclude,omitempty" bson:"exclude,omitempty"`
	Description string            `json:"description,omitempty" bson:"description,omitempty"`
	SyncTime    string            `json:"sync_time" bson:"sync_time,omitempty"`

	Addition interface{} `json:"addition,omitempty" bson:"-"`
}