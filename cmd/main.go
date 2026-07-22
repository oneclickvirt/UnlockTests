package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	req "github.com/imroc/req/v3"
	"github.com/oneclickvirt/UnlockTests/executor"
	"github.com/oneclickvirt/UnlockTests/model"
	"github.com/oneclickvirt/UnlockTests/utils"
	. "github.com/oneclickvirt/defaultset"
)

type cliOptions struct {
	mode                                          int
	showVersion, help, showIP, useBar, cache, log bool
	jsonOutput                                    bool
	iface, dnsServers, httpProxy, socksProxy      string
	language, selection, testNames                string
	concurrency                                   uint64
	timeout                                       time.Duration
}

type cliStructuredOutput struct {
	SchemaVersion    string                          `json:"schema_version"`
	Results          []executor.StructuredResult     `json:"results"`
	ProviderMetadata []executor.ProviderMetadata     `json:"provider_metadata"`
	MetadataSource   executor.ProviderMetadataSource `json:"provider_metadata_source"`
	Error            string                          `json:"error,omitempty"`
}

func parseCLI(args []string) (cliOptions, error) {
	opts := cliOptions{}
	fs := newFlagSet(&opts, io.Discard)
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	if opts.timeout < 0 {
		return opts, fmt.Errorf("timeout must not be negative")
	}
	return opts, nil
}

func newFlagSet(opts *cliOptions, output io.Writer) *flag.FlagSet {
	fs := flag.NewFlagSet("ut", flag.ContinueOnError)
	fs.SetOutput(output)
	fs.BoolVar(&opts.help, "h", false, "show help information")
	fs.IntVar(&opts.mode, "m", 0, "mode: 0 (both), 4 (only), or 6 (only); default is 0, example: -m 4")
	fs.BoolVar(&opts.showVersion, "v", false, "show version")
	fs.BoolVar(&opts.showIP, "s", true, "show IP address status; to disable, use: -s=false")
	fs.BoolVar(&opts.useBar, "b", true, "use progress bar; to disable, use: -b=false")
	fs.BoolVar(&opts.log, "log", false, "enable logging")
	fs.StringVar(&opts.selection, "f", "", "specify selection option in menu; example: -f 0")
	fs.StringVar(&opts.testNames, "test", "", "run specific providers by name or function, comma-separated; example: -test \"Coze,Poe\"")
	fs.StringVar(&opts.iface, "I", "", "bind IP address or network interface; example: -I 192.168.1.100 or -I eth0")
	fs.StringVar(&opts.dnsServers, "dns-servers", "", "specify DNS servers; example: -dns-servers \"1.1.1.1:53\"")
	fs.StringVar(&opts.httpProxy, "http-proxy", "", "specify HTTP proxy; example: -http-proxy \"http://username:password@127.0.0.1:1080\"")
	fs.StringVar(&opts.socksProxy, "socks-proxy", "", "specify SOCKS5 proxy; example: -socks-proxy \"socks5://username:password@127.0.0.1:1080\"")
	fs.Uint64Var(&opts.concurrency, "conc", 0, "max concurrent tests (0=unlimited); example: -conc 50")
	fs.BoolVar(&opts.cache, "cache", false, "enable duplicate test result caching; example: -cache")
	fs.StringVar(&opts.language, "L", "zh", "language; specify 'en' for English or 'zh' for Chinese")
	fs.BoolVar(&opts.jsonOutput, "json", false, "print structured provider results as JSON")
	fs.BoolVar(&opts.jsonOutput, "structured", false, "print structured provider results as JSON")
	fs.DurationVar(&opts.timeout, "timeout", 0, "structured run timeout (for example 2m)")
	return fs
}

