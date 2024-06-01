# UnlockTests

[![Hits](https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Foneclickvirt%2FUnlockTests&count_bg=%2323E01C&title_bg=%23555555&icon=sonarcloud.svg&icon_color=%23E7E7E7&title=hits&edge_flat=false)](https://hits.seeyoufarm.com)

解锁测试模块 (Unlock Tests Module)

开发中，勿要使用

## TODO

<details>

### 同态检测

可能需要拆分检测

```
TLCGO 和 NBCTV 
```

### 无效检测

需要重新构建检测逻辑

```
HBO_Nordic

HBO_Portugal

ElevenSportsTW

MathsSpot

MegogoTV

KPLUS - ssoToken 已过期

TV360 - 登录认证已过期

Crackle - Platform Key is not specified

Eurosport - Tokem 已过期 且 api 官网已升级至于 v3

HBOGOEurope - api.ugw.hbogo.eu 已经 host 为空了 查询不到内容

HBOSpain - api-discovery.hbo.eu 的 host 已经为空了

Salto - Get remote error: tls: unrecognized name

PCRJP - stream error: stream ID 1; INTERNAL_ERROR; received from peer

PrettyDerby - stream error: stream ID 1; INTERNAL_ERROR; received from peer

WorldFlipper - stream error: stream ID 1; INTERNAL_ERROR; received from peer
```

### 地区失效

不需要再做支持

```
KBSAmerican - 不再支持本地区

Paravi - 已迁移并集成到 U-NEXT 中。由于整合，除了传统的Paravi作品外，现在还有电影、动漫、亚洲和外国戏剧等等可以无限观看。
```

### 部分失效

有替代的检测，但仍保留失效检测的部分，未知是否完全失效

```
Au7plus - 7plus-sevennetwork.akamaized.net 无论如何请求都失败

BilibiliID - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

BilibiliTH - 对应URL请求无论如何都返回为空 {"code":10004001,"message":"10004001","ttl":1,"data":null}

TVer - get platform-api.tver.jp failed with code: 400
```

</details>
