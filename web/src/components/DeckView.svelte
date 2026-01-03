<script>
  import { createEventDispatcher } from 'svelte';
  import {
    getCards,
    updateDeck,
    deleteDeck,
    createCard,
    updateCard,
    deleteCard,
    importCards,
  } from '../lib/api.js';
  import Modal from './Modal.svelte';

  export let deck;
  export let cards;

  const dispatch = createEventDispatcher();

  let showDeckModal = false;
  let showCardModal = false;
  let showImportModal = false;
  let editingCard = null;

  let deckName = '';
  let deckDescription = '';
  let cardFront = '';
  let cardBack = '';
  let importText = '';
  let modalError = '';

  function goBack() {
    dispatch('back');
  }

  function startStudy() {
    if (cards.length > 0) {
      dispatch('study');
    }
  }

  // Deck editing
  function openDeckModal() {
    deckName = deck.name;
    deckDescription = deck.description || '';
    modalError = '';
    showDeckModal = true;
  }

  async function handleUpdateDeck() {
    modalError = '';
    if (!deckName.trim()) {
      modalError = 'Name is required';
      return;
    }
    try {
      const updated = await updateDeck(deck.id, deckName, deckDescription);
      deck = { ...deck, ...updated };
      showDeckModal = false;
    } catch (err) {
      modalError = err.message;
    }
  }

  async function handleDeleteDeck() {
    if (confirm('Are you sure you want to delete this deck? All cards will be deleted.')) {
      try {
        await deleteDeck(deck.id);
        dispatch('back');
      } catch (err) {
        console.error('Failed to delete deck:', err);
      }
    }
  }

  // Card editing
  function openCardModal(card = null) {
    editingCard = card;
    cardFront = card?.front || '';
    cardBack = card?.back || '';
    modalError = '';
    showCardModal = true;
  }

  async function handleSaveCard() {
    modalError = '';
    if (!cardFront.trim() || !cardBack.trim()) {
      modalError = 'Front and back are required';
      return;
    }
    try {
      if (editingCard) {
        await updateCard(editingCard.id, cardFront, cardBack);
      } else {
        await createCard(deck.id, cardFront, cardBack);
      }
      cards = await getCards(deck.id);
      dispatch('cardsUpdate', cards);
      showCardModal = false;
    } catch (err) {
      modalError = err.message;
    }
  }

  async function handleDeleteCard(cardId) {
    if (confirm('Delete this card?')) {
      try {
        await deleteCard(cardId);
        cards = await getCards(deck.id);
        dispatch('cardsUpdate', cards);
      } catch (err) {
        console.error('Failed to delete card:', err);
      }
    }
  }

  // Import cards
  function openImportModal() {
    importText = '';
    modalError = '';
    showImportModal = true;
  }

  async function handleImport() {
    modalError = '';
    if (!importText.trim()) {
      modalError = 'Please paste some content to import';
      return;
    }

    const lines = importText.trim().split('\n');
    const cardsToImport = [];

    for (const line of lines) {
      const parts = line.split('\t');
      if (parts.length >= 2) {
        cardsToImport.push({ front: parts[0].trim(), back: parts[1].trim() });
      }
    }

    if (cardsToImport.length === 0) {
      modalError = 'No valid cards found. Use tab-separated format: Question<tab>Answer';
      return;
    }

    try {
      const result = await importCards(deck.id, cardsToImport);
      cards = await getCards(deck.id);
      dispatch('cardsUpdate', cards);
      showImportModal = false;
      alert(`Successfully imported ${result.imported} cards`);
    } catch (err) {
      modalError = err.message;
    }
  }
</script>

