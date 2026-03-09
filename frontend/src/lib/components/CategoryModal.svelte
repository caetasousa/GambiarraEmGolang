<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import Modal from "$lib/components/Modal.svelte";
    import { validateCategoryName } from "$lib/utils/validation";

    /** Dispatcher para eventos do componente */
    const dispatch = createEventDispatcher<{
        add: string;
    }>();

    /** Controla visibilidade do modal */
    export let show: boolean = false;

    /** Lista de categorias existentes */
    export let existingCategories: string[] = [];

    /** Nome da nova categoria */
    let newName: string = "";

    /** Mensagem de erro */
    let error: string = "";

    /**
     * Reseta o estado do modal
     */
    function reset() {
        newName = "";
        error = "";
    }

    /**
     * Manipula o fechamento do modal
     */
    function handleClose() {
        show = false;
        reset();
    }

    /**
     * Adiciona nova categoria
     */
    function handleAdd() {
        const validation = validateCategoryName(newName, existingCategories);

        if (!validation.valid) {
            error = validation.error || "";
            return;
        }

        dispatch("add", newName.trim());
        handleClose();
    }

    /**
     * Permite adicionar com Enter
     */
    function handleKeydown(event: KeyboardEvent) {
        if (event.key === "Enter") {
            handleAdd();
        }
    }
</script>

<Modal {show} title="Nova Categoria" on:close={handleClose}>
    <div slot="body" class="space-y-4">
        <p class="text-sm text-gray-500 dark:text-gray-400">
            Digite o nome da nova categoria para organizar seus serviços.
        </p>
        <div>
            <label
                class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5"
                for="new-category"
            >
                Nome da Categoria
            </label>
            <input
                class="block w-full px-3 py-2 border border-border-light dark:border-border-dark rounded-lg bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-brand-orange transition-all sm:text-sm"
                id="new-category"
                placeholder="Ex: Penteados"
                type="text"
                bind:value={newName}
                on:keydown={handleKeydown}
            />
            {#if error}
                <p class="mt-1 text-xs text-red-500">{error}</p>
            {/if}
        </div>
    </div>
    <div slot="footer" class="flex space-x-3 w-full sm:w-auto">
        <button
            on:click={handleClose}
            class="flex-1 sm:flex-none px-6 py-2.5 border border-border-light dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 font-semibold hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
        >
            Cancelar
        </button>
        <button
            on:click={handleAdd}
            class="flex-1 sm:flex-none px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-lg font-semibold shadow-md transition-all active:scale-95"
        >
            Adicionar
        </button>
    </div>
</Modal>
