<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { auth } from "$lib/stores/auth.svelte";

    interface ServicoInfo {
        id: string;
        nome: string;
        duracao: number;
        preco: number;
        categoria: string;
    }

    interface PrestadorInfo {
        id: string;
        nome: string;
        telefone: string;
    }

    interface Agendamento {
        id: string;
        prestador: PrestadorInfo;
        servico: ServicoInfo;
        data_inicio: string;
        data_fim: string;
        status: "pendente" | "confirmado" | "cancelado" | "concluido";
        notas?: string;
    }

    let clienteId = $derived(auth.user?.id ?? "");

    let agendamentos = $state<Agendamento[]>([]);
    let loading = $state(true);
    let erro = $state("");

    // Filtro por status
    let filtroStatus = $state("todos");
    let agendamentosFiltrados = $derived(
        filtroStatus === "todos"
            ? agendamentos
            : agendamentos.filter(a => a.status === filtroStatus)
    );

    function formatarData(isoStr: string): string {
        if (!isoStr) return "";
        const d = new Date(isoStr);
        return d.toLocaleDateString("pt-BR", {
            day: "2-digit",
            month: "long",
            year: "numeric",
        });
    }

    function formatarHora(isoStr: string): string {
        if (!isoStr) return "";
        const d = new Date(isoStr);
        return d.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit" });
    }

    function formatarPreco(centavos: number): string {
        return (centavos / 100).toLocaleString("pt-BR", {
            style: "currency",
            currency: "BRL",
        });
    }

    const STATUS_MAP: Record<number, string> = { 1: "pendente", 2: "confirmado", 3: "cancelado", 4: "concluido" };
    function normalizarStatus(s: any): string {
        if (typeof s === "number") return STATUS_MAP[s] ?? "pendente";
        return s ?? "pendente";
    }

    function getStatusLabel(status: string): string {
        const labels: Record<string, string> = {
            pendente: "Aguardando Confirmação",
            confirmado: "Confirmado",
            cancelado: "Cancelado",
            concluido: "Concluído",
        };
        return labels[status] || status;
    }

    function getStatusClasses(status: string) {
        const map: Record<string, { border: string; badge: string; icon: string }> = {
            pendente: {
                border: "border-l-yellow-400",
                badge: "bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400",
                icon: "hourglass_top",
            },
            confirmado: {
                border: "border-l-green-500",
                badge: "bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400",
                icon: "check_circle",
            },
            cancelado: {
                border: "border-l-gray-400",
                badge: "bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300",
                icon: "cancel",
            },
            concluido: {
                border: "border-l-blue-400",
                badge: "bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400",
                icon: "done_all",
            },
        };
        return map[status] || map.pendente;
    }

    async function carregarAgendamentos() {
        loading = true;
        erro = "";
        try {
            // Busca agendamentos a partir de 1 ano atrás para listar histórico
            const umAnoAtras = new Date();
            umAnoAtras.setFullYear(umAnoAtras.getFullYear() - 1);
            const dataStr = umAnoAtras.toISOString().split("T")[0];

            const res = await fetchApi(`/api/v1/agendamentos/cliente/${clienteId}?data=${dataStr}`);
            if (res.ok) {
                const json = await res.json();
                const lista = json.data ?? json ?? [];
                // Normalizar status numérico para string e ordenar do mais recente
                agendamentos = lista
                    .map((a: any) => ({ ...a, status: normalizarStatus(a.status) }))
                    .sort((a: Agendamento, b: Agendamento) =>
                        new Date(b.data_inicio).getTime() - new Date(a.data_inicio).getTime()
                    );
            } else {
                erro = "Não foi possível carregar seus agendamentos.";
            }
        } catch (e) {
            erro = "Erro ao conectar com o servidor.";
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        carregarAgendamentos();
    });
</script>

