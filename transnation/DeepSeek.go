package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func DeepSeek(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "DeepSeek",
		hostname: "chat.deepseek.com",
		url:      "https://chat.deepseek.com/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}
