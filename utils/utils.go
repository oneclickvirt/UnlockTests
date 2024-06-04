package utils

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"time"

	"golang.org/x/net/context"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

var ClientProxy = http.ProxyFromEnvironment
var AutoTransport = &http.Transport{
	Proxy:       ClientProxy,
	DialContext: (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
}
var AutoHttpClient = &http.Client{
	Timeout:   30 * time.Second,
	Transport: AutoTransport,
}
var Dialer = &net.Dialer{}
var Ipv4Transport = &http.Transport{
	Proxy: ClientProxy,
	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 强制使用IPv4
		return Dialer.DialContext(ctx, "tcp4", addr)
	},
	// ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   30 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
var Ipv4HttpClient = &http.Client{
	Timeout:   30 * time.Second,
	Transport: Ipv4Transport,
}
var Ipv6Transport = &http.Transport{
	Proxy: ClientProxy,
	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 强制使用IPv4
		return Dialer.DialContext(ctx, "tcp6", addr)
	},
	// ForceAttemptHTTP2:     true,
	MaxIdleConns:           100,
	IdleConnTimeout:        90 * time.Second,
	TLSHandshakeTimeout:    30 * time.Second,
	ExpectContinueTimeout:  1 * time.Second,
	MaxResponseHeaderBytes: 262144,
}
var Ipv6HttpClient = &http.Client{
	Timeout:   30 * time.Second,
	Transport: Ipv6Transport,
}

// ParseInterface 解析网卡IP地址
func ParseInterface(ifaceName, ipAddr, netType string) (*http.Client, error) {
	var localIP net.IP
	if ifaceName != "" {
		// 获取指定网卡的 IP 地址
		iface, err := net.InterfaceByName(ifaceName)
		if err != nil {
			return nil, err
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if (netType == "tcp4" && ipNet.IP.To4() != nil) || (netType == "tcp6" && ipNet.IP.To4() == nil) {
					localIP = ipNet.IP
					break
				}
			}
		}
	} else if ipAddr != "" {
		localIP = net.ParseIP(ipAddr)
		if (netType == "tcp4" && localIP.To4() == nil) || (netType == "tcp6" && localIP.To4() != nil) {
			return nil, fmt.Errorf("IP address does not match the specified netType")
		}
	}
	var dialContext func(ctx context.Context, network, addr string) (net.Conn, error)
	if localIP != nil {
		dialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   12 * time.Second,
				KeepAlive: 12 * time.Second,
				LocalAddr: &net.TCPAddr{
					IP: localIP,
				},
			}).DialContext(ctx, netType, addr)
		}
	} else {
		dialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, netType, addr)
		}
	}
	c := &http.Client{
		Timeout: 12 * time.Second,
		Transport: &http.Transport{
			DialContext: dialContext,
		}}
	return c, nil
}

// Req
// 为 req 设置请求
func Req(c *http.Client) *req.Client {
	client := req.DefaultClient()
	client.ImpersonateChrome()
	client.Transport.DialContext = c.Transport.(*http.Transport).DialContext
	client.SetProxy(c.Transport.(*http.Transport).Proxy)
	client.R().
		SetRetryCount(2).
		SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetRetryFixedInterval(2 * time.Second)
	return client
}

// Gorequest
// 为 gorequest 设置请求
func Gorequest(c *http.Client) *gorequest.SuperAgent {
	request := gorequest.New()
	request.Transport.DialContext = c.Transport.(*http.Transport).DialContext
	request.Transport.Proxy = c.Transport.(*http.Transport).Proxy
	request.Retry(2, 5)
	request.Timeout(12 * time.Second)
	return request
}

// SetGoRequestHeaders
func SetGoRequestHeaders(request *gorequest.SuperAgent, headers map[string]string) *gorequest.SuperAgent {
	for _, i := range headers {
		request = request.Set(i, headers[i])
	}
	return request
}

