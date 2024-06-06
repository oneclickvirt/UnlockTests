package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
	pb "github.com/schollz/progressbar/v3"
)

var (
	total                                           int64
	bar                                             *pb.ProgressBar
	wg                                              *sync.WaitGroup
	IPV4, IPV6                                      = true, true
	R                                               []*model.Result
	Names                                           []string
	ifaceName, ipAddr, netType                      string
	M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT = false, false, false, false, false, false, false, false, false, false, false
	Version                                         = "0.0.2"
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
		// if r.Err != nil {
		// 	s += ": " + r.Err.Error()
		// }
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

func excute(F func(c *http.Client) model.Result, c *http.Client) {
	wg.Add(1)
	total++
	go func() {
		defer wg.Done()
		res := F(c)
		R = append(R, &res)
		bar.Describe(res.Name + " " + ShowResult(&res))
		bar.Add(1)
	}()
}

func preProcess(FuncList [](func(c *http.Client) model.Result)) {
	// 生成顺序输出的名字
	for _, f := range FuncList {
		tp := f(nil)
		if tp.Status != model.PrintHead {
			Names = append(Names, tp.Name)
		}
	}
}

func processFunction(FuncList [](func(c *http.Client) model.Result), c *http.Client) {
	// 实际开始任务
	for _, f := range FuncList {
		excute(f, c)
	}
}

func NorthAmerica() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		us.Fox,
		us.Hulu,
		us.NFLPlus,
		us.ESPNPlus,
		us.Epix,
		us.Starz,
		us.Philo,
		us.FXNOW,
		us.HBOMax,
		us.Shudder,
		uk.BritBox,
		// us.Crackle,
		us.AETV,
		us.NBCTV,
		us.CWTV,
		us.NBATV,
		us.FuboTV,
		us.TubiTV,
		// us.NBCTV,
		us.TLCGO,
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
		eu.MathsSpot,
		// CA
		utils.PrintCA,
		asia.HotStar,
		ca.CBCGem,
		ca.Crave,
	}
	preProcess(FuncList)
	return FuncList
}

func Europe() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		eu.RakutenTV,
		eu.SkyShowTime,
		us.HBOMax,
		eu.SetantaSports,
		eu.MathsSpot,
		// GB
		utils.PrintGB,
		asia.HotStar,
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
	preProcess(FuncList)
	return FuncList
}

func HongKong() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		hk.NowE,
		hk.ViuTV,
		hk.MyTvSuper,
		asia.HBOGO,
		hk.BilibiliHKMO,
		tw.BahamutAnime,
	}
	preProcess(FuncList)
	return FuncList
}

func Africa() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		africa.DSTV,
		africa.Showmax,
		africa.BeinConnect,
	}
	preProcess(FuncList)
	return FuncList
}

func India() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.HotStar,
		in.Zee5,
		in.JioCinema,
		in.MXPlayer,
		us.NBATV,
	}
	preProcess(FuncList)
	return FuncList
}

func Taiwan() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		tw.KKTV,
		tw.LiTV,
		tw.MyVideo,
		tw.Tw4gtv,
		tw.LineTV,
		tw.HamiVideo,
		tw.Catchplay,
		tw.BahamutAnime,
		asia.HBOGO,
		tw.BilibiliTW,
	}
	preProcess(FuncList)
	return FuncList
}

func Japan() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
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
		jp.NETRIDE,
		// Music
		utils.PrintMusic,
		jp.Mora,
		jp.MusicBook,
		jp.KaraokeDam,
		// Forum
		utils.PrintForum,
		jp.EroGameSpace,
	}
	preProcess(FuncList)
	return FuncList
}

func Multination() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
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
		transnation.OneTrust,
		transnation.GoogleSearch,
	}
	preProcess(FuncList)
	return FuncList
}

func SouthAmerica() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.StarPlus,
		us.HBOMax,
		us.DirecTVGO,
	}
	preProcess(FuncList)
	return FuncList
}

func Oceania() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
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
	preProcess(FuncList)
	return FuncList
}

func Korea() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		kr.Wavve,
		kr.Tving,
		kr.Watcha,
		kr.CoupangPlay,
		kr.SPOTVNOW,
		kr.NaverTV,
		kr.Afreeca,
		kr.KBSDomestic,
	}
	preProcess(FuncList)
	return FuncList
}

func SouthEastAsia() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.HotStar,
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
	preProcess(FuncList)
	return FuncList
}

