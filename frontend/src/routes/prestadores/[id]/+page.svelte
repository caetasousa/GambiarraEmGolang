<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import { page } from "$app/stores";
    import { fetchApi } from "$lib/utils/api";

    interface CatalogItem {
        id: string;
        nome: string;
        duracao_padrao: number;
        preco: number;
        categoria: string;
        image_url: string;
    }

    interface Intervalo {
        id: string;
        hora_inicio: string;
        hora_fim: string;
    }

    interface AgendaItem {
        id: string;
        data: string;
        intervalos: Intervalo[];
    }

    interface Provider {
        id: string;
        nome: string;
        email: string;
        telefone: string;
        ativo: boolean;
        foto?: string;
        image_url?: string;
        biografia?: string;
        especialidade?: string;
        rating?: number;
        reviews?: number;
        catalogo: CatalogItem[];
        agenda: AgendaItem[];
    }

    let provider = $state<Provider | null>(null);
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        const id = $page.params.id;
        try {
            const res = await fetchApi(`/api/v1/prestadores/${id}`);
            if (res.ok) {
                provider = await res.json();

                if (provider) {
                    console.log("Provider Data:", provider);

                    provider.foto =
                        provider.image_url ||
                        provider.foto ||
                        `https://ui-avatars.com/api/?name=${encodeURIComponent(provider.nome)}&background=random`;
                    provider.biografia =
                        provider.biografia || "Sem biografia disponível.";
                    provider.especialidade =
                        provider.especialidade || "Profissional BellaVita";
                    provider.rating = 5.0;
                    provider.reviews = 0;

                    if (!provider.agenda) {
                        provider.agenda = [];
                    }
                }
            } else {
                error = "Prestador não encontrado.";
            }
        } catch (err) {
            console.error(err);
            error = "Erro ao carregar dados do prestador.";
        } finally {
            loading = false;
        }
    });

    function formatPrice(cents: number): string {
        return (cents / 100).toLocaleString("pt-BR", {
            style: "currency",
            currency: "BRL",
        });
    }

    function formatDate(dateString: string): string {
        const options: Intl.DateTimeFormatOptions = {
            weekday: "short",
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        };
        return new Date(dateString).toLocaleDateString("pt-BR", options);
    }
</script>

