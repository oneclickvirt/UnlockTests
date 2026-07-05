package utils

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func TestReqHandlesNilHTTPClient(t *testing.T) {
	client := Req(nil)
	if client == nil {
		t.Fatalf("expected req client")
	}
}

func TestReqCopiesHTTPTransportSettings(t *testing.T) {
	proxyURL, err := url.Parse("http://127.0.0.1:8080")
	if err != nil {
		t.Fatalf("parse proxy: %v", err)
	}
	dialer := func(ctx context.Context, network, addr string) (net.Conn, error) {
		return nil, context.Canceled
	}
	source := &http.Client{
		Transport: &http.Transport{
			Proxy:                  http.ProxyURL(proxyURL),
			DialContext:            dialer,
			DisableCompression:     true,
			MaxIdleConns:           77,
			MaxIdleConnsPerHost:    7,
			MaxConnsPerHost:        9,
			IdleConnTimeout:        11 * time.Second,
			ResponseHeaderTimeout:  12 * time.Second,
			ExpectContinueTimeout:  13 * time.Second,
			MaxResponseHeaderBytes: 12345,
			WriteBufferSize:        4096,
			ReadBufferSize:         8192,
		},
	}
	client := Req(source)
	if client.Transport.DialContext == nil {
		t.Fatalf("expected DialContext to be copied")
	}
	if client.Transport.Proxy == nil {
		t.Fatalf("expected Proxy to be copied")
	}
	if client.Transport.DisableCompression != true ||
		client.Transport.MaxIdleConns != 77 ||
		client.Transport.MaxIdleConnsPerHost != 7 ||
		client.Transport.MaxConnsPerHost != 9 ||
		client.Transport.IdleConnTimeout != 11*time.Second ||
		client.Transport.ResponseHeaderTimeout != 12*time.Second ||
		client.Transport.ExpectContinueTimeout != 13*time.Second ||
		client.Transport.MaxResponseHeaderBytes != 12345 ||
		client.Transport.WriteBufferSize != 4096 ||
		client.Transport.ReadBufferSize != 8192 {
		t.Fatalf("transport settings were not copied: %#v", client.Transport)
	}
}

func TestEffectiveReqTimeoutUsesShorterHTTPClientTimeout(t *testing.T) {
	source := &http.Client{Timeout: time.Second}
	if got := effectiveReqTimeout(source, 14*time.Second); got != time.Second {
		t.Fatalf("expected shorter client timeout, got %s", got)
	}
	if got := effectiveReqTimeout(source, 500*time.Millisecond); got != 500*time.Millisecond {
		t.Fatalf("expected shorter fallback timeout, got %s", got)
	}
}
