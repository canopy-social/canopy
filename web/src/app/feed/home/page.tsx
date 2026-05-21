'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import styles from './feed.module.css';

interface Post {
  id: string;
  authorName: string;
  authorHandle: string;
  avatar: string;
  content: string;
  createdAt: string;
  likes: number;
  boosts: number;
  replies: number;
  isLiked?: boolean;
  isBoosted?: boolean;
}

export default function HomeFeedPage() {
  const [content, setContent] = useState('');
  const [posts, setPosts] = useState<Post[]>([
    {
      id: 'post_01hkapc888',
      authorName: 'Sumi Dev',
      authorHandle: '@sumi_dev',
      avatar: 'S',
      content: 'Bootstrapping the Next.js 14 frontend scaffolding for Canopy! Extremely clean minimalist layout with 100% Vanilla CSS. Black, grey, and white monochrome aesthetic feels incredibly premium.',
      createdAt: '2 hrs ago',
      likes: 12,
      boosts: 4,
      replies: 2,
    },
    {
      id: 'post_01hkapc999',
      authorName: 'Linus B',
      authorHandle: '@linus@kernel.org',
      avatar: 'L',
      content: 'Tested the Go moderation service and admin APIs today. The zero-comment code standard in internal/accounts works flawlessly. Looking forward to compiling this on the build farm.',
      createdAt: '5 hrs ago',
      likes: 84,
      boosts: 32,
      replies: 9,
    },
    {
      id: 'post_01hkapc000',
      authorName: 'Alice Johnson',
      authorHandle: '@alice@canopy.social',
      avatar: 'A',
      content: 'Just created an instance block rule for spamdomain.net to protect our local feed. Admin panel settings are super responsive.',
      createdAt: '1 day ago',
      likes: 5,
      boosts: 1,
      replies: 0,
    }
  ]);

  const handlePostSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!content.trim()) return;

    const newPost: Post = {
      id: `post_${Math.random().toString(36).substr(2, 9)}`,
      authorName: 'Sumi Dev',
      authorHandle: '@sumi_dev',
      avatar: 'S',
      content: content,
      createdAt: 'Just now',
      likes: 0,
      boosts: 0,
      replies: 0,
    };

    setPosts([newPost, ...posts]);
    setContent('');
  };

  const handleLike = (id: string) => {
    setPosts(posts.map(post => {
      if (post.id === id) {
        return {
          ...post,
          likes: post.isLiked ? post.likes - 1 : post.likes + 1,
          isLiked: !post.isLiked
        };
      }
      return post;
    }));
  };

  const handleBoost = (id: string) => {
    setPosts(posts.map(post => {
      if (post.id === id) {
        return {
          ...post,
          boosts: post.isBoosted ? post.boosts - 1 : post.boosts + 1,
          isBoosted: !post.isBoosted
        };
      }
      return post;
    }));
  };

  return (
    <div className="animate-fade-in">
      <div className={styles.feedHeader}>
        <h2 className={styles.feedTitle}>Home Feed</h2>
      </div>

      {/* Post Composer */}
      <form onSubmit={handlePostSubmit} className={styles.composer}>
        <textarea
          className={styles.composerTextarea}
          placeholder="What's happening?"
          maxLength={500}
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
        <div className={styles.composerFooter}>
          <span className={styles.charCount}>
            {content.length} / 500
          </span>
          <button type="submit" className="btn-primary" disabled={!content.trim()}>
            POST
          </button>
        </div>
      </form>

      {/* Posts list */}
      <div className={styles.postList}>
        {posts.map((post) => (
          <article key={post.id} className={styles.postCard}>
            <div className={styles.avatar}>{post.avatar}</div>
            <div className={styles.postContent}>
              <div className={styles.postHeader}>
                <div className={styles.authorInfo}>
                  <Link href={`/feed/profile/${post.authorHandle.replace('@', '')}`} className={styles.authorName}>
                    {post.authorName}
                  </Link>
                  <span className={styles.authorHandle}>{post.authorHandle}</span>
                </div>
                <span className={styles.postTime}>{post.createdAt}</span>
              </div>
              <div className={styles.postBody}>{post.content}</div>
              <div className={styles.postActions}>
                {/* Reply */}
                <Link href={`/feed/posts/${post.id}`} className={styles.actionBtn}>
                  <svg className={styles.actionIcon} viewBox="0 0 24 24">
                    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10z" />
                  </svg>
                  <span>{post.replies}</span>
                </Link>

                {/* Boost */}
                <button 
                  onClick={() => handleBoost(post.id)} 
                  className={`${styles.actionBtn} ${post.isBoosted ? styles.actionBtnActive : ''}`}
                >
                  <svg className={`${styles.actionIcon} ${post.isBoosted ? styles.actionIconFilled : ''}`} viewBox="0 0 24 24">
                    <path d="M17 1l4 4-4 4M21 5H9a5 5 0 0 0-5 5v3M7 23l-4-4 4-4M3 19h12a5 5 0 0 0 5-5v-3" />
                  </svg>
                  <span>{post.boosts}</span>
                </button>

                {/* Like */}
                <button 
                  onClick={() => handleLike(post.id)} 
                  className={`${styles.actionBtn} ${post.isLiked ? styles.actionBtnActive : ''}`}
                >
                  <svg className={`${styles.actionIcon} ${post.isLiked ? styles.actionIconFilled : ''}`} viewBox="0 0 24 24">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z" />
                  </svg>
                  <span>{post.likes}</span>
                </button>
              </div>
            </div>
          </article>
        ))}
      </div>
    </div>
  );
}
