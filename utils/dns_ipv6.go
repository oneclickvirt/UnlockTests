package utils

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
)

func lookupIP(ctx context.Context, network, hostname string) ([]net.IP, error) {
	resolver := net.DefaultResolver
	if Dialer.Resolver != nil {
		resolver = Dialer.Resolver
	}
	return resolver.LookupIP(ctx, network, hostname)
}

// CheckIPv6Support 检查域名是否有 AAAA 记录（IPv6 支持）
func CheckIPv6Support(hostname string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 查询 AAAA 记录（IPv6）
	addrs, err := lookupIP(ctx, "ip6", hostname)
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
	addrs, err := lookupIP(ctx, "ip4", hostname)
	if err != nil || len(addrs) == 0 {
		return false
	}
	return true
}

func IsIPv6Client(client interface{}) bool {
	httpClient, ok := client.(*http.Client)
	return ok && httpClient.Transport == Ipv6Transport
}

func extractHostnameFromError(err error) string {
	if err == nil {
		return ""
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		if u, parseErr := url.Parse(urlErr.URL); parseErr == nil {
			return u.Hostname()
		}
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return dnsErr.Name
	}
	errMsg := err.Error()
	for _, pattern := range []string{
		`https?://([^/"\s]+)`,
		`lookup ([^\s:]+)`,
	} {
		re := regexp.MustCompile(pattern)
		if match := re.FindStringSubmatch(errMsg); len(match) > 1 {
			if host := strings.Trim(match[1], "[]"); host != "" {
				if h, _, splitErr := net.SplitHostPort(host); splitErr == nil {
					return h
				}
				return host
			}
		}
	}
	return ""
}

func isDNSResolveError(err error) bool {
	if err == nil {
		return false
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "no such host") ||
		strings.Contains(errMsg, "temporary failure in name resolution") ||
		strings.Contains(errMsg, "name or service not known") ||
		strings.Contains(errMsg, "server misbehaving")
}

func isNoAddressError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "no suitable address") ||
		strings.Contains(errMsg, "no address associated") ||
		strings.Contains(errMsg, "nodename nor servname provided")
}

// HandleNetworkError 智能处理网络错误，在 IPv6 模式下检测是否是因为不支持 IPv6
// client: 当前使用的 HTTP 客户端
// hostname: 要检测的域名
// err: 原始错误
// name: 服务名称
func HandleNetworkError(client interface{}, hostname string, err error, name string) model.Result {
	if hostname == "" {
		hostname = extractHostnameFromError(err)
	}

	if IsIPv6Client(client) {
		if hostname != "" && (isDNSResolveError(err) || isNoAddressError(err)) &&
			!CheckIPv6Support(hostname) && CheckIPv4Support(hostname) {
			return model.Result{Name: name, Status: model.StatusNoIPv6}
		}
	}

	if isDNSResolveError(err) {
		return model.Result{Name: name, Status: model.StatusDNSFailed, Err: err}
	}

	// 返回标准网络错误
	return model.Result{Name: name, Status: model.StatusNetworkErr, Err: err}
}

func NormalizeResult(client interface{}, result model.Result, fallbackName string) model.Result {
	if result.Name == "" {
		result.Name = fallbackName
	}
	if result.Status == model.StatusNetworkErr {
		normalized := HandleNetworkError(client, "", result.Err, result.Name)
		if result.Err == nil {
			normalized.Err = nil
		}
		return normalized
	}
	return result
}
