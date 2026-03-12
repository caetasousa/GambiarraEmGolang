<script lang="ts">
  import { user } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  export let allowedRoles: string[] = [];

  onMount(() => {
    if ($user && allowedRoles.length > 0 && !allowedRoles.includes($user.role)) {
      goto('/');
    }
  });
</script>

{#if !$user || allowedRoles.length === 0 || allowedRoles.includes($user.role)}
  <slot />
{/if}
