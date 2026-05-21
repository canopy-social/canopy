'use client';

import React, { useState, useRef, useCallback } from 'react';
import styles from './layout-editor.module.css';

interface Widget {
  id: string;
  type: string;
  label: string;
  x: number;
  y: number;
  width: number;
  height: number;
}

const WIDGET_TYPES = [
  { type: 'text', label: 'Text Block', icon: '¶' },
  { type: 'image', label: 'Image', icon: '◻' },
  { type: 'bio', label: 'Bio Card', icon: '⊡' },
  { type: 'posts', label: 'Recent Posts', icon: '≡' },
  { type: 'links', label: 'Link List', icon: '⊞' },
  { type: 'stats', label: 'Stats', icon: '▤' },
];

const WIDGET_CONTENT: Record<string, string> = {
  text: 'Editable text block. Double-click to modify content.',
  image: 'Drag an image here or click to upload.',
  bio: 'Your bio card with avatar, display name, and description.',
  posts: 'A feed of your recent posts displayed inline.',
  links: 'A curated list of external links.',
  stats: 'Follower count, post count, and join date.',
};

export default function LayoutEditorPage() {
  const [widgets, setWidgets] = useState<Widget[]>([
    { id: 'w1', type: 'bio', label: 'Bio Card', x: 40, y: 40, width: 320, height: 180 },
    { id: 'w2', type: 'posts', label: 'Recent Posts', x: 400, y: 40, width: 360, height: 280 },
    { id: 'w3', type: 'links', label: 'Link List', x: 40, y: 260, width: 320, height: 160 },
  ]);

  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [draggingId, setDraggingId] = useState<string | null>(null);
  const dragOffset = useRef({ x: 0, y: 0 });
  const canvasRef = useRef<HTMLDivElement>(null);

  const addWidget = (type: string, label: string) => {
    const newWidget: Widget = {
      id: `w_${Date.now()}`,
      type,
      label,
      x: 40 + Math.random() * 200,
      y: 40 + Math.random() * 200,
      width: 280,
      height: 160,
    };
    setWidgets(prev => [...prev, newWidget]);
    setSelectedId(newWidget.id);
  };

  const deleteWidget = (id: string) => {
    setWidgets(prev => prev.filter(w => w.id !== id));
    if (selectedId === id) setSelectedId(null);
  };

  const handleMouseDown = useCallback((e: React.MouseEvent, id: string) => {
    e.preventDefault();
    setDraggingId(id);
    setSelectedId(id);
    const widget = widgets.find(w => w.id === id);
    if (!widget || !canvasRef.current) return;
    const rect = canvasRef.current.getBoundingClientRect();
    dragOffset.current = {
      x: e.clientX - rect.left - widget.x,
      y: e.clientY - rect.top - widget.y,
    };
  }, [widgets]);

  const handleMouseMove = useCallback((e: React.MouseEvent) => {
    if (!draggingId || !canvasRef.current) return;
    const rect = canvasRef.current.getBoundingClientRect();
    const x = Math.max(0, e.clientX - rect.left - dragOffset.current.x);
    const y = Math.max(0, e.clientY - rect.top - dragOffset.current.y);
    setWidgets(prev => prev.map(w => w.id === draggingId ? { ...w, x, y } : w));
  }, [draggingId]);

  const handleMouseUp = useCallback(() => {
    setDraggingId(null);
  }, []);

  const selectedWidget = widgets.find(w => w.id === selectedId);

  const updateWidgetProperty = (field: keyof Widget, value: number) => {
    if (!selectedId) return;
    setWidgets(prev => prev.map(w => w.id === selectedId ? { ...w, [field]: value } : w));
  };

  return (
    <div className={styles.layoutEditor}>
      <div className={styles.editorToolbar}>
        <span className={styles.toolbarTitle}>LAYOUT EDITOR</span>
        <div className={styles.toolbarActions}>
          <button className="btn-secondary" style={{ fontSize: '12px', padding: '6px 16px' }}>RESET</button>
          <button className="btn-primary" style={{ fontSize: '12px', padding: '6px 16px' }}>SAVE LAYOUT</button>
        </div>
      </div>

      <div className={styles.componentPalette}>
        {WIDGET_TYPES.map(wt => (
          <button key={wt.type} className={styles.paletteItem} onClick={() => addWidget(wt.type, wt.label)}>
            <span className={styles.paletteIcon}>{wt.icon}</span>
            {wt.label}
          </button>
        ))}
      </div>

      <div
        ref={canvasRef}
        className={styles.canvas}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
        onClick={() => setSelectedId(null)}
      >
        <div className={styles.canvasGrid} />
        {widgets.map(w => (
          <div
            key={w.id}
            className={`${styles.widget} ${draggingId === w.id ? styles.widgetDragging : ''}`}
            style={{
              left: w.x,
              top: w.y,
              width: w.width,
              height: w.height,
              outline: selectedId === w.id ? '2px solid var(--text-primary)' : 'none',
            }}
            onMouseDown={(e) => { e.stopPropagation(); handleMouseDown(e, w.id); }}
            onClick={(e) => { e.stopPropagation(); setSelectedId(w.id); }}
          >
            <div className={styles.widgetHeader}>
              <span>{w.label}</span>
              <button className={styles.widgetDeleteBtn} onClick={(e) => { e.stopPropagation(); deleteWidget(w.id); }}>×</button>
            </div>
            <div className={styles.widgetBody}>
              {WIDGET_CONTENT[w.type] || 'Widget content'}
            </div>
            <div className={styles.resizeHandle} />
          </div>
        ))}
      </div>

      {selectedWidget && (
        <div className={styles.propertiesPanel}>
          <div className={styles.propertiesTitle}>Widget Properties — {selectedWidget.label}</div>
          <div className={styles.propertyRow}>
            <span className={styles.propertyLabel}>X</span>
            <input type="number" className={styles.propertyInput} value={Math.round(selectedWidget.x)}
              onChange={(e) => updateWidgetProperty('x', parseInt(e.target.value) || 0)} />
          </div>
          <div className={styles.propertyRow}>
            <span className={styles.propertyLabel}>Y</span>
            <input type="number" className={styles.propertyInput} value={Math.round(selectedWidget.y)}
              onChange={(e) => updateWidgetProperty('y', parseInt(e.target.value) || 0)} />
          </div>
          <div className={styles.propertyRow}>
            <span className={styles.propertyLabel}>Width</span>
            <input type="number" className={styles.propertyInput} value={selectedWidget.width}
              onChange={(e) => updateWidgetProperty('width', parseInt(e.target.value) || 100)} />
          </div>
          <div className={styles.propertyRow}>
            <span className={styles.propertyLabel}>Height</span>
            <input type="number" className={styles.propertyInput} value={selectedWidget.height}
              onChange={(e) => updateWidgetProperty('height', parseInt(e.target.value) || 60)} />
          </div>
        </div>
      )}
    </div>
  );
}
