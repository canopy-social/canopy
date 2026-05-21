'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import styles from '../home/feed.module.css';

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
  isLocal: boolean;
  isLiked?: boolean;
  isBoosted?: boolean;
}

export default function PublicTimelinePage() {
  const [filter, setFilter] = useState<'local' | 'federated'>('federated');
  const [posts, setPosts] = useState<Post[]>([
    {
      id: 'post_pub_1',
      authorName: 'Sumi Dev',
      authorHandle: '@sumi_dev',
      avatar: 'S',
      content: 'Local Post: Designing monochrome inputs with solid border variables makes style overriding super elegant. Pure CSS all the way.',
      createdAt: '2 hrs ago',
      likes: 15,
      boosts: 3,
      replies: 1,
      isLocal: true,
    },
    {
      id: 'post_pub_2',
      authorName: 'Griesemer',
      authorHandle: '@gri@go.dev',
      avatar: 'G',
      content: 'Federated Post: Go 1.22 compiler optimizations have reduced memory footprint on big AST parsing. Really nice stability improvements.',
      createdAt: '4 hrs ago',
      likes: 120,
      boosts: 45,
      replies: 12,
      isLocal: false,
    },
    {
      id: 'post_pub_3',
      authorName: 'Alice Johnson',
      authorHandle: '@alice@canopy.social',
      avatar: 'A',
      content: 'Local Post: Successfully connected our timelines to a local Redis server. Fan-out write is incredibly rapid.',
      createdAt: '1 day ago',
      likes: 8,
      boosts: 2,
      replies: 0,
      isLocal: true,
    },
    {
      id: 'post_pub_4',
      authorName: 'Dan Abramov',
      authorHandle: '@dan@react.dev',
      avatar: 'D',
      content: 'Federated Post: Server Actions in Next.js 14 really simplify data fetching without writing standard REST boilerplate endpoints.',
      createdAt: '2 days ago',
      likes: 210,
      boosts: 78,
      replies: 34,
      isLocal: false,
    }
  ]);

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

  const filteredPosts = posts.filter(post => {
    if (filter === 'local') return post.isLocal;
    return true; // Federated shows everything
  });

  return (
    <div className="animate-fade-in">
      <div className={styles.feedHeader}>
        <h2 className={styles.feedTitle}>Public Timeline</h2>

        <div className={styles.publicToggle}>
          <button 
            onClick={() => setFilter('local')} 
            className={`${styles.toggleBtn} ${filter === 'local' ? styles.toggleBtnActive : ''}`}
          >
            Local
          </button>
          <button 
            onClick={() => setFilter('federated')} 
            className={`${styles.toggleBtn} ${filter === 'federated' ? styles.toggleBtnActive : ''}`}
          >
            Federated
          </button>
        </div>
      </div>

      {/* Posts list */}
      <div className={styles.postList}>
        {filteredPosts.map((post) => (
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
