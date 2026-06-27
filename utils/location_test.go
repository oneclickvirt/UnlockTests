package utils

import "testing"

func TestCountryCodeToAlpha2(t *testing.T) {
	tests := map[string]string{
		"+852": "HK",
		"853":  "MO",
		"886":  "TW",
		"66":   "TH",
		"62":   "ID",
		"84":   "VN",
		"60":   "MY",
		"65":   "SG",
		"63":   "PH",
		"673":  "BN",
		"855":  "KH",
		"856":  "LA",
		"95":   "MM",
		"670":  "TL",
		"GBR":  "GB",
		"gb":   "GB",
		"999":  "",
	}
	for code, want := range tests {
		if got := CountryCodeToAlpha2(code); got != want {
			t.Fatalf("CountryCodeToAlpha2(%q) = %q, want %q", code, got, want)
		}
	}
}
