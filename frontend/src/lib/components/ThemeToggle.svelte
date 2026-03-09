<script lang="ts">
    import { theme, toggleTheme } from "$lib/stores/theme";
</script>

<button
    on:click={toggleTheme}
    class="p-2 rounded-full text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 focus:outline-none transition-colors"
    aria-label="Toggle Dark Mode"
>
    {#if $theme === "light"}
        <span class="material-symbols-outlined">dark_mode</span>
    {:else}
        <span class="material-symbols-outlined">light_mode</span>
    {/if}
</button>
