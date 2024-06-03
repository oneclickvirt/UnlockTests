package th

import (
	"github.com/oneclickvirt/UnlockTests/asia"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// BilibiliTH
// 检测泰国B站是否可用
func BilibiliTH(request *gorequest.SuperAgent) model.Result {
	return asia.Bilibili(request, "BilibiliTH", "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=10077726")
}
