<script lang="ts">
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import ImageUpload from "$lib/components/ImageUpload.svelte";
    import { uploadPrestadorImageToSupabase } from "$lib/utils/imageUpload";

    // Interface do Serviço (Catálogo)
    interface CatalogService {
        ID: string;
        Nome: string;
        Categoria: string;
        DuracaoPadrao: number;
    }

    interface CatalogResponse {
        data: CatalogService[];
        total: number;
    }

    // Estado do formulário
    let nome = $state("");
    let cpf = $state("");
    let phone = $state("");
    let email = $state("");
    let selectedCatalogIds = $state<string[]>([]);

    // Estado da imagem
    let selectedImage = $state<File | null>(null);
    let imagePreview = $state<string | null>(null);

    // Estado da lista de serviços
    let services = $state<CatalogService[]>([]);
    let loadingServices = $state(true);
    let page = $state(1);
    let limit = $state(12);
    let totalServices = $state(0);

    // Estado de envio
    let isSubmitting = $state(false);

    // Alert Modal
    let alertState = $state({
        show: false,
        title: "",
        message: "",
        type: "info" as "success" | "error" | "warning" | "info",
        onConfirm: () => {
            alertState.show = false;
        },
    });

    function showAlert(
        title: string,
        message: string,
        type: "success" | "error" | "warning" | "info" = "info",
        onConfirm?: () => void,
    ) {
        alertState = {
            show: true,
            title,
            message,
            type,
            onConfirm:
                onConfirm ||
                (() => {
                    alertState.show = false;
                }),
        };
    }

    async function fetchServices(reset = false) {
        if (reset) {
            page = 1;
            services = [];
            loadingServices = true;
        }

        try {
            const res = await fetch(
                `/api/v1/catalogos?page=${page}&limit=${limit}`,
            );
            if (res.ok) {
                const data: CatalogResponse = await res.json();
                if (reset) {
                    services = data.data || [];
                } else {
                    services = [...services, ...(data.data || [])];
                }
                totalServices = data.total;
            } else {
                console.error("Erro ao buscar serviços:", await res.text());
                showAlert(
                    "Erro",
                    "Falha ao carregar lista de serviços.",
                    "error",
                );
            }
        } catch (error) {
            console.error(error);
            showAlert("Erro", "Erro de conexão ao buscar serviços.", "error");
        } finally {
            loadingServices = false;
        }
    }

    function loadMore() {
        page++;
        fetchServices();
    }

    async function handleSubmit() {
        // Validação do Nome (min=3, max=100)
        if (!nome || nome.trim().length < 3) {
            showAlert(
                "Atenção",
                "O nome deve ter no mínimo 3 caracteres.",
                "warning",
            );
            return;
        }
        if (nome.trim().length > 100) {
            showAlert(
                "Atenção",
                "O nome deve ter no máximo 100 caracteres.",
                "warning",
            );
            return;
        }

        // Validação do CPF (11 dígitos numéricos)
        const cpfDigits = cpf.replace(/\D/g, "");
        if (!cpfDigits || cpfDigits.length !== 11) {
            showAlert(
                "Atenção",
                "O CPF deve conter exatamente 11 dígitos.",
                "warning",
            );
            return;
        }

        // Validação do Email (formato válido)
        if (!email || email.trim() === "") {
            showAlert("Atenção", "O email é obrigatório.", "warning");
            return;
        }
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(email)) {
            showAlert(
                "Atenção",
                "Por favor, insira um email válido.",
                "warning",
            );
            return;
        }

        // Validação do Telefone (min=8, max=15 dígitos)
        const phoneDigits = phone.replace(/\D/g, "");
        if (!phoneDigits || phoneDigits.length < 8) {
            showAlert(
                "Atenção",
                "O telefone deve ter no mínimo 8 dígitos.",
                "warning",
            );
            return;
        }
        if (phoneDigits.length > 15) {
            showAlert(
                "Atenção",
                "O telefone deve ter no máximo 15 dígitos.",
                "warning",
            );
            return;
        }

        // Validação da Foto
        if (!selectedImage) {
            showAlert(
                "Atenção",
                "Selecione uma foto para o profissional.",
                "warning",
            );
            return;
        }

        // Validação dos Serviços (pelo menos um)
        if (selectedCatalogIds.length === 0) {
            showAlert(
                "Atenção",
                "Selecione pelo menos um serviço para o profissional.",
                "warning",
            );
            return;
        }

        isSubmitting = true;

        // 1. Upload da imagem
        let imageUrl = "";
        try {
            const uploadedUrl =
                await uploadPrestadorImageToSupabase(selectedImage);
            if (uploadedUrl) {
                imageUrl = uploadedUrl;
            } else {
                showAlert("Erro", "Erro ao fazer upload da imagem.", "error");
                isSubmitting = false;
                return;
            }
        } catch (error) {
            console.error(error);
            showAlert("Erro", "Erro ao fazer upload da imagem.", "error");
            isSubmitting = false;
            return;
        }

        // 2. Criar payload com a URL da imagem
        const payload = {
            catalogo_ids: selectedCatalogIds,
            cpf: cpfDigits,
            email: email.trim(),
            image_url: imageUrl,
            nome: nome.trim(),
            telefone: phoneDigits,
            senha: "admin123",
        };

        try {
            const res = await fetch("/api/v1/prestadores", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (res.ok) {
                const responseData = await res.json();
                const newId = responseData.id || responseData.ID;

                showAlert(
                    "Sucesso",
                    "Profissional cadastrado com sucesso!",
                    "success",
                    () => {
                        if (newId) {
                            goto(`/prestadores/${newId}`);
                        } else {
                            resetForm();
                        }
                    },
                );
            } else {
                const text = await res.text();
                let errorMessage = text;

                if (text.includes('duplicate key value violates unique constraint "uq_prestadores_email"')) {
                     errorMessage = "Não foi possível realizar o cadastro. Verifique os dados informados.";
                } else if (
                    text.toLowerCase().includes("cpf") ||
                    text.toLowerCase().includes("invalid")
                ) {
                    errorMessage =
                        "CPF inválido. Por favor, verifique o número digitado.";
                } else if (
                    text.toLowerCase().includes("catalogo") ||
                    text.toLowerCase().includes("serviço")
                ) {
                    errorMessage =
                        "O profissional deve ter pelo menos um serviço cadastrado.";
                }

                showAlert("Erro", errorMessage, "error");
            }
        } catch (err) {
            console.error(err);
            showAlert("Erro", "Erro de conexão ao salvar profissional.", "error");
        } finally {
            isSubmitting = false;
        }
    }

    function resetForm() {
        nome = "";
        cpf = "";
        phone = "";
        email = "";
        selectedCatalogIds = [];
        selectedImage = null;
        imagePreview = null;
        window.scrollTo(0, 0);
    }

    function handleImageSelect(data: { file: File; preview: string }) {
        selectedImage = data.file;
        imagePreview = data.preview;
    }

    function handleImageClear() {
        selectedImage = null;
        imagePreview = null;
    }

    onMount(() => {
        fetchServices(true);
    });
