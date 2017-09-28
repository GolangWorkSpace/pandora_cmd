package main

import (
	"strconv"
	"io/ioutil"
	"github.com/go-hayden-base/fs"
	"path/filepath"
	"bytes"
	"strings"
	"os"
)

func PodfileGenerate(parentMenu *CmdMenu) {
	aProject, aTeam, err := ProjectSelect()
	if err != nil {
		PrintThenExit(err.Error())
	}

	var podfile *ProjectPodfileModel
	for {
		println("查询最新发布版的Podfile ...")
		url := GenURL("/api/logic/project/podfile/latest", "team", aTeam.Name, "project", aProject.Name)
		var res *ProjectPodfileResponse
		if err := GETParse(url, &res); err != nil {
			PrintThenExit(err.Error())
		}
		if res.Errno != 0 {
			PrintThenExit("发生错误, errno:", strconv.Itoa(res.Errno), " msg:", res.Msg)
		}

		tv := res.NeedUpgradVersion
		if tv > 0 {
			confirm := SimpleInputString("Podfile模板已变动，是否重新生成Podfile? (y/N):", true)
			if confirm != "y" {
				PrintThenExit("放弃操作，程序退出！")
			}

			url = GenURL("/api/logic/project/podfile/evolution", "team", aTeam.Name, "project", aProject.Name)
			res = nil
			if err := GETParse(url, &res); err != nil {
				PrintThenExit(err.Error())
			}
			if res.Errno != 0 {
				PrintThenExit("发生错误, errno:", strconv.Itoa(res.Errno), " msg:", res.Msg)
			}
			if res.Version < 1 {
				PrintThenExit("发生错误，无法解析生成的Podfile版本！")
			}

			confirm = SimpleInputString("Podfile草稿已生成，是否发布? (Y/n):", true)
			if confirm != "Y" && confirm != "" {
				PrintThenExit("放弃操作，程序退出！")
			}

			url = GenURLWithParam("/api/logic/project/podfile/release", map[string]string{
				"team":    aTeam.Name,
				"project": aProject.Name,
				"version": strconv.FormatInt(res.Version, 10),
			})

			res = nil
			if err := GETParse(url, &res); err != nil {
				PrintThenExit(err.Error())
			}
			if res.Errno != 0 {
				PrintThenExit("发生错误, errno:", strconv.Itoa(res.Errno), " msg:", res.Msg)
			}
			continue
		} else if res.Podfile != nil {
			podfile = res.Podfile
			break
		} else {
			PrintThenExit("发生错误，未能解析出Podfile！")
		}
	}

	podfPath := filepath.Join(_Config.CurrentDir, "Podfile")
	alert := ""
	if fs.FileExists(podfPath) {
		alert = "当前目录已存在Podfile是否覆盖?(Y/n):"
	} else {
		alert = "在当前目录生成Podfile?[Y/n]:"
	}
	confirm := SimpleInputString(alert, true)
	if confirm != "Y" && confirm != "" {
		PrintThenExit("放弃下载Podfile，程序退出!")
	}

	var buffer bytes.Buffer
	buffer.WriteString("# Pandora\n")
	buffer.WriteString("# 团队：" + podfile.Team + " 项目：" + podfile.Project + "\n")
	buffer.WriteString("# Podfile版本：" + strconv.FormatInt(podfile.Version, 10) + "\n")
	buffer.WriteString("# 基于模板版本：" + strconv.FormatInt(podfile.TemplateVersion, 10) + "\n")
	buffer.WriteString("# 创建人：" + podfile.CreateUser + " 时间：" + podfile.CreateTime + "\n")
	if len(podfile.Tags) > 0 {
		buffer.WriteString("# 标签：" + strings.Join(podfile.Tags, ", ") + "\n")
	}
	buffer.WriteString("\n")

	if podfile.Prefix != "" {
		buffer.WriteString(podfile.Prefix + "\n\n")
	}
	buffer.WriteString("target \"" + podfile.Target + "\" do\n")
	for _, h := range podfile.Hierarchies {
		buffer.WriteString("\n    # 层级：" + h.Name + "\n")
		for _, aPod := range h.Pods {
			line := prifGenPodLine(aPod)
			buffer.WriteString(line)
		}
		if len(h.ImplicitPods) > 0 {
			buffer.WriteString("\n    # 层级" + h.Name + "的隐性依赖\n")
			for _, aPod := range h.ImplicitPods {
				line := prifGenPodLine(aPod)
				buffer.WriteString(line)
			}
		}
	}
	buffer.WriteString("end\n\n")
	if podfile.Suffix != "" {
		buffer.WriteString(podfile.Suffix)
	}
	if err := ioutil.WriteFile(podfPath, buffer.Bytes(), os.ModePerm); err != nil {
		PrintThenExit("保存Podfile发生错误:", err.Error())
	}
	println("Podfile已保存:", podfPath)
}

func prifGenPodLine(aPod *PodModule) string {
	line := "    pod '" + aPod.Name + "'"
	if aPod.Version != "" {
		line += ", '" + aPod.Version + "'"
	}
	l := len(aPod.Subspecs)
	if l > 0 {
		line += ", :subspecs => ['"
		line += strings.Join(aPod.Subspecs, "', '")
		line += "']"
	}
	des := strings.TrimSpace(aPod.Description)
	if des != "" {
		des = " Pod描述：" + des
	}
	if aPod.Addition != nil {
		if aPod.Addition.ReferModuleViersion == aPod.Version {
			des += " 参照：" + aPod.Addition.ReferName + "(" + aPod.Addition.ReferVersion + ")"
		}
		if aPod.Addition.NewestVersion != "" && aPod.Addition.NewestVersion != aPod.Version {
			des += " 当前最新版本：" + aPod.Addition.NewestVersion
		}
	}

	if des != "" {
		l = len(line)
		if l >= 60 {
			line += " #" + des
		} else {
			line += (GenSpaceString(60-l) + "#" + des)
		}
	}
	line += "\n"
	return line
}
