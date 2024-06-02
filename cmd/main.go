package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/oneclickvirt/UnlockTests/defaultset"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/transnation"
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

func FormarPrint(language string) {
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
	for _, r := range R {
		if r.Status == "" && r.Name != "" {
			s := ""
			check := false
			realLength := Length - len(r.Name) + 6
			if realLength %2 != 0 {
				realLength += 1
			}
			for i := realLength; i >= 0; i-- {
				s += "="
				if i < (realLength/2) && !check {
					s += " [ " + r.Name + " ] "
					check = true
				}
			}
			if r.Name == "" {
				s = "\n"
			}
			fmt.Println(s)
		} else {
			result := ShowResult(r)
			if r.Status == model.StatusYes && strings.HasSuffix(r.Name, "CDN") {
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

func Multination(ifaceName, ipAddr, netType string) {
	R = append(R, &model.Result{Name: "Multination"})
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
	for _, f := range FuncList {
		excute(f)
	}
}

func main() {
	wg = &sync.WaitGroup{}
	bar = NewBar(0)
	Multination("", "", "tcp4")
	bar.ChangeMax64(total)
	wg.Wait()
	bar.Finish()
	fmt.Println()
	FormarPrint("zh")
	fmt.Println()
}
