package stuc

import "encoding/json"

const (
	NeverExpire = 0x01 // 不过期
	NeedsToExpire = 0x02 // 需要过期
)

/*
{
"app": "dinsight",
"version": "0.1.1",
"award": "减灾中心",
"type": 1,
"expire": 360000000,
"permissions": 1
}
*/

type License struct {
	App string `json:"app"` // app名称
	Version string `json:"version"` // app版本
	Award string `json:"award"` // 颁发给用户的名称
	Type int `json:"type"` // 授权类型
	StartTime int64 `json:"start_time"` // 授权开始时间
	Expire int `json:"expire"` // 过期时间
	Permissions string `json:"permissions"` // 权限信息
}

// New Generate JSON encoded authorization data
func New(app, ver, awd string, typ, exp int, start int64, per string) (ln []byte, err error) {
	var lcs License
	lcs.App = app
	lcs.Version = ver
	lcs.Award = awd
	lcs.Type = typ
	lcs.Expire = exp
	lcs.StartTime = start
	lcs.Permissions = per
	ln, err = json.Marshal(&lcs)
	return
}
