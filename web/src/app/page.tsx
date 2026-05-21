'use client';

import React, { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import styles from './landing.module.css';

export default function LandingPage() {
  const router = useRouter();
  const [checkingSession, setCheckingSession] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('canopy_token');
    if (token) {
      router.push('/feed/home');
    } else {
      setCheckingSession(false);
    }
  }, [router]);

  if (checkingSession) {
    return (
      <div className={styles.landingWrapper} style={{ justifyContent: 'center' }}>
        <span className={styles.loadingText}>SYNCHRONIZING...</span>
      </div>
    );
  }

  return (
    <div className={styles.landingWrapper}>
      <header className={styles.landingHeader}>
        <div className={styles.logo}>CANOPY</div>
        <span className={styles.tagline}>DECENRALIZED / MINIMAL / MONOCHROME</span>
      </header>

      <main className={`${styles.landingCard} animate-slide-up`}>
        <h1 className={styles.cardTitle}>Federated Social System</h1>
        <p className={styles.cardDescription}>
          Canopy is an elite, high-performance federated social microblogging and long-form platform built on a pure Go backend and a highly polished monochrome Next.js frontend.
        </p>

        <div className={styles.actions}>
          <Link href="/login" className="btn-primary" style={{ width: '100%', textAlign: 'center' }}>
            ENTER SERVER
          </Link>
          <Link href="/register" className="btn-secondary" style={{ width: '100%', textAlign: 'center' }}>
            CREATE ACCOUNT
          </Link>
        </div>
      </main>

      <footer className={styles.landingFooter}>
        <span>INSTANCE: canopy.social</span>
        <span>•</span>
        <span>VERSION 1.12.0</span>
      </footer>
    </div>
  );
}
