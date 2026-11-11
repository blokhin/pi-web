package ui

import (
	"bytes"
	_ "embed"
	"html/template"
	"io"
)

//go:embed embedded/app.html
var appTmplStr string

var appTmpl = template.Must(template.New("app").Parse(appTmplStr))

func appStylesheets() template.HTML {
	return template.HTML("<style>\n" + liveThemeCss + "\n" + indexCSS + "\n" + settingsCSS + "\n" + schedulesCSS + "\n" + liveSessionCss + "\n" + liveMenuCss + "\n" + livePaletteCss + "\n</style>")
}

// RenderAppShell renders the Svelte SPA host document. It deliberately reuses
// the same live-document boot path as the existing Go-rendered pages so the
// installed PWA keeps its viewport, theme, WCO, font, and service-worker
// behavior while routes migrate into Svelte incrementally.
//
// bootstrap, when non-empty, is the base64 session payload the SPA reads to
// render the first paint without fetching /api/session — see the session route.
//
// prefix is the request's validated X-Forwarded-Prefix ("" = origin root): the
// shell's asset URLs are rewritten and window.__PI_BASE_PATH__ is declared so
// the SPA prefixes its own fetch/SSE/navigation URLs. An empty prefix renders
// byte-identical output to the pre-prefix shell.
func RenderAppShell(w io.Writer, bootstrap, prefix string) error {
	scriptSrc := template.HTMLEscapeString(PrefixPath(prefix, appScriptPath))
	preload := template.HTML(`<link rel="modulepreload" href="` + scriptSrc + `">`)
	bootstrapTag := template.HTML("")
	if bootstrap != "" {
		// base64 only (A-Za-z0-9+/=), so it cannot break out of the script tag.
		bootstrapTag = template.HTML(`<script id="pi-session-bootstrap" type="application/json">` + template.HTMLEscapeString(bootstrap) + `</script>`)
	}
	data := struct {
		LiveDocumentStart template.HTML
		ThemeBoot         template.HTML
		Bootstrap         template.HTML
		AppScript         template.HTML
		ServiceWorker     template.HTML
		LiveDocumentEnd   template.HTML
	}{
		LiveDocumentStart: template.HTML(renderLiveDocumentStart(liveDocumentData{
			Title:   "pi-web",
			Preload: preload,
			Styles:  appStylesheets(),
			Prefix:  prefix,
		})),
		ThemeBoot:       liveThemeBootScript(),
		Bootstrap:       bootstrapTag,
		AppScript:       template.HTML(`<script type="module" src="` + scriptSrc + `"></script>`),
		ServiceWorker:   liveServiceWorkerScript(prefix),
		LiveDocumentEnd: liveDocumentEnd(),
	}
	var buf bytes.Buffer
	if err := appTmpl.Execute(&buf, data); err != nil {
		return err
	}
	_, err := w.Write(buf.Bytes())
	return err
}
