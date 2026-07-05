package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func Kimi(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:       "Kimi",
		hostname:   "kimi.moonshot.cn",
		url:        "https://kimi.moonshot.cn/",
		noRedirect: true,
		okCodes: map[int]bool{
			http.StatusOK:               true,
			http.StatusMovedPermanently: true,
			http.StatusFound:            true,
		},
		noCodes: map[int]bool{http.StatusForbidden: true},
	})
}
