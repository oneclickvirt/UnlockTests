# UnlockTests

[![Hits](https://hits.spiritlhl.net/UnlockTests.svg?action=hit&title=Hits&title_bg=%23555555&count_bg=%230eecf8&edge_flat=false)](https://hits.spiritlhl.net)

[![Build and Release](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml/badge.svg)](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml)

解锁测试模块 (Unlock Tests Module)

用于检测多类流媒体、AI、直播、论坛、游戏平台在当前网络出口下的访问和区域解锁状态。

## 安装

```shell
curl https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/ut_install.sh -sSf | bash
```

或使用 CDN 入口：

```shell
curl https://cdn.spiritlhl.net/https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/ut_install.sh -sSf | bash
```

安装脚本当前支持 Linux、macOS(Darwin)、FreeBSD、OpenBSD 的主流架构。更多构建产物可查看 release 输出：

https://github.com/oneclickvirt/UnlockTests/releases/tag/output

## 使用

交互式选择检测项目：

```shell
ut
```

在当前目录直接运行二进制：

```shell
./ut
```

查看帮助：

```shell
ut -h
```

```text
Usage: ut [options]
  -I string
        bind IP address or network interface; example: -I 192.168.1.100 or -I eth0
  -L string
        language; specify 'en' for English or 'zh' for Chinese (default "zh")
  -b    use progress bar; to disable, use: -b=false (default true)
  -cache
        enable duplicate test result caching; example: -cache
  -conc uint
        max concurrent tests (0=unlimited); example: -conc 50
  -dns-servers string
        specify DNS servers; example: -dns-servers "1.1.1.1:53"
  -f string
        specify selection option in menu; example: -f 0
  -h    show help information
  -http-proxy string
        specify HTTP proxy; example: -http-proxy "http://username:password@127.0.0.1:1080"
  -log
        enable logging
  -m int
        mode: 0 (both), 4 (only), or 6 (only); default is 0, example: -m 4
  -s    show IP address status; to disable, use: -s=false (default true)
  -socks-proxy string
        specify SOCKS5 proxy; example: -socks-proxy "socks5://username:password@127.0.0.1:1080"
  -test string
        run specific providers by name or function, comma-separated; example: -test "Coze,Poe"
  -v    show version
```

## 检测项目选择

启动后会显示菜单。也可以用 `-f` 直接指定菜单编号，多个编号用空格分隔并加引号。

| 编号 | 检测范围 |
|------|----------|
| `0` | 跨国平台 |
| `1` - `9` | 跨国平台 + 指定地区平台 |
| `10` - `18` | 仅指定地区平台 |
| `19` | 仅体育平台 |
| `20` | 全部平台 |
| `21` | 仅 AI 平台 |

地区平台包括台湾、香港、日本、韩国、北美、南美、欧洲、非洲、大洋洲。

单独运行指定平台可使用 `-test`，多个平台使用半角逗号分隔，例如：

```shell
ut -test "Coze,Poe"
```

## 命令行参数

### 基本参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-m` | 连接模式：0=IPv4 和 IPv6，4=仅 IPv4，6=仅 IPv6 | `ut -m 4` |
| `-I` | 绑定本机 IP 地址或网络接口 | `ut -I 192.168.1.100` 或 `ut -I eth0` |
| `-v` | 显示版本信息并退出 | `ut -v` |
| `-L` | 输出语言：`zh` 或 `en` | `ut -L en` |
| `-f` | 指定菜单编号，多个编号用空格分隔 | `ut -f "0 10"` |
| `-s` | 是否显示本机出口 IP 状态 | `ut -s=false` |
| `-b` | 是否使用进度条 | `ut -b=false` |
| `-log` | 启用日志记录 | `ut -log` |

### 代理和 DNS

| 参数 | 说明 | 示例 |
|------|------|------|
| `-http-proxy` | 设置 HTTP 代理 | `ut -http-proxy "http://127.0.0.1:8080"` |
| `-socks-proxy` | 设置 SOCKS5 代理 | `ut -socks-proxy "socks5://127.0.0.1:1080"` |
| `-dns-servers` | 指定 DNS 服务器 | `ut -dns-servers "1.1.1.1:53"` |

### 执行控制

| 参数 | 说明 | 示例 |
|------|------|------|
| `-conc` | 限制最大并发检测数量，`0` 表示不额外限制 | `ut -conc 50` |
| `-cache` | 启用同名检测结果缓存。同一进程内重复执行时可复用结果；组合或全平台检测默认会在发起请求前按名称去重 | `ut -cache` |

## 环境变量

以下变量用于覆盖源码内置的历史默认参数；未设置时仍沿用默认值以保持旧版行为。部分默认值属于临时授权或签名参数，服务方变更或过期后可通过环境变量替换。

| 变量 | 说明 |
|------|------|
| `UNLOCKTESTS_KPLUS_SSO_TOKEN` | 覆盖 K+ 检测所需的 SSO token；未设置时沿用内置默认值。 |
| `UNLOCKTESTS_TV360_AUTHORIZATION` | 覆盖 TV360 检测所需的 Bearer token；可填写完整 `Bearer ...` 或仅填写 token。 |
| `UNLOCKTESTS_VIDEOMARKET_AUTHORIZATION` | 覆盖 VideoMarket 检测所需的 Bearer token；可填写完整 `Bearer ...` 或仅填写 token。 |
| `UNLOCKTESTS_NLZIET_AUTHORIZATION` | 覆盖 NLZIET 认证检测所需的 Bearer token；认证检测失败时仍会回退到 CDN 地区检测。 |
| `UNLOCKTESTS_KPLUS_PROVISION_DATA` | 覆盖 K+ 检测所需的 provisionData。 |
| `UNLOCKTESTS_NAVERTV_KEY` | 覆盖 Naver TV 请求签名密钥。 |
| `UNLOCKTESTS_LEMINO_SERVICE_TOKEN` | 覆盖 Lemino 服务 token。 |
| `UNLOCKTESTS_SKYGO_NZ_AUTH_URL` | 覆盖 SkyGo NZ fallback 授权 URL。 |
| `UNLOCKTESTS_NEONTV_AUTHORIZATION` | 覆盖 Neon TV 检测所需的 Bearer token；可填写完整 `Bearer ...` 或仅填写 token。 |
| `UNLOCKTESTS_CHANNEL5_AUTH` | 覆盖 Channel 5 live media auth 参数。 |
| `UNLOCKTESTS_MGMPLUS_API_KEY` | 覆盖 MGM+ 会话 API key。 |
| `UNLOCKTESTS_ESPN_SUBJECT_TOKEN` | 覆盖 ESPN+ token exchange 的 subject token。 |
| `UNLOCKTESTS_ESPN_AUTHORIZATION` | 覆盖 ESPN+ API authorization；可填写完整 `Bearer ...` 或仅填写 token。 |
| `UNLOCKTESTS_LITV_SIGNATURE_KEY` | 覆盖 LiTV 旧签名检测 fallback 的签名密钥；默认优先使用新版 RPC 检测。 |
| `UNLOCKTESTS_HBOMAX_VPN_CHECK_ST` | 覆盖 HBO Max VPN 二次校验 token。 |
| `UNLOCKTESTS_AETV_VIDEO_META_TOKEN` | 覆盖 A&E TV video meta token。 |
| `UNLOCKTESTS_DISNEY_AUTHORIZATION` | 覆盖 Disney+ 设备和 token exchange authorization；可填写完整 `Bearer ...` 或仅填写 token。 |
| `UNLOCKTESTS_STARPLUS_AUTHORIZATION` | 覆盖 Star+ CDN Relay 二次检测 authorization。 |
| `UNLOCKTESTS_JOYN_API_KEY` | 覆盖 Joyn entitlement token API key。 |
| `UNLOCKTESTS_SETANTA_API_KEY` | 覆盖 Setanta Sports consent API key。 |
| `UNLOCKTESTS_SDGGGE_TOKEN` | 覆盖 SD Gundam G 检测所需 `x-token`。 |
| `UNLOCKTESTS_CATCHPLAY_AUTHORIZATION` | 覆盖 CatchPlay+ geo API authorization。 |

## 示例

```bash
# 交互式选择检测项目
ut

# 直接检测跨国平台
ut -f 0

# 检测跨国平台和台湾平台
ut -f "0 10"

# 仅检测 IPv4
ut -m 4 -f 0

# 仅检测 IPv6
ut -m 6 -f 0

# 使用 HTTP 代理并限制并发
ut -http-proxy "http://127.0.0.1:8080" -conc 20 -f 20

# 使用 SOCKS5 代理、指定 DNS，并关闭进度条
ut -socks-proxy "socks5://127.0.0.1:1080" -dns-servers "1.1.1.1:53" -b=false -f 0

# 绑定网卡
ut -I eth0 -f 0
```

## 输出说明

常见结果状态：

| 状态 | 含义 |
|------|------|
| `YES` | 可访问或可解锁 |
| `NO` | 不可解锁 |
| `Restricted` | 仅部分内容可用 |
| `Banned` | 当前出口被服务方封禁或限制 |
| `TIMEOUT` | 单个平台检测超时 |
| `Failed (Network Error)` | 网络连接错误 |
| `N/A (DNS Resolve Failed)` | DNS 解析失败 |
| `N/A (No IPv6 Support)` | 当前 IPv6 检测中，目标域名无 IPv6 支持 |
| `CDN Relay Available` | 检测到 CDN Relay 可用 |
| `Unknown` | 响应不符合已知判断逻辑 |

部分支持区域判断的项目会显示 `Region`，部分项目会根据 DNS 检测结果标注 `Native`、`Via DNS` 或 `In Proxy`。

## 卸载

```shell
rm -f /usr/bin/ut
rm -f ./ut
```

## 在 Go 中使用

```shell
go get github.com/oneclickvirt/UnlockTests@v0.0.40-20260705091135
```

结构化接口适合在 goecs 等项目中直接调用，支持 `context.Context`、菜单编号选择、IPv4/IPv6 独立检测和并发上限：

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/oneclickvirt/UnlockTests/executor"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	results, err := executor.RunStructured(ctx, executor.RunOptions{
		Selection:   "0 10",
		IPVersion:   "ipv6",
		Concurrency: 20,
		UseCache:    true,
	})
	if err != nil {
		fmt.Println("partial result:", err)
	}
	for _, result := range results {
		fmt.Printf("%s %s %s %s\n", result.Name, result.Status, result.Region, result.Info)
	}
}
```

也可以用 `executor.ListPlatforms("0")` 先获取指定菜单编号会执行的平台名列表。结构化结果包含 `Name`、`Status`、`Region`、`Info`、`UnlockType` 和 `Error` 字段；旧的 `executor.RunTests` 字符串输出接口保持兼容。

## 开发与测试

默认单元测试不访问真实流媒体站点：

```shell
go test ./...
```

各地区目录中的真实解锁探测测试使用 `live` build tag 隔离，运行时会访问外部服务：

```shell
go test -tags live ./...
```

## Thanks

https://github.com/nkeonkeo/MediaUnlockTest

https://github.com/HsukqiLee/MediaUnlockTest

https://github.com/lmc999/RegionRestrictionCheck

https://github.com/betaxab/RegionRestrictionCheck/tree/refactor-1
