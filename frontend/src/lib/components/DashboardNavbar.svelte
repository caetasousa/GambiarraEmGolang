<script lang="ts">
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { openSidebar } from "$lib/stores/sidebar.svelte";

    let { title = "Dashboard", onMenuClick = openSidebar }: { title?: string; onMenuClick?: () => void } = $props();
</script>

<header
    class="h-16 bg-white dark:bg-gray-900 border-b border-gray-100 dark:border-gray-800 shadow-sm flex items-center justify-between px-6 z-10 flex-shrink-0"
>
    <div class="flex items-center flex-1 gap-4">
        <button
            class="md:hidden p-2 rounded-xl text-slate-500 dark:text-slate-400 hover:bg-orange-50 dark:hover:bg-orange-900/20 hover:text-orange-600 transition-all duration-200 cursor-pointer"
            onclick={onMenuClick}
            aria-label="Abrir menu"
        >
            <span class="material-icons text-[22px]">menu</span>
        </button>
        <div>
            <h2
                class="text-2xl font-bold text-slate-900 dark:text-white leading-tight"
                style="font-family: 'Cormorant', serif;"
            >
                {title}
            </h2>
        </div>
    </div>

    <div class="flex items-center gap-2">
        <button
            class="p-2 rounded-xl text-slate-400 dark:text-slate-500 hover:text-orange-500 hover:bg-orange-50 dark:hover:bg-orange-900/20 transition-all duration-200 cursor-pointer focus:outline-none"
            aria-label="Notificações"
        >
            <span class="material-icons text-[22px]">notifications_none</span>
        </button>
        <ThemeToggle />
    </div>
</header>
