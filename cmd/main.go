package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/oneclickvirt/UnlockTests/executor"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
)

func main() {
	go func() {
		http.Get("https://hits.spiritlhl.net/UnlockTests.svg?action=hit&title=Hits&title_bg=%23555555&count_bg=%230eecf8&edge_flat=false")
	}()
	client := utils.AutoHttpClient
	mode := 0
	var showVersion, help, showIP, useBar bool
	var Iface, DnsServers, httpProxy, language, flagString string
	utFlag := flag.NewFlagSet("ut", flag.ContinueOnError)
	utFlag.BoolVar(&help, "h", false, "show help information")
	utFlag.IntVar(&mode, "m", 0, "mode: 0 (both), 4 (only), or 6 (only); default is 0, example: -m 4")
	utFlag.BoolVar(&showVersion, "v", false, "show version")
	utFlag.BoolVar(&showIP, "s", true, "show IP address status; to disable, use: -s=false")
	utFlag.BoolVar(&useBar, "b", true, "use progress bar; to disable, use: -b=false")
	utFlag.BoolVar(&model.EnableLoger, "log", false, "enable logging")
	utFlag.StringVar(&flagString, "f", "", "specify selection option in menu; example: -f 0")
	utFlag.StringVar(&Iface, "I", "", "specify source IP/interface")
	utFlag.StringVar(&DnsServers, "dns-servers", "", "specify DNS servers")
	utFlag.StringVar(&httpProxy, "http-proxy", "", "specify HTTP proxy")
	utFlag.StringVar(&language, "L", "zh", "language; specify 'en' for English or 'zh' for Chinese")
	utFlag.Parse(os.Args[1:])
	if help {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		utFlag.PrintDefaults()
		return
	}
	if showVersion {
		fmt.Println(model.UnlockTestsVersion)
		return
	}
	if Iface != "" {
		executor.SetupInterface(Iface)
	}
	if DnsServers != "" {
		executor.SetupDnsServers(DnsServers)
	}
	if httpProxy != "" {
		fmt.Println(httpProxy)
		executor.SetupHttpProxy(httpProxy)
	}
	if mode == 4 {
		client = utils.Ipv4HttpClient
		executor.IPV6 = false
	}
	if mode == 6 {
		client = utils.Ipv6HttpClient
		executor.IPV4 = false
	}
	if language == "zh" {
		fmt.Println("项目地址: " + Blue("https://github.com/oneclickvirt/UnlockTests"))
	} else {
		fmt.Println("Github Repo: " + Blue("https://github.com/oneclickvirt/UnlockTests"))
	}
	executor.GetIpv4Info(showIP)
	executor.GetIpv6Info(showIP)
	readStatus := executor.ReadSelect(language, flagString)
	if !readStatus {
		return
	}
	if language == "zh" {
		fmt.Println("测试时间: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	} else {
		fmt.Println("Test time: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	}
	if executor.IPV4 {
		fmt.Println(Blue("IPV4:"))
		fmt.Print(executor.RunTests(client, "ipv4", language, useBar))
	}
	if executor.IPV6 {
		fmt.Println(Blue("IPV6:"))
		if mode == 6 {
			fmt.Print(executor.RunTests(client, "ipv6", language, useBar))
		} else {
			fmt.Print(executor.RunTests(utils.Ipv6HttpClient, "ipv6", language, useBar))
		}
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}
}
