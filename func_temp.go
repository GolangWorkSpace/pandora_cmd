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
	"strings"
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

func TemplateCreate(parentMenu *CmdMenu) {
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}
	if !aTeam.CheckRole(RoleAdmin) {
		PrintThenExit("创建Podfile模板需要Admin及以上的团队权限！")
	}
	if exists, err := TemplateExist(aTeam.Name, aProject.Name); err != nil {
		PrintThenExit(err.Error())
	} else if exists {
		confirm := SimpleInputString("警告：项目"+aProject.Name+"已存在Podfile模板，是否覆盖?(y/N):", true)
		if confirm != "y" {
			PrintlnRed("未创建目标，程序退出！")
			return
		}
	}
	var aTemplate *TemplateModel
	for {
		templatePath := SimpleInputString("请输入模板文件(JSON格式)路径:", false)
		if !filepath.IsAbs(templatePath) {
			templatePath = filepath.Join(_Config.CurrentDir, templatePath)
		}
		if !fs.FileExists(templatePath) {
			PrintlnRed("模板文件[" + templatePath + "]不存在，请重新输入!")
			continue
		}
		fileData, err := ioutil.ReadFile(templatePath)
		if err != nil {
			PrintlnError("读取模板文件"+templatePath+"失败", err)
			PrintlnRed("请重新输入!")
			continue
		}
		if err = json.Unmarshal(fileData, &aTemplate); err != nil {
			PrintlnRed("模板文件" + templatePath + "不是有效的JSON格式，请重新输入!")
			continue
		}
		break
	}
	aTemplate.Team = aTeam.Name
	aTemplate.Project = aProject.Name
	if versionId, err := TemplateRequestCreate(aTemplate); err != nil {
		PrintThenExit(err.Error())
	} else {
		println("创建模板成功! 模板版本:", strconv.FormatInt(versionId, 10))
	}
}

func TemplatePodModifyVersion(parentMenu *CmdMenu) {
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}

	var pods []*PodModule
	for {
		pods = TemplatePodsInput()
		if !TemplateConfirmInputPods(pods) {
			PrintlnRed("您取消了之前的输入，请重新输入!")
			continue
		}
		break
	}

	url := GenURL("/api/logic/template/pods/modify")
	param := new(TemplateModel)
	param.Team = aTeam.Name
	param.Project = aProject.Name
	param.Pods = pods
	var res *TemplateResp
	if err := POSTParse(url, param, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.HasError() {
		PrintThenExit(res.Error("修改模板失败").Error())
	}
	print("修改成功")
	if res.VersionId > 0 {
		print("，模板已生成新版本，版本号:" + strconv.FormatInt(res.VersionId, 10)+"\n")
	} else {
		print("\n")
	}
	if res.Summary != "" {
		println("=== 详细信息 ===\n" + res.Summary)
	}
}

func TemplatePodAdd(parentMenu *CmdMenu)  {
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}

	aTemplate := new(TemplateModel)
	aTemplate.Team = aTeam.Name
	aTemplate.Project = aProject.Name
	aTemplate.Hierarchies = make([]*HierarchyModel, 0, 2)
	PrintlnYellow("Tips:为模板添加Pod需要分层级添加!")
	for {
		hName, err := TemplateHierarchySelect(aTeam.Name, aProject.Name)
		if err != nil {
			PrintThenExit(err.Error())
		}
		var aHierarchy *HierarchyModel
		for _, h := range aTemplate.Hierarchies {
			if h.Name == hName {
				aHierarchy = h
				break
			}
		}
		if aHierarchy != nil {
			PrintlnRed("您已经未层"+aHierarchy.Name+"添加了Pod，不可再添加，请选择其他层!")
			continue
		}

		aHierarchy = new(HierarchyModel)
		aHierarchy.Name = hName
		println("开始未层"+aHierarchy.Name+"添加Pod ...")
		aHierarchy.Pods = TemplatePodsInput()
		aTemplate.Hierarchies = append(aTemplate.Hierarchies, aHierarchy)

		confirm := SimpleInputString("继续喂其他层添加Pod?(y/N):", true)
		if confirm == "y" {
			continue
		}
		break
	}

	url := GenURL("/api/logic/template/pods/add")
	var res *TemplateResp
	if err := POSTParse(url, aTemplate, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.HasError() {
		PrintThenExit(res.Error("修改模板失败").Error())
	}

	print("添加Pod成功")
	if res.VersionId > 0 {
		print("，模板已生成新版本，版本号:" + strconv.FormatInt(res.VersionId, 10)+"\n")
	} else {
		print("\n")
	}
	if res.Summary != "" {
		println("=== 详细信息 ===\n" + res.Summary)
	}
}

