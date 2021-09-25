package api

import (
	"gochat/model"

	"github.com/gogf/gf/encoding/gjson"
)

var SystemApi systemApi

type systemApi struct{}

// 心跳
func (this systemApi) HeartJump(user *model.UserInfo, data *gjson.Json) {
	user.Conn.ResetHeartJump()
}
