<script lang="ts">
    import Modal from "$lib/components/Modal.svelte";
    import { validateCategoryName } from "$lib/utils/validation";
    import type { Snippet } from "svelte";

    let {
        show = $bindable(false),
        existingCategories = [],
        onadd,
    }: {
        show?: boolean;
        existingCategories?: string[];
        onadd?: (name: string) => void;
    } = $props();

    let newName = $state("");
    let error = $state("");

    function reset() {
        newName = "";
        error = "";
    }

    function handleClose() {
        show = false;
        reset();
    }

    function handleAdd() {
        const validation = validateCategoryName(newName, existingCategories);
        if (!validation.valid) {
            error = validation.error || "";
            return;
        }
        onadd?.(newName.trim());
        handleClose();
    }

    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Enter") handleAdd();
    }
</script>

{#snippet bodyContent()}
    <div class="space-y-4">
        <p class="text-sm text-gray-500 dark:text-gray-400">
            Digite o nome da nova categoria para organizar seus serviços.
        </p>
        <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5" for="new-category">
                Nome da Categoria
            </label>
            <input
                class="block w-full px-3 py-2 border border-border-light dark:border-border-dark rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-brand-orange transition-all sm:text-sm"
                id="new-category"
                placeholder="Ex: Penteados"
                type="text"
                bind:value={newName}
                onkeydown={handleKeydown}
            />
            {#if error}
                <p class="mt-1 text-xs text-red-500">{error}</p>
            {/if}
        </div>
    </div>
{/snippet}

{#snippet footerContent()}
    <div class="flex space-x-3 w-full sm:w-auto">
        <button
            onclick={handleClose}
            class="flex-1 sm:flex-none px-6 py-2.5 border border-border-light dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 font-semibold hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
            Cancelar
        </button>
        <button
            onclick={handleAdd}
            class="flex-1 sm:flex-none px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-lg font-semibold shadow-md transition-all active:scale-95"
        >
            Adicionar
        </button>
    </div>
{/snippet}

<Modal {show} title="Nova Categoria" onclose={handleClose} children={bodyContent} footer={footerContent} />
