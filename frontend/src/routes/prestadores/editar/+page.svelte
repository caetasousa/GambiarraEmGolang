<script lang="ts">
  import Sidebar from "$lib/components/Sidebar.svelte";
  import ThemeToggle from "$lib/components/ThemeToggle.svelte";
  import AlertModal from "$lib/components/AlertModal.svelte";
  import Modal from "$lib/components/Modal.svelte";
  import ImageUpload from "$lib/components/ImageUpload.svelte";
  import {
    uploadPrestadorImageToSupabase,
    deleteImageFromSupabase,
  } from "$lib/utils/imageUpload";
  import { onMount } from "svelte";
    import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
  import { fetchApi } from "$lib/utils/api";

  // Updated Provider Interface to match API (ID as string)
  interface Provider {
    id: string;
    nome: string;
    email: string;
    telefone: string;
    ativo: boolean;
    foto?: string;
    image_url?: string;
    especialidade?: string;
    rating?: number;
    reviews?: number;
    status?: "disponivel" | "ocupado" | "ausente";
    biografia?: string;
    cpf?: string;
    servicos?: string[];
    catalogo?: any[];
  }

  interface CatalogService {
    ID: string;
    Nome: string;
  }

  // State
  let providers: Provider[] = [];
  let availableServices: CatalogService[] = []; // Now storing full objects
  let loadingProviders = true;
  let loadingServices = true;

  // Pagination for providers
  let page = 1;
  let limit = 9; // Changed to 9 as requested
  let totalProviders = 0;
  let totalPages = 1; // Track total pages

  let searchQuery = "";
  let statusFilter: "ativos" | "inativos" = "ativos";

  // --- Modal Editing State ---
  let showEditModal = false;
  let editingProvider: Provider | null = null;
  let isSaving = false;

  // Form Fields
  let editNome = "";
  let editCpf = "";
  let editTelefone = "";
  let editEmail = "";
  let editServiceIds: string[] = []; // Changed to store IDs

  // Image Upload
  let selectedImage: File | null = null;
  let imagePreview: string | null = null;

  // Fetch Providers from API
  async function fetchProviders() {
    loadingProviders = true;
    try {
      const ativoParam = statusFilter === "ativos";
      const res = await fetchApi(
        `/api/v1/prestadores?page=${page}&limit=${limit}&ativo=${ativoParam}`,
      );
      if (res.ok) {
        const data = await res.json();
        // Map API data to our Provider interface
        providers = (data.data || []).map((p: any) => ({
          id: p.ID,
          nome: p.Nome,
          email: p.Email,
          telefone: p.Telefone,
          cpf: p.Cpf,
          ativo: p.Ativo,
          foto:
            p.ImagemUrl ||
            `https://ui-avatars.com/api/?name=${encodeURIComponent(p.Nome)}&background=random`,
          especialidade: "Profissional",
          rating: 5.0,
          reviews: 0,
          status: "disponivel",
          biografia: "Sem biografia.",
          // Map Catalogo to local services array of names for display in card
          servicos: p.Catalogo ? p.Catalogo.map((c: any) => c.Nome) : [],
          catalogo: p.Catalogo || [],
        }));
        totalProviders = data.total || 0;
        // Calculate total pages based on API total
        totalPages = Math.ceil(totalProviders / limit);
      } else {
        console.log("Resposta do servidor:", await res.text());
        console.error("Erro ao buscar prestadores:", await res.text());
        showAlert("Erro", "Falha ao carregar lista de profissionais.", "error");
      }
    } catch (error) {
      console.error("Erro de conexão (prestadores):", error);
      showAlert("Erro", "Erro de conexão ao buscar profissionais.", "error");
    } finally {
      loadingProviders = false;
    }
  }

  // Fetch Services from API
  async function fetchServices() {
    loadingServices = true;
    try {
      const res = await fetch("/api/v1/catalogos?limit=100");
      if (res.ok) {
        const data = await res.json();
        if (data.data && Array.isArray(data.data)) {
          // Store full objects now
          availableServices = data.data.map((s: any) => ({
            ID: s.ID,
            Nome: s.Nome,
          }));
        }
      }
    } catch (error) {
      console.error("Erro de conexão (serviços):", error);
    } finally {
      loadingServices = false;
    }
  }

  onMount(() => {
    fetchProviders();
    fetchServices();
  });

  function openEditModal(provider: Provider) {
    editingProvider = provider;
    editNome = provider.nome;
    editCpf = provider.cpf || "";
    editTelefone = provider.telefone || "";
    editEmail = provider.email || "";

    // Reset service IDs
    editServiceIds = [];

    // Map existing services to IDs
    if (provider.catalogo && Array.isArray(provider.catalogo)) {
      editServiceIds = provider.catalogo.map((c: any) => c.ID);
    }

    // Reset Image
    selectedImage = null;
    imagePreview = provider.foto || null;

    showEditModal = true;
  }

  function handleImageSelect(
    event: CustomEvent<{ file: File; preview: string }>,
  ) {
    selectedImage = event.detail.file;
    imagePreview = event.detail.preview;
  }

  function handleImageClear() {
    selectedImage = null;
    imagePreview = null;
  }

  function toggleService(serviceId: string) {
    if (editServiceIds.includes(serviceId)) {
      editServiceIds = editServiceIds.filter((id) => id !== serviceId);
    } else {
      editServiceIds = [...editServiceIds, serviceId];
    }
  }

  async function handleUpdate() {
    if (!editingProvider) return;

    // --- Validation ---
    if (!editNome || editNome.trim().length < 3) {
      showAlert(
        "Atenção",
        "O nome deve ter no mínimo 3 caracteres.",
        "warning",
      );
      return;
    }
    if (editNome.trim().length > 100) {
      showAlert(
        "Atenção",
        "O nome deve ter no máximo 100 caracteres.",
        "warning",
      );
      return;
    }

    // Email
    if (!editEmail || editEmail.trim() === "") {
      showAlert("Atenção", "O email é obrigatório.", "warning");
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(editEmail)) {
      showAlert("Atenção", "Por favor, insira um email válido.", "warning");
      return;
    }

    // Telefone
    const phoneDigits = editTelefone.replace(/\D/g, "");
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

    // Validação de Serviços (Pelo menos um)
    if (editServiceIds.length === 0) {
      showAlert(
        "Atenção",
        "O profissional deve ter no mínimo um serviço atribuído.",
        "warning",
      );
      return;
    }

    // Photo is strictly strictly required? Registration says yes. Edit might be 'keep existing'.
    // Logic below: use existing URL unless new image uploaded.

    isSaving = true;

    let imageUrl = editingProvider.foto || "";

    // Upload new image if selected
    if (selectedImage) {
      try {
        const uploadedUrl = await uploadPrestadorImageToSupabase(selectedImage);
        if (uploadedUrl) {
          imageUrl = uploadedUrl;
        } else {
          showAlert("Erro", "Erro ao fazer upload da imagem.", "error");
          isSaving = false;
          return;
        }
      } catch (error) {
        console.error(error);
        showAlert("Erro", "Erro ao fazer upload da imagem.", "error");
        isSaving = false;
        return;
      }
    }

    const payload = {
      catalogo_ids: editServiceIds,
      email: editEmail.trim(),
      image_url: imageUrl,
      nome: editNome.trim(),
      telefone: phoneDigits,
    };

    try {
      const res = await fetchApi(`/api/v1/prestadores/${editingProvider.id}`, {
        method: "PUT",
        body: JSON.stringify(payload),
      });

      if (res.ok) {
        showAlert(
          "Sucesso",
          "Profissional atualizado com sucesso!",
          "success",
          false,
          () => {
            // Check if we need to delete the old image
            // Only delete if a new image was uploaded (selectedImage exists)
            // AND there was an old image
            // AND the new URL is different from the old one (to avoid deleting just-uploaded same image)
            const oldImageUrl = editingProvider?.foto;
            if (selectedImage && oldImageUrl && oldImageUrl !== imageUrl) {
              console.log("Deleting old provider image:", oldImageUrl);
              // Fire and forget, don't await/block UI
              deleteImageFromSupabase(oldImageUrl).then((success) => {
                if (!success)
                  console.warn("Failed to delete old image in background");
              });
            }

            alertState.show = false;
            showEditModal = false;
            fetchProviders(); // Refresh list to show changes
          },
        );
      } else {
        const text = await res.text();
        console.error("Erro ao atualizar prestador:", text);
        showAlert("Erro", "Falha ao atualizar profissional.", "error");
      }
    } catch (error) {
      console.error("Erro de conexão ao atualizar:", error);
      showAlert("Erro", "Erro de conexão ao atualizar profissional.", "error");
    } finally {
      isSaving = false;
    }
  }

  function changePage(newPage: number) {
    if (newPage >= 1 && newPage <= totalPages) {
      page = newPage;
      fetchProviders();
    }
  }

  function changeFilter(newFilter: "ativos" | "inativos") {
    statusFilter = newFilter;
    page = 1;
    fetchProviders();
  }

  // --- Alert Modal State ---
  let alertState = {
    show: false,
    title: "",
    message: "",
    type: "info" as "success" | "error" | "warning" | "info",
    showCancel: false,
    confirmText: "OK",
    onConfirm: () => {},
  };

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

  async function handleToggleStatus(provider: Provider) {
    const action = provider.ativo ? "Inativar" : "Ativar";
    const endpoint = provider.ativo ? "inativar" : "ativar";
    const confirmMessage = `Tem certeza que deseja ${action.toLowerCase()} o profissional "${provider.nome}"?`;

    showAlert(
      `${action} Profissional`,
      confirmMessage,
      "warning",
      true,
      async () => {
        try {
          // Dismiss alert
          alertState.show = false;

          const res = await fetchApi(
            `/api/v1/prestadores/${provider.id}/${endpoint}`,
            { method: "PUT" },
          );

          if (res.ok) {
            // Refresh list to adhere to new filter
            fetchProviders();
            const successTerm = provider.ativo ? "Inativado" : "Ativado";
            showAlert(
              "Sucesso",
              `Profissional ${successTerm} com sucesso!`,
              "success",
            );
          } else {
            const text = await res.text();
            console.error(`Erro ao ${action.toLowerCase()} prestador:`, text);
            showAlert(
              "Erro",
              `Falha ao ${action.toLowerCase()} profissional.`,
              "error",
            );
          }
        } catch (error) {
          console.error(`Erro de conexão ao ${action.toLowerCase()}:`, error);
          showAlert(
            "Erro",
            `Erro de conexão ao ${action.toLowerCase()} profissional.`,
            "error",
          );
        }
      },
    );
  }
