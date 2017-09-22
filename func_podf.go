package main

import (
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"
	"github.com/go-hayden-base/foundation"
	"errors"
	"github.com/go-hayden-base/fs"
	"path/filepath"
	"bytes"
	"strings"
	"os"
)

func cmd_podfile(args *Args) {
	team := args.GetSubargs("-t")[0]
	project := args.GetSubargs("-p")[0]

	if !VaildStringParams(team, project) {
		PrintThenExit("参数错误！")
	}

	var podfile *ProjectPodfileModel
	for {
		println("查询最新发布版的Podfile ...")
		url := GenURL("/api/logic/project/podfile/latest", "team", team, "project", project)
		var res *ProjectPodfileResponse
		if err := GETParse(url, &res); err != nil {
			PrintThenExit(err.Error())
		}
		if res.Errno != 0 {
			PrintThenExit("发生错误, errno:", strconv.Itoa(res.Errno), " msg:", res.Msg)
		}

		tv := res.NeedUpgradVersion
		if tv > 0 {
			gen := false
			foundation.IArgInput("Podfile模板已变动，是否重新生成Podfile? (y/n):", func(arg string) foundation.IArgAction {
				if arg == "y" {
					gen = true
				}
				return foundation.IArgActionNext
			})

			if !gen {
				PrintThenExit("放弃操作，程序退出！")
			}

			url = GenURL("/api/logic/project/podfile/evolution", "team", team, "project", project)
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
			release := false
			foundation.IArgInput("Podfile草稿已生成，是否发布? (y/n):", func(arg string) foundation.IArgAction {
				if arg == "y" {
					release = true
				}
				return foundation.IArgActionNext
			})

			if !gen {
				PrintThenExit("放弃操作，程序退出！")
			}

			url = GenURL("/api/logic/project/podfile/release", "team", team, "project", project, "version", strconv.FormatInt(res.Version, 10))
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

	currentPath, err := fs.CurrentDir()
	if err != nil {
		PrintThenExit("获取当前目录失败:", err.Error())
	}
	currentPath = filepath.Join(currentPath, "Podfile")
	download := false
	foundation.IArgInput("是否在当前目录("+currentPath+")创建新Podfile? (y/n):", func(arg string) foundation.IArgAction {
		if arg == "y" {
			download = true
		}
		return foundation.IArgActionNext
	})

	if !download {
		PrintThenExit("放弃操作，程序退出！")
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
					line += (GenSpaceString(60 - l)+"#"+des)
				}
			}
			line += "\n"
			buffer.WriteString(line)
		}
	}
	buffer.WriteString("end\n\n")
	if podfile.Suffix != "" {
		buffer.WriteString(podfile.Suffix)
	}
	if err := ioutil.WriteFile(currentPath, buffer.Bytes(), os.ModePerm); err != nil {
		PrintThenExit("保存Podfile发生错误:", err.Error())
	}
	println("Podfile已保存:", currentPath)
}

func intify(info map[string]interface{}, key string) (int, error) {
	x, ok := info[key]
	if !ok {
		return 0, errors.New("无法获取" + key + "的值！")
	}
	r, ok := x.(float64)
	if !ok {
		return 0, errors.New("无法转换" + key + "的值！")
	}
	return int(r), nil
}

func errorInfo(info map[string]interface{}) (int, string, error) {
	en, ok := info["errno"]
	if !ok {
		return 0, "", errors.New("无法获取errno!")
	}
	errno, ok := en.(float64)
	if !ok {
		return 0, "", errors.New("无法转换errno!")
	}

	m, ok := info["msg"]
	message := ""
	if ok {
		message = m.(string)
	}
	return int(errno), message, nil
}

func GET(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return res, nil
}
