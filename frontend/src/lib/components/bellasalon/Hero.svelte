<script lang="ts">
  import Sparkles from "lucide-svelte/icons/sparkles";
  import ArrowRight from "lucide-svelte/icons/arrow-right";
  import Star from "lucide-svelte/icons/star";
  import type { HeroData } from "$lib/config/landingConfig";

  export let data: HeroData;
</script>

<section
  id="home"
  class="relative min-h-screen flex items-center pt-20 overflow-hidden bg-white dark:bg-black"
>
  <!-- Background gradients -->
  <div
    class="absolute top-0 right-0 w-3/5 h-full bg-gradient-to-l from-orange-500/10 to-transparent pointer-events-none"
  ></div>
  <div
    class="absolute top-20 right-20 w-96 h-96 bg-orange-500/20 rounded-full blur-[120px] pointer-events-none"
  ></div>

  <div class="container mx-auto px-4 relative z-10">
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
      <!-- Hero text -->
      <div class="flex flex-col gap-8 fade-in-up">
        <!-- Badge -->
        <div
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-orange-50 dark:bg-white/5 border-0 dark:border dark:border-white/10 w-fit"
        >
          <span class="w-1.5 h-1.5 rounded-full bg-orange-500"></span>
          <span
            class="text-xs font-semibold text-gray-600 dark:text-gray-400 uppercase tracking-wider"
          >
            {data.badge}
          </span>
        </div>

        <!-- Title -->
        <h1
          class="text-5xl md:text-6xl lg:text-7xl font-bold leading-tight tracking-tight text-gray-900 dark:text-white font-sans-bs"
        >
          {data.title} <br />
          <span class="text-orange-500 font-bold">
            {data.titleAccent}
          </span>
        </h1>

        <!-- Description -->
        <p
          class="text-xl text-gray-600 dark:text-gray-400 max-w-lg leading-relaxed font-normal"
        >
          {data.description}
        </p>

        <!-- Buttons -->
        <div class="flex flex-col sm:flex-row gap-4 pt-4">
          <a
            href="/#"
            class="inline-flex items-center justify-center gap-2 px-7 py-3.5 text-sm font-semibold rounded-full bg-orange-500 text-white shadow-md hover:scale-105 transition-all duration-300 hover:shadow-lg"
          >
            {data.primaryButtonText}
          </a>
          <a
            href="/listagem"
            class="inline-flex items-center justify-center gap-2 px-7 py-3.5 text-sm font-semibold rounded-full bg-white dark:bg-transparent text-gray-700 dark:text-white border border-gray-300 dark:border-white/10 hover:bg-gray-100 dark:hover:bg-white/5 transition-all duration-300"
          >
            {data.secondaryButtonText}
            <ArrowRight class="w-4 h-4" />
          </a>
        </div>

        <!-- Stats -->
        <div
          class="flex items-center gap-8 pt-8 border-t border-gray-300 dark:border-white/10"
        >
          {#each data.stats as stat}
            <div class="flex flex-col">
              <h4 class="text-3xl font-bold text-gray-900 dark:text-white">
                {stat.value}
              </h4>
              <p class="text-sm text-gray-600 dark:text-gray-400">
                {stat.label}
              </p>
            </div>
          {/each}
        </div>
      </div>

      <!-- Hero image -->
      <div class="relative h-[500px] md:h-[700px] w-full fade-in-scale">
        <div
          class="absolute inset-0 bg-gradient-to-b from-transparent via-[hsla(var(--bs-background),0.2)] to-[hsl(var(--bs-background))] z-10"
        ></div>
        <img
          src={data.imageUrl}
          alt="Salon Experience"
          class="w-full h-full object-cover rounded-2xl md:rounded-[2.5rem] shadow-2xl"
          style="mask-image: linear-gradient(to bottom, black 85%, transparent 100%); -webkit-mask-image: linear-gradient(to bottom, black 85%, transparent 100%);"
          loading="eager"
        />

        <!-- Floating card -->
        <div class="hidden md:block absolute bottom-32 -left-12 z-20">
          <div
            class="p-5 rounded-2xl max-w-xs shadow-2xl bg-orange-50/50 dark:bg-zinc-900/80 backdrop-blur-md border border-orange-50 dark:border-zinc-800"
          >
            <div class="flex items-center gap-3 mb-2.5">
              <div
                class="w-9 h-9 rounded-full bg-orange-500/20 flex items-center justify-center text-orange-500"
              >
                <Sparkles class="w-4 h-4" />
              </div>
              <div>
                <h5 class="font-bold text-sm text-gray-900 dark:text-white">
                  {data.floatingCard.title}
                </h5>
                <p class="text-xs text-gray-600 dark:text-gray-400">
                  {data.floatingCard.subtitle}
                </p>
              </div>
            </div>
            <p class="text-xs text-gray-700 dark:text-gray-300 italic mb-2.5">
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
