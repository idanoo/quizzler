<script>
  import { createEventDispatcher } from 'svelte';

  export let deck;
  export let cards;

  const dispatch = createEventDispatcher();

  // Fisher-Yates shuffle
  function shuffle(array) {
    const shuffled = [...array];
    for (let i = shuffled.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1));
      [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
    }
    return shuffled;
  }

  let shuffledCards = shuffle(cards);
  let currentIndex = 0;
  let displayIndex = 0; // The card currently displayed
  let prevIndex = null; // The card sliding out
  let prevFlipped = false; // Whether the outgoing card was flipped
  let isFlipped = false;
  let isComplete = false;
  let isTransitioning = false;
  let slideDirection = null; // 'next' or 'prev'

  const SLIDE_DURATION = 400; // slide animation duration

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
    
    // Clear previous card after animation completes
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

  $: currentCard = shuffledCards[displayIndex];
  $: prevCard_ = prevIndex !== null ? shuffledCards[prevIndex] : null;

  function restart() {
    shuffledCards = shuffle(cards);
    currentIndex = 0;
    isFlipped = false;
    isComplete = false;
  }

  function exit() {
    dispatch('exit');
  }
</script>

<div class="study-view">
  <div class="content-header">
    <button class="btn btn-ghost" onclick={exit}>← Exit Study</button>
    <h2>{deck.name}</h2>
    <span class="study-progress">{currentIndex + 1} / {shuffledCards.length}</span>
  </div>

  {#if !isComplete}
    <div class="study-area">
      <div class="flashcard-container">
        <!-- Outgoing card (slides out) -->
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
        
        <!-- Current card -->
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
      <p class="flip-hint">Click card to flip</p>
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
      <button class="btn btn-primary" onclick={restart}>Study Again</button>
    </div>
  {/if}
</div>

<style>
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

  /* Slide animations */
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
    from {
      transform: translateX(100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  @keyframes slideInFromLeft {
    from {
      transform: translateX(-100%);
      opacity: 0;
    }
    to {
      transform: translateX(0);
      opacity: 1;
    }
  }

  @keyframes slideOutToLeft {
    from {
      transform: translateX(0);
      opacity: 1;
    }
    to {
      transform: translateX(-100%);
      opacity: 0;
    }
  }

  @keyframes slideOutToRight {
    from {
      transform: translateX(0);
      opacity: 1;
    }
    to {
      transform: translateX(100%);
      opacity: 0;
    }
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

  /* Disable flip transition during slide animations */
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

  .complete-icon {
    font-size: 4rem;
    margin-bottom: 20px;
  }

  .study-complete h3 {
    font-family: var(--font-display);
    font-size: 2rem;
    margin-bottom: 12px;
  }

  .study-complete p {
    color: var(--text-secondary);
    margin-bottom: 24px;
  }

  @media (max-width: 768px) {
    .flashcard-container {
      height: 300px;
    }

    .flashcard-front p,
    .flashcard-back p {
      font-size: 1.3rem;
    }
  }
</style>

