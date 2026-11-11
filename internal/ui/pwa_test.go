package ui

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestRewriteManifestRootIsUntouched(t *testing.T) {
	if got := string(rewriteManifest("")); got != manifestJSON {
		t.Fatalf("empty prefix must serve the embedded manifest bytes untouched")
	}
}

func TestRewriteManifestHonorsPrefix(t *testing.T) {
	out := rewriteManifest("/agent/pi")
	var doc manifestDoc
	if err := json.Unmarshal(out, &doc); err != nil {
		t.Fatalf("rewritten manifest is not valid JSON: %v\n%s", err, out)
	}
	if doc.StartURL != "/agent/pi" {
		t.Fatalf("start_url = %q, want %q", doc.StartURL, "/agent/pi")
	}
	if doc.Scope != "/agent/pi/" {
		t.Fatalf("scope = %q, want %q", doc.Scope, "/agent/pi/")
	}
	if len(doc.Icons) == 0 {
		t.Fatalf("rewritten manifest lost its icons")
	}
	for _, icon := range doc.Icons {
		if !strings.HasPrefix(icon.Src, "/agent/pi/") {
			t.Fatalf("icon src = %q, want it under the prefix", icon.Src)
		}
	}
}

func TestRewriteServiceWorkerRootIsUntouched(t *testing.T) {
	if got := string(rewriteServiceWorker("")); got != swJS {
		t.Fatalf("empty prefix must serve the embedded service worker bytes untouched")
	}
}

func TestRewriteServiceWorkerHonorsPrefix(t *testing.T) {
	out := string(rewriteServiceWorker("/agent/pi"))
	for _, want := range []string{
		`icon: '/agent/pi/icon.svg'`,
		`badge: '/agent/pi/icon.svg'`,
		"`/agent/pi/session?id=${encodeURIComponent(sessionId)}`",
		`: '/agent/pi'`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("prefixed service worker missing %q\n%s", want, out)
		}
	}
	if strings.Contains(out, "'/icon.svg'") {
		t.Fatalf("prefixed service worker must not leak root-absolute icon\n%s", out)
	}
	if strings.Contains(out, ": '/'") {
		t.Fatalf("prefixed service worker must not leak root-absolute openWindow target\n%s", out)
	}
}
