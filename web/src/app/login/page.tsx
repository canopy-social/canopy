'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import styles from './login.module.css';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!email || !password) {
      setError('Please fill in all fields.');
      return;
    }
    setError('');
    setIsLoading(true);

    try {
      // Real or mock authentication handling
      // We store mock session details in localStorage for testing
      localStorage.setItem('canopy_token', 'mock_jwt_token');
      localStorage.setItem('canopy_user', JSON.stringify({
        id: 'user_01hkabc123',
        username: 'sumi_dev',
        email: email,
        role: 'admin' // default mock role as admin to unlock all parts of scaffolding
      }));

      // Simulate network request delay
      setTimeout(() => {
        setIsLoading(false);
        router.push('/feed/home');
      }, 800);
    } catch {
      setIsLoading(false);
      setError('Invalid credentials. Please try again.');
    }
  };

  return (
    <div className={styles.authWrapper}>
      <div className={`${styles.authCard} animate-slide-up`}>
        <div className={styles.logoHeader}>
          <h1>CANOPY</h1>
          <p className="text-subtle">LOG IN TO YOUR INSTANCE</p>
        </div>

        {error && (
          <div className={styles.errorBox}>
            <span>{error}</span>
          </div>
        )}

        <form onSubmit={handleSubmit} className={styles.authForm}>
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

          <button 
            type="submit" 
            className="btn-primary" 
            style={{ width: '100%', marginTop: '8px' }}
            disabled={isLoading}
          >
            {isLoading ? 'VERIFYING...' : 'ENTER'}
          </button>
        </form>

        <div className={styles.footerLink}>
          <span>No account yet? </span>
          <Link href="/register">CREATE ACCOUNT</Link>
        </div>
      </div>
    </div>
  );
}
