<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import ImageUpload from "$lib/components/ImageUpload.svelte";
    import type { ProviderProfile } from "$lib/types/profile";

    export let profile: ProviderProfile;
    export let isEditing = false;
    export let services: any[] = [];
    export let servicesLoading = false;

    $: isPrestadorOrAdmin = profile.role === "prestador" || profile.role === "admin";

    const dispatch = createEventDispatcher<{
        toggleEdit: void;
        cancelEdit: void;
        save: void;
        imageSelect: { file: File; preview: string };
        imageClear: void;
        addService: void;
        removeService: { service: any };
    }>();

    function handleImageSelect(
        event: CustomEvent<{ file: File; preview: string }>,
    ) {
        dispatch("imageSelect", event.detail);
    }
</script>

<div
    class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark p-6 md:p-8 animate-fade-in text-left space-y-8"
>
    <!-- Meus Dados Section -->
    <div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6">
            Meus Dados
        </h3>

        {#if isPrestadorOrAdmin}
            {#if isEditing}
                <div class="mb-6">
                    <label
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 text-center"
                        for="avatar-upload"
                    >
                        Foto de Perfil
                    </label>
                    <ImageUpload
                        id="avatar-upload"
                        mode="avatar"
                        helpText="A imagem será salva ao atualizar o perfil"
                        previewUrl={profile.avatarUrl}
                        on:select={handleImageSelect}
                        on:clear={() => dispatch("imageClear")}
                    />
                </div>
            {:else}
                <div class="flex items-center space-x-4 mb-8">
                    <img
                        alt={profile.name}
                        class="h-24 w-24 rounded-full border-4 border-gray-100 dark:border-gray-700 object-cover"
                        src={profile.avatarUrl ||
                            `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.name)}&background=random`}
                    />
                </div>
            {/if}
        {:else}
            <div class="flex items-center space-x-4 mb-8">
                <div class="h-24 w-24 rounded-full border-4 border-gray-100 dark:border-gray-700 bg-brand-orange/20 flex items-center justify-center">
                    <span class="material-icons text-brand-orange text-4xl">person</span>
                </div>
            </div>
        {/if}

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Nome: todos os roles -->
            <div class="col-span-1 md:col-span-2">
                <label
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    for="name">Nome Completo</label
                >
                <input
                    class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 p-2.5 disabled:opacity-60 disabled:cursor-not-allowed"
                    type="text"
                    disabled={!isEditing}
                    bind:value={profile.name}
                />
            </div>

            <!-- CPF: prestador e admin -->
            {#if isPrestadorOrAdmin}
                <div>
                    <label
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                        for="cpf">CPF</label
                    >
                    <input
                        class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 p-2.5 disabled:opacity-60 disabled:cursor-not-allowed"
                        type="text"
                        disabled={true}
                        title="O CPF não pode ser alterado"
                        bind:value={profile.cpf}
                    />
                </div>
            {/if}

            <!-- Telefone: todos os roles -->
            <div>
                <label
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    for="phone">Telefone</label
                >
                <input
                    class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 p-2.5 disabled:opacity-60 disabled:cursor-not-allowed"
                    type="text"
                    disabled={!isEditing}
                    bind:value={profile.phone}
                />
            </div>

            <!-- Email: todos os roles -->
            <div class="col-span-1 md:col-span-2">
                <label
                    class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    for="email">Email</label
                >
                <input
                    class="block w-full rounded-md border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-[hsl(var(--bs-muted))]/20 p-2.5 disabled:opacity-60 disabled:cursor-not-allowed"
                    type="email"
                    disabled={!isEditing}
                    bind:value={profile.email}
                />
            </div>

            <!-- Status Ativo: prestador e admin -->
            {#if isPrestadorOrAdmin}
                <div class="col-span-1 md:col-span-2">
                    <p
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
                    >Status</p>
                    <span
                        class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium {profile.ativo
                            ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                            : 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'}"
                    >
                        <span class="material-icons text-[14px]">{profile.ativo ? 'check_circle' : 'cancel'}</span>
                        {profile.ativo ? "Ativo" : "Inativo"}
                    </span>
                </div>
            {/if}
        </div>

        <div class="mt-8 flex justify-end gap-3">
            {#if !isEditing}
                <button
                    on:click={() => dispatch("toggleEdit")}
                    class="px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-lg font-bold shadow-lg shadow-brand-orange/20 transition-all"
                >
                    Editar Dados
                </button>
            {:else}
                <button
                    on:click={() => dispatch("cancelEdit")}
                    class="px-6 py-2.5 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-medium hover:bg-gray-50 dark:hover:bg-gray-700 transition-all"
                >
                    Cancelar
                </button>
                <button
                    on:click={() => dispatch("save")}
                    class="px-6 py-2.5 bg-brand-orange hover:bg-brand-orange-hover text-white rounded-lg font-bold shadow-lg shadow-brand-orange/20 transition-all"
                >
                    Salvar Alterações
                </button>
            {/if}
        </div>
    </div>

    {#if profile.role === "prestador"}
    <!-- Divider -->
    <div class="border-t border-border-light dark:border-border-dark"></div>

    <!-- Meus Serviços Section -->
    <div>
        <div class="flex justify-between items-center mb-6">
            <h3 class="text-xl font-bold text-gray-900 dark:text-white">
                Meus Serviços
            </h3>
            <button
                on:click={() => dispatch("addService")}
                class="text-primary font-medium hover:underline text-sm flex items-center"
            >
                <span class="material-symbols-outlined text-sm mr-1">add</span>
                Adicionar Serviço
            </button>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {#if servicesLoading}
                <div class="col-span-full flex justify-center py-8">
                    <span
                        class="material-symbols-outlined animate-spin text-brand-orange text-3xl"
                        >sync</span
                    >
                </div>
            {:else if services.length === 0}
                <div class="col-span-full text-center py-8 text-gray-500">
                    Nenhum serviço encontrado.
                </div>
            {:else}
                {#each services as service}
                    <div
                        class="relative flex items-start p-4 rounded-xl border border-border-light dark:border-border-dark bg-[hsl(var(--bs-card))] shadow-sm transition-all hover:border-brand-orange/50 group"
                    >
                        {#if service.ImagemUrl}
                            <div
                                class="h-12 w-12 rounded-lg bg-gray-100 overflow-hidden shrink-0 mr-3 border border-gray-200 dark:border-gray-700"
                            >
                                <img
                                    src={service.ImagemUrl}
                                    alt={service.Nome}
                                    class="h-full w-full object-cover"
                                />
                            </div>
                        {/if}

                        <div class="min-w-0 flex-1">
                            <div class="flex justify-between items-start mb-1">
                                <h4
                                    class="font-bold text-gray-900 dark:text-white text-sm line-clamp-1"
                                    title={service.Nome}
                                >
                                    {service.Nome}
                                </h4>
                                <span
                                    class="text-brand-orange font-bold text-sm whitespace-nowrap ml-2"
                                >
                                    R$ {(service.Preco / 100)
                                        .toFixed(2)
                                        .replace(".", ",")}
                                </span>
                            </div>
                            <p class="text-gray-500 dark:text-gray-400 text-xs">
                                {service.DuracaoPadrao} min • {service.Categoria}
                            </p>
                        </div>
                        <div class="ml-2 flex items-center self-center">
                            <button
                                on:click={() =>
                                    dispatch("removeService", { service })}
                                class="p-1.5 text-gray-400 hover:text-red-500 rounded-lg hover:bg-red-50 dark:hover:bg-red-900/10 transition-colors"
                                title="Remover Serviço"
                            >
                                <span
                                    class="material-symbols-outlined text-[20px]"
                                    >delete</span
                                >
                            </button>
                        </div>
                    </div>
                {/each}
            {/if}
        </div>
    </div>
    {/if}
</div>

<style>
    .animate-fade-in {
        animation: fadeIn 0.3s ease-in-out;
    }
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }
</style>
