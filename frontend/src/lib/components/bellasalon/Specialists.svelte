<script lang="ts">
  import CircleCheck from "lucide-svelte/icons/circle-check";
  import type { SpecialistData } from "$lib/config/landingConfig";

  export let specialists: SpecialistData[];
</script>

<section id="specialists" class="py-32 bg-white dark:bg-black">
  <div class="container mx-auto px-4">
    <!-- Section header -->
    <div class="text-center mb-16">
      <span
        class="text-xs font-bold text-orange-500 uppercase tracking-[0.3em]"
      >
        Nossa Equipe
      </span>
      <h2
        class="text-5xl md:text-6xl lg:text-7xl font-bold mt-6 text-gray-900 dark:text-white leading-tight tracking-tight font-sans-bs"
      >
        Conheça Nossos <br />
        <span
          class="text-orange-500 underline decoration-orange-500/20 underline-offset-8"
        >
          Especialistas
        </span>
      </h2>
      <p
        class="text-gray-600 dark:text-gray-400 mt-8 text-xl leading-relaxed max-w-2xl mx-auto"
      >
        Profissionais altamente qualificados e apaixonados pela arte da beleza
      </p>
    </div>

    <!-- Specialists grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-12 max-w-6xl mx-auto">
      {#each specialists as specialist}
        <div
          class="p-12 rounded-[3rem] bg-white dark:bg-zinc-900/50 border border-gray-200 dark:border-zinc-800/50 flex flex-col items-center text-center transition-all duration-300 hover:border-orange-500/30 hover:-translate-y-2 cursor-pointer relative overflow-hidden group shadow-lg dark:shadow-none"
        >
          <!-- Image wrapper -->
          <div class="relative mb-8">
            <!-- Blur effect on hover -->
            <div
              class="absolute inset-0 bg-orange-500/20 rounded-full blur-[2rem] scale-125 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
            ></div>

            <!-- Image border with gradient -->
            <div
              class="relative p-1 bg-gradient-to-tr from-orange-500/50 to-transparent rounded-full"
            >
              <img
                src={specialist.imageUrl}
                alt={specialist.name}
                class="w-40 h-40 rounded-full object-cover transition-all duration-500"
                loading="lazy"
              />
            </div>
          </div>

          <!-- Content -->
          <h3 class="text-2xl font-bold text-gray-900 dark:text-white mb-1">
            {specialist.name}
          </h3>
          <p
            class="text-orange-500 text-sm font-semibold uppercase tracking-wider mb-6"
          >
            {specialist.role}
          </p>

          <!-- Specialty badge -->
          <div
            class="flex items-center gap-3 text-gray-600 dark:text-gray-400 text-xs px-6 py-2.5 rounded-full bg-gray-100 dark:bg-zinc-800/50 border border-gray-200 dark:border-zinc-700/50"
          >
            <CircleCheck class="w-4 h-4 text-orange-500" />
            <span>{specialist.specialty}</span>
          </div>
        </div>
      {/each}
    </div>
  </div>
</section>
