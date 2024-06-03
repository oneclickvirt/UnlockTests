package vn

import (
	"github.com/oneclickvirt/UnlockTests/asia"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// BilibiliVN
// 检测越南B站是否可用
func BilibiliVN(request *gorequest.SuperAgent) model.Result {
	return asia.Bilibili(request, "BilibiliVN", "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11405745")
}
