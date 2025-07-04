package model

const UnlockTestsVersion = "v0.0.27"

var EnableLoger = false

type Result struct {
	Name       string
	Status     string
	Err        error
	Region     string
	Info       string
	UnlockType string
}

const (
	StatusUnexpected = "Unknown"
	StatusNetworkErr = "NetworkError"
	StatusErr        = "Error"
	StatusRestricted = "Restricted"
	StatusYes        = "Yes"
	StatusNo         = "No"
	StatusBanned     = "Banned"
	PrintHead        = "PrintHead"
	UA_Browser       = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36"
	UA_SecCHUA       = "\"Chromium\";v=\"124\", \"Google Chrome\";v=\"124\", \"Not-A.Brand\";v=\"99\""
	UA_Dalvik        = "Mozilla/5.0 (Linux; Android 10; Pixel 4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Mobile Safari/537.36"
	UA_Pjsekai       = "pjsekai/48 CFNetwork/1240.0.4 Darwin/20.6.0"
)

var PrivateIPv4Ranges = []string{
	"10.0.0.0/8",     // RFC 1918
	"172.16.0.0/12",  // RFC 1918
	"192.168.0.0/16", // RFC 1918
	"169.254.0.0/16", // 链路本地地址
	"127.0.0.0/8",    // 回环地址
	"0.0.0.0/8",      // 本网络
	"100.64.0.0/10",  // RFC 6598 (CGNAT)
	"224.0.0.0/4",    // 组播地址
}

var CommonPublicDNS = map[string]bool{
	"8.8.8.8":         true, // Google Public DNS
	"8.8.4.4":         true, // Google Public DNS
	"1.1.1.1":         true, // Cloudflare DNS
	"1.0.0.1":         true, // Cloudflare DNS
	"9.9.9.9":         true, // Quad9
	"149.112.112.112": true, // Quad9
	"208.67.222.222":  true, // OpenDNS
	"208.67.220.220":  true, // OpenDNS
	"64.6.64.6":       true, // Verisign DNS
	"64.6.65.6":       true, // Verisign DNS
	"94.140.14.14":    true, // AdGuard DNS
	"94.140.15.15":    true, // AdGuard DNS
	"185.228.168.9":   true, // CleanBrowsing
	"185.228.169.9":   true, // CleanBrowsing
	"76.76.19.19":     true, // Alternate DNS
	"76.223.122.150":  true, // Alternate DNS
	"77.88.8.8":       true, // Yandex.DNS
	"77.88.8.1":       true, // Yandex.DNS
	"156.154.70.1":    true, // Neustar DNS
	"156.154.71.1":    true, // Neustar DNS
	"8.26.56.26":      true, // Comodo Secure DNS
	"8.20.247.20":     true, // Comodo Secure DNS
	"84.200.69.80":    true, // DNS.WATCH
	"84.200.70.40":    true, // DNS.WATCH
}

