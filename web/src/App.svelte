<script>
  import { isAuthenticated } from './lib/stores.js';
  import Auth from './components/Auth.svelte';
  import Main from './components/Main.svelte';
  import PublicStudy from './components/PublicStudy.svelte';

  // Check if we're on a public study route
  function getPublicStudyDeckId() {
    const match = window.location.pathname.match(/^\/study\/(\d+)$/);
    return match ? parseInt(match[1]) : null;
  }

  let publicDeckId = getPublicStudyDeckId();

  // Update on navigation
  if (typeof window !== 'undefined') {
    window.addEventListener('popstate', () => {
      publicDeckId = getPublicStudyDeckId();
    });
  }
</script>

{#if $isAuthenticated}
  <Main />
{:else if publicDeckId}
  <PublicStudy deckId={publicDeckId} />
{:else}
  <Auth />
{/if}

