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
  let isFlipped = false;
  let isComplete = false;

  function flipCard() {
    isFlipped = !isFlipped;
  }

  function nextCard() {
    if (currentIndex < shuffledCards.length - 1) {
      currentIndex++;
      isFlipped = false;
    } else {
      isComplete = true;
    }
  }

  function prevCard() {
    if (currentIndex > 0) {
      currentIndex--;
      isFlipped = false;
    }
  }

  function restart() {
    shuffledCards = shuffle(cards);
    currentIndex = 0;
    isFlipped = false;
    isComplete = false;
  }

  function exit() {
    dispatch('exit');
  }

  $: currentCard = shuffledCards[currentIndex];
</script>

<div class="study-view">
  <div class="content-header">
    <button class="btn btn-ghost" onclick={exit}>← Exit Study</button>
    <h2>{deck.name}</h2>
    <span class="study-progress">{currentIndex + 1} / {shuffledCards.length}</span>
  </div>

  {#if !isComplete}
    <div class="study-area">
      <button class="flashcard" class:flipped={isFlipped} onclick={flipCard}>
        <div class="flashcard-inner">
          <div class="flashcard-front">
            <p>{currentCard.front}</p>
          </div>
          <div class="flashcard-back">
            <p>{currentCard.back}</p>
          </div>
        </div>
      </button>
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

  .flashcard {
    width: 100%;
    max-width: 600px;
    height: 380px;
    perspective: 1000px;
    cursor: pointer;
    margin-bottom: 16px;
    background: none;
    border: none;
    padding: 0;
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
    .flashcard {
      height: 300px;
    }

    .flashcard-front p,
    .flashcard-back p {
      font-size: 1.3rem;
    }
  }
</style>

