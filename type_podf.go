package main

import (
	"gopkg.in/mgo.v2/bson"
)

type ProjectPodfileResponse struct {
	Errno             int                  `json:"errno,omitempty"`
	Msg               string               `json:"msg,omitempty"`
	Version           int64                `json:"version"`
	NeedUpgradVersion int64                `json:"need_upgrade_version"`
	Podfile           *ProjectPodfileModel `json:"podfile,omitempty"`
}

type ProjectPodfileModel struct {
	Id              bson.ObjectId     `json:"_id,omitempty" bson:"_id,omitempty"`
	Team            string            `json:"team,omitempty" bson:"team,omitempty"`
	Project         string            `json:"project,omitempty" bson:"project,omitempty"`
	Version         int64             `json:"version" bson:"version"`
	TemplateVersion int64             `json:"template_version,omitempty" bson:"template_version"`
	Release         bool              `json:"release" bson:"release"`
	Prefix          string            `json:"prefix,omitempty" bson:"prefix,omitempty"`
	Suffix          string            `json:"suffix,omitempty" bson:"suffix,omitempty"`
	Target          string            `json:"target,omitempty" bson:"target,omitempty"`
	Hierarchies     []*HierarchyModel `json:"hierarchies,omitempty" bson:"hierarchies,omitempty"`
	CreateUser      string            `json:"create_user,omitempty" bson:"create_user,omitempty"`
	CreateTime      string            `json:"create_time,omitempty" bson:"create_time,omitempty"`
	RelaseUser      string            `json:"release_user,omitempty" bson:"release_user,omitempty"`
	ReleaseTIme     string            `json:"release_time,omitempty" bson:"release_time,omitempty"`
	Tags            []string          `json:"tags,omitempty" bson:"tags,omitempty"`
}

type HierarchyModel struct {
	Name             string       `json:"name" bson:"name"`
	AllowInterDepend bool         `json:"allow_interdepend" bson:"allow_interdepend"`
	Pods             []*PodModule `json:"pods,omitempty" bson:"pods,omitempty"`
	ImplicitPods     []*PodModule `json:"implicit_pods,omitempty" bson:"implicit_pods,omitempty"`
}

type PodModule struct {
	Name        string       `json:"name,omitempty" bson:"name,omitempty"`
	Version     string       `json:"version,omitempty" bson:"version,omitempty"`
	VersionType int          `json:"version_type" bson:"version_type"`
	Subspecs    []string     `json:"subspecs,omitempty" bson:"subspecs,omitempty"`
	Description string       `json:"description,omitempty" bson:"description,omitempty"`
	Addition    *PodAddition `json:"addition,omitempty" bson:"description,omitempty"`
}

type PodAddition struct {
	ReferName           string `json:"refer_name,omitempty"`
	ReferVersion        string `json:"refer_version,omitempty"`
	ReferModuleViersion string `json:"refer_module_version,omitempty"`
	NewestVersion       string `json:"newest_version,omitempty"`
}

type ProjectProfileDiffResponse struct {
	Errno int                      `json:"errno,omitempty"`
	Msg   string                   `json:"msg,omitempty"`
	Diff  *ProjectProfileDiffModel `json:"diff,omitempty"`
}

type ProjectProfileDiffModel struct {
	Change map[string][]string `json:"change,omitempty"`
	New    map[string]string   `json:"new,omitempty"`
	Remove map[string]string   `json:"remove,omitempty"`
}
