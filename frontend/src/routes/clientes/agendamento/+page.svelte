<script lang="ts">
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import { onMount } from "svelte";
    import { page } from "$app/stores";
    import { fetchApi } from "$lib/utils/api";
    import { auth } from "$lib/stores/auth.svelte";

    // --- Configuração ---
    let CLIENTE_ID = $derived(auth.user?.id ?? "");

    // --- Estado ---
    let professionals = $state<any[]>([]);
    let services = $state<any[]>([]);
    let selectedProfessional = $state<any>(null);
    let selectedService = $state<any>(null);

    let preSelectedServiceId = $state<string | null>(null);
    let fetchedAppointments = $state<any[]>([]);
    let fetchedClientAppointments = $state<any[]>([]);

    // Estado do Formulário
    let loading = $state(true);
    let loadingProfessionals = $state(false);
    let submitting = $state(false);
    let error = $state("");
    let successMessage = $state("");
    let showSuccessModal = $state(false);

    // Estado do Modal de Erro
    let showErrorModal = $state(false);
    let errorTitle = $state("");
    let errorMessage = $state("");

    // Addons State
    let selectedAddonIds = $state<string[]>([]);

    let availableAddons = $derived((() => {
        if (!selectedProfessional || !selectedService) return [];
        const catalog =
            selectedProfessional.Catalogo ||
            selectedProfessional.catalogo ||
            [];
        const targetId = selectedService.ID || selectedService.id;

        return catalog.filter((s: any) => (s.ID || s.id) !== targetId);
    })());

    function toggleAddon(addon: any) {
        const id = addon.ID || addon.id;
        if (selectedAddonIds.includes(id)) {
            selectedAddonIds = selectedAddonIds.filter((i) => i !== id);
        } else {
            selectedAddonIds = [...selectedAddonIds, id];
        }
        generateSlots();
    }

    // Dynamic Slots
    interface Slot {
        time: string;
        disabled: boolean;
    }
    let morningSlots = $state<Slot[]>([]);
    let afternoonSlots = $state<Slot[]>([]);
    let eveningSlots = $state<Slot[]>([]);
    let selectedSlot = $state<string | null>(null);

    // --- Indicador de dias com agendamento no mês ---
    let monthAppointmentDays = $state<Set<string>>(new Set());

    async function fetchMonthAppointments() {
        if (!CLIENTE_ID) return;
        const startDate = `${currentYear}-${String(currentMonth + 1).padStart(2, "0")}-01`;
        const lastDay = new Date(currentYear, currentMonth + 1, 0).getDate();
        const endDate = `${currentYear}-${String(currentMonth + 1).padStart(2, "0")}-${String(lastDay).padStart(2, "0")}`;
        try {
            const res = await fetchApi(
                `/api/v1/agendamentos/cliente/${CLIENTE_ID}/periodo?data_inicio=${startDate}&data_fim=${endDate}&limit=100`,
            );
            if (res.ok) {
                const json = await res.json();
                const items: any[] = json.data ?? [];
                const days = new Set<string>();
                items.forEach((apt: any) => {
                    const s = apt.status;
                    if (s === 3 || s === "cancelado") return;
                    const dateStr = (apt.data_inicio || "").substring(0, 10);
                    if (dateStr) days.add(dateStr);
                });
                monthAppointmentDays = days;
            }
        } catch (e) {
            console.error("Erro ao buscar agendamentos do mês:", e);
        }
    }

    function dayHasAppointment(day: number): boolean {
        const month = String(currentMonth + 1).padStart(2, "0");
        const dayStr = String(day).padStart(2, "0");
        return monthAppointmentDays.has(`${currentYear}-${month}-${dayStr}`);
    }

    // --- Calendar Logic ---
    const monthNames = [
        "Janeiro",
        "Fevereiro",
        "Março",
        "Abril",
        "Maio",
        "Junho",
        "Julho",
        "Agosto",
        "Setembro",
        "Outubro",
        "Novembro",
        "Dezembro",
    ];
    const today = new Date();
    let currentMonth = $state(today.getMonth());
    let currentYear = $state(today.getFullYear());
    let selectedDay = $state<number | null>(null);

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
        if (currentMonth > 11) {
            currentMonth = 0;
            currentYear++;
        } else if (currentMonth < 0) {
            currentMonth = 11;
            currentYear--;
        }
        selectedDay = null;
        selectedSlot = null;
        morningSlots = []; afternoonSlots = []; eveningSlots = [];
        fetchMonthAppointments();
    }

    let dayRequestId = 0;

    async function selectDay(day: number) {
        if (!day) return;
        const date = new Date(currentYear, currentMonth, day);
        const now = new Date();
        now.setHours(0, 0, 0, 0);
        if (date < now) return;

        const thisRequest = ++dayRequestId;
        selectedDay = day;
        selectedSlot = null;
        morningSlots = []; afternoonSlots = []; eveningSlots = [];

        const month = String(currentMonth + 1).padStart(2, "0");
        const dayStr = String(day).padStart(2, "0");
        const dateIso = `${currentYear}-${month}-${dayStr}`;

        await fetchAvailableProfessionals(dateIso);
        if (thisRequest !== dayRequestId) return;

        if (selectedProfessional) {
            const profId = selectedProfessional.id || selectedProfessional.ID;
            await Promise.all([
                fetchProviderAppointments(profId, dateIso),
                fetchClientAppointments(dateIso),
            ]);
        }
        if (thisRequest !== dayRequestId) return;

        error = "";
        generateSlots();
    }

    // --- Slot Generation Logic ---
    function timeToMinutes(timeStr: string) {
        const [hh, mm] = timeStr.split(":").map(Number);
        return hh * 60 + mm;
    }

    function minutesToTime(minutes: number) {
        const hh = Math.floor(minutes / 60);
        const mm = minutes % 60;
        return `${String(hh).padStart(2, "0")}:${String(mm).padStart(2, "0")}`;
    }

    // Extrai minutos do dia de uma string ISO ou datetime sem conversão de timezone
    // Ex: "2026-03-17T08:00:00Z" → 480, "2026-03-17 14:30:00" → 870
    function extractMinutesFromDateStr(str: string): number {
        if (!str) return -1;
        let timePart = "";
        if (str.includes("T")) {
            timePart = str.split("T")[1] || "";
        } else if (str.includes(" ")) {
            timePart = str.split(" ")[1] || "";
        }
        if (!timePart) return -1;
        const [hh, mm] = timePart.split(":").map(Number);
        if (isNaN(hh) || isNaN(mm)) return -1;
        return hh * 60 + mm;
    }

    let totalDuration = $derived((() => {
        if (!selectedService) return 60;
        let duration =
            selectedService.DuracaoPadrao ||
            selectedService.duracao_padrao ||
            60;

        selectedAddonIds.forEach((id) => {
            const addon = availableAddons.find(
                (a: any) => (a.ID || a.id) === id,
            );
            if (addon) {
                duration += addon.DuracaoPadrao || addon.duracao_padrao || 0;
            }
        });

        return duration;
    })());

    function generateSlots() {
        const oldSlot = selectedSlot;

        morningSlots = [];
        afternoonSlots = [];
        eveningSlots = [];
        selectedSlot = null;

        if (!selectedProfessional || !selectedService || !selectedDay) return;

        const month = String(currentMonth + 1).padStart(2, "0");
        const dayStr = String(selectedDay).padStart(2, "0");
        const dateStr = `${currentYear}-${month}-${dayStr}`;

        const agenda =
            selectedProfessional.Agenda || selectedProfessional.agenda;
        const dayAgenda = agenda?.find(
            (a: any) => (a.Data || a.data) === dateStr,
        );

        if (!dayAgenda) return;

        const duracao = totalDuration;
        const intervalos = dayAgenda.Intervalos || dayAgenda.intervalos || [];

        const allSlots: string[] = [];

        intervalos.forEach((intervalo: any) => {
            const start = timeToMinutes(
                intervalo.HoraInicio || intervalo.horaInicio,
            );
            const end = timeToMinutes(intervalo.HoraFim || intervalo.horaFim);

            let current = start;
            while (current + duracao <= end) {
                allSlots.push(minutesToTime(current));
                current += 15;
            }
        });

        const getStartStr = (ap: any) =>
            ap.data_inicio || ap.DataInicio || ap.start_time || ap.StartTime || "";

        const dayProviderAppointments = fetchedAppointments.filter(
            (ap: any) => getStartStr(ap).startsWith(dateStr),
        );
        const dayClientAppointments = fetchedClientAppointments.filter(
            (ap: any) => getStartStr(ap).startsWith(dateStr),
        );

        const notCancelled = (ap: any) => {
            const s = ap.status ?? ap.Status;
            return s !== 3 && s !== 'cancelado' && s !== 'Cancelado';
        };
        const allBlockers = [
            ...dayProviderAppointments.filter(notCancelled),
            ...dayClientAppointments.filter(notCancelled),
        ];

        allSlots.forEach((slot) => {
            const slotStart = timeToMinutes(slot);
            const slotEnd = slotStart + duracao;

            const hasOverlap = allBlockers.some((ap: any) => {
                let apStartMinutes = -1;
                let apEndMinutes = -1;

                const startStr = getStartStr(ap);
                const endStr =
                    ap.data_fim || ap.DataFim || ap.end_time || ap.EndTime || "";

                // Extrair HH:MM diretamente da string ISO (sem conversão de timezone)
                apStartMinutes = extractMinutesFromDateStr(startStr);
                apEndMinutes = extractMinutesFromDateStr(endStr);

                // Fallback: campo isolado de hora
                if (apStartMinutes === -1) {
                    const isolatedTime =
                        ap.HoraInicio || ap.hora_inicio || ap.Hora || ap.hora;
                    if (isolatedTime) {
                        apStartMinutes = timeToMinutes(isolatedTime);
                    }
                }

                if (apStartMinutes === -1) return false;

                // Se não temos data_fim, calcular a partir da duração do serviço
                if (apEndMinutes === -1) {
                    const apDuracaoMinutes = Number(
                        ap.servico?.duracao ||
                            ap.servico?.DuracaoPadrao ||
                            ap.servico?.duracao_padrao ||
                            ap.Servico?.Duracao ||
                            ap.Servico?.DuracaoPadrao ||
                            ap.Duracao ||
                            ap.duracao ||
                            ap.duracao_minutos ||
                            60,
                    );
                    apEndMinutes = apStartMinutes + apDuracaoMinutes;
                }

                return slotStart < apEndMinutes && slotEnd > apStartMinutes;
            });

            const hour = parseInt(slot.split(":")[0]);
            if (hour < 12) {
                morningSlots.push({ time: slot, disabled: hasOverlap });
            } else if (hour < 18) {
                afternoonSlots.push({ time: slot, disabled: hasOverlap });
            } else {
                eveningSlots.push({ time: slot, disabled: hasOverlap });
            }
        });

        if (
            oldSlot &&
            [...morningSlots, ...afternoonSlots, ...eveningSlots].some(
                (s) => s.time === oldSlot && !s.disabled,
            )
        ) {
            selectedSlot = oldSlot;
        } else {
            const allAvailable = [
                ...morningSlots,
                ...afternoonSlots,
                ...eveningSlots,
            ].filter((s) => !s.disabled);

            if (allAvailable.length > 0) {
                selectedSlot = allAvailable[0].time;
            }
        }

        morningSlots = morningSlots;
        afternoonSlots = afternoonSlots;
        eveningSlots = eveningSlots;
    }

    // generateSlots() é chamado manualmente nos handlers: selectDay, selectProfessional, toggleAddon

    function isPastAndNotToday(day: number) {
        if (!day) return false;
        const date = new Date(currentYear, currentMonth, day);
        const now = new Date();
        now.setHours(0, 0, 0, 0);
        return date < now;
    }

    function dayHasAgenda(day: number): boolean {
        if (!selectedProfessional) return false;
        const agenda = selectedProfessional.Agenda || selectedProfessional.agenda;
        if (!agenda) return false;
        const month = String(currentMonth + 1).padStart(2, "0");
        const dayStr = String(day).padStart(2, "0");
        const dateStr = `${currentYear}-${month}-${dayStr}`;
        return agenda.some((a: any) => (a.Data || a.data) === dateStr);
    }

    async function fetchAvailableProfessionals(dateIso: string) {
        loadingProfessionals = true;
        try {
            const resp = await fetchApi(
                `/api/v1/prestadores/disponiveis?data=${dateIso}`,
            );
            if (resp.ok) {
                const data = await resp.json();
                let allFromDate = data.data || [];

                if (selectedService) {
                    const targetId = selectedService.id || selectedService.ID;
                    professionals = allFromDate.filter((p: any) => {
                        const catalogo = p.Catalogo || p.catalogo || p.catalog;
                        if (!catalogo || !Array.isArray(catalogo)) {
                            console.warn(
                                `Prestador ${p.id || p.ID} sem catálogo na resposta de disponíveis`,
                            );
                            return true;
                        }

                        return catalogo.some(
                            (s: any) =>
                                String(s.id || s.ID) === String(targetId),
                        );
                    });
                } else {
                    professionals = allFromDate;
                }

                if (professionals.length > 0) {
                    selectedProfessional = professionals[0];
                    error = "";
                } else {
                    selectedProfessional = null;
                }
            }
        } catch (e) {
            console.error("Erro ao carregar prestadores:", e);
        } finally {
            loadingProfessionals = false;
        }
    }

    async function fetchProviderAppointments(profId: string, dateStr: string) {
        try {
            const res = await fetchApi(
                `/api/v1/agendamentos/prestador/${profId}?data=${dateStr}`,
            );
            if (res.ok) {
                const json = await res.json();
                fetchedAppointments = json.data || [];
            } else {
                fetchedAppointments = [];
            }
        } catch (e) {
            console.error("Error fetching provider appointments:", e);
            fetchedAppointments = [];
        }
    }

    async function fetchClientAppointments(dateStr: string) {
        try {
            const res = await fetchApi(
                `/api/v1/agendamentos/cliente/${CLIENTE_ID}?data=${dateStr}`,
            );
            if (res.ok) {
                const json = await res.json();
                fetchedClientAppointments = json.data || [];
            } else {
                fetchedClientAppointments = [];
            }
        } catch (e) {
            console.error("Error fetching client appointments:", e);
            fetchedClientAppointments = [];
        }
    }

    async function handleProfessionalChange(prof: any) {
        selectedProfessional = prof;
        if (selectedDay) {
            const month = String(currentMonth + 1).padStart(2, "0");
            const dayStr = String(selectedDay).padStart(2, "0");
            const dateIso = `${currentYear}-${month}-${dayStr}`;
            const profId = prof.id || prof.ID;
            await fetchProviderAppointments(profId, dateIso);
            generateSlots();
        }
    }

    // --- API Interactions ---

    async function fetchData() {
        loading = true;
        try {
            preSelectedServiceId = $page.url.searchParams.get("serviceId");

            if (preSelectedServiceId) {
                const respService = await fetchApi(
                    `/api/v1/catalogos/${preSelectedServiceId}`,
                );
                if (respService.ok) {
                    const data = await respService.json();
                    selectedService = data;
                }
            } else {
                const respServicos = await fetchApi("/api/v1/catalogos?limit=5");
                if (respServicos.ok) {
                    const data = await respServicos.json();
                    services = data.data || [];
                    if (services.length > 0) {
                        selectedService = services[0];
                    }
                }
            }

            const respPrestadores = await fetchApi("/api/v1/prestadores");
            if (respPrestadores.ok) {
                const data = await respPrestadores.json();
                let allProfessionals = data.data || [];

                if (selectedService) {
                    const targetId = selectedService.id || selectedService.ID;
                    professionals = allProfessionals.filter((p: any) => {
                        const catalogo = p.Catalogo || p.catalogo || p.catalog;
                        if (!catalogo || !Array.isArray(catalogo)) {
                            return true;
                        }
                        return catalogo.some(
                            (s: any) =>
                                String(s.id || s.ID) === String(targetId),
                        );
                    });
                } else {
                    professionals = allProfessionals;
                }

                if (professionals.length > 0) {
                    selectedProfessional = professionals[0];
                }
            }
        } catch (e) {
            console.error("Erro ao carregar dados:", e);
            error = "Não foi possível carregar as informações.";
        } finally {
            loading = false;
        }
    }

    // Função para exibir erros amigáveis
    function showError(apiError: string) {
        let title = "Erro ao Agendar";
        let message = "Não foi possível realizar o agendamento.";

        try {
            const errorObj = JSON.parse(apiError);
            const errorText = errorObj.error || errorObj.message || apiError;

            if (
                errorText.toLowerCase().includes("ja existe um agendamento") ||
                errorText.toLowerCase().includes("já existe um agendamento")
            ) {
                title = "Agendamento Duplicado";
                message =
                    "Você já possui um agendamento para este serviço neste dia. Por favor, escolha outra data ou horário.";
            } else if (
                errorText.toLowerCase().includes("horário") ||
                errorText.toLowerCase().includes("horario") ||
                errorText.toLowerCase().includes("ocupado")
            ) {
                title = "Horário Indisponível";
                message =
                    "Este horário não está mais disponível. Por favor, selecione outro horário.";
            } else if (
                errorText.toLowerCase().includes("profissional") ||
                errorText.toLowerCase().includes("prestador")
            ) {
                title = "Profissional Indisponível";
                message =
                    "O profissional selecionado não está disponível neste momento. Por favor, escolha outro profissional.";
            } else {
                message =
                    "Ocorreu um erro ao processar seu agendamento. Por favor, tente novamente ou entre em contato conosco.";
            }
        } catch {
            if (
                apiError.toLowerCase().includes("ja existe") ||
                apiError.toLowerCase().includes("já existe")
            ) {
                title = "Agendamento Duplicado";
                message =
                    "Você já possui um agendamento para este serviço neste dia. Por favor, escolha outra data ou horário.";
            } else if (
                apiError.toLowerCase().includes("conexão") ||
                apiError.toLowerCase().includes("connection")
            ) {
                title = "Erro de Conexão";
                message =
                    "Não foi possível conectar ao servidor. Verifique sua conexão e tente novamente.";
            }
        }

        errorTitle = title;
        errorMessage = message;
        showErrorModal = true;
    }

    async function handleConfirmAppointment() {
        if (
            !selectedProfessional ||
            !selectedService ||
            !selectedDay ||
            !selectedSlot
        ) {
            showError(
                "Por favor, selecione todas as opções antes de confirmar o agendamento.",
            );
            return;
        }

        submitting = true;
        error = "";
        successMessage = "";

        const year = currentYear;
        const month = String(currentMonth + 1).padStart(2, "0");
        const day = String(selectedDay).padStart(2, "0");
        const dateTimeString = `${year}-${month}-${day}T${selectedSlot}:00Z`;

        const payload = {
            catalogo_id: selectedService.ID || selectedService.id,
            cliente_id: CLIENTE_ID,
            data_hora_inicio: dateTimeString,
            notas: "Agendado via Web",
            prestador_id: selectedProfessional.id || selectedProfessional.ID,
        };

        try {
            const response = await fetchApi("/api/v1/agendamentos", {
                method: "POST",
                body: JSON.stringify(payload),
            });

            if (response.ok) {
                showSuccessModal = true;

                const dateOnly = dateTimeString.split("T")[0];
                await fetchClientAppointments(dateOnly);
                if (selectedProfessional) {
                    const profId =
                        selectedProfessional.id || selectedProfessional.ID;
                    await fetchProviderAppointments(profId, dateOnly);
                }
                fetchMonthAppointments();
            } else {
                const text = await response.text();
                console.error("Erro no agendamento:", text);
                showError(text);
            }
        } catch (e) {
            console.error("Erro de conexão:", e);
            showError("Erro de conexão ao tentar agendar.");
        } finally {
            submitting = false;
        }
    }

    onMount(() => {
        fetchData();
        fetchMonthAppointments();
    });
