<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import ProfileAgendaTab from "$lib/components/profile/ProfileAgendaTab.svelte";
    import AvailabilityModal from "$lib/components/profile/AvailabilityModal.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { auth } from "$lib/stores/auth.svelte";
    import type { ProviderData, AgendaInterval } from "$lib/types/profile";

    let providerId = $derived(auth.user?.id ?? "");

    let providerData = $state<ProviderData | null>(null);
    let loading = $state(true);

    // Availability modal state
    let showAvailabilityModal = $state(false);
    let initialDate = $state("");
    let initialIntervals = $state<{ start: string; end: string }[]>([]);

    // Alert modal state
    let showAlert = $state(false);
    let alertMessage = $state("");
    let alertType = $state<"success" | "error">("success");

    let availabilityMap = $derived(buildAvailabilityMap(providerData));
    let dayExceptions = $derived(buildDayExceptions(providerData));

    function buildAvailabilityMap(data: ProviderData | null): Record<string, AgendaInterval[]> {
        if (!data?.agenda) return {};
        const map: Record<string, AgendaInterval[]> = {};
        for (const day of data.agenda) {
            map[day.data] = day.intervalos ?? [];
        }
        return map;
    }

    function buildDayExceptions(data: ProviderData | null): Record<string, any> {
        if (!data?.agenda) return {};
        const map: Record<string, any> = {};
        for (const day of data.agenda) {
            map[day.data] = day;
        }
        return map;
    }

    async function fetchProvider() {
        if (!providerId) return;
        loading = true;
        try {
            const res = await fetchApi(`/api/v1/prestadores/${providerId}`);
            if (res.ok) {
                providerData = await res.json();
            }
        } finally {
            loading = false;
        }
    }

    function handleEditAvailability(data: { day: number; month: number; year: number }) {
        const { day, month, year } = data;
        const dateKey = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
        initialDate = dateKey;
        const existing = availabilityMap[dateKey] ?? [];
        initialIntervals = existing.map(i => ({ start: i.hora_inicio, end: i.hora_fim }));
        showAvailabilityModal = true;
    }

    async function handleSaveAvailability() {
        await fetchProvider();
        alertMessage = "Disponibilidade salva com sucesso!";
        alertType = "success";
        showAlert = true;
        showAvailabilityModal = false;
    }

    onMount(() => {
        fetchProvider();
    });
</script>

