package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"git.sabertrain.com/vector-tech/golib/sec/license/cmd/raw_license/stuc"
	"time"
)

const (
	_OneDay = 3600*24
	_OneYear = 365
	_Dinsight = "dinsight"
)

var (
	help bool
	app string
	appName string
	appVersion string
	startTime string
	award string
	isExp bool
	expire int
	perName string
	outputName string
)

func init()  {
	flag.BoolVar(&help, "help", false, "help")
	flag.StringVar(&app, "a", "dinsight", "Need to create authorized software name")
	flag.StringVar(&appName, "an", "", "`Software Name`")
	flag.StringVar(&appVersion, "av", "", "`Software version number`")
	flag.StringVar(&award, "ad", "", "`Authorized target user name`")
	flag.BoolVar(&isExp, "ie", false, "`Whether to expire or not`")
	flag.IntVar(&expire, "exp", _OneYear, "`Expiration time unit days`")
	flag.StringVar(&perName, "per", "", "`permission file path`")
	flag.StringVar(&outputName, "o", "license.raw", "`Original license file generation path`")
	flag.StringVar(&startTime, "st", "", "`Authorization start time, format: 2009-01-01, default is empty, empty parameter will directly take system time`")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	var startTimestamp int64
	if startTime == "" {
		startTimestamp = time.Now().Unix()
	} else {
		t, er := time.Parse("2006-01-02", startTime)
		if er != nil {
			fmt.Printf("\n时间解析错误: %v\n", er)
			return
		}
		startTimestamp = t.Unix()
	}
	if perName == "" {
		fmt.Printf("\n%s\n","权限文件不能为空")
		return
	}
	perFile, err := os.Open(perName)
	if err != nil {
		fmt.Printf("\n打开权限文件失败: %v\n", err)
		return
	}
	defer perFile.Close()
	per, err := ioutil.ReadAll(perFile)
	if err != nil {
		fmt.Printf("\n读取权限文件失败: %v\n", err)
		return
	}
	outputFile, err := os.OpenFile(outputName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("\n创建原始授权文件失败: %v\n", err)
		return
	}
	expireSed := _OneDay * expire
	var typNum int
	if isExp {
		typNum = stuc.NeedsToExpire
	} else {
		typNum = stuc.NeverExpire
	}
	ln, err := gen(app, per, typNum, expireSed, startTimestamp, appName, appVersion, award)
	if err != nil {
		fmt.Printf("\n密钥生成失败: %v\n", err)
		return
	}
	_, err = outputFile.Write(ln)
	if err != nil {
		fmt.Printf("\n写入密钥文件失败: %v\n", err)
		return
	}
	fmt.Printf("\n密钥文件生成成功,路径: %s\n", outputName)
}

func gen(app string,
	per []byte,
	typ int,
	exp int,
	start int64,
	appName string,
	appVersion string,
	award string) (ln []byte, err error) {
	switch app {
	case _Dinsight:
		ln, err = stuc.Dinsight(appName, appVersion, award, typ, exp, start, string(per))
		return
	default:
		err = errors.New("Unsupported software type")
		return
	}
}