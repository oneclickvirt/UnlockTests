package utils

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"sync"
	"time"

	"github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/model"
	. "github.com/oneclickvirt/defaultset"
)

const (
	DefaultRequestTimeout = 14 * time.Second
	DefaultClientTimeout  = 30 * time.Second
	DefaultRetryCount     = 1
)

// ClientProxy is nil by default. Proxy selection is an explicit API/CLI
// option; ambient HTTP(S)_PROXY variables must not alter test results.
var ClientProxy func(*http.Request) (*url.URL, error)
var Dialer = &net.Dialer{Timeout: DefaultRequestTimeout, KeepAlive: 30 * time.Second}
var AutoTransport = &http.Transport{
	Proxy:       ClientProxy,
	DialContext: Dialer.DialContext,
}
var AutoHttpClient = &http.Client{
	Timeout:   DefaultClientTimeout,
	Transport: AutoTransport,
}
var Ipv4Transport = &http.Transport{
	Proxy: ClientProxy,
	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 强制使用IPv4
		return Dialer.DialContext(ctx, "tcp4", addr)
	},
	// ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   DefaultRequestTimeout,
	ResponseHeaderTimeout: DefaultRequestTimeout,
	ExpectContinueTimeout: 1 * time.Second,
}
var Ipv4HttpClient = &http.Client{
	Timeout:   DefaultClientTimeout,
	Transport: Ipv4Transport,
}
var Ipv6Transport = &http.Transport{
	Proxy: ClientProxy,
	DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
		// 强制使用IPv6
		return Dialer.DialContext(ctx, "tcp6", addr)
	},
	// ForceAttemptHTTP2:     true,
	MaxIdleConns:           100,
	IdleConnTimeout:        90 * time.Second,
	TLSHandshakeTimeout:    DefaultRequestTimeout,
	ResponseHeaderTimeout:  DefaultRequestTimeout,
	ExpectContinueTimeout:  1 * time.Second,
	MaxResponseHeaderBytes: 262144,
}
var Ipv6HttpClient = &http.Client{
	Timeout:   DefaultClientTimeout,
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
				Timeout:   DefaultRequestTimeout,
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
		Timeout: DefaultClientTimeout,
		Transport: &http.Transport{
			DialContext:           dialContext,
			TLSHandshakeTimeout:   DefaultRequestTimeout,
			ResponseHeaderTimeout: DefaultRequestTimeout,
		}}
	return c, nil
}

// Req
// 为 req 设置请求
func Req(c *http.Client) *req.Client {
	client := req.C().Clone()
	client.ImpersonateChrome()
	configureReqTransport(client, c)
	client.R().
		SetRetryCount(DefaultRetryCount).
		SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetRetryFixedInterval(2 * time.Second)
	client.SetTimeout(effectiveReqTimeout(c, DefaultRequestTimeout))
	return client
}

// ReqDefault
// 为 req 设置请求
func ReqDefault(c *http.Client) *req.Client {
	client := req.C().Clone()
	if client.Headers == nil {
		client.Headers = make(http.Header)
	}
	configureReqTransport(client, c)
	client.R().
		SetRetryCount(DefaultRetryCount).
		SetRetryBackoffInterval(1*time.Second, 5*time.Second).
		SetRetryFixedInterval(2 * time.Second)
	client.SetTimeout(effectiveReqTimeout(c, DefaultRequestTimeout))
	return client
}

func effectiveReqTimeout(c *http.Client, fallback time.Duration) time.Duration {
	if c != nil && c.Timeout > 0 && c.Timeout < fallback {
		return c.Timeout
	}
	return fallback
}

type callerContextRoundTripper struct {
	base   http.RoundTripper
	caller context.Context
}

func (t *callerContextRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return roundTripWithCallerContext(t.base, request, t.caller)
}

// WithCallerContext clones an HTTP client and makes every request inherit the
// caller context, including requests created by provider helpers that do not
// expose context.Context in their public signature.
func WithCallerContext(client *http.Client, caller context.Context) *http.Client {
	if client == nil || caller == nil {
		return client
	}
	clone := *client
	base := client.Transport
	if base == nil {
		base = http.DefaultTransport
	} else {
		base, _ = unwrapCallerContextTransport(base)
	}
	clone.Transport = &callerContextRoundTripper{base: base, caller: caller}
	return &clone
}

func unwrapCallerContextTransport(roundTripper http.RoundTripper) (http.RoundTripper, context.Context) {
	var caller context.Context
	for {
		wrapped, ok := roundTripper.(*callerContextRoundTripper)
		if !ok {
			return roundTripper, caller
		}
		if caller == nil {
			caller = wrapped.caller
		}
		roundTripper = wrapped.base
	}
}

