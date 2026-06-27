package asia

import "testing"

func TestBilibiliAnimeNilClient(t *testing.T) {
	result := BilibiliAnime(nil)
	if result.Name != "Bilibili Anime" {
		t.Fatalf("expected Bilibili Anime name, got %q", result.Name)
	}
	if result.Status != "" {
		t.Fatalf("expected empty status for nil client, got %q", result.Status)
	}
}
