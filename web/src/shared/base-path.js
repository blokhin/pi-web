// Runtime base path for reverse-proxy deployments.
//
// A reverse proxy can mount pi-web under a URL prefix (nginx:
// `location /agent/pi/ { proxy_pass http://127.0.0.1:31415/; }`), which strips
// the prefix before the request reaches the server and declares it via the
// X-Forwarded-Prefix header. The Go server validates that header and injects
// `window.__PI_BASE_PATH__` into the served HTML shell — this module is the
// single place the frontend reads it, mirroring the server-side rules in
// internal/ui/prefix.go. An empty/missing value means "served at the origin
// root" and every helper below is a no-op, so root deployments behave
// exactly as before.

/**
 * Normalized base path from a window-like object: leading slash guaranteed,
 * trailing slashes stripped, '' when unset/empty. Mirrors the server's
 * NormalisePrefix contract.
 */
export function readBasePath(windowImpl) {
  const win = windowImpl ?? (typeof window !== 'undefined' ? window : undefined);
  if (!win) return '';
  const raw = win.__PI_BASE_PATH__ ?? '';
  if (!raw) return '';
  const withLead = raw.startsWith('/') ? raw : `/${raw}`;
  return withLead.replace(/\/+$/, '');
}

/** The deployment's base path — the one central constant used everywhere. */
export const BASE_PATH = readBasePath();

function hasBasePathPrefix(url, base) {
  return url === base || url.startsWith(`${base}/`);
}

/**
 * Prefix a root-absolute URL ('/api/…', '/session?id=…') with the deployment
 * base path. Idempotent: a URL that already carries the prefix (e.g. one read
 * back from location.pathname) passes through unchanged, so route-space and
 * browser-space URLs can both flow through the same call. Full and
 * protocol-relative URLs pass through untouched.
 */
export function withBasePath(url, basePath = BASE_PATH) {
  if (!basePath || typeof url !== 'string' || !url.startsWith('/')) return url;
  if (url.startsWith('//') || hasBasePathPrefix(url, basePath)) return url;
  return basePath + url;
}

/**
 * Strip the base path from the browser's real pathname so the SPA router can
 * keep matching unprefixed routes ('/session') under a prefixed deployment.
 * Unprefixed pathnames pass through unchanged.
 */
export function stripBasePath(pathname, basePath = BASE_PATH) {
  if (!basePath || typeof pathname !== 'string') return pathname;
  if (!hasBasePathPrefix(pathname, basePath)) return pathname;
  const stripped = pathname.slice(basePath.length) || '/';
  return stripped.startsWith('/') ? stripped : `/${stripped}`;
}
