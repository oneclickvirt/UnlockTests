package utils

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
)

var lookupIP = func(ctx context.Context, network, hostname string) ([]net.IP, error) {
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
	if !ok {
		return getDNSIPVersion() == "ipv6"
	}
	if httpClient.Transport == Ipv6Transport {
		return true
	}
	if httpClient.Transport == Ipv4Transport {
		return false
	}
	return getDNSIPVersion() == "ipv6"
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

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "timed out") ||
		strings.Contains(errMsg, "deadline exceeded")
}

func IsWAFStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusForbidden, http.StatusNotAcceptable, http.StatusRequestTimeout,
		http.StatusMisdirectedRequest, http.StatusTooManyRequests, http.StatusServiceUnavailable:
		return true
	}
	return statusCode >= 520 && statusCode <= 527
}

func IsUnavailableStatusCode(statusCode int) bool {
	switch statusCode {
	case http.StatusBadRequest, http.StatusNotFound, http.StatusUnavailableForLegalReasons, 452:
		return true
	}
	return false
}

func StatusCodeFromError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	re := regexp.MustCompile(`(?i)(?:code|status)\D{0,12}(\d{3})`)
	matches := re.FindStringSubmatch(err.Error())
	if len(matches) < 2 {
		return 0, false
	}
	statusCode, convErr := strconv.Atoi(matches[1])
	if convErr != nil {
		return 0, false
	}
	return statusCode, true
}

func WAFStatusCodeFromError(err error) (int, bool) {
	if err == nil {
		return 0, false
	}
	statusCode, ok := StatusCodeFromError(err)
	return statusCode, ok && IsWAFStatusCode(statusCode)
}

func IsWAFBlockError(err error) bool {
	if err == nil {
		return false
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return false
	}
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) ||
		errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) ||
		errors.Is(err, syscall.ECONNRESET) || errors.Is(err, syscall.ECONNABORTED) {
		return true
	}
	if _, ok := WAFStatusCodeFromError(err); ok {
		return true
	}
	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "stream error") ||
		strings.Contains(errMsg, "handshake failure") ||
		strings.Contains(errMsg, "connection reset") ||
		strings.Contains(errMsg, "connection was forcibly closed") ||
		strings.Contains(errMsg, "forcibly closed by the remote host") ||
		strings.Contains(errMsg, "unexpected eof") ||
		(strings.Contains(errMsg, "tls:") && strings.Contains(errMsg, "handshake")) ||
		strings.Contains(errMsg, "client.timeout exceeded")
}

func isBannedInfo(info string) bool {
	info = strings.ToLower(strings.TrimSpace(info))
	if info == "" {
		return false
	}
	if strings.Contains(info, "banned") {
		return true
	}
	if strings.Contains(info, "cloudflare") && (strings.Contains(info, "blocked") || strings.Contains(info, "forbidden")) {
		return true
	}
	return strings.Contains(info, "request blocked")
}

func isRateLimitInfo(info string) bool {
	info = strings.ToLower(strings.TrimSpace(info))
	return strings.Contains(info, "429") && strings.Contains(info, "rate limit")
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

	if IsWAFBlockError(err) {
		return model.Result{Name: name, Status: model.StatusBanned, Err: err}
	}

	if isTimeoutError(err) {
		return model.Result{Name: name, Status: model.StatusTimeout, Err: err}
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
	if result.Name == "" {
		result.Name = "Unknown"
	}
	if result.Status == model.StatusNetworkErr {
		normalized := HandleNetworkError(client, "", result.Err, result.Name)
		if result.Err == nil {
			normalized.Err = nil
		}
		return normalized
	}
	if result.Status == model.StatusNo && isRateLimitInfo(result.Info) {
		result.Status = model.StatusRestricted
		return result
	}
	if (result.Status == model.StatusNo || result.Status == model.StatusUnexpected) && isBannedInfo(result.Info) {
		result.Status = model.StatusBanned
		return result
	}
	if result.Status == model.StatusUnexpected && result.Err != nil && IsWAFBlockError(result.Err) {
		return model.Result{
			Name:       result.Name,
			Status:     model.StatusBanned,
			Err:        result.Err,
			Region:     result.Region,
			Info:       result.Info,
			UnlockType: result.UnlockType,
		}
	}
	if result.Status == model.StatusUnexpected && result.Err != nil {
		if statusCode, ok := StatusCodeFromError(result.Err); ok && IsUnavailableStatusCode(statusCode) {
			result.Status = model.StatusNo
			return result
		}
		if strings.Contains(strings.ToLower(result.Err.Error()), "token get null") {
			result.Status = model.StatusNo
			return result
		}
	}
	if result.Status == "" {
		result.Status = model.StatusUnexpected
	}
	return result
}
