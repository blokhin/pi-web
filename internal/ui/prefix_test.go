package ui

import (
	"net/http"
	"strings"
	"testing"
)

func TestNormalisePrefix(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{"empty is root", "", ""},
		{"whitespace only is root", "   ", ""},
		{"plain prefix", "/agent/pi", "/agent/pi"},
		{"adds missing leading slash", "agent/pi", "/agent/pi"},
		{"strips trailing slash", "/agent/pi/", "/agent/pi"},
		{"strips many trailing slashes", "/agent/pi///", "/agent/pi"},
		{"root slash only is root", "/", ""},
		{"trims surrounding whitespace", "  /agent/pi  ", "/agent/pi"},
		{"nested segments keep inner slashes", "/a/b/c", "/a/b/c"},
		{"rejects double slash", "/agent//pi", ""},
		{"rejects traversal", "/agent/../pi", ""},
		{"rejects leading traversal", "..", ""},
		{"rejects bare traversal segment", "/..", ""},
		{"rejects quote", `/agent/"pi`, ""},
		{"rejects single quote", "/agent/'pi'", ""},
		{"rejects angle bracket", "/agent/<script>", ""},
		{"rejects space", "/agent pi", ""},
		{"rejects newline", "/agent\n/pi", ""},
		{"rejects tab", "/agent\tpi", ""},
		{"rejects overlong value", "/" + strings.Repeat("a", 300), ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalisePrefix(tt.raw); got != tt.want {
				t.Fatalf("NormalisePrefix(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestPrefixFromRequest(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{"no header", nil, ""},
		{"simple header", map[string]string{"X-Forwarded-Prefix": "/agent/pi"}, "/agent/pi"},
		{"malformed header falls back to root", map[string]string{"X-Forwarded-Prefix": "/a/../b"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			for k, v := range tt.headers {
				r.Header.Set(k, v)
			}
			if got := PrefixFromRequest(r); got != tt.want {
				t.Fatalf("PrefixFromRequest() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPrefixPath(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
		path   string
		want   string
	}{
		{"no prefix keeps path", "", "/icon.svg", "/icon.svg"},
		{"prefix joins root-absolute path", "/agent/pi", "/icon.svg", "/agent/pi/icon.svg"},
		{"prefix keeps query on path", "/agent/pi", "/session?id=x", "/agent/pi/session?id=x"},
		{"relative path untouched", "/agent/pi", "icon.svg", "icon.svg"},
		{"absolute URL untouched", "/agent/pi", "https://example.com/x", "https://example.com/x"},
		{"root untouched without prefix", "", "/", "/"},
		{"root gains prefix", "/agent/pi", "/", "/agent/pi/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrefixPath(tt.prefix, tt.path); got != tt.want {
				t.Fatalf("PrefixPath(%q, %q) = %q, want %q", tt.prefix, tt.path, got, tt.want)
			}
		})
	}
}