func Sport() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		transnation.DAZN,
		asia.StarPlus,
		us.ESPNPlus,
		us.NBATV,
		us.NBCTV,
		us.FuboTV,
		asia.MolaTV,
		eu.SetantaSports,
		au.OptusSports,
		africa.BeinConnect,
		eu.Eurosport,
	}
	preProcess(FuncList)
	return FuncList
}

func IPV6Multination() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.HotStar,
		transnation.DisneyPlus,
		transnation.Netflix,
		transnation.NetflixCDN,
		transnation.Youtube,
		transnation.YoutubeCDN,
		transnation.WikipediaEditable,
		transnation.Bing,
	}
	preProcess(FuncList)
	return FuncList
}

func GetIpv4Info() {
	resp, err := utils.Req(utils.Ipv4HttpClient).R().Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		IPV4 = false
		fmt.Println("Can not detect IPv4 Address")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		IPV4 = false
		fmt.Println("Can not detect IPv4 Address")
		return
	}
	body := string(b)
	if body != "" && strings.Contains(body, "ip=") {
		s := body
		i := strings.Index(s, "ip=")
		s = s[i+3:]
		i = strings.Index(s, "\n")
		fmt.Println("Your IPV4 address:", Blue(s[:i]))
	}
}

func GetIpv6Info() {
	resp, err := utils.ReqDefault(utils.Ipv6HttpClient).R().Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		IPV6 = false
		fmt.Println("Can not detect IPv6 Address")
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		IPV4 = false
		fmt.Println("Can not detect IPv6 Address")
		return
	}
	body := string(b)
	fmt.Println(body)
	if body != "" && strings.Contains(body, "ip=") {
		s := body
		i := strings.Index(s, "ip=")
		s = s[i+3:]
		i = strings.Index(s, "\n")
		fmt.Println("Your IPV6 address:", Blue(s[:i]))
	}
}

func finallyPrintResult(language, netType string) {
	getPlatformName := func(multi bool, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT bool) string {
		if multi {
			if TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 台湾平台"
			} else if !TW && HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 香港平台"
			} else if !TW && !HK && JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 日本平台"
			} else if !TW && !HK && !JP && KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 韩国平台"
			} else if !TW && !HK && !JP && !KR && NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 北美平台"
			} else if !TW && !HK && !JP && !KR && !NA && SA && !EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 南美平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && EU && !AFR && !OCEA && !SPORT {
				return "跨国平台 + 欧洲平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && !EU && AFR && !OCEA && !SPORT {
				return "跨国平台 + 非洲平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && OCEA && !SPORT {
				return "跨国平台 + 大洋洲平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && SPORT {
				return "跨国平台 + 体育平台"
			} else {
				return "跨国平台"
			}
		} else {
			if TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "台湾平台"
			} else if !TW && HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "香港平台"
			} else if !TW && !HK && JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "日本平台"
			} else if !TW && !HK && !JP && KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "韩国平台"
			} else if !TW && !HK && !JP && !KR && NA && !SA && !EU && !AFR && !OCEA && !SPORT {
				return "北美平台"
			} else if !TW && !HK && !JP && !KR && !NA && SA && !EU && !AFR && !OCEA && !SPORT {
				return "南美平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && EU && !AFR && !OCEA && !SPORT {
				return "欧洲平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && !EU && AFR && !OCEA && !SPORT {
				return "非洲平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && OCEA && !SPORT {
				return "大洋洲平台"
			} else if !TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && SPORT {
				return "体育平台"
			} else {
				return ""
			}
		}
	}

	platformName := getPlatformName(M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT)

	if language == "zh" {
		if netType == "ipv4" || netType == "" {
			FormarPrint(language, platformName)
		} else if netType == "ipv6" {
			FormarPrint(language, "跨国平台")
		}
	} else if language == "en" {
		if netType == "ipv4" || netType == "" {
			enPlatformName := map[string]string{
				"跨国平台":         "Global",
				"跨国平台 + 台湾平台":  "Global + Taiwan",
				"跨国平台 + 香港平台":  "Global + Hong Kong",
				"跨国平台 + 日本平台":  "Global + Japan",
				"跨国平台 + 韩国平台":  "Global + Korea",
				"跨国平台 + 北美平台":  "Global + North America",
				"跨国平台 + 南美平台":  "Global + South America",
				"跨国平台 + 欧洲平台":  "Global + Europe",
				"跨国平台 + 非洲平台":  "Global + Africa",
				"跨国平台 + 大洋洲平台": "Global + Oceania",
				"跨国平台 + 体育平台":  "Global + Sports",
				"台湾平台":         "Taiwan",
				"香港平台":         "Hong Kong",
				"日本平台":         "Japan",
				"韩国平台":         "Korea",
				"北美平台":         "North America",
				"南美平台":         "South America",
				"欧洲平台":         "Europe",
				"非洲平台":         "Africa",
				"大洋洲平台":        "Oceania",
				"体育平台":         "Sports",
			}
			FormarPrint(language, enPlatformName[platformName])
		} else if netType == "ipv6" {
			FormarPrint(language, "Global")
		}
	}
}

