package utils

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/oneclickvirt/UnlockTests/model"
)

func get_nameserver_from_resolv() []string {
	clientConfig, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil
	}
	return clientConfig.Servers
}

// CheckDNSIP 检测IP地址是否同子网/在内网
func CheckDNSIP(ipStr string, referenceIP string) int {
	// 解析输入的IP地址字符串
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 1 // 如果IP地址无效，返回1
	}
	if ip.To4() != nil {
		// 处理IPv4地址
		privateIPv4Ranges := []string{
			"10.0.0.0/8",
			"172.16.0.0/12",
			"169.254.0.0/16",
			"192.168.0.0/16",
		}
		// 检查IP是否在私有IPv4地址范围内
		for _, cidr := range privateIPv4Ranges {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			if ipNet.Contains(ip) {
				return 0 // 如果IP在私有地址范围内，返回0
			}
		}
		// 检查IP是否与参考IP在同一子网内
		refIP := net.ParseIP(referenceIP)
		if refIP != nil && ip.Mask(net.CIDRMask(24, 32)).Equal(refIP.Mask(net.CIDRMask(24, 32))) {
			return 0 // 如果在同一子网内，返回0
		}
	} else {
		// 处理IPv6地址
		// 检查IP是否在特殊IPv6地址范围内
		if strings.HasPrefix(ipStr, "fe8") || strings.HasPrefix(ipStr, "FE8") ||
			strings.HasPrefix(ipStr, "fc") || strings.HasPrefix(ipStr, "FC") ||
			strings.HasPrefix(ipStr, "fd") || strings.HasPrefix(ipStr, "FD") ||
			strings.HasPrefix(ipStr, "ff") || strings.HasPrefix(ipStr, "FF") {
			return 0 // 如果IP在特殊IPv6地址范围内，返回0
		}
	}
	return 1 // 如果IP不符合上述条件，返回1
}

// // lookupHostWithTimeout 检测网址的IP地址
// func lookupHostWithTimeout(hostname string, timeout time.Duration) ([]string, error) {
// 	// 创建带有超时的上下文
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	// 使用默认解析器查找主机地址
// 	return net.DefaultResolver.LookupHost(ctx, hostname)
// }

// isPossibleCDNIP 检查是否可能是CDN IP
func isPossibleCDNIP(ip string) bool {
	for _, prefix := range model.CdnPrefixes {
		if strings.HasPrefix(ip, prefix) {
			return true
		}
	}
	return false
}

// CheckDNS 三个检测DNS的逻辑并发检测
func CheckDNS(hostname string) (string, string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	var result1, result2, result3 string
	wg.Add(3)

	// 内网/同网IP检测
	go func() {
		defer wg.Done()
		addrs, err := net.DefaultResolver.LookupHost(ctx, hostname)
		if err != nil || len(addrs) == 0 {
			result1 = ""
			return
		}
		result1 = "1"
		for i := 0; i < len(addrs); i++ {
			for j := i + 1; j < len(addrs); j++ {
				if CheckDNSIP(addrs[i], addrs[j]) == 0 {
					result1 = "0"
					return
				}
			}
		}
	}()

	// 主域名DNS解析检测
	go func() {
		defer wg.Done()
		addrs, err := net.DefaultResolver.LookupHost(ctx, hostname)
		if err != nil {
			result2 = ""
			return
		}
		cdnCount := 0
		for _, addr := range addrs {
			if isPossibleCDNIP(addr) {
				cdnCount++
			}
		}
		switch {
		case len(addrs) <= 2:
			result2 = "0" // 可能是原生IP
		case cdnCount > 0:
			result2 = "2" // 可能是CDN
		default:
			result2 = "1" // 多个非CDN的IP，可能是负载均衡
		}
	}()

	// 随机前缀DNS解析检测 - 是否存在通配符DNS记录
	go func() {
		defer wg.Done()
		testDomain := fmt.Sprintf("test%d.%s", rand.Int(), hostname)
		addrs, err := net.DefaultResolver.LookupHost(ctx, testDomain)
		if err != nil || len(addrs) == 0 {
			result3 = "1" // 正常情况
			return
		}
		if len(addrs) > 0 {
			result3 = "0" // 可能存在通配符DNS记录
		}
	}()
	wg.Wait()
	return result1, result2, result3
}

// GetUnlockType 获取解锁的类型
func GetUnlockType(results ...string) string {
	// 检查结果中是否有空值
	for _, result := range results {
		if result == "" {
			return ""
		}
	}
	// 检测是否只有常见的nameserver，此时去判断是否原生解锁无意义
	var status bool = true
	nameservers := get_nameserver_from_resolv()
	if nameservers != nil {
		for _, k := range nameservers {
			// 去除IPV6地址
			if strings.Count(k, ":") > 4 {
				continue
			}
			// 检测非V6地址是不是都是公共DNS
			_, exists := model.CommonPublicDNS[k]
			if !exists {
				status = false
				break
			}
		}
	}
	if status {
		return ""
	}
	// 检查结果中是否有"0"
	for _, result := range results {
		if result == "0" {
			return "Native"
		}
	}
	return "Via DNS"
}
