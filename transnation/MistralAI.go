package transnation

import (
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
)

func MistralAI(c *http.Client) model.Result {
	return checkAIStatus(c, aiStatusProbe{
		name:     "Mistral AI",
		hostname: "chat.mistral.ai",
		url:      "https://chat.mistral.ai/",
		okCodes:  map[int]bool{http.StatusOK: true},
		noCodes:  map[int]bool{http.StatusForbidden: true},
	})
}
