package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/vicnoah/golib/pkg/sec/vrsa"
)

const (
	_PRIVATE_KEY = "priv.pem"
	_PUBLIC_KEY  = "pub.pem"
)

var (
	help bool
	bits int
	dir  string
)

func init() {
	flag.BoolVar(&help, "help", false, "help")
	flag.IntVar(&bits, "b", 1024, "Key length, default is 1024 bits")
	flag.StringVar(&dir, "dir", "", "`Key export directory`")
}

func main() {

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	// 路径一致处理
	dir = path.Clean(dir)

	privName := ""
	pubName := ""
	if dir == "" {
		privName = _PRIVATE_KEY
		pubName = _PUBLIC_KEY
	} else {
		er := os.MkdirAll(dir, os.ModePerm)
		if er != nil {
			fmt.Printf("%v\r\n", er)
			return
		}
		privName = dir + "/" + _PRIVATE_KEY
		pubName = dir + "/" + _PUBLIC_KEY
	}

	privFile, err := os.OpenFile(privName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("%v\r\n", err)
		return
	}
	defer privFile.Close()
	pubFile, err := os.OpenFile(pubName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("%v\r\n", err)
		return
	}
	defer pubFile.Close()

	err = vrsa.GenRsaKey(bits, privFile, pubFile)
	if err != nil {
		fmt.Printf("%v\r\n", err)
		return
	}
	fmt.Printf("密钥生成成功->路径: %s, 私钥: %s, 公钥: %s\r\n", dir, _PRIVATE_KEY, _PUBLIC_KEY)
}
