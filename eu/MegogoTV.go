package eu

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"io"
	"net/http"
	"strings"
)

// MegogoTV
// megogo.net 仅 ipv4 且 get 请求
func MegogoTV(c *http.Client) model.Result {
    name := "MegogoTV"
    hostname := "megogo.net"
    if c == nil {
        return model.Result{Name: name}
    }
    url := "https://megogo.net/en"
    client := utils.Req(c)
    resp, err := client.R().Get(url)
    if err != nil {
        return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
    }
    defer resp.Body.Close()
    b, err := io.ReadAll(resp.Body)
    if err != nil {
        return model.Result{Name: name, Status: model.StatusNetworkErr, Err: fmt.Errorf("can not parse body")}
    }
    body := string(b)
    if strings.Contains(strings.ToLower(body), "vpn") {
        return model.Result{Name: name, Status: model.StatusNo, Err: fmt.Errorf("vpn detected")}
    }
    if resp.StatusCode == 200 {
        result1, result2, result3 := utils.CheckDNS(hostname)
        unlockType := utils.GetUnlockType(result1, result2, result3)
        return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
    }
    return model.Result{Name: name, Status: model.StatusUnexpected, 
        Err: fmt.Errorf("get megogo.net failed with code: %d", resp.StatusCode)}
}