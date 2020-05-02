package wr

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"git.sabertrain.com/vector-dev/golib/os/path/vhome"
)

// Write 写入授权
func Write(input io.Reader) (err error) {
	appPath := path.Join(vhome.UserAppDataPath(), APP_CONFIG_PATH)
	appPath = filepath.ToSlash(appPath)
	err = os.MkdirAll(appPath, os.ModePerm)
	if err != nil {
		return
	}
	licenseData, err := ioutil.ReadAll(input)
	if err != nil {
		return
	}
	fn := path.Join(appPath, LicenseName)
	licenseFile, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer licenseFile.Close()
	_, err = licenseFile.Write(licenseData)
	return
}
