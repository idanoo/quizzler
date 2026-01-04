<script>
  import { login, register } from '../lib/api.js';

  let activeTab = 'login';
  let email = '';
  let password = '';
  let error = '';
  let loading = false;

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

<div class="auth-view">
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

<style>
  .auth-view {
    min-height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
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
</style>

