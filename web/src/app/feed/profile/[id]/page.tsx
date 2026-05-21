'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import styles from '../profile.module.css';
import feedStyles from '../../home/feed.module.css';

interface UserProfile {
  username: string;
  displayName: string;
  handle: string;
  avatar: string;
  bio: string;
  postsCount: number;
  followingCount: number;
  followersCount: number;
  isFollowing?: boolean;
}

interface Post {
  id: string;
  content: string;
  createdAt: string;
  likes: number;
  boosts: number;
  replies: number;
  isLiked?: boolean;
  isBoosted?: boolean;
}

export default function ProfilePage({ params }: { params: { id: string } }) {
  const isMe = params.id === 'me' || params.id === 'sumi_dev';
  
  // Resolve user info dynamically
  const [profile, setProfile] = useState<UserProfile>(
    isMe ? {
      username: 'sumi_dev',
      displayName: 'Sumi Dev',
      handle: '@sumi_dev@canopy.social',
      avatar: 'S',
      bio: 'Full-stack developer building Canopy. Pair programming with Antigravity AI. Exploring Next.js, pure CSS, and Go systems engineering.',
      postsCount: 18,
      followingCount: 142,
      followersCount: 890,
    } : params.id === 'linus' ? {
      username: 'linus',
      displayName: 'Linus B',
      handle: '@linus@kernel.org',
      avatar: 'L',
      bio: 'I do Git and Linux kernel development. Sometimes I compile things. No comment.',
      postsCount: 2420,
      followingCount: 12,
      followersCount: 95400,
      isFollowing: true,
    } : {
      username: params.id,
      displayName: params.id.charAt(0).toUpperCase() + params.id.slice(1),
      handle: `@${params.id}@canopy.social`,
      avatar: params.id.charAt(0).toUpperCase(),
      bio: `Hello! I am ${params.id}. Welcome to my Canopy feed. Let's federate!`,
      postsCount: 4,
      followingCount: 30,
      followersCount: 15,
      isFollowing: false,
    }
  );

  const [posts, setPosts] = useState<Post[]>([
    {
      id: 'p1',
      content: `Hello fediverse! Just finished setting up my profile page on Canopy. No gradients, no glassmorphism - just pure minimalist layout.`,
      createdAt: '2 hrs ago',
      likes: 12,
      boosts: 3,
      replies: 1
    },
    {
      id: 'p2',
      content: `Standard Go module and Postgres migrations compile beautifully on Windows inside the sandbox.`,
      createdAt: '1 day ago',
      likes: 8,
      boosts: 2,
      replies: 0
    }
  ]);

  const handleFollowToggle = () => {
    setProfile({
      ...profile,
      isFollowing: !profile.isFollowing,
      followersCount: profile.isFollowing ? profile.followersCount - 1 : profile.followersCount + 1
    });
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
    <div className={`${styles.profileContainer} animate-fade-in`}>
      {/* Profile Card Info */}
      <div className={styles.profileCard}>
        <div className={styles.banner} />
        <div className={styles.profileInfoWrapper}>
          <div className={styles.avatarLarge}>{profile.avatar}</div>
          
          <div className={styles.profileHeader}>
            <div className={styles.nameHandle}>
              <h1 className={styles.displayName}>{profile.displayName}</h1>
              <span className={styles.handle}>{profile.handle}</span>
            </div>
            
            {!isMe && (
              <button 
                onClick={handleFollowToggle} 
                className={profile.isFollowing ? 'btn-secondary' : 'btn-primary'}
                style={{ padding: '8px 20px', fontSize: '13px' }}
              >
                {profile.isFollowing ? 'UNFOLLOW' : 'FOLLOW'}
              </button>
            )}
          </div>

          <p className={styles.bio}>{profile.bio}</p>

          <div className={styles.statsRow}>
            <div className={styles.statItem}>
              <span className={styles.statVal}>{profile.postsCount}</span>posts
            </div>
            <div className={styles.statItem}>
              <span className={styles.statVal}>{profile.followingCount}</span>following
            </div>
            <div className={styles.statItem}>
              <span className={styles.statVal}>{profile.followersCount}</span>followers
            </div>
          </div>
        </div>
      </div>

      {/* User Posts Timeline */}
      <div>
        <h3 className={styles.sectionHeader}>Posts</h3>
        <div className={feedStyles.postList} style={{ marginTop: '16px' }}>
          {posts.map((post) => (
            <article key={post.id} className={feedStyles.postCard}>
              <div className={feedStyles.avatar}>{profile.avatar}</div>
              <div className={feedStyles.postContent}>
                <div className={feedStyles.postHeader}>
                  <div className={feedStyles.authorInfo}>
                    <span className={feedStyles.authorName}>{profile.displayName}</span>
                    <span className={feedStyles.authorHandle}>{profile.handle}</span>
                  </div>
                  <span className={feedStyles.postTime}>{post.createdAt}</span>
                </div>
                <div className={feedStyles.postBody}>{post.content}</div>
                <div className={feedStyles.postActions}>
                  {/* Reply */}
                  <Link href={`/feed/posts/${post.id}`} className={feedStyles.actionBtn}>
                    <svg className={feedStyles.actionIcon} viewBox="0 0 24 24">
                      <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2v10z" />
                    </svg>
                    <span>{post.replies}</span>
                  </Link>

                  {/* Boost */}
                  <button 
                    onClick={() => handleBoost(post.id)} 
                    className={`${feedStyles.actionBtn} ${post.isBoosted ? feedStyles.actionBtnActive : ''}`}
                  >
                    <svg className={`${feedStyles.actionIcon} ${post.isBoosted ? feedStyles.actionIconFilled : ''}`} viewBox="0 0 24 24">
                      <path d="M17 1l4 4-4 4M21 5H9a5 5 0 0 0-5 5v3M7 23l-4-4 4-4M3 19h12a5 5 0 0 0 5-5v-3" />
                    </svg>
                    <span>{post.boosts}</span>
                  </button>

                  {/* Like */}
                  <button 
                    onClick={() => handleLike(post.id)} 
                    className={`${feedStyles.actionBtn} ${post.isLiked ? feedStyles.actionBtnActive : ''}`}
                  >
                    <svg className={`${feedStyles.actionIcon} ${post.isLiked ? feedStyles.actionIconFilled : ''}`} viewBox="0 0 24 24">
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
    </div>
  );
}
