package main

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/go-hayden-base/foundation"
	"strconv"
)

func cmd_follow(args *Args) {
	team := args.GetSubargs("-t")[0]
	project := args.GetSubargs("-p")[0]

	url := _Host + "/api/ref/version/list"
	data, err := GET_DATA(url)
	if err != nil {
		println(err.Error())
		return
	}

	var aReferListRes *FollowReferRes
	err = json.Unmarshal(data, &aReferListRes)
	if err != nil {
		println(err.Error())
		return
	}

	if len(aReferListRes.Refs) == 0 {
		println("无可选参照！")
		return
	}

	alert := "请选择要跟进的版本：\n"
	for idx, aRefer := range aReferListRes.Refs {
		alert += ( strconv.Itoa(idx+1) + "." + aRefer["refer_name"] + "  版本:" + aRefer["refer_version"] + "\n" )
	}

	selected := 0
	new(foundation.IArag).Input(alert, func(arg string) foundation.IArgAction {
		i, err := strconv.Atoi(arg)
		if err != nil {
			println("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		} else if i < 1 || i > len(aReferListRes.Refs) {
			println("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		selected = i - 1
		return foundation.IArgActionNext
	}).Run()

	aRefer := aReferListRes.Refs[selected]
	url = _Host + "/api/logic/template/follow"
	url += "?team=" + team
	url += "&project=" + project
	url += "&refer_name=" + aRefer["refer_name"]
	url += "&refer_version=" + aRefer["refer_version"]

	aRes, err := GET(url)
	errno, msg, err := errorInfo(aRes)
	if err != nil {
		println(err.Error())
		return
	} else if errno != 0 {
		println("发生错误:" + msg)
		return
	}
	println("跟进成功！模板已跟进到："+aRefer["refer_name"]+" "+ aRefer["refer_version"])
}

type FollowReferRes struct {
	Errno float64 `json:"errno,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Refs  []map[string]string `json:"refs,omitempty"`
}

func GET_DATA(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
