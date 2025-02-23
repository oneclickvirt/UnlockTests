package utils

import "strings"

func TwoToThreeCode(code string) string {
	countryCodes := map[string]string{
		"ad": "and", "ae": "are", "af": "afg", "ag": "atg", "ai": "aia", "al": "alb", "am": "arm", "ao": "ago", "aq": "ata", "ar": "arg",
		"as": "asm", "at": "aut", "au": "aus", "aw": "abw", "ax": "ala", "az": "aze", "ba": "bih", "bb": "brb", "bd": "bgd", "be": "bel",
		"bf": "bfa", "bg": "bgr", "bh": "bhr", "bi": "bdi", "bj": "ben", "bl": "blm", "bm": "bmu", "bn": "brn", "bo": "bol", "bq": "bes",
		"br": "bra", "bs": "bhs", "bt": "btn", "bv": "bvt", "bw": "bwa", "by": "blr", "bz": "blz", "ca": "can", "cc": "cck", "cd": "cod",
		"cf": "caf", "cg": "cog", "ch": "che", "ci": "civ", "ck": "cok", "cl": "chl", "cm": "cmr", "cn": "chn", "co": "col", "cr": "cri",
		"cu": "cub", "cv": "cpv", "cw": "cuw", "cx": "cxr", "cy": "cyp", "cz": "cze", "de": "deu", "dj": "dji", "dk": "dnk", "dm": "dma",
		"do": "dom", "dz": "dza", "ec": "ecu", "ee": "est", "eg": "egy", "eh": "esh", "er": "eri", "es": "esp", "et": "eth", "fi": "fin",
		"fj": "fji", "fk": "flk", "fm": "fsm", "fo": "fro", "fr": "fra", "ga": "gab", "gb": "gbr", "gd": "grd", "ge": "geo", "gf": "guf",
		"gg": "ggy", "gh": "gha", "gi": "gib", "gl": "grl", "gm": "gmb", "gn": "gin", "gp": "glp", "gq": "gnq", "gr": "grc", "gs": "sgs",
		"gt": "gtm", "gu": "gum", "gw": "gnb", "gy": "guy", "hk": "hkg", "hm": "hmd", "hn": "hnd", "hr": "hrv", "ht": "hti", "hu": "hun",
		"id": "idn", "ie": "irl", "il": "isr", "im": "imn", "in": "ind", "io": "iot", "iq": "irq", "ir": "irn", "is": "isl", "it": "ita",
		"je": "jey", "jm": "jam", "jo": "jor", "jp": "jpn", "ke": "ken", "kg": "kgz", "kh": "khm", "ki": "kir", "km": "com", "kn": "kna",
		"kp": "prk", "kr": "kor", "kw": "kwt", "ky": "cym", "kz": "kaz", "la": "lao", "lb": "lbn", "lc": "lca", "li": "lie", "lk": "lka",
		"lr": "lbr", "ls": "lso", "lt": "ltu", "lu": "lux", "lv": "lva", "ly": "lby", "ma": "mar", "mc": "mco", "md": "mda", "me": "mne",
		"mf": "maf", "mg": "mdg", "mh": "mhl", "mk": "mkd", "ml": "mli", "mm": "mmr", "mn": "mng", "mo": "mac", "mp": "mnp", "mq": "mtq",
		"mr": "mrt", "ms": "msr", "mt": "mlt", "mu": "mus", "mv": "mdv", "mw": "mwi", "mx": "mex", "my": "mys", "mz": "moz", "na": "nam",
		"nc": "ncl", "ne": "ner", "nf": "nfk", "ng": "nga", "ni": "nic", "nl": "nld", "no": "nor", "np": "npl", "nr": "nru", "nu": "niu",
		"nz": "nzl", "om": "omn", "pa": "pan", "pe": "per", "pf": "pyf", "pg": "png", "ph": "phl", "pk": "pak", "pl": "pol", "pm": "spm",
		"pn": "pcn", "pr": "pri", "ps": "pse", "pt": "prt", "pw": "plw", "py": "pry", "qa": "qat", "re": "reu", "ro": "rou", "rs": "srb",
		"ru": "rus", "rw": "rwa", "sa": "sau", "sb": "slb", "sc": "syc", "sd": "sdn", "se": "swe", "sg": "sgp", "sh": "shn", "si": "svn",
		"sj": "sjm", "sk": "svk", "sl": "sle", "sm": "smr", "sn": "sen", "so": "som", "sr": "sur", "ss": "ssd", "st": "stp", "sv": "slv",
		"sx": "sxm", "sy": "syr", "sz": "swz", "tc": "tca", "td": "tcd", "tf": "atf", "tg": "tgo", "th": "tha", "tj": "tjk", "tk": "tkl",
		"tl": "tls", "tm": "tkm", "tn": "tun", "to": "ton", "tr": "tur", "tt": "tto", "tv": "tuv", "tw": "twn", "tz": "tza", "ua": "ukr",
		"ug": "uga", "um": "umi", "us": "usa", "uy": "ury", "uz": "uzb", "va": "vat", "vc": "vct", "ve": "ven", "vg": "vgb", "vi": "vir",
		"vn": "vnm", "vu": "vut", "wf": "wlf", "ws": "wsm", "ye": "yem", "yt": "myt", "za": "zaf", "zm": "zmb", "zw": "zwe",
	}
	return countryCodes[strings.ToLower(code)]
}

