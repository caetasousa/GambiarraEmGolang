<script lang="ts">
    import Modal from "./Modal.svelte";

    let {
        show = $bindable(false),
        title = "Atenção",
        message = "",
        type = "info" as "success" | "error" | "warning" | "info" | "primary",
        confirmText = "Confirmar",
        cancelText = "Cancelar",
        showCancel = false,
        onconfirm,
        oncancel,
    }: {
        show?: boolean;
        title?: string;
        message?: string;
        type?: "success" | "error" | "warning" | "info" | "primary";
        confirmText?: string;
        cancelText?: string;
        showCancel?: boolean;
        onconfirm?: () => void;
        oncancel?: () => void;
    } = $props();

    let iconColor = $derived(
        type === "error" ? "text-red-500"
        : type === "success" ? "text-green-500"
        : type === "warning" || type === "primary" ? "text-brand-orange"
        : "text-blue-500"
    );

    let iconName = $derived(
        type === "error" ? "error"
        : type === "success" ? "check_circle"
        : type === "warning" ? "warning"
        : type === "primary" ? "check_circle"
        : "info"
    );

    let confirmBtnClass = $derived(
        type === "error" ? "bg-red-500 hover:bg-red-600 focus:ring-red-500"
        : type === "success" ? "bg-green-500 hover:bg-green-600 focus:ring-green-500"
        : type === "warning" || type === "primary" ? "bg-brand-orange hover:bg-brand-orange-hover focus:ring-brand-orange"
        : "bg-blue-500 hover:bg-blue-600 focus:ring-blue-500"
    );

    function handleConfirm() {
        onconfirm?.();
        show = false;
    }

    function handleCancel() {
        oncancel?.();
        show = false;
    }
</script>

{#snippet footerContent()}
    <div class="w-full flex justify-end gap-3">
        {#if showCancel}
            <button
                onclick={handleCancel}
                class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
                {cancelText}
            </button>
        {/if}
        <button
            onclick={handleConfirm}
            class="px-4 py-2 text-sm font-medium text-white rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors {confirmBtnClass}"
        >
            {confirmText}
        </button>
    </div>
{/snippet}

<Modal
    {show}
    {title}
    maxWidth="max-w-sm"
    zIndex="z-[60]"
    onclose={handleCancel}
    footer={footerContent}
>
    <div class="text-center py-2">
        <div class="mb-4 flex justify-center">
            <span class="material-symbols-outlined text-5xl {iconColor}">
                {iconName}
            </span>
        </div>
        <p class="text-gray-600 dark:text-gray-300 text-base leading-relaxed">
            {message}
        </p>
    </div>
</Modal>
