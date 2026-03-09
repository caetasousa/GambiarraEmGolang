<script lang="ts">
    import { createEventDispatcher, onMount } from "svelte";
    import Modal from "$lib/components/Modal.svelte";

    export let show = false;
    export let providerId: string;
    export let existingServiceIds: string[] = []; // IDs of services the provider already has

    const dispatch = createEventDispatcher();

    let catalogServices: any[] = [];
    let loading = false;
    let error = "";
    let searchTerm = "";

    async function fetchCatalog() {
        loading = true;
        error = "";
        try {
            const response = await fetch("/api/v1/catalogos?limit=100");
            if (response.ok) {
                const data = await response.json();
                catalogServices = data.data || [];
            } else {
                error = "Erro ao carregar catálogo.";
            }
        } catch (e) {
            console.error(e);
            error = "Erro de conexão.";
        } finally {
            loading = false;
        }
    }

    // Fetch catalog when modal opens
    $: if (show && catalogServices.length === 0) {
        fetchCatalog();
    }

    function close() {
        show = false;
        searchTerm = "";
        dispatch("close");
    }

    async function addService(service: any) {
        // Optimistic update handled by parent or re-fetch?
        // Let's assume we call API here.
        try {
            // Check if endpoint exists, otherwise mock or adapt
            const response = await fetch(`/api/v1/prestadores/${providerId}/servicos`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ servico_id: service.ID }) // Adjust payload as needed
            });

            if (!response.ok) {
                // Determine if error or just mock success
                console.warn("API add service failed or not implemented");
            }
            
            dispatch("add", service);
        } catch (e) {
            console.error(e);
        }
    }

    $: filteredServices = catalogServices.filter(s => {
        const matchesSearch = s.Nome.toLowerCase().includes(searchTerm.toLowerCase()) || 
                              s.Categoria.toLowerCase().includes(searchTerm.toLowerCase());
        const isNotOwned = !existingServiceIds.includes(s.ID);
        return matchesSearch && isNotOwned;
    });
</script>

<Modal {show} title="Adicionar Serviço" on:close={close} maxWidth="max-w-2xl">
    <div slot="body" class="space-y-4">
        <!-- Search -->
        <div class="relative">
            <span class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <span class="material-symbols-outlined text-gray-400">search</span>
            </span>
            <input
                type="text"
                bind:value={searchTerm}
                placeholder="Buscar serviços..."
                class="block w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
            />
        </div>

        <!-- List -->
        <div class="max-h-96 overflow-y-auto space-y-2">
            {#if loading}
                <div class="flex justify-center py-4">
                    <span class="material-symbols-outlined animate-spin text-brand-orange">sync</span>
                </div>
            {:else if error}
                <p class="text-red-500 text-center text-sm">{error}</p>
            {:else if filteredServices.length === 0}
                <p class="text-gray-500 text-center py-8">Nenhum serviço disponível encontrado.</p>
            {:else}
                {#each filteredServices as service}
                    <div class="flex items-center justify-between p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-brand-orange/50 transition-colors">
                        <div class="flex items-center space-x-3">
                            {#if service.ImagemUrl}
                                <img src={service.ImagemUrl} alt={service.Nome} class="w-10 h-10 rounded-md object-cover bg-gray-100" />
                            {:else}
                                <div class="w-10 h-10 rounded-md bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-400">
                                    <span class="material-symbols-outlined font-size-20">spa</span>
                                </div>
                            {/if}
                            <div>
                                <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{service.Nome}</h4>
                                <p class="text-xs text-gray-500">{service.Categoria} • R$ {(service.Preco / 100).toFixed(2).replace('.', ',')}</p>
                            </div>
                        </div>
                        <button
                            on:click={() => addService(service)}
                            class="text-sm font-medium text-brand-orange hover:bg-orange-50 dark:hover:bg-orange-900/10 px-3 py-1.5 rounded-md transition-colors"
                        >
                            Adicionar
                        </button>
                    </div>
                {/each}
            {/if}
        </div>
    </div>
    <div slot="footer">
        <button
            on:click={close}
            class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
        >
            Fechar
        </button>
    </div>
</Modal>
