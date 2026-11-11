import { describe, expect, it } from 'vitest';
import { readBasePath, stripBasePath, withBasePath } from './base-path.js';

describe('readBasePath', () => {
  it('returns empty for a window without the global', () => {
    expect(readBasePath({})).toBe('');
  });

  it('returns empty for undefined window', () => {
    expect(readBasePath(undefined)).toBe('');
  });

  it('returns empty string as-is', () => {
    expect(readBasePath({ __PI_BASE_PATH__: '' })).toBe('');
  });

  it('keeps a normalized prefix', () => {
    expect(readBasePath({ __PI_BASE_PATH__: '/agent/pi' })).toBe('/agent/pi');
  });

  it('adds a missing leading slash', () => {
    expect(readBasePath({ __PI_BASE_PATH__: 'agent/pi' })).toBe('/agent/pi');
  });

  it('strips trailing slashes', () => {
    expect(readBasePath({ __PI_BASE_PATH__: '/agent/pi/' })).toBe('/agent/pi');
    expect(readBasePath({ __PI_BASE_PATH__: '/agent/pi///' })).toBe('/agent/pi');
  });

  it('collapses a lone slash to empty (root)', () => {
    expect(readBasePath({ __PI_BASE_PATH__: '/' })).toBe('');
  });
});

describe('withBasePath', () => {
  it('passes URLs through when no base path is set', () => {
    expect(withBasePath('/api/sessions', '')).toBe('/api/sessions');
  });

  it('prefixes root-absolute URLs', () => {
    expect(withBasePath('/api/sessions', '/agent/pi')).toBe('/agent/pi/api/sessions');
  });

  it('keeps query strings', () => {
    expect(withBasePath('/session?id=a b', '/agent/pi')).toBe('/agent/pi/session?id=a b');
  });

  it('does not double-prefix an already-prefixed URL', () => {
    expect(withBasePath('/agent/pi/session', '/agent/pi')).toBe('/agent/pi/session');
    expect(withBasePath('/agent/pi', '/agent/pi')).toBe('/agent/pi');
  });

  it('prefixes a path that merely starts with the base string', () => {
    expect(withBasePath('/agent/pix', '/agent/pi')).toBe('/agent/pi/agent/pix');
  });

  it('leaves full URLs and protocol-relative URLs alone', () => {
    expect(withBasePath('https://example.com/x', '/agent/pi')).toBe('https://example.com/x');
    expect(withBasePath('//example.com/x', '/agent/pi')).toBe('//example.com/x');
  });

  it('leaves relative paths alone', () => {
    expect(withBasePath('icon.svg', '/agent/pi')).toBe('icon.svg');
  });

  it('prefixes root', () => {
    expect(withBasePath('/', '/agent/pi')).toBe('/agent/pi/');
  });
});

describe('stripBasePath', () => {
  it('passes pathnames through when no base path is set', () => {
    expect(stripBasePath('/session', '')).toBe('/session');
  });

  it('strips the prefix from real pathnames', () => {
    expect(stripBasePath('/agent/pi/session', '/agent/pi')).toBe('/session');
  });

  it('maps the mount root to the app root', () => {
    expect(stripBasePath('/agent/pi', '/agent/pi')).toBe('/');
    expect(stripBasePath('/agent/pi/', '/agent/pi')).toBe('/');
  });

  it('passes through pathnames outside the mount', () => {
    expect(stripBasePath('/other/app', '/agent/pi')).toBe('/other/app');
  });

  it('does not strip a path that merely starts with the base string', () => {
    expect(stripBasePath('/agent/pix/session', '/agent/pi')).toBe('/agent/pix/session');
  });
});
