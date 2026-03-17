<script lang="ts">
    import Sidebar from "$lib/components/Sidebar.svelte";
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { page } from "$app/stores";

    interface Service {
        id: string;
        nome: string;
        categoria: string;
        duracao_padrao: number;
        preco: number;
        image_url?: string;
    }

    let service = $state<Service | null>(null);
    let loading = $state(true);
    let error = $state("");

    function getCategoryColor(categoria: string): string {
        const colors: Record<string, string> = {
            "Cabelo": "bg-purple-100 dark:bg-purple-900 text-purple-700 dark:text-purple-300",
            "Unhas": "bg-pink-100 dark:bg-pink-900 text-pink-700 dark:text-pink-300",
            "Estética Facial": "bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300",
            "Estética Corporal": "bg-emerald-100 dark:bg-emerald-900 text-emerald-700 dark:text-emerald-300",
            "Massagem": "bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300",
            "Depilação": "bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300",
            "Maquiagem": "bg-rose-100 dark:bg-rose-900 text-rose-700 dark:text-rose-300",
            "Tratamentos": "bg-cyan-100 dark:bg-cyan-900 text-cyan-700 dark:text-cyan-300"
        };
        return colors[categoria] || "bg-orange-100 dark:bg-orange-900 text-orange-700 dark:text-orange-300";
    }

    function formatPrice(priceInCents: number): string {
        return (priceInCents / 100).toFixed(2).replace(".", ",");
    }

    async function fetchServiceDetails(id: string) {
        console.log("Fetching service details for ID:", id);
        loading = true;
        error = "";
        try {
            const url = `/api/v1/catalogos/${id}`;
            console.log("Request URL:", url);
            const response = await fetch(url);
            console.log("Response status:", response.status);

            if (response.ok) {
                const data = await response.json();
                console.log("Service data received:", data);
                service = data;
            } else if (response.status === 404) {
                error = "Serviço não encontrado.";
                console.error("Service not found (404)");
            } else {
                const errorText = await response.text();
                error = "Erro ao carregar detalhes do serviço.";
                console.error("Error response:", response.status, errorText);
            }
        } catch (err) {
            console.error("Erro na requisição:", err);
            error = "Serviço indisponível no momento.";
        } finally {
            loading = false;
        }
    }

    // Reactive statement that fetches data when the ID changes
    $effect(() => {
        const id = $page.params.id;
        if (id) {
            fetchServiceDetails(id);
        }
    });
</script>