func main() {
	opts, err := parseCLI(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, sanitizeErrorText(err.Error()))
		os.Exit(2)
	}
	model.EnableLoger = opts.log
	if opts.help {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		newFlagSet(&cliOptions{}, os.Stdout).PrintDefaults()
		return
	}
	if opts.showVersion {
		fmt.Println(model.UnlockTestsVersion)
		return
	}
	mode, showIP, useBar, cache := opts.mode, opts.showIP, opts.useBar, opts.cache
	Iface, DnsServers, httpProxy, socksProxy := opts.iface, opts.dnsServers, opts.httpProxy, opts.socksProxy
	language, flagString, testString, conc := opts.language, opts.selection, opts.testNames, opts.concurrency
	if mode != 0 && mode != 4 && mode != 6 {
		fmt.Fprintf(os.Stderr, "invalid mode: %d; expected 0, 4, or 6\n", mode)
		os.Exit(2)
	}
	if language != "zh" && language != "en" {
		fmt.Fprintf(os.Stderr, "invalid language: %s; expected zh or en\n", language)
		os.Exit(2)
	}
	if opts.jsonOutput {
		output := cliStructuredOutput{SchemaVersion: "goecs.unlocktests/v1", Results: []executor.StructuredResult{}}
		if conc > uint64(^uint(0)>>1) {
			output.Error = "concurrency exceeds the platform integer range"
			encoded, _ := json.Marshal(output)
			fmt.Println(string(encoded))
			os.Exit(2)
		}
		ipVersion := "auto"
		if mode == 4 {
			ipVersion = "ipv4"
		} else if mode == 6 {
			ipVersion = "ipv6"
		}
		timeout := opts.timeout
		if timeout <= 0 {
			timeout = 2 * time.Minute
		}
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		metadata, metadataSource, metadataErr := executor.LoadProviderMetadata(ctx, nil)
		if metadataErr != nil {
			output.Error = metadataErr.Error()
			encoded, _ := json.Marshal(output)
			fmt.Println(string(encoded))
			os.Exit(1)
		}
		output.ProviderMetadata = metadata
		output.MetadataSource = metadataSource
		runOptions := executor.RunOptions{
			Selection: flagString, IPVersion: ipVersion, Concurrency: int(conc), Interface: Iface,
			DNSServers: DnsServers, HTTPProxy: httpProxy, SOCKSProxy: socksProxy, UseCache: cache,
		}
		var results []executor.StructuredResult
		var runErr error
		if testString != "" {
			results, runErr = executor.RunNamedStructured(ctx, runOptions, testString)
		} else {
			results, runErr = executor.RunStructured(ctx, runOptions)
		}
		output.Results = results
		if runErr != nil {
			output.Error = runErr.Error()
		}
		encoded, marshalErr := json.Marshal(output)
		if marshalErr != nil {
			fmt.Fprintln(os.Stderr, marshalErr)
			return
		}
		fmt.Println(string(encoded))
		if runErr != nil {
			os.Exit(1)
		}
		return
	}
	if Iface != "" {
		if err := executor.SetupInterface(Iface); err != nil {
			fmt.Fprintln(os.Stderr, sanitizeErrorText(err.Error()))
			return
		}
	}
	if DnsServers != "" {
		executor.SetupDnsServers(DnsServers)
	}
	if httpProxy != "" {
		executor.SetupHttpProxy(httpProxy)
	}
	if socksProxy != "" {
		executor.SetupSocksProxy(socksProxy)
	}
	if conc > 0 {
		executor.SetupConcurrency(conc)
	}
	if cache {
		executor.EnableCache()
	}
	if mode == 4 {
		executor.IPV6 = false
	}
	if mode == 6 {
		executor.IPV4 = false
	}
	if language == "zh" {
		fmt.Fprintln(utils.ColorStdout, "项目地址: "+Blue("https://github.com/oneclickvirt/UnlockTests"))
	} else {
		fmt.Fprintln(utils.ColorStdout, "Github Repo: "+Blue("https://github.com/oneclickvirt/UnlockTests"))
	}
	if testString == "" {
		readStatus := executor.ReadSelect(language, flagString)
		if !readStatus {
			return
		}
	}
	trackHit()
	if executor.IPV4 {
		executor.GetIpv4Info(showIP)
	}
	if executor.IPV6 {
		executor.GetIpv6Info(showIP)
	}
	if language == "zh" {
		fmt.Fprintln(utils.ColorStdout, "测试时间: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	} else {
		fmt.Fprintln(utils.ColorStdout, "Test time: ", Yellow(time.Now().Format("2006-01-02 15:04:05")))
	}
	if executor.IPV4 {
		if testString != "" {
			result, err := executor.RunNamedTests(utils.Ipv4HttpClient, "ipv4", language, useBar, testString)
			if err != nil {
				fmt.Fprintln(os.Stderr, sanitizeErrorText(err.Error()))
				return
			}
			fmt.Fprint(utils.ColorStdout, indentLegacyOutput(result))
		} else {
			fmt.Fprint(utils.ColorStdout, indentLegacyOutput(executor.RunTests(utils.Ipv4HttpClient, "ipv4", language, useBar)))
		}
	}
	if executor.IPV6 {
		if testString != "" {
			result, err := executor.RunNamedTests(utils.Ipv6HttpClient, "ipv6", language, useBar, testString)
			if err != nil {
				fmt.Fprintln(os.Stderr, sanitizeErrorText(err.Error()))
				return
			}
			fmt.Fprint(utils.ColorStdout, indentLegacyOutput(result))
		} else {
			fmt.Fprint(utils.ColorStdout, indentLegacyOutput(executor.RunTests(utils.Ipv6HttpClient, "ipv6", language, useBar)))
		}
	}
	if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
		fmt.Println("Press Enter to exit...")
		fmt.Scanln()
	}
}

func trackHit() {
	go func() {
		client := req.C()
		client.SetTimeout(2 * time.Second)
		_, _ = client.R().Get("https://hits.spiritlhl.net/UnlockTests.svg?action=hit&title=Hits&title_bg=%23555555&count_bg=%230eecf8&edge_flat=false")
	}()
}
