<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";

    interface Service {
        ID: string;
        Nome: string;
        Categoria: string;
        DuracaoPadrao: number;
        Preco: number;
        ImagemUrl: string;
    }

    interface ApiResponse {
        data: Service[];
        limit: number;
        page: number;
        total: number;
    }

    // Estado dos serviços
    let services: Service[] = [];
    let featuredServices: Service[] = [];
    let loading = true;
    let error = "";

    // Paginação
    let page = 1;
    let limit = 5; // Quantos serviços queremos mostrar por página
    let apiLimit = 9; // Quantos buscamos da API (5 + 4 featured que serão filtrados)
    let totalItems = 0;

    // Calcula o total de páginas considerando que 4 serviços serão filtrados da primeira página
    $: totalPages = Math.ceil((totalItems - 4) / limit);
    $: hasMore = page < totalPages;

    // Filtra os serviços para remover os que estão em destaque
    $: filteredServices = services.filter(
        (service) =>
            !featuredServices.some((featured) => featured.ID === service.ID),
    );

    // Filtro de categoria
    let selectedCategory = "Todos os Serviços";
    let categories = [
        "Todos os Serviços",
        "Cabelo",
        "Unhas",
        "Estética Facial",
        "Estética Corporal",
        "Depilação",
        "Maquiagem",
        "Massagem",
        "Tratamentos",
    ];

    async function fetchFeaturedServices() {
        try {
            const response = await fetch(`/api/v1/catalogos?page=1&limit=4`);
            if (response.ok) {
                const data: ApiResponse = await response.json();
                featuredServices = data.data || [];
            }
        } catch (err) {
            console.error("Erro ao carregar serviços em destaque:", err);
        }
    }

    async function fetchServices() {
        loading = true;
        error = "";
        try {
            // Busca mais serviços para compensar os que serão filtrados
            let url = `/api/v1/catalogos?page=${page}&limit=${apiLimit}`;

            // Adiciona filtro de categoria se não for "Todos os Serviços"
            if (selectedCategory !== "Todos os Serviços") {
                url += `&categoria=${encodeURIComponent(selectedCategory)}`;
            }

            const response = await fetch(url);

            if (response.ok) {
                const data: ApiResponse = await response.json();
                services = data.data || [];
                totalItems = data.total || 0;
            } else {
                error = "Erro ao carregar serviços: " + response.statusText;
            }
        } catch (err) {
            console.error("Erro na requisição:", err);
            error = "Serviço indisponível no momento.";
        } finally {
            loading = false;
        }
    }

    function loadMore() {
        if (hasMore) {
            page++;
            fetchServices();
        }
    }

    function selectCategory(category: string) {
        selectedCategory = category;
        page = 1; // Reset para primeira página ao mudar categoria
        fetchServices();
    }

    function getCategoryIcon(categoria: string): string {
        const icons: Record<string, string> = {
            Cabelo: "content_cut",
            Unhas: "front_hand",
            "Estética Facial": "face",
            "Estética Corporal": "spa",
            Massagem: "self_improvement",
            Depilação: "spa",
            Maquiagem: "palette",
            Tratamentos: "healing",
        };
        return icons[categoria] || "spa";
    }

    function getCategoryColor(categoria: string): string {
        const colors: Record<string, string> = {
            Cabelo: "bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300",
            Unhas: "bg-pink-100 dark:bg-pink-900 text-pink-700 dark:text-pink-300",
            "Estética Facial":
                "bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300",
            "Estética Corporal":
                "bg-emerald-100 dark:bg-emerald-900 text-emerald-700 dark:text-emerald-300",
            Massagem:
                "bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300",
            Depilação:
                "bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300",
            Maquiagem:
                "bg-rose-100 dark:bg-rose-900 text-rose-700 dark:text-rose-300",
            Tratamentos:
                "bg-cyan-100 dark:bg-cyan-900 text-cyan-700 dark:text-cyan-300",
        };
        return (
            colors[categoria] ||
            "bg-orange-100 dark:bg-orange-900 text-orange-700 dark:text-orange-300"
        );
    }

    function formatPrice(priceInCents: number): string {
        return (priceInCents / 100).toFixed(2).replace(".", ",");
    }

    onMount(() => {
        fetchFeaturedServices();
        fetchServices();
    });
</script>

<div
    class="bg-[hsl(var(--bs-background))] text-gray-800 dark:text-gray-200 font-sans antialiased transition-colors duration-200"
