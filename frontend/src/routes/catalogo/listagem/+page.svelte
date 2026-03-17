<script lang="ts">
    import { onMount } from "svelte";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import Modal from "$lib/components/Modal.svelte";
    import ImageUpload from "$lib/components/ImageUpload.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import { fetchApi } from "$lib/utils/api";
    import {
        uploadImageToSupabase,
        deleteImageFromSupabase,
    } from "$lib/utils/imageUpload";
    import { validateEditServiceForm } from "$lib/utils/validation";

    interface Service {
        ID: string;
        Nome: string;
        Categoria: string;
        DuracaoPadrao: number;
        Preco: number;
        ImagemUrl: string;
    }

    interface ApiResponse {
        data: Service[];
        limit: number;
        page: number;
        total: number;
    }

    let services = $state<Service[]>([]);
    let loading = $state(true);
    let error = $state("");

    // Paginação
    let page = $state(1);
    let limit = $state(12);
    let totalItems = $state(0);

    let totalPages = $derived(Math.ceil(totalItems / limit));

    // Edição
    let showEditModal = $state(false);
    let editingService = $state<Service | null>(null);
    let isUpdating = $state(false);

    // Campos do formulário de edição
    let editNome = $state("");
    let editCategoria = $state("");
    let editDuracao = $state(60);
    let editPreco = $state(""); // Preço visual em Reais ("120,00")

    // Imagem do modal de edição
    let selectedImage = $state<File | null>(null);
    let imagePreview = $state<string | null>(null);

    // Erros de validação
    let editErrors = $state<Record<string, string>>({});

    // Cache busting para imagens (timestamp por serviço)
    let imageTimestamps = $state(new Map<string, number>());

    // Estado do AlertModal
    let alertState = $state({
        show: false,
        title: "",
        message: "",
        type: "info" as "success" | "error" | "warning" | "info",
        showCancel: false,
        confirmText: "OK",
        onConfirm: () => {},
    });

    function showAlert(
        title: string,
        message: string,
        type: "success" | "error" | "warning" | "info" = "info",
        showCancel = false,
        onConfirm: () => void = () => {
            alertState.show = false;
        },
    ) {
        alertState = {
            show: true,
            title,
            message,
            type,
            showCancel,
            confirmText: showCancel ? "Confirmar" : "OK",
            onConfirm,
        };
    }

    function closeAlert() {
        alertState.show = false;
    }

    async function fetchServices() {
        loading = true;
        error = "";
        try {
            const response = await fetch(
                `/api/v1/catalogos?page=${page}&limit=${limit}`,
            );

            if (response.ok) {
                const data: ApiResponse = await response.json();
                services = data.data || [];
                totalItems = data.total || 0;
            } else {
                const text = await response.text();
                error = "Erro ao carregar serviços: " + response.statusText;
            }
        } catch (err) {
            error = "Serviço indisponível no momento.";
        } finally {
            loading = false;
        }
    }

    function nextPage() {
        if (page < totalPages) {
            page++;
            fetchServices();
        }
    }

    function prevPage() {
        if (page > 1) {
            page--;
            fetchServices();
        }
    }

    function openEditModal(service: Service) {
        editingService = service;
        editNome = service.Nome;
        editCategoria = service.Categoria;
        editDuracao = service.DuracaoPadrao;
        // Preço vem em centavos (12000), converte para string visual (120,00)
        editPreco = (service.Preco / 100).toFixed(2).replace(".", ",");

        // Reset imagem e erros
        selectedImage = null;
        imagePreview = service.ImagemUrl;
        editErrors = {};

        showEditModal = true;
    }

    function handleImageSelect(data: { file: File; preview: string }) {
        selectedImage = data.file;
        imagePreview = data.preview;
    }

    function handleImageClear() {
        selectedImage = null;
        imagePreview = null;
    }

    async function handleUpdate() {
        if (!editingService) return;

        // Validação usando utilitário extraído
        editErrors = validateEditServiceForm(
            {
                nome: editNome,
                categoria: editCategoria,
                duracao: editDuracao,
                preco: editPreco,
            },
            !!(imagePreview || selectedImage),
        );

        if (Object.keys(editErrors).length > 0) {
            return;
        }

        isUpdating = true;

        try {
            // Converter preço visual "120,00" para centavos 12000
            const precoCentavos = Math.round(
                parseFloat(editPreco.replace(",", ".")) * 100,
            );

            // Upload de imagem se houver nova
            let imageUrl = editingService.ImagemUrl;

            if (selectedImage) {
                const newImageUrl = await uploadImageToSupabase(selectedImage);
                if (newImageUrl) {
                    if (
                        editingService.ImagemUrl &&
                        editingService.ImagemUrl !== newImageUrl
                    ) {
                        await deleteImageFromSupabase(editingService.ImagemUrl);
                    }
                    imageUrl = newImageUrl;
                } else {
                    showAlert(
                        "Erro",
                        "Erro ao fazer upload da imagem.",
                        "error",
                    );
                    isUpdating = false;
                    return;
                }
            } else if (!imagePreview) {
                showAlert("Erro", "A imagem é obrigatória.", "error");
                isUpdating = false;
                return;
            }

            const payload = {
                categoria: editCategoria,
                duracao_padrao: editDuracao,
                image_url: imageUrl,
                nome: editNome,
                preco: precoCentavos,
            };

            const response = await fetchApi(
                `/api/v1/catalogos/${editingService.ID}`,
                {
                    method: "PUT",
                    body: JSON.stringify(payload),
                },
            );

            if (response.ok) {
                // Atualiza timestamp da imagem se foi alterada
                if (selectedImage) {
                    const newMap = new Map(imageTimestamps);
                    newMap.set(editingService!.ID, Date.now());
                    imageTimestamps = newMap;
                }

                // Atualiza lista localmente (Otimista)
                services = services.map((s) => {
                    if (s.ID === editingService!.ID) {
                        return {
                            ...s,
                            Nome: editNome,
                            Categoria: editCategoria,
                            DuracaoPadrao: editDuracao,
                            Preco: precoCentavos,
                            ImagemUrl: imageUrl,
                        };
                    }
                    return s;
                });
                showEditModal = false;
            } else {
                showAlert(
                    "Erro",
                    "Erro ao atualizar serviço: " + response.statusText,
                    "error",
                );
            }
        } catch (err) {
            console.error(err);
            showAlert("Erro", "Erro de conexão ao atualizar.", "error");
        } finally {
            isUpdating = false;
        }
    }

    async function handleDelete(service: Service) {
        showAlert(
            "Excluir Serviço",
            `Tem certeza que deseja excluir o serviço "${service.Nome}"?`,
            "warning",
            true, // showCancel
            async () => {
                try {
                    const response = await fetchApi(
                        `/api/v1/catalogos/${service.ID}`,
                        { method: "DELETE" },
                    );

                    if (response.ok) {
                        if (service.ImagemUrl) {
                            const deleted = await deleteImageFromSupabase(
                                service.ImagemUrl,
                            );
                            if (!deleted) {
                                console.warn(
                                    "Aviso: Não foi possível deletar a imagem do storage.",
                                );
                            }
                        }

                        services = services.filter((s) => s.ID !== service.ID);
                        totalItems--;
                        if (services.length === 0 && page > 1) {
                            page--;
                            fetchServices();
                        }
                        alertState.show = false;
                    } else {
                        const errorData = await response.json().catch(() => ({}));
                        const errorMessage = errorData.error || await response.text();

                        let friendlyMessage = `Erro ao excluir serviço: ${errorMessage}`;

                        if (errorMessage.includes("falha na infraestrutura")) {
                            friendlyMessage = "Este serviço não pode ser excluído porque está vinculado a profissionais ou agendamentos.";
                        }

                        showAlert(
                            "Erro ao Excluir",
                            friendlyMessage,
                            "error",
                        );
                    }
                } catch (err) {
                    console.error("Erro ao excluir:", err);
                    showAlert(
                        "Erro",
                        "Erro de conexão ao tentar excluir.",
                        "error",
                    );
                }
            },
        );
    }

    onMount(fetchServices);
