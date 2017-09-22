package main

import (
	"os/user"
	"path/filepath"
	"github.com/go-hayden-base/fs"
	"os"
)

var _Config *Config

func SetupConfig() {
	aConfig := new(Config)
	aConfig.SetupCache()
	aConfig.SetupPath()
	aConfig.SetupEnv()
	_Config = aConfig
}

type Config struct {
	Host string
	IsDev bool
	CacheRoot string
	Token string
	Account string
	CurrentDir string
	AuthDir string
	AuthTokenFile string
	AuthUserFile string
}

func (s *Config) SetupPath()  {
	currentDir, err := fs.CurrentDir()
	if err != nil {
		PrintThenExit(err.Error())
	}
	s.CurrentDir = currentDir
	s.AuthDir = filepath.Join(s.CacheRoot, ".a")
	s.AuthTokenFile = filepath.Join(s.AuthDir, "tk")
	s.AuthUserFile = filepath.Join(s.AuthDir, "user")
}

func (s *Config) SetupCache() {
	if s.CacheRoot == "" {
		aUser, err := user.Current()
		if err != nil {
			PrintThenExit(err.Error())
		}
		s.CacheRoot = filepath.Join(aUser.HomeDir, ".pandora", "cache")
	}
	if !fs.DirectoryExists(s.CacheRoot) {
		if err := os.MkdirAll(s.CacheRoot, os.ModePerm); err != nil {
			PrintThenExit(err.Error())
		}
	}
}

func (s *Config) SetupEnv()  {
	env := os.Getenv("PANDORA_CMD_CLIENT_ENV")
	if env == "dev" {
		s.Host = "http://127.0.0.1:3001"
		s.IsDev = true
	} else {
		s.Host = "http://172.24.41.12:3001"
		s.IsDev = false
	}
}
