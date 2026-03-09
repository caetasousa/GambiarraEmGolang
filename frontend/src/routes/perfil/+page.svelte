<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import AvailabilityModal from "$lib/components/profile/AvailabilityModal.svelte";
    import ServiceSelectionModal from "$lib/components/profile/ServiceSelectionModal.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { user } from "$lib/stores/auth";
    
    // Import shared types
    import type { 
        Service, 
        AgendaInterval, 
        AgendaDay, 
        ProviderProfile, 
        ProviderData 
    } from '$lib/types/profile';

    // Import Refactored Components
    import ProfileAgendaTab from "$lib/components/profile/ProfileAgendaTab.svelte";

    // --- Calendar Logic (for modal only, main logic moved to component) ---
    const today = new Date();
    let currentMonth = today.getMonth();
    let currentYear = today.getFullYear();

    // Mock Data for "Exceptions" (Days with specific status)
    let dayExceptions: Record<
        string,
        { status: "off" | "partial"; label?: string }
    > = {
        "2025-07-28": { status: "off", label: "Férias" }, // Example future date
    };





    // --- Provider Data ---
    $: providerId = $user?.id ?? "";
    let availabilityMap: Record<string, AgendaInterval[]> = {};
    let currentProviderData: ProviderData | null = null;

    async function fetchProviderData() {
        try {
            const response = await fetchApi(`/api/v1/prestadores/${providerId}`);
            if (response.ok) {
                const data: ProviderData = await response.json();
                currentProviderData = data;

                // Populate Availability (Agenda)
                const agendaData = data.agenda;
                if (agendaData && Array.isArray(agendaData)) {
                    availabilityMap = agendaData.reduce(
                        (acc: Record<string, AgendaInterval[]>, item: AgendaDay) => {
                            const dateKey = item.data;
                            const intervalos = item.intervalos || [];

                            acc[dateKey] = intervalos.map((inv: AgendaInterval) => ({
                                id: inv.id,
                                hora_inicio: inv.hora_inicio,
                                hora_fim: inv.hora_fim,
                            }));
                            return acc;
                        },
                        {},
                    );
                }
            } else {
                console.error("Failed to fetch provider data");
            }
        } catch (error) {
            console.error("Error fetching provider data:", error);
        }
    }

    onMount(() => {
        fetchProviderData();
    });

    // --- Availability Modal ---
    let showAvailabilityModal = false;
    let selectedDateForModal = "";
    let selectedIntervalsForModal: { start: string; end: string }[] = [];

    function openAvailabilityModal(day?: number, month?: number, year?: number) {
        const targetDate = new Date();
        let targetDay = day;

        if (typeof targetDay !== "number") {
             selectedDateForModal = "";
             selectedIntervalsForModal = [];
             showAvailabilityModal = true;
             return;
        }

        const m = typeof month === 'number' ? month : currentMonth;
        const y = typeof year === 'number' ? year : currentYear;

        const dateObj = new Date(y, m, targetDay);
        const today = new Date();
        today.setHours(0, 0, 0, 0);

        if (dateObj < today) return;

        const paddedMonth = String(m + 1).padStart(2, "0");
        const paddedDay = String(targetDay).padStart(2, "0");
        selectedDateForModal = `${y}-${paddedMonth}-${paddedDay}`;

        if (availabilityMap[selectedDateForModal]) {
            const rawIntervals = availabilityMap[selectedDateForModal];
            selectedIntervalsForModal = rawIntervals.map((i: AgendaInterval) => ({
                start: i.hora_inicio,
                end: i.hora_fim,
            }));
        } else {
            selectedIntervalsForModal = [];
        }
        showAvailabilityModal = true;
    }
    
    function handleEditAvailability(event: CustomEvent<{ day: number; month: number; year: number }>) {
        openAvailabilityModal(event.detail.day, event.detail.month, event.detail.year);
    }

</script>

<div
    class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <Sidebar />

    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar />

        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-6xl mx-auto space-y-6">
                <!-- Agenda Content -->
                <ProfileAgendaTab 
                    {availabilityMap}
                    {dayExceptions}
                    on:editAvailability={handleEditAvailability}
                />
            </div>
        </div>
        <!-- Availability Modal -->
        <AvailabilityModal
            bind:show={showAvailabilityModal}
            {providerId}
            initialDate={selectedDateForModal}
            initialIntervals={selectedIntervalsForModal}
            on:success={fetchProviderData}
        />


    </main>
</div>
