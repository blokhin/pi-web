package ui

import (
	"net/http"
	"strings"
)

// X-Forwarded-Prefix support. A reverse proxy can mount pi-web under a URL
// prefix (nginx: `location /agent/pi/ { proxy_pass http://127.0.0.1:31415/; }`)
// which strips the prefix before the request reaches this server and declares
// it via the X-Forwarded-Prefix header. The server then re-attaches the prefix
// to every URL it hands to the browser (asset hrefs, redirects, the injected
// SPA base path) so the app works under the proxied path without a rebuild.
// An empty/absent/malformed header means "served at the origin root" — the
// pre-existing behavior, byte for byte.

// ForwardedPrefixHeader is the header reverse proxies use to declare the URL
// prefix they stripped from the request path.
const ForwardedPrefixHeader = "X-Forwarded-Prefix"

const (
	// maxPrefixLength bounds the header budget; real mounts are far shorter.
	maxPrefixLength = 256
	// prefixRejectChars never appear in a sane mount path. Any of them means a
	// typo or a header-injection attempt: reject the whole value, never
	// sanitise.
	prefixRejectChars = "\"'<> \n\r\t"
)

// NormalisePrefix validates and normalises a raw X-Forwarded-Prefix value.
// It returns the prefix in "/agent/pi" form (leading slash, no trailing
// slash) or "" when the value is empty or malformed. "//", ".." and the
// characters in prefixRejectChars are rejected so a hostile proxy cannot
// smuggle HTML or path traversal into served pages.
func NormalisePrefix(raw string) string {
	p := strings.TrimSpace(raw)
	if p == "" {
		return ""
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	p = strings.TrimRight(p, "/")
	if p == "" {
		return ""
	}
	if strings.Contains(p, "//") || strings.Contains(p, "..") || strings.ContainsAny(p, prefixRejectChars) {
		return ""
	}
	if len(p) > maxPrefixLength {
		return ""
	}
	return p
}

// PrefixFromRequest returns the normalised X-Forwarded-Prefix of a request,
// or "" when the request carries no (valid) prefix.
func PrefixFromRequest(r *http.Request) string {
	if r == nil {
		return ""
	}
	return NormalisePrefix(r.Header.Get(ForwardedPrefixHeader))
}

// PrefixPath joins a validated prefix onto a root-absolute path. An empty
// prefix, or a path that is not root-absolute, is returned unchanged — so
// no-prefix rendering is byte-identical to the pre-prefix code. Exported so
// other packages that build browser-facing URLs (auth redirects, cookies)
// share the one join rule.
func PrefixPath(prefix, path string) string {
	if prefix == "" || !strings.HasPrefix(path, "/") {
		return path
	}
	return prefix + path
}
