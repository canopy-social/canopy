'use client';

import React, { useState } from 'react';
import styles from './messages.module.css';

interface Message {
  id: string;
  sender: 'me' | 'them';
  content: string;
  time: string;
}

interface Conversation {
  id: string;
  name: string;
  handle: string;
  avatar: string;
  messages: Message[];
}

export default function DirectMessagesPage() {
  const [activeConvId, setActiveConvId] = useState('c1');
  const [inputText, setInputText] = useState('');

  const [conversations, setConversations] = useState<Conversation[]>([
    {
      id: 'c1',
      name: 'Alice Johnson',
      handle: '@alice@canopy.social',
      avatar: 'A',
      messages: [
        { id: 'm1', sender: 'them', content: 'Hey Sumi! Did you check out the new timeline performance checks?', time: '10:05 AM' },
        { id: 'm2', sender: 'me', content: 'Yes! The Redis fan-out is performing incredibly well.', time: '10:08 AM' },
        { id: 'm3', sender: 'them', content: 'Fantastic. I will start testing multi-threading concurrency next.', time: '10:10 AM' }
      ]
    },
    {
      id: 'c2',
      name: 'Dave H',
      handle: '@dave@fediverse.org',
      avatar: 'D',
      messages: [
        { id: 'm4', sender: 'them', content: 'Do you support instance-wide blocks for external domains?', time: 'Yesterday' },
        { id: 'm5', sender: 'me', content: 'Yes, the moderation service exposes domain blocks via admin panel APIs.', time: 'Yesterday' }
      ]
    },
    {
      id: 'c3',
      name: 'Linus B',
      handle: '@linus@kernel.org',
      avatar: 'L',
      messages: [
        { id: 'm6', sender: 'them', content: 'Make sure no comments are added to internal/accounts repository.', time: '2 days ago' },
        { id: 'm7', sender: 'me', content: 'Absolutely, I verified it with a strict git check script.', time: '2 days ago' }
      ]
    }
  ]);

  const activeConv = conversations.find(c => c.id === activeConvId) || conversations[0];

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputText.trim()) return;

    const newMessage: Message = {
      id: `msg_${Math.random().toString(36).substr(2, 9)}`,
      sender: 'me',
      content: inputText,
      time: 'Just now'
    };

    setConversations(conversations.map(c => {
      if (c.id === activeConvId) {
        return {
          ...c,
          messages: [...c.messages, newMessage]
        };
      }
      return c;
    }));

    setInputText('');
  };

  return (
    <div className={`${styles.dmContainer} animate-fade-in`}>
      {/* Conversations List Left */}
      <div className={styles.conversationsList}>
        <div className={styles.conversationsHeader}>Direct Messages</div>
        {conversations.map((conv) => {
          const lastMsg = conv.messages[conv.messages.length - 1];
          return (
            <div 
              key={conv.id} 
              onClick={() => setActiveConvId(conv.id)}
              className={`${styles.conversationItem} ${conv.id === activeConvId ? styles.conversationActive : ''}`}
            >
              <div className={styles.avatarSmall}>{conv.avatar}</div>
              <div className={styles.convDetails}>
                <span className={styles.convName}>{conv.name}</span>
                <span className={styles.convSnippet}>{lastMsg ? lastMsg.content : 'No messages'}</span>
              </div>
            </div>
          );
        })}
      </div>

      {/* Chat Area Right */}
      <div className={styles.chatArea}>
        <div className={styles.chatHeader}>
          <span className={styles.chatTitle}>{activeConv.name}</span>
          <span className={styles.chatSubtitle}>{activeConv.handle}</span>
        </div>

        <div className={styles.messagesStream}>
          {activeConv.messages.map((msg) => (
            <div 
              key={msg.id} 
              className={`${styles.messageRow} ${msg.sender === 'me' ? styles.rowSent : styles.rowReceived}`}
            >
              <div className={`${styles.bubble} ${msg.sender === 'me' ? styles.bubbleSent : styles.bubbleReceived}`}>
                {msg.content}
                <span className={styles.messageTime}>{msg.time}</span>
              </div>
            </div>
          ))}
        </div>

        <form onSubmit={handleSendMessage} className={styles.composerArea}>
          <input
            type="text"
            className={styles.composerInput}
            placeholder="Type your message..."
            value={inputText}
            onChange={(e) => setInputText(e.target.value)}
            required
          />
          <button type="submit" className="btn-primary">SEND</button>
        </form>
      </div>
    </div>
  );
}
