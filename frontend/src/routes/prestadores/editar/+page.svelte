<script lang="ts">
  import Sidebar from "$lib/components/Sidebar.svelte";
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
  let providers = $state<Provider[]>([]);
  let availableServices = $state<CatalogService[]>([]);
  let loadingProviders = $state(true);
  let loadingServices = $state(true);

  // Pagination for providers
  let page = $state(1);
  let limit = $state(9);
  let totalProviders = $state(0);
  let totalPages = $state(1);

  let searchQuery = $state("");
  let statusFilter = $state<"ativos" | "inativos">("ativos");

  // --- Modal Editing State ---
  let showEditModal = $state(false);
  let editingProvider = $state<Provider | null>(null);
  let isSaving = $state(false);

  // Form Fields
  let editNome = $state("");
  let editCpf = $state("");
  let editTelefone = $state("");
  let editEmail = $state("");
  let editServiceIds = $state<string[]>([]);

  // Image Upload
  let selectedImage = $state<File | null>(null);
  let imagePreview = $state<string | null>(null);

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
          servicos: p.Catalogo ? p.Catalogo.map((c: any) => c.Nome) : [],
          catalogo: p.Catalogo || [],
        }));
        totalProviders = data.total || 0;
        totalPages = Math.ceil(totalProviders / limit);
      } else {
        showAlert("Erro", "Falha ao carregar lista de profissionais.", "error");
      }
    } catch (error) {
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

    editServiceIds = [];

    if (provider.catalogo && Array.isArray(provider.catalogo)) {
      editServiceIds = provider.catalogo.map((c: any) => c.ID);
    }

    selectedImage = null;
    imagePreview = provider.foto || null;

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

  function toggleService(serviceId: string) {
    if (editServiceIds.includes(serviceId)) {
      editServiceIds = editServiceIds.filter((id) => id !== serviceId);
    } else {
      editServiceIds = [...editServiceIds, serviceId];
    }
  }

  async function handleUpdate() {
    if (!editingProvider) return;

    if (!editNome || editNome.trim().length < 3) {
      showAlert("Atenção", "O nome deve ter no mínimo 3 caracteres.", "warning");
      return;
    }
    if (editNome.trim().length > 100) {
      showAlert("Atenção", "O nome deve ter no máximo 100 caracteres.", "warning");
      return;
    }
    if (!editEmail || editEmail.trim() === "") {
      showAlert("Atenção", "O email é obrigatório.", "warning");
      return;
    }
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(editEmail)) {
      showAlert("Atenção", "Por favor, insira um email válido.", "warning");
      return;
    }
    const phoneDigits = editTelefone.replace(/\D/g, "");
    if (!phoneDigits || phoneDigits.length < 8) {
      showAlert("Atenção", "O telefone deve ter no mínimo 8 dígitos.", "warning");
      return;
    }
    if (phoneDigits.length > 15) {
      showAlert("Atenção", "O telefone deve ter no máximo 15 dígitos.", "warning");
      return;
    }
    if (editServiceIds.length === 0) {
      showAlert("Atenção", "O profissional deve ter no mínimo um serviço atribuído.", "warning");
      return;
    }

    isSaving = true;

    let imageUrl = editingProvider.foto || "";

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
            const oldImageUrl = editingProvider?.foto;
            if (selectedImage && oldImageUrl && oldImageUrl !== imageUrl) {
              deleteImageFromSupabase(oldImageUrl).then((success) => {
                if (!success)
                  console.warn("Failed to delete old image in background");
              });
            }
            alertState.show = false;
            showEditModal = false;
            fetchProviders();
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
          alertState.show = false;
          const res = await fetchApi(
            `/api/v1/prestadores/${provider.id}/${endpoint}`,
            { method: "PUT" },
          );
          if (res.ok) {
            fetchProviders();
            const successTerm = provider.ativo ? "Inativado" : "Ativado";
            showAlert("Sucesso", `Profissional ${successTerm} com sucesso!`, "success");
          } else {
            showAlert("Erro", `Falha ao ${action.toLowerCase()} profissional.`, "error");
          }
        } catch (error) {
          showAlert("Erro", `Erro de conexão ao ${action.toLowerCase()} profissional.`, "error");
        }
      },
    );
  }
