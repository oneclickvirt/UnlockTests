package tw

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// BahamutAnime
// ani.gamer.com.tw 仅 ipv4 且 get 请求 有问题
// 存在 cloudflare 的质询防御，非5秒盾，无法突破，需要js动态加载
func BahamutAnime(c *http.Client) model.Result {
	name := "Bahamut Anime"
	if c == nil {
		return model.Result{Name: name}
	}
	url := "https://ani.gamer.com.tw/ajax/getdeviceid.php"
	headers := map[string]string{
		"User-Agent": model.UA_Browser,
	}
	client := utils.Req(c)
	client = utils.SetReqHeaders(client, headers)
	resp, err := client.R().Get(url)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body := string(b)
	tempList := strings.Split(body, "\n")
	for _, line := range tempList {
		if strings.Contains(line, "deviceid") {
			fmt.Println(line)
		}
	}
	var res struct {
		Deviceid string `json:"deviceid"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	var res2 struct {
		AnimeSn int `json:"animeSn"`
	}
	json.Unmarshal([]byte(body), &res2)
	resp2, err2 := client.R().Get("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667&device=" + res.Deviceid)
	if err2 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
	}
	defer resp2.Body.Close()
	b, err = io.ReadAll(resp2.Body)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	body2 := string(b)
	if err := json.Unmarshal([]byte(body2), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	if res2.AnimeSn != 0 {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if res2.AnimeSn == 0 || resp2.StatusCode == 403 || resp2.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ani.gamer.com.tw failed with code: %d", resp.StatusCode)}
}
