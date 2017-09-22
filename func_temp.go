package main

import (
	"strconv"
	"errors"
	"bytes"
	"encoding/json"
	"github.com/go-hayden-base/fs"
	"path/filepath"
	"io/ioutil"
	"os"
)

func TemplateList(parentMenu *CmdMenu) {
	_, aProject, templates, err := TemplateRequestList()
	if err != nil {
		PrintlnRed(err.Error())
		return
	}
	println("项目" + aProject.Name + "模板列表如下:")
	for _, aTemplate := range templates {
		line := " - 版本:" + strconv.FormatInt(aTemplate.Version, 10)
		line += " 参照:" + aTemplate.ReferName + "(" + aTemplate.ReferVersion + ")"
		line += " 创建人:" + aTemplate.CreateUser + " 创建时间:" + aTemplate.CreateTime
		println(line)
	}
}

func TemplateShowOne(parentMenu *CmdMenu) {
	_, _, aTemplate, err := TemplateSelect()
	if err != nil {
		PrintlnRed(err.Error())
		return
	}
	b, err := json.MarshalIndent(aTemplate, "", "  ")
	if err != nil {
		PrintlnRed(err.Error())
		return
	}
	alert := "请选择操作:\n1.输出到终端\n2.输出到文件\n:"
	selected := SimpleInputSelectNum(alert, 1, 2)
	if selected == 0 {
		println(string(b))
		return
	}

	for {
		outputPath := SimpleInputString("请输入文件输出路径:", false)
		if !filepath.IsAbs(outputPath) {
			outputPath = filepath.Join(_Config.CurrentDir, outputPath)
		}
		if fs.FileExists(outputPath) {
			confirm := SimpleInputString("文件存在["+outputPath+"], 确认覆盖吗?(y/n):", true)
			if confirm != "y" {
				continue
			}
		}
		if err := ioutil.WriteFile(outputPath, b, os.ModePerm); err != nil {
			PrintlnError("输出文件发生错误", err)
			break
		}
		println("模板已输出到文件: " + outputPath)
		break
	}
}

func TemplateFollow(parentMenu *CmdMenu) {
	PrintlnYellow("Tips:模板跟进需要有Developer及以上的团队权限!")
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		PrintlnRed(err.Error())
		return
	}

	if !aTeam.CheckRole(RoleDeveloper) {
		PrintlnRed("模板跟进需要有Developer及以上的团队权限!")
		return
	}

	PrintlnYellow("开始跟进项目" + aProject.Name + "的模板...")

	aRefer, err := ReferSelect()
	if err != nil {
		PrintlnRed(err.Error())
		return
	}

	url := GenURLWithParam("/api/logic/template/follow", map[string]string{
		"team":          aTeam.Name,
		"project":       aProject.Name,
		"refer_name":    aRefer.ReferName,
		"refer_version": aRefer.ReferVersion,
		"version_id":    strconv.Itoa(aRefer.VersionId),
	})
	var res *Response
	if err := GETParse(url, &res); err != nil {
		PrintlnError("模板跟进失败", err)
		return
	}
	if res.Errno != 0 {
		PrintlnErrorFormat("模板跟进失败", res.Msg, res.Errno)
		return
	}
	println("跟进成功，模板已跟进到：" + aRefer.ReferName + " " + aRefer.ReferVersion + ", 请重新生成Podfile！")
}

func TemplateSelect() (*TeamModel, *ProjectModel, *TemplateModel, error) {
	aTeam, aProject, templates, err := TemplateRequestList()
	if err != nil {
		return nil, nil, nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString("请选择模板:\n")
	for idx, aTemplate := range templates {
		buffer.WriteString(strconv.Itoa(idx+1) + ".版本:" + strconv.FormatInt(aTemplate.Version, 10) + " ")
		buffer.WriteString("参照:" + aTemplate.ReferName + "(" + aTemplate.ReferVersion + ") ")
		buffer.WriteString("创建人:" + aTemplate.CreateUser + " 创建时间:" + aTemplate.CreateTime + "\n")
	}
	buffer.WriteString(":")
	selected := SimpleInputSelectNum(buffer.String(), 1, len(templates))

	aTemaplate := templates[selected]
	aReqTemplate, err := TemplateRequestOne(aTemaplate.Team, aTemaplate.Project, aTemaplate.Version)
	if err != nil {
		return nil, nil, nil, err
	}
	return aTeam, aProject, aReqTemplate, nil
}

func TemplateRequestOne(team, project string, version int64) (*TemplateModel, error) {
	url := GenURLWithParam("/api/logic/template/one", map[string]string{
		"team":    team,
		"project": project,
		"version": strconv.FormatInt(version, 10),
	})
	var res *TemplateResp
	if err := GETParse(url, &res); err != nil {
		return nil, err
	}
	if res.Errno != 0 {
		return nil, errors.New(FormatResError("查询模板失败", res.Msg, res.Errno))
	}
	return res.Template, nil
}

func TemplateRequestList() (*TeamModel, *ProjectModel, []*TemplateModel, error) {
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		return nil, nil, nil, err
	}
	limit := SimpleInputInt("请输入要获取的模板条数(最大为100，最小为1，默认20):", 1, 100, 20)
	url := GenURLWithParam("/api/logic/template/list", map[string]string{
		"team":    aTeam.Name,
		"project": aProject.Name,
		"limit":   strconv.Itoa(limit),
	})
	println(url)
	var res *TemplateResp
	if err := GETParse(url, &res); err != nil {
		return nil, nil, nil, err
	}
	if res.Errno != 0 {
		return nil, nil, nil, errors.New(FormatResError("查询模板失败", res.Msg, res.Errno))
	}
	if len(res.Templates) == 0 {
		return nil, nil, nil, errors.New("没有模板！")
	}
	return aTeam, aProject, res.Templates, nil
}
