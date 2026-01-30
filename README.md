# UnlockTests

[![Hits](https://hits.spiritlhl.net/UnlockTests.svg?action=hit&title=Hits&title_bg=%23555555&count_bg=%230eecf8&edge_flat=false)](https://hits.spiritlhl.net)

[![Build and Release](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml/badge.svg)](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml)

解锁测试模块 (Unlock Tests Module)

各种流媒体、AI、直播、论坛、游戏平台的访问解锁测试模块

## 使用

下载、安装、升级

```shell
curl https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/ut_install.sh -sSf | bash
```

或

```
curl https://cdn.spiritlhl.net/https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/ut_install.sh -sSf | bash
```

以后需要使用时使用

```
ut
```

或

```
./ut
```

唤起菜单进行选择

无环境依赖，理论上适配所有系统和主流架构，更多架构请查看 https://github.com/oneclickvirt/UnlockTests/releases/tag/output

```
Usage of ut:
  -I string
        bind IP address or network interface; example: -I 192.168.1.100 or -I eth0
  -L string
        language, specify to en or zh (default "zh")
  -b    use progress bar, disable example: -b=false (default true)
  -cache
        enable caching and sequential region execution; example: -cache
  -conc uint
        max concurrent tests (0=unlimited); example: -conc 50
  -dns-servers string
        specify DNS servers; example: -dns-servers "1.1.1.1:53"
  -f string
        specify select option in menu, example: -f 0
  -http-proxy string
        specify HTTP proxy; example: -http-proxy "http://username:password@127.0.0.1:1080"
  -log
        enable logging
  -m int
        mode 0(both)/4(only)/6(only), default to 0, example: -m 4
  -s    show ip address status, disable example: -s=false (default true)
  -socks-proxy string
        specify SOCKS5 proxy; example: -socks-proxy "socks5://username:password@127.0.0.1:1080"
  -v    show version
```

## 命令行参数详解

### 基本参数

| 参数 | 说明 | 示例 |
|------|------|------|
| `-m` | 连接模式：0=自动（默认），4=仅IPv4，6=仅IPv6 | `-m 4` 仅测试IPv4 |
| `-I` | 绑定的 IP / 网络接口 | `-I 192.168.1.100` 或 `-I eth0` |
| `-v` | 显示版本信息并退出 | `-v` |
| `-L` | 语言选择：zh=中文，en=英文 | `-L en` |
| `-f` | 指定菜单选项 | `-f 0` 跨国平台 |
| `-s` | 显示IP地址状态 | `-s=false` 关闭IP显示 |
| `-b` | 使用进度条 | `-b=false` 关闭进度条 |
| `-log` | 启用日志记录 | `-log` |

### 代理设置

| 参数 | 说明 | 示例 |
|------|------|------|
| `-http-proxy` | 设置 HTTP 代理 | `-http-proxy "http://username:password@127.0.0.1:1080"` |
| `-socks-proxy` | 设置 SOCKS5 代理 | `-socks-proxy "socks5://username:password@127.0.0.1:1080"` |
| `-dns-servers` | 指定 DNS 服务器 | `-dns-servers "1.1.1.1:53"` |

### 性能优化

| 参数 | 说明 | 示例 |
|------|------|------|
| `-conc` | 最大并发测试数量（0=无限制） | `-conc 50` 限制最大50个并发测试 |
| `-cache` | 启用缓存和串行地区执行 | `-cache` 启用缓存模式 |

## 使用示例

### 基本使用

```bash
# 默认检测所有项目
ut

# 仅检测IPv4项目
ut -m 4

# 仅检测IPv6项目
ut -m 6

# 指定菜单选项（跨国平台）
ut -f 0
```

### 代理使用

```bash
# 使用HTTP代理
ut -http-proxy "http://127.0.0.1:8080"

# 使用SOCKS5代理
ut -socks-proxy "socks5://127.0.0.1:1080"

# 使用带认证的代理
ut -http-proxy "http://user:pass@127.0.0.1:8080"

# 指定DNS服务器
ut -dns-servers "1.1.1.1:53"
```

### 性能优化

```bash
# 限制并发数量为30
ut -conc 30

# 启用缓存模式（串行执行，减少网络压力）
ut -cache

# 组合使用：限制并发+缓存
ut -conc 20 -cache
```

### 组合使用

```bash
# 使用代理，仅IPv4，限制并发
ut -http-proxy "http://127.0.0.1:8080" -m 4 -conc 20

# 使用SOCKS5代理，仅IPv6，启用缓存
ut -socks-proxy "socks5://127.0.0.1:1080" -m 6 -cache

# 绑定网卡，使用代理，限制并发
ut -I eth0 -http-proxy "http://127.0.0.1:8080" -conc 30
```

## 卸载

```
rm -rf /root/ut
rm -rf /usr/bin/ut
```

## TODO

<details>

### 无效检测

需要重新构建检测逻辑

```
ElevenSportsTW

CineMax

NPO Start Plus                   Unknown: Token get null

KPLUS - ssoToken 已过期

TV360 - 登录认证已过期

Salto - Get remote error: tls: unrecognized name

PCRJP - stream error: stream ID 1; INTERNAL_ERROR; received from peer
```

### 部分失效

有替代的检测，但仍保留失效检测的部分，未知是否完全失效

```
TikTok - 在 hk、jp 上测试时不时测不出，在 tw 上失效的概率更大，其他地区没有问题

BilibiliID - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

BilibiliTH - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

BilibiliVN - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}
```

### 无需支持

不需要再做支持

```
KBSAmerican - 不再支持本地区

Paravi - 已迁移并集成到 U-NEXT 中。由于整合，除了传统的Paravi作品外，现在还有电影、动漫、亚洲和外国戏剧等等可以无限观看。

HBOGOEurope - api.ugw.hbogo.eu 已经 host 为空了 查询不到内容

HBOSpain - api-discovery.hbo.eu 的 host 已经为空了

HBOGO - 被 HBOMax 替代合并了

HBO_Nordic - 被合并了

HBO_Portugal - 被合并了

PopcornFlix - 已关服

WorldFlipper - 已关服

KonosubaFD - 已关服
```

</details>

## 在Golang中使用

```
go get github.com/oneclickvirt/UnlockTests@v0.0.34-20260130055000
```

## Thanks

https://github.com/nkeonkeo/MediaUnlockTest

https://github.com/HsukqiLee/MediaUnlockTest

https://github.com/lmc999/RegionRestrictionCheck

https://github.com/betaxab/RegionRestrictionCheck/tree/refactor-1