func roundTripWithCallerContext(roundTripper http.RoundTripper, request *http.Request, caller context.Context) (*http.Response, error) {
	if caller == nil {
		return roundTripper.RoundTrip(request)
	}
	ctx, cancel := context.WithCancel(request.Context())
	stop := context.AfterFunc(caller, cancel)
	if caller.Err() != nil {
		cancel()
	}
	var once sync.Once
	cleanup := func() {
		once.Do(func() {
			stop()
			cancel()
		})
	}
	response, err := roundTripper.RoundTrip(request.Clone(ctx))
	if err != nil {
		cleanup()
		return nil, err
	}
	if response.Body == nil {
		cleanup()
		return response, nil
	}
	response.Body = &contextBody{ReadCloser: response.Body, cleanup: cleanup}
	return response, nil
}

type contextBody struct {
	io.ReadCloser
	cleanup func()
}

func (b *contextBody) Close() error {
	err := b.ReadCloser.Close()
	b.cleanup()
	return err
}

func configureReqTransport(client *req.Client, c *http.Client) {
	if client == nil || c == nil {
		return
	}
	roundTripper, caller := unwrapCallerContextTransport(c.Transport)
	transport, ok := roundTripper.(*http.Transport)
	if !ok || transport == nil {
		return
	}
	cloned := transport.Clone()
	client.Transport.Proxy = cloned.Proxy
	client.Transport.OnProxyConnectResponse = cloned.OnProxyConnectResponse
	client.Transport.DialContext = cloned.DialContext
	client.Transport.DialTLSContext = cloned.DialTLSContext
	client.Transport.TLSClientConfig = cloned.TLSClientConfig
	client.Transport.TLSHandshakeTimeout = cloned.TLSHandshakeTimeout
	client.Transport.DisableKeepAlives = cloned.DisableKeepAlives
	client.Transport.DisableCompression = cloned.DisableCompression
	client.Transport.MaxIdleConns = cloned.MaxIdleConns
	client.Transport.MaxIdleConnsPerHost = cloned.MaxIdleConnsPerHost
	client.Transport.MaxConnsPerHost = cloned.MaxConnsPerHost
	client.Transport.IdleConnTimeout = cloned.IdleConnTimeout
	client.Transport.ResponseHeaderTimeout = cloned.ResponseHeaderTimeout
	client.Transport.ExpectContinueTimeout = cloned.ExpectContinueTimeout
	client.Transport.ProxyConnectHeader = cloned.ProxyConnectHeader
	client.Transport.GetProxyConnectHeader = cloned.GetProxyConnectHeader
	client.Transport.MaxResponseHeaderBytes = cloned.MaxResponseHeaderBytes
	client.Transport.WriteBufferSize = cloned.WriteBufferSize
	client.Transport.ReadBufferSize = cloned.ReadBufferSize
	if caller != nil {
		client.GetTransport().WrapRoundTripFunc(func(roundTripper http.RoundTripper) req.HttpRoundTripFunc {
			return func(request *http.Request) (*http.Response, error) {
				return roundTripWithCallerContext(roundTripper, request, caller)
			}
		})
	}
}

// SetReqHeaders
func SetReqHeaders(client *req.Client, headers map[string]string) *req.Client {
	for key, value := range headers {
		client.Headers.Set(key, value)
	}
	return client
}

// PostJson 向指定的 URL 发送 JSON 格式的 POST 请求，并返回响应、响应体和错误信息
// url: 目标 URL
// payload: 要发送的 JSON 格式的请求体
// headers: 可选的 HTTP 头信息
func PostJson(c *http.Client, url string, payload string, headers map[string]string) (*req.Response, string, error) {
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	// 构建 POST 请求，设置请求类型为 JSON 并添加请求体
	request := ReqDefault(c)
	// 添加可选的 HTTP 头信息
	if headers != nil {
		request = SetReqHeaders(request, headers)
	}
	resp, err := request.R().SetBodyJsonString(payload).Post(url)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("PostJson failed: " + err.Error())
		}
		return resp, "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if model.EnableLoger {
			Logger.Info("read resp.Body failed: " + err.Error())
		}
		return resp, "", err
	}
	body := string(b)
	return resp, body, err
}

// GetRegion
// 判断地址是否在允许的地区范围内
func GetRegion(loc string, locationList []string) bool {
	return slices.Contains(locationList, loc)
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
	return model.Result{Name: "Canada", Status: model.PrintHead, Info: "CBC Gem"}
}

func PrintGB(c *http.Client) model.Result {
	return model.Result{Name: "England", Status: model.PrintHead, Info: "BBC iPLAYER"}
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

func PrintVN(c *http.Client) model.Result {
	return model.Result{Name: "Vietnam", Status: model.PrintHead, Info: "Galaxy Play"}
}

func PrintMY(c *http.Client) model.Result {
	return model.Result{Name: "Malaysia", Status: model.PrintHead, Info: "Sooka"}
}

func PrintIN(c *http.Client) model.Result {
	return model.Result{Name: "India", Status: model.PrintHead, Info: "Tata Play"}
}

func PrintTW(c *http.Client) model.Result {
	return model.Result{Name: "Taiwan", Status: model.PrintHead, Info: "CatchPlay+"}
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
