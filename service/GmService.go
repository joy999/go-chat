package service

import (
	"fmt"
	"gochat/model"
	"gochat/utils"
	"strconv"
	"strings"
	"time"
)

var GmService gmService

type gmService struct {
}

func (this *gmService) DoCmd(user *model.UserInfo, msg string) {
	args := strings.Split(msg, " ")
	switch args[0] {
	case "popular":
		if n, err := strconv.Atoi(args[1]); err == nil {
			//执行处理
			user.Conn.SendSystemMsg(fmt.Sprintf("%ds内出现频率最高的词：%s", n, FilterService.PopularWords(n)))
		}
	case "stats":
		uname := args[1]
		info := UserService.Find(uname)
		long := time.Now().Unix() - info.CreateTime.Unix()
		user.Conn.SendSystemMsg(fmt.Sprintf("用户【%s】在线：%s", uname, utils.FormatLong(long)))
	}
}
