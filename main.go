package main

var _Version = "2.0.0"

func init() {
	SetupConfig()
	SetupMenu()
}

func main() {
	ShowInfo()
	ShowMainMenu()
}

func ShowInfo()  {
	info := "\n当前环境:"
	if _Config.IsDev {
		info += "Debug "
	} else {
		info += "Release "
	}
	info += "版本:" + _Version
	PrintlnYellow(info)
}

