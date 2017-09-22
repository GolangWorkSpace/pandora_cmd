package main

import (
	"github.com/go-hayden-base/foundation"
	"strconv"
	"errors"
	"bytes"
)

func TeamList(parentMenu *CmdMenu) {
	teams, err := TeamRequestList()
	if err != nil {
		PrintlnRed(err.Error())
		return
	}
	if len(teams) == 0 {
		PrintlnRed("您未创建任何团队并且您不属于任何团队!")
		return
	}
	println("您的团队如下:")
	for _, aTeam := range teams {
		println("  - 团队名称:", aTeam.Name, " 您的团队角色:", aTeam.RoleName())
	}
}

func TeamCreate(parentMenu *CmdMenu) {
	for {
		teamName := ""
		foundation.IArgInput("请输入团队名称(建议用纯英文):", func(arg string) foundation.IArgAction {
			if arg == "" {
				PrintlnRed("输入错误，请重新输入!")
				return foundation.IArgActionRepet
			}
			teamName = arg
			return foundation.IArgActionNext
		})
		teamDesc := ""
		foundation.IArgInput("请输入团队描述:", func(arg string) foundation.IArgAction {
			teamDesc = arg
			return foundation.IArgActionNext
		})

		url := GenURL("/api/team/add")
		param := map[string]string{
			"name":        teamName,
			"description": teamDesc,
		}
		var res *Response
		if err := POSTParse(url, param, &res); err != nil {
			PrintlnRed("添加团队发生错误:", err.Error())
			PrintlnRed("请重试!")
			continue
		}
		if res.Errno != 0 {
			PrintlnRed("添加团队发生错误，错误码:", strconv.Itoa(res.Errno), "信息:", res.Msg)
			PrintlnRed("请重试!")
			continue
		}
		println("添加团队成功!")
		break
	}
}

func TeamMembers(parentMenu *CmdMenu) {
	println("")
	aTeam, err := TeamSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}
	if aTeam == nil {
		PrintThenExit("您没有选择团队，程序退出！")
	}
	println("\n" + aTeam.Name + "的团队成员:")
	println("  - 帐号:"+aTeam.Agent, "角色:"+RoleNameMap[RoleAgent])
	for _, aMember := range aTeam.Member {
		println("  - 帐号:"+aMember.Account, "角色:"+aMember.RoleName())
	}
}

