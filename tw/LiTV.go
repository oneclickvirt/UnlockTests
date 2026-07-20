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

const defaultLiTVSignatureKey = "7f4a9c2e8b6d1f3a5e9c7b4d2f8a6e1c"

func LiTV(c *http.Client) model.Result {
	name := "LiTV"
	hostname := "litv.tv"
	if c == nil {
		return model.Result{Name: name}
	}

	deviceID, err := getLiTVDeviceID(c)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}

	rpcResult := liTVRPC(c, hostname, name, deviceID)
	if rpcResult.Status != model.StatusUnexpected {
		return rpcResult
	}

	if legacyResult := liTVLegacy(c, hostname, name, deviceID); legacyResult.Status != model.StatusUnexpected {
		return legacyResult
	}
	return rpcResult
}

func liTVRPC(c *http.Client, hostname, name, deviceID string) model.Result {
	payload := fmt.Sprintf(
		`{"jsonrpc":"2.0","id":0,"method":"CCCService.GetProgramInformation","params":{"version":"2.0","project_num":"LTWEB02","device_id":"%s","swver":"LTWEB0210000WEB20190612185813","content_id":"VOD00328856","content_type":"drama"}}`,
		deviceID,
	)
	headers := map[string]string{
		"Origin":  "https://www.litv.tv",
		"Referer": "https://www.litv.tv/drama/watch/VOD00328856",
	}
	resp, body, err := utils.PostJson(c, "https://proxy.svc.litv.tv/cdi/v2/rpc", payload, headers)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	var res struct {
		Result *struct {
			Data *struct {
				ContentID string `json:"content_id"`
			} `json:"data"`
		} `json:"result"`
		Error *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return model.Result{Name: name, Status: model.StatusUnexpected, Err: err}
	}
	if res.Error != nil {
		return model.Result{Name: name, Status: model.StatusNo}
	}
	if res.Result != nil && res.Result.Data != nil && res.Result.Data.ContentID != "" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusNo}
}

func liTVLegacy(c *http.Client, hostname, name, deviceID string) model.Result {
	puid, err := getLiTVPUID(c)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}

	assetID := "vod70810-000001M001_1500K"
	mediaType := "vod"
	t := time.Now()
	timestamp := t.UnixMilli()
	nonce := genLiTVNonce(t)
	signature := genLiTVSignature(assetID, mediaType, nonce, timestamp)
	payload := map[string]interface{}{
		"AssetId":   assetID,
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
	resp, body, err := utils.PostJson(c, "https://www.litv.tv/api/get-urls-no-auth", string(jsonBytes), headers)
	if err != nil {
		return utils.HandleNetworkError(c, hostname, err, name)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		if strings.Contains(body, "OutsideRegionError") {
			return model.Result{Name: name, Status: model.StatusNo}
		}
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("legacy LiTV check failed with code: %d", resp.StatusCode)}
}

func getLiTVDeviceID(c *http.Client) (string, error) {
	headers := map[string]string{
		"Origin":  "https://www.litv.tv",
		"Referer": "https://www.litv.tv/",
	}
	resp, body, err := utils.PostJson(c, "https://www.litv.tv/api/generate-device-id", "", headers)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		DeviceID       string `json:"deviceId"`
		DeviceIDLegacy string `json:"device-id"`
	}
	var parseErr error
	if body != "" {
		parseErr = json.Unmarshal([]byte(body), &res)
		if parseErr == nil {
			if res.DeviceID != "" {
				return res.DeviceID, nil
			}
			if res.DeviceIDLegacy != "" {
				return res.DeviceIDLegacy, nil
			}
		}
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "device-id" {
			return cookie.Value, nil
		}
	}
	if parseErr != nil {
		return "", parseErr
	}
	return "", fmt.Errorf("deviceId not found in response")
}

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

func genLiTVNonce(t time.Time) string {
	return genLiTVBase36(13) + genLiTVBase36(13) + strconv.FormatInt(t.UnixMilli(), 36)
}

func genLiTVSignature(assetID, mediaType, nonce string, timestamp int64) string {
	key := defaultLiTVSignatureKey
	data := assetID + mediaType + strconv.FormatInt(timestamp, 10) + nonce + key
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func genLiTVBase36(length int) string {
	const base36Chars = "0123456789abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := range result {
		result[i] = base36Chars[rand.Intn(len(base36Chars))]
	}
	return string(result)
}
