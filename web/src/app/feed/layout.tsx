'use client';

import React from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import styles from './layout.module.css';

export default function FeedLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();

  const handleLogout = (e: React.FormEvent) => {
    e.preventDefault();
    localStorage.removeItem('canopy_token');
    localStorage.removeItem('canopy_user');
    router.push('/login');
  };

  const isActive = (path: string) => pathname?.startsWith(path);

  return (
    <div className={styles.layoutContainer}>
      {/* Left Sidebar */}
      <aside className={styles.sidebar}>
        <div>
          <Link href="/feed/home" className={styles.logo}>
            <span>C</span>
            <span className={styles.logoText}>ANOPY</span>
          </Link>

          <nav className={styles.navLinks}>
            <Link 
              href="/feed/home" 
              className={`${styles.navItem} ${isActive('/feed/home') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <path d="M3 9.5L12 3l9 6.5V20a1 1 0 0 1-1 1h-5v-6h-4v6H4a1 1 0 0 1-1-1V9.5z" />
              </svg>
              <span className={styles.navLabel}>Home Feed</span>
            </Link>

            <Link 
              href="/feed/public" 
              className={`${styles.navItem} ${isActive('/feed/public') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="9" />
                <path d="M3.6 9h16.8M3.6 15h16.8M12 3a15.3 15.3 0 0 1 4 9 15.3 15.3 0 0 1-4 9 15.3 15.3 0 0 1-4-9 15.3 15.3 0 0 1 4-9z" />
              </svg>
              <span className={styles.navLabel}>Public Timeline</span>
            </Link>

            <Link 
              href="/feed/messages" 
              className={`${styles.navItem} ${isActive('/feed/messages') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10z" />
              </svg>
              <span className={styles.navLabel}>Direct Messages</span>
            </Link>

            <Link 
              href="/feed/notifications" 
              className={`${styles.navItem} ${isActive('/feed/notifications') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9M13.73 21a2 2 0 0 1-3.46 0" />
              </svg>
              <span className={styles.navLabel}>Notifications</span>
            </Link>

            <Link 
              href="/feed/profile/me" 
              className={`${styles.navItem} ${isActive('/feed/profile') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                <circle cx="12" cy="7" r="4" />
              </svg>
              <span className={styles.navLabel}>My Profile</span>
            </Link>

            <Link 
              href="/feed/theme-editor" 
              className={`${styles.navItem} ${isActive('/feed/theme-editor') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
              </svg>
              <span className={styles.navLabel}>Theme Editor</span>
            </Link>

            <Link 
              href="/admin" 
              className={`${styles.navItem} ${isActive('/admin') ? styles.navItemActive : ''}`}
            >
              <svg className={styles.navIcon} viewBox="0 0 24 24">
                <path d="M12 22c5.523 0 10-4.477 10-10S17.523 2 12 2 2 6.477 2 12s4.477 10 10 10z" />
                <path d="M12 16v-4M12 8h.01" />
              </svg>
              <span className={styles.navLabel}>Admin Panel</span>
            </Link>
          </nav>
        </div>

        <div className={styles.logoutSection}>
          <button onClick={handleLogout} className={styles.logoutBtn}>
            <svg className={styles.navIcon} viewBox="0 0 24 24">
              <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9" />
            </svg>
            <span className={styles.logoutLabel}>Log Out</span>
          </button>
        </div>
      </aside>

      {/* Main Panel */}
      <main className={styles.mainPanel}>
        {children}
      </main>

      {/* Right Widget Panel */}
      <aside className={styles.widgetsPanel}>
        <div className={styles.widgetCard}>
          <h3 className={styles.widgetTitle}>Server Status</h3>
          <div className={styles.serverStatus}>
            <div className={styles.statusRow}>
              <span>Instance Name</span>
              <span className={styles.statusVal}>canopy.social</span>
            </div>
            <div className={styles.statusRow}>
              <span>Federated Domains</span>
              <span className={styles.statusVal}>42 blocked</span>
            </div>
            <div className={styles.statusRow}>
              <span>Total Local Users</span>
              <span className={styles.statusVal}>2,840</span>
            </div>
            <div className={styles.statusRow}>
              <span>Uptime</span>
              <span className={styles.statusVal}>100% stable</span>
            </div>
          </div>
        </div>

        <div className={styles.widgetCard}>
          <h3 className={styles.widgetTitle}>Trends</h3>
          <div className={styles.trendList}>
            <div className={styles.trendItem}>
              <Link href="/feed/public" className={styles.trendTag}>#golang</Link>
              <span className={styles.trendCount}>1,240 posts</span>
            </div>
            <div className={styles.trendItem}>
              <Link href="/feed/public" className={styles.trendTag}>#nextjs14</Link>
              <span className={styles.trendCount}>842 posts</span>
            </div>
            <div className={styles.trendItem}>
              <Link href="/feed/public" className={styles.trendTag}>#monochrome</Link>
              <span className={styles.trendCount}>310 posts</span>
            </div>
            <div className={styles.trendItem}>
              <Link href="/feed/public" className={styles.trendTag}>#minimalism</Link>
              <span className={styles.trendCount}>194 posts</span>
            </div>
          </div>
        </div>
      </aside>
    </div>
  );
}
