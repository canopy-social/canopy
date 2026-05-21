'use client';

import React, { useState } from 'react';

import feedStyles from '../../home/feed.module.css';
import styles from './post.module.css';

interface Comment {
  id: string;
  authorName: string;
  authorHandle: string;
  avatar: string;
  content: string;
  createdAt: string;
}

export default function PostThreadPage({ params }: { params: { id: string } }) {
  const [replyText, setReplyText] = useState('');
  
  // Simulated thread dataset
  const mainPost = {
    id: params.id,
    authorName: 'Linus B',
    authorHandle: '@linus@kernel.org',
    avatar: 'L',
    content: 'Tested the Go moderation service and admin APIs today. The zero-comment code standard in internal/accounts works flawlessly. Looking forward to compiling this on the build farm.',
    createdAt: 'May 21, 2026 • 10:14 AM',
    likes: 84,
    boosts: 32,
    repliesCount: 2,
    isLiked: false,
    isBoosted: false
  };

  const parentPost = {
    id: 'parent_999',
    authorName: 'Sumi Dev',
    authorHandle: '@sumi_dev',
    avatar: 'S',
    content: 'Just deployed the new accounts module to Canopy staging. All checks are fully green!',
    createdAt: 'May 20, 2026 • 9:02 PM',
  };

  const [replies, setReplies] = useState<Comment[]>([
    {
      id: 'r1',
      authorName: 'Alice Johnson',
      authorHandle: '@alice@canopy.social',
      avatar: 'A',
      content: 'Excellent verification, Linus! The backend compilation benchmarks on local machines are down by 15%.',
      createdAt: '1 hr ago'
    },
    {
      id: 'r2',
      authorName: 'Dave H',
      authorHandle: '@dave@fediverse.org',
      avatar: 'D',
      content: 'Nice monochrome design on the new frontend scaffold. Simple and high contrast.',
      createdAt: '30 mins ago'
    }
  ]);

  const handleReplySubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!replyText.trim()) return;

    const newReply: Comment = {
      id: `rep_${Math.random().toString(36).substr(2, 9)}`,
      authorName: 'Sumi Dev',
      authorHandle: '@sumi_dev',
      avatar: 'S',
      content: replyText,
      createdAt: 'Just now'
    };

    setReplies([...replies, newReply]);
    setReplyText('');
  };

  return (
    <div className="animate-fade-in">
      <div className={feedStyles.feedHeader}>
        <h2 className={feedStyles.feedTitle}>Thread</h2>
      </div>

      <div className={styles.threadContainer}>
        {/* Parent Post (if any) */}
        {parentPost && (
          <div className={styles.parentWrapper}>
            <div className={feedStyles.postCard} style={{ borderBottom: 'none' }}>
              <div className={feedStyles.avatar}>{parentPost.avatar}</div>
              <div className={feedStyles.postContent}>
                <div className={feedStyles.postHeader}>
                  <div className={feedStyles.authorInfo}>
                    <span className={feedStyles.authorName}>{parentPost.authorName}</span>
                    <span className={feedStyles.authorHandle}>{parentPost.authorHandle}</span>
                  </div>
                </div>
                <div className={feedStyles.postBody}>{parentPost.content}</div>
              </div>
            </div>
            <div className={styles.threadLine} />
          </div>
        )}

        {/* Focused Main Post */}
        <div className={styles.focusedPost}>
          <div className={feedStyles.postCard} style={{ borderColor: 'var(--border-active)' }}>
            <div className={feedStyles.avatar}>{mainPost.avatar}</div>
            <div className={feedStyles.postContent}>
              <div className={feedStyles.postHeader}>
                <div className={feedStyles.authorInfo}>
                  <span className={feedStyles.authorName}>{mainPost.authorName}</span>
                  <span className={feedStyles.authorHandle}>{mainPost.authorHandle}</span>
                </div>
              </div>
              <div className={styles.focusedBody}>{mainPost.content}</div>
              <div className={styles.focusedTime}>{mainPost.createdAt}</div>
              
              <div className={styles.statsDivider} />
              
              <div className={styles.focusedStats}>
                <span><strong>{mainPost.boosts}</strong> Boosts</span>
                <span><strong>{mainPost.likes}</strong> Likes</span>
              </div>
            </div>
          </div>
        </div>

        {/* Reply Composer */}
        <form onSubmit={handleReplySubmit} className={styles.replyComposer}>
          <textarea
            placeholder="Write your reply..."
            className={styles.replyTextarea}
            value={replyText}
            onChange={(e) => setReplyText(e.target.value)}
            required
          />
          <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: '12px' }}>
            <button type="submit" className="btn-primary">REPLY</button>
          </div>
        </form>

        {/* Replies List */}
        <div className={styles.repliesList}>
          <h3 className={styles.repliesHeader}>Replies</h3>
          {replies.map((reply) => (
            <div key={reply.id} className={feedStyles.postCard} style={{ borderLeft: '3px solid var(--border)' }}>
              <div className={feedStyles.avatar}>{reply.avatar}</div>
              <div className={feedStyles.postContent}>
                <div className={feedStyles.postHeader}>
                  <div className={feedStyles.authorInfo}>
                    <span className={feedStyles.authorName}>{reply.authorName}</span>
                    <span className={feedStyles.authorHandle}>{reply.authorHandle}</span>
                  </div>
                  <span className={feedStyles.postTime}>{reply.createdAt}</span>
                </div>
                <div className={feedStyles.postBody}>{reply.content}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