var CdnPrefixes = []string{
	// Cloudflare
	"103.21.244.", "103.22.200.", "103.31.4.", "104.16.", "104.17.", "104.18.", "104.19.",
	"104.20.", "104.21.", "104.22.", "104.23.", "104.24.", "104.25.", "104.26.", "104.27.",
	"104.28.", "108.162.192.", "141.101.64.", "141.101.65.", "172.64.", "172.65.", "172.66.",
	"172.67.", "173.245.48.", "188.114.96.", "188.114.97.", "188.114.98.", "188.114.99.",
	"190.93.240.", "190.93.241.", "190.93.242.", "190.93.243.", "197.234.240.", "198.41.128.",

	// Akamai
	"23.32.", "23.33.", "23.34.", "23.35.", "23.36.", "23.37.", "23.38.", "23.39.",
	"23.40.", "23.41.", "23.42.", "23.43.", "23.44.", "23.45.", "23.46.", "23.47.",
	"23.48.", "23.49.", "23.50.", "23.51.", "23.52.", "23.53.", "23.54.", "23.55.",
	"23.56.", "23.57.", "23.58.", "23.59.", "23.60.", "23.61.", "23.62.", "23.63.",
	"23.64.", "23.65.", "23.66.", "23.67.", "104.64.", "104.65.", "104.66.", "104.67.",
	"104.68.", "104.69.", "104.70.", "104.71.", "104.72.", "104.73.", "104.74.", "104.75.",

	// Fastly
	"23.235.", "43.249.72.", "103.244.50.", "103.245.222.", "103.245.224.", "104.156.80.",
	"140.248.64.", "140.248.128.", "146.75.", "151.101.", "157.52.", "167.82.", "167.83.",
	"172.111.", "185.31.16.", "185.31.17.", "199.27.72.", "199.232.0.",

	// Amazon CloudFront
	"13.32.", "13.33.", "13.34.", "13.35.", "13.48.", "13.54.", "13.59.", "13.113.",
	"13.124.", "13.225.", "13.226.", "13.227.", "13.228.", "13.249.", "13.250.", "18.64.",
	"52.46.", "52.47.", "52.84.", "52.85.", "52.86.", "52.124.", "52.125.", "52.192.",
	"52.212.", "52.222.", "54.182.", "54.192.", "54.230.", "54.233.", "54.239.", "54.240.",
	"99.84.", "99.86.", "205.251.192.", "205.251.249.", "205.251.250.", "205.251.251.",
	"205.251.252.", "205.251.253.", "205.251.254.", "205.251.255.",

	// Microsoft Azure CDN
	"13.65.", "13.66.", "13.67.", "13.68.", "13.69.", "13.70.", "13.71.", "13.72.",
	"13.73.", "13.74.", "13.75.", "13.76.", "13.77.", "13.78.", "13.79.", "13.80.",
	"13.81.", "13.82.", "13.83.", "13.84.", "13.85.", "13.86.", "13.87.", "13.88.",
	"13.89.", "13.90.", "13.91.", "13.92.", "13.93.", "13.94.", "13.95.", "13.107.",

	// Google Cloud CDN
	"34.64.", "34.65.", "34.66.", "34.67.", "34.68.", "34.69.", "34.70.", "34.71.",
	"34.96.", "34.97.", "34.98.", "34.99.", "34.100.", "34.101.", "34.102.", "34.103.",
	"34.104.", "34.105.", "34.106.", "34.107.", "34.149.", "34.150.", "34.151.", "34.152.",
}

var StarPlusSupportCountry = []string{
	"br", "mx", "ar", "cl", "co", "pe", "uy", "ec", "pa", "cr", "py", "bo",
	"gt", "ni", "do", "sv", "hn", "ve",
}

var GptSupportCountry = []string{
	"al", "dz", "ad", "ao", "ag", "ar", "am", "au", "at", "az", "bs", "bd",
	"bb", "be", "bz", "bj", "bt", "ba", "bw", "br", "bg", "bf", "cv", "ca",
	"cl", "co", "km", "cr", "hr", "cy", "dk", "dj", "dm", "do", "ec", "sv",
	"ee", "fj", "fi", "fr", "ga", "gm", "ge", "de", "gh", "gr", "gd", "gt",
	"gn", "gw", "gy", "ht", "hn", "hu", "is", "in", "id", "iq", "ie", "il",
	"it", "jm", "jp", "jo", "kz", "ke", "ki", "kw", "kg", "lv", "lb", "ls",
	"lr", "li", "lt", "lu", "mg", "mw", "my", "mv", "ml", "mt", "mh", "mr",
	"mu", "mx", "mc", "mn", "me", "ma", "mz", "mm", "na", "nr", "np", "nl",
	"nz", "ni", "ne", "ng", "mk", "no", "om", "pk", "pw", "pa", "pg", "pe",
	"ph", "pl", "pt", "qa", "ro", "rw", "kn", "lc", "vc", "ws", "sm", "st",
	"sn", "rs", "sc", "sl", "sg", "sk", "si", "sb", "za", "es", "lk", "sr",
	"se", "ch", "th", "tg", "to", "tt", "tn", "tr", "tv", "ug", "ae", "us",
	"uy", "vu", "zm", "bo", "bn", "cg", "cz", "va", "fm", "md", "ps", "kr",
	"tw", "tz", "tl", "gb",
}

