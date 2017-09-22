package main

import (
	"strconv"
	"github.com/go-hayden-base/fs"
	"time"
	"path/filepath"
	"os"
	"io/ioutil"
	"strings"
	"github.com/go-hayden-base/pod"
)

func cmd_podfile_diff(args *Args) {
	team := args.GetSubargs("-t")[0]
	project := args.GetSubargs("-p")[0]
	oldv := args.GetSubargs("-v")[0]
	newv := args.GetSubargs("-v")[1]

	var res *ProjectProfileDiffResponse
	url := GenURL("/api/logic/project/podfile/diff", "team", team, "project", project, "old_version", oldv, "new_version", newv)
	if err := GETParse(url, &res); err != nil {
		PrintThenExit("对比失败：" + err.Error())
	}
	if res.Errno != 0 {
		PrintThenExit("对比失败，errno:", strconv.Itoa(res.Errno), "msg:", res.Msg)
	}
	if res.Diff == nil {
		println("无改动！")
	}
	change := false
	if len(res.Diff.Change) > 0 {
		change = true
		println("\n版本改动：")
		for name, vs := range res.Diff.Change {
			println("  -", name, ":", vs[0], "->", vs[1])
		}
	}

	if len(res.Diff.New) > 0 {
		change = true
		println("\n新增：")
		for name, v := range res.Diff.New {
			println("  -", name, ":", v)
		}
	}

	if len(res.Diff.Remove) > 0 {
		change = true
		println("\n移除：")
		for name, v := range res.Diff.Remove {
			println("  -", name, ":", v)
		}
	}

	if !change {
		println("无改动！")
	} else {
		println("\n")
	}
}

func cmd_diff_local(args *Args) {
	old := args.GetFirstSubArgs("-old")
	nw := args.GetFirstSubArgs("-new")
	if old == "" || !fs.FileExists(old) {
		PrintThenExit("old路径不存在！")
	}
	if nw == "" || !fs.FileExists(nw) {
		PrintThenExit("new路径不存在！")
	}

	timestamp := time.Now().Local().Format("20060102150405")
	oldDir := filepath.Join(_Config.CacheRoot, "diff", timestamp+"_old")
	if !fs.DirectoryExists(oldDir) {
		if err := os.MkdirAll(oldDir, os.ModePerm); err != nil {
			PrintThenExit(err.Error())
		}
	}
	newDir := filepath.Join(_Config.CacheRoot, "diff", timestamp+"_new")
	if !fs.DirectoryExists(newDir) {
		if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
			PrintThenExit(err.Error())
		}
	}

	oldPodfilePath := filepath.Join(oldDir, "Podfile")
	oldPodfileContent, err := ioutil.ReadFile(old)
	if err != nil {
		PrintThenExit(err.Error())
	}
	if err := ioutil.WriteFile(oldPodfilePath, oldPodfileContent, os.ModePerm); err != nil {
		PrintThenExit(err.Error())
	}

	newPodfilePath := filepath.Join(newDir, "Podfile")
	newPodfileContent, err := ioutil.ReadFile(nw)
	if err != nil {
		PrintThenExit(err.Error())
	}
	if err := ioutil.WriteFile(newPodfilePath, newPodfileContent, os.ModePerm); err != nil {
		PrintThenExit(err.Error())
	}

	aOldPodfile, err := pod.NewPodfile(oldPodfilePath)
	if err != nil {
		PrintThenExit(err.Error())
	}
	aNewPodfile, err := pod.NewPodfile(newPodfilePath)
	if err != nil {
		PrintThenExit(err.Error())
	}

	aDiff := new(tDiffModel)

	for _, aTarget := range aNewPodfile.Targets {
		for _, aModule := range aTarget.Modules {
			name := aModule.N
			if idx := strings.Index(name, "/"); idx > -1 {
				name = name[:idx]
			}
			v, ok := aOldPodfile.VersionOfModule(name, nil)
			if !ok {
				aDiff.AppendNew(name, aModule.V)
				continue
			}
			if v != aModule.V {
				aDiff.AppendChange(name, v, aModule.V)
			}
		}
	}

	for _, aOldTarget := range aOldPodfile.Targets {
		for _, aOldModule := range aOldTarget.Modules {
			name := aOldModule.N
			if idx := strings.Index(name, "/"); idx > -1 {
				name = name[:idx]
			}
			_, ok := aNewPodfile.VersionOfModule(name, nil)
			if !ok {
				aDiff.AppendRemove(name, aOldModule.V)
			}
		}
	}

	aDiff.Print()
}
