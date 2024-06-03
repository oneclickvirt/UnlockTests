package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/oneclickvirt/UnlockTests/africa"
	"github.com/oneclickvirt/UnlockTests/asia"
	"github.com/oneclickvirt/UnlockTests/au"
	"github.com/oneclickvirt/UnlockTests/ca"
	"github.com/oneclickvirt/UnlockTests/ch"
	"github.com/oneclickvirt/UnlockTests/de"
	. "github.com/oneclickvirt/UnlockTests/defaultset"
	"github.com/oneclickvirt/UnlockTests/es"
	"github.com/oneclickvirt/UnlockTests/eu"
	"github.com/oneclickvirt/UnlockTests/fr"
	"github.com/oneclickvirt/UnlockTests/hk"
	"github.com/oneclickvirt/UnlockTests/in"
	"github.com/oneclickvirt/UnlockTests/it"
	"github.com/oneclickvirt/UnlockTests/jp"
	"github.com/oneclickvirt/UnlockTests/kr"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/nl"
	"github.com/oneclickvirt/UnlockTests/nz"
	"github.com/oneclickvirt/UnlockTests/ru"
	"github.com/oneclickvirt/UnlockTests/sg"
	"github.com/oneclickvirt/UnlockTests/th"
	"github.com/oneclickvirt/UnlockTests/transnation"
	"github.com/oneclickvirt/UnlockTests/tw"
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

func printCenteredMessage(message string, totalLength int) string {
	if totalLength == 0 {
		totalLength = 40
	}
	messageLength := len(message)
	if messageLength > totalLength {
		message = message[:totalLength]
		messageLength = totalLength
	}
	paddingLength := (totalLength - messageLength) / 2
	leftPadding := strings.Repeat("=", paddingLength)
	rightPadding := strings.Repeat("=", totalLength-messageLength-paddingLength)
	return (leftPadding + message + rightPadding + "\n")
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
	head := printCenteredMessage("[ "+message+" ]", 0)
	// 构建一个以 r.Name 为键的字典
	resultMap := make(map[string]*model.Result)
	for _, r := range R {
		resultMap[r.Name] = r
	}
	// 根据 Names 中的 name 顺序输出结果，重新排序结果
	tempList := []string{head}
	for _, name := range Names {
		if r, found := resultMap[name]; found {
			result := ShowResult(r)
			if r.Status == "Yes" && strings.HasSuffix(r.Name, "CDN") {
				result = Blue(r.Region)
			}
			tempList = append(tempList, fmt.Sprintf("%-"+strconv.Itoa(Length)+"s %s\n", r.Name, result))
		}
	}
	// 插入小分区的head行
	for _, r := range R {
		if r.Status == model.PrintHead {
			anotherList := []string{}
			for _, i := range tempList {
				if strings.Contains(i, r.Info) {
					tpHead := printCenteredMessage("[ "+r.Name+" ]", 20)
					anotherList = append(anotherList, tpHead)
				}
				anotherList = append(anotherList, i)
			}
			tempList = anotherList
		}
	}
	// 打印整体文本
	for _, i := range tempList {
		fmt.Printf(i)
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
		tp := f(nil)
		if tp.Status != model.PrintHead {
			Names = append(Names, tp.Name)
		}
	}
	// 实际开始任务
	for _, f := range FuncList {
		excute(f)
	}
}

func NorthAmerica(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		us.Fox,
		us.Hulu,
		us.NFLPlus,
		us.ESPNPlus,
		us.Epix,
		us.Starz,
		us.Philo,
		us.FXNOW,
		us.HBOMax,
		asia.TLCGO,
		us.Shudder,
		uk.BritBox,
		// us.Crackle,
		us.AETV,
		us.CWTV,
		us.NBATV,
		us.FuboTV,
		us.TubiTV,
		// us.NBCTV,
		us.SlingTV,
		us.PlutoTV,
		us.AcornTV,
		us.SHOWTIME,
		us.EncoreTVB,
		us.DiscoveryPlus,
		us.PeacockTV,
		us.Popcornflix,
		us.Crunchyroll,
		us.DirectvStream,
		// CA
		utils.PrintCA,
		asia.Hotstar,
		ca.CBCGem,
		ca.Crave,
	}
	processFunction(FuncList)
}

func Europe(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		eu.RakutenTV,
		eu.SkyShowTime,
		us.HBOMax,
		eu.SetantaSports,
		// GB
		utils.PrintGB,
		asia.Hotstar,
		uk.SkyGo,
		uk.BritBox,
		uk.ITVX,
		uk.Channel4,
		uk.Channel5,
		uk.BBCiPlayer,
		uk.DiscoveryPlus,
		// FR
		utils.PrintFR,
		fr.CanalPlus,
		fr.Molotov,
		// DE
		utils.PrintDE,
		de.Joyn,
		de.SkyDe,
		de.ZDF,
		// NL
		utils.PrintNL,
		nl.NLZIET,
		nl.VideoLand,
		nl.NPOStartPlus,
		// ES
		utils.PrintES,
		es.MoviStarPlus,
		// IT
		utils.PrintIT,
		it.RaiPlay,
		// ch
		utils.PrintCH,
		ch.SkyCh,
		// ru
		utils.PrintRU,
		ru.Amediateka,
	}
	processFunction(FuncList)
}

