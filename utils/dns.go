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
	. "github.com/oneclickvirt/defaultset"
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
		// 检查IP是否在私有IPv4地址范围内
		for _, cidr := range model.PrivateIPv4Ranges {
			_, ipNet, err := net.ParseCIDR(cidr)
			if err != nil {
				continue
			}
			if ipNet.Contains(ip) {
				return 2 // 如果IP在私有地址范围内，返回2，与解锁测试判断无关，可能在通过Proxy检测
			}
		}
		// 检查IP是否与参考IP在同一子网内
		refIP := net.ParseIP(referenceIP)
		if refIP != nil && ip.Mask(net.CIDRMask(24, 32)).Equal(refIP.Mask(net.CIDRMask(24, 32))) {
			return 0 // 如果在同一子网内，返回0
		}
		return 1 // 如果IP不符合上述条件，返回1，意味着多数据中心，可能是DNS解锁
	} else {
		// 检查IP是否为 链路本地地址、唯一本地地址和多播地址
		if strings.HasPrefix(ipStr, "fe8") || strings.HasPrefix(ipStr, "FE8") ||
			strings.HasPrefix(ipStr, "fc") || strings.HasPrefix(ipStr, "FC") ||
			strings.HasPrefix(ipStr, "fd") || strings.HasPrefix(ipStr, "FD") ||
			strings.HasPrefix(ipStr, "ff") || strings.HasPrefix(ipStr, "FF") {
			return 2 // 可能在Proxy中
		}
		return 1
	}
}

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
		totalChecks := 0
		sameSubnetOrPrivateCount := 0
		for i := 0; i < len(addrs); i++ {
			for j := i + 1; j < len(addrs); j++ {
				totalChecks++
				if CheckDNSIP(addrs[i], addrs[j]) == 0 {
					sameSubnetOrPrivateCount++
				}
			}
		}
		if totalChecks > 0 && sameSubnetOrPrivateCount > totalChecks/2 {
			result1 = "0" // 大多数IP在同一子网或是内网IP
		} else {
			result1 = "1" // 大多数IP不在同一子网且不是内网IP，可能是DNS解锁
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
		// 根据解析结果进行判断
		switch {
		case len(addrs) <= 2:
			result2 = "0" // 解析到2个或更少的IP，可能是原生IP - 大多数原生服务通常只有少量IP
		case cdnCount > 0:
			result2 = "1" // 检测到至少一个可能的CDN的IP - CDN的使用通常与DNS解锁相关，而不是原生解锁
		default:
			result2 = "0" // 解析到多个非CDN的IP，可能是使用负载均衡的原生解锁 - 多个非CDN的IP可能表示服务提供商使用了自己的负载均衡系统
		}
	}()
	// 随机前缀DNS解析检测 - 是否存在通配符DNS记录
	go func() {
		defer wg.Done()
		testDomain := fmt.Sprintf("test%d.%s", rand.Int(), hostname)
		addrs, err := net.DefaultResolver.LookupHost(ctx, testDomain)
		if err != nil || len(addrs) == 0 {
			result3 = "0" // 正常情况，不通配
			return
		}
		if len(addrs) > 0 {
			result3 = "1" // 可能存在通配符DNS记录，可能是DNS解锁
		}
	}()
	wg.Wait()
	return result1, result2, result3
}

// GetUnlockType 获取解锁的类型
func GetUnlockType(results ...string) string {
	if model.EnableLoger {
		InitLogger()
		defer Logger.Sync()
	}
	// 检查结果中是否有空值
	for _, result := range results {
		if result == "" {
			return ""
		}
	}
	// 检测是否只有常见的nameserver，此时去判断是否原生解锁无意义
	// 识别不出nameserver时，不做是否DNS解锁的判断
	var status bool = true
	nameservers := get_nameserver_from_resolv()
	if nameservers != nil {
		if model.EnableLoger {
			Logger.Info("Name servers: ")
		}
		for _, k := range nameservers {
			if model.EnableLoger {
				Logger.Info(k)
			}
			ip := net.ParseIP(strings.TrimSpace(k))
			if ip == nil {
				// 无效的 IP 地址跳过检测
				continue
			}
			if ip.To4() != nil {
				// 检测非V6地址是不是都是公共DNS
				_, exists := model.CommonPublicDNS[k]
				if !exists {
					status = false
					break
				}
			} else {
				// 去除IPV6地址的检测
				continue
			}
		}
	} else {
		return ""
	}
	if status {
		return ""
	}
	// 检查结果中是原生解锁的判断为多数
	zeroCount := 0
	for _, result := range results {
		if result == "2" {
			return "In Proxy"
		}
		if result == "0" {
			zeroCount++
		}
	}
	if zeroCount >= 2 {
		return "Native"
	}
	return "Via DNS"
}

// // lookupHostWithTimeout 检测网址的IP地址
// func lookupHostWithTimeout(hostname string, timeout time.Duration) ([]string, error) {
// 	// 创建带有超时的上下文
// 	ctx, cancel := context.WithTimeout(context.Background(), timeout)
// 	defer cancel()
// 	// 使用默认解析器查找主机地址
// 	return net.DefaultResolver.LookupHost(ctx, hostname)
// }
