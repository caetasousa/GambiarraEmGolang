<script lang="ts">
  import Menu from 'lucide-svelte/icons/menu';
  import X from 'lucide-svelte/icons/x';
  import Sun from 'lucide-svelte/icons/sun';
  import Moon from 'lucide-svelte/icons/moon';
  import LogOut from 'lucide-svelte/icons/log-out';
  import { theme, toggleTheme } from '$lib/stores/theme.svelte';
  import { auth, logout } from '$lib/stores/auth.svelte';
  import { goto } from '$app/navigation';

  let isMobileMenuOpen = $state(false);
  let scrolled = $state(false);

  function toggleMobileMenu() {
    isMobileMenuOpen = !isMobileMenuOpen;
  }

  function getPanelRoute(role: string | undefined): string {
    switch (role) {
      case 'admin': return '/prestadores/editar';
      case 'prestador': return '/prestadores/agenda';
      case 'cliente': return '/clientes/meus-agendamentos';
      default: return '/login';
    }
  }

  function handleLogout() {
    logout();
    goto('/login');
  }

  $effect(() => {
    const handleScroll = () => {
      scrolled = window.scrollY > 20;
    };
    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  });
</script>

<nav
  class="fixed top-0 left-0 right-0 z-50 transition-all duration-300 {scrolled ? 'bg-white/80 dark:bg-black/80 backdrop-blur-md py-4' : 'py-6'}"
