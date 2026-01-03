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
    navigate(`/deck/${currentDeck.id}`);
  }

  function startStudy() {
    navigate(`/deck/${currentDeck.id}/study`);
  }

  function backToDecks() {
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
    <div class="header-left">
      <span class="logo-icon">◈</span>
      <span class="header-title">Quizzler</span>
    </div>
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
  }

  .logo-icon {
    font-size: 1.5rem;
    color: var(--accent-primary);
    text-shadow: 0 0 30px var(--accent-glow);
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

