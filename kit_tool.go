package main

import (
	"os"
	cp "github.com/fatih/color"
	"strings"
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

func PrintlnYellow(msg ...string)  {
	cp.Yellow(strings.Join(msg, " "))
}

func PrintlnRed(msg ...string)  {
	cp.Red(strings.Join(msg, " "))
}

func PrintThenExit(msg ...string)  {
	PrintlnRed(msg...)
	os.Exit(1)
}
