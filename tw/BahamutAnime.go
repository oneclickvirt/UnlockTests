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
	client := utils.Req(c)
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
	// fmt.Println(body)
	//tempList := strings.Split(body, "\n")
	//for _, line := range tempList {
	//	if strings.Contains(line, "deviceid") {
	//		fmt.Println(line)
	//	}
	//}
	var res struct {
		Deviceid string `json:"deviceid"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	fmt.Println(res.Deviceid)

	// 14667
	sn := "37783"
	resp2, err2 := client.R().Get("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=" + sn + "&device=" + res.Deviceid)
	if err2 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
	}
	defer resp2.Body.Close()
	b2, err2 := io.ReadAll(resp2.Body)
	if err2 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err2}
	}
	body2 := string(b2)

	resp3, err3 := client.R().Get("https://ani.gamer.com.tw/")
	if err3 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err3}
	}
	defer resp3.Body.Close()
	b3, err3 := io.ReadAll(resp3.Body)
	if err3 != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err3}
	}
	body3 := string(b3)
	// var res3 struct {
	// 	AnimeSn int `json:"animeSn"`
	// }
	// if err := json.Unmarshal(b3, &res3); err != nil {
	// 	return model.Result{Name: name, Status: model.StatusErr, Err: err}
	// }
	// fmt.Println(res3.AnimeSn)
	fmt.Println(body3)
	if !strings.Contains(body2, "animeSn") {
		return model.Result{Name: name, Status: model.StatusYes}
	} else if resp2.StatusCode == 403 || resp2.StatusCode == 404 {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get ani.gamer.com.tw failed with code: %d", resp.StatusCode)}
}
