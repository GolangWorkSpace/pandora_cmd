package main

import "gopkg.in/mgo.v2/bson"

const (
	RoleUnknown   = Role(0)
	RoleViewer    = Role(1)
	RoleDeveloper = Role(2)
	RoleAdmin     = Role(3)

	RoleAgent = Role(9)
)

var RoleNameMap = map[Role]string{
	RoleViewer:    "Viewer",
	RoleDeveloper: "Developer",
	RoleAdmin:     "Admin",
	RoleAgent:     "Agent",
}

var SelectRoles = []*RoleInfo {
	&RoleInfo{
		Role:RoleViewer,
		Name:"Viewer",
	},
	&RoleInfo{
		Role:RoleDeveloper,
		Name:"Developer",
	},
	&RoleInfo{
		Role:RoleAdmin,
		Name:"Admin",
	},
}

type RoleInfo struct {
	Role Role
	Name string
}

type Role uint8

type UserResp struct {
	Errno int `json:"errno,omitempty" bson:"errno,omitempty"`
	Msg   string `json:"msg,omitempty" bson:"msg,omitempty"`
	User  *UserModel `json:"user,omitempty" bson:"user,omitempty"`
}

type UserModel struct {
	Id       bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Account  string        `json:"account,omitempty" bson:"account,omitempty"`
	Password string        `json:"password,omitempty" bson:"password,omitempty"`
	CMDToken string        `json:"cmd_token,omitempty" bson:"cmd_token,omitempty"`
	Name     string        `json:"name,omitempty" bson:"name,omitempty"`
	Role     Role          `json:"role" bson:"role"`
}