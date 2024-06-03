package hk

import (
	"github.com/oneclickvirt/UnlockTests/asia"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

func BilibiliHKMO(request *gorequest.SuperAgent) model.Result {
	return asia.Bilibili(request, "BiliBili HongKong/Macau Only", "https://api.bilibili.com/pgc/player/web/playurl?avid=473502608&cid=845838026&qn=0&type=&otype=json&ep_id=678506&fourk=1&fnver=0&fnval=16&module=bangumi")
}
