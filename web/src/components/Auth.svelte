<script>
  import { onMount } from 'svelte';
  import { login, register, getPublicDecksBrowse } from '../lib/api.js';

  let activeTab = 'login';
  let email = '';
  let password = '';
  let error = '';
  let loading = false;
  let publicDecks = [];
  let loadingDecks = true;

  onMount(async () => {
    try {
      publicDecks = await getPublicDecksBrowse();
    } catch (err) {
      console.error('Failed to load public decks:', err);
    } finally {
      loadingDecks = false;
    }
  });

  async function handleSubmit() {
    error = '';
    loading = true;
    try {
      if (activeTab === 'login') {
        await login(email, password);
      } else {
        await register(email, password);
      }
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function switchTab(tab) {
    activeTab = tab;
    error = '';
  }
</script>

<div class="auth-page">
  <div class="auth-section">
    <div class="auth-container">
      <div class="auth-header">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="2" y="4" width="20" height="16" rx="2"/>
            <path d="M9 9a3 3 0 1 1 3 3v1"/>
            <circle cx="12" cy="16" r="0.5" fill="currentColor"/>
          </svg>
          <h1>Quizzler</h1>
        </div>
        <p class="tagline">Master anything, one card at a time</p>
      </div>

      <div class="auth-tabs">
        <button
          class="auth-tab"
          class:active={activeTab === 'login'}
          onclick={() => switchTab('login')}
        >
          Sign In
        </button>
        <button
          class="auth-tab"
          class:active={activeTab === 'register'}
          onclick={() => switchTab('register')}
        >
          Create Account
        </button>
      </div>

      <form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
        <div class="form-group">
          <label for="email">Email</label>
          <input
            type="email"
            id="email"
            bind:value={email}
            required
            placeholder="you@example.com"
          />
        </div>
        <div class="form-group">
          <label for="password">Password</label>
          <input
            type="password"
            id="password"
            bind:value={password}
            required
            placeholder="••••••••"
            minlength="6"
          />
        </div>
        <button type="submit" class="btn btn-primary submit-btn" disabled={loading}>
          {#if loading}
            Loading...
          {:else if activeTab === 'login'}
            Sign In
          {:else}
            Create Account
          {/if}
        </button>
        {#if error}
          <p class="error-message">{error}</p>
        {/if}
      </form>
    </div>
  </div>

  {#if publicDecks.length > 0 || loadingDecks}
    <div class="public-decks-section">
      <div class="public-decks-container">
        <h2 class="section-title">Browse Public Decks</h2>
        <p class="section-subtitle">Click any deck to start studying — no account required!</p>

        {#if loadingDecks}
          <div class="loading">Loading decks...</div>
        {:else}
          <div class="decks-grid">
            {#each publicDecks as deck}
              <a href="/study/{deck.id}" class="deck-card">
                <h3>{deck.name}</h3>
                <p>{deck.description || 'No description'}</p>
                <div class="deck-card-footer">
                  <span class="card-count">🃏 {deck.card_count} cards</span>
                  <span class="study-badge">Study →</span>
                </div>
              </a>
            {/each}
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .auth-page {
    min-height: 100vh;
  }

  .auth-section {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
  }

  .auth-container {
    width: 100%;
    max-width: 420px;
    animation: fadeIn 0.4s ease;
  }

  .auth-header {
    text-align: center;
    margin-bottom: 48px;
  }

  .logo {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    margin-bottom: 16px;
  }

  .logo-icon {
    width: 48px;
    height: 48px;
    color: var(--accent-primary);
    filter: drop-shadow(0 0 12px var(--accent-glow));
    animation: pulse 3s ease-in-out infinite;
  }

  .logo h1 {
    font-family: var(--font-display);
    font-size: 2.5rem;
    font-weight: 600;
    letter-spacing: -0.02em;
    background: linear-gradient(135deg, var(--text-primary), var(--accent-primary));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .tagline {
    color: var(--text-secondary);
    font-size: 1.1rem;
  }

  .auth-tabs {
    display: flex;
    background: var(--bg-surface);
    border-radius: var(--radius-md);
    padding: 4px;
    margin-bottom: 32px;
    border: 1px solid var(--border-subtle);
  }

  .auth-tab {
    flex: 1;
    padding: 12px 16px;
    border: none;
    background: transparent;
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.95rem;
    font-weight: 500;
    cursor: pointer;
    border-radius: var(--radius-sm);
    transition: all var(--transition-normal);
  }

  .auth-tab:hover {
    color: var(--text-primary);
  }

  .auth-tab.active {
    background: var(--bg-elevated);
    color: var(--text-primary);
    box-shadow: var(--shadow-sm);
  }

  .submit-btn {
    width: 100%;
    margin-top: 8px;
    padding: 16px;
  }

  .submit-btn:disabled {
    opacity: 0.7;
    cursor: not-allowed;
  }

  /* Public Decks Section */
  .public-decks-section {
    background: var(--bg-surface);
    border-top: 1px solid var(--border-subtle);
    padding: 60px 20px;
  }

  .public-decks-container {
    max-width: 1000px;
    margin: 0 auto;
  }

  .section-title {
    font-family: var(--font-display);
    font-size: 1.8rem;
    font-weight: 600;
    text-align: center;
    margin-bottom: 8px;
  }

  .section-subtitle {
    color: var(--text-secondary);
    text-align: center;
    margin-bottom: 40px;
  }

  .loading {
    text-align: center;
    padding: 40px;
    color: var(--text-secondary);
  }

  .decks-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 20px;
  }

  .deck-card {
    display: block;
    text-decoration: none;
    background: var(--bg-card);
    border: 1px solid var(--border-subtle);
    border-radius: var(--radius-lg);
    padding: 24px;
    transition: all var(--transition-normal);
    position: relative;
    overflow: hidden;
    cursor: pointer;
  }

  .deck-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: linear-gradient(90deg, var(--accent-secondary), var(--accent-primary));
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

  .study-badge {
    background: var(--accent-primary);
    color: var(--bg-deep);
    padding: 4px 10px;
    border-radius: var(--radius-sm);
    font-size: 0.8rem;
    font-weight: 600;
    opacity: 0;
    transform: translateX(-8px);
    transition: all var(--transition-normal);
  }

  .deck-card:hover .study-badge {
    opacity: 1;
    transform: translateX(0);
  }

  @media (max-width: 768px) {
    .auth-section {
      padding: 40px 16px;
    }

    .public-decks-section {
      padding: 40px 16px;
    }
  }
</style>

