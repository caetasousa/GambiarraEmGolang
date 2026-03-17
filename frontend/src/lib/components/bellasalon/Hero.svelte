<script lang="ts">
  import Sparkles from "lucide-svelte/icons/sparkles";
  import ArrowRight from "lucide-svelte/icons/arrow-right";
  import Star from "lucide-svelte/icons/star";
  import type { HeroData } from "$lib/config/landingConfig";

  let { data }: { data: HeroData } = $props();
</script>

<section
  id="home"
  class="relative min-h-screen flex items-center pt-20 overflow-hidden bg-gradient-to-br from-orange-50 via-pink-50/30 to-white dark:bg-black dark:from-black dark:via-black dark:to-black"
>
  <!-- Decorative floating blobs -->
  <div
    class="absolute top-16 right-16 w-80 h-80 bg-gradient-to-br from-orange-300/25 to-pink-300/20 rounded-full blur-[100px] pointer-events-none animate-float-blob"
  ></div>
  <div
    class="absolute bottom-24 left-8 w-64 h-64 bg-gradient-to-tr from-violet-300/15 to-pink-300/15 rounded-full blur-[80px] pointer-events-none animate-float-blob"
    style="animation-delay: 3s;"
  ></div>
  <div
    class="absolute top-1/3 left-1/4 w-48 h-48 bg-orange-400/10 rounded-full blur-[60px] pointer-events-none animate-float-blob"
    style="animation-delay: 6s;"
  ></div>

  <!-- Subtle mesh gradient overlay -->
  <div
    class="absolute inset-0 bg-gradient-to-br from-transparent via-pink-50/20 to-orange-50/30 dark:via-transparent dark:to-transparent pointer-events-none"
  ></div>

  <div class="container mx-auto px-4 relative z-10">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
      <!-- Hero text -->
      <div class="flex flex-col gap-8">
        <!-- Badge -->
        <div
          class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-white/80 dark:bg-white/5 border border-pink-200/60 dark:border-white/10 w-fit shadow-sm backdrop-blur-sm animate-fade-in-up-1"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-gradient-to-r from-orange-500 to-pink-500"></span>
          <span
            class="text-xs font-semibold font-montserrat text-gray-600 dark:text-gray-400 uppercase tracking-wider"
          >
            {data.badge}
          </span>
        </div>

        <!-- Title with Cormorant -->
        <h1
          class="text-6xl md:text-7xl lg:text-8xl font-bold leading-[1.05] tracking-tight text-slate-800 dark:text-white font-cormorant animate-fade-in-up-2"
        >
          {data.title} <br />
          <span class="text-gradient-warm italic">
            {data.titleAccent}
          </span>
        </h1>

        <!-- Description -->
        <p
          class="text-lg text-gray-600 dark:text-gray-400 max-w-lg leading-relaxed font-montserrat font-normal animate-fade-in-up-3"
        >
          {data.description}
        </p>

        <!-- Buttons -->
        <div class="flex flex-col sm:flex-row gap-4 pt-2 animate-fade-in-up-4">
          <a
            href="/listagem"
            class="inline-flex items-center justify-center gap-2 px-8 py-3.5 text-sm font-semibold font-montserrat rounded-full bg-gradient-to-r from-orange-500 to-orange-600 text-white shadow-lg shadow-orange-500/30 hover:scale-105 hover:-translate-y-0.5 transition-all duration-200 hover:shadow-orange-500/40 cursor-pointer"
          >
            {data.primaryButtonText}
          </a>
          <a
            href="/listagem"
            class="inline-flex items-center justify-center gap-2 px-8 py-3.5 text-sm font-semibold font-montserrat rounded-full bg-white/80 dark:bg-transparent text-gray-700 dark:text-white border border-pink-200 dark:border-white/10 hover:bg-white hover:border-orange-300 dark:hover:bg-white/5 transition-all duration-200 cursor-pointer backdrop-blur-sm shadow-sm"
          >
            {data.secondaryButtonText}
            <ArrowRight class="w-4 h-4" />
          </a>
        </div>

        <!-- Stats -->
        <div
          class="flex items-center gap-8 pt-8 border-t border-pink-100 dark:border-white/10 animate-fade-in-up-5"
        >
          {#each data.stats as stat}
            <div class="flex flex-col">
              <h4 class="text-3xl font-bold font-cormorant text-slate-800 dark:text-white">
                {stat.value}
              </h4>
              <p class="text-sm font-montserrat text-gray-500 dark:text-gray-400">
                {stat.label}
              </p>
            </div>
          {/each}
        </div>
      </div>

      <!-- Hero image -->
      <div class="relative h-[500px] md:h-[680px] w-full fade-in-scale">
        <div
          class="absolute inset-0 bg-gradient-to-b from-transparent via-orange-50/10 to-orange-50/60 dark:to-black/60 z-10"
        ></div>
        <img
          src={data.imageUrl}
          alt="Salon Experience"
          class="w-full h-full object-cover rounded-3xl md:rounded-[2.5rem] shadow-2xl shadow-orange-200/30 dark:shadow-black/50"
          style="mask-image: linear-gradient(to bottom, black 80%, transparent 100%); -webkit-mask-image: linear-gradient(to bottom, black 80%, transparent 100%);"
          loading="eager"
        />

        <!-- Floating review card -->
        <div class="hidden md:block absolute bottom-28 -left-10 z-20">
          <div
            class="p-5 rounded-2xl max-w-xs shadow-2xl shadow-pink-200/30 dark:shadow-black/50 bg-white/90 dark:bg-zinc-900/80 backdrop-blur-md border border-pink-100 dark:border-zinc-800"
          >
            <div class="flex items-center gap-3 mb-3">
              <div
                class="w-9 h-9 rounded-full bg-gradient-to-br from-orange-500/20 to-pink-500/20 flex items-center justify-center text-orange-500"
              >
                <Sparkles class="w-4 h-4" />
              </div>
              <div>
                <h5 class="font-bold text-sm font-cormorant text-gray-900 dark:text-white">
                  {data.floatingCard.title}
                </h5>
                <p class="text-xs font-montserrat text-gray-500 dark:text-gray-400">
                  {data.floatingCard.subtitle}
                </p>
              </div>
            </div>
            <p class="text-xs font-montserrat text-gray-700 dark:text-gray-300 italic mb-3">
              "{data.floatingCard.text}"
            </p>
            <div class="flex gap-0.5 text-orange-500">
              {#each Array(data.floatingCard.rating) as _}
                <Star class="w-3 h-3 fill-current" />
              {/each}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</section>
