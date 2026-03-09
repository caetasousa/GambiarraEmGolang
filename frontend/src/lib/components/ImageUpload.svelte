<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import {
        validateImageFile,
        createImagePreview,
    } from "$lib/utils/imageUpload";

    /** Dispatcher para eventos do componente */
    const dispatch = createEventDispatcher<{
        select: { file: File; preview: string };
        clear: void;
    }>();

    /** URL de preview da imagem (passada pelo componente pai) */
    export let previewUrl: string | null = null;

    /** ID do input de arquivo (para acessibilidade/labels externos) */
    export let id: string = "image-upload";

    /** Modo de exibição: 'banner' (padrão) ou 'avatar' (circular, menor) */
    export let mode: "banner" | "avatar" = "banner";

    /** Texto de ajuda exibido abaixo da imagem */
    export let helpText: string = "A imagem será enviada ao salvar o serviço";

    /** Mensagem de erro (gerenciada internamente) */
    let error: string = "";

    /**
     * Manipula a seleção de arquivo de imagem
     */
    async function handleFileSelect(event: Event) {
        const input = event.target as HTMLInputElement;
        const file = input.files?.[0];

        if (!file) return;

        // Valida o arquivo
        const validation = validateImageFile(file);
        if (!validation.valid) {
            error = validation.error || "";
            return;
        }

        // Limpa erro e cria preview
        error = "";

        try {
            const preview = await createImagePreview(file);
            dispatch("select", { file, preview });
        } catch (err) {
            error = "Erro ao criar preview da imagem.";
            console.error(err);
        }
    }

    /**
     * Limpa a imagem selecionada
     */
    function handleClear() {
        dispatch("clear");
    }
</script>

{#if previewUrl}
    <!-- Preview da Imagem -->
    <!-- Preview da Imagem -->
    <div class="relative group {mode === 'avatar' ? 'w-32 h-32 mx-auto' : ''}">
        <img
            src={previewUrl}
            alt="Preview da imagem"
            class="w-full object-cover border border-gray-200 dark:border-gray-700 {mode ===
            'avatar'
                ? 'h-32 rounded-full'
                : 'h-48 rounded-lg'}"
        />

        <!-- Overlay com botão de alterar -->
        <div
            class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center justify-center {mode ===
            'avatar'
                ? 'rounded-full'
                : 'rounded-lg'}"
        >
            <label
                for={id}
                class="cursor-pointer bg-white/20 backdrop-blur-sm border border-white/50 text-white hover:bg-white/30 transition-colors font-medium flex items-center justify-center {mode ===
                'avatar'
                    ? 'p-2 rounded-full'
                    : 'px-4 py-2 rounded-md'}"
                title="Alterar Imagem"
            >
                <span class="material-icons text-sm {mode === 'banner' ? 'mr-2' : ''}">edit</span>
                {#if mode === 'banner'}
                    Alterar Imagem
                {/if}
                <input
                    {id}
                    type="file"
                    accept="image/*"
                    class="sr-only"
                    on:change={handleFileSelect}
                />
            </label>
        </div>

        <button
            type="button"
            on:click={handleClear}
            class="absolute flex items-center justify-center bg-red-500 hover:bg-red-600 text-white rounded-full transition-colors shadow-sm z-10 {mode ===
            'avatar'
                ? 'top-0 right-0 h-6 w-6'
                : 'top-2 right-2 h-8 w-8'}"
            title="Remover imagem"
        >
            <span class="material-icons text-lg">close</span>
        </button>
    </div>

    <div
        class="mt-3 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-md"
    >
        <div class="flex items-center text-blue-700 dark:text-blue-400">
            <span class="material-icons text-sm mr-2">info</span>
            <span class="text-sm"
                >{helpText}</span
            >
        </div>
    </div>
{:else}
    <!-- Área de Upload -->
    <label
        for={id}
        class="mt-2 flex justify-center border-dashed border-gray-300 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors cursor-pointer group {mode ===
        'avatar'
            ? 'rounded-full w-32 h-32 mx-auto items-center border-2'
            : 'rounded-lg px-6 py-10 border'}"
    >
        <div class="text-center">
            <span
                class="material-icons mx-auto text-gray-300 dark:text-gray-500 group-hover:text-brand-orange transition-colors {mode ===
                'avatar'
                    ? 'text-3xl'
                    : 'h-12 w-12 text-5xl'}"
                >{mode === "avatar" ? "add_a_photo" : "image"}</span
            >
            {#if mode === 'banner'}
            <div
                class="mt-4 flex text-sm leading-6 text-gray-600 dark:text-gray-400 justify-center"
            >
                <span
                    class="relative cursor-pointer rounded-md font-semibold text-brand-orange hover:text-brand-orange-hover focus-within:outline-none focus-within:ring-2 focus-within:ring-brand-orange focus-within:ring-offset-2"
                >
                    <span>Enviar arquivo</span>
                    <input
                        {id}
                        type="file"
                        accept="image/*"
                        class="sr-only"
                        on:change={handleFileSelect}
                    />
                </span>
            </div>
            <p class="text-xs leading-5 text-gray-500 dark:text-gray-500">
                PNG, JPG até 5MB
            </p>
            {:else}
             <!-- Avatar mode simple input -->
             <input
                {id}
                type="file"
                accept="image/*"
                class="sr-only"
                on:change={handleFileSelect}
            />
            {/if}
        </div>
    </label>
{/if}

{#if error}
    <div
        class="mt-3 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md"
    >
        <p class="text-sm text-red-700 dark:text-red-400">{error}</p>
    </div>
{/if}
