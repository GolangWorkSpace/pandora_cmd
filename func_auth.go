package main

import (
	"github.com/go-hayden-base/fs"
	"io/ioutil"
	"os"
)

func ShowMainMenu() {
	if !fs.FileExists(_Config.AuthTokenFile) || !fs.FileExists(_Config.AuthUserFile) {
		PrintlnYellow("\nTips:您尚未登录，无法使用远程服务!")
		_LocalMenu.Run()
	} else {
		ReadAuthInfo()
		PrintlnYellow("\n您当前的登录帐号为:" + _Config.Account)
		_Menu.Run()
	}
}

func ReadAuthInfo() {
	b, err := ioutil.ReadFile(_Config.AuthTokenFile)
	if err != nil {
		PrintThenExit(err.Error())
	}
	_Config.Token = string(b)

	b, err = ioutil.ReadFile(_Config.AuthUserFile)
	if err != nil {
		PrintThenExit(err.Error())
	}
	_Config.Account = string(b)
}

func WriteAuthInfo(token, account string) {
	if !fs.DirectoryExists(_Config.AuthDir) {
		if err := os.MkdirAll(_Config.AuthDir, os.ModePerm); err != nil {
			PrintThenExit(err.Error())
		}
	}
	if err := ioutil.WriteFile(_Config.AuthTokenFile, []byte(token), os.ModePerm); err != nil {
		PrintThenExit(err.Error())
	}
	if err := ioutil.WriteFile(_Config.AuthUserFile, []byte(account), os.ModePerm); err != nil {
		PrintThenExit(err.Error())
	}
}

func RemoveAuthInfo() {
	if fs.DirectoryExists(_Config.AuthDir) {
		if err := os.RemoveAll(_Config.AuthDir); err != nil {
			PrintThenExit(err.Error())
		}
	}
}
