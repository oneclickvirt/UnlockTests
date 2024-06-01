package us

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
	"strings"
)

// Hulu
// www.hulu.com 仅 ipv4 且 post 请求
func Hulu(request *gorequest.SuperAgent) model.Result {
	name := "Hulu"
	resp, body, errs := request.Get("https://www.hulu.com").
		//Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9").
		Set("User-Agent", model.UA_Browser).
		Set("Accept-Encoding", "gzip, deflate, br").
		Set("Cache-Control", "no-cache").
		Set("DNT", "1").
		Set("Pragma", "no-cache").
		Set("Sec-CH-UA", `"Chromium";v="106", "Google Chrome";v="106", "Not;A=Brand";v="99"`).
		Set("Sec-CH-UA-Mobile", "?0").
		Set("Sec-CH-UA-Platform", "Windows").
		Set("Sec-Fetch-Dest", "document").
		Set("Sec-Fetch-Mode", "navigate").
		Set("Sec-Fetch-Site", "none").
		Set("Sec-Fetch-User", "?1").
		Set("Upgrade-Insecure-Requests", "1").
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		request = gorequest.New()
		request = request.Set("User-Agent", model.UA_Browser).
			Set("Accept", "application/json").
			Set("Accept-Language", "zh-CN,zh;q=0.9").
			Set("Connection", "keep-alive").
			Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8").
			Set("Origin", "https://www.hulu.com").
			Set("Referer", "https://www.hulu.com/welcome").
			Set("Sec-Fetch-Dest", "empty").
			Set("Sec-Fetch-Mode", "cors").
			Set("Sec-Fetch-Site", "same-site").
			Set("sec-ch-ua", model.UA_SecCHUA).
			Set("sec-ch-ua-mobile", "?0").
			Set("sec-ch-ua-platform", "\"Windows\"").
			Set("Cookie", "_hulu_at=eyJhbGciOiJSUzI1NiJ9.eyJhc3NpZ25tZW50cyI6ImV5SjJNU0k2VzExOSIsInJlZnJlc2hfaW50ZXJ2YWwiOjg2NDAwMDAwLCJ0b2tlbl9pZCI6IjQyZDk0YzA5LWYyZTEtNDdmNC1iYzU4LWUwNTA2NGNhYTdhZCIsImFub255bW91c19pZCI6IjYzNDUzMjA2LWFmYzgtNDU4Yi1iODBkLWNiMzk2MmYzZGQyZCIsImlzc3VlZF9hdCI6MTcwNjYwOTUzODc5MiwidHRsIjozMTUzNjAwMDAwMCwiZGV2aWNlX3VwcGVyIjoxfQ.e7sRCOndgn1j30XYkenLcLSQ7vwc2PXk-gFHMIF2gu_3UNEJ3pp3xNOZMN0n7DQRw5Jv68WiGxIvf65s8AetOoD4NLt4sZUDDz9HCRmFHzpmAJdtXWZ-HZ4fYucENuqDDDrsdQ-FCc0mgIe2IXkmQJ6tpIN3Zgcgmpmbeoq6jYyLlqg6f8eMsI1bNAsBGGj-9DXw2PMotlYHWB22pw2NRfJw1TjWXwywRBodAOve7rsu2Vhx-A2-OH4GplRvxLqzCpl2pcjkYg9atmUB7jnNIf_jHqlek4oRRawahWq-2vWnWmb1eMQcH-v2IHs3YdVk7I-t4iS19auPQrdgo6jPaA; _hulu_assignments=eyJ2MSI6W119; bm_mi=8896E057E2FC39F20852615A0C46A2B4~YAAQZB0gFyQrnlSNAQAAU/naWRaCMaJFIPi3L8KYSiCL3AN7XHCPw0gKvdIl0XZ/VE3QiKEr31qjm9sPulHbdQ4XXIXPXZ53DpIK43fLybrT6WxIpmGz3iThk6+xefI2dPLzwBAdoTrsbAbHC2q4LDx0SBM+n21LvTD7UnT2+DyVBK75YCDJJKHlJ5jzB3Q81JIlmqfTzibjgVmPIxXrFdTs5Ll8mtp6WzE3VDISmjGjTRTrSOVYM0YGpyhye1nsm3zBCO13vDjKMCJ/6oAsVqBfgfW07e7sWkWeUiDYLUifRDymc4GaMhavenBvCma/G1qW~1; bm_sv=FEE04D9D797D0237C312D77F57DABBFD~YAAQZB0gFyUrnlSNAQAAU/naWRaMNI8KmoGX9XNJkm9x9VeeGzGQyPfu49M9MnLObz8D4ZYk9Td+3Y8Z/Jfx+kl2qOPXmtOC5GZpA++9bxUKV0SwaoGhivl+ibIJSQTc7lw4kzdM/2w8b3rwItRaHXFa+shMtD3eiKvBePrqCiezucqrcss1U4ojLKEOvcsKJGt6ZTGGs2H+Qu6cyns9BVN0BprMHRY3njHXyxbFIcGy8Lq7aPn6nuZ0ehfZ9Q==~1; ak_bmsc=55F791116713DDB91AB0978225853B77~000000000000000000000000000000~YAAQZB0gF6ErnlSNAQAAHALbWRaA625r4bWVW8g2gHV797RN8bfCwNy6KfnGEucUPiPt4QKjJUldR6lyaM7sarag6A7WLqxEFr/zAFlPQI12Uxsqdzg3IgU0R8g2eMQRnRoGMNSUPyt4rdCWWwGjEcM+dQ8TI+y1vKw9dLXoBJAHofaWe/dZhY4fx2mYKhKFibvdpwJT6UPe4rBz8igd9oTQBn69Ebi6/9YFykqGuKsllxa5+QZWczb0+HLLDRKV4CkZdhbFj0yljEOyz4GHqqP8qg3Xa3lCKzdzsrmPn6zdFbgzCE8HsyPjsmy+/rRfFxagH5rYudLqFXg5o5dXFFJPTiLXtZ/S30ckc/OUWk4JP2ywAQVm/zbp8nlRVMFDEdjIPh/F+5QXfYBV+yL4a85ThlBEXSr54/QWXiHxBRiOwhv2ydoZDfT78r9bUHbMOra37C0xutfo37fbYEw9LWlLdZCub9U5HA/zSeIN3KxrZr0yNKfJjOau7BqdHL+AuvDj134ZPZPVig==; _customer_type=anonymous; s_fid=66C80912997F4CF8-2D3140F8EDC76274; s_cc=true; _rdt_uuid=1706609517486.d5b309e4-2b0b-440f-9817-cf619e4ce15d; _gcl_au=1.1.602757068.1706609518; _scid=cc980fef-26dc-479a-b9a8-b0e531c87cd3; _scid_r=cc980fef-26dc-479a-b9a8-b0e531c87cd3; _tt_enable_cookie=1; _ttp=1h5M9exzlSz7wAFDR78KCHCsnDC; utag_main=v_id:018d59da9a5c00215e601dada5700507d001c07500bd0$_sn:1$_ss:0$_st:1706611329541$ses_id:1706609515101%3Bexp-session$_pn:1%3Bexp-session$_prevpage:%2Fwelcome%3Bexp-1706613129564$trial_duration:undefined%3Bexp-session$program_id:undefined%3Bexp-session$vapi_domain:hulu.com$g_sync_ran:1%3Bexp-session$dc_visit:1$dc_event:1%3Bexp-session$dc_region:ap-east-1%3Bexp-session; _hulu_metrics_context_v1_=%7B%22cookie_session_guid%22%3A%227dc4f3a6826f2c35125268f5ddab1849%22%2C%22referrer_url%22%3A%22%22%2C%22curr_page_uri%22%3A%22www.hulu.com%2Fwelcome%22%2C%22primary_ref_page_uri%22%3Anull%2C%22secondary_ref_page_uri%22%3Anull%2C%22curr_page_type%22%3A%22landing%22%2C%22primary_ref_page_type%22%3Anull%2C%22secondary_ref_page_type%22%3Anull%2C%22secondary_ref_click%22%3A%22%22%2C%22primary_ref_click%22%3A%22%22%7D; metrics_tracker_session_manager=%7B%22session_id%22%3A%22B26515EB8A7952D4D35F374465362A72-529671c4-c8c2-4c7c-8bff-cc201bcd4075%22%2C%22creation_time%22%3A1706609513429%2C%22visit_count%22%3A1%2C%22session_seq%22%3A4%2C%22idle_time%22%3A1706609529579%7D; guid=B26515EB8A7952D4D35F374465362A72; JSESSIONID=ED7031784C3B1843BFC9AACBB156C6BA; s_sq=wdghuluwebprod%3D%2526c.%2526a.%2526activitymap.%2526page%253Dwelcome%2526link%253DLOG%252520IN%2526region%253Dlogin-modal%2526pageIDType%253D1%2526.activitymap%2526.a%2526.c%2526pid%253Dwelcome%2526pidt%253D1%2526oid%253Dfunctionsn%252528%252529%25257B%25257D%2526oidt%253D2%2526ot%253DBUTTON; XSRF-TOKEN=bcfa1766-1f73-442d-a71b-e1cf6c275f45; _h_csrf_id=2a52618e9d006ac2e0b3e65740aa55e2584359553466051c3b01a2f1fb91726a")
		resp, body, errs = request.Post("https://auth.hulu.com/v4/web/password/authenticate").
			Send("user_email=me%40jamchoi.cc&password=Jam0.5cm~&recaptcha_type=web_invisible&rrventerprise=03AFcWeA6UFet_b_82RUmGfFWJCWuqy6kIn854Rhqjwd7vrkjH6Vku1wBZy8-FBA3Efx1p2cuNnKeJRuk7yJWm-xZgFfUx0Wdj2OAhiGvIdWrcpfeuQSXEqaXH4FKdmAHVZ3EqHwe5-h_zgtcyIxq-Nn1-sjeUfx1Y7QyVkb_GWJcr0GLoKgTFLzbF4kmJ8Qsi4IFx9hyYo9TFbBqtYdgxCI2q9DnEzOHrxK-987PEY8qzsR08Hrb9oDvumqLp1gs4uEVTwDKWt37aNB3CMVBKL2lHj7n768kXpgkXFDIhcM2eiJJ-H22qxKzNUpg-Q_N1xzkYsavCJG3ckQgsCTRRP2NU3nIERTWDTVXRBEq52-_ZQWu_Ds4W4UZyP0hEhCD2gambN4YJqEzaeHdGPwOR943nFbG6GILBx4vY-UUc7zjMf2HRjkNPvPpQiHMIYo21JXq6l8-IWyTeHY26NU6M4vCCbzwZEsdSln48rXM_fdBcDHC-8AxUFuBR8j3DMsB6Q3xMS2EHeGVrmhDY1izDNJZsVC_cN0W2tRneOJmni7ZU1iAYoBAGBBM5FDTE4UbYUTnuUn-htm9Q0RzukpYTumF_WwQ3HnEL0JK1Q1xea-hteI8lB4oAkhVOBOHVPii9atdZR9ZLpxRh1pdy3Lwmr1ltsubxE05wqmrmt33P2WsvH_3nBJXC_FhTD06BxT60RuiGtFr2gscHjjl_NCa1F-Dv9Hgi5ek2nLHK37a84bRSoKwLL3Lnpi9byuBntlpf-UXj7nveawKZmZTUBOSc7j6Vmmf124DTPJXsFeofMfUXkqTauPTWJBOz0OdKnLKDHMSsk7oSJVKsDUEeq0iKMdtCMBPvQBaPYAb79LDRwv_ereqyklKcUKQxeZRZmEXLKIWp8BS4U9uTXA2w8hwZWe7goLnUBQATIwojeHKpypSLnzQBu9JCwMU4aXfKIplL8sXuAx3QFD52eGZSCEyuFXP3ACN53QOlTAjjlP2eDT9fEwWHT4o8eJfviyjvm8xDmzKtq4F3u5XB3tL86-dK40XYbGcTI0Irw1nz1UTcxplFgHQgb6i8WEAqb69CQkpGWAUlmnknBirRAv2adqPaW2d_lv6L3Eo-ZupWcZ9Cu4PibM5BruVNXifBwPNPXHKw-sWBj-UP1g9VtxHVEVwoTXrbB-lT8EvjDEDQKrvOwnri4_tzVzn6YKvQMELbxSegvmc2w7xypT2qFzKRFXqwTMLT9d0rf2p9tbwbe39REMR8oI7wPfbjyJjK2XF4DmEAyVvBMuJlBaBsKBs5VynITHFWs4xvkAOe4jO_fzkKXzB6F6DB03ldasxbrNK_cepUOF6FD39-pHvbAGcoTrDrx6FSfecYXwSvc3GxM3IHSKwISKWav2iqPMtIt6ClCgUPgTCBDng2ZptXeVG8FckGIGMEdVlgGt5DG2tdMO2p8Hs5tKXuu8anc_csaaSfLIQ1_kav0dp8vpSXhCxeg899o5coXderUoIBcUsfaBJJm80YnCAc4LaM8HmYtJBcKqCC_uwCckPDOuC0SQy3d07LEi6wyifvY0Kv_-ER6wXvhNWnDZIXJNlH2369X7y8o3y2HMisOwAhfmKN7_ZAaODEOO-5x9JHocAYnt4a8_focwU9JQ_hUQgtdzYpP1ACEqxVjJb0A0NlABpm-CG8V9n9y6XpZkGQiMYJIH3jr6VilHSEM9rQSEv6LN8NFigl3-5Y4Ri7W4joz3LUMQcjFj3qXd3AXonarXhwglVNB9BYquCdA5eq4wVUeAkm3R-e56TK5IZwpb5wNJDO3PhuXHSMwv1k-NEAIfI9_w&scenario=web_password_login&csrf=c2c20e89ce4e314771dcda79994b2cd020b9c30fc25faccdc1ebef3351a5b36b").
			End()
		if strings.Contains(body, "GEO_BLOCKED") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		if strings.Contains(body, "LOGIN_FORBIDDEN") || strings.Contains(body, "LOGIN_BAD_REQUEST") ||
			strings.Contains(body, "Your login is invalid. Please refresh the page.") {
			return model.Result{Name: name, Status: model.StatusYes}
		}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get www.hulu.com failed with code: %d", resp.StatusCode)}
}
