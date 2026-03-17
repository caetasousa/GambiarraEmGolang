<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { auth, logout } from "$lib/stores/auth.svelte";
    import { sidebarOpen, closeSidebar } from "$lib/stores/sidebar.svelte";

    interface NavItem {
        label: string;
        href: string;
        icon: string;
        category?: string;
        roles?: string[];
    }

    const allItems: NavItem[] = [
        { label: "Agenda do Profissional", href: "/prestadores/agenda", icon: "calendar_month", roles: ["prestador", "admin"] },
        { label: "Meus Agendamentos", href: "/clientes/meus-agendamentos", icon: "event_available", roles: ["cliente"] },
        { label: "Catálogo de Serviços", href: "/listagem", icon: "storefront", roles: ["cliente"] },
        { label: "Listagem de Profissionais", href: "/prestadores/editar", icon: "group", category: "Prestadores", roles: ["admin"] },
        { label: "Cadastrar Profissional", href: "/prestadores/cadastro", icon: "person_add", roles: ["admin"] },
        { label: "Listagem de Serviços", href: "/catalogo/listagem", icon: "list_alt", category: "Serviços", roles: ["admin"] },
        { label: "Cadastrar Serviço", href: "/catalogo", icon: "add_box", roles: ["admin"] },
    ];

    let role = $derived(auth.user?.role ?? "");
    let items = $derived(allItems.filter(item => !item.roles || item.roles.includes(role)));
    let roleLabel = $derived(
        role === "admin" ? "Admin" : role === "prestador" ? "Profissional" : "Cliente"
    );

    const staticItems = ["/perfil/meus-dados", "/perfil"];

    function isActive(href: string): boolean {
        const path = $page.url.pathname;
        if (href === "/") return path === "/";
        if (path === href) return true;
        if (path.startsWith(href + "/")) {
            const allHrefs = [...items.map(i => i.href), ...staticItems];
            const hasMoreSpecificMatch = allHrefs.some(
                h => h !== href && path.startsWith(h) && h.length > href.length,
            );
            return !hasMoreSpecificMatch;
        }
        return false;
    }

    function handleLogout() {
        logout();
        goto("/login");
    }

    function closeMenu() {
        closeSidebar();
    }


</script>

