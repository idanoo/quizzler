<script>
  import { onMount } from 'svelte';
  import { user, logout } from '../lib/stores.js';
  import { getDeck, getCards } from '../lib/api.js';
  import Decks from './Decks.svelte';
  import DeckView from './DeckView.svelte';
  import Study from './Study.svelte';

  let view = 'decks';
  let currentDeck = null;
  let currentCards = [];
  let isPublicDeck = false;
  let loading = false;

  onMount(() => {
    handleRoute();
    window.addEventListener('popstate', handleRoute);
    return () => window.removeEventListener('popstate', handleRoute);
  });

  async function handleRoute() {
    const path = window.location.pathname;
    const deckMatch = path.match(/^\/deck\/(\d+)(\/study)?$/);

    if (deckMatch) {
      const deckId = parseInt(deckMatch[1]);
      const isStudy = !!deckMatch[2];

      // Only load if different deck or no deck loaded
      if (!currentDeck || currentDeck.id !== deckId) {
        loading = true;
        try {
          currentDeck = await getDeck(deckId);
          currentCards = await getCards(deckId);
        } catch (err) {
          console.error('Failed to load deck:', err);
          navigate('/');
          return;
        } finally {
          loading = false;
        }
      }
      view = isStudy ? 'study' : 'deck';
    } else {
      view = 'decks';
      currentDeck = null;
      currentCards = [];
    }
  }

  function navigate(path, replace = false) {
    if (replace) {
      history.replaceState({}, '', path);
    } else {
      history.pushState({}, '', path);
    }
    handleRoute();
  }

  function openDeck(event) {
    currentDeck = event.detail.deck;
    currentCards = event.detail.cards;
    isPublicDeck = event.detail.isPublic || false;
    navigate(`/deck/${currentDeck.id}`);
  }

  function startStudy() {
    navigate(`/deck/${currentDeck.id}/study`);
  }

  function backToDecks() {
    isPublicDeck = false;
    navigate('/');
  }

  function backToDeck() {
    navigate(`/deck/${currentDeck.id}`);
  }

  function handleCardsUpdate(event) {
    currentCards = event.detail;
  }
</script>

<div class="main-view">
  <header class="main-header">
    <button class="header-left" onclick={backToDecks}>
      <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="2" y="4" width="20" height="16" rx="2"/>
        <path d="M9 9a3 3 0 1 1 3 3v1"/>
        <circle cx="12" cy="16" r="0.5" fill="currentColor"/>
      </svg>
      <span class="header-title">Quizzler</span>
    </button>
    <div class="header-right">
      <span class="user-email">{$user?.email}</span>
      <button class="btn btn-ghost" onclick={logout}>Sign Out</button>
    </div>
  </header>

  <main class="main-content">
    {#if loading}
      <div class="loading">Loading...</div>
    {:else if view === 'decks'}
      <Decks on:openDeck={openDeck} />
    {:else if view === 'deck'}
      <DeckView
        deck={currentDeck}
        cards={currentCards}
        readOnly={isPublicDeck}
        on:back={backToDecks}
        on:study={startStudy}
        on:cardsUpdate={handleCardsUpdate}
      />
    {:else if view === 'study'}
      <Study
        deck={currentDeck}
        cards={currentCards}
        on:exit={backToDeck}
      />
    {/if}
  </main>
</div>

<style>
  .main-view {
    min-height: 100vh;
  }

  .main-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 32px;
    background: rgba(20, 24, 33, 0.8);
    backdrop-filter: blur(20px);
    border-bottom: 1px solid var(--border-subtle);
    position: sticky;
    top: 0;
    z-index: 100;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 10px;
    background: none;
    border: none;
    padding: 8px 12px;
    margin: -8px -12px;
    border-radius: var(--radius-md);
    cursor: pointer;
    transition: background var(--transition-fast);
  }

  .header-left:hover {
    background: var(--bg-elevated);
  }

  .logo-icon {
    width: 28px;
    height: 28px;
    color: var(--accent-primary);
    filter: drop-shadow(0 0 8px var(--accent-glow));
  }

  .header-title {
    font-family: var(--font-display);
    font-size: 1.3rem;
    font-weight: 600;
    color: var(--text-primary);
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .user-email {
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  .main-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 40px 32px;
  }

  .loading {
    text-align: center;
    padding: 60px;
    color: var(--text-secondary);
  }

  @media (max-width: 768px) {
    .main-header {
      padding: 12px 16px;
    }

    .user-email {
      display: none;
    }

    .main-content {
      padding: 24px 16px;
    }
  }
</style>

