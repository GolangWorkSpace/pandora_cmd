package main

import (
	"errors"
	"bytes"
	"strconv"
	"github.com/go-hayden-base/foundation"
	"os"
)

func ProjectList(parentMenu *CmdMenu)  {
	projects, aTeam, err := ProjectRequestList()
	if err != nil {
		PrintThenExit(err.Error())
	}
	println("团队"+aTeam.Name+"有如下项目:")
	for _, aProject := range projects {
		println("  -", aProject.Name, "("+aProject.Git+")")
	}
}

func ProjectCreate(parentMenu *CmdMenu) {
	PrintlnYellow("Tips: 必须具有团队Admin权限才能添加项目")
	var aTeam *TeamModel
	var err error
	for {
		aTeam, err = TeamSelect()
		if err != nil {
			PrintThenExit(err.Error())
		}
		if aTeam.CheckRole(RoleAdmin) {
			break
		}
		PrintlnRed("您没有为团队"+aTeam.Name+"添加项目的权限，需要Admin及以上权限！请重新选择！")
	}

	projectName := SimpleInputString("请输入项目名称(建议英文和数字的组合):", false)
	projectGit := SimpleInputString("请输入项目Git地址:", false)
	projectDesc := SimpleInputString("请输入项目描述(选填):", true)

	param := map[string]string{
		"team":        aTeam.Name,
		"name":        projectName,
		"git":         projectGit,
		"description": projectDesc,
	}
	url := GenURL("/api/project/add")
	var res *ProjectResp
	if err := POSTParse(url, param, &res); err != nil {
		PrintlnError("添加项目失败", err)
		return
	}
	if res.Errno != 0 {
		PrintlnErrorFormat("添加项目失败", res.Msg, res.Errno)
		os.Exit(1)
	}
	println("添加项目成功！")
}

func ProjectSelect() (*ProjectModel, *TeamModel, error) {
	projects, aTeam, err := ProjectRequestList()
	if err != nil {
		return nil, nil, err
	}
	l := len(projects)
	if l == 0 {
		return nil, nil, errors.New("团队" + aTeam.Name + "暂无项目，请添加项目！")
	}
	var buffer bytes.Buffer
	buffer.WriteString("请选择项目:\n")
	for idx, aProject := range projects {
		buffer.WriteString(strconv.Itoa(idx+1) + ".项目名称:" + aProject.Name + " Git地址:" + aProject.Git + "\n")
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
	return projects[selected], aTeam, nil
}

func ProjectRequestList() ([]*ProjectModel, *TeamModel, error) {
	aTeam, err := TeamSelect()
	if err != nil {
		return nil, nil, err
	}

	url := GenURL("/api/project/list", "team", aTeam.Name)
	var res *ProjectResp
	if err := GETParse(url, &res); err != nil {
		return nil, nil, err
	}
	if res.Errno != 0 {
		return nil, nil, errors.New(FormatResError("获取项目失败", res.Msg, res.Errno))
	}
	if len(res.Projects) == 0 {
		return nil, nil, errors.New("团队" + aTeam.Name + "没有任何项目!")
	}
	return res.Projects, aTeam, nil
}
