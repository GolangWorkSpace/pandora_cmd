package main

import (
	"bytes"
	"strconv"
)

type CmdMenu struct {
	Title      string
	Func       func(parentMenu *CmdMenu)
	SubMenu    []*CmdMenu
	ParentMenu *CmdMenu
}

func (s *CmdMenu) Run() {
	l := len(s.SubMenu)
	if l == 0 {
		if s.Func != nil {
			s.Func(s)
		} else {
			PrintThenExit("该功能尚未实现，敬请期待！")
		}
		return
	}

	var buffer bytes.Buffer
	buffer.WriteString("\n")
	menuPath := s.MenuPath()
	if menuPath != "" {
		buffer.WriteString("=== 菜单路径：" + menuPath + " ===\n")
	}
	buffer.WriteString("请选择:\n")
	for idx, aSubMenu := range s.SubMenu {
		buffer.WriteString(strconv.Itoa(idx+1) + "." + aSubMenu.Title + "\n")
	}

	if s.ParentMenu != nil {
		l += 1
		buffer.WriteString(strconv.Itoa(l) + ".<<返回上级菜单\n")
	}

	buffer.WriteString(":")
	selected := SimpleInputSelectNum(buffer.String(), 1, l)

	if s.ParentMenu != nil && selected == l-1 {
		s.ParentMenu.Run()
		return
	}
	aSubMenu := s.SubMenu[selected]
	aSubMenu.Run()
}

func (s *CmdMenu) MenuPath() string {
	if s.ParentMenu != nil {
		parentPath := s.ParentMenu.MenuPath()
		if parentPath != "" {
			return parentPath + "/" + s.Title
		} else {
			return s.Title
		}
	}
	return s.Title
}

func (s *CmdMenu) PrepareParentMenu(parentMenu *CmdMenu) {
	s.ParentMenu = parentMenu
	if len(s.SubMenu) == 0 {
		return
	}
	for _, aSubMenu := range s.SubMenu {
		aSubMenu.PrepareParentMenu(s)
	}
}
