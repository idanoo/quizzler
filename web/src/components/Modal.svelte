<script>
  import { createEventDispatcher } from 'svelte';

  export let title = '';

  const dispatch = createEventDispatcher();

  function close() {
    dispatch('close');
  }

  function handleBackdropClick(e) {
    if (e.target === e.currentTarget) {
      close();
    }
  }

  function handleKeydown(e) {
    if (e.key === 'Escape') {
      close();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="modal-backdrop" onclick={handleBackdropClick} role="dialog" aria-modal="true">
  <div class="modal-content">
    <div class="modal-header">
      <h3>{title}</h3>
      <button class="modal-close" onclick={close} aria-label="Close">&times;</button>
    </div>
    <div class="modal-body">
      <slot />
    </div>
  </div>
</div>

<style>
  .modal-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.7);
    backdrop-filter: blur(4px);
    z-index: 1000;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 20px;
    animation: fadeIn 0.2s ease;
  }

  .modal-content {
    background: var(--bg-surface);
    border: 1px solid var(--border-normal);
    border-radius: var(--radius-lg);
    width: 100%;
    max-width: 500px;
    max-height: 90vh;
    overflow: auto;
    box-shadow: var(--shadow-lg);
    animation: slideUp 0.3s ease;
  }

  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .modal-header h3 {
    font-family: var(--font-display);
    font-size: 1.3rem;
    font-weight: 600;
  }

  .modal-close {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 1.5rem;
    cursor: pointer;
    padding: 4px 8px;
    transition: color var(--transition-fast);
  }

  .modal-close:hover {
    color: var(--text-primary);
  }

  .modal-body {
    padding: 24px;
  }
</style>