</script>

<div
    class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <Sidebar />

    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar />

        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-5xl mx-auto space-y-6">
                <div class="flex items-center justify-between">
                    <div>
                        <nav
                            class="flex text-sm text-gray-500 dark:text-gray-400 mb-2"
                        >
                            <a
                                class="hover:text-brand-orange transition-colors"
                                href="/">Início</a
                            >
                            <span class="mx-2">/</span>
                            <span
                                class="text-gray-900 dark:text-white font-medium"
                                >Serviços Cadastrados</span
                            >
                        </nav>
                        <h1
                            class="text-2xl font-bold text-gray-900 dark:text-white"
                        >
                            Catálogo de Serviços
                        </h1>
                        <p class="text-gray-500 dark:text-gray-400 mt-1">
                            Gerencie os serviços oferecidos no salão.
                        </p>
                    </div>
                    <a
                        href="/catalogo"
                        class="px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-md font-medium shadow-md transition-colors flex items-center"
                    >
                        <span class="material-symbols-outlined text-sm mr-2"
                            >add</span
                        >
                        Novo Serviço
                    </a>
                </div>

                {#if loading}
                    <div class="flex justify-center py-12">
                        <span
                            class="material-icons animate-spin text-brand-orange text-4xl"
                            >sync</span
                        >
                    </div>
                {:else if error}
                    <div
                        class="bg-red-50 text-red-700 p-4 rounded-md border border-red-200"
                    >
                        {error}
                    </div>
                {:else if services.length === 0}
                    <div
                        class="bg-[hsl(var(--bs-card))] rounded-lg border-2 border-dashed border-border-light dark:border-border-dark p-12 text-center text-gray-500 dark:text-gray-400"
                    >
                        <span class="material-icons text-5xl mb-4 opacity-20"
                            >inventory_2</span
                        >
                        <p>Nenhum serviço cadastrado ainda.</p>
                        <a
                            href="/catalogo"
                            class="text-brand-orange hover:underline mt-2 inline-block"
                            >Cadastrar primeiro serviço</a
                        >
                    </div>
                {:else}
                    <div
                        class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8"
                    >
                        {#each services as service}
                            <div
                                class="bg-[hsl(var(--bs-card))] rounded-2xl border border-border-light dark:border-border-dark overflow-hidden hover:shadow-xl transition-all duration-300 group flex flex-col relative"
                            >
                                <!-- Top Action Link (Image Area) -->
                                <a
                                    href="/catalogo/listagem/{service.ID}"
                                    class="absolute inset-x-0 top-0 h-48 z-10"
                                    aria-label="Ver detalhes de {service.Nome}"
                                ></a>

                                {#if service.ImagemUrl}
                                    <div
                                        class="h-48 overflow-hidden relative bg-gray-100 dark:bg-gray-800"
                                    >
                                        <img
                                            src={imageTimestamps.has(service.ID)
                                                ? `${service.ImagemUrl}?t=${imageTimestamps.get(service.ID)}`
                                                : service.ImagemUrl}
                                            alt={service.Nome}
                                            class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                                            onerror={(e) => {
                                                (
                                                    e.target as HTMLImageElement
                                                ).src =
                                                    "https://images.unsplash.com/photo-1521590832167-7bcbfaa6381f?auto=format&fit=crop&q=80&w=600";
                                            }}
                                        />
                                        <div
                                            class="absolute inset-0 bg-gradient-to-t from-black/40 via-transparent to-transparent"
                                        ></div>
                                    </div>
                                {:else}
                                    <div
                                        class="h-48 bg-gray-50 dark:bg-gray-800 flex items-center justify-center text-gray-400"
                                    >
                                        <span
                                            class="material-symbols-outlined text-6xl opacity-20"
                                            >image</span
                                        >
                                    </div>
                                {/if}

                                <div
                                    class="p-6 pt-0 -mt-6 relative z-20 flex-1 flex flex-col"
                                >
                                    <div
                                        class="flex items-start justify-between mb-4"
                                    >
                                        <!-- Category Icon Box -->
                                        <div
                                            class="h-14 w-14 bg-white dark:bg-gray-800 rounded-xl shadow-lg border border-border-light dark:border-border-dark flex items-center justify-center text-brand-orange transform group-hover:-translate-y-1 transition-transform duration-300"
                                        >
                                            <span
                                                class="material-symbols-outlined text-3xl"
                                            >
                                                {service.Categoria === "Cabelo"
                                                    ? "content_cut"
                                                    : service.Categoria ===
                                                        "Unhas"
                                                      ? "front_hand"
                                                      : service.Categoria ===
                                                          "Massagem"
                                                        ? "self_improvement"
                                                        : service.Categoria ===
                                                            "Estética Facial"
                                                          ? "face"
                                                          : "spa"}
                                            </span>
                                        </div>

                                        <!-- Duration only in top -->
                                        <div class="text-right">
                                            <p
                                                class="text-[10px] font-bold text-gray-400 dark:text-gray-500 uppercase tracking-widest mt-1"
                                            >
                                                {service.DuracaoPadrao} MIN
                                            </p>
                                        </div>
                                    </div>

                                    <div class="mt-2">
                                        <div
                                            class="flex items-start justify-between gap-3"
                                        >
                                            <a
                                                href="/catalogo/listagem/{service.ID}"
                                                class="block group/link flex-1"
                                            >
                                                <h3
                                                    class="text-xl font-bold text-gray-900 dark:text-white group-hover/link:text-brand-orange transition-colors duration-300 leading-tight"
                                                >
                                                    {service.Nome}
                                                </h3>
                                            </a>
                                            <p
                                                class="text-xl font-bold text-brand-orange whitespace-nowrap"
                                            >
                                                R$ {(service.Preco / 100)
                                                    .toFixed(2)
                                                    .replace(".", ",")}
                                            </p>
                                        </div>
                                        <p
                                            class="text-sm font-semibold text-gray-400 dark:text-gray-500 mt-1 uppercase tracking-wide"
                                        >
                                            {service.Categoria}
                                        </p>
                                    </div>

                                    <!-- Bottom Info -->
                                    <div
                                        class="mt-auto pt-6 flex items-center justify-between border-t border-border-light/60 dark:border-border-dark/60"
                                    >
                                        <span
                                            class="text-[11px] font-mono text-gray-400 dark:text-gray-500 font-medium"
                                            >ID: #{service.ID?.substring(
                                                0,
                                                8,
                                            ) || "N/A"}</span
                                        >
                                        <div class="flex items-center gap-1">
                                            <button
                                                class="p-2 text-gray-400 hover:text-brand-orange dark:hover:text-brand-orange transition-colors rounded-lg hover:bg-brand-orange/5"
                                                title="Editar"
                                                onclick={() =>
                                                    openEditModal(service)}
                                            >
                                                <span
                                                    class="material-symbols-outlined text-[20px]"
                                                    >edit</span
                                                >
                                            </button>
                                            <button
                                                class="p-2 text-gray-400 hover:text-red-500 dark:hover:text-red-500 transition-colors rounded-lg hover:bg-red-500/5"
                                                title="Excluir"
                                                onclick={() =>
                                                    handleDelete(service)}
                                            >
                                                <span
                                                    class="material-symbols-outlined text-[20px]"
                                                    >delete</span
                                                >
                                            </button>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        {/each}
                    </div>

                    <!-- Pagination -->
                    {#if totalPages > 1}
                        <div
                            class="mt-8 flex justify-center items-center space-x-4"
                        >
                            <button
                                class="px-4 py-2 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                onclick={prevPage}
                                disabled={page === 1}
                            >
                                Anterior
                            </button>
                            <span class="text-gray-600 dark:text-gray-400">
                                Página {page} de {totalPages}
                            </span>
                            <button
                                class="px-4 py-2 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                onclick={nextPage}
                                disabled={page === totalPages}
                            >
                                Próxima
                            </button>
                        </div>
                    {/if}
                {/if}
            </div>

            <!-- Modal de Edição -->
            {#snippet editModalBody()}
                <div class="space-y-4">
                    <div>
                        <div class="mb-4">
                            <span
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
                            >
                                Imagem do Serviço
                            </span>
                            <ImageUpload
                                previewUrl={imagePreview}
                                onselect={handleImageSelect}
                                onclear={handleImageClear}
                            />
                            {#if editErrors.imagem}
                                <p class="mt-1 text-xs text-red-500">
                                    {editErrors.imagem}
                                </p>
                            {/if}
                        </div>

                        <label
                            for="edit-nome"
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                        >
                            Nome do Serviço
                        </label>
                        <input
                            id="edit-nome"
                            type="text"
                            bind:value={editNome}
                            class="block w-full px-3 py-2 border {editErrors.nome
                                ? 'border-red-500'
                                : 'border-border-light dark:border-border-dark'} rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
                        />
                        {#if editErrors.nome}
                            <p class="mt-1 text-xs text-red-500">
                                {editErrors.nome}
                            </p>
                        {/if}
                    </div>

                    <div>
                        <label
                            for="edit-categoria"
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                        >
                            Categoria
                        </label>
                        <select
                            id="edit-categoria"
                            bind:value={editCategoria}
                            class="block w-full px-3 py-2 border {editErrors.categoria
                                ? 'border-red-500'
                                : 'border-border-light dark:border-border-dark'} rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
                        >
                            <option value="Cabelo">Cabelo</option>
                            <option value="Unhas">Unhas</option>
                            <option value="Estética Facial"
                                >Estética Facial</option
                            >
                            <option value="Massagem">Massagem</option>
                            <option value="Outros">Outros</option>
                        </select>
                        {#if editErrors.categoria}
                            <p class="mt-1 text-xs text-red-500">
                                {editErrors.categoria}
                            </p>
                        {/if}
                    </div>

                    <div class="grid grid-cols-2 gap-4">
                        <div>
                            <label
                                for="edit-preco"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                            >
                                Preço (R$)
                            </label>
                            <input
                                id="edit-preco"
                                type="text"
                                bind:value={editPreco}
                                placeholder="0,00"
                                class="block w-full px-3 py-2 border {editErrors.preco
                                    ? 'border-red-500'
                                    : 'border-border-light dark:border-border-dark'} rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
                            />
                            {#if editErrors.preco}
                                <p class="mt-1 text-xs text-red-500">
                                    {editErrors.preco}
                                </p>
                            {/if}
                        </div>

                        <div>
                            <label
                                for="edit-duracao"
                                class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                            >
                                Duração (min)
                            </label>
                            <input
                                id="edit-duracao"
                                type="number"
                                bind:value={editDuracao}
                                class="block w-full px-3 py-2 border {editErrors.duracao
                                    ? 'border-red-500'
                                    : 'border-border-light dark:border-border-dark'} rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
                            />
                            {#if editErrors.duracao}
                                <p class="mt-1 text-xs text-red-500">
                                    {editErrors.duracao}
                                </p>
                            {/if}
                        </div>
                    </div>
                </div>
            {/snippet}

            {#snippet editModalFooter()}
                <button
                    onclick={() => (showEditModal = false)}
                    class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                >
                    Cancelar
                </button>
                <button
                    onclick={handleUpdate}
                    disabled={isUpdating}
                    class="px-4 py-2 text-sm font-medium text-white bg-brand-orange rounded-md hover:bg-brand-orange-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-orange disabled:opacity-50 transition-colors"
                >
                    {isUpdating ? "Salvando..." : "Salvar Alterações"}
                </button>
            {/snippet}

            <Modal
                show={showEditModal}
                title="Editar Serviço"
                onclose={() => (showEditModal = false)}
                children={editModalBody}
                footer={editModalFooter}
            />

            <!-- Alert Modal -->
            <AlertModal
                show={alertState.show}
                title={alertState.title}
                message={alertState.message}
                type={alertState.type}
                showCancel={alertState.showCancel}
                confirmText={alertState.confirmText}
                onconfirm={alertState.onConfirm}
                oncancel={closeAlert}
            />

            <footer
                class="mt-12 text-center text-sm text-gray-500 dark:text-gray-400 pb-6"
            >
                <p>© 2023 BellaVita Salon. Todos os direitos reservados.</p>
            </footer>
        </div>
    </main>
</div>

<style>
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
