package main

import (
	"os"
	"time"
	"strconv"
)


var _Args *Args

var _Host string = "http://127.0.0.1:3001"

func init() {
	_Args = NewArgs()
	if _Args == nil {
		print("无法解析参数！")
		os.Exit(0)
	}
	_Args.RegisterFunc("podfile", cmd_podfile)
	_Args.RegisterFunc("podfile-diff", cmd_podfile_diff)
	_Args.RegisterFunc("follow", cmd_follow)
	_Args.RegisterFunc("acc", cmd_acc)
}

func main() {
	start := time.Now().UnixNano()
	ok := _Args.Exec()
	if !ok {
		print("未执行任何操作！")
	}
	end := time.Now().UnixNano()
	cost := float64(end-start) / float64(1000000000)
	println("程序耗时:  " + strconv.FormatFloat(cost, 'f', -1, 64) + " 秒")
}
