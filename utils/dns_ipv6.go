package utils

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
)

// CheckIPv6Support 检查域名是否有 AAAA 记录（IPv6 支持）
func CheckIPv6Support(hostname string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 查询 AAAA 记录（IPv6）
	addrs, err := net.DefaultResolver.LookupIP(ctx, "ip6", hostname)
	if err != nil || len(addrs) == 0 {
		return false
	}
	return true
}

// CheckIPv4Support 检查域名是否有 A 记录（IPv4 支持）
func CheckIPv4Support(hostname string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// 查询 A 记录（IPv4）
	addrs, err := net.DefaultResolver.LookupIP(ctx, "ip4", hostname)
	if err != nil || len(addrs) == 0 {
		return false
	}
	return true
}

// HandleNetworkError 智能处理网络错误，在 IPv6 模式下检测是否是因为不支持 IPv6
// client: 当前使用的 HTTP 客户端
// hostname: 要检测的域名
// err: 原始错误
// name: 服务名称
func HandleNetworkError(client interface{}, hostname string, err error, name string) model.Result {
	// 检查是否在使用 IPv6 客户端
	isIPv6Client := false
	if httpClient, ok := client.(*http.Client); ok {
		if httpClient.Transport == Ipv6Transport {
			isIPv6Client = true
		}
	}
	
	if isIPv6Client {
		// 在 IPv6 模式下，检查是否是因为没有 AAAA 记录
		if !CheckIPv6Support(hostname) {
			return model.Result{Name: name, Status: model.StatusNoIPv6}
		}
		
		// 检查错误信息，判断是否是 DNS 解析失败
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "no such host") || 
			   strings.Contains(errMsg, "Temporary failure in name resolution") ||
			   strings.Contains(errMsg, "Name or service not known") {
				return model.Result{Name: name, Status: model.StatusDNSFailed, Err: err}
			}
		}
	}
	
	// 返回标准网络错误
	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
}
