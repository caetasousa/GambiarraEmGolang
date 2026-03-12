<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import ProfileAgendaTab from "$lib/components/profile/ProfileAgendaTab.svelte";
    import AvailabilityModal from "$lib/components/profile/AvailabilityModal.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { user } from "$lib/stores/auth";
    import type { ProviderData, AgendaInterval } from "$lib/types/profile";

    $: providerId = $user?.id ?? "";

    let providerData: ProviderData | null = null;
    let loading = true;

    // Availability modal state
    let showAvailabilityModal = false;
    let initialDate = "";
    let initialIntervals: { start: string; end: string }[] = [];

    // Alert modal state
    let showAlert = false;
    let alertMessage = "";
    let alertType: "success" | "error" = "success";

    $: availabilityMap = buildAvailabilityMap(providerData);
    $: dayExceptions = buildDayExceptions(providerData);

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

    function handleEditAvailability(event: CustomEvent<{ day: number; month: number; year: number }>) {
        const { day, month, year } = event.detail;
        const dateKey = `${year}-${String(month + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
        initialDate = dateKey;
        const existing = availabilityMap[dateKey] ?? [];
        initialIntervals = existing.map(i => ({ start: i.hora_inicio, end: i.hora_fim }));
        showAvailabilityModal = true;
    }

    async function handleSaveAvailability(event: CustomEvent<{ date: string; intervals: { start: string; end: string }[] }>) {
        const { date, intervals } = event.detail;
        try {
            const res = await fetchApi(`/api/v1/prestadores/${providerId}/agenda`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ data: date, intervalos: intervals.map(i => ({ hora_inicio: i.start, hora_fim: i.end })) }),
            });
            if (res.ok) {
                alertMessage = "Disponibilidade salva com sucesso!";
                alertType = "success";
                await fetchProvider();
            } else {
                alertMessage = "Erro ao salvar disponibilidade.";
                alertType = "error";
            }
        } catch {
            alertMessage = "Erro ao salvar disponibilidade.";
            alertType = "error";
        }
        showAlert = true;
        showAvailabilityModal = false;
    }

    onMount(() => {
        fetchProvider();
    });
</script>

<div class="h-screen flex overflow-hidden bg-[hsl(var(--bs-background))]">
    <Sidebar />
    <div class="flex-1 flex flex-col overflow-hidden">
        <DashboardNavbar title="Agenda Diária" />
        <main class="flex-1 overflow-y-auto p-6 md:p-8">
            {#if loading}
                <div class="flex justify-center items-center h-64">
                    <span class="material-icons animate-spin text-brand-orange text-4xl">refresh</span>
                </div>
            {:else if providerData}
                <ProfileAgendaTab
                    availabilityMap={availabilityMap}
                    dayExceptions={dayExceptions}
                    on:editAvailability={handleEditAvailability}
                />
            {:else}
                <p class="text-gray-500 dark:text-gray-400">Não foi possível carregar sua agenda.</p>
            {/if}
        </main>
    </div>
</div>

<AvailabilityModal
    bind:show={showAvailabilityModal}
    {providerId}
    {initialDate}
    {initialIntervals}
    on:save={handleSaveAvailability}
    on:close={() => (showAvailabilityModal = false)}
/>

<AlertModal
    bind:show={showAlert}
    message={alertMessage}
    type={alertType}
    on:close={() => (showAlert = false)}
/>
