package main

import (
	"strconv"
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
