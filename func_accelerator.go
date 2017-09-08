package main

import (
	"github.com/go-hayden-base/fs"
	"path/filepath"
	"github.com/go-hayden-base/pod"
	"github.com/go-hayden-base/foundation"
	"bytes"
	"strconv"
	"github.com/go-hayden-base/cmd"
	"os"
	"errors"
	"io/ioutil"
	"regexp"
	"time"
	user2 "os/user"
	"strings"
)

func cmd_acc(args *Args) {
	currentPath, err := fs.CurrentDir()
	if err != nil {
		PrintThenExit("无法获取当前路径:", err.Error())
	}
	println("解析Podfile ...")
	podfilePath := filepath.Join(currentPath, "Podfile")
	if !fs.FileExists(podfilePath) {
		PrintThenExit("无法在当前目录找到Podfile文件!")
	}

	aPodfile, err := pod.NewPodfile(podfilePath)
	if err != nil {
		PrintThenExit("无法解析Podfile:", err.Error())
	}
	funcCheckPodfile(aPodfile)

	repoPahts := funcGetRepoPaths()
	if len(repoPahts) == 0 {
		PrintThenExit("未找到Podspec仓库，程序退出！")
	}

	repoPath, repoName := funcGetSelectRepo(repoPahts)
	if err := funcUpdateRepo(repoName); err != nil {
		PrintThenExit("更新仓库失败！")
	}

	needsParseSpecPaths, err := funcGetNeedsParseSpecPaths(repoPath, aPodfile)
	if err != nil {
		PrintThenExit("获取需要解析的Podspec失败:", err.Error())
	}

	if err = funcParseSpecAndGen(needsParseSpecPaths); err != nil {
		PrintThenExit("解析Podspec失败: " + err.Error())
	}

	buffer, err := funcPodfileReplaceThenGenBuffer(podfilePath)
	if err != nil {
		PrintThenExit("解析替换Podfile失败: " + err.Error())
	}

	msg, err := funcPodfileGen(buffer)
	if err != nil {
		PrintThenExit("生成Podfile失败: " + err.Error())
	}
	println(msg)
}

func funcCheckPodfile(aPodfile *pod.Podfile) {
	for _, aTarget := range aPodfile.Targets {
		for _, aPod := range aTarget.Modules {
			if aPod.V == "" && aPod.SpecPath == ""{
				PrintlnYellow("警告:", "未指定版本 ->", aPod.N)
			} else if aPod.V == "" && aPod.SpecPath != "" {
				PrintlnBlue("提示:", "指定路径 ", )
			}
		}
	}
}

func funcGetCocoapodsRepoRoot() string {
	var repoRoot string
	user, err := user2.Current()
	if err == nil {
		repoRoot = filepath.Join(user.HomeDir, ".cocoapods", "repos")
	}
	if !fs.DirectoryExists(repoRoot) {
		new(foundation.IArag).Input("无法找到CocoaPods的仓库根目录("+repoRoot+")，请手动指定(绝对路径):", func(arg string) foundation.IArgAction {
			if arg != "" && !filepath.IsAbs(arg) {
				print("路径错误，请重新指定!")
				return foundation.IArgActionRepet
			}
			if !fs.DirectoryExists(arg) {
				print("路径不存在，请重新指定！")
				return foundation.IArgActionRepet
			}
			repoRoot = arg
			return foundation.IArgActionNext
		}).Run()
	}
	return repoRoot
}

func funcGetRepoPaths() []string {
	repoRoot := funcGetCocoapodsRepoRoot()
	repoPaths := make([]string, 0, 5)
	fs.ListDirectory(repoRoot, false, func(file fs.FileInfo, err error) {
		if err != nil || !file.IsDir() || file.Name() == "master" {
			return
		}
		gitPath := filepath.Join(file.FilePath(), ".git")
		if !fs.DirectoryExists(gitPath) {
			return
		}
		repoPaths = append(repoPaths, file.FilePath())
	})
	return repoPaths
}

func funcGetSelectRepo(repoPaths []string) (string, string) {
	var buffer bytes.Buffer
	buffer.WriteString("请选择要使用的仓库：\n")
	for idx, p := range repoPaths {
		buffer.WriteString(strconv.Itoa(idx+1) + ".")
		buffer.WriteString(filepath.Base(p) + "\n")
	}
	l := len(repoPaths)
	selectIdx := 0
	new(foundation.IArag).Input(buffer.String(), func(arg string) foundation.IArgAction {
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < 1 || idx > l {
			println("输入错误！")
			return foundation.IArgActionRepet
		}
		selectIdx = idx
		return foundation.IArgActionNext
	}).Run()
	selectedPath := repoPaths[selectIdx-1]
	return selectedPath, filepath.Base(selectedPath)
}

func funcUpdateRepo(name string) error {
	println("更新仓库:", name, "...")
	cmdStr := "pod repo update " + name
	aCmd := cmd.Cmd(cmdStr)
	aCmd.Stdout = os.Stdout
	err := aCmd.Run()
	if err != nil {
		return err
	}
	print("等待Git完成 ...")
	time.Sleep(3 * time.Second)
	print(" ok!\n")
	return nil
}

