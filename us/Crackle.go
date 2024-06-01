package us

// Crackle
// prod-api.crackle.com 双栈 get 请求 有问题
// {"path":"/appconfig","version":"v2.0.0","status":"400","timestamp":"2024-05-31T10:28:34.542Z","error":{"message":"Platform Key is not specified","type":"ApiError","code":121,"details":{}}}
//func Crackle(request *gorequest.SuperAgent) model.Result {
//	name := "Crackle"
//	url := "https://prod-api.crackle.com/appconfig"
//	request = request.Set("User-Agent", model.UA_Browser)
//	resp, body, errs := request.Get(url).Retry(2, 5).End()
//	if len(errs) > 0 {
//		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
//	}
//	defer resp.Body.Close()
//	// fmt.Println(body)
//	if strings.Contains(body, "302 Found") || resp.StatusCode == 403 || resp.StatusCode == 451 {
//		return model.Result{Name: name, Status: model.StatusNo}
//	} else if resp.StatusCode == 200 {
//		return model.Result{Name: name, Status: model.StatusYes}
//	}
//	return model.Result{Name: name, Status: model.StatusUnexpected,
//		Err: fmt.Errorf("get prod-api.crackle.com failed with code: %d", resp.StatusCode)}
//}
