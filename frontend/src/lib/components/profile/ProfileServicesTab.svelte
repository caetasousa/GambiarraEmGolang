<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import type { Service } from "$lib/types/profile"; // Import shared type

    // Use Service instead of ServiceView for consistency if possible, or map
    export let services: any[] = []; // Using any[] for now as the mapping logic is in parent
    export let servicesLoading = false;

    const dispatch = createEventDispatcher<{
        addService: void;
        removeService: { service: any };
    }>();
</script>

<div class="animate-fade-in space-y-6">
    <div class="flex justify-between items-center">
        <h3 class="text-lg font-bold text-gray-900 dark:text-white">
            Serviços Habilitados
        </h3>
        <!-- Catalog management restored -->
        <button
            on:click={() => dispatch("addService")}
            class="text-primary font-medium hover:underline text-sm flex items-center"
        >
            <span class="material-symbols-outlined text-sm mr-1">add</span>
            Adicionar Serviço
        </button>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {#if servicesLoading}
            <div class="col-span-full flex justify-center py-8">
                <span
                    class="material-symbols-outlined animate-spin text-brand-orange text-3xl"
                    >sync</span
                >
            </div>
        {:else if services.length === 0}
            <div class="col-span-full text-center py-8 text-gray-500">
                Nenhum serviço encontrado.
            </div>
        {:else}
            {#each services as service}
                <div
                    class="relative flex items-start p-4 rounded-xl border border-border-light dark:border-border-dark bg-[hsl(var(--bs-card))] shadow-sm transition-all hover:border-brand-orange/50 group"
                >
                    {#if service.ImagemUrl}
                        <div
                            class="h-12 w-12 rounded-lg bg-gray-100 overflow-hidden shrink-0 mr-3 border border-gray-200 dark:border-gray-700"
                        >
                            <img
                                src={service.ImagemUrl}
                                alt={service.Nome}
                                class="h-full w-full object-cover"
                            />
                        </div>
                    {/if}

                    <div class="min-w-0 flex-1">
                        <div class="flex justify-between items-start mb-1">
                            <h4
                                class="font-bold text-gray-900 dark:text-white text-sm line-clamp-1"
                                title={service.Nome}
                            >
                                {service.Nome}
                            </h4>
                            <span
                                class="text-brand-orange font-bold text-sm whitespace-nowrap ml-2"
                            >
                                R$ {(service.Preco / 100)
                                    .toFixed(2)
                                    .replace(".", ",")}
                            </span>
                        </div>
                        <p class="text-gray-500 dark:text-gray-400 text-xs">
                            {service.DuracaoPadrao} min • {service.Categoria}
                        </p>
                    </div>
                    <div class="ml-2 flex items-center self-center">
                        <button
                            on:click={() =>
                                dispatch("removeService", { service })}
                            class="p-1.5 text-gray-400 hover:text-red-500 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/10 transition-colors"
                            title="Remover Serviço"
                        >
                            <span class="material-symbols-outlined text-[20px]"
                                >delete</span
                            >
                        </button>
                    </div>
                </div>
            {/each}
        {/if}
    </div>
</div>

<style>
    .animate-fade-in {
        animation: fadeIn 0.3s ease-in-out;
    }
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
