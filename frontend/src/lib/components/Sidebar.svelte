<script lang="ts">
    import { page } from "$app/stores";
    import { goto } from "$app/navigation";
    import { user, logout, isAuthenticated } from "$lib/stores/auth";

    interface NavItem {
        label: string;
        href: string;
        icon: string;
        category?: string;
        roles?: string[];
    }

    const allItems: NavItem[] = [
        {
            label: "Agenda do Profissional",
            href: "/prestadores/agenda",
            icon: "calendar_month",
            roles: ["prestador", "admin"],
        },
        {
            label: "Meus Agendamentos",
            href: "/clientes/meus-agendamentos",
            icon: "event_available",
            roles: ["cliente"],
        },
        {
            label: "Catálogo de Serviços",
            href: "/listagem",
            icon: "storefront",
            roles: ["cliente"],
        },
        {
            label: "Listagem de Profissionais",
            href: "/prestadores/editar",
            icon: "group",
            category: "Prestadores",
            roles: ["admin"],
        },
        {
            label: "Cadastrar Profissional",
            href: "/prestadores/cadastro",
            icon: "person_add",
            roles: ["admin"],
        },
        {
            label: "Listagem de Serviços",
            href: "/catalogo/listagem",
            icon: "list_alt",
            category: "Serviços",
            roles: ["admin"],
        },
        {
            label: "Cadastrar Serviço",
            href: "/catalogo",
            icon: "add_box",
            roles: ["admin"],
        },
    ];

    $: role = $user?.role ?? "";
    $: items = allItems.filter(
        (item) => !item.roles || item.roles.includes(role),
    );

    const staticItems = ["/perfil/meus-dados", "/perfil"];

    $: isActive = (href: string) => {
        const path = $page.url.pathname;
        if (href === "/") return path === "/";
        if (path === href) return true;
        if (path.startsWith(href + "/")) {
            const allHrefs = [...items.map(i => i.href), ...staticItems];
            const hasMoreSpecificMatch = allHrefs.some(
                (h) =>
                    h !== href &&
                    path.startsWith(h) &&
                    h.length > href.length,
            );
            return !hasMoreSpecificMatch;
        }
        return false;
    };

    function handleLogout() {
        logout();
        goto("/login");
    }

    $: roleLabel =
        role === "admin"
            ? "Admin"
            : role === "prestador"
              ? "Profissional"
              : "Cliente";
</script>

{#if $isAuthenticated}
<aside
    class="w-64 bg-[hsl(var(--bs-card))] flex-col hidden md:flex h-full flex-shrink-0"
>
    <a href="/" class="h-16 flex items-center px-6 gap-2 group">
        <div class="w-9 h-9 bg-orange-500 rounded-xl flex items-center justify-center text-white font-bold text-base shadow-lg shadow-orange-500/25 transition-transform duration-300 group-hover:scale-105">
            B
        </div>
        <span class="text-base font-bold text-gray-900 dark:text-white">
            Bella<span class="text-orange-500">Vita</span>
        </span>
    </a>

    <nav class="flex-1 overflow-y-auto py-4 px-3 space-y-1">
        {#each items as item}
            {#if item.category}
                <div
                    class="pt-4 pb-2 px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider"
                >
                    {item.category}
                </div>
            {/if}
            <a
                class="group flex items-center px-3 py-2 text-sm font-medium rounded-md transition-all duration-200 {isActive(
                    item.href,
                )
                    ? 'bg-brand-orange/10 text-brand-orange'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-brand-orange/10 hover:text-brand-orange'}"
                href={item.href}
            >
                <span
                    class="material-symbols-outlined mr-3 {isActive(item.href)
                        ? 'text-brand-orange'
                        : 'text-gray-400 group-hover:text-brand-orange transition-colors'}"
                >
                    {item.icon}
                </span>
                {item.label}
            </a>
        {/each}
    </nav>

    <div
        class="border-t border-border-light dark:border-border-dark p-4 space-y-1"
    >
        <h3
            class="px-3 text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2"
        >
            Conta
        </h3>
        <a
            class="group flex items-center px-3 py-2 text-sm font-medium rounded-md transition-all duration-200 {isActive(
                '/perfil/meus-dados',
            )
                ? 'bg-brand-orange/10 text-brand-orange'
                : 'text-gray-700 dark:text-gray-300 hover:bg-brand-orange/10 hover:text-brand-orange'}"
            href="/perfil/meus-dados"
        >
            <span
                class="material-symbols-outlined mr-3 {isActive(
                    '/perfil/meus-dados',
                )
                    ? 'text-brand-orange'
                    : 'text-gray-400 group-hover:text-brand-orange transition-colors'}"
                >person</span
            >
            Meus Dados
        </a>
        {#if role !== 'cliente'}
        <a
            class="group flex items-center px-3 py-2 text-sm font-medium rounded-md transition-all duration-200 {isActive(
                '/perfil',
            )
                ? 'bg-brand-orange/10 text-brand-orange'
                : 'text-gray-700 dark:text-gray-300 hover:bg-brand-orange/10 hover:text-brand-orange'}"
            href="/perfil"
        >
            <span
                class="material-symbols-outlined mr-3 {isActive('/perfil')
                    ? 'text-brand-orange'
                    : 'text-gray-400 group-hover:text-brand-orange transition-colors'}"
                >edit_calendar</span
            >
            Agenda Diária
        </a>
        {/if}
        <button
            on:click={handleLogout}
            class="group flex w-full items-center px-3 py-2 text-sm font-medium rounded-md transition-all duration-200 text-gray-700 dark:text-gray-300 hover:bg-red-50 dark:hover:bg-red-900/20 hover:text-red-600 dark:hover:text-red-400"
        >
            <span
                class="material-symbols-outlined mr-3 text-gray-400 group-hover:text-red-500 transition-colors"
                >logout</span
            >
            Sair
        </button>
    </div>

    <div class="border-t border-border-light dark:border-border-dark p-4">
        <div class="flex items-center">
            {#if $user?.image_url}
                <img
                    alt="User avatar"
                    class="h-8 w-8 rounded-full object-cover"
                    src={$user.image_url}
                />
            {:else}
                <div
                    class="h-8 w-8 rounded-full bg-orange-500 flex items-center justify-center flex-shrink-0"
                >
                    <span class="text-white text-xs font-bold">
                        {($user?.nome ?? "U").split(" ").map((n: string) => n[0]).slice(0, 2).join("").toUpperCase()}
                    </span>
                </div>
            {/if}
            <div class="ml-3 min-w-0">
                <p
                    class="text-sm font-medium text-gray-700 dark:text-gray-200 truncate"
                >
                    {$user?.nome ?? "Usuário"}
                </p>
                <p class="text-xs text-gray-500 dark:text-gray-400">
                    {roleLabel}
                </p>
            </div>
        </div>
    </div>
</aside>
{/if}