{#if auth.isAuthenticated}
    <!-- Overlay mobile -->
    {#if sidebarOpen.value}
        <div
            class="fixed inset-0 bg-black/40 backdrop-blur-sm z-40 md:hidden"
            onclick={closeMenu}
            role="button"
            tabindex="-1"
            onkeydown={(e) => e.key === 'Escape' && closeMenu()}
        ></div>
    {/if}

    <!-- Sidebar -->
    <aside
        class="w-72 bg-white dark:bg-gray-900 border-r border-gray-100 dark:border-gray-800 flex-col h-full flex-shrink-0
               md:flex shadow-sm
               {sidebarOpen.value ? 'flex fixed inset-y-0 left-0 z-50 shadow-2xl' : 'hidden'}"
    >
        <!-- Logo -->
        <a href="/" class="h-16 flex items-center px-6 gap-3 group border-b border-gray-100 dark:border-gray-800 flex-shrink-0" onclick={closeMenu}>
            <div class="w-10 h-10 bg-orange-500 rounded-xl flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-orange-500/30 transition-all duration-300 group-hover:scale-105 group-hover:shadow-orange-500/40 flex-shrink-0">
                <span class="font-serif italic">B</span>
            </div>
            <span class="text-lg font-bold text-slate-900 dark:text-white tracking-tight" style="font-family: 'Cormorant', serif; font-style: italic;">
                Bella<span class="text-orange-500">Vita</span>
            </span>
        </a>

        <!-- Navigation -->
        <nav class="flex-1 overflow-y-auto py-5 px-3 space-y-0.5">
            {#each items as item}
                {#if item.category}
                    <div class="pt-5 pb-2 px-4 text-xs font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-widest" aria-hidden="true">
                        {item.category}
                    </div>
                {/if}
                <a
                    class="group flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 cursor-pointer
                        {isActive(item.href)
                            ? 'bg-orange-50 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400 border-l-4 border-orange-500 pl-3'
                            : 'text-slate-600 dark:text-slate-300 hover:bg-orange-50 dark:hover:bg-orange-900/20 hover:text-orange-600 dark:hover:text-orange-400 border-l-4 border-transparent pl-3'}"
                    href={item.href}
                    onclick={closeMenu}
                >
                    <span class="material-symbols-outlined text-[20px] flex-shrink-0 transition-colors duration-200 {isActive(item.href) ? 'text-orange-500 dark:text-orange-400' : 'text-slate-400 dark:text-slate-500 group-hover:text-orange-500 dark:group-hover:text-orange-400'}">
                        {item.icon}
                    </span>
                    <span class="truncate">{item.label}</span>
                </a>
            {/each}
        </nav>

        <!-- Account section -->
        <div class="border-t border-gray-100 dark:border-gray-800 px-3 pt-4 pb-2 space-y-0.5">
            <p class="px-4 text-xs font-semibold text-slate-400 dark:text-slate-500 uppercase tracking-widest mb-2">Conta</p>
            <a
                class="group flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 cursor-pointer
                    {isActive('/perfil/meus-dados')
                        ? 'bg-orange-50 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400 border-l-4 border-orange-500 pl-3'
                        : 'text-slate-600 dark:text-slate-300 hover:bg-orange-50 dark:hover:bg-orange-900/20 hover:text-orange-600 dark:hover:text-orange-400 border-l-4 border-transparent pl-3'}"
                href="/perfil/meus-dados"
                onclick={closeMenu}
            >
                <span class="material-symbols-outlined text-[20px] flex-shrink-0 transition-colors duration-200 {isActive('/perfil/meus-dados') ? 'text-orange-500 dark:text-orange-400' : 'text-slate-400 dark:text-slate-500 group-hover:text-orange-500'}">person</span>
                <span>Meus Dados</span>
            </a>
            {#if role !== 'cliente'}
            <a
                class="group flex items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 cursor-pointer
                    {isActive('/perfil')
                        ? 'bg-orange-50 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400 border-l-4 border-orange-500 pl-3'
                        : 'text-slate-600 dark:text-slate-300 hover:bg-orange-50 dark:hover:bg-orange-900/20 hover:text-orange-600 dark:hover:text-orange-400 border-l-4 border-transparent pl-3'}"
                href="/perfil"
                onclick={closeMenu}
            >
                <span class="material-symbols-outlined text-[20px] flex-shrink-0 transition-colors duration-200 {isActive('/perfil') ? 'text-orange-500 dark:text-orange-400' : 'text-slate-400 dark:text-slate-500 group-hover:text-orange-500'}">edit_calendar</span>
                <span>Agenda Diária</span>
            </a>
            {/if}
            <button
                onclick={handleLogout}
                class="group flex w-full items-center gap-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 cursor-pointer text-slate-600 dark:text-slate-300 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 dark:hover:text-red-400 border-l-4 border-transparent pl-3"
            >
                <span class="material-symbols-outlined text-[20px] flex-shrink-0 text-slate-400 dark:text-slate-500 group-hover:text-red-500 transition-colors duration-200">logout</span>
                <span>Sair</span>
            </button>
        </div>

        <!-- User info footer -->
        <div class="border-t border-gray-100 dark:border-gray-800 p-4">
            <div class="flex items-center gap-3">
                {#if auth.user?.image_url}
                    <img
                        alt="Avatar"
                        class="h-10 w-10 rounded-full object-cover ring-2 ring-orange-200 ring-offset-1 flex-shrink-0"
                        src={auth.user.image_url}
                    />
                {:else}
                    <div class="h-10 w-10 rounded-full bg-gradient-to-br from-orange-400 to-pink-500 flex items-center justify-center flex-shrink-0 ring-2 ring-orange-200 ring-offset-1">
                        <span class="text-white text-xs font-bold">
                            {(auth.user?.nome ?? "U").split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase()}
                        </span>
                    </div>
                {/if}
                <div class="min-w-0 flex-1">
                    <p class="text-sm font-semibold text-slate-800 dark:text-slate-100 truncate" style="font-family: 'Montserrat', sans-serif;">
                        {auth.user?.nome ?? "Usuário"}
                    </p>
                    <p class="text-xs text-slate-400 dark:text-slate-500 font-medium">{roleLabel}</p>
                </div>
            </div>
        </div>
    </aside>
{/if}
