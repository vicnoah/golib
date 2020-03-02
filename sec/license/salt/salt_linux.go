package salt

import (
	"crypto/sha1"
	"encoding/hex"
	"os/exec"
	"strings"
)

// Get Obtain operating system information to generate a salt for authorization binding
// The current implementation is not perfect, you need to rely on wmic
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