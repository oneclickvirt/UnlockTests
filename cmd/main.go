package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oneclickvirt/UnlockTests/asia"
	. "github.com/oneclickvirt/UnlockTests/defaultset"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/transnation"
	"github.com/oneclickvirt/UnlockTests/uk"
	"github.com/oneclickvirt/UnlockTests/us"
	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/parnurzeal/gorequest"
	pb "github.com/schollz/progressbar/v3"
)

var (
	total                      int64
	bar                        *pb.ProgressBar
	wg                         *sync.WaitGroup
	IPV4                       = true
	IPV6                       = true
	R                          []*model.Result
	Names                      []string
	ifaceName, ipAddr, netType string
)

func NewBar(count int64) *pb.ProgressBar {
	return pb.NewOptions64(
		count,
		pb.OptionSetDescription("testing"),
		pb.OptionSetWriter(os.Stderr),
		pb.OptionSetWidth(20),
		pb.OptionThrottle(100*time.Millisecond),
		pb.OptionShowCount(),
		pb.OptionClearOnFinish(),
		pb.OptionEnableColorCodes(true),
		pb.OptionSpinnerType(14),
	)
}

func ShowResult(r *model.Result) (s string) {
	formatResult := func(colorFunc func(string) string, status string, r model.Result) string {
		s := colorFunc(status)
		if r.Info != "" {
			s += colorFunc(" (" + r.Info + ")")
		}
		if r.Region != "" {
			s += colorFunc(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		return s
	}
	switch r.Status {
	case model.StatusYes:
		return formatResult(Green, "YES", *r)
	case model.StatusNetworkErr:
		return Red("NO") + Yellow(" (Network Err)")
	case model.StatusRestricted:
		return formatResult(Yellow, "Restricted", *r)
	case model.StatusErr:
		s = Yellow("Error")
		if r.Err != nil {
			s += ": " + r.Err.Error()
		}
		return s
	case model.StatusNo:
		return formatResult(Red, "NO", *r)
	case model.StatusBanned:
		s = Red("Banned")
		if r.Info != "" {
			s += Yellow(" (" + r.Info + ")")
		}
		return s
	case model.StatusUnexpected:
		s = Purple("Unknown")
		if r.Err != nil {
			s += ": " + r.Err.Error()
		}
		return s
	default:
		return ""
	}
}

func printCenteredMessage(message string) {
	totalLength := 40
	messageLength := len(message)
	if messageLength > totalLength {
		message = message[:totalLength]
		messageLength = totalLength
	}
	paddingLength := (totalLength - messageLength) / 2
	leftPadding := strings.Repeat("=", paddingLength)
	rightPadding := strings.Repeat("=", totalLength-messageLength-paddingLength)
	fmt.Println(leftPadding + message + rightPadding)
}

func FormarPrint(language, message string) {
	if language == "zh" {
		fmt.Println("测试时间: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	} else {
		fmt.Println("Test time: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	}
	Length := 25
	for _, r := range R {
		if len(r.Name) > Length {
			Length = len(r.Name)
		}
	}
	printCenteredMessage("[ " + message + " ]")
	// 构建一个以 r.Name 为键的字典
	resultMap := make(map[string]*model.Result)
	for _, r := range R {
		resultMap[r.Name] = r
	}
	// 根据 Names 中的 name 顺序输出结果，重新排序结果
	for _, name := range Names {
		if r, found := resultMap[name]; found {
			result := ShowResult(r)
			if r.Status == "Yes" && strings.HasSuffix(r.Name, "CDN") {
				result = Blue(r.Region)
			}
			fmt.Printf("%-"+strconv.Itoa(Length)+"s %s\n", r.Name, result)
		}
	}
}

func excute(F func(request *gorequest.SuperAgent) model.Result) {
	wg.Add(1)
	total++
	go func() {
		defer wg.Done()
		req, err := utils.ParseInterface(ifaceName, ipAddr, netType)
		if err == nil {
			res := F(req)
			R = append(R, &res)
			bar.Describe(res.Name + " " + ShowResult(&res))
			bar.Add(1)
		} else {
			bar.Describe(err.Error())
			bar.Add(1)
		}
	}()
}

func processFunction(FuncList [](func(request *gorequest.SuperAgent) model.Result)) {
	// 生成顺序输出的名字
	for _, f := range FuncList {
		Names = append(Names, f(nil).Name)
	}
	// 实际开始任务
	for _, f := range FuncList {
		excute(f)
	}
}

func Multination(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		transnation.DAZN,
		transnation.DisneyPlus,
		transnation.Netflix,
		transnation.Youtube,
		transnation.PrimeVideo,
		transnation.TVBAnywhere,
		transnation.IQiYi,
		transnation.YoutubeCDN,
		transnation.NetflixCDN,
		transnation.Spotify,
		transnation.OpenAI,
		transnation.Bing,
		transnation.WikipediaEditable,
		transnation.Instagram,
		transnation.Steam,
		transnation.Reddit,
	}
	processFunction(FuncList)
}

func SouthAmerica(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		asia.StarPlus,
		us.HBOMax,
		us.DirecTVGO,
	}
	processFunction(FuncList)
}

func Oceania(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		us.NBATV,
		us.AcornTV,
		uk.BritBox,
		transnation.ParamountPlus,
		transnation.SonyLiv,
	}
	processFunction(FuncList)
}

func main() {
	wg = &sync.WaitGroup{}
	bar = NewBar(0)
	// Multination("", "", "tcp4")
	// SouthAmerica("", "", "tcp4")
	Oceania("", "", "tcp4")
	bar.ChangeMax64(total)
	wg.Wait()
	bar.Finish()
	fmt.Println()
	// FormarPrint("zh", "Multination")
	// FormarPrint("zh", "SouthAmerica")
	FormarPrint("zh", "Oceania")
	fmt.Println()
}
