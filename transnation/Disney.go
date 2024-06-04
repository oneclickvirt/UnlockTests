package transnation

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"net/http"
	"net/url"
	"strings"
)

// DisneyPlus
// www.disneyplus.com 双栈 且 post 请求
func DisneyPlus(c *http.Client) model.Result {
	name := "Disney+"
	if c == nil {
		return model.Result{Name: name}
	}
	QueryResult := requestDisney(c, "https://www.disneyplus.com", "query")
	if QueryResult.Status != model.StatusYes {
		return QueryResult
	}
	VerifyResult := requestDisney(c, "https://disney.api.edge.bamgrid.com/token", "auth")
	if VerifyResult.Status != model.StatusYes {
		return VerifyResult
	}
	return QueryResult
}

// 实际进行请求
func requestDisney(c *http.Client, URL string, method string) model.Result {
	name := "Disney+"
	if c == nil {
		return model.Result{Name: name}
	}
	data := url.Values{
		"grant_type":         {"refresh_token"},
		"refresh_token":      {"eyJ6aXAiOiJERUYiLCJraWQiOiJLcTYtNW1Ia3BxOXdzLUtsSUUyaGJHYkRIZFduRjU3UjZHY1h6aFlvZi04IiwiY3R5IjoiSldUIiwiZW5jIjoiQzIwUCIsImFsZyI6ImRpciJ9..OdwL8TEIFZouLDJe.wLz6zEC3PlPAGxx4X4qyP837lUbFrI_DQGnrJDMtEaQd5gsjHwaYshscoDXCYjMioU8JvsH_HKZga3fzSDEoWuMA5lgv4dyJpoB4Cqi91JjPSkqsRHKZ1I-nRoTmnSkcW3RHE-0coAqDWgK7IZ5cPiHQ-9KVRqqZkmTbEHynBdgH2y-FJP8zK0-dAynzR2krlUahhcykp7J7VqhZj_l5HVZZkPylZ6eKoK4J8fQvuGJoqMaRZTzrIH4Yk9J3GMbKnYqEG3SKRp5qAuWTtqLDOoGN0wWsUE5VRuCZxRKpxayJWABq2u4ABkAtIqUx8CPx77ZXxZVlcjRN1Xa8F2-e2mTxZq_1FgzmWECFg6onkDj_TpfBdeFoxDzhnRNceoQ-iyyNf3sgxJ_nz_bwztVZf0Vt3OR8yBnXfbkuEY7GQ4pvCuy-peW0mwJJCd2eJ9ADwDEGmoY4F47W-8rxdBhgna-0hu0FuLxt9MlmH_tGCmM_T-61xsxymLO9tlkwBnxNw4u6T9X2hcvC7-4uzr5cJiaJ3sGPMNo_ixTrP8SG9zCIse-X6_Lq0v3Uo-QOKhcD4N3gIfwZFYEvf-HVGWzFpU683q9CJfTTEXhsufj1URhSis7GdAa3nLZVt7CScsMPcYrMI317PmU-Brdvl_Ic4QeHTeF8-57kzD3mm5mrlQ7kQIXQzzQPqHYt70MzxL_scfT90cpYaSOBQnB1l--226h7X51XxSbrOcO-25zS7OSyedya8eMG6zAmgkk1zvZUzdCHZyzYD8-t0KYcfA5AwiLIFHxgqL4ni9fVy-SpYTKRwCmkp_pZOPaFwJh8zkhw8QaSLHq7ubko7H1kjJZxzsG1l4Bla1QRlj_-FVoY8GZ6okFk3Ts6A2qOK6v8UT7sL_w2zaHDQH1q2o05vsLwqIOxg3Xyey0tahzPbl-In_i1JGGvqGXOiPcKL5uOcTOo1luk32AbCS9i5mkopTS401YYYMH-Sx_krW_VJd2czpFefc0dlagtzBytqlcyscscFwq6IE6VHwG2Ij-WfO44G5hGDJFkZMZLeDUnTIyNrLe9hcfJp73koOSFnURsFWFjM2lgUIayiREAl02oh2alUyqnG09gdXufT_2W0DjA4i7qYuv6ol5NIVc389dF3x4a_7dPBvsMU3ppA1rlV04FlK6_fRv-Dk_jclXRZiQ5ul2ZO2CQ96LmrzmkdeNxFxcwaNXCJGBiRWXfMunoddIRg_LrVGuqWRgxj4DEnngZ2-qI_dliGiYraIehsHvtWeXIUWNF_FQSnQgZLg4WPekcluCecE4Iv7Sz36k9GUDqqs8hRWddirhufYem6RC84PyNqafCnwczrx5pOacmVzDl9Oi8OIhdDasdJa7gvsDoFzf6bv5st7EvbORkgPs6MK46mDMlwkL7TqjrJnSJzozCX4zLbYeyiWK6EXCehOpImMN262KLYQxnf5ugvk11gIA4NXpTbzyo4hp2LS7u8UMs5_w3t02vizxSQGokp-3qkEWmViy3pup1IXMPrcpS6KWHX0AYi1oRDZB5B8vM04pRHwYjsgMp2L-w4PMaDC4QDRU81IdvQ5VRkyLT7CL5hDlq5smXw_7wSFTWxs9vc5PmnrykSAkwFPocORC2j4T96uiu3z4gNoBu_dwKNcPi-dV7myC4iRRTmpm0V5A9IW510RGTyso_b-1hUeGvToYl9VwNgN7Impt3PjEQO2HXMU3p96tdulDEA_8bbyPdEGfxxVK3k2n_dxj_GzPKA8V4ESoNMRrV1vCuxPnrzfAOhqmNOEewTHqlxENSsZFGvfzVj1KemR7zLky14JMVslILnvxl6vuX7SbfIQ5JDktq9qKtTKo1mFrA-mBS3n00FacjPi364nnugiWQN7EwhNdEDH_KtWXGZVh-u2NM5cdoS1kAsOKSLxFTnTDG738LhoB3i_ZOjHFASKiZcsX6yD5csIP21jG5nFF9Qw2qsnqmxRuDLilIoGczEMt2Pfo180CG8Dyr7XtOYNeVU7__h9zBm9CvaAHDoQQhU4KlXM4LsljFeajw5f2wn08OmsdfkSYYl45O718QgzR_RRqwDpQH2pyKDJZ9yZt5OCyxcbnCgepjUyp6S-Pigfw73ASoCknhLLheb2mqkWIC-s3NmClpMoK-IyE57AiHHCatZfPGPnNofVioN5SbVR08mV7pdyQEhQGxGFM_LTAFFpwC48gOFTq-FWdV58muDULTqO3ImbGG6X3vV-PVbher1oJx0CFnelGGIx9lwM-yHbpVZGq9IXnKqoblCHiwuaJgbCKBnTjia2gYPNlN0Ql1ia3vQc7bybDVHyLePAVbOk10MdwHprwMGE__wsXqagElQCGJpU3ytPDktncRPCSQBQ3mw94CCIOQYEyhnA1Vik127AznwbR10Xm59diGBtix0Ao-VIrjKzQNw2hXqC_H-IgY46OT5ZndZ02SAe6AVyipq6kTui_ZyuQhy-zAOiat4t6qh-LyL1xImBuOZ7e79737LYiLHEIgHOIQ68DKcSmsIuA.gwrRhM5AiYUQ6iAbRZhxlw"},
		"subject_token_type": {"urn:bamtech:params:oauth:token-type:device"},
	}
	headers := map[string]string{
		"User-Agent":    model.UA_Browser,
		"authorization": "Bearer ZGlzbmV5JmJyb3dzZXImMS4wLjA.Cu56AgSfBTDag5NiRA81oLHkDZfu5L3CKadnefEAY84",
	}
	request := utils.Gorequest(c)
	request = utils.SetGoRequestHeaders(request, headers)
	resp, body, errs := request.Post(URL).
		Type("form").
		SendString(data.Encode()).
		End()
	if len(errs) > 0 {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: errs[0]}
	}
	switch method {
	case "auth":
		if strings.Contains(body, "unauthorized") {
			return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("unauthorized")}
		}
		if strings.Contains(body, "403 ERROR") || strings.Contains(body, "forbidden-location") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		return model.Result{Name: name, Status: model.StatusYes}
	case "query":
		if location := resp.Header.Get("Location"); location == "" {
			for _, c := range resp.Request.Cookies() {
				if c.Name == "x-dss-country" {
					return model.Result{
						Name: name, Status: model.StatusYes,
						Region: strings.ToLower(c.Value),
					}
				}
			}
			return model.Result{
				Name: name, Status: model.StatusYes,
			}
		} else if location == "https://disneyplus.disney.co.jp/" {
			return model.Result{
				Name: name, Status: model.StatusYes,
				Region: "jp",
			}
		} else if location == "https://preview.disneyplus.com/unavailable/" {
			return model.Result{
				Name: name, Status: model.StatusNo, Info: "unavailable",
			}
		}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}
