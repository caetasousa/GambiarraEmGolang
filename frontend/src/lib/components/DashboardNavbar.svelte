<script lang="ts">
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { goto } from "$app/navigation";
    import { logout } from "$lib/stores/auth";

    async function handleLogout() {
        logout();
        await goto("/login");
    }
</script>

<header
    class="h-16 bg-[hsl(var(--bs-card))] border-b border-border-light dark:border-border-dark flex items-center justify-between px-6 z-10 flex-shrink-0"
>
    <div class="flex items-center flex-1">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
            Dashboard
        </h2>
    </div>
    <div class="flex items-center space-x-2">
        <button
            class="p-2 rounded-full text-gray-400 hover:text-gray-500 dark:hover:text-gray-300 focus:outline-none transition-colors"
        >
            <span class="material-icons">notifications</span>
        </button>
        <ThemeToggle />
        <button
            onclick={handleLogout}
            class="flex items-center gap-1 px-3 py-2 rounded-lg text-sm font-medium text-gray-600 dark:text-gray-300 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 dark:hover:text-red-400 transition-colors"
            title="Sair"
        >
            <span class="material-icons text-[18px]">logout</span>
            <span class="hidden sm:inline">Sair</span>
        </button>
    </div>
</header>
