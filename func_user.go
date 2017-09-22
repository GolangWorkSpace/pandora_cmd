package main

import (
	"strconv"
	"regexp"
	"strings"
	"github.com/go-hayden-base/foundation"
	"bytes"
)

func UserRegister(parentMenu *CmdMenu) {
	account := ""
	foundation.IArgInput("请输入注册的帐号（必须是dianping.com或meituan.com的邮箱）:", func(arg string) foundation.IArgAction {
		if !regexp.MustCompile(`^\S+@((dianping.com)|(meituan.com))$`).MatchString(arg) {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		account = arg
		return foundation.IArgActionNext
	})

	password := ""
	foundation.IArgInput("请输入登录密码(6-20位字母或数字的组合):", func(arg string) foundation.IArgAction {
		if !regexp.MustCompile(`^[0-9a-zA-Z]{6,20}$`).MatchString(arg) {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		password = arg
		return foundation.IArgActionNext
	})

	name := ""
	foundation.IArgInput("请输入您的名称（默认为帐号）:", func(arg string) foundation.IArgAction {
		if arg == "" {
			if idx := strings.Index(account, "@"); idx > 0 {
				name = account[:idx]
			} else {
				name = account
			}
		} else {
			name = arg
		}
		return foundation.IArgActionNext
	})

	param := map[string]string{
		"account":  account,
		"password": foundation.StrMD5(password),
		"name":     name,
	}

	var userResp *UserResp
	url := GenURL("/api/user/register")
	if err := POSTParse(url, param, &userResp); err != nil {
		PrintThenExit("注册失败:", err.Error())
	}

	if userResp.Errno != 0 {
		PrintThenExit("注册失败，错误码:", strconv.Itoa(userResp.Errno), "信息:", userResp.Msg)
	}

	PrintlnBlue("恭喜！注册成功，请登录！")
	UserLogin(parentMenu)
}

func UserLogin(parentMenu *CmdMenu) {
	for {
		account := ""
		foundation.IArgInput("请输入登录的帐号：", func(arg string) foundation.IArgAction {
			account = arg
			return foundation.IArgActionNext
		})
		password := ""
		foundation.IArgInput("请输入登录的密码：", func(arg string) foundation.IArgAction {
			password = arg
			return foundation.IArgActionNext
		})
		password = foundation.StrMD5(password)
		param := map[string]string{
			"account":  account,
			"password": password,
		}
		url := GenURL("/api/user/login")
		var userRes *UserResp
		if err := POSTParse(url, param, &userRes); err != nil {
			PrintThenExit(err.Error())
		}
		if userRes.Errno != 0 {
			PrintlnRed("登录失败，错误码:", strconv.Itoa(userRes.Errno), "信息:", userRes.Msg)
			continue
		}
		WriteAuthInfo(userRes.User.CMDToken, userRes.User.Account)
		println("登录成功！")
		ShowMainMenu()
		break
	}
}

func UserSelectRole() Role {
	var buffer bytes.Buffer
	buffer.WriteString("请选择角色:\n")
	for idx, aRoleInfo := range SelectRoles {
		buffer.WriteString(strconv.Itoa(idx+1) + "." + aRoleInfo.Name + "\n")
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
	return SelectRoles[selected].Role
}

func UserChangePassword(parentMenu *CmdMenu) {
	for {
		oldPassword := ""
		foundation.IArgInput("请输入旧密码:", func(arg string) foundation.IArgAction {
			oldPassword = foundation.StrMD5(arg)
			return foundation.IArgActionNext
		})
		newPassword := ""
		foundation.IArgInput("请输入新密码(6-20位字母或数字的组合):", func(arg string) foundation.IArgAction {
			if !regexp.MustCompile(`^[0-9a-zA-Z]{6,20}$`).MatchString(arg) {
				PrintlnRed("输入新密码错误，请重新输入!")
				return foundation.IArgActionRepet
			}
			newPassword = foundation.StrMD5(arg)
			return foundation.IArgActionNext
		})
		url := GenURL("/api/user/changepassword", "old_password", oldPassword, "new_password", newPassword)
		var res *UserResp
		if err := GETParse(url, &res); err != nil {
			PrintThenExit(err.Error())
		}
		if res.Errno != 0 {
			PrintlnRed("修改密码失败，错误码:", strconv.Itoa(res.Errno), "信息:", res.Msg)
			PrintlnRed("请重试!")
			continue
		}
		println("修改密码成功!")
		break
	}
}

func UserLogout(parentMenu *CmdMenu) {
	RemoveAuthInfo()
	println("已注销！")
}
