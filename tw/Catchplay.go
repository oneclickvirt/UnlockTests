package tw

// Catchplay
// sunapi.catchplay.com 仅 ipv4 且 get 请求
// unauthorized 有问题
// func Catchplay(request *gorequest.SuperAgent) model.Result {
// 	name := "CatchPlay+"
// 	url := "https://sunapi.catchplay.com/geo"
// 	request.Set("authorization",
// 		"Basic NTQ3MzM0NDgtYTU3Yi00MjU2LWE4MTEtMzdlYzNkNjJmM2E0Ok90QzR3elJRR2hLQ01sSDc2VEoy")
// 	resp, body, errs := request.Get(url).End()
// 	if len(errs) > 0 {
// 		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
// 	}
// 	defer resp.Body.Close()
// 	var res struct {
// 		Code string `json:"code"`
// 	}
// 	if err := json.Unmarshal([]byte(body), &res); err != nil {
// 		if strings.Contains(body, "is not allowed") && strings.Contains(body, "The location") {
// 			return model.Result{Name: name, Status: model.StatusNo}
// 		}
// 		fmt.Println(body)
// 		return model.Result{Name: name, Status: model.StatusErr, Err: err}
// 	}
// 	if res.Code == "100016" {
// 		return model.Result{Name: name, Status: model.StatusNo}
// 	} else if res.Code == "0" {
// 		return model.Result{Name: name, Status: model.StatusYes}
// 	}
// 	return model.Result{Name: name, Status: model.StatusUnexpected,
// 		Err: fmt.Errorf("get sunapi.catchplay.com failed with code: %d", resp.StatusCode)}
// }
