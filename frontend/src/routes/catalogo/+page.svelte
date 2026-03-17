<script lang="ts">
  import Sidebar from "$lib/components/Sidebar.svelte";
  import DashboardNavbar from "$lib/components/DashboardNavbar.svelte";
  import ImageUpload from "$lib/components/ImageUpload.svelte";
  import CategoryModal from "$lib/components/CategoryModal.svelte";
  import { fetchApi } from "$lib/utils/api";
  import {
    uploadImageToSupabase,
    deleteImageFromSupabase,
  } from "$lib/utils/imageUpload";
  import {
    validateServiceForm,
    convertPriceTocents,
    type ValidationErrors,
  } from "$lib/utils/validation";

  // ========================================
  // VARIÁVEIS DE ESTADO DO FORMULÁRIO
  // ========================================

  /** Nome do serviço */
  let nome = $state<string>("");

  /** Lista de categorias disponíveis */
  let categories = $state<string[]>([
    "Cabelo",
    "Unhas",
    "Estética Facial",
    "Estética Corporal",
    "Depilação",
    "Maquiagem",
    "Massagem",
    "Tratamentos",
  ]);

  /** Categoria selecionada */
  let categoria = $state<string>("");

  /** Duração padrão em minutos */
  let duracao_padrao = $state<number>(60);

  /** Preço em formato string */
  let preco = $state<string>("");

  /** Descrição do serviço */
  let description = $state<string>("");

  /** Indica se está enviando */
  let loading = $state<boolean>(false);

  /** Mensagem de feedback */
  let message = $state<string>("");

  /** Tipo da mensagem */
  let messageType = $state<"success" | "error" | "">("");

  /** Erros de validação */
  let errors = $state<ValidationErrors>({});

  // ========================================
  // VARIÁVEIS DO MODAL DE CATEGORIA
  // ========================================

  /** Controla visibilidade do modal */
  let showCategoryModal = $state<boolean>(false);

  // ========================================
  // VARIÁVEIS DE UPLOAD DE IMAGEM
  // ========================================

  /** Arquivo selecionado */
  let selectedImage = $state<File | null>(null);

  /** URL de preview */
  let imagePreview = $state<string | null>(null);

  /** Indica se está fazendo upload */
  let uploadingImage = $state<boolean>(false);

  /** URL pública da imagem */
  let uploadedImageUrl = $state<string | null>(null);

  /** Erro de imagem */
  let imageError = $state<string>("");

  // ========================================
  // FUNÇÕES
  // ========================================

  /**
   * Abre modal de nova categoria
   */
  function openCategoryModal(): void {
    showCategoryModal = true;
  }

  /**
   * Adiciona nova categoria
   */
  function handleAddCategory(newCategory: string): void {
    categories = [...categories, newCategory];
    categoria = newCategory;
  }

  /**
   * Manipula seleção de imagem
   */
  function handleImageSelect(file: File, preview: string): void {
    selectedImage = file;
    imagePreview = preview;
    imageError = "";
  }

  /**
   * Limpa imagem selecionada
   */
  function handleImageClear(): void {
    selectedImage = null;
    imagePreview = null;
    uploadedImageUrl = null;
    imageError = "";
  }

  /**
   * Valida formulário
   */
  function validate(): boolean {
    const formData = { nome, categoria, duracao_padrao, preco, description };
    const hasImage = !!(selectedImage || uploadedImageUrl);

    errors = validateServiceForm(formData, hasImage);

    if (!hasImage) {
      imageError = "A imagem é obrigatória.";
    }

    return Object.keys(errors).length === 0;
  }

  /**
   * Envia formulário
   */
  async function handleSubmit(event: SubmitEvent): Promise<void> {
    event.preventDefault();
    loading = true;
    message = "";
    imageError = "";

    if (!validate()) {
      loading = false;
      return;
    }

    try {
      // Upload da imagem
      let imageUrl = uploadedImageUrl;

      if (selectedImage && !uploadedImageUrl) {
        uploadingImage = true;
        try {
          imageUrl = await uploadImageToSupabase(selectedImage);
        } catch (uploadErr: any) {
          imageError = uploadErr?.message ?? "Erro ao enviar imagem.";
          loading = false;
          uploadingImage = false;
          return;
        }
        uploadingImage = false;

        if (!imageUrl) {
          imageError = "Erro ao enviar imagem.";
          loading = false;
          return;
        }

        uploadedImageUrl = imageUrl;
      }

      // Prepara payload
      const precoEmCentavos = convertPriceTocents(preco);
      const payload = {
        categoria,
        duracao_padrao: Math.floor(duracao_padrao),
        nome,
        preco: precoEmCentavos,
        image_url: imageUrl,
      };

      // Envia para backend
      const url = "/api/v1/catalogos/";
      console.log(`Sending POST request to ${url}`, payload);

      const response = await fetchApi(url, {
        method: "POST",
        body: JSON.stringify(payload),
      });

      if (response.ok) {
        message = "Serviço cadastrado com sucesso!";
        messageType = "success";

        // Reset
        nome = "";
        categoria = "";
        duracao_padrao = 60;
        preco = "";
        description = "";
        handleImageClear();

        // Limpa mensagem após 3 segundos
        setTimeout(() => {
          message = "";
          messageType = "";
        }, 3000);
      } else {
        // Se falhar no backend, deleta a imagem que subiu
        if (uploadedImageUrl) {
          console.log(
            "Falha ao salvar serviço, deletando imagem órfã...",
            uploadedImageUrl,
          );
          await deleteImageFromSupabase(uploadedImageUrl);
          uploadedImageUrl = null;
        }

        const errorData = await response.json().catch(() => ({}));
        message = `Erro ao cadastrar: ${errorData.message || response.statusText}`;
        messageType = "error";
      }
    } catch (error) {
      message = "Erro de conexão com o servidor.";
      messageType = "error";
      console.error(error);
    } finally {
      loading = false;
    }
  }
