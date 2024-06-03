package utils

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/parnurzeal/gorequest"
)

// ParseInterface 解析网卡IP地址
func ParseInterface(ifaceName, ipAddr, netType string) (*gorequest.SuperAgent, error) {
	var localIP net.IP
	var request *gorequest.SuperAgent
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
	request = gorequest.New()
	defaultTransport := http.DefaultTransport.(*http.Transport)
	customTransport := defaultTransport.Clone()
	if localIP != nil {
		customTransport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{
				Timeout:   20 * time.Second,
				KeepAlive: 20 * time.Second,
				LocalAddr: &net.TCPAddr{
					IP: localIP,
				},
			}).DialContext(ctx, netType, addr)
		}
		request.Client.Transport = customTransport
	} else {
		customTransport.DialContext = func(ctx context.Context, network string, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, netType, addr)
		}
		request.Client.Transport = customTransport
	}
	request.Timeout(20 * time.Second)
	return request, nil
}

// PostJson 向指定的 URL 发送 JSON 格式的 POST 请求，并返回响应、响应体和错误信息
// request: gorequest.SuperAgent 实例，用于构建请求
// url: 目标 URL
// payload: 要发送的 JSON 格式的请求体
// headers: 可选的 HTTP 头信息
func PostJson(request *gorequest.SuperAgent, url string, payload string, headers ...map[string]string) (gorequest.Response, []byte, []error) {
	// 构建 POST 请求，设置请求类型为 JSON 并添加请求体
	req := request.Post(url).
		Type("json").
		Send(payload)
	// 添加可选的 HTTP 头信息
	for _, header := range headers {
		for k, v := range header {
			req = req.Set(k, v)
		}
	}
	// 发送请求并接收响应、响应体和错误信息
	resp, body, errs := req.EndBytes()
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

// 通过Info标记要被插入的行的下一行包含什么文本内容
func PrintCA(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Canada", Status: model.PrintHead, Info: "Hotstar"}
}

func PrintGB(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "England", Status: model.PrintHead, Info: "Hotstar"}
}

func PrintFR(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "France", Status: model.PrintHead, Info: "Canal+"}
}

func PrintDE(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Germany", Status: model.PrintHead, Info: "Joyn"}
}

func PrintNL(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Netherlands", Status: model.PrintHead, Info: "NLZIET"}
}

func PrintES(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Spain", Status: model.PrintHead, Info: "Movistar+"}
}

func PrintIT(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Italy", Status: model.PrintHead, Info: "Rai Play"}
}

func PrintCH(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Switzerland", Status: model.PrintHead, Info: "SKY CH"}
}

func PrintRU(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Russia", Status: model.PrintHead, Info: "Amediateka"}
}

func PrintAU(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Australia", Status: model.PrintHead, Info: "Stan"}
}

func PrintNZ(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "New Zealand", Status: model.PrintHead, Info: "Neon TV"}
}

func PrintSG(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Singapore", Status: model.PrintHead, Info: "MeWatch"}
}

func PrintTH(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Thailand", Status: model.PrintHead, Info: "AIS Play"}
}

func PrintGame(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Game", Status: model.PrintHead, Info: "Kancolle Japan"}
}

func PrintMusic(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Music", Status: model.PrintHead, Info: "Mora"}
}

func PrintForum(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Forum", Status: model.PrintHead, Info: "EroGameSpace"}
}

func PrintRForum(request *gorequest.SuperAgent) model.Result {
	return model.Result{Name: "Forum", Status: model.PrintHead, Info: "Reddit"}
}
