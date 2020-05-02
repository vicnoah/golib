package license

import (
	"bytes"
	"git.sabertrain.com/vector-dev/golib/sec/rsa"
	"io/ioutil"
	"os"
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

func TestDecryptionLicense(t *testing.T)  {
	pubKey := `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCt7uZyFteXRPZdXTqqMw7ND2IC
p0Si36x/WXMXTNlI2goE4qrTkoWU7s1aNJxIqJItVBUgfzQ6ICkKG/vwcgQ9qf2s
l9iwl+yVDykrL9H3eVrVqE2Zcq3TQnB9g0cXUexZCYgwnxgwXJWxArl5iYU/jH0D
01FPN0H8XqwDb5kGwwIDAQAB
-----END PUBLIC KEY-----`
	result := `{"app":"","version":"","award":"","type":1,"start_time":1582625678,"expire":31536000,"permissions":"dsfjfjfsdjf\n"}`
	salt := "vector"
	licenseFile, err := os.Open("license.vec")
	if err != nil {
		return
	}
	defer licenseFile.Close()
	data, err := ioutil.ReadAll(licenseFile)
	if err != nil {
		return
	}
	Convey("test read license", t, func() {
		ret, er := Decrypt(data, []byte(salt), []byte(pubKey))
		So(er, ShouldBeNil)
		So(string(ret), ShouldEqual, result)
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