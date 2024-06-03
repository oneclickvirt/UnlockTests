package fr

// Salto
// geo.salto.fr 双栈 get 请求 有问题
// tls验证失败，识别失效，未知原因
// func Salto(request *gorequest.SuperAgent) model.Result {
// 	name := "Salto"
// 	url := "https://www.salto.fr/"
// 	client := req.DefaultClient()
// 	client.ImpersonateChrome()
// 	client.Headers.Set("User-Agent", model.UA_Browser)
// 	resp, err := req.R().
		// SetRetryCount(2).
		// SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		// SetRetryFixedInterval(2 * time.Second).Get(url)
// 	if err != nil {
// 		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
// 	}
// 	defer resp.Body.Close()
// 	b, err := io.ReadAll(resp.Body)
// 	if err == nil {
// 		fmt.Println(string(b), resp.StatusCode)
// 	}
// 	return model.Result{Name: name, Status: model.StatusUnexpected,
// 		Err: fmt.Errorf("get geo.salto.fr failed with code: %d", resp.StatusCode)}
// }
