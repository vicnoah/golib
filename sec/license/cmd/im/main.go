package main

import (
	"e.coding.net/vector-tech/golib/sec/license/salt"
	"e.coding.net/vector-tech/golib/sec/license/wr"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	_ACTION_COMPUTER_CODE = "code"
	_ACTION_IMPORT        = "import"
	_ACTION_Read          = "read"
	_CODE_OUTPUT_NAME     = "computer.ivec"
)

var (
	help       bool
	action     string
	codeOutput string
	inputFile  string
	pubKeyFile string
)

func init() {
	flag.BoolVar(&help, "help", false, "")
	flag.StringVar(&action, "ac", _ACTION_COMPUTER_CODE, "`生成计算机签名(code)，导入授权(import)，读取授权(read)`")
	flag.StringVar(&codeOutput, "co", "", "计算机信息文件创建路径，默认为执行路径")
	flag.StringVar(&inputFile, "if", "license.vec", "导入授权时授权文件的路径")
	flag.StringVar(&pubKeyFile, "pf", "pub.pem", "读取授权时的公钥文件路径")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	f := func() error {
		switch action {
		case _ACTION_COMPUTER_CODE:
			return actionCode()
		case _ACTION_IMPORT:
			return actionImport()
		case _ACTION_Read:
			return actionRead()
		default:
			return errors.New("不支持的操作")
		}
	}
	err := f()
	if err != nil {
		fmt.Printf("\n发生错误: %v\n", err)
		return
	}
}

func actionCode() (err error) {
	sl, err := salt.Get()
	if err != nil {
		return
	}
	f, err := os.OpenFile(_CODE_OUTPUT_NAME, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(sl)
	fmt.Printf("\n计算机信息文件生成成功: %s\n", _CODE_OUTPUT_NAME)
	return
}

func actionImport() (err error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return
	}
	defer f.Close()
	err = wr.Write(f)
	fmt.Println("授权导入成功")
	return
}

func actionRead() (err error) {
	salt := "vector"
	f, err := os.Open(pubKeyFile)
	if err != nil {
		return
	}
	pubKey, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	ret, err := wr.Read([]byte(salt), pubKey)
	if err != nil {
		return
	}
	fmt.Printf("\n授权读取成功: %s\n", string(ret))
	return
}