<div class="h-screen flex overflow-hidden bg-slate-50 dark:bg-gray-950">
    <Sidebar />
    <div class="flex-1 flex flex-col overflow-hidden">
        <DashboardNavbar title="Agenda Diária" />
        <main class="flex-1 overflow-y-auto p-6 md:p-8">

            {#if loading}
                <!-- Loading state -->
                <div class="flex flex-col justify-center items-center h-64 gap-4">
                    <div class="w-14 h-14 rounded-2xl bg-orange-50 flex items-center justify-center">
                        <span class="material-icons animate-spin text-orange-500 text-3xl">refresh</span>
                    </div>
                    <p class="text-sm text-slate-400 font-medium">Carregando agenda...</p>
                </div>

            {:else if providerData}
                <!-- Profile header card -->
                <div class="mb-8 rounded-2xl overflow-hidden shadow-sm bg-white dark:bg-gray-900 border border-gray-100 dark:border-gray-800">
                    <!-- Cover gradient -->
                    <div class="h-28 bg-gradient-to-r from-orange-500 to-pink-500 relative">
                        <div class="absolute inset-0 bg-black/10"></div>
                    </div>

                    <!-- Avatar + info row -->
                    <div class="px-6 pb-5 relative">
                        <div class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
                            <div class="flex items-end gap-4">
                                <!-- Avatar overlapping cover -->
                                <div class="-mt-10 flex-shrink-0">
                                    {#if providerData.image_url}
                                        <img
                                            src={providerData.image_url}
                                            alt={providerData.nome}
                                            class="h-20 w-20 rounded-2xl object-cover ring-4 ring-white dark:ring-gray-900 shadow-xl"
                                        />
                                    {:else}
                                        <div class="h-20 w-20 rounded-2xl bg-gradient-to-br from-orange-400 to-pink-500 flex items-center justify-center ring-4 ring-white shadow-xl">
                                            <span class="text-white text-2xl font-bold">
                                                {(providerData.nome ?? "P").split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase()}
                                            </span>
                                        </div>
                                    {/if}
                                </div>
                                <div class="pb-1">
                                    <h1
                                        class="text-2xl font-bold text-slate-900 dark:text-white leading-tight"
                                        style="font-family: 'Cormorant', serif;"
                                    >
                                        {providerData.nome ?? "Profissional"}
                                    </h1>
                                    <p class="text-sm text-orange-500 font-medium mt-0.5">
                                        Profissional
                                    </p>
                                </div>
                            </div>

                            <!-- Quick stats pills -->
                            <div class="flex items-center gap-3 pb-1 flex-wrap">
                                <div class="flex items-center gap-1.5 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-full px-3 py-1.5">
                                    <span class="material-symbols-outlined text-[16px] text-orange-500">calendar_month</span>
                                    <span class="text-xs font-semibold text-slate-600 dark:text-slate-300">
                                        {providerData.agenda?.length ?? 0} dias configurados
                                    </span>
                                </div>
                                {#if providerData.catalogo?.length}
                                    <div class="flex items-center gap-1.5 bg-slate-50 dark:bg-slate-800 border border-slate-100 dark:border-slate-700 rounded-full px-3 py-1.5">
                                        <span class="material-symbols-outlined text-[16px] text-pink-500">spa</span>
                                        <span class="text-xs font-semibold text-slate-600 dark:text-slate-300">
                                            {providerData.catalogo.length} {providerData.catalogo.length === 1 ? 'serviço' : 'serviços'}
                                        </span>
                                    </div>
                                {/if}
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Agenda content -->
                <div class="bg-white dark:bg-gray-900 rounded-2xl shadow-sm border border-gray-100 dark:border-gray-800 p-6">
                    <div class="flex items-center gap-2 mb-6">
                        <div class="w-8 h-8 bg-orange-50 rounded-lg flex items-center justify-center">
                            <span class="material-symbols-outlined text-[18px] text-orange-500">edit_calendar</span>
                        </div>
                        <div>
                            <h2
                                class="text-lg font-bold text-slate-900 dark:text-white leading-tight"
                                style="font-family: 'Cormorant', serif;"
                            >
                                Gestão de Disponibilidade
                            </h2>
                            <p class="text-xs text-slate-400">Clique em um dia para editar seus horários</p>
                        </div>
                    </div>

                    <ProfileAgendaTab
                        availabilityMap={availabilityMap}
                        dayExceptions={dayExceptions}
                        oneditAvailability={handleEditAvailability}
                    />
                </div>

            {:else}
                <!-- Empty/error state -->
                <div class="flex flex-col items-center justify-center h-64 gap-4">
                    <div class="w-16 h-16 rounded-2xl bg-slate-100 dark:bg-slate-800 flex items-center justify-center">
                        <span class="material-symbols-outlined text-slate-300 text-3xl">calendar_off</span>
                    </div>
                    <div class="text-center">
                        <p class="font-semibold text-slate-600 dark:text-slate-300">Agenda indisponível</p>
                        <p class="text-sm text-slate-400 mt-1">Não foi possível carregar sua agenda.</p>
                    </div>
                </div>
            {/if}
        </main>
    </div>
</div>

<AvailabilityModal
    bind:show={showAvailabilityModal}
    {providerId}
    {initialDate}
    {initialIntervals}
    onsuccess={handleSaveAvailability}
/>

<AlertModal
    bind:show={showAlert}
    message={alertMessage}
    type={alertType}
    oncancel={() => (showAlert = false)}
/>
