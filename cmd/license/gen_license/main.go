package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/vicnoah/golib/pkg/sec/license"
)

var (
	help       bool
	privateKey string
	rawLicense string
	output     string
	fileName   string
	password   string
	salt       string
)

func init() {
	flag.BoolVar(&help, "help", false, "help")
	flag.StringVar(&password, "p", "", "`Authorization file password`")
	flag.StringVar(&salt, "s", "", "`Information used to authorize binding with the computer`")
	flag.StringVar(&rawLicense, "rl", "", "`Original License Information File Path`")
	flag.StringVar(&privateKey, "pk", "", "`Encryption private key`")
	flag.StringVar(&output, "o", "", "`license output directory`")
	flag.StringVar(&fileName, "fn", "license.vec", "`Key file name without path`")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if rawLicense == "" {
		fmt.Printf("\n%s\n", "原始授权文件不能为空")
		return
	}
	if password == "" {
		fmt.Printf("\n%s\n", "授权密码不能为空")
		return
	}
	if salt == "" {
		fmt.Printf("\n%s\n", "加密salt不能为空")
		return
	}
	if privateKey == "" {
		fmt.Printf("\n%s\n", "加密私钥不能为空")
		return
	}
	var outputName string
	if output != "" {
		outputName = path.Clean(output) + "/" + fileName
	} else {
		outputName = fileName
	}
	rFile, err := os.Open(rawLicense)
	if err != nil {
		fmt.Printf("\n打开原始授权文件失败: %v\n", err)
		return
	}
	defer rFile.Close()
	rLicense, err := io.ReadAll(rFile)
	if err != nil {
		fmt.Printf("\n读取原始授权文件失败: %v\n", err)
		return
	}
	pFile, err := os.Open(privateKey)
	if err != nil {
		fmt.Printf("\n打开私钥失败: %v\n", err)
		return
	}
	defer pFile.Close()
	pKey, err := io.ReadAll(pFile)
	if err != nil {
		fmt.Printf("\n读取密钥失败: %v\n", err)
		return
	}
	oFile, err := os.OpenFile(outputName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("\n创建密钥文件失败: %v\n", err)
		return
	}
	defer oFile.Close()
	ret, err := license.Encrypt(rLicense, []byte(password), []byte(salt), pKey)
	if err != nil {
		fmt.Printf("\n授权加密失败: %v\n", err)
		return
	}
	_, err = oFile.Write(ret)
	if err != nil {
		fmt.Printf("\n生成授权失败: %v\n", err)
		return
	}
	fmt.Printf("\n授权生成成功.授权文件路径: %s, 加密密钥: %s, 加密salt: %s,请保存相关信息.\n", outputName, password, salt)
}
