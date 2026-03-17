<script lang="ts">
  import { auth } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';

  let { allowedRoles = [], children }: { allowedRoles?: string[]; children: import('svelte').Snippet } = $props();

  $effect(() => {
    if (auth.user && allowedRoles.length > 0 && !allowedRoles.includes(auth.user.role)) {
      goto('/');
    }
  });
</script>

{#if !auth.user || allowedRoles.length === 0 || allowedRoles.includes(auth.user.role)}
  {@render children()}
{/if}
