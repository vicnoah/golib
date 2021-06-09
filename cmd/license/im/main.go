package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/vicnoah/golib/helper/vbrowser"

	"github.com/vicnoah/golib/pkg/sec/license/salt"
	"github.com/vicnoah/golib/pkg/sec/license/wr"
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
	runHTTP()
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
	sl, err := salt.Get()
	if err != nil {
		return
	}
	f, err := os.Open(pubKeyFile)
	if err != nil {
		return
	}
	pubKey, err := io.ReadAll(f)
	if err != nil {
		return
	}
	ret, err := wr.Read(sl, pubKey)
	if err != nil {
		return
	}
	fmt.Printf("\n授权读取成功: %s\n", string(ret))
	return
}

func runHTTP() {
	http.HandleFunc("/", html)
	http.HandleFunc("/download", download)
	http.HandleFunc("/upload", upload)
	go http.ListenAndServe("127.0.0.1:8081", nil)
	vbrowser.OpenURL("http://localhost:8081")
	select {}
}

func download(w http.ResponseWriter, r *http.Request) {
	sl, err := salt.Get()
	if err != nil {
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+_CODE_OUTPUT_NAME)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(sl)))
	w.Write(sl)
}

func upload(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	f, _, err := r.FormFile("upload")
	if err != nil {
		httpStatus502(&w)
	}
	err = wr.Write(f)
	if err != nil {
		w.Write([]byte("授权导入失败"))
		return
	}
	w.Write([]byte("授权导入成功"))
}

func html(w http.ResponseWriter, r *http.Request) {
	var context = `
<html>
<head>
<title>授权导入导出工具</title>
</head>
<body>
<a href="/download">下载签名</a>
<form action="/upload" method="post" enctype="multipart/form-data">
    <p><input type="file" name="upload"></p>
    <p><input type="submit" value="上传授权"></p>
</form>
</body>
</html>
	`
	_, err := w.Write([]byte(context))
	if err != nil {
		httpStatus502(&w)
	}
}

func httpStatus502(w *http.ResponseWriter) {
	(*w).WriteHeader(http.StatusBadGateway)
}
