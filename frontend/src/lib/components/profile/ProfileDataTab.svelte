<script lang="ts">
    import { createEventDispatcher } from "svelte";
    import ImageUpload from "$lib/components/ImageUpload.svelte";
    import type { ProviderProfile, ProviderData } from "$lib/types/profile"; // Import shared type

    export let profile: ProviderProfile;
    export let isEditing = false;

    const dispatch = createEventDispatcher<{
        toggleEdit: void;
        cancelEdit: void;
        save: void;
        imageSelect: { file: File; preview: string };
        imageClear: void;
    }>();

    // Helper to emit imageSelect
    function handleImageSelect(
        event: CustomEvent<{ file: File; preview: string }>,
    ) {
        dispatch("imageSelect", event.detail);
    }
</script>

<div
    class="bg-[hsl(var(--bs-card))] rounded-xl border border-border-light dark:border-border-dark p-6 md:p-8 animate-fade-in text-left"
>
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

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
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
    </div>

    <div class="mt-8 flex justify-end gap-3">
        {#if !isEditing}
            <button
                on:click={() => dispatch("toggleEdit")}
                class="px-6 py-2.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-bold shadow-md transition-all"
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