</script>

<div class="font-sans bg-slate-50 dark:bg-gray-950 text-slate-800 dark:text-slate-100 antialiased h-screen flex overflow-hidden">
  <Sidebar />

  <main class="flex-1 flex flex-col h-full overflow-hidden relative">
    <DashboardNavbar title="Gerenciar Profissionais" />

    <div class="flex-1 overflow-y-auto p-6 md:p-8">
      <div class="max-w-6xl mx-auto space-y-6">

        <!-- Page subheader + actions -->
        <div class="flex flex-col md:flex-row md:items-center justify-between gap-4">
          <div>
            <nav class="flex items-center gap-1.5 text-xs text-slate-400 mb-1">
              <a href="/" class="hover:text-orange-500 transition-colors cursor-pointer">Início</a>
              <span class="material-symbols-outlined text-[12px]">chevron_right</span>
              <span class="text-slate-600 dark:text-slate-300 font-medium">Equipe de Profissionais</span>
            </nav>
            <p class="text-sm text-slate-400">Adicione, edite ou gerencie membros da equipe.</p>
          </div>

          <div class="flex flex-col sm:flex-row items-center gap-3">
            <!-- Status filter toggle -->
            <div class="flex items-center gap-1 bg-white dark:bg-gray-900 border border-gray-100 dark:border-gray-800 shadow-sm p-1 rounded-xl">
              <button
                onclick={() => changeFilter("ativos")}
                class="px-4 py-1.5 text-sm font-semibold rounded-lg transition-all duration-200 cursor-pointer
                  {statusFilter === 'ativos'
                    ? 'bg-orange-500 text-white shadow-sm'
                    : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}"
              >
                Ativos
              </button>
              <button
                onclick={() => changeFilter("inativos")}
                class="px-4 py-1.5 text-sm font-semibold rounded-lg transition-all duration-200 cursor-pointer
                  {statusFilter === 'inativos'
                    ? 'bg-orange-500 text-white shadow-sm'
                    : 'text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200'}"
              >
                Inativos
              </button>
            </div>

            <a
              href="/prestadores/cadastro"
              class="w-full sm:w-auto px-5 py-2.5 bg-orange-500 hover:bg-orange-600 text-white rounded-xl font-semibold shadow-md shadow-orange-500/20 transition-all duration-200 active:scale-95 flex items-center justify-center gap-2 cursor-pointer"
            >
              <span class="material-symbols-outlined text-[18px]">person_add</span>
              Novo Profissional
            </a>
          </div>
        </div>

        {#if loadingProviders}
          <div class="flex flex-col justify-center items-center h-64 gap-4">
            <div class="w-14 h-14 rounded-2xl bg-orange-50 flex items-center justify-center">
              <span class="material-symbols-outlined animate-spin text-orange-500 text-3xl">sync</span>
            </div>
            <p class="text-sm text-slate-400 font-medium">Carregando profissionais...</p>
          </div>

        {:else if providers.length === 0}
          <div class="bg-white dark:bg-gray-900 rounded-2xl border-2 border-dashed border-slate-200 dark:border-slate-700 p-14 text-center">
            <div class="w-16 h-16 bg-slate-50 dark:bg-slate-800 rounded-2xl flex items-center justify-center mx-auto mb-4">
              <span class="material-symbols-outlined text-slate-300 text-4xl">person_search</span>
            </div>
            <p class="font-semibold text-slate-600 dark:text-slate-300">Nenhum profissional encontrado</p>
            <p class="text-sm text-slate-400 mt-1">
              {statusFilter === 'ativos' ? 'Não há profissionais ativos.' : 'Não há profissionais inativos.'}
            </p>
          </div>

        {:else}
          <!-- Provider cards grid -->
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
            {#each providers as provider}
              <div class="bg-white dark:bg-gray-900 rounded-2xl border border-gray-100 dark:border-gray-800 overflow-hidden hover:shadow-lg transition-all duration-300 group flex flex-col shadow-sm">
                <!-- Card body -->
                <div class="p-6 flex flex-col items-center text-center flex-1">
                  <!-- Avatar with status indicator -->
                  <div class="relative mb-4">
                    <img
                      src={provider.foto}
                      alt={provider.nome}
                      class="h-20 w-20 rounded-2xl object-cover ring-4 ring-white dark:ring-gray-900 shadow-md group-hover:shadow-lg transition-all duration-300"
                    />
                    <span
                      class="absolute -bottom-0.5 -right-0.5 w-4 h-4 rounded-full border-2 border-white dark:border-gray-900 shadow-sm {provider.ativo ? 'bg-emerald-500' : 'bg-slate-300'}"
                      title={provider.ativo ? "Ativo" : "Inativo"}
                    ></span>
                  </div>

                  <h3
                    class="text-base font-bold text-slate-900 dark:text-white group-hover:text-orange-500 transition-colors duration-200 leading-tight"
                    style="font-family: 'Cormorant', serif; font-size: 1.1rem;"
                  >
                    {provider.nome}
                  </h3>
                  <p class="text-xs text-orange-500 font-semibold mt-1 uppercase tracking-wider">
                    {provider.especialidade}
                  </p>

                  <!-- Rating row -->
                  <div class="flex items-center gap-1 mt-2.5 text-sm">
                    <span class="material-icons text-amber-400 text-sm">star</span>
                    <span class="font-bold text-slate-700 dark:text-slate-200">{provider.rating}</span>
                    <span class="text-slate-400 text-xs">({provider.reviews} avaliações)</span>
                  </div>

                  <!-- Services tags -->
                  {#if provider.servicos && provider.servicos.length > 0}
                    <div class="flex flex-wrap gap-1.5 justify-center mt-3">
                      {#each provider.servicos.slice(0, 3) as servico}
                        <span class="text-xs px-2 py-0.5 bg-orange-50 text-orange-600 rounded-full font-medium border border-orange-100">
                          {servico}
                        </span>
                      {/each}
                      {#if provider.servicos.length > 3}
                        <span class="text-xs px-2 py-0.5 bg-slate-50 dark:bg-slate-800 text-slate-400 rounded-full font-medium">
                          +{provider.servicos.length - 3}
                        </span>
                      {/if}
                    </div>
                  {/if}
                </div>

                <!-- Card footer -->
                <div class="px-5 py-3.5 bg-slate-50 dark:bg-gray-800 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between">
                  <div class="flex items-center gap-1.5">
                    <span class="w-2 h-2 rounded-full flex-shrink-0 {provider.ativo ? 'bg-emerald-500' : 'bg-slate-300'}"></span>
                    <span class="text-xs font-semibold text-slate-500 dark:text-slate-400 uppercase tracking-wider">
                      {provider.ativo ? "Ativo" : "Inativo"}
                    </span>
                  </div>
                  <div class="flex items-center gap-1">
                    <a
                      href={`/prestadores/${provider.id}`}
                      class="p-2 text-slate-400 hover:text-blue-500 transition-all duration-200 rounded-xl hover:bg-blue-50 cursor-pointer"
                      title="Ver Detalhes"
                    >
                      <span class="material-symbols-outlined text-[18px]">visibility</span>
                    </a>
                    <button
                      class="p-2 text-slate-400 hover:text-orange-500 transition-all duration-200 rounded-xl hover:bg-orange-50 cursor-pointer"
                      title="Editar"
                      onclick={() => openEditModal(provider)}
                    >
                      <span class="material-symbols-outlined text-[18px]">edit</span>
                    </button>
                    <button
                      class="p-2 text-slate-400 transition-all duration-200 rounded-xl cursor-pointer
                        {provider.ativo
                          ? 'hover:text-red-500 hover:bg-red-50'
                          : 'hover:text-emerald-500 hover:bg-emerald-50'}"
                      title={provider.ativo ? "Inativar" : "Ativar"}
                      onclick={() => handleToggleStatus(provider)}
                    >
                      <span class="material-symbols-outlined text-[18px]">
                        {provider.ativo ? "block" : "check_circle"}
                      </span>
                    </button>
                  </div>
                </div>
              </div>
            {/each}
          </div>

          <!-- Pagination Controls -->
          <div class="flex items-center justify-between bg-white dark:bg-gray-900 px-5 py-3.5 rounded-2xl border border-gray-100 dark:border-gray-800 shadow-sm mt-2">
            <!-- Mobile -->
            <div class="flex flex-1 justify-between sm:hidden">
              <button
                onclick={() => changePage(page - 1)}
                disabled={page === 1}
                class="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-semibold rounded-xl border border-gray-200 dark:border-gray-700 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 disabled:opacity-40 transition-all cursor-pointer"
              >
                <span class="material-symbols-outlined text-[16px]">chevron_left</span>
                Anterior
              </button>
              <button
                onclick={() => changePage(page + 1)}
                disabled={page === totalPages}
                class="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-semibold rounded-xl border border-gray-200 dark:border-gray-700 text-slate-600 dark:text-slate-300 hover:bg-slate-50 dark:hover:bg-slate-800 disabled:opacity-40 transition-all cursor-pointer"
              >
                Próximo
                <span class="material-symbols-outlined text-[16px]">chevron_right</span>
              </button>
            </div>

            <!-- Desktop -->
            <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
              <p class="text-sm text-slate-500 dark:text-slate-400">
                Mostrando
                <span class="font-semibold text-slate-700 dark:text-slate-200">{(page - 1) * limit + 1}</span>
                a
                <span class="font-semibold text-slate-700 dark:text-slate-200">{Math.min(page * limit, totalProviders)}</span>
                de
                <span class="font-semibold text-slate-700 dark:text-slate-200">{totalProviders}</span>
                resultados
              </p>
              <nav class="flex items-center gap-1" aria-label="Pagination">
                <button
                  onclick={() => changePage(page - 1)}
                  disabled={page === 1}
                  class="p-2 rounded-xl text-slate-400 hover:text-orange-500 hover:bg-orange-50 disabled:opacity-40 transition-all duration-200 cursor-pointer"
                >
                  <span class="material-symbols-outlined text-[18px]">chevron_left</span>
                </button>

                {#each Array(totalPages) as _, i}
                  <button
                    onclick={() => changePage(i + 1)}
                    class="w-9 h-9 rounded-xl text-sm font-semibold transition-all duration-200 cursor-pointer
                      {page === i + 1
                        ? 'bg-orange-500 text-white shadow-sm shadow-orange-500/25'
                        : 'text-slate-600 dark:text-slate-300 hover:bg-orange-50 dark:hover:bg-orange-900/20 hover:text-orange-600'}"
                  >
                    {i + 1}
                  </button>
                {/each}

                <button
                  onclick={() => changePage(page + 1)}
                  disabled={page === totalPages}
                  class="p-2 rounded-xl text-slate-400 hover:text-orange-500 hover:bg-orange-50 disabled:opacity-40 transition-all duration-200 cursor-pointer"
                >
                  <span class="material-symbols-outlined text-[18px]">chevron_right</span>
                </button>
              </nav>
            </div>
          </div>
        {/if}

        <footer class="mt-8 text-center text-xs text-slate-400 pb-4">
          <p>© 2025 BellaVita. Todos os direitos reservados.</p>
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
    onconfirm={alertState.onConfirm}
    oncancel={closeAlert}
  />

  {#snippet editModalBody()}
    <div class="space-y-4">
      <div>
        <span class="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-2">Foto do Profissional</span>
        <ImageUpload
          previewUrl={imagePreview}
          onselect={handleImageSelect}
          onclear={handleImageClear}
        />
      </div>

      <div>
        <label for="edit-nome" class="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-1">Nome Completo</label>
        <input
          id="edit-nome"
          type="text"
          bind:value={editNome}
          class="block w-full px-3.5 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-white dark:bg-gray-800 text-slate-900 dark:text-white focus:ring-2 focus:ring-orange-300 focus:border-orange-400 outline-none text-sm transition-all"
          placeholder="Ex: Ana Maria Costa"
        />
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label for="edit-cpf" class="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-1">CPF</label>
          <input
            id="edit-cpf"
            type="text"
            bind:value={editCpf}
            disabled
            class="block w-full px-3.5 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-slate-50 dark:bg-gray-800 text-slate-400 text-sm cursor-not-allowed"
            placeholder="000.000.000-00"
          />
        </div>
        <div>
          <label for="edit-telefone" class="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-1">Telefone / WhatsApp</label>
          <input
            id="edit-telefone"
            type="text"
            bind:value={editTelefone}
            class="block w-full px-3.5 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-white dark:bg-gray-800 text-slate-900 dark:text-white focus:ring-2 focus:ring-orange-300 focus:border-orange-400 outline-none text-sm transition-all"
            placeholder="(00) 00000-0000"
          />
        </div>
      </div>

      <div>
        <label for="edit-email" class="block text-sm font-semibold text-slate-700 dark:text-slate-200 mb-1">Endereço de Email</label>
        <input
          id="edit-email"
          type="email"
          bind:value={editEmail}
          class="block w-full px-3.5 py-2.5 border border-gray-200 dark:border-gray-700 rounded-xl bg-white dark:bg-gray-800 text-slate-900 dark:text-white focus:ring-2 focus:ring-orange-300 focus:border-orange-400 outline-none text-sm transition-all"
          placeholder="profissional@bellavita.com"
        />
      </div>

      <div>
        <h3 class="flex items-center text-sm font-bold text-slate-800 dark:text-slate-100 mb-1">
          <span class="material-symbols-outlined text-orange-500 mr-2 text-[18px]">spa</span>
          Serviços Atribuídos
        </h3>
        <p class="text-xs text-slate-400 mb-3">
          Selecione todos os serviços que este profissional está qualificado para realizar.
        </p>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-2.5 max-h-44 overflow-y-auto pr-1">
          {#each availableServices as service}
            <button
              type="button"
              class="flex items-center justify-between p-3 rounded-xl border transition-all duration-200 text-left group cursor-pointer
                {editServiceIds.includes(service.ID)
                  ? 'bg-orange-50 border-orange-300 text-orange-700 shadow-sm'
                  : 'bg-slate-50 dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-slate-600 dark:text-slate-300 hover:border-orange-200 hover:bg-orange-50/50'}"
              onclick={() => toggleService(service.ID)}
            >
              <span class="text-sm font-medium">{service.Nome}</span>
              {#if editServiceIds.includes(service.ID)}
                <span class="material-symbols-outlined text-orange-500 text-lg">check_circle</span>
              {:else}
                <span class="material-symbols-outlined text-slate-300 group-hover:text-slate-400 text-lg">radio_button_unchecked</span>
              {/if}
            </button>
          {/each}
        </div>
      </div>
    </div>
  {/snippet}

  {#snippet editModalFooter()}
    <button
      onclick={() => (showEditModal = false)}
      class="px-4 py-2.5 text-sm font-semibold text-slate-700 dark:text-slate-200 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl hover:bg-slate-50 dark:hover:bg-gray-700 transition-all duration-200 cursor-pointer"
    >
      Cancelar
    </button>
    <button
      onclick={handleUpdate}
      disabled={isSaving}
      class="px-5 py-2.5 text-sm font-semibold text-white bg-orange-500 rounded-xl hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-400 disabled:opacity-50 transition-all duration-200 shadow-sm shadow-orange-500/20 cursor-pointer"
    >
      {isSaving ? "Salvando..." : "Salvar Alterações"}
    </button>
  {/snippet}

  <Modal
    show={showEditModal}
    title="Editar Profissional"
    maxWidth="max-w-2xl"
    onclose={() => (showEditModal = false)}
    children={editModalBody}
    footer={editModalFooter}
  />
</div>
