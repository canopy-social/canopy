'use client';

import React, { useState, useRef, useCallback } from 'react';

interface GardenItem {
  id: string;
  type: 'post' | 'note' | 'image' | 'link';
  content: string;
  x: number;
  y: number;
  width: number;
}

const INITIAL_ITEMS: GardenItem[] = [
  { id: 'g1', type: 'post', content: 'Just shipped the theme engine. 20 tests passing, CSS gen pipeline solid.', x: 60, y: 60, width: 260 },
  { id: 'g2', type: 'note', content: 'TODO: Add gradient presets to the theme editor', x: 380, y: 100, width: 220 },
  { id: 'g3', type: 'image', content: 'Profile photo placeholder', x: 120, y: 320, width: 160 },
  { id: 'g4', type: 'link', content: 'github.com/canopy-social/canopy', x: 400, y: 340, width: 240 },
  { id: 'g5', type: 'post', content: 'Decentralized social is about owning your space. Garden mode makes it spatial.', x: 60, y: 500, width: 300 },
];

const TYPE_STYLES: Record<string, { border: string; label: string }> = {
  post: { border: '1px solid var(--border)', label: 'POST' },
  note: { border: '1px dashed var(--text-secondary)', label: 'NOTE' },
  image: { border: '2px solid var(--border)', label: 'IMAGE' },
  link: { border: '1px solid var(--text-secondary)', label: 'LINK' },
};

export default function GardenModePage() {
  const [items, setItems] = useState<GardenItem[]>(INITIAL_ITEMS);
  const [draggingId, setDraggingId] = useState<string | null>(null);
  const canvasRef = useRef<HTMLDivElement>(null);
  const dragOffset = useRef({ x: 0, y: 0 });

  const handleMouseDown = useCallback((e: React.MouseEvent, id: string) => {
    e.preventDefault();
    setDraggingId(id);
    const item = items.find(i => i.id === id);
    if (!item || !canvasRef.current) return;
    const rect = canvasRef.current.getBoundingClientRect();
    dragOffset.current = { x: e.clientX - rect.left - item.x, y: e.clientY - rect.top - item.y };
  }, [items]);

  const handleMouseMove = useCallback((e: React.MouseEvent) => {
    if (!draggingId || !canvasRef.current) return;
    const rect = canvasRef.current.getBoundingClientRect();
    const x = Math.max(0, e.clientX - rect.left - dragOffset.current.x);
    const y = Math.max(0, e.clientY - rect.top - dragOffset.current.y);
    setItems(prev => prev.map(i => i.id === draggingId ? { ...i, x, y } : i));
  }, [draggingId]);

  const handleMouseUp = useCallback(() => { setDraggingId(null); }, []);

  const addItem = (type: GardenItem['type']) => {
    const newItem: GardenItem = {
      id: `g_${Date.now()}`,
      type,
      content: type === 'note' ? 'New note...' : type === 'link' ? 'https://example.com' : type === 'image' ? 'Image placeholder' : 'New post content...',
      x: 100 + Math.random() * 200,
      y: 100 + Math.random() * 200,
      width: type === 'note' ? 180 : 240,
    };
    setItems(prev => [...prev, newItem]);
  };

  return (
    <div style={{ padding: '24px' }}>
      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '20px', paddingBottom: '16px', borderBottom: '1px solid var(--border)' }}>
        <h1 style={{ fontFamily: 'var(--font-heading)', fontSize: '18px', fontWeight: 700, letterSpacing: '-0.02em' }}>GARDEN MODE</h1>
        <div style={{ display: 'flex', gap: '8px' }}>
          <button className="btn-secondary" style={{ fontSize: '11px', padding: '6px 12px' }} onClick={() => addItem('post')}>+ POST</button>
          <button className="btn-secondary" style={{ fontSize: '11px', padding: '6px 12px' }} onClick={() => addItem('note')}>+ NOTE</button>
          <button className="btn-secondary" style={{ fontSize: '11px', padding: '6px 12px' }} onClick={() => addItem('image')}>+ IMAGE</button>
          <button className="btn-secondary" style={{ fontSize: '11px', padding: '6px 12px' }} onClick={() => addItem('link')}>+ LINK</button>
          <button className="btn-primary" style={{ fontSize: '11px', padding: '6px 16px' }}>SAVE</button>
        </div>
      </div>

      <p style={{ fontSize: '13px', color: 'var(--text-secondary)', marginBottom: '16px' }}>
        Drag items freely on your 2D spatial canvas. Your page visitors will see this exact layout.
      </p>

      <div
        ref={canvasRef}
        style={{
          position: 'relative',
          width: '100%',
          minHeight: '700px',
          border: '2px dashed var(--border)',
          borderRadius: '8px',
          background: 'var(--bg-secondary)',
          overflow: 'hidden',
          cursor: draggingId ? 'grabbing' : 'default',
        }}
        onMouseMove={handleMouseMove}
        onMouseUp={handleMouseUp}
        onMouseLeave={handleMouseUp}
      >
        <div style={{
          position: 'absolute', inset: 0, pointerEvents: 'none',
          backgroundImage: 'radial-gradient(circle, var(--border) 1px, transparent 1px)',
          backgroundSize: '24px 24px',
          opacity: 0.4,
        }} />

        {items.map(item => (
          <div
            key={item.id}
            style={{
              position: 'absolute',
              left: item.x,
              top: item.y,
              width: item.width,
              border: TYPE_STYLES[item.type].border,
              borderRadius: '6px',
              background: 'var(--bg-primary)',
              padding: '10px 12px',
              cursor: draggingId === item.id ? 'grabbing' : 'grab',
              userSelect: 'none',
              zIndex: draggingId === item.id ? 100 : 1,
              transition: draggingId === item.id ? 'none' : 'box-shadow 0.15s ease',
              boxShadow: draggingId === item.id ? '0 4px 12px rgba(0,0,0,0.15)' : 'none',
            }}
            onMouseDown={(e) => handleMouseDown(e, item.id)}
          >
            <div style={{ fontSize: '9px', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.08em', color: 'var(--text-secondary)', marginBottom: '6px' }}>
              {TYPE_STYLES[item.type].label}
            </div>
            <div style={{ fontSize: '13px', lineHeight: '1.5', color: 'var(--text-primary)' }}>
              {item.content}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
