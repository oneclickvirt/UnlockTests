package executor

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"sort"
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
	"github.com/oneclickvirt/UnlockTests/vn"
	. "github.com/oneclickvirt/defaultset"
	pb "github.com/schollz/progressbar/v3"
)

var (
	total                                               int64
	bar                                                 *pb.ProgressBar
	wg                                                  *sync.WaitGroup
	IPV4, IPV6                                          = true, true
	R                                                   []*model.Result
	resultMutex                                         sync.Mutex
	Names                                               []string
	M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT, AI = false, false, false, false, false, false, false, false, false, false, false, false
	sem                                                 chan struct{}
	cacheEnabled                                        = false
	resultCache                                         = make(map[string]model.Result)
	cacheMutex                                          sync.RWMutex
	runTestsMutex                                       sync.Mutex
)

const testExecutionTimeout = 30 * time.Second

func NewBar(count int64) *pb.ProgressBar {
	return pb.NewOptions64(
		count,
		pb.OptionSetDescription("testing"),
		pb.OptionSetWriter(utils.ColorStderr),
		pb.OptionSetWidth(20),
		pb.OptionThrottle(100*time.Millisecond),
		pb.OptionShowCount(),
		pb.OptionClearOnFinish(),
		pb.OptionEnableColorCodes(true),
		pb.OptionSpinnerType(14),
	)
}

func ShowResult(r *model.Result) (s string) {
	logErr := func() {
		if model.EnableLoger && r.Err != nil {
			InitLogger()
			defer Logger.Sync()
			Logger.Info(r.Name + " " + r.Err.Error())
		}
	}
	formatResult := func(colorFunc func(string) string, status string, r model.Result) string {
		s := colorFunc(status)
		if r.Info != "" {
			s += colorFunc(" (" + r.Info + ")")
		}
		if r.Region != "" {
			s += colorFunc(" (Region: " + strings.ToUpper(r.Region) + ")")
		}
		if r.UnlockType != "" {
			s += colorFunc(" [" + r.UnlockType + "]")
		}
		return s
	}
	switch r.Status {
	case model.StatusYes:
		return formatResult(Green, "YES", *r)
	case model.StatusNetworkErr:
		logErr()
		return Red("Failed") + Yellow(" (Network Error)")
	case model.StatusNoIPv6:
		return Yellow("N/A") + White(" (No IPv6 Support)")
	case model.StatusDNSFailed:
		logErr()
		return Yellow("N/A") + White(" (DNS Resolve Failed)")
	case model.StatusRestricted:
		return formatResult(Yellow, "Restricted", *r)
	case model.StatusErr:
		s = Yellow("Error")
		logErr()
		return s
	case model.StatusNo:
		return formatResult(Red, "NO", *r)
	case model.StatusBanned:
		s = Red("Banned")
		if r.Info != "" {
			s += Yellow(" (" + r.Info + ")")
		}
		return s
	case model.StatusTimeout:
		logErr()
		return Yellow("TIMEOUT")
	case model.StatusCDNRelay:
		return Blue(model.StatusCDNRelay)
	case model.StatusUnexpected:
		s = Purple("Unknown")
		if r.Err != nil {
			s += ": " + r.Err.Error()
			logErr()
		}
		return s
	default:
		s = Purple("Unknown")
		if r.Status != "" {
			s += White(" (" + r.Status + ")")
		}
		if r.Err != nil {
			s += ": " + r.Err.Error()
			logErr()
		}
		return s
	}
}

func RemoveDuplicates(input []string) []string {
	// 创建一个映射来跟踪已经遇到的字符串
	seen := make(map[string]bool)
	// 创建一个新的切片来存储去重后的结果
	result := []string{}
	// 遍历输入的字符串切片
	for _, str := range input {
		// 如果字符串没有在映射中出现过，则添加到结果切片中
		if !seen[str] {
			result = append(result, str)
			// 将字符串标记为已出现
			seen[str] = true
		}
	}
	return result
}

