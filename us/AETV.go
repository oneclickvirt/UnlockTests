package us

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// extractAETVCountryCode extracts country code from AETN meta tag
func extractAETVCountryCode(html string) string {
	re := regexp.MustCompile(`<meta\s+name=["']aetn:countryCode["']\s+content=["']([A-Z]{2})["']\s*/?>`)
	matches := re.FindStringSubmatch(html)
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

// AETV
// 仅 ipv4 且 post 请求
func AETV(c *http.Client) model.Result {
	name := "A&E TV"
	hostname := "aetv.com"
	if c == nil {
		return model.Result{Name: name}
	}
	// 检测 Google DAI + Argus API
	step1Url := "https://dai.google.com/ondemand/hls/content/2540935/vid/2400478275868/streams"
	step1Data := `dai-dlid=default_delivery_60101&afid=59946479&adobe_id=17218268117322661452685981206066538688&asnw=171213&caid=305489&imafw__fw_ae=nonauthenticated&imafw__fw_vcid2=2395239744222918&devicename=Desktop&imafw_csid=aetv.desktop.video&imafw__fw_player_height=901&imafw__fw_player_width=474&imafw__fw_us_privacy=1---&imafw__fw_site_page=https%3A%2F%2Fplay.aetv.com%2Fshows%2Fozark-law%2Fseason-1%2Fepisode-1&ltd=1&imafw__fw_h_user_agent=Mozilla%2F5.0%20(Windows%20NT%2010.0%3B%20Win64%3B%20x64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F144.0.0.0%20Safari%2F537.36%20Edg%2F144.0.0.0&imafw_prof=171213%3Aaetn_desktop_ae_googlessai_vod&metr=1031&nw=171213&pvrn=7469863039993294&ssnw=171213&vprn=7469863039993294&flag=%2Bsltp%2Bvicb%2Bslcb%2Bamsl%2Bamcb%2Bssus%2Bemcr%2Bfbad%2Bdtrd%2Bplay&resp=vmap1&cld=1&ctv=0&correlator=3412237903465914&ptt=20&osd=2&sdr=1&sdki=41&sdkv=h.3.740.0&uach=WyJXaW5kb3dzIiwiMTUuMC4wIiwieDg2IiwiIiwiMTQ0LjAuMzcxOS45MiIsbnVsbCwwLG51bGwsIjY0IixbWyJOb3QoQTpCcmFuZCIsIjguMC4wLjAiXSxbIkNocm9taXVtIiwiMTQ0LjAuNzU1OS45NyJdLFsiTWljcm9zb2Z0IEVkZ2UiLCIxNDQuMC4zNzE5LjkyIl1dLDBd&ua=Mozilla%2F5.0%20(Windows%20NT%2010.0%3B%20Win64%3B%20x64)%20AppleWebKit%2F537.36%20(KHTML%2C%20like%20Gecko)%20Chrome%2F144.0.0.0%20Safari%2F537.36%20Edg%2F144.0.0.0&eid=44751890%2C95322027%2C95331589%2C95332046&frm=0&omid_p=Google1%2Fh.3.740.0&sdk_apis=7&sid=6C436384-F48F-4AEA-8FCA-9F3633757A74&ssss=gima&ref=https%3A%2F%2Fplay.aetv.com%2Fshows%2Fozark-law%2Fseason-1%2Fepisode-1&url=https%3A%2F%2Fplay.aetv.com%2Fshows%2Fozark-law%2Fseason-1%2Fepisode-1&wta=0&us_privacy=1---&gpp_sid=-1&eoidce=1`
	client1 := utils.Req(c)
	resp1, err1 := client1.R().
		SetHeader("accept", "*/*").
		SetHeader("accept-language", "zh-HK,zh;q=0.9").
		SetHeader("content-type", "application/x-www-form-urlencoded;charset=UTF-8").
		SetHeader("origin", "https://play.aetv.com").
		SetHeader("referer", "https://play.aetv.com/").
		SetBody(step1Data).
		Post(step1Url)
	if err1 == nil && resp1 != nil {
		defer resp1.Body.Close()
		if resp1.StatusCode == 200 || resp1.StatusCode == 201 {
			b1, err := io.ReadAll(resp1.Body)
			if err == nil {
				var res1 struct {
					StreamManifest string `json:"stream_manifest"`
				}
				if err := json.Unmarshal(b1, &res1); err == nil && res1.StreamManifest != "" {
					manifestUrl, err := url.Parse(res1.StreamManifest)
					if err == nil {
						originPath := manifestUrl.Query().Get("originpath")
						if originPath != "" {
							argusUrl := "https://argus.media.aetnd.com/"
							u2, _ := url.Parse(argusUrl)
							q2 := u2.Query()
							q2.Set("bestCdn", "fastly")
							q2.Set("brand", "aetv")
							q2.Set("dfpOp", originPath)
							q2.Set("client", "tve-web-theo")
							q2.Set("sig", "00697b29892c60831e3de250beaec91bd6ab73901a0ffb723b533130425058484d6c62")
							u2.RawQuery = q2.Encode()
							client2 := utils.Req(c)
							resp2, err2 := client2.R().
								SetHeader("accept", "*/*").
								SetHeader("accept-language", "zh-HK,zh;q=0.9").
								SetHeader("origin", "https://play.aetv.com").
								SetHeader("referer", "https://play.aetv.com/").
								SetHeader("sec-ch-ua", `"Not(A:Brand";v="8", "Chromium";v="144", "Microsoft Edge";v="144"`).
								SetHeader("sec-ch-ua-mobile", "?0").
								SetHeader("sec-ch-ua-platform", `"Windows"`).
								SetHeader("sec-fetch-dest", "empty").
								SetHeader("sec-fetch-mode", "cors").
								SetHeader("sec-fetch-site", "cross-site").
								SetHeader("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0").
								SetHeader("x-video-meta-token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3Njk2NzU2NDAsImV4cCI6MTc2OTY3OTI0MCwicHBsSWQiOiIzMDU0ODkiLCJpc0JlaGluZFdhbGwiOmZhbHNlLCJpc0xvbmdGb3JtIjp0cnVlLCJlbmNvZGVyIjoiYml0bW92aW5fdjEiLCJyZW5kaXRpb25zTGlzdCI6bnVsbCwicmVuZGl0aW9uc1BhdGhQcmVmaXgiOiJiaXRtb3Zpbi9BRVROLUFFVFZfVk1TL0FFTl9PWlJLXzMwNTQ4OV9HTEJfNDkzNjY5XzIzOThfNjBfMjAyNTAxMDRfMDFfQUVUTi1BRVRWX1ZNUyIsInJlZ2lvbnNBdmFpbGFibGUiOlsiVVMiLCJDQSIsIkFTIiwiR1UiLCJNUCIsIlBSIiwiVkkiLCJVTSJdfQ._IBIJ-Yh8X9hGOGglCNQWcW6WZKrUbjTNQ8Rdon8u_A").
								Get(u2.String())
							if err2 == nil && resp2 != nil {
								defer resp2.Body.Close()
								if resp2.StatusCode == 403 {
									b2, err := io.ReadAll(resp2.Body)
									if err == nil {
										var res2 struct {
											Exception string `json:"exception"`
										}
										if err := json.Unmarshal(b2, &res2); err == nil {
											if res2.Exception == "GeoLocationBlocked" {
												return model.Result{Name: name, Status: model.StatusNo}
											}
											if res2.Exception == "JWTExpiredSignature" {
												// JWT过期说明通过了地理位置检查
												result1, result2, result3 := utils.CheckDNS(hostname)
												unlockType := utils.GetUnlockType(result1, result2, result3)
												return model.Result{
													Name:       name,
													Status:     model.StatusYes,
													Region:     "us",
													UnlockType: unlockType,
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	// 检查主页 meta 标签
	url3 := "https://www.aetv.com/"
	client3 := utils.Req(c)
	resp3, err3 := client3.R().Get(url3)
	if err3 == nil {
		defer resp3.Body.Close()
		b3, err3 := io.ReadAll(resp3.Body)
		if err3 == nil {
			body3 := string(b3)
			region := extractAETVCountryCode(body3)
			switch region {
			case "US":
				result1, result2, result3 := utils.CheckDNS(hostname)
				unlockType := utils.GetUnlockType(result1, result2, result3)
				return model.Result{
					Name:       name,
					Status:     model.StatusYes,
					Region:     "us",
					UnlockType: unlockType,
				}
			case "":
				// 继续
			default:
				return model.Result{Name: name, Status: model.StatusNo}
			}
		}
	}
	// 检查隐私管理器地理位置
	url0 := "https://geo.privacymanager.io/"
	client0 := utils.Req(c)
	resp0, err0 := client0.R().Get(url0)
	if err0 == nil {
		defer resp0.Body.Close()
		b0, err0 := io.ReadAll(resp0.Body)
		if err0 == nil {
			var geoRes struct {
				Country string `json:"country"`
			}
			if err := json.Unmarshal(b0, &geoRes); err == nil {
				if geoRes.Country == "US" || geoRes.Country == "CA" {
					result1, result2, result3 := utils.CheckDNS(hostname)
					unlockType := utils.GetUnlockType(result1, result2, result3)
					return model.Result{
						Name:       name,
						Status:     model.StatusYes,
						Region:     strings.ToLower(geoRes.Country),
						UnlockType: unlockType,
					}
				}
			}
		}
	}
	// 检查 theplatform API (检测地理位置封锁)
	url1 := "https://link.theplatform.com/s/xc6n8B/UR27JDU0bu2s/"
	client1b := utils.Req(c)
	resp1b, err1b := client1b.R().Post(url1)
	if err1b == nil {
		defer resp1b.Body.Close()
		b1, err1b := io.ReadAll(resp1b.Body)
		if err1b == nil {
			body1 := string(b1)
			if strings.Contains(body1, "GeoLocationBlocked") {
				return model.Result{Name: name, Status: model.StatusNo}
			}
		}
	}
	// 检查 play.aetv.com cookie
	url2 := "https://play.aetv.com/"
	client2b := utils.Req(c)
	resp2b, err2b := client2b.R().Post(url2)
	if err2b == nil {
		defer resp2b.Body.Close()
		b2, err2b := io.ReadAll(resp2b.Body)
		if err2b == nil {
			body2 := string(b2)
			if body2 != "" {
				tp := utils.ReParse(body2, `AETN-Country-Code=([A-Z]+)`)
				if tp != "" {
					region := strings.ToLower(tp)
					if region == "ca" || region == "us" {
						result1, result2, result3 := utils.CheckDNS(hostname)
						unlockType := utils.GetUnlockType(result1, result2, result3)
						return model.Result{
							Name:       name,
							Status:     model.StatusYes,
							Region:     region,
							UnlockType: unlockType,
						}
					} else {
						return model.Result{Name: name, Status: model.StatusNo}
					}
				}
			}
		}
	}
	return model.Result{
		Name:   name,
		Status: model.StatusUnexpected,
	}
}