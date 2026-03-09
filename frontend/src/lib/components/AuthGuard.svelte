<script lang="ts">
    import { isAuthenticated } from "$lib/stores/auth";
    import { goto } from "$app/navigation";
    import { onMount } from "svelte";

    onMount(() => {
        const unsubscribe = isAuthenticated.subscribe((auth) => {
            if (!auth) {
                goto("/login");
            }
        });
        return unsubscribe;
    });
</script>

{#if $isAuthenticated}
    <slot />
{/if}