func TeamMemberAdd(parentMenu *CmdMenu) {
	println("")
	var aTeam *TeamModel
	var err error
	for {
		aTeam, err = TeamSelect()
		if err != nil {
			PrintThenExit(err.Error())
		}
		if aTeam == nil {
			PrintThenExit("您没有选择团队，程序退出！")
		}
		if !aTeam.CheckRole(RoleAdmin) {
			PrintlnRed("添加成员需要Admin及以上的权限，您没有为团队" + aTeam.Name + "添加成员的权限! 请重新选择!")
			continue
		}
		break
	}
	println("已选团队:" + aTeam.Name)
	account := ""
	foundation.IArgInput("请输入要添加的成员帐号:", func(arg string) foundation.IArgAction {
		if arg == "" {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		if aTeam.MemeberExist(arg) {
			PrintlnRed("该用户已经是团队成员，请重新输入!")
			return foundation.IArgActionRepet
		}
		account = arg
		return foundation.IArgActionNext
	})
	role := UserSelectRole()
	aNewMember := new(TeamMember)
	aNewMember.Account = account
	aNewMember.Role = role

	url := GenURL("/api/team/update/member/force")
	if aTeam.Member != nil {
		aTeam.Member = append(aTeam.Member, aNewMember)
	} else {
		aTeam.Member = []*TeamMember{aNewMember}
	}
	aTeam.Description = ""
	aTeam.Id = ""
	var res *Response
	if err := POSTParse(url, aTeam, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.Errno != 0 {
		PrintlnErrorFormat("添加团队成员失败", res.Msg, res.Errno)
		return
	}
	println("添加团队成员成功！")
}

func TemMemberModifyRole(parentMenu *CmdMenu)  {
	println("")
	var aTeam *TeamModel
	var err error
	for {
		aTeam, err = TeamSelect()
		if err != nil {
			PrintThenExit(err.Error())
		}
		if aTeam == nil {
			PrintThenExit("您没有选择团队，程序退出！")
		}
		if !aTeam.CheckRole(RoleAdmin) {
			PrintlnRed("添加成员需要Admin及以上的权限，您没有为团队" + aTeam.Name + "添加成员的权限! 请重新选择!")
			continue
		}
		break
	}
	println("已选团队:" + aTeam.Name)
	l := len(aTeam.Member)
	if l == 0 {
		PrintlnRed("该团队没有成员，请先添加成员！")
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("请选择成员:\n")
	for idx, aMember := range aTeam.Member {
		buffer.WriteString(strconv.Itoa(idx+1)+"."+aMember.Account+"\n")
	}
	buffer.WriteString(":")
	selected := 0
	foundation.IArgInput(buffer.String(), func(arg string) foundation.IArgAction {
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < 1 || idx > 3 {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		selected = idx - 1
		return foundation.IArgActionNext
	})
	aMember := aTeam.Member[selected]
	role := UserSelectRole()
	aMember.Role = role

	url := GenURL("/api/team/update/member/force")
	aTeam.Description = ""
	aTeam.Id = ""
	var res *Response
	if err := POSTParse(url, aTeam, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.Errno != 0 {
		PrintlnErrorFormat("修改成员角色失败", res.Msg, res.Errno)
		return
	}
	println("修改成员角色成功！")
}

func TeamMemberRemove(parentMenu *CmdMenu)  {
	println("")
	var aTeam *TeamModel
	var err error
	for {
		aTeam, err = TeamSelect()
		if err != nil {
			PrintThenExit(err.Error())
		}
		if aTeam == nil {
			PrintThenExit("您没有选择团队，程序退出！")
		}
		if !aTeam.CheckRole(RoleAdmin) {
			PrintlnRed("添加成员需要Admin及以上的权限，您没有为团队" + aTeam.Name + "添加成员的权限! 请重新选择!")
			continue
		}
		break
	}
	println("已选团队:" + aTeam.Name)
	l := len(aTeam.Member)
	if l == 0 {
		PrintlnRed("该团队没有成员，请先添加成员！")
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("请选择成员:\n")
	for idx, aMember := range aTeam.Member {
		buffer.WriteString(strconv.Itoa(idx+1)+"."+aMember.Account+"\n")
	}
	buffer.WriteString(":")
	selected := 0
	foundation.IArgInput(buffer.String(), func(arg string) foundation.IArgAction {
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < 1 || idx > 3 {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		selected = idx - 1
		return foundation.IArgActionNext
	})
	aTeam.Member = append(aTeam.Member[:selected], aTeam.Member[selected+1:]...)

	url := GenURL("/api/team/update/member/force")
	aTeam.Description = ""
	aTeam.Id = ""
	var res *Response
	if err := POSTParse(url, aTeam, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.Errno != 0 {
		PrintlnErrorFormat("删除成员角色失败", res.Msg, res.Errno)
		return
	}
	println("删除成员角色成功！")
}

func TeamSelect() (*TeamModel, error) {
	teams, err := TeamRequestList()
	if err != nil {
		return nil, err
	}
	l := len(teams)
	if l == 0 {
		return nil, errors.New("您没有在任何团队，您可以创建一个团队或申请加入一个团队!")
	}
	var buffer bytes.Buffer
	buffer.WriteString("请选择团队:\n")
	for idx, aTeam := range teams {
		buffer.WriteString(strconv.Itoa(idx+1) + ".团队名称:" + aTeam.Name + " 您的团队角色:" + aTeam.RoleName() + "\n")
	}
	buffer.WriteString(":")
	selected := 0
	foundation.IArgInput(buffer.String(), func(arg string) foundation.IArgAction {
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < 1 || idx > l {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		selected = idx - 1
		return foundation.IArgActionNext
	})
	return teams[selected], nil
}

func TeamRequestList() ([]*TeamModel, error) {
	url := GenURL("/api/team/list")
	var teamResp *TeamResp
	if err := GETParse(url, &teamResp); err != nil {
		return nil, errors.New(FormtError("查询团队列表失败", err))
	}

	if teamResp.Errno != 0 {
		return nil, errors.New(FormatResError("查询团队列表失败", teamResp.Msg, teamResp.Errno))
	}

	return teamResp.Teams, nil
}
