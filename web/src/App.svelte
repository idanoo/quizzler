<script>
  import { isAuthenticated } from './lib/stores.js';
  import Auth from './components/Auth.svelte';
  import Main from './components/Main.svelte';
</script>

{#if $isAuthenticated}
  <Main />
{:else}
  <Auth />
{/if}

