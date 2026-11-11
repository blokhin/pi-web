package ui

import (
	"bytes"
	_ "embed"
	"html/template"
	"io"
	"strings"
)

//go:embed embedded/auth.html
var authTmplStr string

var authTmpl = template.Must(template.New("auth").Parse(authTmplStr))

// RenderAuthPrompt writes the standalone token-gate page. It reuses the shared
// theme stylesheet and boot script so the login screen honors every theme
// instead of a hardcoded dark/light copy that drifts from theme.css. The
// runtime "custom" theme resolves too: auth.html pulls /custom-themes.css,
// which is served publicly for this pre-auth case.
//
// cookieTheme is the request's pi-web-theme cookie (the live theme toggle's
// carrier); when empty it falls back to the server-persisted theme so the login
// screen matches the user's choice even before any client storage exists.
//
// prefix is the request's validated X-Forwarded-Prefix ("" = origin root).
// auth.html is a static template whose two root-absolute references (icon,
// custom-themes.css) are re-pointed at the prefixed mount after rendering; an
// empty prefix renders the template untouched.
func RenderAuthPrompt(w io.Writer, cookieTheme, prefix string) error {
	theme := cookieTheme
	if theme == "" {
		theme = themeProvider()
	}
	data := struct {
		ServerTheme string
		ThemeCss    template.HTML
		ThemeBoot   template.HTML
	}{
		ServerTheme: theme,
		ThemeCss:    template.HTML("<style>\n" + liveThemeCss + "\n</style>"),
		ThemeBoot:  liveThemeBootScript(),
	}
	var buf bytes.Buffer
	if err := authTmpl.Execute(&buf, data); err != nil {
		return err
	}
	html := buf.String()
	if prefix != "" {
		html = strings.ReplaceAll(html, `"/icon.svg"`, `"`+prefix+`/icon.svg"`)
		html = strings.ReplaceAll(html, `"/custom-themes.css"`, `"`+prefix+`/custom-themes.css"`)
	}
	_, err := w.Write([]byte(html))
	return err
}