func PrintCenteredMessage(message string, totalLength int) string {
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

func FormarPrint(message string) string {
	Length := 25
	for _, r := range R {
		if len(r.Name) > Length {
			Length = len(r.Name)
		}
	}
	head := PrintCenteredMessage("[ "+message+" ]", 0)
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
			tempList = append(tempList, fmt.Sprintf("%-"+strconv.Itoa(Length)+"s %s\n", r.Name, result))
		}
	}
	// 插入小分区的head行
	if !M || !TW || !HK || !JP || !KR || !NA || !SA || !EU || !AFR || !OCEA || !SPORT || !AI {
		for _, r := range R {
			if r.Status == model.PrintHead {
				anotherList := []string{}
				for _, i := range tempList {
					if strings.Contains(i, r.Info) {
						tpHead := PrintCenteredMessage("[ "+r.Name+" ]", 20)
						anotherList = append(anotherList, tpHead)
					}
					anotherList = append(anotherList, i)
				}
				tempList = anotherList
			}
		}
	}
	// 去重
	tempList = RemoveDuplicates(tempList)
	var res strings.Builder
	for _, i := range tempList {
		res.WriteString(i)
	}
	return res.String()
}

func formatIPVersionLabel(netType string) string {
	switch strings.ToLower(strings.TrimSpace(netType)) {
	case "ipv4", "tcp4", "4":
		return "IPV4"
	case "ipv6", "tcp6", "6":
		return "IPV6"
	default:
		return strings.ToUpper(strings.TrimSpace(netType))
	}
}

func formatVersionedHeader(netType, title string) string {
	ipLabel := formatIPVersionLabel(netType)
	title = strings.TrimSpace(title)
	switch {
	case ipLabel == "":
		return title
	case title == "":
		return ipLabel
	default:
		return ipLabel + " " + title
	}
}

func Excute(F func(c *http.Client) model.Result, c *http.Client, useProgressBar bool, ipVersion string) {
	testInfo := F(nil)
	testName := testInfo.Name
	wg.Add(1)
	go func() {
		defer wg.Done()
		// panic 恢复机制
		defer func() {
			if r := recover(); r != nil {
				panicResult := model.Result{
					Name:   testName,
					Status: model.StatusErr,
					Err:    fmt.Errorf("panic recovered: %v", r),
				}
				resultMutex.Lock()
				R = append(R, &panicResult)
				resultMutex.Unlock()
				if useProgressBar {
					bar.Describe(panicResult.Name + " " + ShowResult(&panicResult))
					bar.Add(1)
				}
			}
		}()

		// 并发控制
		if sem != nil {
			sem <- struct{}{} // 获取一个通道资源
			defer func() {
				<-sem // 释放一个通道资源
			}()
		}
		// 检查缓存（区分 IPv4/IPv6）
		if cacheEnabled {
			cacheKey := resultCacheKey(testName, ipVersion, c)
			cacheMutex.RLock()
			if cachedResult, exists := resultCache[cacheKey]; exists {
				cacheMutex.RUnlock()
				resultMutex.Lock()
				R = append(R, &cachedResult)
				resultMutex.Unlock()
				if useProgressBar {
					bar.Describe(cachedResult.Name + " " + ShowResult(&cachedResult))
					bar.Add(1)
				}
				return
			}
			cacheMutex.RUnlock()
		}

		// 执行测试
		res := runTestWithTimeout(F, c, testName)
		res = utils.NormalizeResult(c, res, testName)

		// 保存到缓存（区分 IPv4/IPv6）
		if cacheEnabled {
			cacheKey := resultCacheKey(testName, ipVersion, c)
			cacheMutex.Lock()
			resultCache[cacheKey] = res
			cacheMutex.Unlock()
		}

		resultMutex.Lock()
		R = append(R, &res)
		resultMutex.Unlock()
		if useProgressBar {
			bar.Describe(res.Name + " " + ShowResult(&res))
			bar.Add(1)
		}
	}()
}

func PreProcess(FuncList [](func(c *http.Client) model.Result)) {
	// 生成顺序输出的名字
	for _, f := range FuncList {
		tp := f(nil)
		if tp.Status != model.PrintHead {
			Names = append(Names, tp.Name)
		}
	}
}

func prepareFuncList(FuncList [](func(c *http.Client) model.Result)) [](func(c *http.Client) model.Result) {
	FuncList = sortedFuncList(FuncList)
	PreProcess(FuncList)
	return FuncList
}

