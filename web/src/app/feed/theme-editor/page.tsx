'use client';

import React, { useState, useCallback } from 'react';
import styles from './editor.module.css';

interface ColorPalette {
  background: string;
  surface: string;
  text_primary: string;
  text_secondary: string;
  accent: string;
  accent_text: string;
  border: string;
  link: string;
}

interface FontConfig {
  family: string;
  size: number;
  weight: number;
  line_height: number;
  letter_spacing: number;
}

interface FontPalette {
  body: FontConfig;
  heading: FontConfig;
  mono: FontConfig;
  display: FontConfig;
}

interface ThemeVersion {
  id: string;
  label?: string;
  auto_saved: boolean;
  created_at: string;
}

const ALLOWED_FONTS = [
  'Inter', 'Outfit', 'Roboto', 'Open Sans', 'Lato', 'Montserrat',
  'Poppins', 'Raleway', 'Playfair Display', 'Merriweather',
  'Source Code Pro', 'Fira Code', 'JetBrains Mono', 'IBM Plex Sans',
  'IBM Plex Mono', 'DM Sans', 'DM Serif Display', 'Space Grotesk',
  'Space Mono', 'Archivo', 'Sora', 'Inconsolata', 'system-ui',
  'monospace', 'serif', 'sans-serif',
];

const DEFAULT_COLORS: ColorPalette = {
  background: '#ffffff',
  surface: '#f8f8f8',
  text_primary: '#111111',
  text_secondary: '#555555',
  accent: '#0066ff',
  accent_text: '#ffffff',
  border: '#e0e0e0',
  link: '#0066ff',
};

const DEFAULT_FONTS: FontPalette = {
  body: { family: 'Inter', size: 16, weight: 400, line_height: 1.6, letter_spacing: 0 },
  heading: { family: 'Outfit', size: 24, weight: 700, line_height: 1.2, letter_spacing: -0.02 },
  mono: { family: 'JetBrains Mono', size: 14, weight: 400, line_height: 1.5, letter_spacing: 0 },
  display: { family: 'Outfit', size: 48, weight: 900, line_height: 1.0, letter_spacing: -0.03 },
};

const COLOR_LABELS: Record<keyof ColorPalette, string> = {
  background: 'Background',
  surface: 'Surface',
  text_primary: 'Text',
  text_secondary: 'Muted',
  accent: 'Accent',
  accent_text: 'Accent Text',
  border: 'Border',
  link: 'Link',
};

type PreviewMode = 'profile' | 'post' | 'essay';

