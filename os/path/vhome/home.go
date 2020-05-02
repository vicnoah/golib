// Package vhome Operations on the home folder of the operating system user
package vhome

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

const (
	// WIN runtime中windows系统名称
	WIN = "windows"
	// WIN_HOME_DRIVE windows HOMEDRIVE环境变量名称
	WIN_HOME_DRIVE = "HOMEDRIVE"
	// WIN_HOME_PATH windows HOME环境变量名称
	WIN_HOME_PATH = "HOMEPATH"
	// WIN_USERPROFILE windows USERPROFILE环境变量名称
	WIN_USERPROFILE = "USERPROFILE"
	// NIX_HOME UNIX 或 LINUX HOME环境变量名称
	NIX_HOME = "HOME"
	// WIN_APP_PATH windows APPDATA根路径
	WIN_APP_PATH = "AppData/Local"
	// NIX_APP_PATH UNIX 或 LINUX APPDATA根路径
	NIX_APP_PATH = ".config"
)

// UserAppDataPath Get user app file storage path
func UserAppDataPath() string {
	homeDir := UserHomeDir()
	homeDir = filepath.ToSlash(homeDir)
	if runtime.GOOS == WIN {
		return path.Join(homeDir, WIN_APP_PATH)
	}
	return path.Join(homeDir, NIX_APP_PATH)
}

// UserHomeDir Get user root directory
func UserHomeDir() string {
	if runtime.GOOS == WIN {
		home := os.Getenv(WIN_HOME_DRIVE) + os.Getenv(WIN_HOME_PATH)
		if home == "" {
			home = os.Getenv(WIN_USERPROFILE)
		}
		return home
	}
	return os.Getenv(NIX_HOME)
}