package main

type TemplateResp struct {
	Errno     int `json:"errno,omitempty" bson:"errno,omitempty"`
	Msg       string `json:"msg,omitempty" bson:"msg,omitempty"`
	Template  *TemplateModel `json:"template,omitempty" bson:"template,omitempty"`
	Templates []*TemplateModel `json:"templates,omitempty" bson:"templates,omitempty"`
}

type TemplateModel struct {
	Team           string `json:"team,omitempty" bson:"team,omitempty"`
	Project        string `json:"project,omitempty" bson:"project,omitempty"`
	Version        int64 `json:"version" bson:"version"`
	ReferName      string `json:"refer_name,omitempty" bson:"refer_name,omitempty"`
	ReferVersion   string `json:"refer_version,omitempty" bson:"refer_version,omitempty"`
	ReferVersionId int `json:"refer_version_id,omitempty" bson:"refer_version_id,omitempty"`
	Prefix         string `json:"prefix,omitempty" bson:"prefix,omitempty"`
	Suffix         string `json:"suffix,omitempty" bson:"suffix,omitempty"`
	Target         string `json:"target,omitempty" bson:"target,omitempty"`
	Hierarchies    []*HierarchyModel `json:"hierarchies,omitempty" bson:"hierarchies,omitempty"`
	CreateUser     string `json:"create_user,omitempty" bson:"create_user,omitempty"`
	CreateTime     string `json:"create_time,omitempty" bson:"create_time,omitempty"`
}
