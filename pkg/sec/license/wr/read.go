package wr

import (
	"git.sabertrain.com/vector-dev/golib/os/path/vhome"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"git.sabertrain.com/vector-dev/golib/pkg/sec/license"
)

// Read 读取授权
func Read(salt, pubKey []byte) (ret []byte, err error) {
	appPath := path.Join(vhome.UserAppDataPath(), APP_CONFIG_PATH)
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