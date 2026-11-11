// Integration coverage for a reverse-proxy mount: the module-level BASE_PATH is
// read from window.__PI_BASE_PATH__ at import time, so these tests set the
// global BEFORE dynamically importing the modules under test. Each vitest file
// runs with a fresh module registry, so the imports below see the declared
// prefix even though other test files import the same modules at root.
import { describe, expect, it, vi } from 'vitest';

describe('mounted deployment (window.__PI_BASE_PATH__ = /agent/pi)', () => {
  it('readBasePath picks up the server-declared prefix', async () => {
    window.__PI_BASE_PATH__ = '/agent/pi';
    const { readBasePath } = await import('./base-path.js');
    expect(readBasePath(window)).toBe('/agent/pi');
  });

  it('getJSON/postJSON fetch through the mount', async () => {
    window.__PI_BASE_PATH__ = '/agent/pi';
    const { getJSON, postJSON } = await import('./api.js');
    const fetchImpl = vi.fn(
      async () => new Response(JSON.stringify({ ok: true }), { status: 200 }),
    );
    await getJSON('/api/sessions', { fetchImpl });
    await postJSON('/api/new-session', { path: '/tmp' }, { fetchImpl });
    expect(fetchImpl.mock.calls[0][0]).toBe('/agent/pi/api/sessions');
    expect(fetchImpl.mock.calls[1][0]).toBe('/agent/pi/api/new-session');
  });

  it('navigate pushes prefixed URLs onto history', async () => {
    window.__PI_BASE_PATH__ = '/agent/pi';
    const { navigate } = await import('./navigation.js');
    const pushState = vi.fn();
    navigate('/session?id=abc', { windowImpl: { history: { pushState } } });
    expect(pushState).toHaveBeenCalledWith({}, '', '/agent/pi/session?id=abc');
  });

  it('navigate does not double-prefix an already-prefixed path', async () => {
    window.__PI_BASE_PATH__ = '/agent/pi';
    const { navigate } = await import('./navigation.js');
    const pushState = vi.fn();
    navigate('/agent/pi/settings', { windowImpl: { history: { pushState } } });
    expect(pushState).toHaveBeenCalledWith({}, '', '/agent/pi/settings');
  });

  it('the router strips the prefix when matching routes', async () => {
    window.__PI_BASE_PATH__ = '/agent/pi';
    const { stripBasePath } = await import('./base-path.js');
    expect(stripBasePath('/agent/pi/session')).toBe('/session');
    expect(stripBasePath('/agent/pi')).toBe('/');
    expect(stripBasePath('/agent/pi/')).toBe('/');
  });
});