func HongKong(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		hk.NowE,
		hk.ViuTV,
		hk.MyTvSuper,
		asia.HBOGO,
		hk.BilibiliHKMO,
	}
	processFunction(FuncList)
}

func Africa(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		africa.DSTV,
		africa.Showmax,
		africa.BeinConnect,
	}
	processFunction(FuncList)
}

func India(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		asia.Hotstar,
		in.Zee5,
		in.JioCinema,
		in.MXPlayer,
		us.NBATV,
	}
	processFunction(FuncList)
}

func Taiwan(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		tw.KKTV,
		tw.LiTV,
		tw.MyVideo,
		tw.Tw4gtv,
		tw.LineTV,
		tw.HamiVideo,
		// tw.Catchplay,
		// tw.BahamutAnime,
		asia.HBOGO,
		tw.BilibiliTW,
	}
	processFunction(FuncList)
}

func Japan(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		jp.DMM,
		jp.DMMTV,
		jp.Abema,
		jp.Niconico,
		jp.Telasa,
		jp.UNext,
		jp.Hulu,
		jp.TVer,
		jp.Lemino,
		jp.Wowow,
		jp.VideoMarket,
		jp.DAnimeStore,
		jp.FOD,
		jp.Radiko,
		jp.RakutenTV,
		jp.J_COM_ON_DEMAND,
		// Game
		utils.PrintGame,
		jp.Kancolle,
		jp.PrettyDerby,
		jp.KonosubaFD,
		jp.PCRJP,
		jp.WorldFlipper,
		jp.ProjectSekai,
		// Music
		utils.PrintMusic,
		jp.Mora,
		jp.MusicBook,
		jp.KaraokeDam,
		// Forum
		utils.PrintForum,
		jp.EroGameSpace,
	}
	processFunction(FuncList)
}

func Multination(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		transnation.DAZN,
		transnation.DisneyPlus,
		transnation.Netflix,
		transnation.NetflixCDN,
		transnation.Youtube,
		transnation.YoutubeCDN,
		transnation.PrimeVideo,
		transnation.ParamountPlus,
		transnation.TVBAnywhere,
		transnation.IQiYi,
		transnation.ViuCom,
		transnation.Spotify,
		transnation.Steam,
		transnation.OpenAI,
		transnation.WikipediaEditable,
		transnation.Reddit,
		transnation.TikTok,
		transnation.Bing,
		transnation.Instagram,
		transnation.KOCOWA,
		transnation.SonyLiv,
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
		// AU
		utils.PrintAU,
		au.Stan,
		au.Binge,
		au.Au7plus,
		au.Channel9,
		au.Channel10,
		au.ABCiView,
		au.OptusSports,
		au.SBSonDemand,
		eu.Docplay,
		au.KayoSports,
		// NZ
		utils.PrintNZ,
		nz.NeonTV,
		nz.SkyGO,
		nz.ThreeNow,
		nz.MaoriTV,
	}
	processFunction(FuncList)
}

func Korean(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		kr.Wavve,
		kr.Tving,
		kr.Watcha,
		kr.CoupangPlay,
		kr.SPOTVNOW,
		kr.NaverTV,
		kr.Afreeca,
		kr.KBSDomestic,
	}
	processFunction(FuncList)
}

func SouthEastAsia(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		asia.Hotstar,
		asia.HBOGO,
		asia.BilibiliSEA,
		// SG
		utils.PrintSG,
		sg.MeWatch,
		// TH
		utils.PrintTH,
		th.AISPlay,
		th.TrueID,
		// th.BilibiliTH, 失效 - 不做检测
		// ID 全失效 - 不做检测
		// VN 全失效 - 不做检测
	}
	processFunction(FuncList)
}

func Sport(ifaceName, ipAddr, netType string) {
	var FuncList = [](func(request *gorequest.SuperAgent) model.Result){
		transnation.DAZN,
		asia.StarPlus,
		us.ESPNPlus,
		us.NBATV,
		us.FuboTV,
		asia.MolaTV,
		eu.SetantaSports,
		au.OptusSports,
		africa.BeinConnect,
		// eu.Eurosport,
	}
	processFunction(FuncList)
}

func main() {
	wg = &sync.WaitGroup{}
	bar = NewBar(0)
	// NorthAmerica("", "", "tcp4")
	// Europe("", "", "tcp4")
	// HongKong("", "", "tcp4")
	// Africa("", "", "tcp4")
	// India("", "", "tcp4")
	// Taiwan("", "", "tcp4")
	// Japan("", "", "tcp4")
	// Multination("", "", "tcp4")
	// SouthAmerica("", "", "tcp4")
	// Oceania("", "", "tcp4")
	// Korean("", "", "tcp4")
	// SouthEastAsia("", "", "tcp4")
	Sport("", "", "tcp4")
	bar.ChangeMax64(total)
	wg.Wait()
	bar.Finish()
	fmt.Println()
	// FormarPrint("zh", "North America")
	// FormarPrint("zh", "Europe")
	// FormarPrint("zh", "HongKong")
	// FormarPrint("zh", "Africa")
	// FormarPrint("zh", "India")
	// FormarPrint("zh", "Taiwan")
	// FormarPrint("zh", "Japan")
	// FormarPrint("zh", "Multination")
	// FormarPrint("zh", "South America")
	// FormarPrint("zh", "Oceania")
	// FormarPrint("zh", "Korean")
	// FormarPrint("zh", "South East Asia")
	FormarPrint("zh", "Sport")
	fmt.Println()
}
