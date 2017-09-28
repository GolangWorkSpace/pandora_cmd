package main

import (
	"errors"
	"bytes"
	"strconv"
)

func RepoList(parentMenu *CmdMenu)  {
	repos, err := RepoRequestList()
	if err != nil {
		PrintThenExit(err.Error())
	}
	println("仓库列表:")
	for _, aRepo := range repos {
		println("  -", aRepo.RepoName, " Git:", aRepo.Git)
	}
}

func RepoSync(parentMenu *CmdMenu)  {
	PrintlnYellow("Tips: 同步仓库需要一定时间，请您耐心等待!")
	aRepo, err := RepoSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}
	println("正在同步仓库"+aRepo.RepoName+" ...")
	url := GenURL("/api/repo/module/block_sync")
	param := []string{aRepo.RepoName}
	var res *RepoResp
	if err := POSTParse(url, param, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.Errno != 0 {
		PrintlnErrorFormat("同步仓库失败", res.Msg, res.Errno)
	}
	println("同步仓库"+aRepo.RepoName+"成功，任务数:"+strconv.Itoa(res.TaskCount))
}

func RepoSelect() (*RepoModel, error)  {
	repos, err := RepoRequestList()
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString("请选择仓库:\n")
	for idx, aRepo := range repos {
		buffer.WriteString(strconv.Itoa(idx+1)+"."+aRepo.RepoName+"  Git: "+aRepo.Git + "\n")
	}
	buffer.WriteString(":")
	selected := SimpleInputSelectNum(buffer.String(), 1, len(repos))
	return repos[selected], nil
}

func RepoRequestList() ([]*RepoModel, error) {
	url := GenURL("/api/repo/list")
	var res *RepoResp
	if err := GETParse(url, &res); err != nil {
		return nil, err
	}
	if res.Errno != 0 {
		return nil, errors.New(FormatResError("查询仓库列表失败", res.Msg, res.Errno))
	}
	if len(res.Repos) == 0 {
		return nil, errors.New("没有仓库!")
	}
	return res.Repos, nil
}