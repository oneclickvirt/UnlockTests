package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/oneclickvirt/UnlockTests/utils"
	"github.com/oneclickvirt/UnlockTests/uts"
	. "github.com/oneclickvirt/defaultset"
)

func main() {
	go func() {
		http.Get("https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Foneclickvirt%2FUnlockTests&count_bg=%2323E01C&title_bg=%23555555&icon=sonarcloud.svg&icon_color=%23E7E7E7&title=hits&edge_flat=false")
	}()
	client := utils.AutoHttpClient
	mode := 0
	var showVersion, showIP, useBar bool
	var Iface, DnsServers, httpProxy, language, flagString string
	flag.IntVar(&mode, "m", 0, "mode 0(both)/4(only)/6(only), default to 0, example: -m 4")
	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showIP, "s", true, "show ip address status, disable example: -s=false")
	flag.BoolVar(&useBar, "b", true, "use progress bar, disable example: -b=false")
	flag.StringVar(&flagString, "f", "", "specify select option in menu, example: -f 0")
	flag.StringVar(&Iface, "I", "", "specify source ip / interface")
	flag.StringVar(&DnsServers, "dns-servers", "", "specify dns servers")
	flag.StringVar(&httpProxy, "http-proxy", "", "specify http proxy")
	flag.StringVar(&language, "L", "zh", "language, specify to en or zh")
	flag.Parse()
	if showVersion {
		fmt.Println(uts.UnlockTestsVersion)
		return
	}
	if Iface != "" {
		uts.SetupInterface(Iface)
	}
	if DnsServers != "" {
		uts.SetupDnsServers(DnsServers)
	}
	if httpProxy != "" {
		fmt.Println(httpProxy)
		uts.SetupHttpProxy(httpProxy)
	}
	if mode == 4 {
		client = utils.Ipv4HttpClient
		uts.IPV6 = false
	}
	if mode == 6 {
		client = utils.Ipv6HttpClient
		uts.IPV4 = false
	}
	if language == "zh" {
		fmt.Println("项目地址: " + Blue("https://github.com/oneclickvirt/UnlockTests"))
	} else {
		fmt.Println("Github Repo: " + Blue("https://github.com/oneclickvirt/UnlockTests"))
	}
	uts.GetIpv4Info(showIP)
	uts.GetIpv6Info(showIP)
	readStatus := uts.ReadSelect(language, flagString)
	if !readStatus {
		return
	}
	if language == "zh" {
		fmt.Println("测试时间: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	} else {
		fmt.Println("Test time: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	}
	if uts.IPV4 {
		fmt.Println(Blue("IPV4:"))
		fmt.Printf(uts.RunTests(client, "ipv4", language, useBar))
	}
	if uts.IPV6 {
		fmt.Println(Blue("IPV6:"))
		if mode == 6 {
			fmt.Printf(uts.RunTests(client, "ipv6", language, useBar))
		} else {
			fmt.Printf(uts.RunTests(utils.Ipv6HttpClient, "ipv6", language, useBar))
		}
	}
}
