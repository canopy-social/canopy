'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import styles from './notifications.module.css';

interface NotificationItem {
  id: string;
  type: 'like' | 'boost' | 'follow' | 'reply';
  actorName: string;
  actorHandle: string;
  postContent?: string;
  time: string;
}

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<NotificationItem[]>([
    {
      id: 'n1',
      type: 'like',
      actorName: 'Alice Johnson',
      actorHandle: 'alice',
      postContent: 'Bootstrapping the Next.js 14 frontend scaffolding for Canopy! Extremely clean minimalist layout with 100% Vanilla CSS.',
      time: '12 mins ago'
    },
    {
      id: 'n2',
      type: 'boost',
      actorName: 'Linus B',
      actorHandle: 'linus',
      postContent: 'Tested the Go moderation service and admin APIs today. The zero-comment code standard in internal/accounts works flawlessly.',
      time: '1 hr ago'
    },
    {
      id: 'n3',
      type: 'follow',
      actorName: 'Dave H',
      actorHandle: 'dave',
      time: '2 hrs ago'
    },
    {
      id: 'n4',
      type: 'reply',
      actorName: 'Alice Johnson',
      actorHandle: 'alice',
      postContent: 'Excellent verification, Linus! The backend compilation benchmarks on local machines are down by 15%.',
      time: '1 day ago'
    }
  ]);

  const handleClearAll = () => {
    setNotifications([]);
  };

  return (
    <div className="animate-fade-in">
      <div className={styles.notifHeader}>
        <h2 className={styles.notifTitle}>Notifications</h2>
        {notifications.length > 0 && (
          <button onClick={handleClearAll} className={styles.clearBtn}>
            CLEAR ALL
          </button>
        )}
      </div>

      {notifications.length === 0 ? (
        <div className="card" style={{ textAlign: 'center', padding: '40px', color: 'var(--foreground-subtle)' }}>
          <span>Your inbox is completely clear.</span>
        </div>
      ) : (
        <div className={styles.notifList}>
          {notifications.map((notif) => (
            <div key={notif.id} className={styles.notifCard}>
              <div className={styles.notifIconWrapper}>
                {notif.type === 'like' && (
                  <svg className={styles.notifIcon} viewBox="0 0 24 24">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
                  </svg>
                )}
                {notif.type === 'boost' && (
                  <svg className={styles.notifIcon} viewBox="0 0 24 24">
                    <path d="M17 1l4 4-4 4M21 5H9a5 5 0 0 0-5 5v3M7 23l-4-4 4-4M3 19h12a5 5 0 0 0 5-5v-3" />
                  </svg>
                )}
                {notif.type === 'follow' && (
                  <svg className={styles.notifIcon} viewBox="0 0 24 24">
                    <path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
                    <circle cx="8.5" cy="7" r="4" />
                  </svg>
                )}
                {notif.type === 'reply' && (
                  <svg className={styles.notifIcon} viewBox="0 0 24 24">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10z" />
                  </svg>
                )}
              </div>

              <div className={styles.notifDetails}>
                <span className={styles.notifMsg}>
                  <Link href={`/feed/profile/${notif.actorHandle}`} className={styles.actorName}>
                    {notif.actorName}
                  </Link>{' '}
                  {notif.type === 'like' && 'liked your post'}
                  {notif.type === 'boost' && 'boosted your post'}
                  {notif.type === 'follow' && 'followed you'}
                  {notif.type === 'reply' && 'replied to you'}
                </span>

                {notif.postContent && (
                  <p className={styles.notifSnippet}>
                    {notif.postContent.length > 80 
                      ? `${notif.postContent.substring(0, 80)}...` 
                      : notif.postContent}
                  </p>
                )}

                <span className={styles.notifTime}>{notif.time}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