var ClaudeSupportCountry = []string{
	"al", "dz", "ad", "ao", "ag", "ar", "am", "au", "at", "az", "bs", "bh",
	"bd", "bb", "be", "bz", "bj", "bt", "bo", "ba", "bw", "br", "bn", "bg",
	"bf", "bi", "cv", "kh", "cm", "ca", "td", "cl", "co", "km", "cg", "cr",
	"ci", "hr", "cy", "cz", "dk", "dj", "dm", "do", "ec", "eg", "sv", "gq",
	"ee", "sz", "fj", "fi", "fr", "ga", "gm", "ge", "de", "gh", "gr", "gd",
	"gt", "gn", "gw", "gy", "ht", "hn", "hu", "is", "in", "id", "iq", "ie",
	"il", "it", "jm", "jp", "jo", "kz", "ke", "ki", "kw", "kg", "la", "lv",
	"lb", "ls", "lr", "li", "lt", "lu", "mg", "mw", "my", "mv", "mt", "mh",
	"mr", "mu", "mx", "fm", "md", "mc", "mn", "me", "ma", "mz", "na", "nr",
	"np", "nl", "nz", "ne", "ng", "mk", "no", "om", "pk", "pw", "ps", "pa",
	"pg", "py", "pe", "ph", "pl", "pt", "qa", "ro", "rw", "kn", "lc", "vc",
	"ws", "sm", "st", "sa", "sn", "rs", "sc", "sl", "sg", "sk", "si", "sb",
	"za", "kr", "es", "lk", "sr", "se", "ch", "tw", "tj", "tz", "th", "tl",
	"tg", "to", "tt", "tn", "tr", "tm", "tv", "ug", "ua", "ae", "gb", "us",
	"uy", "uz", "vu", "va", "vn", "zm", "zw",
}

var AppleSupportCountry = []string{
	"alb", "dza", "afg", "arg", "are", "abw", "omn", "aze", "egy", "eth", "irl", "est", "and", "ago", "aia", "atg",
	"aut", "ala", "aus", "mac", "brb", "png", "bhs", "pak", "pry", "pse", "bhr", "pan", "bra", "blr", "bmu", "bgr",
	"mnp", "mkd", "ben", "bel", "isl", "bol", "pri", "pol", "bih", "bwa", "blz", "btn", "bfa", "bdi", "bvt", "iot",
	"gnq", "dnk", "deu", "tls", "tgo", "dom", "dma", "rus", "ecu", "eri", "fra", "fro", "pyf", "guf", "atf", "phl",
	"fin", "cpv", "flk", "gmb", "cog", "cod", "col", "cri", "grd", "grl", "geo", "ggy", "glp", "gum", "guy", "kaz",
	"hti", "kor", "nld", "bes", "sxm", "hmd", "mne", "hnd", "kir", "dji", "kgz", "gin", "gnb", "can", "gha", "gab",
	"khm", "cze", "zwe", "cmr", "qat", "cym", "cck", "com", "xks", "civ", "kwt", "hrv", "ken", "cok", "cuw", "lva",
	"lso", "lao", "lbn", "lbr", "lby", "ltu", "lie", "reu", "lux", "rwa", "rou", "mdg", "mlt", "mdv", "mwi", "mys",
	"mli", "mhl", "mtq", "myt", "imn", "mus", "mrt", "usa", "umi", "asm", "vir", "mng", "msr", "bgd", "per", "fsm",
	"mmr", "mda", "mar", "mco", "moz", "mex", "nam", "zaf", "ata", "sgs", "ssd", "npl", "nic", "ner", "nga", "niu",
	"nor", "nfk", "plw", "pcn", "prt", "jpn", "swe", "che", "slv", "wsm", "srb", "sle", "sen", "cyp", "syc", "sau",
	"blm", "cxr", "stp", "shn", "kna", "lca", "maf", "smr", "spm", "vct", "lka", "svk", "svn", "sjm", "swz", "sdn",
	"sur", "som", "slb", "tjk", "twn", "tha", "tza", "ton", "tca", "tto", "tun", "tuv", "tur", "tkm", "tkl", "wlf",
	"vut", "gtm", "ven", "brn", "uga", "ukr", "ury", "uzb", "esp", "esh", "grc", "hkg", "sgp", "ncl", "nzl", "hun",
	"jam", "arm", "yem", "irq", "isr", "ita", "ind", "idn", "gbr", "vgb", "jor", "vnm", "zmb", "jey", "tcd", "gib",
	"chl", "caf", "chn", "nru", "vat", "fji",
}