</script>

<div
    class="flex h-screen bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark font-body overflow-hidden transition-colors duration-200"
>
    <!-- Sidebar -->
    <Sidebar />

    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar />

        <!-- Content -->
        <div class="flex-1 overflow-y-auto p-4 md:p-6 lg:p-8">
            <div class="max-w-6xl mx-auto space-y-6">
                <!-- Breadcrumbs & Title -->
                <div>
                    <nav
                        class="flex text-sm text-gray-500 dark:text-gray-400 mb-2"
                    >
                        <a
                            href="/"
                            class="hover:text-brand-orange transition-colors"
                            >Início</a
                        >
                        <span class="mx-2">/</span>
                        <span class="text-gray-900 dark:text-white font-medium"
                            >Agendar Serviço</span
                        >
                    </nav>
                    <h1
                        class="text-2xl font-bold text-gray-900 dark:text-white"
                    >
                        Selecione Data e Horário
                    </h1>
                    <p class="text-gray-500 dark:text-gray-400 mt-1">
                        Finalize o agendamento para o seu momento de cuidado.
                    </p>
                </div>

                {#if loading}
                    <div class="flex justify-center py-12">
                        <span
                            class="material-icons animate-spin text-brand-orange text-4xl"
                            >sync</span
                        >
                    </div>
                {:else}
                    <!-- Main Grid -->
                    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
                        <!-- Left Column (Service & Professional) -->
                        <div class="lg:col-span-1 space-y-6">
                            <!-- Service Card (Display Selected Service) -->
                            {#if selectedService}
                                <div
                                    class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark overflow-hidden"
                                >
                                    <div
                                        class="h-32 bg-gray-200 dark:bg-gray-700 relative group"
                                    >
                                        {#if selectedService.ImagemUrl || selectedService.image_url}
                                            <img
                                                src={selectedService.ImagemUrl ||
                                                    selectedService.image_url}
                                                alt={selectedService.Nome ||
                                                    selectedService.nome}
                                                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
                                            />
                                        {:else}
                                            <div
                                                class="w-full h-full flex items-center justify-center text-gray-400"
                                            >
                                                <span
                                                    class="material-symbols-outlined text-4xl"
                                                    >image</span
                                                >
                                            </div>
                                        {/if}
                                        <div
                                            class="absolute inset-0 bg-black/30 flex items-center justify-center"
                                        >
                                            <h3
                                                class="text-white font-bold text-lg drop-shadow-md text-center px-4"
                                            >
                                                {selectedService.Nome ||
                                                    selectedService.nome}
                                            </h3>
                                        </div>
                                    </div>
                                    <div class="p-5">
                                        <div
                                            class="flex justify-between items-start mb-4"
                                        >
                                            <div>
                                                <p
                                                    class="text-sm text-gray-500 dark:text-gray-400"
                                                >
                                                    Categoria
                                                </p>
                                                <p
                                                    class="font-medium text-gray-900 dark:text-white"
                                                >
                                                    {selectedService.Categoria ||
                                                        selectedService.categoria}
                                                </p>
                                            </div>
                                            <div class="text-right">
                                                <p
                                                    class="text-sm text-gray-500 dark:text-gray-400"
                                                >
                                                    Duração
                                                </p>
                                                <p
                                                    class="font-medium text-gray-900 dark:text-white"
                                                >
                                                    {selectedService.DuracaoPadrao ||
                                                        selectedService.duracao_padrao}
                                                    min
                                                </p>
                                            </div>
                                        </div>
                                        <div
                                            class="flex justify-between items-center border-t border-border-light dark:border-border-dark pt-4"
                                        >
                                            <span
                                                class="text-lg font-bold text-brand-orange"
                                                >R$ {(
                                                    (selectedService.Preco ||
                                                        selectedService.preco) /
                                                    100
                                                )
                                                    .toFixed(2)
                                                    .replace(".", ",")}</span
                                            >
                                            <button
                                                onclick={() => {
                                                    selectedService = null;
                                                    selectedProfessional = null;
                                                    selectedDay = null;
                                                    selectedSlot = null;
                                                    morningSlots = []; afternoonSlots = []; eveningSlots = [];
                                                }}
                                                class="text-sm text-brand-orange hover:text-brand-orange-hover font-medium transition-colors cursor-pointer"
                                                >Alterar Serviço</button
                                            >
                                        </div>
                                    </div>
                                </div>
                            {/if}

                            <!-- Professional Selection -->
                            <div
                                class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-5"
                            >
                                <h3
                                    class="font-semibold text-gray-900 dark:text-white mb-4"
                                >
                                    Profissional
                                </h3>
                                {#if !selectedDay}
                                    <div
                                        class="flex flex-col items-center justify-center py-8 text-center px-4"
                                    >
                                        <span
                                            class="material-icons text-4xl text-gray-300 dark:text-gray-600 mb-2"
                                            >event</span
                                        >
                                        <p
                                            class="text-gray-500 dark:text-gray-400 font-medium"
                                        >
                                            Escolha um dia no calendário
                                        </p>
                                        <p
                                            class="text-xs text-gray-400 dark:text-gray-500 mt-1"
                                        >
                                            Para ver os profissionais e horários
                                            disponíveis.
                                        </p>
                                    </div>
                                {:else if loadingProfessionals}
                                    <div class="flex justify-center py-8">
                                        <span
                                            class="material-icons animate-spin text-brand-orange text-3xl"
                                            >sync</span
                                        >
                                    </div>
                                {:else if professionals.length === 0}
                                    <p class="text-sm text-gray-500">
                                        Nenhum profissional encontrado.
                                    </p>
                                {:else}
                                    <div class="space-y-3">
                                        {#each professionals as prof}
                                            {@const profId =
                                                prof?.id || prof?.ID}
                                            {@const selectedId =
                                                selectedProfessional?.id ||
                                                selectedProfessional?.ID}
                                            <label
                                                class="flex items-center p-3 rounded-md border cursor-pointer transition-all duration-200
                                                    {selectedId &&
                                                profId &&
                                                selectedId === profId
                                                    ? 'border-brand-orange bg-brand-orange/5'
                                                    : 'border-border-light dark:border-border-dark hover:bg-gray-50 dark:hover:bg-gray-800'}"
                                            >
                                                <input
                                                    type="radio"
                                                    name="professional"
                                                    value={prof}
                                                    bind:group={
                                                        selectedProfessional
                                                    }
                                                    onchange={() => handleProfessionalChange(prof)}
                                                    class="form-radio h-4 w-4 text-brand-orange border-gray-300 focus:ring-brand"
                                                />
                                                <div
                                                    class="ml-3 flex items-center"
                                                >
                                                    {#if prof?.image_url || prof?.avatar || prof?.imagemUrl || prof?.ImagemUrl}
                                                        <img
                                                            src={prof?.image_url ||
                                                                prof?.avatar ||
                                                                prof?.imagemUrl ||
                                                                prof?.ImagemUrl}
                                                            alt={prof?.name ||
                                                                prof?.nome ||
                                                                prof?.Nome}
                                                            class="h-8 w-8 rounded-full object-cover border border-gray-200 dark:border-gray-700"
                                                            onerror={(e) => {
                                                                (
                                                                    e.target as HTMLImageElement
                                                                ).src =
                                                                    "https://ui-avatars.com/api/?name=" +
                                                                    (prof?.name ||
                                                                        prof?.nome ||
                                                                        prof?.Nome);
                                                            }}
                                                        />
                                                    {:else}
                                                        <div
                                                            class="h-8 w-8 rounded-full bg-gray-300 flex items-center justify-center text-gray-600 font-bold text-xs"
                                                        >
                                                            {(
                                                                prof?.name ||
                                                                prof?.nome ||
                                                                prof?.Nome ||
                                                                "?"
                                                            )
                                                                .substring(0, 2)
                                                                .toUpperCase()}
                                                        </div>
                                                    {/if}
                                                    <div class="ml-2">
                                                        <span
                                                            class="block text-sm font-medium text-gray-900 dark:text-white"
                                                            >{prof?.name ||
                                                                prof?.nome ||
                                                                prof?.Nome}</span
                                                        >
                                                        <span
                                                            class="block text-xs text-gray-500 dark:text-gray-400"
                                                            >{prof?.role ||
                                                                prof?.especialidade ||
                                                                prof?.Especialidade ||
                                                                "Especialista"}</span
                                                        >
                                                    </div>
                                                </div>
                                            </label>
                                        {/each}
                                    </div>
                                {/if}
                            </div>
                        </div>

                        <!-- Right Column (Calendar & Slots) -->
                        <div class="lg:col-span-2 space-y-6">
                            <!-- Calendar -->
                            <div
                                class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark p-6"
                            >
                                <div
                                    class="flex items-center justify-between mb-6"
                                >
                                    <h2
                                        class="text-xl font-bold text-gray-900 dark:text-white flex items-center"
                                    >
                                        <span class="capitalize"
                                            >{monthNames[currentMonth]}
                                            {currentYear}</span
                                        >
                                    </h2>
                                    <div class="flex space-x-2">
                                        <button
                                            onclick={() => changeMonth(-1)}
                                            class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-400 transition-colors"
                                        >
                                            <span class="material-icons"
                                                >chevron_left</span
                                            >
                                        </button>
                                        <button
                                            onclick={() => changeMonth(1)}
                                            class="p-2 rounded-full hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-400 transition-colors"
                                        >
                                            <span class="material-icons"
                                                >chevron_right</span
                                            >
                                        </button>
                                    </div>
                                </div>

                                <!-- Calendar Grid -->
                                <div
                                    class="grid grid-cols-7 gap-px lg:gap-2 text-center"
                                >
                                    <!-- Weekday Headers -->
                                    {#each ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"] as day}
                                        <div
                                            class="text-xs font-semibold text-gray-400 uppercase tracking-wider py-2"
                                        >
                                            {day}
                                        </div>
                                    {/each}

                                    <!-- Days Grid -->
                                    {#each calendarDays as day}
                                        {#if day}
                                            {@const dateObj = new Date(
                                                currentYear,
                                                currentMonth,
                                                day,
                                            )}
                                            {@const dayOfWeek =
                                                dateObj.getDay()}
                                            {@const isPast =
                                                isPastAndNotToday(day)}
                                            {@const isSelected =
                                                selectedDay === day}
                                            {@const isSunday = dayOfWeek === 0}
                                            {@const isSaturday =
                                                dayOfWeek === 6}

                                            <button
                                                disabled={isPast || isSunday}
                                                onclick={() => selectDay(day)}
                                                class="aspect-square relative flex items-start justify-end p-2 rounded-lg border transition-all group
                                            {isPast
                                                    ? 'opacity-40 cursor-not-allowed grayscale'
                                                    : 'hover:border-brand-orange'}
                                            {isSelected
                                                    ? 'bg-blue-100 dark:bg-blue-900/30 border-blue-200 dark:border-blue-500'
                                                    : isSunday
                                                      ? 'bg-red-50 dark:bg-red-900/10 border-transparent'
                                                      : isSaturday
                                                        ? 'bg-orange-50 dark:bg-orange-900/10 border-transparent'
                                                        : 'bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 border-transparent'}
                                            "
                                            >
                                                <span
                                                    class="z-10 text-sm font-medium {isSunday
                                                        ? 'text-red-500'
                                                        : 'text-gray-700 dark:text-gray-300'}"
                                                >
                                                    {day}
                                                </span>

                                                <!-- Status indicator / Label -->
                                                <div
                                                    class="absolute inset-x-1 bottom-1 flex flex-col items-center"
                                                >
                                                    {#if isSunday && !isPast}
                                                        <div
                                                            class="w-1.5 h-1.5 rounded-full  bg-red-500 mb-0.5"
                                                        ></div>
                                                    {:else if isSaturday && !isPast}
                                                        <div
                                                            class="w-1.5 h-1.5 rounded-full bg-orange-400 mb-0.5"
                                                        ></div>
                                                    {:else if !isPast}
                                                        <div class="flex gap-0.5 items-center">
                                                            <div
                                                                class="w-1.5 h-1.5 rounded-full {dayHasAppointment(day)
                                                                    ? 'bg-violet-500'
                                                                    : 'bg-emerald-500'}"
                                                            ></div>
                                                        </div>
                                                    {/if}
                                                </div>
                                            </button>
                                        {:else}
                                            <div></div>
                                        {/if}
                                    {/each}
                                </div>
                            </div>

                            <!-- Slots -->
                            {#if selectedDay}
                                <div
                                    class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-6"
                                >
                                    <h3
                                        class="text-lg font-bold text-gray-900 dark:text-white mb-4 flex items-center"
                                    >
                                        <span
                                            class="material-icons mr-2 text-brand-orange"
                                            >schedule</span
                                        >
                                        Horários Disponíveis
                                        <span
                                            class="ml-2 text-sm font-normal text-gray-500"
                                            >{selectedDay}/{currentMonth +
                                                1}</span
                                        >
                                    </h3>

                                    <div class="space-y-6">
                                        {#if morningSlots.length > 0}
                                            <div>
                                                <h4
                                                    class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                                                >
                                                    Manhã
                                                </h4>
                                                <div
                                                    class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-3"
                                                >
                                                    {#each morningSlots as slot}
                                                        <button
                                                            disabled={slot.disabled}
                                                            class="px-3 py-2 text-sm font-medium border rounded-md transition-all duration-200
                                                            {selectedSlot ===
                                                            slot.time
                                                                ? 'border-brand-orange bg-brand-orange text-white shadow-sm'
                                                                : 'border-border-light dark:border-border-dark text-gray-700 dark:text-gray-300 hover:border-brand-orange hover:text-brand-orange dark:hover:text-brand-orange'}
                                                            {slot.disabled
                                                                ? 'cursor-not-allowed bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-400 dark:text-red-500 line-through hover:border-red-200 dark:hover:border-red-800 hover:text-red-400 dark:hover:text-red-500'
                                                                : ''}"
                                                            onclick={() =>
                                                                (selectedSlot =
                                                                    slot.time)}
                                                        >
                                                            {slot.time}
                                                        </button>
                                                    {/each}
                                                </div>
                                            </div>
                                        {/if}

                                        {#if afternoonSlots.length > 0}
                                            <div>
                                                <h4
                                                    class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                                                >
                                                    Tarde
                                                </h4>
                                                <div
                                                    class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-3"
                                                >
                                                    {#each afternoonSlots as slot}
                                                        <button
                                                            disabled={slot.disabled}
                                                            class="px-3 py-2 text-sm font-medium border rounded-md transition-all duration-200
                                                            {selectedSlot ===
                                                            slot.time
                                                                ? 'border-brand-orange bg-brand-orange text-white shadow-sm'
                                                                : 'border-border-light dark:border-border-dark text-gray-700 dark:text-gray-300 hover:border-brand-orange hover:text-brand-orange dark:hover:text-brand-orange'}
                                                            {slot.disabled
                                                                ? 'cursor-not-allowed bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-400 dark:text-red-500 line-through hover:border-red-200 dark:hover:border-red-800 hover:text-red-400 dark:hover:text-red-500'
                                                                : ''}"
                                                            onclick={() =>
                                                                (selectedSlot =
                                                                    slot.time)}
                                                        >
                                                            {slot.time}
                                                        </button>
                                                    {/each}
                                                </div>
                                            </div>
                                        {/if}

                                        {#if eveningSlots.length > 0}
                                            <div>
                                                <h4
                                                    class="text-xs font-bold text-gray-500 uppercase tracking-wider mb-2"
                                                >
                                                    Noite
                                                </h4>
                                                <div
                                                    class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-3"
                                                >
                                                    {#each eveningSlots as slot}
                                                        <button
                                                            disabled={slot.disabled}
                                                            class="px-3 py-2 text-sm font-medium border rounded-md transition-all duration-200
                                                            {selectedSlot ===
                                                            slot.time
                                                                ? 'border-brand-orange bg-brand-orange text-white shadow-sm'
                                                                : 'border-border-light dark:border-border-dark text-gray-700 dark:text-gray-300 hover:border-brand-orange hover:text-brand-orange dark:hover:text-brand-orange'}
                                                            {slot.disabled
                                                                ? 'cursor-not-allowed bg-red-50 dark:bg-red-900/20 border-red-200 dark:border-red-800 text-red-400 dark:text-red-500 line-through hover:border-red-200 dark:hover:border-red-800 hover:text-red-400 dark:hover:text-red-500'
                                                                : ''}"
                                                            onclick={() =>
                                                                (selectedSlot =
                                                                    slot.time)}
                                                        >
                                                            {slot.time}
                                                        </button>
                                                    {/each}
                                                </div>
                                            </div>
                                        {/if}

                                        {#if morningSlots.length === 0 && afternoonSlots.length === 0 && eveningSlots.length === 0}
                                            <div class="py-4 text-center">
                                                <p
                                                    class="text-sm text-gray-500 italic"
                                                >
                                                    Sem horários disponíveis
                                                    para este profissional neste
                                                    dia.
                                                </p>
                                            </div>
                                        {/if}
                                    </div>
                                </div>
                            {/if}

                            <!-- Actions -->
                            <div
                                class="flex items-center justify-end space-x-4 pt-4 border-t border-border-light dark:border-border-dark"
                            >
                                <button
                                    onclick={() => {
                                        selectedDay = null;
                                        selectedSlot = null;
                                        morningSlots = []; afternoonSlots = []; eveningSlots = [];
                                        selectedProfessional = null;
                                        error = "";
                                    }}
                                    class="px-6 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer"
                                >
                                    Cancelar
                                </button>
                                <button
                                    onclick={handleConfirmAppointment}
                                    disabled={submitting}
                                    class="px-8 py-2 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-md font-medium shadow-md transition-all duration-200 flex items-center hover:shadow-lg disabled:opacity-70 disabled:cursor-not-allowed"
                                >
                                    {#if submitting}
                                        <span
                                            class="material-icons animate-spin mr-2 text-sm"
                                            >sync</span
                                        >
                                        Enviando...
                                    {:else}
                                        Confirmar Agendamento
                                        <span
                                            class="material-icons text-sm ml-2"
                                            >arrow_forward</span
                                        >
                                    {/if}
                                </button>
                            </div>
                        </div>
                    </div>
                {/if}
            </div>

            <!-- Footer inside main content to scroll together -->
            <footer
                class="mt-12 text-center text-sm text-gray-500 dark:text-gray-400 pb-6"
            >
                <p>© 2023 BellaVita. Todos os direitos reservados.</p>
                <div class="mt-2 space-x-4">
                    <a href="/" class="hover:text-brand transition-colors"
                        >Termos</a
                    >
                    <a href="/" class="hover:text-brand transition-colors"
                        >Privacidade</a
                    >
                </div>
            </footer>
        </div>
    </main>

    <!-- Success Modal -->
    {#if showSuccessModal}
        <div
            class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm transition-opacity duration-300"
        >
            <div
                class="bg-[hsl(var(--bs-card))] rounded-2xl shadow-2xl max-w-md w-full p-8 transform transition-all duration-300 scale-100 border border-border-light dark:border-border-dark text-center"
            >
                <div
                    class="w-20 h-20 bg-green-100 dark:bg-green-900/30 rounded-full flex items-center justify-center mx-auto mb-6"
                >
                    <span class="material-icons text-green-500 text-4xl"
                        >check_circle</span
                    >
                </div>

                <h2
                    class="text-2xl font-bold text-gray-900 dark:text-white mb-2"
                >
                    Agendamento Realizado!
                </h2>
                <p class="text-gray-600 dark:text-gray-400 mb-8">
                    Seu agendamento foi realizado com sucesso! Estamos
                    aguardando a confirmação do profissional.
                </p>

                <div
                    class="bg-gray-50 dark:bg-gray-800/50 rounded-xl p-4 mb-8 text-left space-y-2 border border-border-light dark:border-border-dark"
                >
                    <div class="flex items-center text-sm">
                        <span
                            class="material-icons text-brand-orange text-lg mr-2"
                            >event</span
                        >
                        <span
                            class="text-gray-700 dark:text-gray-300 font-medium"
                        >
                            {selectedDay}/{currentMonth + 1}/{currentYear} às {selectedSlot}
                        </span>
                    </div>
                    <div class="flex items-center text-sm">
                        <span
                            class="material-icons text-brand-orange text-lg mr-2"
                            >person</span
                        >
                        <span
                            class="text-gray-700 dark:text-gray-300 font-medium"
                        >
                            {selectedProfessional?.name ||
                                selectedProfessional?.nome ||
                                selectedProfessional?.Nome}
                        </span>
                    </div>
                </div>

                <div class="space-y-3">
                    <button
                        onclick={() => (window.location.href = "/clientes/meus-agendamentos")}
                        class="w-full py-3 px-4 bg-brand-orange hover:bg-brand-orange/90 text-white font-bold rounded-xl transition-all duration-200 shadow-lg shadow-brand-orange/20"
                    >
                        Ver Meus Agendamentos
                    </button>
                    <button
                        onclick={() => (showSuccessModal = false)}
                        class="w-full py-3 px-4 bg-transparent hover:bg-gray-100 dark:hover:bg-gray-800 text-gray-600 dark:text-gray-400 font-medium rounded-xl transition-colors"
                    >
                        Fechar
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <!-- Error Modal -->
    <AlertModal
        show={showErrorModal}
        title={errorTitle}
        message={errorMessage}
        type="error"
        onconfirm={() => (showErrorModal = false)}
        oncancel={() => (showErrorModal = false)}
    />
</div>
