<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { auth } from "$lib/stores/auth.svelte";

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

    // Filter mode: "date" (default calendar), "todos", "pendentes", "confirmados"
    type FilterMode = "date" | "todos" | "pendentes" | "confirmados";

    // State
    let providerId = $derived(auth.user?.id ?? "");
    let appointments = $state<Appointment[]>([]);
    let todayAppointments = $state<Appointment[]>([]);
    let filteredAppointments = $state<Appointment[]>([]);
    let activeFilter = $state<FilterMode>("date");
    let loading = $state(true);
    let actionLoading = $state<string | null>(null);
    let weekAppointmentDays = $state<Set<string>>(new Set());
    // Global stats: all appointments from today onward
    let allFutureAppointments = $state<Appointment[]>([]);

    // Calendar state
    let selectedDate = $state(new Date());
    let currentWeekStart = $state(getWeekStart(new Date()));

    const monthNames = [
        "Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
        "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"
    ];
    const dayNames = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"];

    // Reactive week days (7 days from Sunday)
    let currentWeekDays = $derived((() => {
        const days: Date[] = [];
        const weekStart = new Date(currentWeekStart);
        for (let i = 0; i < 7; i++) {
            const day = new Date(weekStart);
            day.setDate(weekStart.getDate() + i);
            days.push(day);
        }
        return days;
    })());

    // Reactive week display header
    let weekStart = $derived(currentWeekDays[0]);
    let weekEnd = $derived(currentWeekDays[6]);
    let weekHeaderLabel = $derived(currentWeekDays.length === 7
        ? `Semana de ${weekStart.getDate().toString().padStart(2, "0")}/${(weekStart.getMonth() + 1).toString().padStart(2, "0")} a ${weekEnd.getDate().toString().padStart(2, "0")}/${(weekEnd.getMonth() + 1).toString().padStart(2, "0")}`
        : "");

    // Reactive stats from all future appointments (today onward)
    let statsTotal = $derived(allFutureAppointments.filter(a => a.status !== "cancelado").length);
    let statsPendentes = $derived(allFutureAppointments.filter(a => a.status === "pendente").length);
    let statsConfirmados = $derived(allFutureAppointments.filter(a => a.status === "confirmado").length);

    // Reactive pending list for right panel
    let pendingAppointments = $derived(
        activeFilter === "date"
            ? todayAppointments.filter(a => a.status === "pendente")
            : filteredAppointments.filter(a => a.status === "pendente")
    );

    // The list displayed depends on activeFilter
    let displayedAppointments = $derived(
        activeFilter === "date" ? todayAppointments : filteredAppointments
    );

    // Reactive title for appointment list
    let listTitle = $derived(
        activeFilter === "pendentes"
            ? "Agendamentos Pendentes"
            : activeFilter === "confirmados"
            ? "Agendamentos Confirmados"
            : activeFilter === "todos"
            ? "Todos os Agendamentos"
            : isToday(selectedDate)
            ? "Atendimentos de Hoje"
            : `Atendimentos de ${selectedDate.getDate().toString().padStart(2, "0")} de ${monthNames[selectedDate.getMonth()]}`
    );

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
        activeFilter = "date";
        selectedDate = new Date(date);
        fetchAppointmentsForDate(date);
    }

    // Fetch all appointments from today onward, optionally filtered by status
    async function fetchFilteredAppointments(filter: FilterMode) {
        activeFilter = filter;
        loading = true;
        try {
            const today = formatDateString(new Date());
            const farFuture = "2099-12-31";
            const response = await fetchApi(
                `/api/v1/agendamentos/prestador/${providerId}/periodo?data_inicio=${today}&data_fim=${farFuture}&limit=100`
            );
            if (response.ok) {
                const data = await response.json();
                const raw: any[] = data.data ?? [];
                const all = raw.map((apt: any) => ({
                    ...apt,
                    status: normalizarStatus(apt.status)
                }));
                let result: Appointment[];
                if (filter === "pendentes") {
                    result = all.filter(a => a.status === "pendente");
                } else if (filter === "confirmados") {
                    result = all.filter(a => a.status === "confirmado");
                } else {
                    // "todos" — show all non-cancelled from today onward
                    result = all.filter(a => a.status !== "cancelado");
                }
                // Sort by date ascending
                filteredAppointments = result.sort((a, b) =>
                    new Date(a.data_inicio).getTime() - new Date(b.data_inicio).getTime()
                );
            } else {
                filteredAppointments = [];
            }
        } catch (error) {
            console.error("Erro ao buscar agendamentos filtrados:", error);
            filteredAppointments = [];
        } finally {
            loading = false;
        }
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

    // Re-fetch data based on current active filter
    async function refreshCurrentView() {
        if (activeFilter === "date") {
            await fetchAppointmentsForDate(selectedDate);
        } else {
            await fetchFilteredAppointments(activeFilter);
        }
        await fetchGlobalStats();
        await fetchWeekAppointments();
    }

    // Confirm appointment
    async function confirmar(id: string) {
        actionLoading = id;
        try {
            const response = await fetchApi(`/api/v1/agendamentos/${id}/confirmar`, {
                method: "PUT"
            });
            if (response.ok) {
                await refreshCurrentView();
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
                await refreshCurrentView();
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
            pendente: "border-l-amber-400",
            confirmado: "border-l-emerald-500",
            cancelado: "border-l-slate-300",
            concluido: "border-l-blue-400"
        };
        return map[status] ?? "border-l-slate-300";
    }

    function getStatusIconBgClass(status: string): string {
        const map: Record<string, string> = {
            pendente: "bg-amber-100 text-amber-600",
            confirmado: "bg-emerald-100 text-emerald-600",
            cancelado: "bg-slate-100 text-slate-500",
            concluido: "bg-blue-100 text-blue-500"
        };
        return map[status] ?? "bg-slate-100 text-slate-500";
    }

    function getStatusBadgeClass(status: string): string {
        const map: Record<string, string> = {
            pendente: "bg-amber-100 text-amber-700 border border-amber-200",
            confirmado: "bg-emerald-100 text-emerald-700 border border-emerald-200",
            cancelado: "bg-slate-100 text-slate-600 border border-slate-200",
            concluido: "bg-blue-100 text-blue-600 border border-blue-200"
        };
        return map[status] ?? "bg-slate-100 text-slate-500";
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

    // Fetch all future appointments for stats cards
    async function fetchGlobalStats() {
        if (!providerId) return;
        try {
            const today = formatDateString(new Date());
            const response = await fetchApi(
                `/api/v1/agendamentos/prestador/${providerId}/periodo?data_inicio=${today}&data_fim=2099-12-31&limit=100`
            );
            if (response.ok) {
                const data = await response.json();
                const raw: any[] = data.data ?? [];
                allFutureAppointments = raw.map((apt: any) => ({
                    ...apt,
                    status: normalizarStatus(apt.status)
                }));
            }
        } catch (e) {
            console.error("Erro ao buscar stats globais:", e);
        }
    }

    // Fetch week appointments to show dot indicators
    async function fetchWeekAppointments() {
        if (!providerId) return;
        const startStr = formatDateString(currentWeekDays[0]);
        const endStr = formatDateString(currentWeekDays[6]);
        try {
            const response = await fetchApi(
                `/api/v1/agendamentos/prestador/${providerId}/periodo?data_inicio=${startStr}&data_fim=${endStr}&limit=100`
            );
            if (response.ok) {
                const data = await response.json();
                const items: any[] = data.data ?? [];
                const days = new Set<string>();
                items.forEach((apt: any) => {
                    const status = apt.status;
                    if (status === 3 || status === 'cancelado') return;
                    const startDate = apt.data_inicio?.substring(0, 10);
                    if (startDate) days.add(startDate);
                });
                weekAppointmentDays = days;
            }
        } catch (e) {
            console.error("Erro ao buscar agendamentos da semana:", e);
        }
    }

    function dayHasAppointments(date: Date): boolean {
        return weekAppointmentDays.has(formatDateString(date));
    }

    // Re-fetch week indicators when the week changes
    $effect(() => {
        void currentWeekStart;
        fetchWeekAppointments();
    });

    onMount(async () => {
        await fetchAppointmentsForDate(selectedDate);
        await fetchGlobalStats();
        await fetchWeekAppointments();
    });
</script>

<div class="font-sans bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-slate-100 antialiased h-screen flex overflow-hidden">
    <Sidebar />
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar title="Agenda de Atendimentos" />
        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-6xl mx-auto space-y-6">

                <!-- Page Header -->
                <div>
                    <p class="text-sm text-slate-400 font-medium">
                        Gerencie seus atendimentos e confirme solicitações pendentes.
                    </p>
                </div>

                <!-- Stats Cards (3 cards) — clickable filters -->
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <!-- Todos -->
                    <button
                        onclick={() => fetchFilteredAppointments("todos")}
                        class="bg-white dark:bg-gray-900 p-5 rounded-2xl shadow-sm hover:shadow-md transition-all duration-200 border-2 flex items-center gap-4 cursor-pointer text-left
                            {activeFilter === 'todos' ? 'border-blue-400 ring-2 ring-blue-100 dark:ring-blue-900' : 'border-gray-100 dark:border-gray-800'}"
                    >
                        <div class="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-500">
                            <span class="material-symbols-outlined text-[22px]">calendar_today</span>
                        </div>
                        <div>
                            <p class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Todos</p>
                            <p class="text-3xl font-bold text-slate-900 dark:text-white mt-0.5" style="font-family: 'Cormorant', serif;">{statsTotal}</p>
                        </div>
                    </button>
                    <!-- Pendentes -->
                    <button
                        onclick={() => fetchFilteredAppointments("pendentes")}
                        class="bg-white dark:bg-gray-900 p-5 rounded-2xl shadow-sm hover:shadow-md transition-all duration-200 border-2 flex items-center gap-4 cursor-pointer text-left
                            {activeFilter === 'pendentes' ? 'border-amber-400 ring-2 ring-amber-100 dark:ring-amber-900' : 'border-gray-100 dark:border-gray-800'}"
                    >
                        <div class="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0 bg-amber-50 text-amber-500">
                            <span class="material-symbols-outlined text-[22px]">hourglass_top</span>
                        </div>
                        <div>
                            <p class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Pendentes</p>
                            <p class="text-3xl font-bold text-slate-900 dark:text-white mt-0.5" style="font-family: 'Cormorant', serif;">{statsPendentes}</p>
                        </div>
                    </button>
                    <!-- Confirmados -->
                    <button
                        onclick={() => fetchFilteredAppointments("confirmados")}
                        class="bg-white dark:bg-gray-900 p-5 rounded-2xl shadow-sm hover:shadow-md transition-all duration-200 border-2 flex items-center gap-4 cursor-pointer text-left
                            {activeFilter === 'confirmados' ? 'border-emerald-400 ring-2 ring-emerald-100 dark:ring-emerald-900' : 'border-gray-100 dark:border-gray-800'}"
                    >
                        <div class="w-12 h-12 rounded-xl flex items-center justify-center flex-shrink-0 bg-emerald-50 text-emerald-500">
                            <span class="material-symbols-outlined text-[22px]">check_circle</span>
                        </div>
                        <div>
                            <p class="text-xs font-semibold text-slate-400 uppercase tracking-wider">Confirmados</p>
                            <p class="text-3xl font-bold text-slate-900 dark:text-white mt-0.5" style="font-family: 'Cormorant', serif;">{statsConfirmados}</p>
                        </div>
                    </button>
                </div>

                <!-- Weekly Calendar -->
                <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 p-5">
                    <div class="flex items-center justify-between mb-4">
                        <h2
                            class="text-base font-bold text-slate-800 dark:text-slate-100"
                            style="font-family: 'Cormorant', serif; font-size: 1.1rem;"
                        >
                            {weekHeaderLabel}
                        </h2>
                        <div class="flex gap-1">
                            <button
                                onclick={previousWeek}
                                class="p-2 rounded-xl hover:bg-orange-50 hover:text-orange-600 text-slate-400 transition-all duration-200 cursor-pointer"
                                aria-label="Semana anterior"
                            >
                                <span class="material-symbols-outlined text-[20px]">chevron_left</span>
                            </button>
                            <button
                                onclick={nextWeek}
                                class="p-2 rounded-xl hover:bg-orange-50 hover:text-orange-600 text-slate-400 transition-all duration-200 cursor-pointer"
                                aria-label="Próxima semana"
                            >
                                <span class="material-symbols-outlined text-[20px]">chevron_right</span>
                            </button>
                        </div>
                    </div>
                    <div class="grid grid-cols-7 gap-2 text-center">
                        {#each currentWeekDays as day, index}
                            {@const isTodayDate = isToday(day)}
                            {@const isSelected = day.getDate() === selectedDate.getDate() && day.getMonth() === selectedDate.getMonth() && day.getFullYear() === selectedDate.getFullYear()}
                            <button
                                onclick={() => selectDate(day)}
                                class="p-2.5 rounded-xl transition-all duration-200 cursor-pointer group
                                    {isSelected
                                        ? 'bg-orange-500 text-white shadow-md shadow-orange-500/25'
                                        : isTodayDate
                                        ? 'bg-blue-50 text-blue-600 border border-blue-100'
                                        : 'hover:bg-orange-50 dark:hover:bg-orange-900/20 hover:text-orange-600'}"
                            >
                                <p class="text-xs font-semibold uppercase tracking-wider {isSelected ? 'text-white/80' : 'text-slate-400 dark:text-slate-500 group-hover:text-orange-400'}">
                                    {dayNames[index]}
                                </p>
                                <p class="text-sm font-bold mt-1
                                    {isSelected
                                        ? 'text-white'
                                        : isTodayDate
                                        ? 'text-blue-600'
                                        : 'text-slate-700 dark:text-slate-200 group-hover:text-orange-600'}">
                                    {day.getDate()}
                                </p>
                                {#if dayHasAppointments(day)}
                                    <div class="w-1.5 h-1.5 rounded-full {isSelected ? 'bg-white/80' : 'bg-orange-400'} mx-auto mt-1"></div>
                                {:else}
                                    <div class="w-1.5 h-1.5 mt-1"></div>
                                {/if}
                            </button>
                        {/each}
                    </div>
                </div>

                <!-- Main content grid: appointment list + right panel -->
                <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">

                    <!-- Appointment List (col-span-2) -->
                    <div class="lg:col-span-2 space-y-3">
                        <div class="flex items-center gap-2 mb-1">
                            <span class="material-symbols-outlined text-orange-500 text-[20px]">schedule</span>
                            <h3
                                class="text-lg font-bold text-slate-900 dark:text-white"
                                style="font-family: 'Cormorant', serif;"
                            >
                                {listTitle}
                            </h3>
                        </div>

                        {#if loading}
                            <div class="flex flex-col justify-center items-center py-20 gap-3">
                                <div class="w-12 h-12 rounded-2xl bg-orange-50 flex items-center justify-center">
                                    <span class="material-symbols-outlined animate-spin text-orange-500 text-2xl">sync</span>
                                </div>
                                <p class="text-sm text-slate-400">Buscando atendimentos...</p>
                            </div>
                        {:else if displayedAppointments.length === 0}
                            <!-- Styled empty state -->
                            <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-dashed border-slate-200 dark:border-slate-700 p-14 text-center">
                                <div class="w-16 h-16 bg-slate-50 dark:bg-slate-800 rounded-2xl flex items-center justify-center mx-auto mb-4">
                                    <span class="material-symbols-outlined text-slate-300 text-4xl">event_busy</span>
                                </div>
                                <p class="font-semibold text-slate-600 dark:text-slate-300">Nenhum atendimento</p>
                                <p class="text-sm text-slate-400 mt-1">
                                    {#if activeFilter === "date"}
                                        Nenhum atendimento para este dia.<br>Selecione outra data no calendário.
                                    {:else if activeFilter === "pendentes"}
                                        Nenhum agendamento pendente.
                                    {:else if activeFilter === "confirmados"}
                                        Nenhum agendamento confirmado.
                                    {:else}
                                        Nenhum agendamento futuro encontrado.
                                    {/if}
                                </p>
                            </div>
                        {:else}
                            {#each displayedAppointments as apt}
                                <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm hover:shadow-md transition-all duration-200 border-l-4 {getStatusBorderClass(apt.status)} border-y border-r border-gray-100 dark:border-gray-800 p-4 {apt.status === 'cancelado' ? 'opacity-60' : ''}">
                                    <div class="flex flex-col sm:flex-row justify-between items-start gap-4">

                                        <!-- Left: icon + info -->
                                        <div class="flex items-start gap-4 flex-1 min-w-0">
                                            <div class="{getStatusIconBgClass(apt.status)} p-2.5 rounded-xl flex-shrink-0">
                                                <span class="material-symbols-outlined text-[20px]">
                                                    {getStatusIcon(apt.status)}
                                                </span>
                                            </div>
                                            <div class="min-w-0 flex-1">
                                                <!-- Time range + service name + badge -->
                                                <div class="flex flex-wrap items-center gap-2">
                                                    <span class="text-base font-bold text-slate-900 dark:text-white">
                                                        {#if activeFilter !== "date"}
                                                            <span class="text-sm font-semibold text-orange-500 mr-1.5">
                                                                {new Date(apt.data_inicio).toLocaleDateString("pt-BR", { day: "2-digit", month: "2-digit" })}
                                                            </span>
                                                        {/if}
                                                        {formatarHora(apt.data_inicio)}
                                                        {#if apt.data_fim}
                                                            <span class="font-normal text-slate-300 dark:text-slate-600 mx-1">–</span>{formatarHora(apt.data_fim)}
                                                        {/if}
                                                    </span>
                                                    <span class="text-sm font-semibold text-slate-700 dark:text-slate-200 truncate">
                                                        {apt.servico?.nome ?? "Serviço"}
                                                    </span>
                                                    <span class="text-xs px-2.5 py-1 rounded-full font-semibold {getStatusBadgeClass(apt.status)}">
                                                        {getStatusLabel(apt.status)}
                                                    </span>
                                                </div>

                                                <!-- Client name + phone -->
                                                <p class="text-sm text-slate-500 dark:text-slate-400 mt-1.5">
                                                    <span class="font-semibold text-slate-700 dark:text-slate-200">
                                                        {apt.cliente?.nome ?? "Cliente"}
                                                    </span>
                                                    {#if apt.cliente?.telefone}
                                                        <span class="text-slate-400 ml-2">{apt.cliente.telefone}</span>
                                                    {/if}
                                                </p>

                                                <!-- Duration + price + notes -->
                                                <div class="flex flex-wrap items-center gap-3 mt-2 text-xs text-slate-400 dark:text-slate-500">
                                                    <span class="flex items-center gap-1">
                                                        <span class="material-symbols-outlined text-[13px]">timer</span>
                                                        {apt.servico?.duracao ?? 60} min
                                                    </span>
                                                    {#if apt.servico?.preco}
                                                        <span class="flex items-center gap-1 text-slate-500 dark:text-slate-400 font-medium">
                                                            <span class="material-symbols-outlined text-[13px]">payments</span>
                                                            {formatarPreco(apt.servico.preco)}
                                                        </span>
                                                    {/if}
                                                    {#if apt.notas}
                                                        <span class="flex items-center gap-1 italic">
                                                            <span class="material-symbols-outlined text-[13px]">note</span>
                                                            {apt.notas}
                                                        </span>
                                                    {/if}
                                                </div>
                                            </div>
                                        </div>

                                        <!-- Right: action buttons -->
                                        <div class="flex items-center gap-2 flex-shrink-0">
                                            {#if apt.status === "pendente"}
                                                <button
                                                    onclick={() => confirmar(apt.id)}
                                                    disabled={actionLoading === apt.id}
                                                    class="flex items-center gap-1.5 px-3.5 py-2 bg-emerald-500 hover:bg-emerald-600 disabled:opacity-60 text-white rounded-xl text-xs font-semibold shadow-sm shadow-emerald-500/20 transition-all duration-200 cursor-pointer"
                                                >
                                                    {#if actionLoading === apt.id}
                                                        <span class="material-symbols-outlined text-[14px] animate-spin">sync</span>
                                                    {:else}
                                                        <span class="material-symbols-outlined text-[14px]">check</span>
                                                    {/if}
                                                    Confirmar
                                                </button>
                                                <button
                                                    onclick={() => cancelar(apt.id)}
                                                    disabled={actionLoading === apt.id}
                                                    class="flex items-center gap-1.5 px-3.5 py-2 border border-red-200 text-red-500 hover:bg-red-50 disabled:opacity-60 rounded-xl text-xs font-semibold transition-all duration-200 cursor-pointer"
                                                >
                                                    {#if actionLoading === apt.id}
                                                        <span class="material-symbols-outlined text-[14px] animate-spin">sync</span>
                                                    {:else}
                                                        <span class="material-symbols-outlined text-[14px]">close</span>
                                                    {/if}
                                                    Rejeitar
                                                </button>
                                            {:else if apt.status === "confirmado"}
                                                <button
                                                    onclick={() => cancelar(apt.id)}
                                                    disabled={actionLoading === apt.id}
                                                    class="flex items-center gap-1.5 px-3.5 py-2 border border-slate-200 dark:border-slate-700 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 disabled:opacity-60 rounded-xl text-xs font-semibold transition-all duration-200 cursor-pointer"
                                                >
                                                    {#if actionLoading === apt.id}
                                                        <span class="material-symbols-outlined text-[14px] animate-spin">sync</span>
                                                    {:else}
                                                        <span class="material-symbols-outlined text-[14px]">cancel</span>
                                                    {/if}
                                                    Cancelar
                                                </button>
                                            {/if}
                                        </div>

                                    </div>
                                </div>
                            {/each}
                        {/if}
                    </div>

                    <!-- Right panel: pending confirmations -->
                    <div class="space-y-4">
                        <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 p-5">
                            <div class="flex items-center gap-2 mb-4">
                                <div class="w-8 h-8 bg-amber-50 rounded-lg flex items-center justify-center flex-shrink-0">
                                    <span class="material-symbols-outlined text-[16px] text-amber-500">hourglass_top</span>
                                </div>
                                <h3
                                    class="text-base font-bold text-slate-800 dark:text-slate-100"
                                    style="font-family: 'Cormorant', serif;"
                                >
                                    Pendentes de Confirmação
                                </h3>
                            </div>

                            {#if loading}
                                <div class="flex justify-center py-6">
                                    <span class="material-symbols-outlined animate-spin text-orange-500 text-2xl">sync</span>
                                </div>
                            {:else if pendingAppointments.length === 0}
                                <div class="text-center py-8">
                                    <div class="w-12 h-12 bg-slate-50 dark:bg-slate-800 rounded-xl flex items-center justify-center mx-auto mb-3">
                                        <span class="material-symbols-outlined text-slate-300 text-2xl">task_alt</span>
                                    </div>
                                    <p class="text-sm font-semibold text-slate-500 dark:text-slate-400">Tudo em dia!</p>
                                    <p class="text-xs text-slate-400 mt-0.5">Nenhuma solicitação pendente</p>
                                </div>
                            {:else}
                                <div class="space-y-2.5">
                                    {#each pendingAppointments as apt}
                                        <div class="flex items-center justify-between gap-3 p-3.5 rounded-xl bg-amber-50 dark:bg-amber-950/40 border border-amber-100 dark:border-amber-900">
                                            <div class="min-w-0 flex-1">
                                                <p class="text-sm font-semibold text-slate-800 dark:text-amber-100 truncate">
                                                    {apt.cliente?.nome ?? "Cliente"}
                                                </p>
                                                <p class="text-xs text-slate-500 dark:text-amber-200/70 mt-0.5 truncate">
                                                    {apt.servico?.nome ?? "Serviço"} · {formatarHora(apt.data_inicio)}
                                                </p>
                                            </div>
                                            <button
                                                onclick={() => confirmar(apt.id)}
                                                disabled={actionLoading === apt.id}
                                                class="flex-shrink-0 flex items-center gap-1 px-3 py-1.5 bg-emerald-500 hover:bg-emerald-600 disabled:opacity-60 text-white rounded-lg text-xs font-semibold shadow-sm transition-all duration-200 cursor-pointer"
                                            >
                                                {#if actionLoading === apt.id}
                                                    <span class="material-symbols-outlined text-[13px] animate-spin">sync</span>
                                                {:else}
                                                    <span class="material-symbols-outlined text-[13px]">check</span>
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
    ::-webkit-scrollbar {
        width: 6px;
        height: 6px;
    }
    ::-webkit-scrollbar-track {
        background: transparent;
    }
    ::-webkit-scrollbar-thumb {
        background: #e2e8f0;
        border-radius: 999px;
    }
    ::-webkit-scrollbar-thumb:hover {
        background: #cbd5e1;
    }
</style>