// SetReqHeaders
func SetReqHeaders(client *req.Client, headers map[string]string) *req.Client {
	for _, i := range headers {
		client.Headers.Set(i, headers[i])
	}
	return client
}

// PostJson 向指定的 URL 发送 JSON 格式的 POST 请求，并返回响应、响应体和错误信息
// url: 目标 URL
// payload: 要发送的 JSON 格式的请求体
// headers: 可选的 HTTP 头信息
func PostJson(c *http.Client, url string, payload string, headers ...map[string]string) (gorequest.Response, []byte, []error) {
	// 构建 POST 请求，设置请求类型为 JSON 并添加请求体
	request := Gorequest(c)
	request = request.Post(url).
		Type("json").
		Send(payload)
	// 添加可选的 HTTP 头信息
	for _, header := range headers {
		request = SetGoRequestHeaders(request, header)
	}
	// 发送请求并接收响应、响应体和错误信息
	resp, body, errs := request.EndBytes()
	// 返回响应、响应体和错误信息
	return resp, body, errs
}

// GetRegion
// 判断地址是否在允许的地区范围内
func GetRegion(loc string, locationList []string) bool {
	for _, s := range locationList {
		if loc == s {
			return true
		}
	}
	return false
}

// ReParse
// 根据正则表达式提取内容
func ReParse(responseBody, rex string) string {
	re := regexp.MustCompile(rex)
	match := re.FindStringSubmatch(responseBody)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// 通过Info标记要被插入的行的下一行包含什么文本内容
func PrintCA(c *http.Client) model.Result {
	return model.Result{Name: "Canada", Status: model.PrintHead, Info: "HotStar"}
}

func PrintGB(c *http.Client) model.Result {
	return model.Result{Name: "England", Status: model.PrintHead, Info: "HotStar"}
}

func PrintFR(c *http.Client) model.Result {
	return model.Result{Name: "France", Status: model.PrintHead, Info: "Canal+"}
}

func PrintDE(c *http.Client) model.Result {
	return model.Result{Name: "Germany", Status: model.PrintHead, Info: "Joyn"}
}

func PrintNL(c *http.Client) model.Result {
	return model.Result{Name: "Netherlands", Status: model.PrintHead, Info: "NLZIET"}
}

func PrintES(c *http.Client) model.Result {
	return model.Result{Name: "Spain", Status: model.PrintHead, Info: "Movistar+"}
}

func PrintIT(c *http.Client) model.Result {
	return model.Result{Name: "Italy", Status: model.PrintHead, Info: "Rai Play"}
}

func PrintCH(c *http.Client) model.Result {
	return model.Result{Name: "Switzerland", Status: model.PrintHead, Info: "SKY CH"}
}

func PrintRU(c *http.Client) model.Result {
	return model.Result{Name: "Russia", Status: model.PrintHead, Info: "Amediateka"}
}

func PrintAU(c *http.Client) model.Result {
	return model.Result{Name: "Australia", Status: model.PrintHead, Info: "Stan"}
}

func PrintNZ(c *http.Client) model.Result {
	return model.Result{Name: "New Zealand", Status: model.PrintHead, Info: "Neon TV"}
}

func PrintSG(c *http.Client) model.Result {
	return model.Result{Name: "Singapore", Status: model.PrintHead, Info: "MeWatch"}
}

func PrintTH(c *http.Client) model.Result {
	return model.Result{Name: "Thailand", Status: model.PrintHead, Info: "AIS Play"}
}

func PrintGame(c *http.Client) model.Result {
	return model.Result{Name: "Game", Status: model.PrintHead, Info: "Kancolle Japan"}
}

func PrintMusic(c *http.Client) model.Result {
	return model.Result{Name: "Music", Status: model.PrintHead, Info: "Mora"}
}

func PrintForum(c *http.Client) model.Result {
	return model.Result{Name: "Forum", Status: model.PrintHead, Info: "EroGameSpace"}
}
