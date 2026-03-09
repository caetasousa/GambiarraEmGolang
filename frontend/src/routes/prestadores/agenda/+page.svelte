<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { user } from "$lib/stores/auth";

    // Types
    interface Appointment {
        id: string;
        cliente: { id: string; nome: string; email: string; telefone: string };
        prestador: { id: string; nome: string };
        servico: { id: string; nome: string; duracao: number; preco: number; categoria: string };
        data_inicio: string;
        data_fim: string;
        status: "pendente" | "confirmado" | "cancelado" | "concluido";
        notas?: string;
    }

    // State
    $: providerId = $user?.id ?? "";
    let appointments: Appointment[] = [];
    let todayAppointments: Appointment[] = [];
    let loading = true;
    let actionLoading: string | null = null;

    // Calendar state
    let selectedDate = new Date();
    let currentWeekStart = getWeekStart(new Date());

    const monthNames = [
        "Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
        "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"
    ];
    const dayNames = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];

    // Reactive week display header
    $: weekStart = currentWeekDays[0];
    $: weekEnd = currentWeekDays[6];
    $: weekHeaderLabel = currentWeekDays.length === 7
        ? `Semana de ${weekStart.getDate().toString().padStart(2, "0")}/${(weekStart.getMonth() + 1).toString().padStart(2, "0")} a ${weekEnd.getDate().toString().padStart(2, "0")}/${(weekEnd.getMonth() + 1).toString().padStart(2, "0")}`
        : "";

    // Reactive stats from todayAppointments
    $: statsTotal = todayAppointments.length;
    $: statsPendentes = todayAppointments.filter(a => a.status === "pendente").length;
    $: statsConfirmados = todayAppointments.filter(a => a.status === "confirmado").length;

    // Reactive pending list for right panel
    $: pendingAppointments = todayAppointments.filter(a => a.status === "pendente");

    // Reactive title for appointment list
    $: listTitle = isToday(selectedDate)
        ? "Atendimentos de Hoje"
        : `Atendimentos de ${selectedDate.getDate().toString().padStart(2, "0")} de ${monthNames[selectedDate.getMonth()]}`;

    // Normalize status number → string
    const STATUS_MAP: Record<number, string> = {
        1: "pendente",
        2: "confirmado",
        3: "cancelado",
        4: "concluido"
    };

    function normalizarStatus(s: any): "pendente" | "confirmado" | "cancelado" | "concluido" {
        if (typeof s === "number") return (STATUS_MAP[s] ?? "pendente") as any;
        return s ?? "pendente";
    }

    // Format ISO to HH:MM in pt-BR
    function formatarHora(isoStr: string): string {
        if (!isoStr) return "";
        return new Date(isoStr).toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit" });
    }

    // Format price from cents to BRL currency string
    function formatarPreco(centavos: number): string {
        return (centavos / 100).toLocaleString("pt-BR", { style: "currency", currency: "BRL" });
    }

    // Format date to YYYY-MM-DD
    function formatDateString(date: Date): string {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        return `${year}-${month}-${day}`;
    }

    // Check if date is today
    function isToday(date: Date): boolean {
        const today = new Date();
        return (
            date.getDate() === today.getDate() &&
            date.getMonth() === today.getMonth() &&
            date.getFullYear() === today.getFullYear()
        );
    }

    // Get the start of the week (Sunday)
    function getWeekStart(date: Date): Date {
        const d = new Date(date);
        const day = d.getDay();
        d.setDate(d.getDate() - day);
        d.setHours(0, 0, 0, 0);
        return d;
    }

    // Reactive current week days (7 days from Sunday)
    $: currentWeekDays = (() => {
        const days: Date[] = [];
        const weekStart = new Date(currentWeekStart);
        for (let i = 0; i < 7; i++) {
            const day = new Date(weekStart);
            day.setDate(weekStart.getDate() + i);
            days.push(day);
        }
        return days;
    })();

    // Navigate to previous week
    function previousWeek() {
        const newStart = new Date(currentWeekStart);
        newStart.setDate(newStart.getDate() - 7);
        currentWeekStart = newStart;
    }

    // Navigate to next week
    function nextWeek() {
        const newStart = new Date(currentWeekStart);
        newStart.setDate(newStart.getDate() + 7);
        currentWeekStart = newStart;
    }

    // Select a date and fetch its appointments (past dates allowed)
    function selectDate(date: Date) {
        selectedDate = new Date(date);
        fetchAppointmentsForDate(date);
    }

    // Fetch appointments for a specific date
    async function fetchAppointmentsForDate(date: Date) {
        loading = true;
        try {
            const dateStr = formatDateString(date);
            const response = await fetchApi(
                `/api/v1/agendamentos/prestador/${providerId}?data=${dateStr}`
            );

            if (response.ok) {
                const data = await response.json();
                const raw: any[] = Array.isArray(data) ? data : (data.data ?? []);
                appointments = raw.map((apt: any) => ({
                    ...apt,
                    status: normalizarStatus(apt.status)
                }));
                todayAppointments = appointments.filter(apt =>
                    apt.data_inicio?.startsWith(dateStr)
                );
            } else {
                appointments = [];
                todayAppointments = [];
            }
        } catch (error) {
            console.error("Erro ao buscar agendamentos:", error);
            appointments = [];
            todayAppointments = [];
        } finally {
            loading = false;
        }
    }

    // Confirm appointment
    async function confirmar(id: string) {
        actionLoading = id;
        try {
            const response = await fetchApi(`/api/v1/agendamentos/${id}/confirmar`, {
                method: "PUT"
            });
            if (response.ok) {
                await fetchAppointmentsForDate(selectedDate);
            }
        } catch (error) {
            console.error("Erro ao confirmar agendamento:", error);
        } finally {
            actionLoading = null;
        }
    }

    // Cancel appointment
    async function cancelar(id: string) {
        actionLoading = id;
        try {
            const response = await fetchApi(`/api/v1/agendamentos/${id}/cancelar`, {
                method: "PUT"
            });
            if (response.ok) {
                await fetchAppointmentsForDate(selectedDate);
            }
        } catch (error) {
            console.error("Erro ao cancelar agendamento:", error);
        } finally {
            actionLoading = null;
        }
    }

    // Status helpers
    function getStatusLabel(status: string): string {
        const labels: Record<string, string> = {
            pendente: "Pendente",
            confirmado: "Confirmado",
            cancelado: "Cancelado",
            concluido: "Concluído"
        };
        return labels[status] ?? status;
    }

    function getStatusBorderClass(status: string): string {
        const map: Record<string, string> = {
            pendente: "border-l-yellow-400",
            confirmado: "border-l-green-500",
            cancelado: "border-l-gray-400",
            concluido: "border-l-blue-400"
        };
        return map[status] ?? "border-l-gray-300";
    }

    function getStatusIconBgClass(status: string): string {
        const map: Record<string, string> = {
            pendente: "bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400",
            confirmado: "bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400",
            cancelado: "bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400",
            concluido: "bg-blue-100 dark:bg-blue-900/30 text-blue-500 dark:text-blue-400"
        };
        return map[status] ?? "bg-gray-100 text-gray-500";
    }

    function getStatusBadgeClass(status: string): string {
        const map: Record<string, string> = {
            pendente: "bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-300",
            confirmado: "bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300",
            cancelado: "bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300",
            concluido: "bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-300"
        };
        return map[status] ?? "bg-gray-100 text-gray-500";
    }

    function getStatusIcon(status: string): string {
        const map: Record<string, string> = {
            pendente: "hourglass_top",
            confirmado: "check_circle",
            cancelado: "cancel",
            concluido: "done_all"
        };
        return map[status] ?? "help";
    }

    onMount(() => {
        fetchAppointmentsForDate(selectedDate);
    });
