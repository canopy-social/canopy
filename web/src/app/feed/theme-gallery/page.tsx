'use client';

import React, { useState } from 'react';

interface GalleryTheme {
  id: string;
  name: string;
  author: string;
  colors: { background: string; surface: string; text: string; accent: string };
  downloads: number;
  remixes: number;
}

const GALLERY_THEMES: GalleryTheme[] = [
  { id: 't1', name: 'Midnight', author: '@sumi_dev', colors: { background: '#0a0a0a', surface: '#1a1a1a', text: '#e0e0e0', accent: '#ffffff' }, downloads: 342, remixes: 18 },
  { id: 't2', name: 'Paper White', author: '@alice', colors: { background: '#faf9f7', surface: '#ffffff', text: '#1a1a1a', accent: '#000000' }, downloads: 567, remixes: 31 },
  { id: 't3', name: 'Terminal', author: '@dev_hacker', colors: { background: '#0d1117', surface: '#161b22', text: '#c9d1d9', accent: '#58a6ff' }, downloads: 891, remixes: 45 },
  { id: 't4', name: 'Warm Gray', author: '@minimal_ist', colors: { background: '#f5f0eb', surface: '#ebe5de', text: '#2d2926', accent: '#8b7355' }, downloads: 234, remixes: 12 },
  { id: 't5', name: 'Ink', author: '@typographer', colors: { background: '#1c1c1e', surface: '#2c2c2e', text: '#f2f2f7', accent: '#ff9f0a' }, downloads: 445, remixes: 22 },
  { id: 't6', name: 'Bone', author: '@cream_soda', colors: { background: '#f8f4ef', surface: '#eee8e0', text: '#3d3833', accent: '#b5651d' }, downloads: 178, remixes: 8 },
  { id: 't7', name: 'Slate', author: '@stone_cold', colors: { background: '#1e293b', surface: '#334155', text: '#e2e8f0', accent: '#94a3b8' }, downloads: 620, remixes: 29 },
  { id: 't8', name: 'Snow', author: '@winter_dev', colors: { background: '#ffffff', surface: '#f9fafb', text: '#111827', accent: '#6366f1' }, downloads: 1024, remixes: 67 },
];

export default function ThemeGalleryPage() {
  const [search, setSearch] = useState('');
  const [sortBy, setSortBy] = useState<'downloads' | 'remixes' | 'name'>('downloads');

  const filtered = GALLERY_THEMES
    .filter(t => t.name.toLowerCase().includes(search.toLowerCase()) || t.author.toLowerCase().includes(search.toLowerCase()))
    .sort((a, b) => {
      if (sortBy === 'downloads') return b.downloads - a.downloads;
      if (sortBy === 'remixes') return b.remixes - a.remixes;
      return a.name.localeCompare(b.name);
    });

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '24px', paddingBottom: '16px', borderBottom: '1px solid var(--border)' }}>
        <h1 style={{ fontFamily: 'var(--font-heading)', fontSize: '18px', fontWeight: 700, letterSpacing: '-0.02em' }}>THEME GALLERY</h1>
        <div style={{ display: 'flex', gap: '8px' }}>
          <input
            type="text"
            className="input-field"
            placeholder="Search themes..."
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            style={{ width: '200px', fontSize: '12px', padding: '6px 10px' }}
          />
          <select className="input-field" style={{ width: '140px', fontSize: '12px' }} value={sortBy} onChange={(e) => setSortBy(e.target.value as 'downloads' | 'remixes' | 'name')}>
            <option value="downloads">Most Used</option>
            <option value="remixes">Most Remixed</option>
            <option value="name">Alphabetical</option>
          </select>
        </div>
      </div>

      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))', gap: '16px' }}>
        {filtered.map(theme => (
          <div key={theme.id} style={{
            border: '1px solid var(--border)',
            borderRadius: '8px',
            overflow: 'hidden',
            transition: 'box-shadow 0.15s ease',
            cursor: 'pointer',
          }}
          onMouseEnter={(e) => { (e.currentTarget as HTMLDivElement).style.boxShadow = '0 0 0 2px var(--text-primary)'; }}
          onMouseLeave={(e) => { (e.currentTarget as HTMLDivElement).style.boxShadow = 'none'; }}
          >
            <div style={{
              height: '120px',
              backgroundColor: theme.colors.background,
              padding: '16px',
              display: 'flex',
              flexDirection: 'column',
              justifyContent: 'space-between',
            }}>
              <div style={{ display: 'flex', gap: '6px' }}>
                <div style={{ width: '20px', height: '20px', borderRadius: '4px', backgroundColor: theme.colors.surface, border: '1px solid rgba(128,128,128,0.2)' }} />
                <div style={{ width: '20px', height: '20px', borderRadius: '4px', backgroundColor: theme.colors.accent }} />
              </div>
              <div>
                <div style={{ color: theme.colors.text, fontSize: '18px', fontWeight: 700 }}>{theme.name}</div>
                <div style={{ color: theme.colors.text, opacity: 0.6, fontSize: '12px' }}>Preview Text</div>
              </div>
            </div>
            <div style={{ padding: '12px 16px', borderTop: '1px solid var(--border)' }}>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <span style={{ fontSize: '13px', fontWeight: 600 }}>{theme.author}</span>
                <div style={{ display: 'flex', gap: '12px', fontSize: '11px', color: 'var(--text-secondary)' }}>
                  <span>{theme.downloads} uses</span>
                  <span>{theme.remixes} remixes</span>
                </div>
              </div>
              <div style={{ display: 'flex', gap: '6px', marginTop: '10px' }}>
                <button className="btn-primary" style={{ flex: 1, fontSize: '11px', padding: '5px 0' }}>APPLY</button>
                <button className="btn-secondary" style={{ flex: 1, fontSize: '11px', padding: '5px 0' }}>REMIX</button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {filtered.length === 0 && (
        <div style={{ textAlign: 'center', padding: '60px 0', color: 'var(--text-secondary)' }}>
          No themes match your search.
        </div>
      )}
    </div>
  );
}
