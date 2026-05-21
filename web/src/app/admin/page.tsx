'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import styles from './admin.module.css';

interface UserAccount {
  id: string;
  username: string;
  email: string;
  role: string;
  status: 'active' | 'silenced' | 'suspended';
}

interface ModerationReport {
  id: string;
  reporter: string;
  targetUser: string;
  reason: string;
  status: 'open' | 'resolved';
}

interface DomainBlock {
  id: string;
  domain: string;
  severity: 'silence' | 'suspend';
}

interface InviteCode {
  id: string;
  token: string;
  maxUses: number;
  uses: number;
  expired: boolean;
}

export default function AdminDashboard() {
  const [activeTab, setActiveTab] = useState<'accounts' | 'reports' | 'blocks' | 'invites' | 'settings'>('accounts');

  // Accounts state
  const [accounts, setAccounts] = useState<UserAccount[]>([
    { id: 'u1', username: 'sumi_dev', email: 'sumi@canopy.social', role: 'admin', status: 'active' },
    { id: 'u2', username: 'linus', email: 'linus@kernel.org', role: 'moderator', status: 'active' },
    { id: 'u3', username: 'spambot99', email: 'spam@botnet.cc', role: 'user', status: 'suspended' },
    { id: 'u4', username: 'trollface', email: 'troll@anonymous.org', role: 'user', status: 'silenced' }
  ]);

  // Reports state
  const [reports, setReports] = useState<ModerationReport[]>([
    { id: 'r1', reporter: 'alice', targetUser: 'trollface', reason: 'Repeated offensive comments under dev posts', status: 'open' },
    { id: 'r2', reporter: 'linus', targetUser: 'spambot99', reason: 'Federated spam advertisements', status: 'resolved' }
  ]);

  // Domain blocks state
  const [blocks, setBlocks] = useState<DomainBlock[]>([
    { id: 'b1', domain: 'spamdomain.net', severity: 'suspend' },
    { id: 'b2', domain: 'annoyingads.com', severity: 'silence' }
  ]);
  const [newDomain, setNewDomain] = useState('');
  const [blockSeverity, setBlockSeverity] = useState<'silence' | 'suspend'>('suspend');

  // Invites state
  const [invites, setInvites] = useState<InviteCode[]>([
    { id: 'i1', token: 'CANOPY-DEV-2026', maxUses: 10, uses: 4, expired: false },
    { id: 'i2', token: 'CANOPY-BETA-TEST', maxUses: 5, uses: 5, expired: true }
  ]);
  const [inviteUses, setInviteUses] = useState(5);

  // Server settings state
  const [settings, setSettings] = useState({
    openRegistration: false,
    requireInvites: true,
    maxImageMB: 10,
    maxVideoMB: 50,
  });

  const toggleUserStatus = (id: string, action: 'suspend' | 'unsuspend' | 'silence' | 'unsilence') => {
    setAccounts(accounts.map(acc => {
      if (acc.id === id) {
        if (action === 'suspend') return { ...acc, status: 'suspended' };
        if (action === 'unsuspend') return { ...acc, status: 'active' };
        if (action === 'silence') return { ...acc, status: 'silenced' };
        if (action === 'unsilence') return { ...acc, status: 'active' };
      }
      return acc;
    }));
  };

  const resolveReport = (id: string) => {
    setReports(reports.map(rep => {
      if (rep.id === id) return { ...rep, status: 'resolved' };
      return rep;
    }));
  };

  const handleAddBlock = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDomain.trim()) return;
    const newBlock: DomainBlock = {
      id: `b_${Math.random().toString(36).substr(2, 9)}`,
      domain: newDomain,
      severity: blockSeverity
    };
    setBlocks([...blocks, newBlock]);
    setNewDomain('');
  };

  const handleRemoveBlock = (id: string) => {
    setBlocks(blocks.filter(b => b.id !== id));
  };

  const handleGenerateInvite = () => {
    const randomToken = `CANOPY-${Math.random().toString(36).substr(2, 4).toUpperCase()}-${Math.random().toString(36).substr(2, 4).toUpperCase()}`;
    const newInvite: InviteCode = {
      id: `inv_${Math.random().toString(36).substr(2, 9)}`,
      token: randomToken,
      maxUses: inviteUses,
      uses: 0,
      expired: false
    };
    setInvites([newInvite, ...invites]);
  };

  return (
    <div className={`${styles.adminContainer} animate-fade-in`}>
      <div className={styles.adminHeader}>
        <h1 className={styles.adminTitle}>Admin Control</h1>
        <Link href="/feed/home" className={styles.backBtn}>
          BACK TO FEED
        </Link>
      </div>

      {/* Tabs */}
      <div className={styles.tabBar}>
        <button onClick={() => setActiveTab('accounts')} className={`${styles.tabBtn} ${activeTab === 'accounts' ? styles.tabActive : ''}`}>Accounts</button>
        <button onClick={() => setActiveTab('reports')} className={`${styles.tabBtn} ${activeTab === 'reports' ? styles.tabActive : ''}`}>Reports</button>
        <button onClick={() => setActiveTab('blocks')} className={`${styles.tabBtn} ${activeTab === 'blocks' ? styles.tabActive : ''}`}>Instance Blocks</button>
        <button onClick={() => setActiveTab('invites')} className={`${styles.tabBtn} ${activeTab === 'invites' ? styles.tabActive : ''}`}>Invite Codes</button>
        <button onClick={() => setActiveTab('settings')} className={`${styles.tabBtn} ${activeTab === 'settings' ? styles.tabActive : ''}`}>Server Settings</button>
      </div>

      {/* Accounts Tab */}
      {activeTab === 'accounts' && (
        <div className={styles.sectionCard}>
          <h2 className={styles.sectionTitle}>User Accounts Management</h2>
          <table className={styles.adminTable}>
            <thead>
              <tr>
                <th>Username</th>
                <th>Email</th>
                <th>Role</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {accounts.map(acc => (
                <tr key={acc.id}>
                  <td><strong>@{acc.username}</strong></td>
                  <td>{acc.email}</td>
                  <td>{acc.role}</td>
                  <td>
                    <span className={`${styles.badge} ${
                      acc.status === 'active' ? styles.badgeActive : 
                      acc.status === 'suspended' ? styles.badgeSuspended : styles.badgeSilenced
                    }`}>
                      {acc.status}
                    </span>
                  </td>
                  <td>
                    <div style={{ display: 'flex', gap: '8px' }}>
                      {acc.status === 'active' ? (
                        <>
                          <button onClick={() => toggleUserStatus(acc.id, 'silence')} className={styles.actionButton}>Silence</button>
                          <button onClick={() => toggleUserStatus(acc.id, 'suspend')} className={styles.actionButton}>Suspend</button>
                        </>
                      ) : acc.status === 'silenced' ? (
                        <>
                          <button onClick={() => toggleUserStatus(acc.id, 'unsilence')} className={styles.actionButton}>Unsilence</button>
                          <button onClick={() => toggleUserStatus(acc.id, 'suspend')} className={styles.actionButton}>Suspend</button>
                        </>
                      ) : (
                        <button onClick={() => toggleUserStatus(acc.id, 'unsuspend')} className={styles.actionButton}>Unsuspend</button>
                      )}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Reports Tab */}
      {activeTab === 'reports' && (
        <div className={styles.sectionCard}>
          <h2 className={styles.sectionTitle}>Content Reports & Moderation</h2>
          <table className={styles.adminTable}>
            <thead>
              <tr>
                <th>ID</th>
                <th>Reporter</th>
                <th>Offending User</th>
                <th>Reason</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {reports.map(rep => (
                <tr key={rep.id}>
                  <td><code>{rep.id}</code></td>
                  <td>@{rep.reporter}</td>
                  <td><strong>@{rep.targetUser}</strong></td>
                  <td>{rep.reason}</td>
                  <td>
                    <span className={`${styles.badge} ${rep.status === 'resolved' ? styles.badgeActive : styles.badgeSuspended}`}>
                      {rep.status}
                    </span>
                  </td>
                  <td>
                    {rep.status === 'open' && (
                      <button onClick={() => resolveReport(rep.id)} className={styles.actionButton}>
                        RESOLVE
                      </button>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Domain Blocks Tab */}
      {activeTab === 'blocks' && (
        <div className={styles.sectionCard}>
          <h2 className={styles.sectionTitle}>Domain Instance Blocks</h2>
          <form onSubmit={handleAddBlock} className={styles.formRow}>
            <div style={{ display: 'flex', gap: '12px', width: '100%' }}>
              <input
                type="text"
                className="input-field"
                placeholder="domain.com to block"
                value={newDomain}
                onChange={(e) => setNewDomain(e.target.value)}
                required
              />
              <select 
                className="input-field" 
                style={{ width: '150px' }}
                value={blockSeverity}
                onChange={(e) => setBlockSeverity(e.target.value as 'silence' | 'suspend')}
              >
                <option value="suspend">Suspend</option>
                <option value="silence">Silence</option>
              </select>
            </div>
            <button type="submit" className="btn-primary">BLOCK</button>
          </form>

          <table className={styles.adminTable}>
            <thead>
              <tr>
                <th>Blocked Domain</th>
                <th>Severity</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {blocks.map(b => (
                <tr key={b.id}>
                  <td><code>{b.domain}</code></td>
                  <td>
                    <span className={`${styles.badge} ${b.severity === 'suspend' ? styles.badgeSuspended : styles.badgeSilenced}`}>
                      {b.severity}
                    </span>
                  </td>
                  <td>
                    <button onClick={() => handleRemoveBlock(b.id)} className={styles.actionButton}>
                      REMOVE BLOCK
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Invite Codes Tab */}
      {activeTab === 'invites' && (
        <div className={styles.sectionCard}>
          <h2 className={styles.sectionTitle}>Registration Invite Tokens</h2>
          <div className={styles.formRow}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <span style={{ fontSize: '14px', fontWeight: '600' }}>Token Limit (uses):</span>
              <input
                type="number"
                className="input-field"
                style={{ width: '100px' }}
                value={inviteUses}
                onChange={(e) => setInviteUses(parseInt(e.target.value) || 1)}
                min={1}
              />
            </div>
            <button onClick={handleGenerateInvite} className="btn-primary">
              GENERATE NEW CODE
            </button>
          </div>

          <table className={styles.adminTable}>
            <thead>
              <tr>
                <th>Invite Token</th>
                <th>Uses Allowed</th>
                <th>Uses Taken</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              {invites.map(inv => (
                <tr key={inv.id}>
                  <td><code>{inv.token}</code></td>
                  <td>{inv.maxUses}</td>
                  <td>{inv.uses}</td>
                  <td>
                    <span className={`${styles.badge} ${inv.expired ? styles.badgeSilenced : styles.badgeActive}`}>
                      {inv.expired ? 'expired' : 'active'}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {/* Settings Tab */}
      {activeTab === 'settings' && (
        <div className={styles.sectionCard}>
          <h2 className={styles.sectionTitle}>Global Instance Server Settings</h2>
          <div style={{ display: 'flex', flexDirection: 'column', gap: '20px' }}>
            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <input
                id="openReg"
                type="checkbox"
                checked={settings.openRegistration}
                onChange={(e) => setSettings({ ...settings, openRegistration: e.target.checked })}
                style={{ width: '18px', height: '18px', accentColor: 'var(--accent)' }}
              />
              <label htmlFor="openReg" style={{ fontSize: '15px', fontWeight: '500' }}>Allow Open Local Registrations</label>
            </div>

            <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
              <input
                id="reqInv"
                type="checkbox"
                checked={settings.requireInvites}
                onChange={(e) => setSettings({ ...settings, requireInvites: e.target.checked })}
                style={{ width: '18px', height: '18px', accentColor: 'var(--accent)' }}
              />
              <label htmlFor="reqInv" style={{ fontSize: '15px', fontWeight: '500' }}>Strictly Require Invite Code Tokens</label>
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px', marginTop: '12px' }}>
              <span style={{ fontSize: '13px', fontWeight: '600', textTransform: 'uppercase', color: 'var(--foreground-subtle)' }}>Max Image Size Limit (MB)</span>
              <input
                type="number"
                className="input-field"
                style={{ width: '150px' }}
                value={settings.maxImageMB}
                onChange={(e) => setSettings({ ...settings, maxImageMB: parseInt(e.target.value) || 0 })}
              />
            </div>

            <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
              <span style={{ fontSize: '13px', fontWeight: '600', textTransform: 'uppercase', color: 'var(--foreground-subtle)' }}>Max Video Size Limit (MB)</span>
              <input
                type="number"
                className="input-field"
                style={{ width: '150px' }}
                value={settings.maxVideoMB}
                onChange={(e) => setSettings({ ...settings, maxVideoMB: parseInt(e.target.value) || 0 })}
              />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
