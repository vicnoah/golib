package salt

import (
	"crypto/sha1"
	"encoding/hex"
	"os/exec"
	"strings"
)

// Get 获取操作系统信息生成salt，用于授权绑定
func Get() (sl []byte, err error) {
	cmd := exec.Command("CMD", "/C", "WMIC DISKDRIVE GET SERIALNUMBER")
	allInfo, err := cmd.Output()
	if err != nil {
		return
	}

	sn := strings.Split(string(allInfo), "\n")[1]
	h := sha1.New()
	_, err = h.Write([]byte(sn))
	if err != nil {
		return
	}
	sl = []byte(hex.EncodeToString(h.Sum(nil)))
	return
}