var DiscoveryPlusSupportCountry = []string{
	"at", "br", "ca", "dk", "fi", "de", "in", "ie", "it", "nl", "no", "es",
	"se", "gb", "us"}

var SonyLivSupportCountry = []string{
	"AE", "AF", "AT", "AU", "BD", "BE", "BH", "BT", "CA", "CH", "CN", "DE", "DK", "ES", "FI",
	"FR", "GB", "GR", "HK", "ID", "IE", "IN", "IT", "KW", "LK", "MO", "MV", "MY", "NL", "NO",
	"NP", "NZ", "OM", "PH", "PK", "PL", "PT", "QA", "SA", "SE", "SG", "TH", "TW", "US",
}

var NLZIETSupportCountry = []string{
	"be", "bg", "cz", "dk", "de", "ee", "ie", "el", "es", "fr",
	"hr", "it", "cy", "lv", "lt", "lu", "hu", "mt", "nl", "at",
	"pl", "pt", "ro", "si", "sk", "fi", "se",
}

var GeminiSupportCountry = []string{
	"ax", "al", "dz", "as", "ad", "ao", "ai", "aq", "ag", "ar",
	"am", "aw", "au", "at", "az", "bh", "bd", "bb", "be", "bz",
	"bj", "bm", "bt", "bo", "ba", "bw", "br", "io", "vg", "bn",
	"bg", "bf", "bi", "cv", "kh", "cm", "ca", "bq", "ky", "cf",
	"td", "cl", "cx", "cc", "co", "km", "ck", "cr", "ci", "hr",
	"cw", "cz", "cd", "dk", "dj", "dm", "do", "ec", "eg", "sv",
	"gq", "er", "ee", "sz", "et", "fk", "fo", "fj", "fi", "fr",
	"gf", "pf", "tf", "ga", "ge", "de", "gh", "gi", "gr", "gl",
	"gd", "gp", "gu", "gt", "gg", "gn", "gw", "gy", "ht", "hm",
	"hn", "hu", "is", "in", "id", "iq", "ie", "im", "il", "it",
	"jm", "jp", "je", "jo", "kz", "ke", "ki", "xk", "kw", "kg",
	"la", "lv", "lb", "ls", "lr", "ly", "li", "lt", "lu", "mg",
	"mw", "my", "mv", "ml", "mt", "mh", "mq", "mr", "mu", "yt",
	"mx", "fm", "md", "mc", "mn", "me", "ms", "ma", "mz", "mm",
	"na", "nr", "np", "nl", "nc", "nz", "ni", "ne", "ng", "nu",
	"nf", "mk", "mp", "no", "om", "pk", "pw", "ps", "pa", "pg",
	"py", "pe", "ph", "pn", "pl", "pt", "pr", "qa", "cy", "cg",
	"re", "ro", "rw", "bl", "sh", "kn", "lc", "mf", "pm", "vc",
	"ws", "sm", "st", "sa", "sn", "rs", "sc", "sl", "sg", "sx",
	"sk", "si", "sb", "so", "za", "gs", "kr", "ss", "es", "lk",
	"sd", "sr", "sj", "se", "ch", "tw", "tj", "tz", "th", "bs",
	"gm", "tl", "tg", "tk", "to", "tt", "tn", "tr", "tm", "tc",
	"tv", "vi", "ug", "ua", "ae", "gb", "us", "um", "uy", "uz",
	"vu", "va", "ve", "vn", "wf", "eh", "ye", "zm", "zw",
}