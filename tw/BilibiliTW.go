package tw

import (
	"github.com/oneclickvirt/UnlockTests/asia"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

func BilibiliTW(request *gorequest.SuperAgent) model.Result {
	return asia.Bilibili(request, "https://api.bilibili.com/pgc/player/web/playurl?avid=50762638&cid=100279344&qn=0&type=&otype=json&ep_id=268176&fourk=1&fnver=0&fnval=16&module=bangumi")
}