func funcGetNeedsParseSpecPaths(repoPath string, aPodfile *pod.Podfile) ([]string, error) {
	parseSpecs := make([]string, 0, 10)
	for _, aTarget := range aPodfile.Targets {
		for _, aPod := range aTarget.Modules {
			if aPod.Type != "" {
				continue
			}
			name := aPod.N
			if idx := strings.Index(name, "/"); idx > -1 {
				name = name[:idx]
			}

			desDir := filepath.Join(".pandora_cache", "spec", name, aPod.V)
			des := filepath.Join(desDir, name+".podspec.json")
			if fs.FileExists(des) {
				continue
			}

			srcDir := filepath.Join(repoPath, name, aPod.V)
			if !fs.DirectoryExists(srcDir) {
				return nil, errors.New("无法找到Podspec目录: " + name + " " + aPod.V + " " + srcDir)
			}

			src := funcGetSpecFilePath(srcDir)
			if src == "" {
				return nil, errors.New("无法找到Podspec文件: " + name + " " + aPod.V)
			}

			if foundation.SliceContainsStr(src, parseSpecs) {
				continue
			}
			parseSpecs = append(parseSpecs, src)
		}
	}
	return parseSpecs, nil
}

func funcGetSpecFilePath(root string) (string) {
	p := ""
	fs.ListDirectory(root, false, func(file fs.FileInfo, err error) {
		if err != nil || file.IsDir() || p != "" {
			return
		}
		ext := filepath.Ext(file.Name())
		if ext != ".json" && ext != ".podspec" {
			return
		}
		p = file.FilePath()
	})
	return p
}

func funcParseSpecAndGen(parseSpecPaths []string) error {
	for _, p := range parseSpecPaths {
		println("解析：" + p + " ...")
		b, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}

		podName, version, name := funcPodVersionName(p)
		desDir := filepath.Join(".pandora_cache", "spec", podName, version)
		if !fs.DirectoryExists(desDir) {
			if err := os.MkdirAll(desDir, os.ModePerm); err != nil {
				return err
			}
		}

		ext := filepath.Ext(name)
		if ext == ".podspec" {
			if b, err = cmd.Exec("pod ipc spec " + p); err != nil {
				return errors.New("无法解析Podspec: " + p)
			}
		}
		if b, err = pod.SpecTrimDependency(b); err != nil {
			return err
		}

		des := filepath.Join(desDir, podName+".podspec.json")
		if err = ioutil.WriteFile(des, b, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func funcPodVersionName(p string) (string, string, string) {
	name := filepath.Base(p)
	p = filepath.Dir(p)
	version := filepath.Base(p)
	p = filepath.Dir(p)
	podName := filepath.Base(p)
	return podName, version, name
}

var regPod = regexp.MustCompile(`^(\s*pod\s+)'(\S+)'\s*,\s*'(\S+)'(.*)`)

func funcPodfileReplaceThenGenBuffer(podfilePath string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	var err error
	fs.ReadLine(podfilePath, func(line string, finished bool, err error, stop *bool) {
		if err != nil {
			return
		}
		if !regPod.MatchString(line) {
			buffer.WriteString(line)
			if !finished {
				buffer.WriteString("\n")
			}
			return
		}
		matchs := regPod.FindAllStringSubmatch(line, -1)
		if len(matchs) == 0 {
			err = errors.New("无法解析 -> " + line)
			*stop = true
			return
		}
		items := matchs[0]
		if len(items) < 5 {
			err = errors.New("无法解析子匹配 -> " + line)
			*stop = true
			return
		}
		prefix, podName, podVersion, suffix := items[1], items[2], items[3], items[4]
		specPath := filepath.Join(".pandora_cache", "spec", podName, podVersion, podName+".podspec.json")
		if !fs.FileExists(specPath) {
			err = errors.New("找不到Podspec -> " + specPath)
			*stop = true
			return
		}
		newPath := prefix + "'" + podName + "', :podspec => '" + specPath + "'" + suffix
		buffer.WriteString(newPath + "\n")
	})

	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func funcPodfileGen(buffer *bytes.Buffer) (string, error) {
	replace := false
	confirmStr := "加速后的Podfile已经生成，请选择操作:\n" +
		"1.替换（替换会将原Podfile重命名为Podfile加时间戳后缀的格式，如Podfile_20170903151211）\n" +
		"2.不替换（选择不替换将会在当前目录生成Podfile_acc的文件）"
	new(foundation.IArag).Input(confirmStr, func(arg string) foundation.IArgAction {
		selected, err := strconv.Atoi(arg)
		if err != nil || selected < 1 || selected > 2 {
			println("输入非法，请重新输入！")
			return foundation.IArgActionRepet
		}
		replace = selected == 1
		return foundation.IArgActionNext
	}).Run()
	if !replace {
		if err := ioutil.WriteFile("Podfile_acc", buffer.Bytes(), os.ModePerm); err != nil {
			return "", err
		}
		return "成功生成加速后的Podfile_acc!", nil
	}

	newName := "Podfile_" + time.Now().Local().Format("20060102150405")
	if err := os.Rename("Podfile", "Podfile_"+newName); err != nil {
		return "", errors.New("重命名原Podfile失败 -> " + err.Error())
	}
	println("重命名原Podfile为" + newName)
	if err := ioutil.WriteFile("Podfile", buffer.Bytes(), os.ModePerm); err != nil {
		return "", errors.New("生成新Podfile失败 -> " + err.Error())
	}
	return "生成加速后的新Podfile成功！", nil
}