<div class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200">
    <Sidebar />
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar />
        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-4xl mx-auto space-y-6">

                <!-- Header -->
                <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
                    <div>
                        <nav class="flex text-sm text-gray-500 dark:text-gray-400 mb-2">
                            <span class="text-gray-900 dark:text-white font-medium">Meus Agendamentos</span>
                        </nav>
                        <h1 class="text-2xl font-bold text-gray-900 dark:text-white">Meus Agendamentos</h1>
                        <p class="text-gray-500 dark:text-gray-400 mt-1">Acompanhe seus agendamentos e histórico de serviços.</p>
                    </div>
                    <a
                        href="/clientes/agendamento"
                        class="flex items-center px-4 py-2 bg-brand-orange hover:bg-brand-orange/90 text-white rounded-lg text-sm font-medium shadow-md transition-colors"
                    >
                        <span class="material-symbols-outlined text-[20px] mr-2">add</span>
                        Novo Agendamento
                    </a>
                </div>

                <!-- Filtros -->
                <div class="flex flex-wrap gap-2">
                    {#each [
                        { value: "todos", label: "Todos" },
                        { value: "pendente", label: "Pendentes" },
                        { value: "confirmado", label: "Confirmados" },
                        { value: "concluido", label: "Concluídos" },
                        { value: "cancelado", label: "Cancelados" },
                    ] as filtro}
                        <button
                            onclick={() => (filtroStatus = filtro.value)}
                            class="px-3 py-1.5 rounded-full text-sm font-medium transition-colors {filtroStatus === filtro.value
                                ? 'bg-brand-orange text-white'
                                : 'bg-white dark:bg-[hsl(var(--bs-card))] text-gray-600 dark:text-gray-300 border border-border-light dark:border-border-dark hover:bg-gray-50 dark:hover:bg-gray-800'}"
                        >
                            {filtro.label}
                        </button>
                    {/each}
                </div>

                <!-- Lista -->
                {#if loading}
                    <div class="flex justify-center py-16">
                        <span class="material-symbols-outlined animate-spin text-brand-orange text-4xl">sync</span>
                    </div>
                {:else if erro}
                    <div class="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-6 text-center">
                        <span class="material-symbols-outlined text-red-500 text-4xl mb-2">error</span>
                        <p class="text-red-600 dark:text-red-400">{erro}</p>
                        <button
                            onclick={carregarAgendamentos}
                            class="mt-4 px-4 py-2 bg-brand-orange text-white rounded-lg text-sm hover:bg-brand-orange/90 transition-colors"
                        >
                            Tentar novamente
                        </button>
                    </div>
                {:else if agendamentosFiltrados.length === 0}
                    <div class="bg-[hsl(var(--bs-card))] rounded-lg border border-border-light dark:border-border-dark p-12 text-center">
                        <span class="material-symbols-outlined text-gray-300 dark:text-gray-600 text-6xl mb-3">event_busy</span>
                        <p class="text-gray-500 dark:text-gray-400 font-medium">
                            {filtroStatus === "todos" ? "Você ainda não tem agendamentos." : "Nenhum agendamento com esse status."}
                        </p>
                        {#if filtroStatus === "todos"}
                            <a
                                href="/clientes/agendamento"
                                class="inline-flex items-center mt-4 px-4 py-2 bg-brand-orange text-white rounded-lg text-sm font-medium hover:bg-brand-orange/90 transition-colors"
                            >
                                <span class="material-symbols-outlined text-[18px] mr-1">add</span>
                                Agendar agora
                            </a>
                        {/if}
                    </div>
                {:else}
                    <div class="space-y-4">
                        {#each agendamentosFiltrados as ag}
                            {@const sc = getStatusClasses(ag.status)}
                            <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border-l-4 {sc.border} border-y border-r border-border-light dark:border-border-dark p-5 transition-all hover:shadow-md {ag.status === 'cancelado' ? 'opacity-70' : ''}">
                                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                                    <!-- Info principal -->
                                    <div class="flex items-start gap-4">
                                        <div class="shrink-0 w-12 h-12 rounded-full bg-brand-orange/10 flex items-center justify-center">
                                            <span class="material-symbols-outlined text-brand-orange">content_cut</span>
                                        </div>
                                        <div>
                                            <div class="flex items-center gap-2 flex-wrap">
                                                <span class="font-semibold text-gray-900 dark:text-white">
                                                    {ag.servico?.nome ?? "Serviço"}
                                                </span>
                                                <span class="text-xs px-2 py-0.5 rounded-full font-medium {sc.badge}">
                                                    <span class="material-symbols-outlined text-[13px] align-middle mr-0.5">{sc.icon}</span>
                                                    {getStatusLabel(ag.status)}
                                                </span>
                                            </div>
                                            <p class="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                                                <span class="material-symbols-outlined text-[15px] align-middle mr-1">person</span>
                                                {ag.prestador?.nome ?? "Profissional"}
                                            </p>
                                            <div class="flex items-center gap-3 mt-1 text-xs text-gray-400 flex-wrap">
                                                <span class="flex items-center gap-1">
                                                    <span class="material-symbols-outlined text-[14px]">calendar_today</span>
                                                    {formatarData(ag.data_inicio)}
                                                </span>
                                                <span class="flex items-center gap-1">
                                                    <span class="material-symbols-outlined text-[14px]">schedule</span>
                                                    {formatarHora(ag.data_inicio)} – {formatarHora(ag.data_fim)}
                                                </span>
                                                <span class="flex items-center gap-1">
                                                    <span class="material-symbols-outlined text-[14px]">timer</span>
                                                    {ag.servico?.duracao ?? "–"} min
                                                </span>
                                            </div>
                                            {#if ag.notas}
                                                <p class="mt-1 text-xs text-gray-400 italic">"{ag.notas}"</p>
                                            {/if}
                                        </div>
                                    </div>

                                    <!-- Preço -->
                                    <div class="shrink-0 text-right">
                                        <p class="text-lg font-bold text-brand-orange">
                                            {formatarPreco(ag.servico?.preco ?? 0)}
                                        </p>
                                        <p class="text-xs text-gray-400">{ag.servico?.categoria ?? ""}</p>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>
                {/if}

            </div>
        </div>
    </main>
</div>
