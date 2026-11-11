package ui

import (
	"strings"
	"testing"
)

func renderAuthPromptToString(t *testing.T, cookieTheme, prefix string) string {
	t.Helper()
	var b strings.Builder
	if err := RenderAuthPrompt(&b, cookieTheme, prefix); err != nil {
		t.Fatalf("RenderAuthPrompt: %v", err)
	}
	return b.String()
}

func TestAuthPromptRootKeepsRootAbsoluteRefs(t *testing.T) {
	html := renderAuthPromptToString(t, "", "")
	for _, want := range []string{`href="/icon.svg"`, `url("/icon.svg")`, `href="/custom-themes.css"`} {
		if !strings.Contains(html, want) {
			t.Fatalf("root auth prompt missing %q", want)
		}
	}
	if strings.Contains(html, "__PI_BASE_PATH__") {
		t.Fatalf("auth prompt must not declare a SPA base path")
	}
}

func TestAuthPromptHonorsPrefix(t *testing.T) {
	html := renderAuthPromptToString(t, "", "/agent/pi")
	for _, want := range []string{
		`href="/agent/pi/icon.svg"`,
		`url("/agent/pi/icon.svg")`,
		`href="/agent/pi/custom-themes.css"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("prefixed auth prompt missing %q\n%s", want, html)
		}
	}
	if strings.Contains(html, `href="/custom-themes.css"`) {
		t.Fatalf("prefixed auth prompt must not leak root-absolute stylesheet href\n%s", html)
	}
}
