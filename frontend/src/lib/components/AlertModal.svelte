<script lang="ts">
    import Modal from "./Modal.svelte";
    import { createEventDispatcher } from "svelte";

    export let show = false;
    export let title = "Atenção";
    export let message = "";
    export let type: "success" | "error" | "warning" | "info" | "primary" =
        "info";
    export let confirmText = "Confirmar";
    export let cancelText = "Cancelar";
    export let showCancel = false;

    const dispatch = createEventDispatcher();

    function handleConfirm() {
        dispatch("confirm");
        show = false;
    }

    function handleCancel() {
        dispatch("cancel");
        show = false;
    }

    // Cores baseadas no tipo
    $: iconColor =
        type === "error"
            ? "text-red-500"
            : type === "success"
              ? "text-green-500"
              : type === "warning" || type === "primary"
                ? "text-brand-orange"
                : "text-blue-500";

    $: iconName =
        type === "error"
            ? "error"
            : type === "success"
              ? "check_circle"
              : type === "warning"
                ? "warning"
                : type === "primary"
                  ? "check_circle"
                  : "info";

    $: confirmBtnClass =
        type === "error"
            ? "bg-red-500 hover:bg-red-600 focus:ring-red-500"
            : type === "success"
              ? "bg-green-500 hover:bg-green-600 focus:ring-green-500"
              : type === "warning" || type === "primary"
                ? "bg-brand-orange hover:bg-brand-orange-hover focus:ring-brand-orange"
                : "bg-blue-500 hover:bg-blue-600 focus:ring-blue-500";
</script>

<Modal
    {show}
    {title}
    maxWidth="max-w-sm"
    zIndex="z-[60]"
    on:close={handleCancel}
>
    <div slot="body" class="text-center py-2">
        <div class="mb-4 flex justify-center">
            <span class="material-symbols-outlined text-5xl {iconColor}">
                {iconName}
            </span>
        </div>
        <p class="text-gray-600 dark:text-gray-300 text-base leading-relaxed">
            {message}
        </p>
    </div>

    <div slot="footer" class="w-full flex justify-end gap-3">
        {#if showCancel}
            <button
                on:click={handleCancel}
                class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-800 border border-border-light dark:border-border-dark rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
            >
                {cancelText}
            </button>
        {/if}
        <button
            on:click={handleConfirm}
            class="px-4 py-2 text-sm font-medium text-white rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors {confirmBtnClass}"
        >
            {confirmText}
        </button>
    </div>
</Modal>
