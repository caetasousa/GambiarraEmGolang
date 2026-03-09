<script lang="ts">
    import { createEventDispatcher, onMount } from "svelte";
    import { fade, fly } from "svelte/transition";

    export let show = false;
    export let title = "";
    export let maxWidth = "max-w-md";

    export let zIndex = "z-50";

    const dispatch = createEventDispatcher();

    function close() {
        dispatch("close");
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Escape" && show) {
            close();
        }
    }
</script>

<svelte:window on:keydown={handleKeydown} />

{#if show}
    <div
        class="fixed inset-0 {zIndex} flex items-center justify-center p-4 sm:p-6"
        role="dialog"
        aria-modal="true"
    >
        <!-- Backdrop -->
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div
            class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity"
            transition:fade={{ duration: 200 }}
            on:click={close}
        ></div>

        <!-- Modal Content -->
        <div
            class="relative w-full {maxWidth} max-h-[90vh] flex flex-col bg-[hsl(var(--bs-card))] rounded-2xl shadow-2xl border border-border-light dark:border-border-dark overflow-hidden z-10"
            transition:fly={{ y: 20, duration: 300, opacity: 0 }}
        >
            <div
                class="px-4 py-3 sm:px-6 sm:py-4 border-b border-border-light dark:border-border-dark flex items-center justify-between shrink-0"
            >
                <h3
                    class="text-lg font-bold text-gray-900 dark:text-white font-display"
                >
                    {title}
                </h3>
                <button
                    on:click={close}
                    class="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors rounded-full hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                    <span class="material-symbols-outlined">close</span>
                </button>
            </div>

            <div class="p-4 sm:p-6 overflow-y-auto">
                <slot name="body" />
            </div>

            {#if $$slots.footer}
                <div
                    class="px-4 py-3 sm:px-6 sm:py-4 bg-[hsl(var(--bs-muted))]/50 border-t border-border-light dark:border-border-dark flex justify-end space-x-3 shrink-0"
                >
                    <slot name="footer" />
                </div>
            {/if}
        </div>
    </div>
{/if}

<style>
    /* Prevent body scroll when modal is open */
    :global(body.modal-open) {
        overflow: hidden;
    }
</style>