<div
    class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <Sidebar />

    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <header
            class="h-16 bg-[hsl(var(--bs-card))] border-b border-border-light dark:border-border-dark flex items-center justify-between px-6 z-10"
        >
            <div class="flex items-center">
                <a
                    href="/catalogo/listagem"
                    class="flex items-center text-gray-500 hover:text-brand-orange transition-colors group"
                >
                    <span
                        class="material-symbols-outlined mr-2 group-hover:-translate-x-1 transition-transform"
                        >arrow_back</span
                    >
                    <span class="hidden sm:inline font-medium"
                        >Voltar para o catálogo</span
                    >
                </a>
            </div>
            <div class="flex items-center space-x-4">
                <ThemeToggle />
            </div>
        </header>

        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            {#if loading}
                <div class="flex justify-center items-center h-full">
                    <span class="material-symbols-outlined animate-spin text-brand-orange text-5xl">sync</span>
                </div>
            {:else if error}
                <div class="max-w-4xl mx-auto">
                    <div class="bg-red-50 dark:bg-red-900/20 text-red-700 dark:text-red-400 p-6 rounded-lg border border-red-200 dark:border-red-800">
                        <p class="text-lg font-semibold mb-2">Erro</p>
                        <p>{error}</p>
                    </div>
                </div>
            {:else if service}
                <div class="max-w-4xl mx-auto space-y-8">
                    <!-- Breadcrumb -->
                    <nav class="flex text-sm text-gray-500 dark:text-gray-400 mb-2">
                        <a href="/" class="hover:text-primary transition-colors"
                            >Início</a
                        >
                        <span class="mx-2">/</span>
                        <a
                            href="/catalogo/listagem"
                            class="hover:text-primary transition-colors">Catálogo</a
                        >
                        <span class="mx-2">/</span>
                        <span class="text-gray-900 dark:text-white font-medium"
                            >{service.nome}</span
                        >
                    </nav>

                    <!-- Service Hero Card -->
                    <div
                        class="bg-[hsl(var(--bs-card))] rounded-3xl shadow-xl border border-border-light dark:border-border-dark overflow-hidden transition-all duration-300"
                    >
                        <div
                            class="relative h-64 sm:h-80 bg-gray-200 dark:bg-gray-800"
                        >
                            {#if service.image_url}
                                <img
                                    src={service.image_url}
                                    alt={service.nome}
                                    class="w-full h-full object-cover"
                                    onerror={(e) => {
                                        (e.target as HTMLImageElement).src = "https://images.unsplash.com/photo-1560066984-138dadb4c035?auto=format&fit=crop&q=80&w=1200";
                                    }}
                                />
                            {:else}
                                <div class="w-full h-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                                    <span class="material-symbols-outlined text-9xl text-gray-400">image</span>
                                </div>
                            {/if}
                            <div
                                class="absolute inset-0 bg-gradient-to-t from-black/60 via-black/20 to-transparent"
                            ></div>
                            <div class="absolute bottom-8 left-8 right-8">
                                <div class="flex items-center gap-3 mb-3">
                                    <span
                                        class="px-3 py-1 bg-brand-orange text-white text-[10px] font-black uppercase tracking-widest rounded-full shadow-lg"
                                    >
                                        {service.categoria}
                                    </span>
                                </div>
                                <h1
                                    class="text-3xl md:text-5xl font-black text-white leading-tight drop-shadow-lg"
                                >
                                    {service.nome}
                                </h1>
                            </div>
                        </div>

                        <div class="p-8">
                            <div class="grid grid-cols-1 md:grid-cols-3 gap-12">
                                <div class="md:col-span-2 space-y-10">
                                    <!-- About Section -->
                                    <section>
                                        <h3
                                            class="text-xl font-black text-gray-900 dark:text-white flex items-center mb-6 uppercase tracking-tight"
                                        >
                                            <span
                                                class="material-symbols-outlined mr-3 text-brand-orange text-3xl"
                                                >description</span
                                            >
                                            Sobre o Serviço
                                        </h3>
                                        <p
                                            class="text-gray-600 dark:text-gray-300 leading-relaxed text-lg"
                                        >
                                            Serviço de {service.categoria.toLowerCase()} profissional com duração de {service.duracao_padrao} minutos.
                                            Oferecemos atendimento de qualidade com profissionais especializados para garantir a melhor experiência.
                                        </p>
                                    </section>

                                    <!-- Service Details -->
                                    <section>
                                        <h3
                                            class="text-xl font-black text-gray-900 dark:text-white flex items-center mb-6 uppercase tracking-tight"
                                        >
                                            <span
                                                class="material-symbols-outlined mr-3 text-brand-orange text-3xl"
                                                >info</span
                                            >
                                            Informações Importantes
                                        </h3>
                                        <div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-6">
                                            <div class="flex items-start gap-3">
                                                <span class="material-symbols-outlined text-blue-600 dark:text-blue-400 text-2xl">verified_user</span>
                                                <div>
                                                    <ul class="text-sm text-blue-800 dark:text-blue-400 space-y-1">
                                                        <li>• Chegue com 10 minutos de antecedência</li>
                                                        <li>• Cancelamentos devem ser feitos com 24h de antecedência</li>
                                                        <li>• Traga documentos de identificação</li>
                                                        <li>• Profissionais qualificados e certificados</li>
                                                    </ul>
                                                </div>
                                            </div>
                                        </div>
                                    </section>
                                </div>

                                <!-- Detail Sidebar -->
                                <div class="space-y-6">
                                    <div
                                        class="bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 rounded-3xl p-6 border border-border-light dark:border-border-dark shadow-inner"
                                    >
                                        <h4
                                            class="text-[10px] font-black text-gray-400 dark:text-gray-500 uppercase tracking-[0.2em] mb-6"
                                        >
                                            Especificações
                                        </h4>
                                        <div class="space-y-6">
                                            <!-- Categoria -->
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="h-10 w-10 rounded-xl bg-white dark:bg-[hsl(var(--bs-card))] flex items-center justify-center text-brand-orange shadow-sm border border-border-light/50 dark:border-border-dark/50"
                                                >
                                                    <span
                                                        class="material-symbols-outlined text-[22px]"
                                                        >category</span
                                                    >
                                                </div>
                                                <div>
                                                    <p
                                                        class="text-[11px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-tighter"
                                                    >
                                                        Categoria
                                                    </p>
                                                    <p
                                                        class="text-md font-black text-gray-900 dark:text-white"
                                                    >
                                                        {service.categoria}
                                                    </p>
                                                </div>
                                            </div>

                                            <!-- Duração -->
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="h-10 w-10 rounded-xl bg-white dark:bg-gray-800 flex items-center justify-center text-brand-orange shadow-sm border border-border-light/50 dark:border-border-dark/50"
                                                >
                                                    <span
                                                        class="material-symbols-outlined text-[22px]"
                                                        >schedule</span
                                                    >
                                                </div>
                                                <div>
                                                    <p
                                                        class="text-[11px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-tighter"
                                                    >
                                                        Duração
                                                    </p>
                                                    <p
                                                        class="text-md font-black text-gray-900 dark:text-white"
                                                    >
                                                        {service.duracao_padrao} min
                                                    </p>
                                                </div>
                                            </div>

                                            <!-- Preço -->
                                            <div class="flex items-start gap-4">
                                                <div
                                                    class="h-10 w-10 rounded-xl bg-white dark:bg-gray-800 flex items-center justify-center text-brand-orange shadow-sm border border-border-light/50 dark:border-border-dark/50"
                                                >
                                                    <span
                                                        class="material-symbols-outlined text-[22px]"
                                                        >payments</span
                                                    >
                                                </div>
                                                <div>
                                                    <p
                                                        class="text-[11px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-tighter"
                                                    >
                                                        Preço Base
                                                    </p>
                                                    <p
                                                        class="text-md font-black text-gray-900 dark:text-white"
                                                    >
                                                        R$ {formatPrice(service.preco)}
                                                    </p>
                                                </div>
                                            </div>
                                        </div>

                                        <div
                                            class="mt-8 pt-6 border-t border-border-light dark:border-border-dark"
                                        >
                                            <a
                                                href="/agendamento"
                                                class="w-full py-3.5 px-6 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-2xl font-black shadow-lg shadow-brand-orange/20 transition-all hover:scale-[1.02] active:scale-95 flex items-center justify-center gap-3"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-2xl"
                                                    >calendar_month</span
                                                >
                                                <span class="tracking-widest"
                                                    >Agendar</span
                                                >
                                            </a>
                                        </div>
                                    </div>

                                    <!-- Security/Policy Card -->
                                    <div
                                        class="bg-brand-orange/5 rounded-3xl p-6 border border-brand-orange/10"
                                    >
                                        <div class="flex items-start gap-4">
                                            <span
                                                class="material-symbols-outlined text-brand-orange"
                                                >verified_user</span
                                            >
                                            <p
                                                class="text-sm text-gray-600 dark:text-gray-400 leading-relaxed font-medium"
                                            >
                                                Cancelamento gratuito até 24h antes
                                                do horário agendado.
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    </main>
</div>
