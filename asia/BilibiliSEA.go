package asia

import (
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// BilibiliSEA
// 检测东南亚B站是否可用
func BilibiliSEA(request *gorequest.SuperAgent) model.Result {
	return Bilibili(request, "https://api.bilibili.tv/intl/gateway/web/playurl?s_locale=en_US&platform=web&ep_id=347666")
}
