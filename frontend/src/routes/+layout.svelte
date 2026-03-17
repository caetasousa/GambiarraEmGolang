<script lang="ts">
  import '../app.css';
  import { theme } from '$lib/stores/theme.svelte';

  let { children } = $props();
</script>

{@render children()}
