package tw

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
)

// LiTV
// www.litv.tv 仅 ipv4 且 post 请求
func LiTV(c *http.Client) model.Result {
	name := "LiTV"
	hostname := "litv.tv"
	if c == nil {
		return model.Result{Name: name}
	}
	// 获取 device-id
	deviceID, err := getLiTVDeviceID(c)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	// 获取 PUID
	puid, err := getLiTVPUID(c)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	assetId := "vod70810-000001M001_1500K"
	mediaType := "vod"
	t := time.Now()
	timestamp := t.UnixMilli()
	nonce := genLiTVNonce(t)
	signature := genLiTVSignature(assetId, mediaType, nonce, timestamp)
	payload := map[string]interface{}{
		"AssetId":   assetId,
		"MediaType": mediaType,
		"puid":      puid,
		"timestamp": timestamp,
		"nonce":     nonce,
		"signature": signature,
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusErr, Err: err}
	}
	headers := map[string]string{
		"Cookie":  fmt.Sprintf("device-id=%s; PUID=%s", deviceID, puid),
		"Origin":  "https://www.litv.tv",
		"Referer": "https://www.litv.tv/drama/watch/VOD00328856",
	}
	resp, body, err := utils.PostJson(c, "https://www.litv.tv/api/get-urls-no-auth",
		string(jsonBytes),
		headers,
	)
	if err != nil {
		return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
	}
	defer resp.Body.Close()
	bodyString := string(body)
	if resp.StatusCode == 200 {
		if strings.Contains(bodyString, "OutsideRegionError") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected}
}

// getLiTVDeviceID 获取设备ID
func getLiTVDeviceID(c *http.Client) (string, error) {
	headers := map[string]string{
		"Origin":  "https://www.litv.tv",
		"Referer": "https://www.litv.tv/",
	}
	resp, _, err := utils.PostJson(c, "https://www.litv.tv/api/generate-device-id", "", headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "device-id" {
			return cookie.Value, nil
		}
	}
	return "", fmt.Errorf("device-id cookie not found")
}

// getLiTVPUID 获取PUID
func getLiTVPUID(c *http.Client) (string, error) {
	payload := `{"jsonrpc":"2.0","id":100,"method":"PustiService.PUID","params":{"version":"2.0","device_id":"","device_category":"LTWEB00","puid":"","aaid":"","idfa":""}}`
	headers := map[string]string{
		"Origin":  "https://www.litv.tv",
		"Referer": "https://www.litv.tv/",
	}
	resp, body, err := utils.PostJson(c, "https://pusti.svc.litv.tv/puid", payload, headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var res struct {
		Result struct {
			Puid string `json:"puid"`
		} `json:"result"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return "", err
	}
	if res.Result.Puid == "" {
		return "", fmt.Errorf("puid not found in response")
	}
	return res.Result.Puid, nil
}

// genLiTVNonce 生成nonce
func genLiTVNonce(t time.Time) string {
	return genBase36(13) + genBase36(13) + strconv.FormatInt(t.UnixMilli(), 36)
}

// genLiTVSignature 生成签名
func genLiTVSignature(assetId, mediaType, nonce string, timestamp int64) string {
	key := "7f4a9c2e8b6d1f3a5e9c7b4d2f8a6e1c"
	// e + t + r + n + i
	data := assetId + mediaType + strconv.FormatInt(timestamp, 10) + nonce + key
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// genBase36 生成指定长度的base36随机字符串
func genBase36(length int) string {
	const base36Chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = base36Chars[rand.Intn(len(base36Chars))]
	}
	return string(result)
}

// AnotherLiTV
// www.litv.tv 的另一个检测逻辑
// func AnotherLiTV(c *http.Client) model.Result {
// 	name := "LiTV"
// 	hostname := "litv.tv"
// 	if c == nil {
// 		return model.Result{Name: name}
// 	}
// 	url := "https://www.litv.tv/vod/ajax/getUrl"
// 	payload := `{"type":"noauth","assetId":"vod44868-010001M001_800K","puid":"6bc49a81-aad2-425c-8124-5b16e9e01337"}`
// 	headers := map[string]string{
// 		"Content-Type": "application/json",
// 	}
// 	resp, body, err := utils.PostJson(c, url, payload, headers)
// 	if err != nil {
// 		return utils.HandleNetworkError(c, hostname, err, name)
// 	}
// 	defer resp.Body.Close()
// 	var jsonResponse map[string]interface{}
// 	err = json.Unmarshal([]byte(body), &jsonResponse)
// 	if err != nil {
// 		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
// 	}
// 	errorMessage, ok := jsonResponse["errorMessage"].(string)
// 	if !ok {
// 		return model.Result{Name: name, Status: model.StatusUnexpected}
// 	}
// 	switch errorMessage {
// 	case "null":
// 		result1, result2, result3 := utils.CheckDNS(hostname)
// 		unlockType := utils.GetUnlockType(result1, result2, result3)
// 		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
// 	case "vod.error.outsideregionerror":
// 		return model.Result{Name: name, Status: model.StatusNo}
// 	default:
// 		return model.Result{Name: name, Status: model.StatusUnexpected}
// 	}
// }