>
    <div class="flex h-screen overflow-hidden">
        <Sidebar />
        <main class="flex-1 flex flex-col h-full overflow-hidden relative">
            <DashboardNavbar title="Catálogo de Serviços" />
            <div class="flex-1 overflow-y-auto p-6 md:p-8">
                <div
                    class="relative rounded-2xl overflow-hidden bg-gradient-to-r from-orange-500 to-amber-600 dark:from-orange-600 dark:to-amber-800 shadow-lg mb-10 text-white"
                >
                    <div
                        class="absolute inset-0 bg-[url('https://www.transparenttextures.com/patterns/cubes.png')] opacity-10"
                    ></div>
                    <div
                        class="relative z-10 p-8 md:p-12 flex flex-col md:flex-row items-center justify-between gap-8"
                    >
                        <div class="max-w-xl space-y-4">
                            <span
                                class="inline-block px-3 py-1 bg-white/20 backdrop-blur-sm rounded-full text-xs font-semibold tracking-wide uppercase"
                                >Oferta da Semana</span
                            >
                            <h1
                                class="text-3xl md:text-5xl font-bold font-display leading-tight"
                            >
                                Renove sua beleza com nossos especialistas
                            </h1>
                            <p class="text-white/90 text-lg max-w-md">
                                Descubra os melhores tratamentos capilares e
                                estéticos selecionados para você.
                            </p>
                            <div class="pt-4 flex gap-3">
                                <button
                                    class="px-6 py-3 bg-white text-orange-600 font-semibold rounded-lg shadow-sm hover:bg-gray-50 transition"
                                    >Agendar Agora</button
                                >
                                <button
                                    class="px-6 py-3 bg-transparent border border-white text-white font-semibold rounded-lg hover:bg-white/10 transition"
                                    >Ver Promoções</button
                                >
                            </div>
                        </div>
                        <div class="hidden md:block relative w-80 h-64">
                            <div
                                class="absolute top-0 right-0 w-64 h-48 bg-gray-900 rounded-lg rotate-3 shadow-2xl overflow-hidden border-4 border-white/20"
                            >
                                <img
                                    alt="Hair styling"
                                    class="w-full h-full object-cover opacity-80"
                                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuCYxtIrdDpKj4dz4C07_cfWm0fxJ2SvIE1Zx3hpSCCNfZ1s81dDodesar6jXm1yeDKsi5KDsz5teIegxCm_3tPCdafgyTDbGD2dvKkGff75y7ghdBAP8YrsSshcpmAl2T8XYzbHFKFhIDtvz-CEFaXRGrsXoNics974oaCLIhQDimGQNRYyCpGZ8FzS9vfNt_NanQ2xDXKAUgG4osP3ifwEvHBY_QFu0pvgHpUjsKpX8_zJarrWTLTaNxWvABA3eivnvNdyPnRmypg"
                                />
                            </div>
                            <div
                                class="absolute bottom-0 left-4 w-48 h-36 bg-orange-700 rounded-lg -rotate-6 shadow-xl overflow-hidden border-4 border-white/20 z-20"
                            >
                                <img
                                    alt="Makeup"
                                    class="w-full h-full object-cover"
                                    src="https://lh3.googleusercontent.com/aida-public/AB6AXuCVqEmwaujVPrnjTd8HxPLpAqQMN9g6nrZiCfnHhptuRvhgO5h7ZS4l0G1CD70g0NkKNEK0bWQJoSCvcMJ3wGHOGJUzK-fER2p6jTMlIKOHvwf-R6VzkS_5shKPV8K79XldmB851IJLkvg_DUvRgrlmx9FJCka0ZuKnQRjLRYEHbitvjZ9idVnNaC3nZ4qoHzW-6en_s56qqbR1NiXH5cRi1xdfwOE7X45e-4lSD3HDhNbGX1ZI4awIB6hNaTxyU5vnkKcng-wQuHg"
                                />
                            </div>
                            <span
                                class="material-icons-outlined absolute -top-4 left-10 text-white/40 text-6xl animate-pulse"
                                >spa</span
                            >
                        </div>
                    </div>
                </div>
                <div
                    class="flex overflow-x-auto pb-4 gap-3 mb-8 no-scrollbar scroll-smooth"
                >
                    {#each categories as category}
                        <button
                            on:click={() => selectCategory(category)}
                            class="whitespace-nowrap px-4 py-2 {selectedCategory ===
                            category
                                ? 'bg-primary text-white shadow-md shadow-primary/30'
                                : 'bg-[hsl(var(--bs-card))] border border-border-light dark:border-border-dark text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800'} text-sm font-medium rounded-full transition-colors"
                        >
                            {category}
                        </button>
                    {/each}
                </div>
                <div class="mb-10">
                    <div class="flex items-center justify-between mb-6">
                        <h2
                            class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2"
                        >
                            <span class="material-symbols-outlined text-primary"
                                >local_fire_department</span
                            >
                            Mais Procurados
                        </h2>
                    </div>

                    {#if loading}
                        <div class="flex justify-center py-12">
                            <span
                                class="material-icons animate-spin text-primary text-4xl"
                                >sync</span
                            >
                        </div>
                    {:else if error}
                        <div
                            class="bg-red-50 text-red-700 p-4 rounded-md border border-red-200"
                        >
                            {error}
                        </div>
                    {:else if featuredServices.length === 0}
                        <div
                            class="text-center py-12 text-gray-500 dark:text-gray-400"
                        >
                            Nenhum serviço disponível no momento.
                        </div>
                    {:else}
                        <div
                            class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6"
                        >
                            {#each featuredServices as service}
                                <a
                                    href="/listagem/{service.ID}"
                                    class="group bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark overflow-hidden hover:shadow-lg transition-all duration-300 flex flex-col cursor-pointer"
                                >
                                    <div class="relative h-48 overflow-hidden">
                                        {#if service.ImagemUrl}
                                            <img
                                                alt={service.Nome}
                                                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500"
                                                src={service.ImagemUrl}
                                                on:error={(e) => {
                                                    (
                                                        e.target as HTMLImageElement
                                                    ).src =
                                                        "https://images.unsplash.com/photo-1560066984-138dadb4c035?auto=format&fit=crop&q=80&w=600";
                                                }}
                                            />
                                        {:else}
                                            <div
                                                class="w-full h-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-6xl text-gray-400"
                                                    >image</span
                                                >
                                            </div>
                                        {/if}
                                        <div
                                            class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/60 to-transparent h-20"
                                        ></div>
                                    </div>
                                    <div class="p-4 flex-1 flex flex-col">
                                        <div
                                            class="mb-1 text-xs font-semibold text-primary uppercase tracking-wide"
                                        >
                                            {service.Categoria}
                                        </div>
                                        <h3
                                            class="text-lg font-bold text-gray-900 dark:text-white mb-2"
                                        >
                                            {service.Nome}
                                        </h3>
                                        <p
                                            class="text-sm text-gray-500 dark:text-gray-400 line-clamp-2 mb-4"
                                        >
                                            Duração: {service.DuracaoPadrao} minutos
                                        </p>
                                        <div
                                            class="mt-auto flex flex-col gap-3 pt-4 border-t border-border-light dark:border-border-dark"
                                        >
                                            <div class="flex items-center justify-between">
                                                <div>
                                                    <span
                                                        class="text-xs text-gray-400 block"
                                                        >A partir de</span
                                                    >
                                                    <span
                                                        class="text-lg font-bold text-gray-900 dark:text-white"
                                                        >R$ {formatPrice(
                                                            service.Preco,
                                                        )}</span
                                                    >
                                                </div>
                                            </div>
                                            <button
                                                on:click|preventDefault|stopPropagation={() => window.location.href = `/clientes/agendamento?serviceId=${service.ID}`}
                                                class="w-full px-4 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-lg font-semibold shadow-md transition-all duration-200 flex items-center justify-center gap-2 group-hover:shadow-lg active:scale-95"
                                            >
                                                <span class="material-symbols-outlined text-[18px]">
                                                    event_available
                                                </span>
                                                Agendar Agora
                                            </button>
                                        </div>
                                    </div>
                                </a>
                            {/each}
                        </div>
                    {/if}
                </div>

                <div class="mb-10">
                    <h2
                        class="text-xl font-bold text-gray-900 dark:text-white mb-6"
                    >
                        Todos os Serviços
                    </h2>

                    {#if loading}
                        <div class="flex justify-center py-12">
                            <span
                                class="material-icons animate-spin text-primary text-4xl"
                                >sync</span
                            >
                        </div>
                    {:else if error}
                        <div
                            class="bg-red-50 text-red-700 p-4 rounded-md border border-red-200"
                        >
                            {error}
                        </div>
                    {:else if filteredServices.length === 0}
                        <div
                            class="text-center py-12 text-gray-500 dark:text-gray-400"
                        >
                            Nenhum serviço cadastrado ainda.
                        </div>
                    {:else}
                        <div
                            class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark divide-y divide-border-light dark:divide-border-dark"
                        >
                            {#each filteredServices as service}
                                <div
                                    class="p-6 flex flex-col md:flex-row gap-6 items-start md:items-center hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors"
                                >
                                    <div
                                        class="w-20 h-20 rounded-lg bg-gray-200 dark:bg-gray-700 overflow-hidden flex-shrink-0"
                                    >
                                        {#if service.ImagemUrl}
                                            <img
                                                alt={service.Nome}
                                                class="w-full h-full object-cover"
                                                src={service.ImagemUrl}
                                                on:error={(e) => {
                                                    (
                                                        e.target as HTMLImageElement
                                                    ).src =
                                                        "https://images.unsplash.com/photo-1560066984-138dadb4c035?auto=format&fit=crop&q=80&w=200";
                                                }}
                                            />
                                        {:else}
                                            <div
                                                class="w-full h-full flex items-center justify-center"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-3xl text-gray-400"
                                                    >image</span
                                                >
                                            </div>
                                        {/if}
                                    </div>
                                    <div class="flex-1">
                                        <div
                                            class="flex items-center gap-2 mb-1"
                                        >
                                            <a href="/listagem/{service.ID}">
                                                <h3
                                                    class="text-lg font-bold text-gray-900 dark:text-white hover:text-primary transition-colors cursor-pointer"
                                                >
                                                    {service.Nome}
                                                </h3>
                                            </a>
                                            <span
                                                class="px-2 py-0.5 {getCategoryColor(
                                                    service.Categoria,
                                                )} text-xs font-semibold rounded-full"
                                                >{service.Categoria}</span
                                            >
                                        </div>
                                        <div
                                            class="flex items-center gap-4 text-xs text-gray-500 mt-2"
                                        >
                                            <span
                                                class="flex items-center gap-1"
                                                ><span
                                                    class="material-symbols-outlined text-sm"
                                                    >schedule</span
                                                >
                                                {service.DuracaoPadrao} min</span
                                            >
                                        </div>
                                    </div>
                                    <div
                                        class="flex flex-col items-end gap-2 w-full md:w-auto"
                                    >
                                        <div
                                            class="text-xl font-bold text-gray-900 dark:text-white"
                                        >
                                            R$ {formatPrice(service.Preco)}
                                        </div>
                                        <a
                                            href="/listagem/{service.ID}"
                                            class="w-full md:w-auto px-4 py-2 bg-primary hover:bg-orange-600 text-white text-sm font-medium rounded-lg transition-colors text-center inline-block"
                                        >
                                            Agendar
                                        </a>
                                    </div>
                                </div>
                            {/each}
                        </div>

                        <!-- Pagination -->
                        {#if totalPages > 1}
                            <div
                                class="mt-8 flex justify-center items-center gap-4"
                            >
                                <button
                                    on:click|preventDefault={() => {
                                        page--;
                                        fetchServices();
                                    }}
                                    disabled={page === 1}
                                    class="px-4 py-2 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                >
                                    Anterior
                                </button>
                                <span class="text-gray-600 dark:text-gray-400">
                                    Página {page} de {totalPages}
                                </span>
                                <button
                                    on:click|preventDefault={() => {
                                        page++;
                                        fetchServices();
                                    }}
                                    disabled={page === totalPages}
                                    class="px-4 py-2 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                >
                                    Próxima
                                </button>
                            </div>
                        {/if}
                    {/if}
                </div>

                <!-- Footer -->
                <div
                    class="mt-12 border-t border-border-light dark:border-border-dark pt-8 pb-4 text-center text-sm text-gray-500"
                >
                    <p>© 2023 BellaVita. Todos os direitos reservados.</p>
                    <div class="flex justify-center gap-4 mt-2">
                        <a class="hover:text-primary" href="/">Termos</a>
                        <a class="hover:text-primary" href="/">Privacidade</a>
                        <a class="hover:text-primary" href="/">Ajuda</a>
                    </div>
                </div>
            </div>
        </main>
    </div>
</div>
