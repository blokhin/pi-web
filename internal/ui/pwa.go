package ui

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

//go:embed embedded/assets/manifest.webmanifest
var manifestJSON string

//go:embed embedded/assets/sw.js
var swJS string

//go:embed embedded/assets/icon.svg
var iconSVG string

//go:embed embedded/assets/icon-maskable.svg
var iconMaskableSVG string

//go:embed embedded/assets/pi-logo.svg
var piLogoSVG string

//go:embed embedded/assets/cat.mp3
var CatMP3 []byte

//go:embed embedded/assets/done.mp3
var DoneMP3 []byte

//go:embed embedded/assets/cat.webm
var catWebm []byte

// indexCSS and settingsCSS are inlined into the SPA shell by
// appStylesheets() (spa_page.go); they are not served as standalone routes.
//go:embed embedded/styles/index.css
var indexCSS string

//go:embed embedded/styles/settings.css
var settingsCSS string

//go:embed embedded/styles/schedules.css
var schedulesCSS string

// manifestDoc mirrors the fields of embedded/assets/manifest.webmanifest that
// contain origin-rooted paths; under an X-Forwarded-Prefix they must be
// re-prefixed so the PWA installs at the proxied path.
type manifestDoc struct {
	StartURL string `json:"start_url"`
	Scope    string `json:"scope"`
	Icons    []struct {
		Src string `json:"src"`
	} `json:"icons"`
}

// rewriteManifest re-points the manifest's start_url, scope and icon paths at
// the prefixed mount. With an empty prefix the embedded bytes are served
// untouched.
func rewriteManifest(prefix string) []byte {
	if prefix == "" {
		return []byte(manifestJSON)
	}
	var doc manifestDoc
	if err := json.Unmarshal([]byte(manifestJSON), &doc); err != nil {
		return []byte(manifestJSON)
	}
	if doc.StartURL != "" {
		doc.StartURL = PrefixPath(prefix, doc.StartURL)
	}
	if doc.Scope != "" {
		doc.Scope = PrefixPath(prefix, doc.Scope)
		if !strings.HasSuffix(doc.Scope, "/") {
			doc.Scope += "/"
		}
	}
	for i := range doc.Icons {
		doc.Icons[i].Src = PrefixPath(prefix, doc.Icons[i].Src)
	}
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return []byte(manifestJSON)
	}
	return out
}

// rewriteServiceWorker re-points the service worker's hardcoded origin-rooted
// URLs (notification icon, notificationclick target) at the prefixed mount.
// The replacements target the exact literals in embedded/assets/sw.js; with an
// empty prefix the embedded bytes are served untouched.
func rewriteServiceWorker(prefix string) []byte {
	if prefix == "" {
		return []byte(swJS)
	}
	sw := strings.ReplaceAll(swJS, "'/icon.svg'", "'"+prefix+"/icon.svg'")
	sw = strings.ReplaceAll(sw, "`/session?id=${encodeURIComponent(sessionId)}`", "`"+prefix+"/session?id=${encodeURIComponent(sessionId)}`")
	sw = strings.ReplaceAll(sw, ": '/'", ": '"+prefix+"'")
	return []byte(sw)
}

// registerPWAHandlers serves the manifest, service worker, and icons.
// Routes are registered without auth: a manifest/icon leaks nothing
// sensitive, and the service worker must be reachable for installability
// even before the user authenticates.
func RegisterPWAHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/manifest.webmanifest", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/manifest+json")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write(rewriteManifest(PrefixFromRequest(r)))
	})
	mux.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		// Allow the SW to control the whole origin.
		w.Header().Set("Service-Worker-Allowed", "/")
		w.Header().Set("Cache-Control", "no-cache")
		_, _ = w.Write(rewriteServiceWorker(PrefixFromRequest(r)))
	})
	mux.HandleFunc("/icon.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write([]byte(iconSVG))
	})
	mux.HandleFunc("/icon-maskable.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write([]byte(iconMaskableSVG))
	})
	mux.HandleFunc("/pi-logo.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		_, _ = w.Write([]byte(piLogoSVG))
	})
	mux.HandleFunc("/cat.webm", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "video/webm")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeContent(w, r, "cat.webm", time.Time{}, bytes.NewReader(catWebm))
	})
}