func ThreeToTwoCode(code string) string {
	countryCodes := map[string]string{
		"and": "ad", "are": "ae", "afg": "af", "atg": "ag", "aia": "ai", "alb": "al", "arm": "am", "ago": "ao", "ata": "aq", "arg": "ar",
		"asm": "as", "aut": "at", "aus": "au", "abw": "aw", "ala": "ax", "aze": "az", "bih": "ba", "brb": "bb", "bgd": "bd", "bel": "be",
		"bfa": "bf", "bgr": "bg", "bhr": "bh", "bdi": "bi", "ben": "bj", "blm": "bl", "bmu": "bm", "brn": "bn", "bol": "bo", "bes": "bq",
		"bra": "br", "bhs": "bs", "btn": "bt", "bvt": "bv", "bwa": "bw", "blr": "by", "blz": "bz", "can": "ca", "cck": "cc", "cod": "cd",
		"caf": "cf", "cog": "cg", "che": "ch", "civ": "ci", "cok": "ck", "chl": "cl", "cmr": "cm", "chn": "cn", "col": "co", "cri": "cr",
		"cub": "cu", "cpv": "cv", "cuw": "cw", "cxr": "cx", "cyp": "cy", "cze": "cz", "deu": "de", "dji": "dj", "dnk": "dk", "dma": "dm",
		"dom": "do", "dza": "dz", "ecu": "ec", "est": "ee", "egy": "eg", "esh": "eh", "eri": "er", "esp": "es", "eth": "et", "fin": "fi",
		"fji": "fj", "flk": "fk", "fsm": "fm", "fro": "fo", "fra": "fr", "gab": "ga", "gbr": "gb", "grd": "gd", "geo": "ge", "guf": "gf",
		"ggy": "gg", "gha": "gh", "gib": "gi", "grl": "gl", "gmb": "gm", "gin": "gn", "glp": "gp", "gnq": "gq", "grc": "gr", "sgs": "gs",
		"gtm": "gt", "gum": "gu", "gnb": "gw", "guy": "gy", "hkg": "hk", "hmd": "hm", "hnd": "hn", "hrv": "hr", "hti": "ht", "hun": "hu",
		"idn": "id", "irl": "ie", "isr": "il", "imn": "im", "ind": "in", "iot": "io", "irq": "iq", "irn": "ir", "isl": "is", "ita": "it",
		"jey": "je", "jam": "jm", "jor": "jo", "jpn": "jp", "ken": "ke", "kgz": "kg", "khm": "kh", "kir": "ki", "com": "km", "kna": "kn",
		"prk": "kp", "kor": "kr", "kwt": "kw", "cym": "ky", "kaz": "kz", "lao": "la", "lbn": "lb", "lca": "lc", "lie": "li", "lka": "lk",
		"lbr": "lr", "lso": "ls", "ltu": "lt", "lux": "lu", "lva": "lv", "lby": "ly", "mar": "ma", "mco": "mc", "mda": "md", "mne": "me",
		"maf": "mf", "mdg": "mg", "mhl": "mh", "mkd": "mk", "mli": "ml", "mmr": "mm", "mng": "mn", "mac": "mo", "mnp": "mp", "mtq": "mq",
		"mrt": "mr", "msr": "ms", "mlt": "mt", "mus": "mu", "mdv": "mv", "mwi": "mw", "mex": "mx", "mys": "my", "moz": "mz", "nam": "na",
		"ncl": "nc", "ner": "ne", "nfk": "nf", "nga": "ng", "nic": "ni", "nld": "nl", "nor": "no", "npl": "np", "nru": "nr", "niu": "nu",
		"nzl": "nz", "omn": "om", "pan": "pa", "per": "pe", "pyf": "pf", "png": "pg", "phl": "ph", "pak": "pk", "pol": "pl", "spm": "pm",
		"pcn": "pn", "pri": "pr", "pse": "ps", "prt": "pt", "plw": "pw", "pry": "py", "qat": "qa", "reu": "re", "rou": "ro", "srb": "rs",
		"rus": "ru", "rwa": "rw", "sau": "sa", "slb": "sb", "syc": "sc", "sdn": "sd", "swe": "se", "sgp": "sg", "shn": "sh", "svn": "si",
		"sjm": "sj", "svk": "sk", "sle": "sl", "smr": "sm", "sen": "sn", "som": "so", "sur": "sr", "ssd": "ss", "stp": "st", "slv": "sv",
		"sxm": "sx", "syr": "sy", "swz": "sz", "tca": "tc", "tcd": "td", "atf": "tf", "tgo": "tg", "tha": "th", "tjk": "tj", "tkl": "tk",
		"tls": "tl", "tkm": "tm", "tun": "tn", "ton": "to", "tur": "tr", "tto": "tt", "tuv": "tv", "twn": "tw", "tza": "tz", "ukr": "ua",
		"uga": "ug", "umi": "um", "usa": "us", "ury": "uy", "uzb": "uz", "vat": "va", "vct": "vc", "ven": "ve", "vgb": "vg", "vir": "vi",
		"vnm": "vn", "vut": "vu", "wlf": "wf", "wsm": "ws", "yem": "ye", "myt": "yt", "zaf": "za", "zmb": "zm", "zwe": "zw",
	}
	return countryCodes[strings.ToLower(code)]
}
