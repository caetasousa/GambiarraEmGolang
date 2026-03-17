<script lang="ts">
    import type { AgendaInterval } from "$lib/types/profile";

    let {
        availabilityMap = {},
        dayExceptions = {},
        oneditAvailability,
    }: {
        availabilityMap?: Record<string, AgendaInterval[]>;
        dayExceptions?: Record<string, unknown>;
        oneditAvailability?: (data: { day: number; month: number; year: number }) => void;
    } = $props();

    const monthNames = ["Janeiro","Fevereiro","Março","Abril","Maio","Junho","Julho","Agosto","Setembro","Outubro","Novembro","Dezembro"];

    const today = new Date();
    let currentMonth = $state(today.getMonth());
    let currentYear = $state(today.getFullYear());

    function getDaysInMonth(month: number, year: number) {
        return new Date(year, month + 1, 0).getDate();
    }

    function getFirstDayOfMonth(month: number, year: number) {
        return new Date(year, month, 1).getDay();
    }

    let daysInCurrentMonth = $derived(getDaysInMonth(currentMonth, currentYear));
    let firstDayOfWeek = $derived(getFirstDayOfMonth(currentMonth, currentYear));
    let calendarDays = $derived([
        ...Array(firstDayOfWeek).fill(null),
        ...Array.from({ length: daysInCurrentMonth }, (_, i) => i + 1),
    ]);

    function changeMonth(step: number) {
        currentMonth += step;
        if (currentMonth > 11) { currentMonth = 0; currentYear++; }
        else if (currentMonth < 0) { currentMonth = 11; currentYear--; }
    }

    function getDayStatus(day: number, map: Record<string, AgendaInterval[]>) {
        if (!day) return null;
        const date = new Date(currentYear, currentMonth, day);
        const todayDate = new Date();
        todayDate.setHours(0, 0, 0, 0);
        const isPast = date < todayDate;
        const dateStr = `${currentYear}-${String(currentMonth + 1).padStart(2, "0")}-${String(day).padStart(2, "0")}`;

        if (map[dateStr]) {
            const rawIntervals = map[dateStr];
            const uniqueIntervals = rawIntervals.filter(
                (v, i, a) => a.findIndex(t => t.hora_inicio === v.hora_inicio && t.hora_fim === v.hora_fim) === i,
            );
            let label = "Disponível";
            if (uniqueIntervals.length > 0) {
                const first = uniqueIntervals[0];
                label = `${first.hora_inicio} - ${first.hora_fim}`;
                if (uniqueIntervals.length > 1) label += "...";
            }
            return { status: "open", label, isReal: true, intervals: uniqueIntervals, isPast };
        }

        const dayOfWeek = date.getDay();
        if (dayOfWeek === 0) return { status: "off", label: "Folga", isReal: false, isPast };
        if (dayOfWeek === 6) return { status: "partial", label: "08:00 - 12:00", isReal: false, isPast };
        return { status: "open", label: "08:00 - 12:00", isReal: false, isPast };
    }
</script>