func sortedFuncList(FuncList [](func(c *http.Client) model.Result)) [](func(c *http.Client) model.Result) {
	sorted := append([](func(c *http.Client) model.Result)(nil), FuncList...)
	sort.SliceStable(sorted, func(i, j int) bool {
		left := functionSortKey(sorted[i])
		right := functionSortKey(sorted[j])
		if left == right {
			return functionDisplayName(sorted[i]) < functionDisplayName(sorted[j])
		}
		return left < right
	})
	return sorted
}

func functionDisplayName(f func(c *http.Client) model.Result) string {
	tp := f(nil)
	if tp.Name != "" {
		return tp.Name
	}
	return fmt.Sprintf("%p", f)
}

func functionSortKey(f func(c *http.Client) model.Result) string {
	return strings.ToLower(functionDisplayName(f))
}

func namesFromFunctions(FuncList [](func(c *http.Client) model.Result)) []string {
	names := make([]string, 0, len(FuncList))
	for _, f := range FuncList {
		tp := f(nil)
		if tp.Status != model.PrintHead && tp.Name != "" {
			names = append(names, tp.Name)
		}
	}
	return names
}

func resultCacheKey(testName, ipVersion string, c *http.Client) string {
	transportID := "<nil>"
	if c != nil && c.Transport != nil {
		transportID = fmt.Sprintf("%p", c.Transport)
	}
	return testName + "_" + ipVersion + "_" + transportID
}

func runTestWithTimeout(F func(c *http.Client) model.Result, c *http.Client, testName string) model.Result {
	ctx, cancel := context.WithTimeout(context.Background(), testExecutionTimeout)
	defer cancel()
	testClient := clientWithContextDeadline(c, ctx)
	resultCh := make(chan model.Result, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultCh <- model.Result{
					Name:   testName,
					Status: model.StatusErr,
					Err:    fmt.Errorf("panic recovered: %v", r),
				}
			}
		}()
		resultCh <- F(testClient)
	}()
	select {
	case res := <-resultCh:
		return res
	case <-ctx.Done():
		return model.Result{
			Name:   testName,
			Status: model.StatusTimeout,
			Err:    fmt.Errorf("test timeout (%v)", testExecutionTimeout),
		}
	}
}

func uniqueFuncList(FuncList [](func(c *http.Client) model.Result)) [](func(c *http.Client) model.Result) {
	seen := make(map[string]bool, len(FuncList))
	unique := make([](func(c *http.Client) model.Result), 0, len(FuncList))
	for _, f := range FuncList {
		tp := f(nil)
		key := tp.Name
		if key == "" {
			key = fmt.Sprintf("%p", f)
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		unique = append(unique, f)
	}
	return unique
}

func ProcessFunction(FuncList [](func(c *http.Client) model.Result), c *http.Client, useProgressBar bool, ipVersion string) {
	// 实际开始任务
	for _, f := range FuncList {
		Excute(f, c, useProgressBar, ipVersion)
	}
}

func NorthAmerica() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		us.AcornTV,
		us.AETV,
		us.AMCPlus,
		uk.BritBox,
		us.CWTV,
		us.Crackle,
		us.Crunchyroll,
		us.DirectvStream,
		us.DiscoveryPlus,
		us.EncoreTVB,
		us.FXNOW,
		us.Fox,
		us.HBOMax,
		asia.HotStar,
		us.Hulu,
		us.MGMPlus,
		us.NFLPlus,
		us.PeacockTV,
		us.Philo,
		us.PlutoTV,
		us.Pornhub,
		us.SlingTV,
		us.Starz,
		us.Shudder,
		us.TLCGO,
		us.TubiTV,
		eu.Viaplay,
		// CA
		utils.PrintCA,
		ca.CBCGem,
		ca.Crave,
	}
	return prepareFuncList(FuncList)
}

func Europe() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		eu.RakutenTV,
		eu.SkyShowTime,
		us.HBOMax,
		eu.MegogoTV,
		eu.TNTSports,
		eu.Viaplay,
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
		fr.FranceTV,
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
	return prepareFuncList(FuncList)
}

func HongKong() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		us.HBOMax,
		hk.HoyTV,
		hk.MyTvSuper,
		hk.NowTV,
		hk.ViuTV,
	}
	return prepareFuncList(FuncList)
}

func Africa() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		africa.DSTV,
		africa.Showmax,
	}
	return prepareFuncList(FuncList)
}

func India() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.HotStar,
		in.TataPlay,
		in.MXPlayer,
		in.Zee5,
	}
	return prepareFuncList(FuncList)
}