func switchOptions(c string) {
	switch c {
	case "0":
		M = true
	case "1":
		M = true
		TW = true
	case "2":
		M = true
		HK = true
	case "3":
		M = true
		JP = true
	case "4":
		M = true
		KR = true
	case "5":
		M = true
		NA = true
	case "6":
		M = true
		SA = true
	case "7":
		M = true
		EU = true
	case "8":
		M = true
		AFR = true
	case "9":
		M = true
		OCEA = true
	case "10":
		TW = true
	case "11":
		HK = true
	case "12":
		JP = true
	case "13":
		KR = true
	case "14":
		NA = true
	case "15":
		SA = true
	case "16":
		EU = true
	case "17":
		AFR = true
	case "18":
		OCEA = true
	case "19":
		SPORT = true
	case "20":
		M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT = true, true, true, true, true, true, true, true, true, true, true
	default:
		M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT = false, false, false, false, false, false, false, false, false, false, false
	}
}

func ReadSelect(language, flagString string) {
	if flagString == "" {
		if language == "zh" {
			fmt.Println("请选择检测项目: ")
			fmt.Println("[0]: 跨国平台")
			fmt.Println("[1]: 跨国平台 + 台湾平台")
			fmt.Println("[2]: 跨国平台 + 香港平台")
			fmt.Println("[3]: 跨国平台 + 日本平台")
			fmt.Println("[4]: 跨国平台 + 韩国平台")
			fmt.Println("[5]: 跨国平台 + 北美平台")
			fmt.Println("[6]: 跨国平台 + 南美平台")
			fmt.Println("[7]: 跨国平台 + 欧洲平台")
			fmt.Println("[8]: 跨国平台 + 非洲平台")
			fmt.Println("[9]: 跨国平台 + 大洋洲平台")
			fmt.Println("[10]: 仅台湾平台")
			fmt.Println("[11]: 仅香港平台")
			fmt.Println("[12]: 仅日本平台")
			fmt.Println("[13]: 仅韩国平台")
			fmt.Println("[14]: 仅北美平台")
			fmt.Println("[15]: 仅南美平台")
			fmt.Println("[16]: 仅欧洲平台")
			fmt.Println("[17]: 仅非洲平台")
			fmt.Println("[18]: 仅大洋洲平台")
			fmt.Println("[19]: 仅体育平台")
			fmt.Println("[20]: 全部平台")
			fmt.Print("请输入对应数字,空格分隔(回车确认): ")
		} else {
			fmt.Println("Please select detection items:")
			fmt.Println("[0]: International platform")
			fmt.Println("[1]: International platform + Taiwan platform")
			fmt.Println("[2]: International platform + Hong Kong platform")
			fmt.Println("[3]: International platform + Japan platform")
			fmt.Println("[4]: International platform + Korea platform")
			fmt.Println("[5]: International platform + North America platform")
			fmt.Println("[6]: International platform + South America platform")
			fmt.Println("[7]: International platform + Europe platform")
			fmt.Println("[8]: International platform + Africa platform")
			fmt.Println("[9]: International platform + Oceania platform")
			fmt.Println("[10]: Taiwan platform only")
			fmt.Println("[11]: Hong Kong platform only")
			fmt.Println("[12]: Japan platform only")
			fmt.Println("[13]: Korea platform only")
			fmt.Println("[14]: North America platform only")
			fmt.Println("[15]: South America platform only")
			fmt.Println("[16]: Europe platform only")
			fmt.Println("[17]: Africa platform only")
			fmt.Println("[18]: Oceania platform only")
			fmt.Println("[19]: Sports platform only")
			fmt.Println("[20]: All platforms")
			fmt.Print("Please enter corresponding numbers, separated by spaces (press Enter to confirm): ")
		}
		r := bufio.NewReader(os.Stdin)
		l, _, err := r.ReadLine()
		if err != nil {
			fmt.Println("Failed to read select option.")
			return
		}
		for _, c := range strings.Split(string(l), " ") {
			switchOptions(c)
		}
	} else {
		switchOptions(flagString)
	}
}

var setSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	return
}

