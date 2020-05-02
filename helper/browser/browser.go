package vbrowser

import (
	"git.sabertrain.com/vector-tech/golib/util/cmd"
	"fmt"
	"runtime"
)

var (
	commands = map[string][]string{
		"windows": {"start", "explorer", "rundll32"},
		"darwin":  {"open"},
		"linux":   {"xdg-open"},
	}
	args = map[string][]string{
		"start":    {},
		"explorer": {},
		"rundll32": {"url.dll,FileProtocolHandler"},
		"open":     {},
		"xdg-open": {},
	}
)

// OpenURL calls the OS default program for uri
func OpenURL(uri string) (err error) {
	runs, ok := commands[runtime.GOOS]
	if !ok {
		err = fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
		return
	}

	for _, v := range runs {
		args := append(args[v], uri)
		_, _, err = vcmd.Command(v, args...)
		if err == nil {
			return
		}
	}
	return
}