export default function ThemeEditorPage() {
  const [colors, setColors] = useState<ColorPalette>(DEFAULT_COLORS);
  const [fonts, setFonts] = useState<FontPalette>(DEFAULT_FONTS);
  const [bgType, setBgType] = useState('color');
  const [pageMaxWidth, setPageMaxWidth] = useState(800);
  const [pagePadding, setPagePadding] = useState(24);
  const [showFollowerCount, setShowFollowerCount] = useState(true);
  const [showFollowingCount, setShowFollowingCount] = useState(true);
  const [gardenMode, setGardenMode] = useState(false);
  const [previewMode, setPreviewMode] = useState<PreviewMode>('profile');
  const [versions, setVersions] = useState<ThemeVersion[]>([
    { id: 'v1', label: 'Initial Setup', auto_saved: false, created_at: '2026-05-20T10:00:00Z' },
    { id: 'v2', auto_saved: true, created_at: '2026-05-21T08:30:00Z' },
    { id: 'v3', label: 'Dark Monochrome', auto_saved: false, created_at: '2026-05-21T09:15:00Z' },
  ]);
  const [saveLabel, setSaveLabel] = useState('');
  const [hasUnsaved, setHasUnsaved] = useState(false);

  const updateColor = useCallback((key: keyof ColorPalette, value: string) => {
    setColors(prev => ({ ...prev, [key]: value }));
    setHasUnsaved(true);
  }, []);

  const updateFont = useCallback((role: keyof FontPalette, field: keyof FontConfig, value: string | number) => {
    setFonts(prev => ({
      ...prev,
      [role]: { ...prev[role], [field]: value },
    }));
    setHasUnsaved(true);
  }, []);

  const handleSaveVersion = () => {
    if (!saveLabel.trim()) return;
    const newVersion: ThemeVersion = {
      id: `v_${Date.now()}`,
      label: saveLabel,
      auto_saved: false,
      created_at: new Date().toISOString(),
    };
    setVersions(prev => [newVersion, ...prev]);
    setSaveLabel('');
    setHasUnsaved(false);
  };

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const handleRestoreVersion = (id: string) => {
    setColors(DEFAULT_COLORS);
    setFonts(DEFAULT_FONTS);
    setHasUnsaved(false);
  };

  const previewStyle: React.CSSProperties = {
    backgroundColor: colors.background,
    color: colors.text_primary,
    fontFamily: `"${fonts.body.family}", system-ui`,
    fontSize: `${fonts.body.size}px`,
    fontWeight: fonts.body.weight,
    lineHeight: fonts.body.line_height,
    maxWidth: `${pageMaxWidth}px`,
    padding: `${pagePadding}px`,
    margin: '0 auto',
    transition: 'all 0.3s ease',
  };

  return (
    <div>
      <div className={styles.editorHeader}>
        <span className={styles.editorTitle}>THEME EDITOR</span>
        <div style={{ display: 'flex', gap: '8px' }}>
          <input
            type="text"
            className="input-field"
            placeholder="Version label..."
            value={saveLabel}
            onChange={(e) => setSaveLabel(e.target.value)}
            style={{ width: '180px', fontSize: '12px', padding: '6px 10px' }}
          />
          <button className="btn-primary" onClick={handleSaveVersion} style={{ fontSize: '12px', padding: '6px 16px' }}>
            {hasUnsaved ? 'SAVE *' : 'SAVE'}
          </button>
        </div>
      </div>

      <div className={styles.editorContainer}>
        {/* Left Panel — Editor Controls */}
        <div className={styles.editorPanel}>
          <div className={styles.sectionTitle}>Colors</div>
          <div className={styles.colorGrid}>
            {(Object.keys(COLOR_LABELS) as Array<keyof ColorPalette>).map((key) => (
              <div key={key} className={styles.colorField}>
                <span className={styles.colorLabel}>{COLOR_LABELS[key]}</span>
                <div className={styles.colorInputRow}>
                  <div className={styles.colorSwatch} style={{ backgroundColor: colors[key], position: 'relative' }}>
                    <input
                      type="color"
                      value={colors[key]}
                      onChange={(e) => updateColor(key, e.target.value)}
                      style={{ opacity: 0, position: 'absolute', top: 0, left: 0, width: '100%', height: '100%', cursor: 'pointer' }}
                    />
                  </div>
                  <input
                    type="text"
                    className={styles.colorHexInput}
                    value={colors[key]}
                    onChange={(e) => updateColor(key, e.target.value)}
                  />
                </div>
              </div>
            ))}
          </div>

          <div className={styles.sectionTitle}>Typography</div>
          {(['body', 'heading', 'mono', 'display'] as Array<keyof FontPalette>).map((role) => (
            <div key={role} className={styles.fontField}>
              <div className={styles.fontFieldLabel}>{role.toUpperCase()}</div>
              <select
                className={styles.fontSelect}
                value={fonts[role].family}
                onChange={(e) => updateFont(role, 'family', e.target.value)}
              >
                {ALLOWED_FONTS.map(f => <option key={f} value={f}>{f}</option>)}
              </select>
              <div className={styles.rangeRow}>
                <span className={styles.rangeLabel}>Size</span>
                <input type="range" className={styles.rangeInput} min={8} max={72} value={fonts[role].size}
                  onChange={(e) => updateFont(role, 'size', parseInt(e.target.value))} />
                <span className={styles.rangeValue}>{fonts[role].size}px</span>
              </div>
              <div className={styles.rangeRow}>
                <span className={styles.rangeLabel}>Weight</span>
                <input type="range" className={styles.rangeInput} min={100} max={900} step={100} value={fonts[role].weight}
                  onChange={(e) => updateFont(role, 'weight', parseInt(e.target.value))} />
                <span className={styles.rangeValue}>{fonts[role].weight}</span>
              </div>
            </div>
          ))}

          <div className={styles.sectionTitle}>Layout</div>
          <div className={styles.layoutField}>
            <div className={styles.layoutLabel}>Background Type</div>
            <select className={styles.bgTypeSelect} value={bgType} onChange={(e) => { setBgType(e.target.value); setHasUnsaved(true); }}>
              <option value="color">Solid Color</option>
              <option value="gradient">Gradient</option>
              <option value="image">Image</option>
            </select>
          </div>

          <div className={styles.rangeRow}>
            <span className={styles.rangeLabel}>Max Width</span>
            <input type="range" className={styles.rangeInput} min={400} max={1600} step={50} value={pageMaxWidth}
              onChange={(e) => { setPageMaxWidth(parseInt(e.target.value)); setHasUnsaved(true); }} />
            <span className={styles.rangeValue}>{pageMaxWidth}</span>
          </div>

          <div className={styles.rangeRow}>
            <span className={styles.rangeLabel}>Padding</span>
            <input type="range" className={styles.rangeInput} min={0} max={64} value={pagePadding}
              onChange={(e) => { setPagePadding(parseInt(e.target.value)); setHasUnsaved(true); }} />
            <span className={styles.rangeValue}>{pagePadding}</span>
          </div>

          <div className={styles.sectionTitle}>Display</div>
          <div className={styles.toggleRow}>
            <span className={styles.toggleLabel}>Show Follower Count</span>
            <label className={styles.toggle}>
              <input type="checkbox" checked={showFollowerCount} onChange={(e) => { setShowFollowerCount(e.target.checked); setHasUnsaved(true); }} />
              <span className={styles.toggleTrack} />
            </label>
          </div>
          <div className={styles.toggleRow}>
            <span className={styles.toggleLabel}>Show Following Count</span>
            <label className={styles.toggle}>
              <input type="checkbox" checked={showFollowingCount} onChange={(e) => { setShowFollowingCount(e.target.checked); setHasUnsaved(true); }} />
              <span className={styles.toggleTrack} />
            </label>
          </div>
          <div className={styles.toggleRow}>
            <span className={styles.toggleLabel}>Garden Mode</span>
            <label className={styles.toggle}>
              <input type="checkbox" checked={gardenMode} onChange={(e) => { setGardenMode(e.target.checked); setHasUnsaved(true); }} />
              <span className={styles.toggleTrack} />
            </label>
          </div>
        </div>

        {/* Center Panel — Live Preview */}
        <div className={styles.previewPanel}>
          <div className={styles.previewToolbar}>
            {(['profile', 'post', 'essay'] as PreviewMode[]).map((mode) => (
              <button
                key={mode}
                className={`${styles.previewBtn} ${previewMode === mode ? styles.previewBtnActive : ''}`}
                onClick={() => setPreviewMode(mode)}
              >
                {mode.toUpperCase()}
              </button>
            ))}
          </div>
          <div className={styles.previewFrame}>
            <div className={styles.previewContent} style={previewStyle}>
              {previewMode === 'profile' && (
                <>
                  <div className={styles.previewProfile}>
                    <div className={styles.previewAvatar} style={{ backgroundColor: colors.accent, color: colors.accent_text }}>S</div>
                    <div className={styles.previewDisplayName} style={{ fontFamily: `"${fonts.heading.family}", system-ui`, fontWeight: fonts.heading.weight }}>
                      Sumi Dev
                    </div>
                    <div className={styles.previewHandle} style={{ color: colors.text_secondary }}>@sumi_dev</div>
                    <div className={styles.previewBio} style={{ color: colors.text_secondary }}>
                      Building decentralized social infrastructure. Canopy contributor. Minimalist design enthusiast.
                    </div>
                    {showFollowerCount && <span style={{ fontSize: '13px', color: colors.text_secondary, marginRight: '16px' }}><strong style={{ color: colors.text_primary }}>1,247</strong> Followers</span>}
                    {showFollowingCount && <span style={{ fontSize: '13px', color: colors.text_secondary }}><strong style={{ color: colors.text_primary }}>312</strong> Following</span>}
                  </div>
                  <div className={styles.previewPost} style={{ borderColor: colors.border, backgroundColor: colors.surface }}>
                    <div className={styles.previewPostHeader}>
                      <div className={styles.previewPostAvatar} style={{ backgroundColor: colors.accent, color: colors.accent_text }}>S</div>
                      <div>
                        <div className={styles.previewPostName} style={{ fontFamily: `"${fonts.body.family}", system-ui` }}>Sumi Dev</div>
                      </div>
                    </div>
                    <div className={styles.previewPostBody}>
                      Just shipped the theme engine backend with CSS generation, validation, and versioning. 20 tests passing.
                    </div>
                  </div>
                </>
              )}
              {previewMode === 'post' && (
                <div>
                  <div className={styles.previewPost} style={{ borderColor: colors.border, backgroundColor: colors.surface }}>
                    <div className={styles.previewPostHeader}>
                      <div className={styles.previewPostAvatar} style={{ backgroundColor: colors.accent, color: colors.accent_text }}>A</div>
                      <div>
                        <div className={styles.previewPostName}>Alice Johnson</div>
                      </div>
                    </div>
                    <div className={styles.previewPostBody}>The monochrome design system really shines with this theme engine. Clean borders, no gradients, just pure structure.</div>
                  </div>
                  <div className={styles.previewPost} style={{ borderColor: colors.border, backgroundColor: colors.surface }}>
                    <div className={styles.previewPostHeader}>
                      <div className={styles.previewPostAvatar} style={{ backgroundColor: '#333', color: '#fff' }}>D</div>
                      <div>
                        <div className={styles.previewPostName}>Dave H</div>
                      </div>
                    </div>
                    <div className={styles.previewPostBody}>Federation support coming soon. The architecture is solid.</div>
                  </div>
                </div>
              )}
              {previewMode === 'essay' && (
                <div>
                  <h1 style={{ fontFamily: `"${fonts.display.family}", system-ui`, fontSize: `${fonts.display.size}px`, fontWeight: fonts.display.weight, lineHeight: fonts.display.line_height, letterSpacing: `${fonts.display.letter_spacing}em`, marginBottom: '16px' }}>
                    On Decentralized Social Design
                  </h1>
                  <p style={{ color: colors.text_secondary, fontSize: '14px', marginBottom: '24px' }}>
                    Published May 21, 2026 · 8 min read
                  </p>
                  <p style={{ marginBottom: '16px' }}>
                    The promise of federated social networks lies not just in data ownership, but in the freedom to define your own aesthetic identity. Unlike centralized platforms where every profile looks the same, decentralized systems can empower users to express themselves through design.
                  </p>
                  <h2 style={{ fontFamily: `"${fonts.heading.family}", system-ui`, fontSize: `${fonts.heading.size}px`, fontWeight: fonts.heading.weight, marginTop: '32px', marginBottom: '12px' }}>
                    The Role of Themes
                  </h2>
                  <p>
                    A theme engine must balance creative freedom with safety constraints. Color validation, font allowlisting, and CSS scoping prevent abuse while enabling genuine self-expression.
                  </p>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Right Panel — Version History */}
        <div className={styles.historyPanel}>
          <div className={styles.sectionTitle}>Version History</div>
          {versions.map((v) => (
            <div key={v.id} className={styles.historyItem} onClick={() => handleRestoreVersion(v.id)}>
              <div className={styles.historyLabel}>
                {v.label || 'Auto-save'}
                {v.auto_saved && <span className={styles.historyBadge}>AUTO</span>}
              </div>
              <div className={styles.historyTime}>
                {new Date(v.created_at).toLocaleString()}
              </div>
            </div>
          ))}

          <div className={styles.actionBar}>
            <button className="btn-secondary" style={{ flex: 1, fontSize: '11px' }}>EXPORT</button>
            <button className="btn-secondary" style={{ flex: 1, fontSize: '11px' }}>IMPORT</button>
          </div>
        </div>
      </div>
    </div>
  );
}
