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
	url := "https://cdi.ofiii.com/ofiii_cdi/video/urls?device_type=pc&device_id=b4e377ac-8870-43a4-957a-07f95549a03d&media_type=comic&asset_id=vod68157-020020M001&project_num=OFWEB00&puid=dcafe020-e335-49fb-b9c7-52bd9a15c305"
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
	case 400, 403:
		return model.Result{Name: name, Status: model.StatusNo}
	case 200:
		return model.Result{Name: name, Status: model.StatusYes}
	default:
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: fmt.Errorf("unexpected code: %d", resp.StatusCode)}
	}
}
