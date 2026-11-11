package ui

import (
	"strings"
	"testing"
)

func TestAppShellPreservesPWAContract(t *testing.T) {
	old := appScriptPath
	appScriptPath = "/static/assets/app-test.js"
	defer func() { appScriptPath = old }()

	var b strings.Builder
	if err := RenderAppShell(&b, "", ""); err != nil {
		t.Fatalf("RenderAppShell: %v", err)
	}
	html := b.String()
	for _, want := range []string{
		`<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, interactive-widget=resizes-content">`,
		`<link rel="icon" type="image/svg+xml" href="/icon.svg">`,
		`<link rel="apple-touch-icon" href="/icon.svg">`,
		`<link rel="manifest" href="/manifest.webmanifest">`,
		`<meta name="theme-color" content="#0e0e13">`,
		`<meta name="mobile-web-app-capable" content="yes">`,
		`<meta name="apple-mobile-web-app-status-bar-style" content="black-translucent">`,
		`<meta name="apple-mobile-web-app-title" content="Pi Sessions">`,
		`<meta name="pi-web-theme"`,
		`navigator.windowControlsOverlay`,
		`<link rel="stylesheet" href="/custom-themes.css">`,
		`<style id="pi-web-fonts">`,
		`<div id="spa-root"></div>`,
		`<script type="module" src="/static/assets/app-test.js"></script>`,
		`navigator.serviceWorker.register('/sw.js',{scope:'/'})`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("app shell missing %q\n%s", want, html)
		}
	}
}

func TestAppShellWithoutPrefixIsByteIdenticalToRoot(t *testing.T) {
	old := appScriptPath
	appScriptPath = "/static/assets/app-test.js"
	defer func() { appScriptPath = old }()

	var b strings.Builder
	if err := RenderAppShell(&b, "", ""); err != nil {
		t.Fatalf("RenderAppShell: %v", err)
	}
	html := b.String()
	if strings.Contains(html, "__PI_BASE_PATH__") {
		t.Fatalf("no-prefix shell must not declare a base path\n%s", html)
	}
	if strings.Contains(html, "agent/pi") {
		t.Fatalf("no-prefix shell must not mention any prefix\n%s", html)
	}
}

func TestAppShellHonorsForwardedPrefix(t *testing.T) {
	old := appScriptPath
	appScriptPath = "/static/assets/app-test.js"
	defer func() { appScriptPath = old }()

	var b strings.Builder
	if err := RenderAppShell(&b, "", "/agent/pi"); err != nil {
		t.Fatalf("RenderAppShell: %v", err)
	}
	html := b.String()
	for _, want := range []string{
		`window.__PI_BASE_PATH__="/agent/pi";`,
		`<link rel="icon" type="image/svg+xml" href="/agent/pi/icon.svg">`,
		`<link rel="apple-touch-icon" href="/agent/pi/icon.svg">`,
		`<link rel="manifest" href="/agent/pi/manifest.webmanifest">`,
		`<link rel="stylesheet" href="/agent/pi/custom-themes.css">`,
		`<script type="module" src="/agent/pi/static/assets/app-test.js"></script>`,
		`navigator.serviceWorker.register('/agent/pi/sw.js',{scope:'/agent/pi/'})`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("prefixed app shell missing %q\n%s", want, html)
		}
	}
	if strings.Contains(html, `href="/icon.svg"`) {
		t.Fatalf("prefixed shell must not leak root-absolute icon href\n%s", html)
	}
}
