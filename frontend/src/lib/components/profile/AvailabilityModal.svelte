<script lang="ts">
    import Modal from "$lib/components/Modal.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import { createEventDispatcher } from "svelte";
    import { fetchApi } from "$lib/utils/api";

    export let show = false;
    export let providerId: string;
    export let initialDate: string = "";
    export let initialIntervals: { start: string; end: string }[] = [];

    const dispatch = createEventDispatcher();

    // State
    let date = new Date().toISOString().split("T")[0]; // Default to today
    let intervals: { start: string; end: string }[] = [
        { start: "08:00", end: "12:00" },
    ];

    // Update date and intervals when modal opens with initial data
    $: if (show) {
        if (initialDate) {
            date = initialDate;
        }
        if (initialIntervals && initialIntervals.length > 0) {
            intervals = initialIntervals.map((i) => ({ ...i })); // Clone to avoid mutation
        } else {
            // Reset to default if no initial intervals provided (for new entries)
            intervals = [{ start: "08:00", end: "12:00" }];
        }
    }

    let loading = false;
    let error = "";
    let successMessage = "";
    let showDeleteConfirm = false;

    // Validation
    function validate() {
        if (!date) return "Selecione uma data.";

        const selectedDate = new Date(date);
        const today = new Date();
        today.setHours(0, 0, 0, 0); // Reset time for comparison
        // Fix timezone issue by treating input date as local
        const [year, month, day] = date.split("-").map(Number);
        const dateObj = new Date(year, month - 1, day);

        if (dateObj < today) {
            return "A data não pode ser no passado.";
        }

        if (intervals.length === 0) {
            // Allows empty intervals (implies deletion)
            return "";
        }

        for (const interval of intervals) {
            if (!interval.start || !interval.end) {
                return "Preencha todos os horários.";
            }
            if (interval.end <= interval.start) {
                return "O horário final deve ser maior que o inicial.";
            }
        }

        return "";
    }

    // Actions
    function addInterval() {
        let newStart = "";
        let newEnd = "";

        // Smart defaults based on requirement
        if (intervals.length === 0) {
            newStart = "08:00";
            newEnd = "12:00";
        } else if (intervals.length === 1) {
            newStart = "14:00";
            newEnd = "18:00";
        } else if (intervals.length === 2) {
            newStart = "19:00";
            newEnd = "22:00";
        }

        intervals = [...intervals, { start: newStart, end: newEnd }];
    }

    function removeInterval(index: number) {
        intervals = intervals.filter((_, i) => i !== index);
    }

    function close() {
        show = false;
        error = "";
        successMessage = "";
        // Reset to default on close? Optional.
    }

    async function handleSubmit() {
        error = "";
        successMessage = "";

        const validationError = validate();
        if (validationError) {
            error = validationError;
            return;
        }

        loading = true;

        try {
            const payload = {
                data: date,
                intervalos: intervals.map((i) => ({
                    hora_inicio: i.start,
                    hora_fim: i.end,
                })),
            };

            let response;

            if (intervals.length === 0) {
                // DELETE if no intervals
                response = await fetchApi(
                    `/api/v1/prestadores/${providerId}/agenda?data=${date}`,
                    { method: "DELETE" },
                );
            } else {
                // PUT if updating/creating
                response = await fetchApi(
                    `/api/v1/prestadores/${providerId}/agenda`,
                    {
                        method: "PUT",
                        body: JSON.stringify(payload),
                    },
                );
            }

            if (!response.ok) {
                const text = await response.text();
                throw new Error(text || "Erro ao salvar disponibilidade.");
            }

            successMessage = "Disponibilidade salva com sucesso!";
            setTimeout(() => {
                close();
                dispatch("success");
            }, 500);
        } catch (err: any) {
            console.error(err);
            error = err.message || "Erro de conexão.";
        } finally {
            loading = false;
        }
    }

    async function handleDelete() {
        showDeleteConfirm = true;
    }

    async function confirmDelete() {
        loading = true;
        error = "";
        successMessage = "";

        try {
            const response = await fetchApi(
                `/api/v1/prestadores/${providerId}/agenda?data=${date}`,
                { method: "DELETE" },
            );

            if (!response.ok) {
                const text = await response.text();
                throw new Error(text || "Erro ao excluir disponibilidade.");
            }

            successMessage = "Agenda excluída com sucesso!";
            setTimeout(() => {
                close();
                dispatch("success");
            }, 300);
        } catch (err: any) {
            console.error(err);
            error = err.message || "Erro ao excluir.";
        } finally {
            loading = false;
        }
    }
