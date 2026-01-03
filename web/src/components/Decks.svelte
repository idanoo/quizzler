<script>
  import { onMount, createEventDispatcher } from 'svelte';
  import { getDecks, getDeck, getCards, createDeck } from '../lib/api.js';
  import Modal from './Modal.svelte';

  const dispatch = createEventDispatcher();

  let decks = [];
  let loading = true;
  let showModal = false;
  let newDeckName = '';
  let newDeckDescription = '';
  let modalError = '';

  onMount(async () => {
    await loadDecks();
  });

  async function loadDecks() {
    loading = true;
    try {
      decks = await getDecks();
    } catch (err) {
      console.error('Failed to load decks:', err);
    } finally {
      loading = false;
    }
  }

  async function openDeck(deckId) {
    try {
      const deck = await getDeck(deckId);
      const cards = await getCards(deckId);
      dispatch('openDeck', { deck, cards });
    } catch (err) {
      console.error('Failed to open deck:', err);
    }
  }

  async function handleCreateDeck() {
    modalError = '';
    if (!newDeckName.trim()) {
      modalError = 'Name is required';
      return;
    }
    try {
      const deck = await createDeck(newDeckName, newDeckDescription);
      const cards = await getCards(deck.id);
      closeModal();
      dispatch('openDeck', { deck, cards });
    } catch (err) {
      modalError = err.message;
    }
  }

  function openModal() {
    newDeckName = '';
    newDeckDescription = '';
    modalError = '';
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }
</script>

<div class="decks-view">
  <div class="content-header">
    <h2>Your Decks</h2>
    <button class="btn btn-primary" onclick={openModal}>
      <span>+</span> New Deck
    </button>
  </div>

  {#if loading}
    <div class="loading">Loading...</div>
  {:else if decks.length === 0}
    <div class="empty-state">
      <div class="empty-icon">📚</div>
      <h3>No decks yet</h3>
      <p>Create your first deck to start learning</p>
    </div>
  {:else}
    <div class="decks-grid">
      {#each decks as deck}
        <button class="deck-card" onclick={() => openDeck(deck.id)}>
          <h3>{deck.name}</h3>
          <p>{deck.description || 'No description'}</p>
          <div class="deck-card-footer">
            <span class="card-count">🃏 {deck.card_count} cards</span>
          </div>
        </button>
      {/each}
    </div>
  {/if}
</div>

{#if showModal}
  <Modal title="New Deck" on:close={closeModal}>
    <form onsubmit={(e) => { e.preventDefault(); handleCreateDeck(); }}>
      <div class="form-group">
        <label for="deck-name">Name</label>
        <input
          type="text"
          id="deck-name"
          bind:value={newDeckName}
          required
          placeholder="e.g., Spanish Vocabulary"
        />
      </div>
      <div class="form-group">
        <label for="deck-desc">Description (optional)</label>
        <textarea
          id="deck-desc"
          bind:value={newDeckDescription}
          placeholder="What is this deck about?"
        ></textarea>
      </div>
      {#if modalError}
        <p class="error-message">{modalError}</p>
      {/if}
      <button type="submit" class="btn btn-primary modal-submit">Create Deck</button>
    </form>
  </Modal>
{/if}

<style>
  .content-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 32px;
  }

  .content-header h2 {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 600;
  }

  .loading {
    text-align: center;
    padding: 60px;
    color: var(--text-secondary);
  }

  .decks-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .deck-card {
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-lg);
    padding: 24px;
    cursor: pointer;
    transition: all var(--transition-normal);
    position: relative;
    overflow: hidden;
    text-align: left;
    width: 100%;
  }

  .deck-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--accent-primary), var(--accent-secondary));
    opacity: 0;
    transition: opacity var(--transition-normal);
  }

  .deck-card:hover {
    transform: translateY(-4px);
    border-color: var(--border-normal);
    box-shadow: var(--shadow-md);
  }

  .deck-card:hover::before {
    opacity: 1;
  }

  .deck-card h3 {
    font-family: var(--font-display);
    font-size: 1.3rem;
    font-weight: 600;
    margin-bottom: 8px;
    color: var(--text-primary);
  }

  .deck-card p {
    color: var(--text-secondary);
    font-size: 0.9rem;
    margin-bottom: 16px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .deck-card-footer {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--text-muted);
    font-size: 0.85rem;
  }

  .card-count {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    background: var(--bg-elevated);
    padding: 4px 10px;
    border-radius: var(--radius-sm);
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
    color: var(--text-secondary);
  }

  .empty-icon {
    font-size: 3rem;
    margin-bottom: 16px;
    opacity: 0.6;
  }

  .empty-state h3 {
    font-family: var(--font-display);
    font-size: 1.5rem;
    margin-bottom: 8px;
    color: var(--text-primary);
  }

  .empty-state p {
    color: var(--text-muted);
  }

  .modal-submit {
    width: 100%;
    margin-top: 8px;
  }
</style>