<div
    class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <Sidebar />

    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar title="Perfil do Profissional" />

        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-4xl mx-auto space-y-6">
                {#if loading}
                    <div class="flex items-center justify-center h-64">
                        <span
                            class="material-symbols-outlined animate-spin text-4xl text-primary"
                            >sync</span
                        >
                        <span class="ml-2 text-lg text-gray-500"
                            >Carregando perfil...</span
                        >
                    </div>
                {:else if error}
                    <div
                        class="bg-red-50 dark:bg-red-900/10 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 p-4 rounded-lg text-center"
                    >
                        <p class="font-bold">Erro</p>
                        <p>{error}</p>
                        <a
                            href="/prestadores"
                            class="text-primary hover:underline mt-2 inline-block"
                            >Voltar para a lista</a
                        >
                    </div>
                {:else if provider}
                    <!-- Breadcrumb -->
                    <nav
                        class="flex text-sm text-gray-500 dark:text-gray-400 mb-2"
                    >
                        <a href="/" class="hover:text-primary transition-colors"
                            >Início</a
                        >
                        <span class="mx-2">/</span>
                        <a
                            href="/prestadores"
                            class="hover:text-primary transition-colors"
                            >Profissionais</a
                        >
                        <span class="mx-2">/</span>
                        <span class="text-gray-900 dark:text-white font-medium"
                            >{provider.nome}</span
                        >
                    </nav>

                    <!-- Profile Card -->
                    <div
                        class="bg-[hsl(var(--bs-card))] rounded-xl shadow-sm border border-border-light dark:border-border-dark overflow-hidden"
                    >
                        <div
                            class="relative h-40 bg-gradient-to-r from-brand-orange to-amber-500"
                        >
                            <div
                                class="absolute -bottom-12 left-8 flex items-end"
                            >
                                <div class="relative">
                                    <img
                                        src={provider.foto}
                                        alt={provider.nome}
                                        class="h-32 w-32 rounded-full border-4 border-white dark:border-surface-dark object-cover bg-white shadow-lg"
                                    />
                                    <span
                                        class="absolute bottom-1 right-1 w-6 h-6 rounded-full border-4 border-white dark:border-surface-dark {provider.ativo
                                            ? 'bg-green-500'
                                            : 'bg-gray-400'}"
                                        title={provider.ativo
                                            ? "Ativo"
                                            : "Inativo"}
                                    ></span>
                                </div>
                                <div class="ml-6 mb-2 hidden sm:block">
                                    <h1
                                        class="text-2xl font-bold text-gray-900 dark:text-white"
                                    >
                                        {provider.nome}
                                    </h1>
                                    <p class="text-brand-orange font-medium">
                                        {provider.especialidade}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div class="pt-16 px-8 pb-8 space-y-8">
                            <!-- Mobile Header -->
                            <div class="sm:hidden text-center">
                                <h1
                                    class="text-2xl font-bold text-gray-900 dark:text-white"
                                >
                                    {provider.nome}
                                </h1>
                                <p class="text-brand-orange font-medium">
                                    {provider.especialidade}
                                </p>
                            </div>

                            <!-- Stats & Rating -->
                            <div
                                class="flex flex-wrap items-center justify-center sm:justify-start gap-6 py-4 border-y border-border-light dark:border-border-dark"
                            >
                                <div class="flex items-center text-amber-500">
                                    <span class="material-icons text-xl mr-1"
                                        >star</span
                                    >
                                    <span class="text-lg font-bold"
                                        >{provider.rating}</span
                                    >
                                    <span class="text-gray-400 ml-1"
                                        >({provider.reviews} avaliações)</span
                                    >
                                </div>
                                <div class="flex items-center text-gray-500">
                                    <span class="material-symbols-outlined mr-2"
                                        >verified</span
                                    >
                                    <span class="text-sm font-medium"
                                        >Verificado</span
                                    >
                                </div>
                                <div class="ml-auto w-full sm:w-auto">
                                    <button
                                        class="w-full sm:w-auto px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-md font-bold shadow-md transition-all flex items-center justify-center"
                                    >
                                        <span
                                            class="material-symbols-outlined mr-2"
                                            >calendar_month</span
                                        >
                                        Agendar com {provider.nome.split(
                                            " ",
                                        )[0]}
                                    </button>
                                </div>
                            </div>

                            <!-- Info Sections -->
                            <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                                <div class="md:col-span-2 space-y-8">
                                    <!-- Biography -->
                                    <section>
                                        <h3
                                            class="text-lg font-bold text-gray-900 dark:text-white flex items-center mb-4"
                                        >
                                            <span
                                                class="material-symbols-outlined mr-2 text-primary"
                                                >info</span
                                            >
                                            Sobre
                                        </h3>
                                        <p
                                            class="text-gray-600 dark:text-gray-300 leading-relaxed text-lg italic"
                                        >
                                            "{provider.biografia}"
                                        </p>
                                    </section>

                                    <!-- Services -->
                                    <section>
                                        <div
                                            class="flex items-center justify-between mb-4"
                                        >
                                            <h3
                                                class="text-lg font-bold text-gray-900 dark:text-white flex items-center"
                                            >
                                                <span
                                                    class="material-symbols-outlined mr-2 text-primary"
                                                    >content_cut</span
                                                >
                                                Serviços Atribuídos
                                            </h3>
                                        </div>
                                        <div class="grid grid-cols-1 gap-3">
                                            {#each provider.catalogo as servico}
                                                <div
                                                    class="flex items-center justify-between p-4 rounded-lg bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 border border-border-light dark:border-border-dark hover:border-brand-orange transition-colors"
                                                >
                                                    <div>
                                                        <p
                                                            class="font-bold text-gray-900 dark:text-white uppercase tracking-tight"
                                                        >
                                                            {servico.nome}
                                                        </p>
                                                        <p
                                                            class="text-xs text-gray-500"
                                                        >
                                                            {servico.duracao_padrao}
                                                            min • {servico.categoria}
                                                        </p>
                                                    </div>
                                                    <div class="text-right">
                                                        <p
                                                            class="text-lg font-bold text-brand-orange"
                                                        >
                                                            {formatPrice(
                                                                servico.preco,
                                                            )}
                                                        </p>
                                                    </div>
                                                </div>
                                            {:else}
                                                <p class="text-gray-500 italic">
                                                    Nenhum serviço atribuído.
                                                </p>
                                            {/each}
                                        </div>
                                    </section>
                                </div>

                                <!-- Sidebar/Contact Info -->
                                <div class="space-y-6">
                                    <section
                                        class="p-6 rounded-xl bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 border border-border-light dark:border-border-dark"
                                    >
                                        <h3
                                            class="text-sm font-bold text-gray-500 uppercase tracking-widest mb-4"
                                        >
                                            Contato
                                        </h3>
                                        <div class="space-y-4">
                                            <div
                                                class="flex items-center text-sm"
                                            >
                                                <span
                                                    class="material-symbols-outlined mr-3 text-gray-400"
                                                    >mail</span
                                                >
                                                <span
                                                    class="text-gray-700 dark:text-gray-300 break-all"
                                                    >{provider.email}</span
                                                >
                                            </div>
                                            <div
                                                class="flex items-center text-sm"
                                            >
                                                <span
                                                    class="material-symbols-outlined mr-3 text-gray-400"
                                                    >call</span
                                                >
                                                <span
                                                    class="text-gray-700 dark:text-gray-300"
                                                    >{provider.telefone}</span
                                                >
                                            </div>
                                        </div>
                                    </section>

                                    <section
                                        class="p-6 rounded-xl bg-brand-orange/5 border border-brand-orange/20"
                                    >
                                        <h3
                                            class="text-sm font-bold text-brand-orange uppercase tracking-widest mb-4"
                                        >
                                            Disponibilidade
                                        </h3>
                                        <div class="space-y-3 text-sm">
                                            {#if provider.agenda && provider.agenda.length > 0}
                                                {#each provider.agenda as dia}
                                                    <div
                                                        class="border-b border-gray-100 dark:border-gray-700/50 pb-2 last:border-0 last:pb-0"
                                                    >
                                                        <div
                                                            class="font-medium text-gray-900 dark:text-white mb-1"
                                                        >
                                                            {formatDate(
                                                                dia.data,
                                                            )}
                                                        </div>
                                                        <div class="space-y-1">
                                                            {#each dia.intervalos as intervalo}
                                                                <div
                                                                    class="flex justify-between text-gray-600 dark:text-gray-400 text-xs"
                                                                >
                                                                    <span
                                                                        >{intervalo.hora_inicio}
                                                                        - {intervalo.hora_fim}</span
                                                                    >
                                                                    <span
                                                                        class="text-green-600 dark:text-green-400 font-medium"
                                                                        >Livre</span
                                                                    >
                                                                </div>
                                                            {/each}
                                                        </div>
                                                    </div>
                                                {/each}
                                            {:else}
                                                <div
                                                    class="text-center py-2 text-gray-500"
                                                >
                                                    <p>Agenda não disponível</p>
                                                    <button
                                                        class="mt-2 text-primary text-xs hover:underline"
                                                        >Solicitar
                                                        disponibilidade</button
                                                    >
                                                </div>
                                            {/if}
                                        </div>
                                    </section>
                                </div>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>
        </div>
    </main>
</div>

<style>
    /* Custom scrollbar to match original style */
    ::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
    ::-webkit-scrollbar-track {
        background: transparent;
    }
    ::-webkit-scrollbar-thumb {
        background: #d1d5db;
        border-radius: 4px;
    }
    :global(.dark) ::-webkit-scrollbar-thumb {
        background: #4b5563;
    }
</style>
