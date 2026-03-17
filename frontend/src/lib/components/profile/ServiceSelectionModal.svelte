<script lang="ts">
    import Modal from "$lib/components/Modal.svelte";
    import { fetchApi } from "$lib/utils/api";
    import type { Snippet } from "svelte";

    let {
        show = $bindable(false),
        providerId,
        existingServiceIds = [],
        onadd,
        onclose,
    }: {
        show?: boolean;
        providerId: string;
        existingServiceIds?: string[];
        onadd?: (service: Record<string, unknown>) => void;
        onclose?: () => void;
    } = $props();

    let catalogServices = $state<Record<string, unknown>[]>([]);
    let loading = $state(false);
    let error = $state("");
    let searchTerm = $state("");

    async function fetchCatalog() {
        loading = true;
        error = "";
        try {
            const response = await fetchApi("/api/v1/catalogos?limit=100");
            if (response.ok) {
                const data = await response.json();
                catalogServices = data.data || [];
            } else {
                error = "Erro ao carregar catálogo.";
            }
        } catch {
            error = "Erro de conexão.";
        } finally {
            loading = false;
        }
    }

    $effect(() => {
        if (show && catalogServices.length === 0) {
            fetchCatalog();
        }
    });

    function close() {
        show = false;
        searchTerm = "";
        onclose?.();
    }

    async function addService(service: Record<string, unknown>) {
        try {
            const response = await fetchApi(`/api/v1/prestadores/${providerId}/servicos`, {
                method: "POST",
                body: JSON.stringify({ servico_id: service.ID }),
            });
            if (!response.ok) {
                console.warn("API add service failed or not implemented");
            }
            onadd?.(service);
        } catch {
            // silently handled
        }
    }

    let filteredServices = $derived(catalogServices.filter(s => {
        const matchesSearch = (s.Nome as string).toLowerCase().includes(searchTerm.toLowerCase()) ||
                              (s.Categoria as string).toLowerCase().includes(searchTerm.toLowerCase());
        const isNotOwned = !existingServiceIds.includes(s.ID as string);
        return matchesSearch && isNotOwned;
    }));
</script>

{#snippet bodyContent()}
    <div class="space-y-4">
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
                                <img src={service.ImagemUrl as string} alt={service.Nome as string} class="w-10 h-10 rounded-md object-cover bg-gray-100" />
                            {:else}
                                <div class="w-10 h-10 rounded-md bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-400">
                                    <span class="material-symbols-outlined font-size-20">spa</span>
                                </div>
                            {/if}
                            <div>
                                <h4 class="text-sm font-semibold text-gray-900 dark:text-white">{service.Nome}</h4>
                                <p class="text-xs text-gray-500">{service.Categoria} • R$ {((service.Preco as number) / 100).toFixed(2).replace('.', ',')}</p>
                            </div>
                        </div>
                        <button
                            onclick={() => addService(service)}
                            class="text-sm font-medium text-brand-orange hover:bg-orange-50 dark:hover:bg-orange-900/10 px-3 py-1.5 rounded-md transition-colors"
                        >
                            Adicionar
                        </button>
                    </div>
                {/each}
            {/if}
        </div>
    </div>
{/snippet}

{#snippet footerContent()}
    <button
        onclick={close}
        class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
    >
        Fechar
    </button>
{/snippet}

<Modal {show} title="Adicionar Serviço" onclose={close} maxWidth="max-w-2xl" children={bodyContent} footer={footerContent} />
