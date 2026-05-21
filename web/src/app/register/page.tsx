'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import styles from '../login/login.module.css';

export default function RegisterPage() {
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [inviteCode, setInviteCode] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!username || !email || !password) {
      setError('Please fill in all required fields.');
      return;
    }
    setError('');
    setIsLoading(true);

    try {
      // Real or mock registration handling
      localStorage.setItem('canopy_token', 'mock_jwt_token');
      localStorage.setItem('canopy_user', JSON.stringify({
        id: 'user_01hkabc123',
        username: username,
        email: email,
        role: 'user'
      }));

      // Simulate network request delay
      setTimeout(() => {
        setIsLoading(false);
        router.push('/feed/home');
      }, 800);
    } catch {
      setIsLoading(false);
      setError('Failed to create account. Please check inputs.');
    }
  };

  return (
    <div className={styles.authWrapper}>
      <div className={`${styles.authCard} animate-slide-up`}>
        <div className={styles.logoHeader}>
          <h1>CANOPY</h1>
          <p className="text-subtle">CREATE YOUR INDIVIDUAL ACCOUNT</p>
        </div>

        {error && (
          <div className={styles.errorBox}>
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className={styles.authForm}>
          <div className={styles.inputGroup}>
            <label htmlFor="username">Username</label>
            <input
              id="username"
              type="text"
              className="input-field"
              placeholder="e.g. sumi_dev"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>

          <div className={styles.inputGroup}>
            <label htmlFor="email">Email Address</label>
            <input
              id="email"
              type="email"
              className="input-field"
              placeholder="name@domain.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div className={styles.inputGroup}>
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              className="input-field"
              placeholder="••••••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          <div className={styles.inputGroup}>
            <label htmlFor="inviteCode">Invite Code (Optional)</label>
            <input
              id="inviteCode"
              type="text"
              className="input-field"
              placeholder="e.g. CANOPY-XXXX-XXXX"
              value={inviteCode}
              onChange={(e) => setInviteCode(e.target.value)}
            />
          </div>

          <button 
            type="submit" 
            className="btn-primary" 
            style={{ width: '100%', marginTop: '8px' }}
            disabled={isLoading}
          >
            {isLoading ? 'CREATING...' : 'REGISTER'}
          </button>
        </form>

        <div className={styles.footerLink}>
          <span>Already have an account? </span>
          <Link href="/login">LOG IN</Link>
        </div>
      </div>
    </div>
  );
}
