package executor

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"syscall"
	"time"

	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
	"golang.org/x/net/proxy"
)

var setSocketOptions = func(network, address string, c syscall.RawConn, interfaceName string) (err error) {
	return
}

func SetupInterface(Iface string) {
	ClearCache()
	Iface = strings.TrimSpace(Iface)
	if Iface == "" {
		utils.Dialer.LocalAddr = nil
		utils.Dialer.Control = nil
		resetTransportDialers()
		return
	}
	if IP := net.ParseIP(Iface); IP != nil {
		utils.Dialer.LocalAddr = &net.TCPAddr{IP: IP}
		utils.Dialer.Control = nil
	} else {
		utils.Dialer.LocalAddr = nil
		utils.Dialer.Control = func(network, address string, c syscall.RawConn) error {
			return setSocketOptions(network, address, c, Iface)
		}
	}
	resetTransportDialers()
}

func SetupDnsServers(DnsServers string) {
	ClearCache()
	DnsServers = strings.TrimSpace(DnsServers)
	utils.SetCustomDNSServers(DnsServers)
	dnsDialAddress := firstDNSServerDialAddress(DnsServers)
	if dnsDialAddress == "" {
		utils.Dialer.Resolver = nil
		resetTransportDialers()
		return
	}
	utils.Dialer.Resolver = &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, "udp", dnsDialAddress)
		},
	}
	resetTransportDialers()
}

func firstDNSServerDialAddress(servers string) string {
	fields := strings.FieldsFunc(servers, func(r rune) bool {
		return r == ',' || r == ';' || r == ' '
	})
	for _, server := range fields {
		server = strings.TrimSpace(server)
		if server == "" {
			continue
		}
		if host, port, err := net.SplitHostPort(server); err == nil {
			host = strings.Trim(host, "[]")
			if host != "" && port != "" {
				return net.JoinHostPort(host, port)
			}
			continue
		}
		host := strings.Trim(server, "[]")
		if host != "" {
			return net.JoinHostPort(host, "53")
		}
	}
	return ""
}

func SetupHttpProxy(httpProxy string) {
	ClearCache()
	httpProxy = strings.TrimSpace(httpProxy)
	resetTransportDialers()
	if httpProxy == "" {
		setTransportProxy(http.ProxyFromEnvironment)
		return
	}
	u, err := url.Parse(httpProxy)
	if err != nil || u.Scheme == "" || u.Host == "" {
		fmt.Printf("Warning: HTTP proxy address is invalid: %s\n", httpProxy)
		setTransportProxy(http.ProxyFromEnvironment)
		return
	}
	setTransportProxy(http.ProxyURL(u))
}

func SetupSocksProxy(socksProxy string) {
	ClearCache()
	socksProxy = strings.TrimSpace(socksProxy)
	resetTransportDialers()
	restoreTransportProxy()
	if socksProxy == "" {
		return
	}
	proxyURL, err := url.Parse(socksProxy)
	if err != nil {
		fmt.Printf("Warning: SOCKS5 proxy address is invalid: %v\n", err)
		return
	}
	if proxyURL.Scheme != "" && proxyURL.Scheme != "socks5" && proxyURL.Scheme != "socks5h" {
		fmt.Printf("Warning: SOCKS5 proxy scheme is invalid: %s\n", proxyURL.Scheme)
		return
	}
	if proxyURL.Host == "" {
		fmt.Println("Warning: SOCKS5 proxy host is empty")
		return
	}
	var auth *proxy.Auth
	if proxyURL.User != nil {
		username := proxyURL.User.Username()
		password, _ := proxyURL.User.Password()
		auth = &proxy.Auth{
			User:     username,
			Password: password,
		}
	}
	dialer, err := proxy.SOCKS5("tcp", proxyURL.Host, auth, utils.Dialer)
	if err != nil {
		fmt.Printf("Warning: Failed to create SOCKS5 connection: %v\n", err)
		return
	}
	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		fmt.Println("Warning: SOCKS5 dialer does not support context")
		return
	}

	utils.AutoTransport.DialContext = contextDialer.DialContext

	originalDialContext := contextDialer.DialContext
	utils.Ipv4Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return originalDialContext(ctx, "tcp4", addr)
	}
	utils.Ipv6Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return originalDialContext(ctx, "tcp6", addr)
	}

	clearTransportProxy()
}