func TemplatePodRemove(parentMenu *CmdMenu) {
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}

	var pods []*PodModule
	for {
		pods = TemplatePodsInput()
		if !TemplateConfirmInputPods(pods) {
			PrintlnRed("您取消了之前的输入，请重新输入!")
			continue
		}
		break
	}

	url := GenURL("/api/logic/template/pods/remove")
	param := new(TemplateModel)
	param.Team = aTeam.Name
	param.Project = aProject.Name
	param.Pods = pods
	var res *TemplateResp
	if err := POSTParse(url, param, &res); err != nil {
		PrintThenExit(err.Error())
	}
	if res.HasError() {
		PrintThenExit(res.Error("修改模板失败").Error())
	}
	print("删除Pod成功")
	if res.VersionId > 0 {
		print("，模板已生成新版本，版本号:" + strconv.FormatInt(res.VersionId, 10)+"\n")
	} else {
		print("\n")
	}
	if res.Summary != "" {
		println("=== 详细信息 ===\n" + res.Summary)
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

func TemplateHierarchySelect(team, project string) (string, error) {
	url := GenURL("/api/logic/template/hierarchy/names", "team", team, "project", project)
	var res *TemplateResp
	if err := GETParse(url, &res); err != nil {
		return "", err
	}
	if res.Errno != 0 {
		return "", errors.New(FormatResError("查询层级发送错误", res.Msg, res.Errno))
	}
	if len(res.HierarchyNames) == 0 {
		return "", errors.New("模板没有层级！")
	}
	var buffer bytes.Buffer
	buffer.WriteString("请选择层:\n")
	for idx, name := range res.HierarchyNames {
		buffer.WriteString(strconv.Itoa(idx+1) + "." + name + "\n")
	}
	buffer.WriteString(":")
	selected := SimpleInputSelectNum(buffer.String(), 1, len(res.HierarchyNames))
	return res.HierarchyNames[selected], nil
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

func TemplateExist(team, project string) (bool, error) {
	url := GenURL("/api/logic/template/exists", "team", team, "project", project)
	var res *TemplateResp
	if err := GETParse(url, &res); err != nil {
		return false, err
	}
	if res.Errno != 0 {
		return false, errors.New(FormatResError("查询模板是否存在失败", res.Msg, res.Errno))
	}
	return res.Exists, nil
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

func TemplateRequestCreate(aTemplate *TemplateModel) (int64, error) {
	url := GenURL("/api/logic/template/add")
	var res *TemplateResp
	if err := POSTParse(url, aTemplate, &res); err != nil {
		return 0, err
	}
	if res.Errno != 0 {
		return 0, errors.New(FormatResError("创建目模板失败", res.Msg, res.Errno))
	}
	return res.VersionId, nil
}

func TemplatePodsInput() []*PodModule {
	pods := make([]*PodModule, 0, 2)
	PrintlnYellow("Tips:Pod输入支持多个，每次(每行)输入一个Pod，包括Pod名称和版本号，版本号是可选，如果要输入版本号，Pod名称和版本号以空格分隔，直接回车（不输入）将结束输入流程!")
	println("请输入Pod名称和版本")
	dup := make(map[string]bool)
	for {
		line := SimpleInputString(":", true)
		if line == "" {
			if len(pods) == 0 {
				PrintlnRed("请至少输入一个Pod！")
				continue
			}
			break
		}
		items := strings.Split(line, " ")
		aPod := new(PodModule)
		aPod.Name = items[0]
		l := len(items)
		if l > 1 {
			for i := 1; i < l; i++ {
				v := items[i]
				if v != "" {
					aPod.Version = v
					break
				}
			}
		}

		if _, ok := dup[aPod.Name]; ok {
			confirm := SimpleInputString("您已经添加了名为"+aPod.Name+"的Pod，确定使用当前输入覆盖之前的输入吗?(y/N):", true)
			if confirm != "y" {
				continue
			}
			findIdx := -1
			for idx, aExistPod := range pods {
				if aPod.Name == aExistPod.Name {
					findIdx = idx
					break
				}
			}
			if findIdx > -1 {
				pods[findIdx] = aPod
				continue
			}
		}
		dup[aPod.Name] = true
		pods = append(pods, aPod)
	}
	return pods
}

func TemplateConfirmInputPods(pods []*PodModule) bool {
	var buffer bytes.Buffer
	buffer.WriteString("请确认您输入的Pods:\n")
	for _, aPod := range pods {
		buffer.WriteString("  - 名称:" + aPod.Name + " 版本:" + aPod.Version + "\n")
	}
	buffer.WriteString("确定使用以上Pods?(Y/n):")
	input := SimpleInputString(buffer.String(), true)
	return input == "Y" || input == "y" || input == ""
}
