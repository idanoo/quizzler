<script>
  import { onMount, onDestroy } from 'svelte';
  import { getPublicDeck, getPublicDeckCards } from '../lib/api.js';

  export let deckId;

  let deck = null;
  let cards = [];
  let loading = true;
  let error = null;

  function handleKeydown(e) {
    if (loading || error || isComplete) return;
    
    if (e.code === 'Space') {
      e.preventDefault();
      flipCard();
    } else if (e.code === 'ArrowRight') {
      e.preventDefault();
      nextCard();
    } else if (e.code === 'ArrowLeft') {
      e.preventDefault();
      prevCard();
    }
  }

  onMount(async () => {
    window.addEventListener('keydown', handleKeydown);
    
    try {
      [deck, cards] = await Promise.all([
        getPublicDeck(deckId),
        getPublicDeckCards(deckId)
      ]);
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  });

  onDestroy(() => {
    window.removeEventListener('keydown', handleKeydown);
  });

  // Fisher-Yates shuffle
  function shuffle(array) {
    const shuffled = [...array];
    for (let i = shuffled.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
    }
    return shuffled;
  }

  let shuffledCards = [];
  let currentIndex = 0;
  let displayIndex = 0;
  let prevIndex = null;
  let prevFlipped = false;
  let isFlipped = false;
  let isComplete = false;
  let isTransitioning = false;
  let slideDirection = null;

  const SLIDE_DURATION = 400;

  $: if (cards.length > 0 && shuffledCards.length === 0) {
    shuffledCards = shuffle(cards);
  }

  $: currentCard = shuffledCards[displayIndex];
  $: prevCard_ = prevIndex !== null ? shuffledCards[prevIndex] : null;

  function flipCard() {
    if (!isTransitioning) {
      isFlipped = !isFlipped;
    }
  }

  function changeCard(delta) {
    prevIndex = currentIndex;
    prevFlipped = isFlipped;
    currentIndex += delta;
    displayIndex = currentIndex;
    slideDirection = delta > 0 ? 'next' : 'prev';
    isFlipped = false;
    isTransitioning = true;
    
    setTimeout(() => {
      prevIndex = null;
      prevFlipped = false;
      slideDirection = null;
      isTransitioning = false;
    }, SLIDE_DURATION);
  }

  function nextCard() {
    if (isTransitioning) return;
    
    if (currentIndex < shuffledCards.length - 1) {
      changeCard(1);
    } else {
      isComplete = true;
    }
  }

  function prevCard() {
    if (isTransitioning) return;
    
    if (currentIndex > 0) {
      changeCard(-1);
    }
  }

  function restart() {
    shuffledCards = shuffle(cards);
    currentIndex = 0;
    displayIndex = 0;
    isFlipped = false;
    isComplete = false;
  }

  function goHome() {
    window.location.href = '/';
  }
</script>

<div class="public-study-page">
  <header class="public-header">
    <button class="header-left" onclick={goHome}>
      <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <rect x="2" y="4" width="20" height="16" rx="2"/>
        <path d="M9 9a3 3 0 1 1 3 3v1"/>
        <circle cx="12" cy="16" r="0.5" fill="currentColor"/>
      </svg>
      <span class="header-title">Quizzler</span>
    </button>
    <div class="header-right">
      <a href="/" class="btn btn-secondary">Sign In</a>
      <a href="/" class="btn btn-primary">Create Account</a>
    </div>
  </header>

  <main class="study-content">
    {#if loading}
      <div class="loading">Loading deck...</div>
    {:else if error}
      <div class="error-state">
        <h2>Deck not found</h2>
        <p>This deck may not exist or isn't public.</p>
        <button class="btn btn-primary" onclick={goHome}>Go Home</button>
      </div>
    {:else if cards.length === 0}
      <div class="empty-state">
        <h2>{deck.name}</h2>
        <p>This deck has no cards yet.</p>
        <button class="btn btn-secondary" onclick={goHome}>← Back</button>
      </div>
    {:else}
      <div class="study-view">
        <div class="content-header">
          <button class="btn btn-ghost" onclick={goHome}>← Back</button>
          <h2>{deck.name}</h2>
          <span class="study-progress">{currentIndex + 1} / {shuffledCards.length}</span>
        </div>

        {#if !isComplete}
          <div class="study-area">
            <div class="flashcard-container">
              {#if prevCard_}
                <div 
                  class="flashcard"
                  class:flipped={prevFlipped}
                  class:slide-out-left={slideDirection === 'next'}
                  class:slide-out-right={slideDirection === 'prev'}
                >
                  <div class="flashcard-inner">
                    <div class="flashcard-front">
                      <p>{prevCard_.front}</p>
                    </div>
                    <div class="flashcard-back">
                      <p>{prevCard_.back}</p>
                    </div>
                  </div>
                </div>
              {/if}
              
              <button 
                class="flashcard" 
                class:flipped={isFlipped}
                class:slide-in-right={slideDirection === 'next'}
                class:slide-in-left={slideDirection === 'prev'}
                onclick={flipCard}
              >
                <div class="flashcard-inner">
                  <div class="flashcard-front">
                    <p>{currentCard.front}</p>
                  </div>
                  <div class="flashcard-back">
                    <p>{currentCard.back}</p>
                  </div>
                </div>
              </button>
            </div>
            <p class="flip-hint">Click or press Space to flip · Arrow keys to navigate</p>
            <div class="study-controls">
              <button class="btn btn-secondary" onclick={prevCard} disabled={currentIndex === 0}>
                Previous
              </button>
              <button class="btn btn-primary" onclick={nextCard}>
                {currentIndex === shuffledCards.length - 1 ? 'Finish' : 'Next'}
              </button>
            </div>
          </div>
        {:else}
          <div class="study-complete">
            <div class="complete-icon">🎉</div>
            <h3>Great job!</h3>
            <p>You've reviewed all cards in this deck</p>
            <div class="complete-actions">
              <button class="btn btn-primary" onclick={restart}>Study Again</button>
              <button class="btn btn-secondary" onclick={goHome}>Back to Home</button>
            </div>
            <div class="signup-cta">
              <p>Want to create your own decks?</p>
              <a href="/" class="btn btn-accent">Create a Free Account</a>
            </div>
          </div>
        {/if}
      </div>
    {/if}
  </main>
</div>

<style>
  .public-study-page {
    min-height: 100vh;
  }

  .public-header {
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
    gap: 12px;
  }

  .header-right a {
    text-decoration: none;
  }

  .study-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 40px 32px;
  }

  .loading, .error-state, .empty-state {
    text-align: center;
    padding: 60px 20px;
  }

  .error-state h2, .empty-state h2 {
    font-family: var(--font-display);
    font-size: 1.5rem;
    margin-bottom: 12px;
  }

  .error-state p, .empty-state p {
    color: var(--text-secondary);
    margin-bottom: 24px;
  }

  .content-header {
    display: flex;
    align-items: center;
    gap: 24px;
    margin-bottom: 32px;
  }

  .content-header h2 {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 600;
    flex: 1;
  }

  .study-progress {
    background: var(--bg-elevated);
    padding: 6px 14px;
    border-radius: var(--radius-sm);
    font-size: 0.9rem;
    color: var(--text-secondary);
  }

  .study-area {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 40px 20px;
  }

  .flashcard-container {
    width: 100%;
    max-width: 600px;
    height: 380px;
    position: relative;
    margin-bottom: 16px;
    overflow: hidden;
  }

  .flashcard {
    width: 100%;
    height: 100%;
    perspective: 1000px;
    cursor: pointer;
    background: none;
    border: none;
    padding: 0;
    position: absolute;
    top: 0;
    left: 0;
  }

  .flashcard.slide-in-right {
    animation: slideInFromRight 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards;
  }

  .flashcard.slide-in-left {
    animation: slideInFromLeft 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards;
  }

  .flashcard.slide-out-left {
    animation: slideOutToLeft 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards;
    pointer-events: none;
  }

  .flashcard.slide-out-right {
    animation: slideOutToRight 0.4s cubic-bezier(0.25, 0.46, 0.45, 0.94) forwards;
    pointer-events: none;
  }

  @keyframes slideInFromRight {
    from { transform: translateX(100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }

  @keyframes slideInFromLeft {
    from { transform: translateX(-100%); opacity: 0; }
    to { transform: translateX(0); opacity: 1; }
  }

  @keyframes slideOutToLeft {
    from { transform: translateX(0); opacity: 1; }
    to { transform: translateX(-100%); opacity: 0; }
  }

  @keyframes slideOutToRight {
    from { transform: translateX(0); opacity: 1; }
    to { transform: translateX(100%); opacity: 0; }
  }

  .flashcard-inner {
    position: relative;
    width: 100%;
    height: 100%;
    transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
    transform-style: preserve-3d;
  }

  .flashcard.flipped .flashcard-inner {
    transform: rotateY(180deg);
  }

  .flashcard.slide-in-right .flashcard-inner,
  .flashcard.slide-in-left .flashcard-inner,
  .flashcard.slide-out-left .flashcard-inner,
  .flashcard.slide-out-right .flashcard-inner {
    transition: none;
  }

  .flashcard-front,
  .flashcard-back {
    position: absolute;
    width: 100%;
    height: 100%;
    backface-visibility: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px;
    border-radius: var(--radius-xl);
    text-align: center;
  }

  .flashcard-front {
    background: linear-gradient(135deg, var(--bg-card), var(--bg-elevated));
    border: 1px solid var(--border-normal);
    box-shadow: var(--shadow-lg), var(--shadow-glow);
  }

  .flashcard-front p {
    font-family: var(--font-display);
    font-size: 1.8rem;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.4;
  }

  .flashcard-back {
    background: linear-gradient(135deg, var(--accent-primary), #38b2ac);
    transform: rotateY(180deg);
  }

  .flashcard-back p {
    font-family: var(--font-display);
    font-size: 1.6rem;
    font-weight: 600;
    color: var(--bg-deep);
    line-height: 1.4;
  }

  .flip-hint {
    color: var(--text-muted);
    font-size: 0.85rem;
    margin-bottom: 32px;
  }

  .study-controls {
    display: flex;
    gap: 16px;
  }

  .study-controls button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .study-complete {
    text-align: center;
    padding: 80px 20px;
    animation: fadeIn 0.4s ease;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .complete-icon {
    font-size: 4rem;
    margin-bottom: 20px;
  }

  .study-complete h3 {
    font-family: var(--font-display);
    font-size: 2rem;
    margin-bottom: 12px;
  }

  .study-complete > p {
    color: var(--text-secondary);
    margin-bottom: 24px;
  }

  .complete-actions {
    display: flex;
    gap: 16px;
    justify-content: center;
    margin-bottom: 48px;
  }

  .signup-cta {
    padding-top: 32px;
    border-top: 1px solid var(--border-subtle);
  }

  .signup-cta p {
    color: var(--text-secondary);
    margin-bottom: 16px;
  }

  .signup-cta a {
    text-decoration: none;
  }

  @media (max-width: 768px) {
    .public-header {
      padding: 12px 16px;
    }

    .header-right .btn-secondary {
      display: none;
    }

    .study-content {
      padding: 24px 16px;
    }

    .flashcard-container {
      height: 300px;
    }

    .flashcard-front p,
    .flashcard-back p {
      font-size: 1.3rem;
    }
  }
</style>