func resetTransportDialers() {
	utils.AutoTransport.DialContext = utils.Dialer.DialContext
	utils.Ipv4Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return utils.Dialer.DialContext(ctx, "tcp4", addr)
	}
	utils.Ipv6Transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return utils.Dialer.DialContext(ctx, "tcp6", addr)
	}
}

func setTransportProxy(proxyFunc func(*http.Request) (*url.URL, error)) {
	utils.ClientProxy = proxyFunc
	restoreTransportProxy()
}

func restoreTransportProxy() {
	for _, transport := range []*http.Transport{utils.Ipv4Transport, utils.Ipv6Transport, utils.AutoTransport} {
		transport.Proxy = utils.ClientProxy
	}
}

func clearTransportProxy() {
	for _, transport := range []*http.Transport{utils.Ipv4Transport, utils.Ipv6Transport, utils.AutoTransport} {
		transport.Proxy = nil
	}
}

func SetupConcurrency(conc uint64) {
	if conc > 0 {
		sem = make(chan struct{}, conc)
		return
	}
	sem = nil
}

func EnableCache() {
	cacheEnabled = true
}

func ClearCache() {
	cacheMutex.Lock()
	resultCache = make(map[string]model.Result)
	cacheMutex.Unlock()
}

func maskIP(ip string) string {
	if net.ParseIP(ip).To4() != nil {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			parts[3] = "xxx"
			return strings.Join(parts, ".")
		}
	} else {
		parts := strings.Split(ip, ":")
		if len(parts) > 1 {
			if len(parts[len(parts)-1]) > 0 {
				parts[len(parts)-1] = "xxx"
			} else {
				parts[len(parts)-2] = "xxx"
			}
			return strings.Join(parts, ":")
		}
	}
	return ip
}

func GetIpv4Info(showIP bool) {
	client := utils.Req(utils.Ipv4HttpClient)
	client.SetTimeout(5 * time.Second)
	resp, err := client.R().Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		IPV4 = false
		if showIP {
			fmt.Println("Can not detect IPv4 Address")
		}
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		IPV4 = false
		if showIP {
			fmt.Println("Can not detect IPv4 Address")
		}
		return
	}
	body := string(b)
	if showIP && body != "" && strings.Contains(body, "ip=") {
		s := body
		i := strings.Index(s, "ip=")
		s = s[i+3:]
		i = strings.Index(s, "\n")
		ip := s[:i]
		maskedIP := maskIP(ip)
		fmt.Fprintln(utils.ColorStdout, "Your IPV4 address:", Blue(maskedIP))
	}
}

func GetIpv6Info(showIP bool) {
	client := utils.Req(utils.Ipv6HttpClient)
	client.SetTimeout(5 * time.Second)
	resp, err := client.R().Get("https://www.cloudflare.com/cdn-cgi/trace")
	if err != nil {
		IPV6 = false
		if showIP {
			fmt.Println("Can not detect IPv6 Address")
		}
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		IPV6 = false
		if showIP {
			fmt.Println("Can not detect IPv6 Address")
		}
		return
	}
	body := string(b)
	if showIP && body != "" && strings.Contains(body, "ip=") {
		s := body
		i := strings.Index(s, "ip=")
		s = s[i+3:]
		i = strings.Index(s, "\n")
		ip := s[:i]
		maskedIP := maskIP(ip)
		fmt.Fprintln(utils.ColorStdout, "Your IPV6 address:", Blue(maskedIP))
	}
}