</script>

<div
    class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <Sidebar />
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar />
        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-6xl mx-auto space-y-6">

                <!-- Page Header -->
                <div>
                    <nav class="flex text-sm text-gray-500 dark:text-gray-400 mb-2">
                        <span class="text-gray-900 dark:text-white font-medium">Agenda</span>
                    </nav>
                    <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
                        Agenda de Atendimentos
                    </h1>
                    <p class="text-gray-500 dark:text-gray-400 mt-1">
                        Gerencie seus atendimentos e confirme solicitações pendentes.
                    </p>
                </div>

                <!-- Stats Cards (3 cards) -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <!-- Total do dia -->
                    <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
                        <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 mr-4">
                            <span class="material-symbols-outlined">calendar_today</span>
                        </div>
                        <div>
                            <p class="text-sm text-gray-500 dark:text-gray-400">Total do Dia</p>
                            <p class="text-2xl font-bold text-gray-900 dark:text-white">{statsTotal}</p>
                        </div>
                    </div>
                    <!-- Pendentes -->
                    <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
                        <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400 mr-4">
                            <span class="material-symbols-outlined">hourglass_top</span>
                        </div>
                        <div>
                            <p class="text-sm text-gray-500 dark:text-gray-400">Pendentes</p>
                            <p class="text-2xl font-bold text-gray-900 dark:text-white">{statsPendentes}</p>
                        </div>
                    </div>
                    <!-- Confirmados -->
                    <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
                        <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 mr-4">
                            <span class="material-symbols-outlined">check_circle</span>
                        </div>
                        <div>
                            <p class="text-sm text-gray-500 dark:text-gray-400">Confirmados</p>
                            <p class="text-2xl font-bold text-gray-900 dark:text-white">{statsConfirmados}</p>
                        </div>
                    </div>
                </div>

                <!-- Weekly Calendar -->
                <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-4">
                    <div class="flex items-center justify-between mb-4">
                        <h2 class="text-base font-semibold text-gray-900 dark:text-white">
                            {weekHeaderLabel}
                        </h2>
                        <div class="flex space-x-1">
                            <button
                                on:click={previousWeek}
                                class="p-1.5 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-300 transition-colors"
                                aria-label="Semana anterior"
                            >
                                <span class="material-symbols-outlined">chevron_left</span>
                            </button>
                            <button
                                on:click={nextWeek}
                                class="p-1.5 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-300 transition-colors"
                                aria-label="Próxima semana"
                            >
                                <span class="material-symbols-outlined">chevron_right</span>
                            </button>
                        </div>
                    </div>
                    <div class="grid grid-cols-7 gap-2 text-center">
                        {#each currentWeekDays as day, index}
                            {@const isTodayDate = isToday(day)}
                            {@const isSelected = day.getDate() === selectedDate.getDate() && day.getMonth() === selectedDate.getMonth() && day.getFullYear() === selectedDate.getFullYear()}
                            <button
                                on:click={() => selectDate(day)}
                                class="p-2 rounded-lg transition-all {isSelected
                                    ? 'bg-brand-orange text-white shadow-md'
                                    : isTodayDate
                                    ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400'
                                    : 'hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer'}"
                            >
                                <p class="text-xs uppercase {isSelected ? 'text-white/80' : 'text-gray-500 dark:text-gray-400'}">
                                    {dayNames[index]}
                                </p>
                                <p class="text-sm font-medium mt-1 {isSelected
                                    ? 'text-white font-bold'
                                    : isTodayDate
                                    ? 'text-blue-600 dark:text-blue-400 font-bold'
                                    : 'text-gray-900 dark:text-white'}">
                                    {day.getDate()}
                                </p>
                            </button>
                        {/each}
                    </div>
                </div>

                <!-- Main content grid: appointment list + right panel -->
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

                    <!-- Appointment List (col-span-2) -->
                    <div class="lg:col-span-2 space-y-4">
                        <h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
                            <span class="material-symbols-outlined text-brand-orange mr-2">schedule</span>
                            {listTitle}
                        </h3>

                        {#if loading}
                            <div class="flex justify-center items-center py-16">
                                <span class="material-symbols-outlined animate-spin text-brand-orange text-4xl">sync</span>
                            </div>
                        {:else if todayAppointments.length === 0}
                            <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-12 text-center">
                                <span class="material-symbols-outlined text-gray-300 dark:text-gray-600 text-6xl mb-3 block">event_busy</span>
                                <p class="text-gray-500 dark:text-gray-400 font-medium">Nenhum atendimento para este dia</p>
                                <p class="text-sm text-gray-400 dark:text-gray-500 mt-1">Selecione outro dia no calendário</p>
                            </div>
                        {:else}
                            {#each todayAppointments as apt}
                                <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border-l-4 {getStatusBorderClass(apt.status)} border-y border-r border-border-light dark:border-border-dark p-4 transition-all hover:shadow-md {apt.status === 'cancelado' ? 'opacity-70' : ''}">
                                    <div class="flex flex-col sm:flex-row justify-between items-start gap-4">

                                        <!-- Left: icon + info -->
                                        <div class="flex items-start space-x-4 flex-1 min-w-0">
                                            <div class="{getStatusIconBgClass(apt.status)} p-2 rounded-full flex-shrink-0">
                                                <span class="material-symbols-outlined text-[20px]">
                                                    {getStatusIcon(apt.status)}
                                                </span>
                                            </div>
                                            <div class="min-w-0 flex-1">
                                                <!-- Time range + service name + badge -->
                                                <div class="flex flex-wrap items-center gap-2">
                                                    <span class="text-base font-bold text-gray-900 dark:text-white">
                                                        {formatarHora(apt.data_inicio)}
                                                        {#if apt.data_fim}
                                                            <span class="font-normal text-gray-400 mx-1">–</span>{formatarHora(apt.data_fim)}
                                                        {/if}
                                                    </span>
                                                    <span class="text-sm font-medium text-gray-700 dark:text-gray-300 truncate">
                                                        {apt.servico?.nome ?? "Serviço"}
                                                    </span>
                                                    <span class="text-xs px-2 py-0.5 rounded-full font-medium {getStatusBadgeClass(apt.status)}">
                                                        {getStatusLabel(apt.status)}
                                                    </span>
                                                </div>

                                                <!-- Client name + phone -->
                                                <p class="text-sm text-gray-600 dark:text-gray-300 mt-1">
                                                    <span class="font-medium text-gray-800 dark:text-gray-200">
                                                        {apt.cliente?.nome ?? "Cliente"}
                                                    </span>
                                                    {#if apt.cliente?.telefone}
                                                        <span class="text-gray-400 ml-2">{apt.cliente.telefone}</span>
                                                    {/if}
                                                </p>

                                                <!-- Duration + notes -->
                                                <div class="flex flex-wrap items-center gap-3 mt-2 text-xs text-gray-400">
                                                    <span class="flex items-center gap-1">
                                                        <span class="material-symbols-outlined text-[14px]">timer</span>
                                                        {apt.servico?.duracao ?? 60} min
                                                    </span>
                                                    {#if apt.servico?.preco}
                                                        <span class="flex items-center gap-1">
                                                            <span class="material-symbols-outlined text-[14px]">payments</span>
                                                            {formatarPreco(apt.servico.preco)}
                                                        </span>
                                                    {/if}
                                                    {#if apt.notas}
                                                        <span class="flex items-center gap-1 text-gray-500 dark:text-gray-400 italic">
                                                            <span class="material-symbols-outlined text-[14px]">note</span>
                                                            {apt.notas}
                                                        </span>
                                                    {/if}
                                                </div>
                                            </div>
                                        </div>

                                        <!-- Right: action buttons -->
                                        <div class="flex items-center space-x-2 flex-shrink-0">
                                            {#if apt.status === "pendente"}
                                                <button
                                                    on:click={() => confirmar(apt.id)}
                                                    disabled={actionLoading === apt.id}
                                                    class="flex items-center gap-1 px-3 py-1.5 bg-green-600 hover:bg-green-700 disabled:opacity-60 text-white rounded text-sm font-medium shadow-sm transition-colors"
                                                >
                                                    {#if actionLoading === apt.id}
                                                        <span class="material-symbols-outlined text-[16px] animate-spin">sync</span>
                                                    {:else}
                                                        <span class="material-symbols-outlined text-[16px]">check</span>
                                                    {/if}
                                                    Confirmar
                                                </button>
                                                <button
                                                    on:click={() => cancelar(apt.id)}
                                                    disabled={actionLoading === apt.id}
                                                    class="flex items-center gap-1 px-3 py-1.5 border border-red-400 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 disabled:opacity-60 rounded text-sm font-medium transition-colors"
                                                >
                                                    {#if actionLoading === apt.id}
                                                        <span class="material-symbols-outlined text-[16px] animate-spin">sync</span>
                                                    {:else}
                                                        <span class="material-symbols-outlined text-[16px]">close</span>
                                                    {/if}
                                                    Rejeitar
                                                </button>
                                            {:else if apt.status === "confirmado"}
                                                <button
                                                    on:click={() => cancelar(apt.id)}
                                                    disabled={actionLoading === apt.id}
                                                    class="flex items-center gap-1 px-3 py-1.5 border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-60 rounded text-sm font-medium transition-colors"
                                                >
                                                    {#if actionLoading === apt.id}
                                                        <span class="material-symbols-outlined text-[16px] animate-spin">sync</span>
                                                    {:else}
                                                        <span class="material-symbols-outlined text-[16px]">cancel</span>
                                                    {/if}
                                                    Cancelar
                                                </button>
                                            {/if}
                                            <!-- cancelado / concluido: no buttons -->
                                        </div>

                                    </div>
                                </div>
                            {/each}
                        {/if}
                    </div>

                    <!-- Right panel: pending confirmations -->
                    <div class="space-y-4">
                        <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-5">
                            <h3 class="text-base font-semibold text-gray-900 dark:text-white mb-4 flex items-center gap-2">
                                <span class="material-symbols-outlined text-yellow-500">hourglass_top</span>
                                Pendentes de Confirmação
                            </h3>

                            {#if loading}
                                <div class="flex justify-center py-6">
                                    <span class="material-symbols-outlined animate-spin text-brand-orange text-3xl">sync</span>
                                </div>
                            {:else if pendingAppointments.length === 0}
                                <div class="text-center py-6">
                                    <span class="material-symbols-outlined text-gray-300 dark:text-gray-600 text-4xl block mb-2">task_alt</span>
                                    <p class="text-sm text-gray-500 dark:text-gray-400">Nenhuma solicitação pendente</p>
                                </div>
                            {:else}
                                <div class="space-y-3">
                                    {#each pendingAppointments as apt}
                                        <div class="flex items-center justify-between gap-3 p-3 rounded-lg bg-yellow-50 dark:bg-yellow-900/10 border border-yellow-100 dark:border-yellow-900/30">
                                            <div class="min-w-0 flex-1">
                                                <p class="text-sm font-medium text-gray-900 dark:text-white truncate">
                                                    {apt.cliente?.nome ?? "Cliente"}
                                                </p>
                                                <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5 truncate">
                                                    {apt.servico?.nome ?? "Serviço"} · {formatarHora(apt.data_inicio)}
                                                </p>
                                            </div>
                                            <button
                                                on:click={() => confirmar(apt.id)}
                                                disabled={actionLoading === apt.id}
                                                class="flex-shrink-0 flex items-center gap-1 px-2.5 py-1 bg-green-600 hover:bg-green-700 disabled:opacity-60 text-white rounded text-xs font-medium transition-colors"
                                            >
                                                {#if actionLoading === apt.id}
                                                    <span class="material-symbols-outlined text-[14px] animate-spin">sync</span>
                                                {:else}
                                                    <span class="material-symbols-outlined text-[14px]">check</span>
                                                {/if}
                                                Confirmar
                                            </button>
                                        </div>
                                    {/each}
                                </div>
                            {/if}
                        </div>
                    </div>

                </div>
            </div>
        </div>
    </main>
</div>

<style>
    /* Custom scrollbar */
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