</script>

<div
    class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <Sidebar />
    <main class="flex-1 flex flex-col h-full overflow-hidden relative">
        <DashboardNavbar />

        <div class="flex-1 overflow-y-auto p-6 md:p-8">
            <div class="max-w-4xl mx-auto space-y-6">
                <div>
                    <nav
                        class="flex text-sm text-gray-500 dark:text-gray-400 mb-2"
                    >
                        <a class="hover:text-primary transition-colors" href="/"
                            >Admin</a
                        >
                        <span class="mx-2">/</span>
                        <span class="text-gray-900 dark:text-white font-medium"
                            >Novo Profissional</span
                        >
                    </nav>
                    <h1
                        class="text-2xl font-bold text-gray-900 dark:text-white"
                    >
                        Registro de Profissionais
                    </h1>
                    <p class="text-gray-500 dark:text-gray-400 mt-1">
                        Adicione um novo profissional e atribua seus serviços.
                    </p>
                </div>

                <div
                    class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark overflow-hidden"
                >
                    <div class="p-6 md:p-8 space-y-8">
                        <div>
                            <h3
                                class="text-lg font-semibold text-gray-900 dark:text-white flex items-center mb-6"
                            >
                                <span
                                    class="material-symbols-outlined mr-2 text-primary"
                                    >person</span
                                >
                                Dados Pessoais
                            </h3>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div class="col-span-1 md:col-span-2">
                                    <label
                                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                                        for="name">Nome Completo</label
                                    >
                                    <input
                                        bind:value={nome}
                                        class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-input-focus focus:ring-input-focus sm:text-sm py-2.5 transition-colors"
                                        id="name"
                                        placeholder="Ex: Ana Maria Costa"
                                        type="text"
                                    />
                                </div>
                                <div>
                                    <label
                                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                                        for="cpf">CPF</label
                                    >
                                    <input
                                        bind:value={cpf}
                                        class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-input-focus focus:ring-input-focus sm:text-sm py-2.5 transition-colors"
                                        id="cpf"
                                        placeholder="000.000.000-00"
                                        type="text"
                                    />
                                </div>
                                <div>
                                    <label
                                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                                        for="phone">Telefone / WhatsApp</label
                                    >
                                    <input
                                        bind:value={phone}
                                        class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-input-focus focus:ring-input-focus sm:text-sm py-2.5 transition-colors"
                                        id="phone"
                                        placeholder="(00) 00000-0000"
                                        type="tel"
                                    />
                                </div>
                                <div class="col-span-1 md:col-span-2">
                                    <label
                                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                                        for="email">Endereço de Email</label
                                    >
                                    <div class="relative">
                                        <div
                                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                                        >
                                            <span
                                                class="material-symbols-outlined text-gray-400 text-[20px]"
                                                >mail</span
                                            >
                                        </div>
                                        <input
                                            bind:value={email}
                                        class="block w-full pl-10 rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-input-focus focus:ring-input-focus sm:text-sm py-2.5 transition-colors"
                                            id="email"
                                            placeholder="profissional@bellavita.com"
                                            type="email"
                                        />
                                    </div>
                                </div>
                                <div class="col-span-1 md:col-span-2">
                                    <div
                                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                                    >
                                        Foto do Profissional
                                    </div>
                                    <ImageUpload
                                        previewUrl={imagePreview}
                                        onselect={handleImageSelect}
                                        onclear={handleImageClear}
                                    />
                                </div>
                            </div>
                        </div>

                        <div
                            class="border-t border-border-light dark:border-border-dark"
                        ></div>

                        <div>
                            <h3
                                class="text-lg font-semibold text-gray-900 dark:text-white flex items-center mb-2"
                            >
                                <span
                                    class="material-symbols-outlined mr-2 text-primary"
                                    >spa</span
                                >
                                Serviços Atribuídos
                            </h3>
                            <p
                                class="text-sm text-gray-500 dark:text-gray-400 mb-6"
                            >
                                Selecione todos os serviços que este
                                profissional está qualificado para realizar.
                            </p>

                            <!-- Services Grid -->
                            <div
                                class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
                            >
                                {#if loadingServices && services.length === 0}
                                    <div
                                        class="col-span-full py-8 text-center text-gray-500"
                                    >
                                        Carregando serviços...
                                    </div>
                                {:else}
                                    {#each services as service}
                                        <label
                                            class="relative flex items-start p-4 rounded-lg border border-border-light dark:border-border-dark bg-gray-50 dark:bg-gray-800/50 hover:bg-white dark:hover:bg-gray-800 cursor-pointer transition-all group"
                                            class:border-primary={selectedCatalogIds.includes(
                                                service.ID,
                                            )}
                                            class:bg-primary-50={selectedCatalogIds.includes(
                                                service.ID,
                                            )}
                                        >
                                            <div class="min-w-0 flex-1 text-sm">
                                                <div
                                                    class="font-medium text-gray-900 dark:text-white"
                                                >
                                                    {service.Nome}
                                                </div>
                                                <p
                                                    class="text-gray-500 dark:text-gray-400 text-xs mt-1"
                                                >
                                                    {service.DuracaoPadrao} min •
                                                    {service.Categoria}
                                                </p>
                                            </div>
                                            <div
                                                class="ml-3 flex items-center h-5"
                                            >
                                                <input
                                                    class="h-5 w-5 text-input-focus border-gray-300 rounded focus:ring-input-focus dark:bg-gray-700 dark:border-gray-600"
                                                    type="checkbox"
                                                    value={service.ID}
                                                    bind:group={
                                                        selectedCatalogIds
                                                    }
                                                />
                                            </div>
                                        </label>
                                    {/each}
                                {/if}
                            </div>

                            <!-- Load More -->
                            {#if services.length < totalServices}
                                <div class="mt-4 text-center">
                                    <button
                                        type="button"
                                        onclick={loadMore}
                                        class="text-sm text-primary hover:underline hover:text-primary-hover font-medium"
                                    >
                                        Carregar mais serviços
                                    </button>
                                </div>
                            {/if}
                        </div>
                    </div>
                    <div
                        class="px-6 py-4 bg-gray-50 dark:bg-gray-800/50 border-t border-border-light dark:border-border-dark flex items-center justify-end space-x-3"
                    >
                        <button
                            type="button"
                            onclick={() => window.history.back()}
                            class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 font-medium hover:bg-white dark:hover:bg-gray-700 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                        >
                            Cancelar
                        </button>
                        <button
                            type="button"
                            onclick={handleSubmit}
                            disabled={isSubmitting}
                            class="px-6 py-2 bg-primary hover:bg-primary-hover text-white rounded-md font-medium shadow-md transition-colors flex items-center focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50"
                        >
                            {#if isSubmitting}
                                <span
                                    class="material-symbols-outlined animate-spin mr-2 text-[20px]"
                                    >sync</span
                                >
                                Salvando...
                            {:else}
                                <span
                                    class="material-symbols-outlined text-[20px] mr-2"
                                    >save</span
                                >
                                Salvar Profissional
                            {/if}
                        </button>
                    </div>
                </div>

                <div
                    class="bg-blue-50 dark:bg-blue-900/10 border border-blue-100 dark:border-blue-900/30 rounded-lg p-4 flex items-start"
                >
                    <span
                        class="material-symbols-outlined text-blue-500 mt-0.5 mr-3"
                        >info</span
                    >
                    <div>
                        <h4
                            class="text-sm font-medium text-blue-800 dark:text-blue-300"
                        >
                            Informação Importante
                        </h4>
                        <p
                            class="text-sm text-blue-600 dark:text-blue-400 mt-1"
                        >
                            O novo profissional receberá um email com as
                            credenciais de acesso provisórias. Certifique-se de
                            que o email inserido está correto.
                        </p>
                    </div>
                </div>
            </div>
            <footer
                class="mt-12 text-center text-sm text-gray-500 dark:text-gray-400 pb-6"
            >
                <p>© 2023 BellaVita Salon. Todos os direitos reservados.</p>
                <div class="mt-2 space-x-4">
                    <a class="hover:text-primary transition-colors" href="/"
                        >Termos</a
                    >
                    <a class="hover:text-primary transition-colors" href="/"
                        >Privacidade</a
                    >
                </div>
            </footer>
        </div>
    </main>
</div>

<AlertModal
    show={alertState.show}
    title={alertState.title}
    message={alertState.message}
    type={alertState.type}
    onconfirm={alertState.onConfirm}
    oncancel={() => (alertState.show = false)}
/>

<style>
    /* Custom scrollbar to match original style */
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
