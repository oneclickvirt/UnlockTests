# UnlockTests

[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Foneclickvirt%2FUnlockTests&count_bg=%2323E01C&title_bg=%23555555&icon=sonarcloud.svg&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com) [![Build and Release](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml/badge.svg)](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml)

解锁测试模块 (Unlock Tests Module)

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
        specify source ip / interface
  -L string
        language, specify to en or zh (default "zh")
  -b    use progress bar, disable example: -b=false (default true)
  -dns-servers string
        specify dns servers
  -f string
        specify select option in menu, example: -f 0
  -http-proxy string
        specify http proxy
  -m int
        mode 0(both)/4(only)/6(only), default to 0, example: -m 4
  -s    show ip address status, disable example: -s=false (default true)
  -v    show version
```

## 卸载

```
rm -rf /root/ut
rm -rf /usr/bin/ut
```

## TODO

<details>

### 同态检测

可能需要拆分检测

```
GYAO 和 LINE VOOM
```

### 无效检测

需要重新构建检测逻辑

```
ElevenSportsTW

MegogoTV

CineMax

KPLUS - ssoToken 已过期

TV360 - 登录认证已过期

Salto - Get remote error: tls: unrecognized name

Catchplay - unauthorized 原 token 已过期

PCRJP - stream error: stream ID 1; INTERNAL_ERROR; received from peer

PrettyDerby - stream error: stream ID 1; INTERNAL_ERROR; received from peer

WorldFlipper - stream error: stream ID 1; INTERNAL_ERROR; received from peer
```

### 部分失效

有替代的检测，但仍保留失效检测的部分，未知是否完全失效

```
TikTok - 在 hk、jp 上测试时不时测不出，在 tw 上失效的概率更大，其他地区没有问题

BilibiliID - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

BilibiliTH - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

BilibiliVN - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

TVer - get platform-api.tver.jp failed with code: 400
```

### 无需支持

不需要再做支持

```
KBSAmerican - 不再支持本地区

Paravi - 已迁移并集成到 U-NEXT 中。由于整合，除了传统的Paravi作品外，现在还有电影、动漫、亚洲和外国戏剧等等可以无限观看。

HBOGOEurope - api.ugw.hbogo.eu 已经 host 为空了 查询不到内容

HBOSpain - api-discovery.hbo.eu 的 host 已经为空了

HBO_Nordic - 被合并了

HBO_Portugal - 被合并了
```

</details>

## 在Golang中使用

```
go get github.com/oneclickvirt/UnlockTests@latest
```

## Thanks

https://github.com/nkeonkeo/MediaUnlockTest

https://github.com/HsukqiLee/MediaUnlockTest

https://github.com/lmc999/RegionRestrictionCheck

https://github.com/betaxab/RegionRestrictionCheck/tree/refactor-1
