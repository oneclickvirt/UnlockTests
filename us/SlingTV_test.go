package us

import "testing"

func TestParseSlingGeoResponse(t *testing.T) {
	region, blocked := parseSlingGeoResponse([]byte(`{"country_code":"US","ip_restricted":false}`))
	if region != "us" || blocked {
		t.Fatalf("expected us/unblocked, got region=%q blocked=%v", region, blocked)
	}

	region, blocked = parseSlingGeoResponse([]byte(`{"countryCode":"CA","blocked":true}`))
	if region != "ca" || !blocked {
		t.Fatalf("expected ca/blocked, got region=%q blocked=%v", region, blocked)
	}

	region, blocked = parseSlingGeoResponse([]byte(`{"country":"USA","ip_restricted":false}`))
	if region != "us" || blocked {
		t.Fatalf("expected usa to normalize to us/unblocked, got region=%q blocked=%v", region, blocked)
	}
}
