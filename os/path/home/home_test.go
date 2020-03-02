package home

import (
	"github.com/smartystreets/assertions/should"
	. "github.com/smartystreets/goconvey/convey"
	"runtime"
	"testing"
)

func TestUserHomeDir(t *testing.T) {
	Convey("Test User Home Dir", t, func() {
		So(func() {
			UserHomeDir()
		}, ShouldNotPanic)
		homeDir := UserHomeDir()
		if runtime.GOOS == "windows" {
			So(homeDir, should.ContainSubstring, "user")
		} else if runtime.GOOS == "linux" {
			So(homeDir, ShouldContainSubstring, "/home/")
		}
	})
}

func TestUserAppDataPath(t *testing.T) {
	Convey("Test User AppData Path", t, func() {
		So(func() {
			UserAppDataPath()
		}, ShouldNotPanic)
		userAppDataPath := UserAppDataPath()
		if runtime.GOOS == "windows" {
			So(userAppDataPath, ShouldContainSubstring, "AppData")
		} else if runtime.GOOS == "linux" {
			println(userAppDataPath)
			So(userAppDataPath, ShouldContainSubstring, ".config")
		}
	})
}