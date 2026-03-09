<script lang="ts">
  import { onMount } from "svelte";
  import Menu from 'lucide-svelte/icons/menu';
  import X from 'lucide-svelte/icons/x';
  import Search from 'lucide-svelte/icons/search';
  import Sun from 'lucide-svelte/icons/sun';
  import Moon from 'lucide-svelte/icons/moon';
  import { theme } from '$lib/stores/theme';

  let isMobileMenuOpen = false;
  let scrolled = false;

  function toggleMobileMenu() {
    isMobileMenuOpen = !isMobileMenuOpen;
  }

  function toggleTheme() {
    theme.update(t => t === 'dark' ? 'light' : 'dark');
  }

  onMount(() => {
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
      <a href="#home" class="flex items-center gap-2 group">
        <div class="w-9 h-9 bg-orange-500 rounded-xl flex items-center justify-center text-white font-bold text-base shadow-lg shadow-orange-500/25 transition-transform duration-300 group-hover:scale-105">
          B
        </div>
        <span class="text-base font-bold text-gray-900 dark:text-white font-sans-bs">
          Bella<span class="text-orange-500">Salon</span>
        </span>
      </a>

      <!-- Desktop Nav Links - Contained -->
      <div class="hidden md:flex items-center gap-1 bg-gray-100/50 dark:bg-white/5 backdrop-blur-sm px-5 py-2 rounded-full border border-gray-200/50 dark:border-white/5">
        <a
          href="#home"
          class="relative px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors duration-300 after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-orange-500 after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-300"
        >
          Início
        </a>
        <a
          href="#services"
          class="relative px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors duration-300 after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-orange-500 after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-300"
        >
          Serviços
        </a>
        <a
          href="#specialists"
          class="relative px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors duration-300 after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-orange-500 after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-300"
        >
          Especialistas
        </a>
        <a
          href="#testimonials"
          class="relative px-3 py-1.5 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors duration-300 after:absolute after:bottom-0 after:left-0 after:right-0 after:h-0.5 after:bg-orange-500 after:scale-x-0 hover:after:scale-x-100 after:transition-transform after:duration-300"
        >
          Depoimentos
        </a>
      </div>

      <!-- Desktop Actions -->
      <div class="hidden md:flex items-center gap-3">
        <button
          on:click={toggleTheme}
          class="w-9 h-9 rounded-full bg-gray-100/50 dark:bg-white/5 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-200/50 dark:hover:bg-white/10 transition-all duration-300 flex items-center justify-center"
          aria-label="Toggle theme"
        >
          {#if $theme === 'dark'}
            <Sun class="w-4 h-4" />
          {:else}
            <Moon class="w-4 h-4" />
          {/if}
        </button>
        <button
          class="w-9 h-9 rounded-full bg-gray-100/50 dark:bg-white/5 text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-white hover:bg-gray-200/50 dark:hover:bg-white/10 transition-all duration-300 flex items-center justify-center"
          aria-label="Search"
        >
          <Search class="w-4 h-4" />
        </button>
        <button class="px-5 py-2 rounded-full bg-orange-500 text-white text-sm font-semibold shadow-md hover:scale-105 hover:bg-orange-600 transition-all duration-300 hover:shadow-lg">
          Agendar Agora
        </button>
      </div>

      <!-- Mobile Actions -->
      <div class="flex md:hidden items-center gap-2">
        <button
          on:click={toggleTheme}
          class="w-9 h-9 rounded-full border border-gray-300 dark:border-white/10 bg-white dark:bg-transparent text-gray-700 dark:text-white hover:bg-gray-100 dark:hover:bg-white/5 transition-all duration-300 flex items-center justify-center"
          aria-label="Toggle theme"
        >
          {#if $theme === 'dark'}
            <Sun class="w-4 h-4" />
          {:else}
            <Moon class="w-4 h-4" />
          {/if}
        </button>
        <button
          on:click={toggleMobileMenu}
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
    <div
      class="md:hidden bg-white/95 dark:bg-black/95 backdrop-blur-3xl border-b border-gray-200 dark:border-white/10 p-4 flex flex-col gap-4 animate-slide-down"
    >
      <a
        href="#home"
        on:click={() => isMobileMenuOpen = false}
        class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300"
      >
        Início
      </a>
      <a
        href="#services"
        on:click={() => isMobileMenuOpen = false}
        class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300"
      >
        Serviços
      </a>
      <a
        href="#specialists"
        on:click={() => isMobileMenuOpen = false}
        class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300"
      >
        Especialistas
      </a>
      <a
        href="#testimonials"
        on:click={() => isMobileMenuOpen = false}
        class="text-base font-medium text-gray-700 dark:text-gray-300 hover:text-orange-500 py-2 border-b border-gray-200 dark:border-white/5 transition-colors duration-300"
      >
        Depoimentos
      </a>
      <button class="w-full mt-4 px-6 py-2.5 rounded-full bg-orange-500 text-white text-sm font-semibold shadow-md">
        Agendar Agora
      </button>
    </div>
  {/if}
</nav>
