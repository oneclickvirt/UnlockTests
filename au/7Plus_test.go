package au

import (
	"testing"

	"github.com/oneclickvirt/UnlockTests/model"
)

func TestEvaluateSevenPlusMarket(t *testing.T) {
	yes := evaluateSevenPlusMarket("7plus", sevenPlusMarketResponse{
		ID:        4,
		PlaceName: "Sydney South",
	}, "Native")
	if yes.Status != model.StatusYes || yes.Region != "au" || yes.Info != "Sydney South" || yes.UnlockType != "Native" {
		t.Fatalf("expected AU market to be unlocked, got %#v", yes)
	}

	no := evaluateSevenPlusMarket("7plus", sevenPlusMarketResponse{
		ID:        1,
		PlaceName: "Outside Australia",
	}, "")
	if no.Status != model.StatusNo || no.Info != "Outside Australia" {
		t.Fatalf("expected non-AU market to be locked, got %#v", no)
	}
}