</script>

<div
  class="font-body bg-[hsl(var(--bs-background))] text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
  <Sidebar />
  <main class="flex-1 flex flex-col h-full overflow-hidden relative">
    <DashboardNavbar />
    <div class="flex-1 overflow-y-auto p-6 md:p-8">
      <div class="max-w-5xl mx-auto space-y-6">
        <div>
          <nav class="flex text-sm text-gray-500 dark:text-gray-400 mb-2">
            <a class="hover:text-primary transition-colors" href="/">Início</a>
            <span class="mx-2">/</span>
            <a class="hover:text-primary transition-colors" href="/catalogo"
              >Serviços</a
            >
            <span class="mx-2">/</span>
            <span class="text-gray-900 dark:text-white font-medium"
              >Novo Serviço</span
            >
          </nav>
          <h1 class="text-2xl font-bold text-gray-900 dark:text-white">
            Cadastro de Serviços
          </h1>
          <p class="text-gray-500 dark:text-gray-400 mt-1">
            Adicione novos serviços ao catálogo do salão. Organize por
            categorias para facilitar a busca.
          </p>
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div class="lg:col-span-2 space-y-6">
            <form
              class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-6 md:p-8"
              onsubmit={handleSubmit}
            >
              <div class="space-y-6">
                <div>
                  <label
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    for="service-name"
                  >
                    Nome do Serviço
                  </label>
                  <input
                    class="block w-full rounded-md {errors.nome
                      ? 'border-red-500'
                      : 'border-border-light dark:border-border-dark'} bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-primary focus:ring-primary sm:text-sm py-2.5 px-3 transition-colors"
                    id="service-name"
                    name="service-name"
                    placeholder="Ex: Corte Feminino, Manicure Completa, Hidratação Profunda"
                    type="text"
                    bind:value={nome}
                    minlength="3"
                    required
                  />
                  {#if errors.nome}
                    <p class="mt-1 text-xs text-red-500">{errors.nome}</p>
                  {/if}
                </div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div class="flex items-end gap-2">
                    <div class="flex-1">
                      <label
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                        for="category"
                      >
                        Categoria
                      </label>
                      <select
                        class="block w-full rounded-md border-border-light dark:border-border-dark bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-primary focus:ring-primary sm:text-sm py-2.5 px-3 transition-colors"
                        id="category"
                        name="category"
                        bind:value={categoria}
                        required
                      >
                        <option
                          value=""
                          class="bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                          >Selecione uma categoria</option
                        >
                        {#each categories as cat}
                          <option
                            value={cat}
                            class="bg-white dark:bg-gray-800 text-gray-900 dark:text-white"
                            >{cat}</option
                          >
                        {/each}
                      </select>
                    </div>
                    <button
                      onclick={openCategoryModal}
                      type="button"
                      class="h-[42px] w-[42px] flex items-center justify-center border border-border-light dark:border-border-dark rounded-md bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-500 hover:text-brand-orange hover:bg-brand-orange/5 hover:border-brand-orange transition-all duration-200 group active:scale-95 shadow-sm"
                      title="Nova Categoria"
                    >
                      <span
                        class="material-symbols-outlined text-[20px] group-hover:rotate-90 transition-transform duration-300"
                        >add</span
                      >
                    </button>
                  </div>
                  <div>
                    <label
                      class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                      for="duration"
                    >
                      Duração Padrão
                    </label>
                    <div class="relative rounded-md shadow-sm">
                      <input
                        class="block w-full rounded-md {errors.duracao
                          ? 'border-red-500'
                          : 'border-border-light dark:border-border-dark'} bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-primary focus:ring-primary sm:text-sm py-2.5 px-3 pr-12 transition-colors"
                        id="duration"
                        name="duration"
                        placeholder="60"
                        type="number"
                        step="1"
                        min="0"
                        onkeypress={(e) => {
                          if (e.key === "." || e.key === ",")
                            e.preventDefault();
                        }}
                        bind:value={duracao_padrao}
                        required
                      />
                      {#if errors.duracao}
                        <p class="mt-1 text-xs text-red-500">
                          {errors.duracao}
                        </p>
                      {/if}
                      <div
                        class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-3"
                      >
                        <span
                          class="text-gray-500 dark:text-gray-400 sm:text-sm"
                          >min</span
                        >
                      </div>
                    </div>
                  </div>
                </div>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <label
                      class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                      for="price"
                    >
                      Preço
                    </label>
                    <div class="relative rounded-md shadow-sm">
                      <div
                        class="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3"
                      >
                        <span
                          class="text-gray-500 dark:text-gray-400 sm:text-sm"
                          >R$</span
                        >
                      </div>
                      <input
                        class="block w-full rounded-md {errors.preco
                          ? 'border-red-500'
                          : 'border-border-light dark:border-border-dark'} bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-primary focus:ring-primary sm:text-sm py-2.5 px-3 pl-10 transition-colors"
                        id="price"
                        name="price"
                        placeholder="0.00"
                        type="text"
                        bind:value={preco}
                        required
                      />
                    </div>
                    {#if errors.preco}
                      <p class="mt-1 text-xs text-red-500">{errors.preco}</p>
                    {/if}
                  </div>
                </div>
                <div>
                  <label
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    for="description"
                  >
                    Descrição
                  </label>
                  <textarea
                    class="block w-full rounded-md border-border-light dark:border-border-dark bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 text-gray-900 dark:text-white shadow-sm focus:border-primary focus:ring-primary sm:text-sm py-2.5 px-3 transition-colors"
                    id="description"
                    name="description"
                    placeholder="Descreva os detalhes do serviço, técnicas utilizadas, benefícios, produtos aplicados..."
                    rows="4"
                    bind:value={description}
                  ></textarea>
                </div>
              </div>
              <div
                class="mt-8 flex items-center justify-end space-x-4 pt-6 border-t border-border-light dark:border-border-dark"
              >
                <a
                  href="/catalogo/listagem"
                  class="px-6 py-2 border border-border-light dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 font-medium hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors"
                >
                  Cancelar
                </a>
                <button
                  class="px-8 py-2 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-md font-medium shadow-md transition-colors flex items-center disabled:opacity-50"
                  type="submit"
                  disabled={loading}
                >
                  <span class="material-icons text-sm mr-2"
                    >{loading ? "sync" : "save"}</span
                  >
                  {#if loading && uploadingImage}
                    Enviando imagem...
                  {:else if loading}
                    Salvando...
                  {:else}
                    Salvar Serviço
                  {/if}
                </button>
              </div>
              {#if message}
                <div
                  class="mt-4 p-4 rounded-md {messageType === 'success'
                    ? 'bg-green-50 text-green-700 border border-green-200'
                    : 'bg-red-50 text-red-700 border border-red-200'}"
                >
                  {message}
                </div>
              {/if}
            </form>
          </div>
          <div class="lg:col-span-1 space-y-6">
            <div
              class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-6"
            >
              <h3 class="font-semibold text-gray-900 dark:text-white mb-4">
                Imagem de Capa
              </h3>

              <ImageUpload
                previewUrl={imagePreview}
                onselect={({ file, preview }) => handleImageSelect(file, preview)}
                onclear={handleImageClear}
              />
              {#if imageError}
                <p class="mt-2 text-xs text-red-500">{imageError}</p>
              {/if}
            </div>
            <div
              class="bg-brand-orange/5 rounded-lg border border-brand-orange/20 p-5"
            >
              <div class="flex items-start">
                <span class="material-icons text-brand-orange mt-0.5 mr-3"
                  >info</span
                >
                <div>
                  <h4 class="text-sm font-medium text-brand-orange">
                    Dica Profissional
                  </h4>
                  <p class="mt-1 text-sm text-gray-600 dark:text-gray-400">
                    Serviços com descrições completas e imagens profissionais
                    atraem mais clientes e facilitam o agendamento.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <footer
        class="mt-12 text-center text-sm text-gray-500 dark:text-gray-400 pb-6"
      >
        <p>© 2023 BellaVita Salon. Todos os direitos reservados.</p>
        <div class="mt-2 space-x-4">
          <a class="hover:text-primary" href="/">Termos</a>
          <a class="hover:text-primary" href="/">Privacidade</a>
        </div>
      </footer>
    </div>
  </main>
</div>

<CategoryModal
  bind:show={showCategoryModal}
  existingCategories={categories}
  onadd={handleAddCategory}
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