<div class="deck-view">
  <div class="content-header">
    <button class="btn btn-ghost" onclick={goBack}>← Back</button>
    <div class="deck-title-area">
      <h2>{deck.name}</h2>
      {#if deck.description}
        <p class="deck-description">{deck.description}</p>
      {/if}
    </div>
    <div class="deck-actions">
      <button
        class="btn btn-accent"
        onclick={startStudy}
        disabled={cards.length === 0}
      >
        Study
      </button>
      <button class="btn btn-ghost" onclick={openDeckModal}>Edit</button>
      <button class="btn btn-danger-ghost" onclick={handleDeleteDeck}>Delete</button>
    </div>
  </div>

  <div class="cards-section">
    <div class="cards-header">
      <h3>Cards</h3>
      <div class="cards-header-actions">
        <button class="btn btn-ghost" onclick={openImportModal}>Import</button>
        <button class="btn btn-secondary" onclick={() => openCardModal()}>+ Add Card</button>
      </div>
    </div>

    {#if cards.length === 0}
      <div class="empty-state">
        <div class="empty-icon">🃏</div>
        <h3>No cards yet</h3>
        <p>Add cards to start building your deck</p>
      </div>
    {:else}
      <div class="cards-list">
        {#each cards as card}
          <div class="card-item">
            <div class="card-side">
              <div class="card-side-label">Front</div>
              <div class="card-side-content">{card.front}</div>
            </div>
            <div class="card-side">
              <div class="card-side-label">Back</div>
              <div class="card-side-content">{card.back}</div>
            </div>
            <div class="card-actions">
              <button class="btn btn-ghost" onclick={() => openCardModal(card)}>Edit</button>
              <button class="btn btn-danger-ghost" onclick={() => handleDeleteCard(card.id)}>
                Delete
              </button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

{#if showDeckModal}
  <Modal title="Edit Deck" on:close={() => (showDeckModal = false)}>
    <form onsubmit={(e) => { e.preventDefault(); handleUpdateDeck(); }}>
      <div class="form-group">
        <label for="deck-name">Name</label>
        <input type="text" id="deck-name" bind:value={deckName} required />
      </div>
      <div class="form-group">
        <label for="deck-desc">Description</label>
        <textarea id="deck-desc" bind:value={deckDescription}></textarea>
      </div>
      {#if modalError}
        <p class="error-message">{modalError}</p>
      {/if}
      <button type="submit" class="btn btn-primary modal-submit">Save Changes</button>
    </form>
  </Modal>
{/if}

{#if showCardModal}
  <Modal title={editingCard ? 'Edit Card' : 'Add Card'} on:close={() => (showCardModal = false)}>
    <form onsubmit={(e) => { e.preventDefault(); handleSaveCard(); }}>
      <div class="form-group">
        <label for="card-front">Front</label>
        <textarea
          id="card-front"
          bind:value={cardFront}
          required
          placeholder="Question or term"
        ></textarea>
      </div>
      <div class="form-group">
        <label for="card-back">Back</label>
        <textarea
          id="card-back"
          bind:value={cardBack}
          required
          placeholder="Answer or definition"
        ></textarea>
      </div>
      {#if modalError}
        <p class="error-message">{modalError}</p>
      {/if}
      <button type="submit" class="btn btn-primary modal-submit">
        {editingCard ? 'Save Changes' : 'Add Card'}
      </button>
    </form>
  </Modal>
{/if}

{#if showImportModal}
  <Modal title="Import Cards" on:close={() => (showImportModal = false)}>
    <form onsubmit={(e) => { e.preventDefault(); handleImport(); }}>
      <div class="form-group">
        <label for="import-text">Paste tab-separated content</label>
        <textarea
          id="import-text"
          bind:value={importText}
          required
          placeholder="Question&#9;Answer&#10;What is 2+2?&#9;4&#10;Capital of France?&#9;Paris"
          class="import-textarea"
        ></textarea>
        <p class="form-hint">Format: Question&lt;tab&gt;Answer (one per line)</p>
      </div>
      {#if modalError}
        <p class="error-message">{modalError}</p>
      {/if}
      <button type="submit" class="btn btn-primary modal-submit">Import Cards</button>
    </form>
  </Modal>
{/if}

<style>
  .content-header {
    display: flex;
    align-items: center;
    gap: 24px;
    margin-bottom: 32px;
    flex-wrap: wrap;
  }

  .content-header h2 {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 600;
  }

  .deck-title-area {
    flex: 1;
  }

  .deck-description {
    color: var(--text-secondary);
    font-size: 0.95rem;
    margin-top: 4px;
  }

  .deck-actions {
    display: flex;
    gap: 12px;
  }

  .deck-actions button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .cards-section {
    margin-top: 20px;
  }

  .cards-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;
  }

  .cards-header h3 {
    font-family: var(--font-display);
    font-size: 1.3rem;
    font-weight: 600;
  }

  .cards-header-actions {
    display: flex;
    gap: 12px;
  }

  .import-textarea {
    min-height: 200px;
    font-family: monospace;
    font-size: 0.9rem;
  }

  .form-hint {
    font-size: 0.8rem;
    color: var(--text-muted);
    margin-top: 8px;
  }

  .cards-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .card-item {
    display: grid;
    grid-template-columns: 1fr 1fr auto;
    gap: 20px;
    align-items: center;
    padding: 20px;
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-md);
    transition: all var(--transition-normal);
  }

  .card-item:hover {
    border-color: var(--border-normal);
    background: var(--bg-elevated);
  }

  .card-side {
    padding: 12px 16px;
    background: var(--bg-surface);
    border-radius: var(--radius-sm);
    min-height: 60px;
  }

  .card-side-label {
    font-size: 0.7rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--text-muted);
    margin-bottom: 4px;
  }

  .card-side-content {
    color: var(--text-primary);
    font-size: 0.95rem;
  }

  .card-actions {
    display: flex;
    gap: 8px;
  }

  .card-actions .btn {
    padding: 8px 12px;
    font-size: 0.85rem;
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

  .modal-submit {
    width: 100%;
    margin-top: 8px;
  }

  @media (max-width: 768px) {
    .content-header {
      flex-direction: column;
      align-items: flex-start;
    }

    .deck-actions {
      width: 100%;
      justify-content: flex-start;
    }

    .card-item {
      grid-template-columns: 1fr;
      gap: 12px;
    }

    .card-actions {
      justify-content: flex-end;
    }
  }
</style>

