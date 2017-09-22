package main

import "gopkg.in/mgo.v2/bson"

type TeamResp struct {
	Errno int `json:"errno,omitempty" bson:"errno,omitempty"`
	Msg string `json:"msg,omitempty" bson:"msg,omitempty"`
	Team *TeamModel `json:"team,omitempty" bson:"team,omitempty"`
	Teams []*TeamModel `json:"teams,omitempty" bson:"team,omitempty"`
}

type TeamModel struct {
	Id          bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string        `json:"name,omitempty" bson:"name,omitempty"`
	Agent       string        `json:"agent,omitempty" bson:"agent,omitempty"`
	Member      []*TeamMember `json:"member" bson:"member"`
	Description string        `json:"description,omitempty" bson:"description,omitempty"`

	Addition interface{} `json:"addition,omitempty" bson:"-"`
}

type TeamMember struct {
	Account string `json:"account,omitempty" bson:"account,omitempty"`
	Role    Role `json:"role,omitempty" bson:"role,omitempty"`
}

func (s *TeamModel)RoleName() string {
	if s.Agent == _Config.Account {
		return RoleNameMap[RoleAgent]
	}
	var aTeamMember *TeamMember
	for _, aMember := range s.Member {
		if aMember.Account == _Config.Account {
			aTeamMember = aMember
			break
		}
	}
	if aTeamMember != nil {
		return aTeamMember.RoleName()
	}
	return "未知"
}

func (s *TeamModel)CheckRole(role Role) bool {
	currentRole := RoleUnknown
	if s.Agent == _Config.Account {
		currentRole = RoleAgent
	} else {
		for _, aMember := range s.Member {
			if aMember.Account == _Config.Account {
				currentRole = aMember.Role
				break
			}
		}
	}
	return currentRole >= role
}

func (s *TeamModel)MemeberExist(account string) bool {
	if s.Agent == account {
		return true
	}
	for _, aMember := range s.Member {
		if aMember.Account == account {
			return true
		}
	}
	return false
}

func (s *TeamMember)RoleName() string  {
	role, ok := RoleNameMap[s.Role]
	if !ok {
		return "未知"
	}
	return role
}
