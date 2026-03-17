<script lang="ts">
    import { auth } from "$lib/stores/auth.svelte";
    import { goto } from "$app/navigation";

    let { children } = $props();

    $effect(() => {
        if (!auth.isAuthenticated) {
            goto("/login");
        }
    });
</script>

{#if auth.isAuthenticated}
    {@render children()}
{/if}