</script>

<AlertModal
    bind:show={showDeleteConfirm}
    title="Excluir Agenda"
    message="Tem certeza que deseja remover todos os horários deste dia? Essa ação não pode ser desfeita."
    type="error"
    confirmText="Sim, excluir"
    cancelText="Cancelar"
    showCancel={true}
    on:confirm={confirmDelete}
/>

<Modal
    {show}
    title="Adicionar Disponibilidade"
    maxWidth="max-w-lg"
    on:close={close}
>
    <div slot="body" class="space-y-6">
        {#if error}
            <div
                class="p-3 bg-red-50 text-red-700 rounded-md text-sm border border-red-200"
            >
                {error}
            </div>
        {/if}

        {#if successMessage}
            <div
                class="p-3 bg-green-50 text-green-700 rounded-md text-sm border border-green-200"
            >
                {successMessage}
            </div>
        {/if}

        <!-- Date Selection -->
        <div>
            <label
                for="date"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >
                Data
            </label>
            <input
                type="date"
                id="date"
                bind:value={date}
                min={new Date().toISOString().split("T")[0]}
                class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white shadow-sm focus:border-brand-orange focus:ring-brand-orange sm:text-sm p-2.5 border"
            />
        </div>

        <!-- Intervals -->
        <div>
            <div class="flex justify-between items-center mb-2">
                <span
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >
                    Intervalos de Horário
                </span>
                <button
                    type="button"
                    on:click={addInterval}
                    class="text-sm text-brand-orange hover:text-brand-orange-hover font-medium flex items-center"
                >
                    <span class="material-symbols-outlined text-base mr-1"
                        >add_circle</span
                    >
                    Adicionar
                </button>
            </div>

            <div class="space-y-3">
                {#each intervals as interval, i}
                    <div class="flex items-center space-x-2 animate-fade-in">
                        <div class="flex-1">
                            <input
                                type="time"
                                bind:value={interval.start}
                                class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white shadow-sm focus:border-brand-orange focus:ring-brand-orange sm:text-sm p-2 border"
                                aria-label="Hora de início"
                            />
                        </div>
                        <span class="text-gray-400">-</span>
                        <div class="flex-1">
                            <input
                                type="time"
                                bind:value={interval.end}
                                class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white shadow-sm focus:border-brand-orange focus:ring-brand-orange sm:text-sm p-2 border"
                                aria-label="Hora de fim"
                            />
                        </div>
                        <button
                            type="button"
                            on:click={() => removeInterval(i)}
                            class="p-2 text-gray-400 hover:text-red-500 rounded-full hover:bg-red-50 dark:hover:bg-red-900/10 transition-colors"
                            title="Remover intervalo"
                            disabled={intervals.length === 1 && false}
                        >
                            <span class="material-symbols-outlined">delete</span
                            >
                        </button>
                    </div>
                {/each}
                {#if intervals.length === 0}
                    <p class="text-sm text-gray-500 italic text-center py-2">
                        Nenhum horário definido.
                    </p>
                {/if}
            </div>
        </div>
    </div>

    <div slot="footer" class="flex justify-between w-full">
        <div>
            {#if initialIntervals.length > 0}
                <button
                    type="button"
                    on:click={handleDelete}
                    disabled={loading}
                    class="px-4 py-2 border border-red-300 rounded-md text-sm font-medium text-red-700 hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 transition-colors flex items-center"
                >
                    <span class="material-symbols-outlined text-sm mr-1"
                        >delete</span
                    >
                    Excluir Agenda
                </button>
            {/if}
        </div>
        <div class="flex space-x-3">
            <button
                type="button"
                on:click={close}
                class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-sm font-medium text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
                Cancelar
            </button>
            <button
                type="button"
                on:click={handleSubmit}
                disabled={loading}
                class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-brand-orange hover:bg-brand-orange-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-orange disabled:opacity-50 transition-colors flex items-center"
            >
                {#if loading}
                    <span
                        class="material-symbols-outlined animate-spin text-sm mr-2"
                        >sync</span
                    >
                    Salvando...
                {:else}
                    Salvar Disponibilidade
                {/if}
            </button>
        </div>
    </div>
</Modal>

<style>
    .animate-fade-in {
        animation: fadeIn 0.3s ease-in-out;
    }
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(-5px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
