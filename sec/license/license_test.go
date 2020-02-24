package license

import (
	"bytes"
	"e.coding.net/vector-tech/golib/sec/rsa"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestLicense(t *testing.T) {
	var datas = []struct{
		Data []byte
		Pass []byte
		Salt []byte
		privateKey []byte
		publicKey []byte
	}{
		{
			Data:[]byte("123"),
			Pass:[]byte("-123456"),
			Salt: []byte("1"),
		},
		{
			Data:[]byte("1234567890123456"),
			Pass:[]byte("-123456"),
			Salt: []byte("1"),
		},
		{
			Data:[]byte(`{
"name": "vector",
"expire": 360000000,
"level": 1
}`),
			Pass:[]byte("wwg9776586516"),
			Salt: []byte("afjife"),
		},
	}
	for index, _ := range datas {
		priv, pub, er := getKey()
		if er != nil {
			return
		}
		datas[index].privateKey = priv
		datas[index].publicKey = pub
	}
	Convey("Test license creation and reading", t, func() {
		for _, data := range datas {
			cipher, er := Encrypt(data.Data, data.Pass, data.Salt, data.privateKey)
			So(er, ShouldBeNil)
			ret, er := Decrypt(cipher, data.Salt, data.publicKey)
			So(er, ShouldBeNil)
			So(ret, ShouldResemble, data.Data)
			ret1, er := Decrypt(cipher, []byte("1243"), data.publicKey)
			So(er, ShouldBeNil)
			So(ret1, ShouldNotResemble, data.Data)
			ret2, er := Decrypt([]byte("dsfsdj"), []byte("1243"), data.publicKey)
			So(er, ShouldNotBeNil)
			So(ret2, ShouldNotResemble, data.Data)
		}
	})
}

func getKey() (privateKey, publicKey []byte, err error) {
	var (
		priv bytes.Buffer
		pub bytes.Buffer
	)
	err = rsa.GenRsaKey(1024, &priv, &pub)
	if err != nil {
		return
	}
	privateKey = priv.Bytes()
	publicKey = pub.Bytes()
	return
}