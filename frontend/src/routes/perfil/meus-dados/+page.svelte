<script lang="ts">
    import Sidebar from "$lib/components/Sidebar.svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
    import ServiceSelectionModal from "$lib/components/profile/ServiceSelectionModal.svelte";
    import AlertModal from "$lib/components/AlertModal.svelte";
    import ProfileDataAndServicesTab from "$lib/components/profile/ProfileDataAndServicesTab.svelte";
    import { fetchApi } from "$lib/utils/api";
    import { user } from "$lib/stores/auth";
    
    import type { 
        Service, 
        ProviderProfile, 
        ProviderData 
    } from '$lib/types/profile';

    import { uploadPrestadorImageToSupabase, deleteImageFromSupabase } from "$lib/utils/imageUpload";

    // --- User Data ---
    $: userId = $user?.id ?? "";
    $: userRole = $user?.role ?? "";
    $: isPrestadorOrAdmin = userRole === "prestador" || userRole === "admin";
    let currentProviderData: ProviderData | null = null;

    // --- Services Logic ---
    let services: any[] = [];
    let servicesLoading = false;
    let showSelectionModal = false;

    // --- Edit Mode Logic ---
    let isEditing = false;
    let selectedFile: File | null = null;

    // --- Form Data ---
    let profile: ProviderProfile = {
        name: "",
        cpf: "",
        phone: "",
        email: "",
        avatarUrl: "",
    };

    async function fetchUserData() {
        try {
            if (isPrestadorOrAdmin) {
                const response = await fetchApi(`/api/v1/prestadores/${userId}`);
                if (response.ok) {
                    const data: ProviderData = await response.json();
                    currentProviderData = data;
                    profile = {
                        name: data.nome,
                        email: data.email,
                        phone: data.telefone,
                        cpf: data.cpf,
                        avatarUrl: data.image_url,
                        ativo: data.ativo,
                        role: userRole,
                    };
                    const catalogData = data.catalogo;
                    if (catalogData && Array.isArray(catalogData)) {
                        services = catalogData.map((s: Service) => ({
                            ID: s.id || (s as any).ID,
                            Nome: s.nome || (s as any).Nome,
                            Preco: s.preco || (s as any).Preco,
                            DuracaoPadrao: s.duracao_padrao || (s as any).DuracaoPadrao,
                            Categoria: s.categoria || (s as any).Categoria,
                            ImagemUrl: s.image_url || (s as any).ImagemUrl || ""
                        }));
                    }
                } else {
                    console.error("Failed to fetch prestador data");
                }
            } else {
                // admin ou cliente
                const response = await fetchApi(`/api/v1/clientes/${userId}`);
                if (response.ok) {
                    const data = await response.json();
                    profile = {
                        name: data.nome,
                        email: data.email,
                        phone: data.telefone,
                        cpf: "",
                        avatarUrl: "",
                        role: data.role,
                    };
                } else {
                    console.error("Failed to fetch cliente/admin data");
                }
            }
        } catch (error) {
            console.error("Error fetching user data:", error);
        }
    }

    $: if (userId && userRole) fetchUserData();

    function toggleEdit() {
        isEditing = true;
    }

    function cancelEdit() {
        isEditing = false;
        selectedFile = null;
        fetchUserData(); // Reset data
    }

    function handleImageSelect(event: CustomEvent<{ file: File; preview: string }>) {
        selectedFile = event.detail.file;
        profile.avatarUrl = event.detail.preview;
    }

    function handleImageClear() {
        selectedFile = null;
        profile.avatarUrl = currentProviderData?.image_url || "";
    }

    function openSelectionModal() {
        showSelectionModal = true;
    }

    async function saveProviderCatalog(updatedServices: any[]) {
        if (!currentProviderData) return;

        const catalogIds = updatedServices.map(s => String(s.ID));

        const updatedProvider: ProviderData = {
            ...currentProviderData,
            catalogo_ids: catalogIds,
            catalogo: currentProviderData.catalogo || [],
            agenda: currentProviderData.agenda || []
        };

        try {
            const response = await fetchApi(`/api/v1/prestadores/${userId}`, {
                method: "PUT",
                body: JSON.stringify(updatedProvider),
            });

            if (response.ok) {
                services = updatedServices;
                currentProviderData = updatedProvider;
            } else {
                const errorText = await response.text();
                console.error("Failed to update catalog", response.status, errorText);
                alert(`Erro ao atualizar catálogo (${response.status}): ${errorText}`);
            }
        } catch (e) {
            console.error("Error saving catalog:", e);
            alert("Erro de conexão.");
        }
    }

    async function saveProfileData() {
        if (!currentProviderData) return;

        try {
            let imageUrl = profile.avatarUrl;

            if (selectedFile) {
                const uploadedUrl = await uploadPrestadorImageToSupabase(selectedFile);
                if (uploadedUrl) {
                    imageUrl = uploadedUrl;
                } else {
                    throw new Error("Falha ao fazer upload da imagem.");
                }
            }

            const catalogIds = services.map((s) => String(s.ID));

            const updatedProvider: ProviderData = {
                ...currentProviderData,
                catalogo_ids: catalogIds,
                email: profile.email,
                image_url: imageUrl,
                nome: profile.name,
                telefone: profile.phone,
                cpf: profile.cpf,
                catalogo: currentProviderData.catalogo || [],
                agenda: currentProviderData.agenda || []
            };

            const response = await fetchApi(`/api/v1/prestadores/${userId}`, {
                method: "PUT",
                body: JSON.stringify(updatedProvider),
            });

            if (response.ok) {
                if (selectedFile && currentProviderData.image_url && currentProviderData.image_url !== imageUrl) {
                    deleteImageFromSupabase(currentProviderData.image_url).then(success => {
                        if (success) console.log("Old image deleted successfully");
                        else console.warn("Failed to delete old image");
                    });
                }

                currentProviderData = { ...currentProviderData, ...updatedProvider };
                isEditing = false;
                
                alertConfig = {
                    title: "Sucesso",
                    message: "Dados atualizados com sucesso!",
                    type: "success",
                    confirmText: "OK",
                    showCancel: false,
                };
                onConfirmAction = () => { showAlert = false; };
                showAlert = true;
            } else {
                const errorText = await response.text();
                console.error("Failed to update profile", response.status, errorText);
                
                alertConfig = {
                    title: "Erro",
                    message: `Erro ao atualizar dados: ${errorText}`,
                    type: "error",
                    confirmText: "Fechar",
                    showCancel: false,
                };
                showAlert = true;
            }
        } catch (e) {
            console.error("Error saving profile:", e);
            alertConfig = {
                title: "Erro",
                message: "Erro de conexão ao tentar salvar.",
                type: "error",
                confirmText: "Fechar",
                showCancel: false,
            };
            showAlert = true;
        }
    }

    // --- Alert Modal Logic ---
    let showAlert = false;
    let alertConfig = {
        title: "",
        message: "",
        type: "info" as "info" | "success" | "warning" | "error" | "primary",
        confirmText: "Confirmar",
        showCancel: true,
    };
    let onConfirmAction: () => void = () => {};

    function handleConfirmAction() {
        if (onConfirmAction) onConfirmAction();
        showAlert = false;
    }

    function handleAddService(event: CustomEvent) {
        const newService = event.detail;
        if (!services.find((s) => s.ID === (newService.ID || newService.id))) {
            alertConfig = {
                title: "Confirmar Adição",
                message: `Deseja adicionar o serviço "${newService.Nome || newService.nome}" aos seus serviços?`,
                type: "primary",
                confirmText: "Adicionar",
                showCancel: true,
            };
            onConfirmAction = () => {
                const serviceToAdd = {
                    ID: newService.ID || newService.id,
                    Nome: newService.Nome || newService.nome,
                    Preco: newService.Preco || newService.preco,
                    DuracaoPadrao:
                        newService.DuracaoPadrao || newService.duracao_padrao,
                    Categoria: newService.Categoria || newService.categoria,
                    ImagemUrl: newService.ImagemUrl || newService.image_url,
                };
                const updatedServices = [...services, serviceToAdd];
                saveProviderCatalog(updatedServices);
            };
            showAlert = true;
        }
    }

    function handleRemoveService(event: CustomEvent<{ service: any}>) {
        const service = event.detail.service;
        alertConfig = {
            title: "Remover Serviço",
            message: `Tem certeza que deseja remover "${service.Nome}"?`,
            type: "error",
            confirmText: "Remover",
            showCancel: true,
        };
        onConfirmAction = () => {
            const updatedServices = services.filter((s) => s.ID !== service.ID);
            saveProviderCatalog(updatedServices);
        };
        showAlert = true;
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
                <!-- Page Header -->
                <div class="mb-6">
                    <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                        Meus Dados
                    </h1>
                    <p class="text-gray-500 dark:text-gray-400 mt-2">
                        Gerencie suas informações pessoais e serviços oferecidos
                    </p>
                </div>

                <!-- Content -->
                <ProfileDataAndServicesTab 
                    {profile} 
                    {isEditing}
                    {services}
                    {servicesLoading}
                    on:toggleEdit={toggleEdit}
                    on:cancelEdit={cancelEdit}
                    on:save={saveProfileData}
                    on:imageSelect={handleImageSelect}
                    on:imageClear={handleImageClear}
                    on:addService={openSelectionModal}
                    on:removeService={handleRemoveService}
                />
            </div>
        </div>

        <!-- Service Selection Modal -->
        <ServiceSelectionModal
            bind:show={showSelectionModal}
            providerId={userId}
            existingServiceIds={services.map(s => s.ID)}
            on:add={handleAddService}
        />

        <AlertModal
            bind:show={showAlert}
            title={alertConfig.title}
            message={alertConfig.message}
            type={alertConfig.type}
            confirmText={alertConfig.confirmText}
            showCancel={alertConfig.showCancel}
            on:confirm={handleConfirmAction}
        />
    </main>
</div>
