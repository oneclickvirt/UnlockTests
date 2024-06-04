# UnlockTests

[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Foneclickvirt%2FUnlockTests&count_bg=%2323E01C&title_bg=%23555555&icon=sonarcloud.svg&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com) [![Build and Release](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml/badge.svg)](https://github.com/oneclickvirt/UnlockTests/actions/workflows/main.yaml)

解锁测试模块 (Unlock Tests Module)

开发中，勿要使用

## 使用

更新时间：2024.06.04

安装

```shell
curl https://raw.githubusercontent.com/oneclickvirt/UnlockTests/main/ut_install.sh -sSf | sh
```

以后需要使用时使用

```
UT
```

唤起菜单

无环境依赖，理论上适配所有系统和主流架构，更多架构请查看 https://github.com/oneclickvirt/UnlockTests/releases/tag/output

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

MetaAI

KPLUS - ssoToken 已过期

TV360 - 登录认证已过期

Crackle - Platform Key is not specified

Salto - Get remote error: tls: unrecognized name

PCRJP - stream error: stream ID 1; INTERNAL_ERROR; received from peer

PrettyDerby - stream error: stream ID 1; INTERNAL_ERROR; received from peer

WorldFlipper - stream error: stream ID 1; INTERNAL_ERROR; received from peer

Catchplay - unauthorized 原 token 已过期

BahamutAnime - 存在 cloudflare 的质询防御，非5秒盾，无法突破，需要js动态加载
```

### 部分失效

有替代的检测，但仍保留失效检测的部分，未知是否完全失效

```
SonyLiv - 获取不到region

AISPlay - {"head": {"title": "Unidentified device", "itype": "item"}, "error": "Unidentified device", "st": {"code": 60400, "info": "Unidentified device"}, "op": {"cache_time": 600000, "timeout": 3000}}

Au7plus - 7plus-sevennetwork.akamaized.net 无论如何请求都失败

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
