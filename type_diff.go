package main

import "strings"

type tDiffModel struct {
	Change map[string][]string `json:"change,omitempty"`
	New    map[string]string `json:"new,omitempty"`
	Remove map[string]string `json:"remove,omitempty"`
}

func (s *tDiffModel) AppendChange(name, oldv, newv string) {
	if s.Change == nil {
		s.Change = make(map[string][]string)
	}
	c := make([]string, 2, 2)
	c[0] = oldv
	c[1] = newv
	s.Change[name] = c
}

func (s *tDiffModel) AppendNew(name, version string) {
	if s.New == nil {
		s.New = make(map[string]string)
	}
	s.New[name] = version
}

func (s *tDiffModel) AppendRemove(name, version string) {
	if s.Remove == nil {
		s.Remove = make(map[string]string)
	}
	s.Remove[name] = version
}

func (s *tDiffModel) Print() {
	lc := len(s.Change)
	lr := len(s.Remove)
	ln := len(s.New)
	if lc == 0 && lr == 0 && ln == 0 {
		println("没有变动！")
		return
	}
	if lc > 0 {
		println("\n版本变动：")
		for name, change := range s.Change {
			println("  -", name, ": ", strings.Join(change, " -> "))
		}
	}
	if ln > 0 {
		println("\n新增：")
		for name, version := range s.New {
			println("  -", name, ": ", version)
		}
	}
	if lr > 0 {
		println("\n移除：")
		for name, version := range s.Remove {
			println("  -", name, ": ", version)
		}
	}
}
