package wr

import (
	"e.coding.net/vector-tech/golib/os/path/home"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"e.coding.net/vector-tech/golib/sec/license"
)

// Read 读取授权
func Read(salt, pubKey []byte) (ret []byte, err error) {
	appPath := path.Join(home.UserAppDataPath(), APP_CONFIG_PATH)
	appPath = filepath.ToSlash(appPath)
	fn := path.Join(appPath, LicenseName)
	f, err := os.Open(fn)
	if err != nil {
		return
	}
	defer f.Close()

	licenseData, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	ret, err = license.Decrypt(licenseData, salt, pubKey)
	return
}