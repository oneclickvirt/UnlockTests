package vn

import (
	"fmt"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	"log"
	"net/http"
	"strings"
)

// KPLUS
// xem.kplus.vn 仅 ipv4 且 post 请求 有问题
// ssoToken 已过期
func KPLUS(c *http.Client) model.Result {
	name := "K+"
	hostname := "kplus.vn"
	if c == nil {
		return model.Result{Name: name}
	}
	ssoToken := "eyJrZXkiOiJ2c3R2IiwiZW5jIjoiQTEyOENCQy1IUzI1NiIsImFsZyI6ImRpciJ9..MWbBlLuci2KNLl9lvMe63g.IbBX7-dg3BWaXzzoxTQz-pJFulm_Y8axWLuG5DcJxQ9jTUPOhA2e6dzOP2hryAFVPFoIRs97ONGTHEYTFQgUtRlvqvx53jyTi3yegU6zWhJnhYZA2sdaj9khsNvVAth0zcWFoWA9GGwfNE5TZLOwczAexIxqC1Ee-tQDILC4XklFrJfvdzoCQBABRXpD_O4HHHIYFs0jBMtYSyD9Vq7dTD61sAVca_83lav7jvpP17PuAo3HHIFQtUdcugpgkB91mJbABIDTPdo0mqdzbgTA_FilwO1Z5qnpwqIZIXy0bhVXFFcwUZPIUxjLEVzP3SyHceFF5N-v7OeYhYZRLYuBKxWj1cRb3LAa3FGJvefqRsBadlsr0cZnOgx0TsL51a2SaIpNyyGtaq8KTTLULIZBb2Zsq2jmBkZtxjoPxUR8ku7J4sL0tfLDoMlWVZkrX4_1tls3E-l8Ael-wd0kbS1i2vpf-Vdh80lRClpDg3ibSSUFPsp3wYMFsuKfyY8vpHrCfYDJDDbYOSv20sfnU7q7gcmizTCFBuiszmXbFX9_aH8UOaCGeqkYDV1ZZ3mQ26TM7JEquuZTV09wdi81ABoM8RZcb2ua0cuocaO4-asMh8KQWNea9BCYlKK5NSPz--oGgGxSdvxZ63qQz1Lr4QZytA2buoQV5OlMoEP7k87fPcig5rPqsK7aeWUXJSmfiOBbSLztoiamvvHClMpds3frv0ud8NWUUoijmS_JUGfF7XYNxWWqEGJuDUoSllV5MVwtIb5wM069gR7zknrr5aRVDi3Nho16KHQ_iB3vxoIr-ExajWLNlvo44CopGhxhgOAKPkULV356uamZpB7twY_iEVrwGMQA1_hEH4usO-UbzuxL_pssLhJKD4NjVcTe86Z08Bfm0IyiNWESmFkA6FVfsxu57Yfd4bXT8mxnfXXmklb7u7vB0RVYRo4i26QGJbPknybHdfgQWEvRCMoAjEG-E2LymBAMwFneWEpPTwBMpfvlTHnGnUtfViA4Zy1xqF2q95g9AF9nF3sE4YpYuSFSkUQB4sZd8emDApIdP6Avqsq809Gg06_R2sUGrD9SQ-XbXhvtAYMcaUcSv54hJvRcSUkygqU8tdg4tJHR23UBb-I.UfpC5BKhvt8EE5gpIFMQoQ"
	firstRequestData := "{\"osVersion\":\"Windows NT 10.0\",\"appVersion\":\"114.0.0.0\",\"deviceModel\":\"Chrome\",\"deviceType\":\"PC\",\"deviceSerial\":\"w39db81c0-a2e9-11ed-952a-49b91c9e6f09\",\"deviceOem\":\"Chrome\",\"devicePrettyName\":\"Chrome\",\"ssoToken\":\"" +
		ssoToken + "\",\"brand\":\"vstv\",\"environment\":\"p\",\"language\":\"en_US\",\"memberId\":\"0\",\"featureLevel\":4,\"provisionData\":\"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYyI6dHJ1ZSwiaWF0IjoxNjg2NTc4NzYyLCJ1cCI6ImNwaSIsImRlIjoiYnJhbmRNYXBwaW5nIiwiYnIiOiJ2c3R2IiwiZHMiOiJ3MzlkYjgxYzAtYTJlOS0xMWVkLTk1MmEtNDliOTFjOWU2ZjA5In0.3mbI7wnJKtRf3493yc_ZEMEvzUXldwDx0sSZdwQnlNk\"}"
	headers := map[string]string{
		"Origin":  "https://xem.kplus.vn",
		"Referer": "https://xem.kplus.vn/",
	}
	url := "https://tvapi-sgn.solocoo.tv/v1/session"
	resp, body, err := utils.PostJson(c, url, firstRequestData, headers)
	if err != nil {
		log.Fatalf("Failed to make first request: %v", err.Error())
	}
	defer resp.Body.Close()
	//fmt.Println(body)
	token := ""
	if strings.Contains(body, "\"token\"") {
		token = strings.Split(strings.Split(body, "\"token\":\"")[1], "\"")[0]
	}
	if token == "" {
		return model.Result{Name: name, Status: model.StatusErr, Err: fmt.Errorf("failed to extract token from response")}
	}

	secondRequestData := `{"player":{"name":"RxPlayer","version":"3.29.0","capabilities":{"mediaTypes":["DASH","DASH"],"drmSystems":["PlayReady","Widevine"],"smartLib":true}}}`
	url2 := "https://tvapi-sgn.solocoo.tv/v1/assets/BJO0h8jMwJWg5Id_4VLxIJ-VscUzRry_myp4aC21/play"
	headers2 := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
	secondResp, secondBody, secondErr := utils.PostJson(c, url2, secondRequestData, headers2)
	if secondErr != nil {
		log.Fatalf("Failed to make second request: %v", secondErr.Error())
	}
	defer secondResp.Body.Close()
	//fmt.Println(secondBody)
	if strings.Contains(secondBody, "geoblock") {
		return model.Result{Name: name, Status: model.StatusNo}
	} else if secondBody != "" {
		result1, result2, result3 := utils.CheckDNS(hostname)
		unlockType := utils.GetUnlockType(result1, result2, result3)
		return model.Result{Name: name, Status: model.StatusYes, UnlockType: unlockType}
	}
	return model.Result{Name: name, Status: model.StatusUnexpected,
		Err: fmt.Errorf("get xem.kplus.vn failed with code: %d", resp.StatusCode)}
}
