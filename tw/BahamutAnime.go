package tw

// BahamutAnime
// ani.gamer.com.tw 仅 ipv4 且 get 请求 有问题
// 存在 cloudflare 的质询防御，非5秒盾，无法突破，需要js动态加载
// func BahamutAnime(request *gorequest.SuperAgent) model.Result {
// 	name := "Bahamut Anime"
// 	url := "https://ani.gamer.com.tw/ajax/getdeviceid.php"
// 	request = request.Set("User-Agent", model.UA_Browser).Timeout(15 * time.Second)
// 	resp, body, errs := request.Get(url).End()
// 	if len(errs) > 0 {
// 		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
// 	}
// 	defer resp.Body.Close()
// 	tempList := strings.Split(body, "\n")
// 	for _, line := range tempList {
// 		if strings.Contains(line, "deviceid") {
// 			fmt.Println(line)
// 		}
// 	}
// 	var res struct {
// 		Deviceid string `json:"deviceid"`
// 	}
// 	if err := json.Unmarshal([]byte(body), &res); err != nil {
// 		return model.Result{Name: name, Status: model.StatusErr, Err: err}
// 	}
// 	var res2 struct {
// 		AnimeSn int `json:"animeSn"`
// 	}
// 	json.Unmarshal([]byte(body), &res2)
// 	cookies := resp.Request.Cookies()
// 	request.AddCookies(cookies)
// 	resp2, body2, errs2 := request.Get("https://ani.gamer.com.tw/ajax/token.php?adID=89422&sn=14667&device=" + res.Deviceid).End()
// 	if len(errs2) > 0 {
// 		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs2[0]}
// 	}
// 	defer resp2.Body.Close()
// 	if err := json.Unmarshal([]byte(body2), &res); err != nil {
// 		return model.Result{Name: name, Status: model.StatusErr, Err: err}
// 	}
// 	if res2.AnimeSn != 0 {
// 		return model.Result{Name: name, Status: model.StatusYes}
// 	} else if res2.AnimeSn == 0 || resp2.StatusCode == 403 || resp2.StatusCode == 404 {
// 		return model.Result{Name: name, Status: model.StatusNo}
// 	}
// 	return model.Result{Name: name, Status: model.StatusUnexpected,
// 		Err: fmt.Errorf("get ani.gamer.com.tw failed with code: %d", resp.StatusCode)}
// }
