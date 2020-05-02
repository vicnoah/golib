package vbrowser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOpenURL(t *testing.T) {
	url := "http://localhost:8081"
	Convey("Test Browser Open URL", t, func() {
		So(func() {
			OpenURL(url)
		}, ShouldNotPanic)
		err := OpenURL(url)
		So(err, ShouldBeNil)
	})
}
