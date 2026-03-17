<script lang="ts">
  import CircleCheck from "lucide-svelte/icons/circle-check";
  import type { SpecialistData } from "$lib/config/landingConfig";

  let { specialists }: { specialists: SpecialistData[] } = $props();
</script>

<section id="specialists" class="py-32 bg-white dark:bg-black">
  <div class="container mx-auto px-4">
    <!-- Section header -->
    <div class="text-center mb-20">
      <span
        class="text-xs font-bold font-montserrat text-orange-500 uppercase tracking-[0.3em]"
      >
        Nossa Equipe
      </span>
      <h2
        class="text-5xl md:text-6xl lg:text-7xl font-bold font-cormorant mt-6 text-slate-800 dark:text-white leading-tight tracking-tight"
      >
        Conheca Nossos <br />
        <span class="text-gradient-warm italic">
          Especialistas
        </span>
      </h2>
      <p
        class="text-gray-500 dark:text-gray-400 mt-8 text-lg font-montserrat leading-relaxed max-w-xl mx-auto"
      >
        Profissionais altamente qualificados e apaixonados pela arte da beleza
      </p>
    </div>

    <!-- Specialists grid — more spacious layout -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-10 max-w-5xl mx-auto">
      {#each specialists as specialist, i}
        <div
          class="group relative p-10 rounded-3xl bg-white dark:bg-zinc-900/50 border border-pink-100 dark:border-zinc-800/50 flex flex-col items-center text-center overflow-hidden shadow-lg shadow-pink-100/50 dark:shadow-none card-hover-lift"
          style="animation-delay: {i * 150}ms;"
        >
          <!-- Soft background blob on hover -->
          <div
            class="absolute inset-0 bg-gradient-to-br from-orange-50/60 to-pink-50/60 dark:from-orange-900/10 dark:to-pink-900/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none rounded-3xl"
          ></div>

          <!-- Image wrapper -->
          <div class="relative mb-8 z-10">
            <!-- Glow on hover -->
            <div
              class="absolute inset-0 bg-gradient-to-br from-orange-400/30 to-pink-400/20 rounded-2xl blur-[2.5rem] scale-110 opacity-0 group-hover:opacity-100 transition-opacity duration-400"
            ></div>

            <!-- Square image with rounded corners (rounded-2xl) -->
            <div
              class="relative p-1 bg-gradient-to-tr from-orange-400/60 via-pink-300/40 to-transparent rounded-2xl"
            >
              <img
                src={specialist.imageUrl}
                alt={specialist.name}
                class="w-36 h-36 rounded-2xl object-cover transition-transform duration-400 group-hover:scale-[1.03]"
                loading="lazy"
              />
            </div>
          </div>

          <!-- Content -->
          <h3 class="text-2xl font-bold font-cormorant text-slate-800 dark:text-white mb-1 z-10">
            {specialist.name}
          </h3>
          <p
            class="text-orange-500 text-xs font-semibold font-montserrat uppercase tracking-widest mb-5 z-10"
          >
            {specialist.role}
          </p>

          <!-- Specialty pill badge -->
          <div
            class="flex items-center gap-2 text-gray-600 dark:text-gray-400 text-xs font-montserrat px-5 py-2 rounded-full bg-pink-50 dark:bg-zinc-800/50 border border-pink-200/60 dark:border-zinc-700/50 z-10"
          >
            <CircleCheck class="w-3.5 h-3.5 text-orange-500 shrink-0" />
            <span>{specialist.specialty}</span>
          </div>
        </div>
      {/each}
    </div>
  </div>
</section>
