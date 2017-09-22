package main

import (
	"os"
	cp "github.com/fatih/color"
	"strings"
	"strconv"
	"github.com/go-hayden-base/foundation"
)

func VaildStringParams(params ...string) (bool) {
	for _, param := range params {
		if param == "" {
			return false
		}
	}
	return true
}

func GenSpaceString(count int) (string) {
	s := ""
	for i := 0; i < count; i++ {
		s += " "
	}
	return s
}

func PrintlnBlue(msg ...string) {
	cp.Blue(strings.Join(msg, " "))
}

func PrintlnYellow(msg ...string) {
	cp.Yellow(strings.Join(msg, " "))
}

func PrintlnRed(msg ...string) {
	cp.Red(strings.Join(msg, " "))
}

func PrintThenExit(msg ...string) {
	PrintlnRed(msg...)
	os.Exit(1)
}

func PrintlnErrorFormat(prefix, msg string, errno int) {
	PrintlnRed(FormatResError(prefix, msg, errno))
}

func FormatResError(prefix, msg string, errno int) string {
	return prefix + "," + " 错误码:" + strconv.Itoa(errno) + " 信息:" + msg
}

func PrintlnError(prefix string, err error) {
	PrintlnRed(FormtError(prefix, err))
}

func FormtError(prefix string, err error) string {
	return prefix + ": " + err.Error()
}

func SimpleInputString(alert string, canBeEmpty bool) string {
	input := ""
	foundation.IArgInput(alert, func(arg string) foundation.IArgAction {
		if !canBeEmpty && arg == "" {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		input = arg
		return foundation.IArgActionNext
	})
	return input
}

func SimpleInputSelectNum(alert string, start, length int) int {
	selected := 0
	foundation.IArgInput(alert, func(arg string) foundation.IArgAction {
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < start || idx > length + start - 1 {
			PrintlnRed("输入错误, 请重新输入!")
			return foundation.IArgActionRepet
		}
		selected = idx - start
		return foundation.IArgActionNext
	})
	return selected
}

func SimpleInputInt(alert string, start, length, defaultVal int) int {
	val := defaultVal
	foundation.IArgInput(alert, func(arg string) foundation.IArgAction {
		if arg == "" {
			return foundation.IArgActionNext
		}
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < start || idx > length + start - 1 {
			PrintlnRed("输入错误, 请重新输入!")
			return foundation.IArgActionRepet
		}
		val = idx
		return foundation.IArgActionNext
	})
	return val
}
