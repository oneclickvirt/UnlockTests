package id

import (
	"github.com/oneclickvirt/UnlockTests/asia"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// BilibiliID
// 检测印度尼西亚B站是否可用
func BilibiliID(request *gorequest.SuperAgent) model.Result {
	return asia.Bilibili(request, "BilibiliID", "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=11130043")
}