func Taiwan() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		tw.BahamutAnime,
		tw.Catchplay,
		tw.FridayVideo,
		tw.HamiVideo,
		us.HBOMax,
		tw.KKTV,
		tw.LiTV,
		tw.LineTV,
		tw.MyVideo,
		tw.Ofiii,
		tw.Tw4gtv,
	}
	return prepareFuncList(FuncList)
}

func Japan() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		jp.Abema,
		jp.DAnimeStore,
		jp.DMM,
		jp.DMMTV,
		jp.FOD,
		jp.Hulu,
		jp.J_COM_ON_DEMAND,
		jp.Lemino,
		jp.MGStage,
		jp.NHKPlus,
		jp.Niconico,
		jp.Radiko,
		jp.RakutenMagazine,
		jp.RakutenTV,
		jp.SDGGGE,
		jp.TVer,
		jp.Telasa,
		jp.UNext,
		jp.VideoMarket,
		jp.Wowow,
		jp.AnimeFesta,
		// Game
		utils.PrintGame,
		jp.Kancolle,
		jp.PrettyDerby,
		jp.PCRJP,
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
	return prepareFuncList(FuncList)
}

func Multination() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		transnation.Apple,
		asia.BilibiliAnime,
		transnation.Bing,
		transnation.Claude,
		transnation.Copilot,
		transnation.Coze,
		transnation.DeepSeek,
		transnation.DisneyPlus,
		transnation.Gemini,
		transnation.GoogleSearch,
		transnation.GooglePlayStore,
		transnation.Grok,
		transnation.IQiYi,
		transnation.Instagram,
		transnation.Kimi,
		transnation.KOCOWA,
		eu.MathsSpot,
		transnation.MetaAI,
		transnation.MistralAI,
		transnation.Netflix,
		transnation.NetflixCDN,
		transnation.OneTrust,
		transnation.OpenAI,
		transnation.ParamountPlus,
		transnation.PerplexityAI,
		transnation.Poe,
		transnation.PrimeVideo,
		transnation.Reddit,
		transnation.SonyLiv,
		transnation.Sora,
		transnation.Spotify,
		transnation.Steam,
		transnation.TVBAnywhere,
		transnation.TikTok,
		transnation.ViuCom,
		transnation.WeTV,
		transnation.WikipediaEditable,
		transnation.Youtube,
		transnation.YoutubeCDN,
	}
	return prepareFuncList(FuncList)
}

func AIPlatforms() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		transnation.Claude,
		transnation.Copilot,
		transnation.Coze,
		transnation.DeepSeek,
		transnation.Gemini,
		transnation.Grok,
		transnation.Kimi,
		transnation.MetaAI,
		transnation.MistralAI,
		transnation.OpenAI,
		transnation.PerplexityAI,
		transnation.Poe,
		transnation.Sora,
	}
	return prepareFuncList(FuncList)
}

func SouthAmerica() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		us.HBOMax,
		us.DirecTVGO,
	}
	return prepareFuncList(FuncList)
}

func Oceania() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		us.AcornTV,
		us.AMCPlus,
		uk.BritBox,
		// AU
		utils.PrintAU,
		au.ABCiView,
		au.Au7plus,
		au.Channel10,
		au.Channel9,
		au.KayoSports,
		au.SBSonDemand,
		au.Stan,
		eu.Docplay,
		// NZ
		utils.PrintNZ,
		nz.NeonTV,
		nz.SkyGO,
		nz.ThreeNow,
		nz.MaoriTV,
	}
	return prepareFuncList(FuncList)
}

func Korea() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		kr.SOOP,
		kr.CoupangPlay,
		kr.KBSDomestic,
		kr.NaverTV,
		kr.PandaTV,
		kr.SPOTVNOW,
		kr.Tving,
		kr.Watcha,
		kr.Wavve,
	}
	return prepareFuncList(FuncList)
}

func SouthEastAsia() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.HotStar,
		us.HBOMax,
		// SG
		utils.PrintSG,
		sg.MeWatch,
		// TH
		utils.PrintTH,
		th.AISPlay,
		th.TrueID,
		// th.BilibiliTH, 失效 - 不做检测
		// ID 全失效 - 不做检测
		utils.PrintVN,
		vn.GalaxyPlay,
		vn.KPLUS,
		vn.TV360,
		utils.PrintMY,
		asia.Sooka,
	}
	return prepareFuncList(FuncList)
}

