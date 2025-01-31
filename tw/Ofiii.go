package tw

import (
	"fmt"
	"net/http"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// Ofiii
func Ofiii(c *http.Client) model.Result {
	name := "Ofiii"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://ntdofifreepc.akamaized.net"
	client := utils.Req(c)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return model.Result{Status: model.StatusNetworkErr, Err: err}
	// }
	switch resp.StatusCode {
	case 403:
		return model.Result{Name: name, Status: model.StatusNo}
	case 451:
		return model.Result{Name: name, Status: model.StatusYes}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("unexpected code: %d", resp.StatusCode)}
	}
}
