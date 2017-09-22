package main

type ReferResp struct {
	Errno  int `json:"errno,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Refers []*ReferVersionModel `json:"refs,omitempty"`
}

type ReferVersionModel struct {
	ReferName    string `json:"refer_name,omitempty" bson:"refer_name,omitempty"`
	ReferVersion string `json:"refer_version,omitempty" bson:"refer_version,omitempty"`
	VersionId    int `json:"version_id,omitempty" bson:"version_id,omitempty"`
	Sha          string `json:"sha,omitempty" bson:"sha,omitempty"`
	Timestamp    int64 `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	TimeString   string `json:"time_string,omitempty" bson:"time_string,omitempty"`
	PodModules   map[string]*PodModule `json:"pod_modules,omitempty" bson:"pod_modules,omitempty"`
}
