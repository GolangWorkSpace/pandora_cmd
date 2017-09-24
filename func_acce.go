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

func PodInstallAccLocal(parentMenu *CmdMenu) {
	currentPath := _Config.CurrentDir
	// 解析Podfile
	print("解析Podfile ...")
	podfilePath := filepath.Join(currentPath, "Podfile")
	if !fs.FileExists(podfilePath) {
		PrintThenExit("无法在当前目录找到Podfile文件!")
	}
	aPodfile, err := pod.NewPodfile(podfilePath)
	if err != nil {
		PrintThenExit("无法解析Podfile: ", err.Error())
	}
	print(" ok!\n")

	print("获取CocoaPods仓库列表 ...")
	repoPahts := funcGetRepoPaths()
	if len(repoPahts) == 0 {
		PrintThenExit("未找到Podspec仓库，程序退出！")
	}
	print(" ok!\n")

	// 选择仓库
	repoPath, repoName := funcGetSelectRepo(repoPahts)
	if err := funcUpdateRepo(repoName); err != nil {
		PrintThenExit("更新仓库失败！")
	}

	// 获取需要解析的Podspec路径
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

func funcGetCocoapodsRepoRoot() string {
	var repoRoot string
	user, err := user2.Current()
	if err == nil {
		repoRoot = filepath.Join(user.HomeDir, ".cocoapods", "repos")
	}
	if !fs.DirectoryExists(repoRoot) {
		foundation.IArgInput("无法找到CocoaPods的仓库根目录("+repoRoot+")，请手动指定(绝对路径):", func(arg string) foundation.IArgAction {
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
		})
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
	foundation.IArgInput(buffer.String(), func(arg string) foundation.IArgAction {
		idx, err := strconv.Atoi(arg)
		if err != nil || idx < 1 || idx > l {
			println("输入错误！")
			return foundation.IArgActionRepet
		}
		selectIdx = idx
		return foundation.IArgActionNext
	})
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
	print("等待Git操作完成 ...")
	time.Sleep(2 * time.Second)
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
	if len(parseSpecPaths) == 0 {
		return nil
	}
	threadCount := 10
	print("有需要解析的Podspec，")
	foundation.IArgInput("请输入解析线程数，建议不超过20，直接回车将使用默认值10:\n", func(arg string) foundation.IArgAction {
		if foundation.StrIsEmpty(arg) {
			return foundation.IArgActionNext
		}
		count, err := strconv.Atoi(arg)
		if err != nil {
			PrintlnRed("输入错误，请重新输入!")
			return foundation.IArgActionRepet
		}
		if count < 1 {
			PrintlnRed("线程数必须大于0，请重新输入！")
			return foundation.IArgActionRepet
		}
		threadCount = count
		return foundation.IArgActionNext
	})
	c := make(chan bool, threadCount)
	for _, p := range parseSpecPaths {
		c <- true
		go funcDoParseSingleSpec(p, c)
	}
	return nil
}

func funcDoParseSingleSpec(specpath string, c chan bool)  {
	println("解析：" + specpath + " ...")
	b, err := ioutil.ReadFile(specpath)
	if err != nil {
		PrintThenExit(err.Error())
	}

	podName, version, name := funcPodVersionName(specpath)
	desDir := filepath.Join(".pandora_cache", "spec", podName, version)
	if !fs.DirectoryExists(desDir) {
		if err := os.MkdirAll(desDir, os.ModePerm); err != nil {
			PrintThenExit("创建缓存目录失败："+desDir)
		}
	}

	ext := filepath.Ext(name)
	if ext == ".podspec" {
		if b, err = cmd.Exec("pod ipc spec " + specpath); err != nil {
			PrintThenExit("无法解析Podspec: " + specpath)
		}
	}

	if b, err = pod.SpecTrimDependency(b); err != nil {
		PrintThenExit("去除依赖失败："+specpath)
	}

	des := filepath.Join(desDir, podName+".podspec.json")
	if err = ioutil.WriteFile(des, b, os.ModePerm); err != nil {
		PrintThenExit("写入缓存spec失败：", des, err.Error())
	}

	<- c
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
	foundation.IArgInput(confirmStr, func(arg string) foundation.IArgAction {
		selected, err := strconv.Atoi(arg)
		if err != nil || selected < 1 || selected > 2 {
			println("输入非法，请重新输入！")
			return foundation.IArgActionRepet
		}
		replace = selected == 1
		return foundation.IArgActionNext
	})
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