</script>

<div
  class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
  <Sidebar />

  <main class="flex-1 flex flex-col h-full overflow-hidden relative">
    <DashboardNavbar />

    <div class="flex-1 overflow-y-auto p-6 md:p-8">
      <div class="max-w-6xl mx-auto space-y-8">
        <div>
          <nav class="flex text-sm text-gray-500 dark:text-gray-400 mb-2">
            <a href="/" class="hover:text-primary transition-colors">Início</a>
            <span class="mx-2">/</span>
            <span class="text-gray-900 dark:text-white font-medium"
              >Equipe de Profissionais</span
            >
          </nav>
          <div
            class="flex flex-col md:flex-row md:items-center justify-between gap-4"
          >
            <div>
              <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                Gerenciar Profissionais
              </h1>
              <p class="text-gray-500 dark:text-gray-400 mt-1">
                Adicione, edite ou remova membros da equipe.
              </p>
            </div>

            <div class="flex flex-col sm:flex-row items-center gap-3">
              <div
                class="flex items-center space-x-1 bg-gray-100 dark:bg-[hsl(var(--bs-muted))]/20 p-1 rounded-xl"
              >
                <button
                  on:click={() => changeFilter("ativos")}
                  class="px-4 py-1.5 text-sm font-semibold rounded-lg transition-all {statusFilter ===
                  'ativos'
                    ? 'bg-white dark:bg-gray-700 shadow-md text-brand-orange'
                    : 'text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'}"
                >
                  Ativos
                </button>
                <button
                  on:click={() => changeFilter("inativos")}
                  class="px-4 py-1.5 text-sm font-semibold rounded-lg transition-all {statusFilter ===
                  'inativos'
                    ? 'bg-white dark:bg-gray-700 shadow-md text-brand-orange'
                    : 'text-gray-500 hover:text-gray-700 dark:hover:text-gray-300'}"
                >
                  Inativos
                </button>
              </div>

              <a
                href="/prestadores/cadastro"
                class="w-full sm:w-auto px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-xl font-bold shadow-lg shadow-brand-orange/20 transition-all active:scale-95 flex items-center justify-center"
              >
                <span class="material-symbols-outlined text-xl mr-2"
                  >person_add</span
                >
                Novo Profissional
              </a>
            </div>
          </div>
        </div>

        {#if loadingProviders}
          <div class="flex justify-center items-center h-64">
            <span
              class="material-symbols-outlined animate-spin text-4xl text-primary"
              >sync</span
            >
            <span class="ml-2 text-gray-500 dark:text-gray-400"
              >Carregando profissionais...</span
            >
          </div>
        {:else if providers.length === 0}
          <div
            class="bg-[hsl(var(--bs-card))] rounded-xl border-2 border-dashed border-border-light dark:border-border-dark p-12 text-center"
          >
            <span class="material-symbols-outlined text-5xl text-gray-300 mb-4"
              >person_search</span
            >
            <p class="text-gray-500">Nenhum profissional encontrado.</p>
          </div>
        {:else}
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each providers as provider}
              <div
                class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark overflow-hidden hover:shadow-xl transition-all duration-300 group flex flex-col"
              >
                <div class="p-6 flex flex-col items-center text-center">
                  <div class="relative mb-4">
                    <img
                      src={provider.foto}
                      alt={provider.nome}
                      class="h-24 w-24 rounded-full object-cover border-4 border-white dark:border-gray-700 shadow-sm"
                    />
                    <span
                      class="absolute bottom-1 right-1 w-5 h-5 rounded-full border-2 border-white dark:border-gray-800 {provider.ativo
                        ? 'bg-green-500'
                        : 'bg-gray-400'}"
                      title={provider.ativo ? "Ativo" : "Inativo"}
                    ></span>
                  </div>
                  <h3
                    class="text-lg font-bold text-gray-900 dark:text-white group-hover:text-brand-orange transition-colors"
                  >
                    {provider.nome}
                  </h3>
                  <p class="text-sm text-brand-orange font-medium mt-1">
                    {provider.especialidade}
                  </p>

                  <div class="flex items-center mt-3 text-sm text-amber-500">
                    <span class="material-icons text-sm mr-1">star</span>
                    <span class="font-bold">{provider.rating}</span>
                    <span class="text-gray-400 ml-1"
                      >({provider.reviews} avaliações)</span
                    >
                  </div>

                  <p
                    class="text-sm text-gray-500 dark:text-gray-400 mt-4 line-clamp-2 italic"
                  >
                    "{provider.biografia}"
                  </p>
                </div>

                <div
                  class="mt-auto px-6 py-4 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 border-t border-border-light dark:border-border-dark flex items-center justify-between"
                >
                  <div class="flex items-center gap-2">
                    <span
                      class="flex h-2 w-2 rounded-full {provider.ativo
                        ? 'bg-green-500'
                        : 'bg-red-500'}"
                    ></span>
                    <span class="text-xs text-gray-400 uppercase tracking-wider"
                      >{provider.ativo ? "Ativo" : "Inativo"}</span
                    >
                  </div>
                  <div class="flex items-center gap-2">
                    <a
                      href={`/prestadores/${provider.id}`}
                      class="p-2 text-gray-400 hover:text-blue-500 dark:hover:text-blue-500 transition-colors rounded-lg hover:bg-blue-500/5"
                      title="Ver Detalhes"
                    >
                      <span class="material-symbols-outlined text-[20px]"
                        >visibility</span
                      >
                    </a>
                    <button
                      class="p-2 text-gray-400 hover:text-brand-orange dark:hover:text-brand-orange transition-colors rounded-lg hover:bg-brand-orange/5"
                      title="Editar"
                      on:click={() => openEditModal(provider)}
                    >
                      <span class="material-symbols-outlined text-[20px]"
                        >edit</span
                      >
                    </button>
                    <button
                      class="p-2 text-gray-400 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 {provider.ativo
                        ? 'hover:text-red-500'
                        : 'hover:text-green-500'}"
                      title={provider.ativo ? "Inativar" : "Ativar"}
                      on:click={() => handleToggleStatus(provider)}
                    >
                      <span class="material-symbols-outlined text-[20px]"
                        >{provider.ativo ? "block" : "check_circle"}</span
                      >
                    </button>
                  </div>
                </div>
              </div>
            {/each}
          </div>

          <!-- Pagination Controls -->
          <div
            class="flex items-center justify-between border-t border-gray-200 dark:border-border-dark bg-[hsl(var(--bs-card))] px-4 py-3 sm:px-6 mt-6 rounded-lg"
          >
            <div class="flex flex-1 justify-between sm:hidden">
              <button
                on:click={() => changePage(page - 1)}
                disabled={page === 1}
                class="relative inline-flex items-center rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50"
              >
                Anterior
              </button>
              <button
                on:click={() => changePage(page + 1)}
                disabled={page === totalPages}
                class="relative ml-3 inline-flex items-center rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 disabled:opacity-50"
              >
                Próximo
              </button>
            </div>
            <div
              class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between"
            >
              <div>
                <p class="text-sm text-gray-700 dark:text-gray-300">
                  Mostrando
                  <span class="font-medium">{(page - 1) * limit + 1}</span>
                  a
                  <span class="font-medium"
                    >{Math.min(page * limit, totalProviders)}</span
                  >
                  de
                  <span class="font-medium">{totalProviders}</span>
                  resultados
                </p>
              </div>
              <div>
                <nav
                  class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                  aria-label="Pagination"
                >
                  <button
                    on:click={() => changePage(page - 1)}
                    disabled={page === 1}
                    class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 dark:ring-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                  >
                    <span class="sr-only">Anterior</span>
                    <span class="material-symbols-outlined text-sm"
                      >chevron_left</span
                    >
                  </button>

                  <!-- Simple page numbers logic -->
                  {#each Array(totalPages) as _, i}
                    <button
                      on:click={() => changePage(i + 1)}
                      class="relative inline-flex items-center px-4 py-2 text-sm font-semibold {page ===
                      i + 1
                        ? 'z-10 bg-brand-orange text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-brand-orange'
                        : 'text-gray-900 dark:text-gray-300 ring-1 ring-inset ring-gray-300 dark:ring-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 focus:z-20 focus:outline-offset-0'}"
                    >
                      {i + 1}
                    </button>
                  {/each}

                  <button
                    on:click={() => changePage(page + 1)}
                    disabled={page === totalPages}
                    class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 dark:ring-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 focus:z-20 focus:outline-offset-0 disabled:opacity-50"
                  >
                    <span class="sr-only">Próximo</span>
                    <span class="material-symbols-outlined text-sm"
                      >chevron_right</span
                    >
                  </button>
                </nav>
              </div>
            </div>
          </div>
        {/if}

        <footer
          class="mt-12 text-center text-sm text-gray-500 dark:text-gray-400 pb-6"
        >
          <p>© 2023 BellaVita Salon. Todos os direitos reservados.</p>
        </footer>
      </div>
    </div>
  </main>

  <AlertModal
    show={alertState.show}
    title={alertState.title}
    message={alertState.message}
    type={alertState.type}
    showCancel={alertState.showCancel}
    confirmText={alertState.confirmText}
    on:confirm={alertState.onConfirm}
    on:cancel={closeAlert}
  />

  <Modal
    show={showEditModal}
    title="Editar Profissional"
    maxWidth="max-w-2xl"
    on:close={() => (showEditModal = false)}
  >
    <div slot="body" class="space-y-4">
      <div>
        <span
          class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
          >Foto do Profissional</span
        >
        <ImageUpload
          previewUrl={imagePreview}
          on:select={handleImageSelect}
          on:clear={handleImageClear}
        />
      </div>

      <div>
        <label
          for="edit-nome"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
          >Nome Completo</label
        >
        <input
          id="edit-nome"
          type="text"
          bind:value={editNome}
          class="block w-full px-3 py-2 border border-border-light dark:border-border-dark rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
          placeholder="Ex: Ana Maria Costa"
        />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label
            for="edit-cpf"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >CPF</label
          >
          <input
            id="edit-cpf"
            type="text"
            bind:value={editCpf}
            disabled
            class="block w-full px-3 py-2 border border-border-light dark:border-border-dark rounded-md bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 focus:ring-brand-orange focus:border-brand-orange sm:text-sm cursor-not-allowed"
            placeholder="000.000.000-00"
          />
        </div>
        <div>
          <label
            for="edit-telefone"
            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
            >Telefone / WhatsApp</label
          >
          <input
            id="edit-telefone"
            type="text"
            bind:value={editTelefone}
            class="block w-full px-3 py-2 border border-border-light dark:border-border-dark rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
            placeholder="(00) 00000-0000"
          />
        </div>
      </div>

      <div>
        <label
          for="edit-email"
          class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
          >Endereço de Email</label
        >
        <input
          id="edit-email"
          type="email"
          bind:value={editEmail}
          class="block w-full px-3 py-2 border border-border-light dark:border-border-dark rounded-md bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-brand-orange focus:border-brand-orange sm:text-sm"
          placeholder="profissional@bellavita.com"
        />
      </div>

      <div>
        <h3
          class="flex items-center text-sm font-bold text-gray-900 dark:text-white mb-3"
        >
          <span class="material-symbols-outlined text-brand-orange mr-2"
            >spa</span
          >
          Serviços Atribuídos
        </h3>
        <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
          Selecione todos os serviços que este profissional está qualificado
          para realizar.
        </p>

        <div
          class="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-40 overflow-y-auto p-1"
        >
          {#each availableServices as service}
            <button
              type="button"
              class="flex items-center justify-between p-3 rounded-lg border transition-all text-left group
                        {editServiceIds.includes(service.ID)
                ? 'bg-brand-orange/10 border-brand-orange text-brand-orange'
                : 'bg-gray-50 dark:bg-gray-800 border-border-light dark:border-border-dark text-gray-700 dark:text-gray-300 hover:border-brand-orange/50'}"
              on:click={() => toggleService(service.ID)}
            >
              <span class="text-sm font-medium">{service.Nome}</span>
              {#if editServiceIds.includes(service.ID)}
                <span class="material-symbols-outlined text-edit-check text-lg"
                  >check_circle</span
                >
              {:else}
                <span
                  class="material-symbols-outlined text-gray-300 group-hover:text-gray-400 text-lg"
                  >radio_button_unchecked</span
                >
              {/if}
            </button>
          {/each}
        </div>
      </div>
    </div>

    <div slot="footer">
      <button
        on:click={() => (showEditModal = false)}
        class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
      >
        Cancelar
      </button>
      <button
        on:click={handleUpdate}
        disabled={isSaving}
        class="px-4 py-2 text-sm font-medium text-white bg-brand-orange rounded-md hover:bg-brand-orange-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-brand-orange disabled:opacity-50 transition-colors"
      >
        {isSaving ? "Salvando..." : "Salvar Alterações"}
      </button>
    </div>
  </Modal>
</div>