>
  <div class="container mx-auto px-4">
    <div class="flex items-center justify-between gap-8">
      <!-- Logo -->
      <a href="#home" class="flex items-center gap-2 group cursor-pointer">
        <div class="w-9 h-9 bg-gradient-to-br from-orange-500 to-pink-500 rounded-xl flex items-center justify-center text-white font-bold text-base shadow-lg shadow-orange-500/30 transition-all duration-300 group-hover:scale-110 group-hover:shadow-orange-500/50">
          B
        </div>
        <span class="text-lg font-bold text-gray-900 dark:text-white font-cormorant tracking-wide">
          Bella<span class="text-gradient-warm">Vita</span>
        </span>
      </a>

      <!-- Desktop Nav Links -->
      <div class="hidden md:flex items-center gap-1 bg-white/60 dark:bg-white/5 backdrop-blur-sm px-5 py-2 rounded-full border border-pink-100/80 dark:border-white/5 shadow-sm">
        <a href="#home" class="relative px-3 py-1.5 text-sm font-medium font-montserrat text-gray-600 dark:text-gray-300 hover:text-orange-500 dark:hover:text-orange-400 transition-colors duration-200 after:absolute after:bottom-0 after:left-2 after:right-2 after:h-0.5 after:bg-gradient-to-r after:from-orange-500 after:to-pink-500 after:rounded-full after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-250">
          Inicio
        </a>
        <a href="#services" class="relative px-3 py-1.5 text-sm font-medium font-montserrat text-gray-600 dark:text-gray-300 hover:text-orange-500 dark:hover:text-orange-400 transition-colors duration-200 after:absolute after:bottom-0 after:left-2 after:right-2 after:h-0.5 after:bg-gradient-to-r after:from-orange-500 after:to-pink-500 after:rounded-full after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-250">
          Servicos
        </a>
        <a href="#specialists" class="relative px-3 py-1.5 text-sm font-medium font-montserrat text-gray-600 dark:text-gray-300 hover:text-orange-500 dark:hover:text-orange-400 transition-colors duration-200 after:absolute after:bottom-0 after:left-2 after:right-2 after:h-0.5 after:bg-gradient-to-r after:from-orange-500 after:to-pink-500 after:rounded-full after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-250">
          Especialistas
        </a>
        <a href="#testimonials" class="relative px-3 py-1.5 text-sm font-medium font-montserrat text-gray-600 dark:text-gray-300 hover:text-orange-500 dark:hover:text-orange-400 transition-colors duration-200 after:absolute after:bottom-0 after:left-2 after:right-2 after:h-0.5 after:bg-gradient-to-r after:from-orange-500 after:to-pink-500 after:rounded-full after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-250">
          Depoimentos
        </a>
      </div>

      <!-- Desktop Actions -->
      <div class="hidden md:flex items-center gap-3">
        <button
          onclick={toggleTheme}
          class="w-9 h-9 rounded-full bg-gray-100/50 dark:bg-white/5 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-200/50 dark:hover:bg-white/10 transition-all duration-300 flex items-center justify-center"
          aria-label="Toggle theme"
        >
          {#if theme.value === 'dark'}
            <Sun class="w-4 h-4" />
          {:else}
            <Moon class="w-4 h-4" />
          {/if}
        </button>

        {#if auth.isAuthenticated}
          <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
            {auth.user?.nome}
          </span>
          <a
            href={getPanelRoute(auth.user?.role)}
            class="px-5 py-2 rounded-full bg-orange-500 text-white text-sm font-semibold shadow-md hover:scale-105 hover:bg-orange-600 transition-all duration-300 hover:shadow-lg"
          >
            Meu Painel
          </a>
          <button
            onclick={handleLogout}
            class="w-9 h-9 rounded-full bg-gray-100/50 dark:bg-white/5 text-gray-600 dark:text-gray-400 hover:text-red-500 hover:bg-gray-200/50 dark:hover:bg-white/10 transition-all duration-300 flex items-center justify-center"
            aria-label="Sair"
          >
            <LogOut class="w-4 h-4" />
          </button>
        {:else}
          <a href="/login" class="px-4 py-2 rounded-full text-sm font-medium font-montserrat text-gray-700 dark:text-gray-300 hover:text-orange-500 dark:hover:text-orange-400 transition-colors duration-200 cursor-pointer">
            Entrar
          </a>
          <a href="/clientes/cadastro" class="px-5 py-2 rounded-full bg-gradient-to-r from-orange-500 to-orange-600 text-white text-sm font-semibold font-montserrat shadow-md shadow-orange-500/25 hover:scale-105 hover:shadow-lg hover:shadow-orange-500/30 transition-all duration-200 cursor-pointer">
            Cadastre-se
          </a>
        {/if}
      </div>

      <!-- Mobile Actions -->
      <div class="flex md:hidden items-center gap-2">
        <button
          onclick={toggleTheme}
          class="w-9 h-9 rounded-full border border-gray-300 dark:border-white/10 bg-white dark:bg-transparent text-gray-700 dark:text-white hover:bg-gray-100 dark:hover:bg-white/5 transition-all duration-300 flex items-center justify-center"
          aria-label="Toggle theme"
        >
          {#if theme.value === 'dark'}
            <Sun class="w-4 h-4" />
          {:else}
            <Moon class="w-4 h-4" />
          {/if}
        </button>
        <button
          onclick={toggleMobileMenu}
          class="w-9 h-9 rounded-full bg-transparent text-gray-700 dark:text-white hover:bg-gray-100 dark:hover:bg-white/5 transition-all duration-300 flex items-center justify-center"
          aria-label="Toggle menu"
        >
          {#if isMobileMenuOpen}
            <X class="w-4 h-4" />
          {:else}
            <Menu class="w-4 h-4" />
          {/if}
        </button>
      </div>
    </div>
  </div>

  <!-- Mobile Menu -->
  {#if isMobileMenuOpen}
    <div class="md:hidden bg-white/95 dark:bg-black/95 backdrop-blur-3xl border-b border-gray-200 dark:border-white/10 p-4 flex flex-col gap-4 animate-slide-down">
      <a href="#home" onclick={() => isMobileMenuOpen = false} class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300">
        Inicio
      </a>
      <a href="#services" onclick={() => isMobileMenuOpen = false} class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300">
        Servicos
      </a>
      <a href="#specialists" onclick={() => isMobileMenuOpen = false} class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300">
        Especialistas
      </a>
      <a href="#testimonials" onclick={() => isMobileMenuOpen = false} class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300">
        Depoimentos
      </a>

      {#if auth.isAuthenticated}
        <div class="flex items-center gap-2 py-2 border-b border-gray-200 dark:border-white/5">
          <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
            {auth.user?.nome}
          </span>
        </div>
        <a
          href={getPanelRoute(auth.user?.role)}
          onclick={() => isMobileMenuOpen = false}
          class="w-full mt-2 px-6 py-2.5 rounded-full bg-orange-500 text-white text-sm font-semibold shadow-md text-center"
        >
          Meu Painel
        </a>
        <button
          onclick={() => { isMobileMenuOpen = false; handleLogout(); }}
          class="w-full px-6 py-2.5 rounded-full border border-gray-300 dark:border-white/10 text-gray-700 dark:text-gray-300 text-sm font-semibold hover:text-red-500 transition-colors duration-300"
        >
          Sair
        </button>
      {:else}
        <a href="/login" onclick={() => isMobileMenuOpen = false} class="w-full mt-2 px-6 py-2.5 rounded-full border border-gray-300 dark:border-white/10 text-gray-700 dark:text-gray-300 text-sm font-semibold text-center">
          Entrar
        </a>
        <a href="/clientes/cadastro" onclick={() => isMobileMenuOpen = false} class="w-full px-6 py-2.5 rounded-full bg-orange-500 text-white text-sm font-semibold shadow-md text-center">
          Cadastre-se
        </a>
      {/if}
    </div>
  {/if}
</nav>