func Sport() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		africa.BeinConnect,
		asia.StarPlus,
		eu.Eurosport,
		eu.SetantaSports,
		transnation.DAZN,
		us.ESPNPlus,
		us.FuboTV,
		us.NBATV,
		us.NBCTV,
	}
	return prepareFuncList(FuncList)
}

func IPV6Multination() [](func(c *http.Client) model.Result) {
	var FuncList = [](func(c *http.Client) model.Result){
		asia.BilibiliAnime,
		transnation.Apple,
		transnation.Bing,
		transnation.OpenAI,
		transnation.Claude,
		transnation.Copilot,
		transnation.Coze,
		transnation.DeepSeek,
		transnation.DisneyPlus,
		transnation.Gemini,
		transnation.GooglePlayStore,
		transnation.Grok,
		transnation.Kimi,
		transnation.MetaAI,
		transnation.MistralAI,
		transnation.Netflix,
		transnation.NetflixCDN,
		transnation.PerplexityAI,
		transnation.Poe,
		transnation.Sora,
		transnation.Spotify,
		transnation.WeTV,
		transnation.WikipediaEditable,
		transnation.Youtube,
		transnation.YoutubeCDN,
	}
	return prepareFuncList(FuncList)
}

func finallyPrintResult(language, netType string) string {
	var result string
	getPlatformName := func(multi bool, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT, AI bool) string {
		if TW && HK && JP && KR && NA && SA && EU && AFR && OCEA && SPORT && AI {
			return "All Platform"
		}
		if AI && !TW && !HK && !JP && !KR && !NA && !SA && !EU && !AFR && !OCEA && !SPORT {
			if multi {
				return "Global + AI"
			}
			return "AI"
		}
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
			} else if TW && HK && JP && KR && NA && SA && EU && AFR && OCEA && SPORT {
				return "所有平台"
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

	platformName := getPlatformName(M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT, AI)

	switch language {
	case "zh":
		result += FormarPrint(formatVersionedHeader(netType, platformName))
	case "en":
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
			"所有平台":         "All Platform",
		}
		displayName := enPlatformName[platformName]
		if displayName == "" {
			displayName = platformName
		}
		result += FormarPrint(formatVersionedHeader(netType, displayName))
	}
	return result
}

func resetOptions() {
	M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT, AI = false, false, false, false, false, false, false, false, false, false, false, false
}

func SwitchOptions(c string) bool {
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
		M, TW, HK, JP, KR, NA, SA, EU, AFR, OCEA, SPORT, AI = true, true, true, true, true, true, true, true, true, true, true, true
	case "21":
		AI = true
	default:
		return false
	}
	return true
}

func parseSelection(flagString string) bool {
	resetOptions()
	fields := strings.FieldsFunc(flagString, func(r rune) bool {
		return r == ' ' || r == '\t' || r == '\n' || r == '\r' || r == ','
	})
	if len(fields) == 0 {
		return false
	}
	for _, c := range fields {
		if !SwitchOptions(c) {
			resetOptions()
			return false
		}
	}
	return true
}

func ReadSelect(language, flagString string) bool {
	if flagString == "" {
		prompt := ""
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
			fmt.Println("[21]: 仅 AI 平台")
			prompt = "请输入对应数字,空格分隔(回车确认): "
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
			fmt.Println("[21]: AI platforms only")
			prompt = "Please enter corresponding numbers, separated by spaces (press Enter to confirm): "
		}
		l, err := readLine(prompt)
		if err != nil {
			fmt.Println("Failed to read select option.")
			return false
		}
		if !parseSelection(l) {
			fmt.Println("Invalid select option.")
			return false
		}
	} else {
		if !parseSelection(flagString) {
			fmt.Println("Invalid select option.")
			return false
		}
	}
	return true
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
	if AI {
		funcList = append(funcList, AIPlatforms()...)
	}
	return funcList
}

func allPlatformFuncList() [](func(c *http.Client) model.Result) {
	var funcList [](func(c *http.Client) model.Result)
	for _, build := range []func() [](func(c *http.Client) model.Result){
		Multination,
		Taiwan,
		HongKong,
		Japan,
		Korea,
		NorthAmerica,
		SouthAmerica,
		Europe,
		Africa,
		Oceania,
		Sport,
		AIPlatforms,
	} {
		funcList = append(funcList, build()...)
	}
	return sortedFuncList(uniqueFuncList(funcList))
}