<div class="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-fade-in">
    <div class="lg:col-span-2 space-y-6">
        <div class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark p-6">
            <div class="flex items-center justify-between mb-6">
                <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center">
                    <span class="capitalize">{monthNames[currentMonth]} {currentYear}</span>
                </h2>
                <div class="flex space-x-2">
                    <button onclick={() => changeMonth(-1)} class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-400">
                        <span class="material-symbols-outlined">chevron_left</span>
                    </button>
                    <button onclick={() => changeMonth(1)} class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-400">
                        <span class="material-symbols-outlined">chevron_right</span>
                    </button>
                </div>
            </div>

            <div class="grid grid-cols-7 gap-px lg:gap-2 text-center">
                {#each ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"] as day}
                    <div class="text-xs font-semibold text-gray-400 uppercase tracking-wider py-2">{day}</div>
                {/each}

                {#each calendarDays as day}
                    {#if day}
                        {@const status = getDayStatus(day, availabilityMap)}
                        <button
                            onclick={() => oneditAvailability?.({ day, month: currentMonth, year: currentYear })}
                            class="aspect-square relative flex items-start justify-end p-2 rounded-lg border border-transparent transition-all group
                            {status?.isPast ? 'opacity-40 cursor-not-allowed grayscale' : 'hover:border-brand-orange'}
                            {status?.isReal ? 'bg-blue-100 dark:bg-blue-900/30 border-blue-200 dark:border-blue-500'
                                : status?.status === 'off' ? 'bg-red-50 dark:bg-red-900/10'
                                : status?.status === 'partial' ? 'bg-orange-50 dark:bg-orange-900/10'
                                : 'bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20'}"
                        >
                            <span class="z-10 text-sm font-medium {status?.status === 'off' ? 'text-red-500' : 'text-gray-700 dark:text-gray-300'}">{day}</span>
                            <div class="absolute inset-x-1 bottom-1 flex flex-col items-center">
                                {#if status?.status === "off"}
                                    <span class="text-[10px] font-bold text-red-500 uppercase tracking-tight truncate w-full">Folga</span>
                                {:else if status?.status === "partial"}
                                    <div class="w-1.5 h-1.5 rounded-full bg-orange-400 mb-0.5"></div>
                                    <span class="text-[10px] text-orange-600 dark:text-orange-400 truncate w-full hidden sm:block">{status?.label}</span>
                                {:else}
                                    <div class="w-1.5 h-1.5 rounded-full {status?.isReal ? 'bg-blue-500' : 'bg-green-400'}"></div>
                                    {#if status?.isReal && status?.intervals}
                                        <div class="flex flex-col w-full px-0.5 mt-0.5 space-y-0.5 overflow-hidden">
                                            {#each (status.intervals as AgendaInterval[]).slice(0, 3) as interval}
                                                <span class="text-[9px] leading-tight text-blue-600 dark:text-blue-400 truncate w-full font-medium block">
                                                    {interval.hora_inicio}-{interval.hora_fim}
                                                </span>
                                            {/each}
                                            {#if (status.intervals as AgendaInterval[]).length > 3}
                                                <span class="text-[8px] leading-none text-blue-500 text-center block">+{(status.intervals as AgendaInterval[]).length - 3}</span>
                                            {/if}
                                        </div>
                                    {:else if status?.isReal}
                                        <span class="text-[10px] text-blue-600 dark:text-blue-400 truncate w-full hidden sm:block font-medium">{status?.label}</span>
                                    {/if}
                                {/if}
                            </div>
                        </button>
                    {:else}
                        <div></div>
                    {/if}
                {/each}
            </div>
        </div>
    </div>

    <div class="space-y-6">
        <div class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark p-6">
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-lg font-bold text-gray-900 dark:text-white flex items-center">
                    <span class="material-symbols-outlined text-brand-orange mr-2">schedule</span>
                    Semana Padrão
                </h3>
                <button class="text-xs text-primary font-medium hover:underline">Editar</button>
            </div>
            <p class="text-xs text-gray-500 mb-4">Seus horários recorrentes.</p>
            <div class="space-y-3">
                <div class="flex justify-between items-center text-sm p-2 rounded hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                    <span class="text-gray-600 dark:text-gray-400">Seg - Sex</span>
                    <div class="flex flex-col text-right">
                        <span class="font-medium text-gray-900 dark:text-white">08:00 - 14:00</span>
                        <span class="font-medium text-gray-900 dark:text-white">14:00 - 18:00</span>
                    </div>
                </div>
                <div class="flex justify-between items-center text-sm p-2 rounded hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                    <span class="text-gray-600 dark:text-gray-400">Sábado</span>
                    <span class="font-medium text-gray-900 dark:text-white">08:00 - 12:00</span>
                </div>
                <div class="flex justify-between items-center text-sm p-2 rounded hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                    <span class="text-gray-600 dark:text-gray-400">Domingo</span>
                    <span class="font-bold text-red-400">Folga</span>
                </div>
            </div>
        </div>

        <div class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark p-6">
            <h3 class="text-lg font-bold text-gray-900 dark:text-white mb-4">Ações Rápidas</h3>
            <div class="grid grid-cols-2 gap-3">
                <button class="flex flex-col items-center justify-center p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-brand-orange hover:bg-orange-50 dark:hover:bg-orange-900/10 transition-all text-center">
                    <span class="material-symbols-outlined text-brand-orange mb-1">block</span>
                    <span class="text-xs font-semibold text-gray-700 dark:text-gray-300">Bloquear Dia</span>
                </button>
                <button
                    onclick={() => oneditAvailability?.({ day: 0, month: currentMonth, year: currentYear })}
                    class="flex flex-col items-center justify-center p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/10 transition-all text-center"
                >
                    <span class="material-symbols-outlined text-blue-500 mb-1">edit_calendar</span>
                    <span class="text-xs font-semibold text-gray-700 dark:text-gray-300">Ajustar Horário</span>
                </button>
                <button class="flex flex-col items-center justify-center p-3 rounded-lg border border-gray-200 dark:border-gray-700 hover:border-purple-500 hover:bg-purple-50 dark:hover:bg-purple-900/10 transition-all text-center">
                    <span class="material-symbols-outlined text-purple-500 mb-1">flight</span>
                    <span class="text-xs font-semibold text-gray-700 dark:text-gray-300">Marcar Férias</span>
                </button>
            </div>
        </div>
    </div>
</div>

<style>
    .animate-fade-in { animation: fadeIn 0.3s ease-in-out; }
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }
</style>