func getFuncList() [](func(c *http.Client) model.Result) {
	var funcList [](func(c *http.Client) model.Result)
	if M {
		funcList = append(funcList, Multination()...)
	}
	if TW {
		funcList = append(funcList, Taiwan()...)
	}
	if HK {
		funcList = append(funcList, HongKong()...)
	}
	if JP {
		funcList = append(funcList, Japan()...)
	}
	if KR {
		funcList = append(funcList, Korea()...)
	}
	if NA {
		funcList = append(funcList, NorthAmerica()...)
	}
	if SA {
		funcList = append(funcList, SouthAmerica()...)
	}
	if EU {
		funcList = append(funcList, Europe()...)
	}
	if AFR {
		funcList = append(funcList, Africa()...)
	}
	if OCEA {
		funcList = append(funcList, Oceania()...)
	}
	if SPORT {
		funcList = append(funcList, Sport()...)
	}
	return funcList
}

func runTests(client *http.Client, ipVersion, language string) {
	Names = []string{}
	total = 0
	wg = &sync.WaitGroup{}
	bar = NewBar(0)
	funcList := getFuncList()
	processFunction(funcList, client)
	bar.ChangeMax64(total)
	wg.Wait()
	bar.Finish()
	finallyPrintResult(language, ipVersion)
}

func setupInterface(Iface string) {
	if IP := net.ParseIP(Iface); IP != nil {
		utils.Dialer.LocalAddr = &net.TCPAddr{IP: IP}
	} else {
		utils.Dialer.Control = func(network, address string, c syscall.RawConn) error {
			return setSocketOptions(network, address, c, Iface)
		}
	}
}

func setupDnsServers(DnsServers string) {
	utils.Dialer.Resolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "udp", DnsServers)
		},
	}
}

func setupHttpProxy(httpProxy string) {
	if u, err := url.Parse(httpProxy); err == nil {
		utils.ClientProxy = http.ProxyURL(u)
		for _, transport := range []*http.Transport{utils.Ipv4Transport, utils.Ipv6Transport, utils.AutoTransport} {
			transport.Proxy = utils.ClientProxy
		}
		for _, httpClient := range []*http.Client{utils.Ipv4HttpClient, utils.Ipv6HttpClient, utils.AutoHttpClient} {
			httpClient.Transport = utils.AutoTransport
		}
	}
}

func main() {
	go func() {
		http.Get("https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Foneclickvirt%2FUnlockTests&count_bg=%2323E01C&title_bg=%23555555&icon=sonarcloud.svg&icon_color=%23E7E7E7&title=hits&edge_flat=false")
	}()
	client := utils.AutoHttpClient
	mode := 0
	showVersion := false
	Iface := ""
	DnsServers := ""
	httpProxy := ""
	language := ""
	showIP := false
	flagString := ""
	flag.IntVar(&mode, "m", 4, "mode 0(both)/4(only)/6(only), default to 4")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showIP, "s", true, "show ip address, example: -s=false")
	flag.StringVar(&flagString, "f", "", "specify select option in menu")
	flag.StringVar(&Iface, "I", "", "specify source ip / interface")
	flag.StringVar(&DnsServers, "dns-servers", "", "specify dns servers")
	flag.StringVar(&httpProxy, "http-proxy", "", "specify http proxy")
	flag.StringVar(&language, "L", "zh", "language, specify to en or zh")
	flag.Parse()
	if showVersion {
		fmt.Println(Version)
		return
	}
	if Iface != "" {
		setupInterface(Iface)
	}
	if DnsServers != "" {
		setupDnsServers(DnsServers)
	}
	if httpProxy != "" {
		fmt.Println(httpProxy)
		setupHttpProxy(httpProxy)
	}
	if mode == 4 {
		client = utils.Ipv4HttpClient
		IPV6 = false
	}
	if mode == 6 {
		client = utils.Ipv6HttpClient
		IPV4 = false
	}
	if language == "zh" {
		fmt.Println("项目地址: " + Blue("https://github.com/oneclickvirt/UnlockTests"))
	} else {
		fmt.Println("Github Repo: " + Blue("https://github.com/oneclickvirt/UnlockTests"))
	}
	if showIP {
		GetIpv4Info()
		GetIpv6Info()
	}
	ReadSelect(language, flagString)
	if IPV4 {
		fmt.Println(Blue("IPV4:"))
		runTests(client, "ipv4", language)
	}
	if IPV6 {
		fmt.Println(Blue("IPV6:"))
		if mode == 6 {
			runTests(client, "ipv6", language)
		} else {
			runTests(utils.Ipv6HttpClient, "ipv6", language)
		}
	}
	fmt.Println()
}