func parseTestNames(testNames string) []string {
	fields := strings.FieldsFunc(testNames, func(r rune) bool {
		return r == ','
	})
	names := make([]string, 0, len(fields))
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field != "" {
			names = append(names, field)
		}
	}
	return names
}

func normalizeTestName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	replacer := strings.NewReplacer(" ", "", "-", "", "_", "", ".", "", "+", "", "'", "", "\"", "")
	return replacer.Replace(name)
}

func functionShortName(f func(c *http.Client) model.Result) string {
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	if fullName == "" {
		return ""
	}
	if idx := strings.LastIndex(fullName, "."); idx >= 0 {
		fullName = fullName[idx+1:]
	}
	return strings.TrimSuffix(fullName, "-fm")
}

func functionMatchesTestName(f func(c *http.Client) model.Result, target string) bool {
	target = normalizeTestName(target)
	if target == "" {
		return false
	}
	info := f(nil)
	for _, candidate := range []string{info.Name, functionShortName(f)} {
		if normalizeTestName(candidate) == target {
			return true
		}
	}
	return false
}

func functionsForTestNamesLocked(testNames string) ([]func(c *http.Client) model.Result, []string, []string) {
	state := snapshotSelectionState()
	defer restoreSelectionState(state)

	Names = nil
	R = nil
	targets := parseTestNames(testNames)
	candidates := allPlatformFuncList()
	funcs := make([]func(c *http.Client) model.Result, 0, len(targets))
	missing := make([]string, 0)
	for _, target := range targets {
		found := false
		for _, f := range candidates {
			if functionMatchesTestName(f, target) {
				funcs = append(funcs, f)
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, target)
		}
	}
	funcs = sortedFuncList(uniqueFuncList(funcs))
	return funcs, namesFromFunctions(funcs), missing
}

func RunNamedTests(client *http.Client, ipVersion, language string, useProgressBar bool, testNames string) (string, error) {
	runTestsMutex.Lock()
	defer runTestsMutex.Unlock()
	Names = []string{}
	resultMutex.Lock()
	R = nil
	resultMutex.Unlock()
	total = 0
	wg = &sync.WaitGroup{}
	utils.SetDNSIPVersion(ipVersion)
	defer utils.SetDNSIPVersion("")
	funcList, names, missing := functionsForTestNamesLocked(testNames)
	if len(missing) > 0 {
		return "", fmt.Errorf("test not found: %s", strings.Join(missing, ", "))
	}
	if len(funcList) == 0 {
		return "", fmt.Errorf("no valid test specified")
	}
	Names = names
	total = int64(len(funcList))
	if useProgressBar {
		bar = NewBar(total)
	}
	ProcessFunction(funcList, client, useProgressBar, ipVersion)
	wg.Wait()
	if useProgressBar {
		bar.Finish()
		fmt.Fprint(utils.ColorStderr, "\r\033[K")
		time.Sleep(50 * time.Millisecond)
	}
	return FormarPrint(formatVersionedHeader(ipVersion, "Selected Tests")), nil
}

func RunTests(client *http.Client, ipVersion, language string, useProgressBar bool) string {
	runTestsMutex.Lock()
	defer runTestsMutex.Unlock()
	Names = []string{}
	resultMutex.Lock()
	R = nil
	resultMutex.Unlock()
	total = 0
	wg = &sync.WaitGroup{}
	utils.SetDNSIPVersion(ipVersion)
	defer utils.SetDNSIPVersion("")
	funcList := getFuncList()
	funcList = sortedFuncList(uniqueFuncList(funcList))
	Names = namesFromFunctions(funcList)
	total = int64(len(funcList))
	if useProgressBar {
		bar = NewBar(total)
	}
	ProcessFunction(funcList, client, useProgressBar, ipVersion)
	wg.Wait()
	if useProgressBar {
		bar.Finish()
		// 确保进度条完全清除后再输出结果，避免显示重叠
		// 先清除当前行，然后给一个短暂延时确保终端状态更新
		fmt.Fprint(utils.ColorStderr, "\r\033[K")
		time.Sleep(50 * time.Millisecond)
	}
	return finallyPrintResult(language, ipVersion)
}
