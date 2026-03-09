<script lang="ts">
  import User from 'lucide-svelte/icons/user';
  import ArrowRight from 'lucide-svelte/icons/arrow-right';
  import type { TestimonialData } from '$lib/config/landingConfig';

  export let testimonials: TestimonialData[];
</script>

<section id="testimonials" class="py-32 bg-gray-50 dark:bg-black">
  <div class="container mx-auto px-4">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-24 items-start">
      <!-- Testimonials header -->
      <div class="md:sticky md:top-32">
        <span class="text-xs font-bold text-orange-500 uppercase tracking-[0.3em]">
          Testemunhos
        </span>
        <h2 class="text-5xl md:text-6xl lg:text-7xl font-bold mt-6 text-gray-900 dark:text-white leading-tight tracking-tight font-sans-bs">
          O que dizem nossas <br />
          <span class="text-orange-500 underline decoration-orange-500/20 underline-offset-8">
            Clientes
          </span>
        </h2>
        <p class="text-gray-600 dark:text-gray-400 mt-8 text-xl leading-relaxed max-w-md">
          A satisfação de quem confia em nossa arte é o nosso maior prêmio. Descubra por que somos referência em São Paulo.
        </p>

        <!-- Navigation buttons -->
        <div class="flex gap-4 mt-12">
          <button
            class="p-4 rounded-full border border-gray-300 dark:border-white/10 bg-transparent text-gray-900 dark:text-white hover:bg-orange-500 hover:text-white hover:border-orange-500 transition-all duration-300 active:scale-90"
            aria-label="Previous"
          >
            <ArrowRight class="w-5 h-5 rotate-180" />
          </button>
          <button
            class="p-4 rounded-full border border-gray-300 dark:border-white/10 bg-transparent text-gray-900 dark:text-white hover:bg-orange-500 hover:text-white hover:border-orange-500 transition-all duration-300 active:scale-90"
            aria-label="Next"
          >
            <ArrowRight class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Testimonials list -->
      <div class="flex flex-col gap-8">
        {#each testimonials as testimonial}
          <div
            class="bg-white dark:bg-zinc-900/50 border border-gray-200 dark:border-zinc-800/50 p-12 rounded-[3rem] relative transition-all duration-300 hover:border-orange-500/30 group shadow-lg dark:shadow-none"
          >
            <!-- Quote icon -->
            <svg
              class="absolute -top-6 -left-6 text-orange-500/10 w-32 h-32 -rotate-12 group-hover:text-orange-500/20 transition-colors duration-300"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M3 21c3 0 7-1 7-8V5c0-1.25-.756-2.017-2-2H4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2 1 0 1 0 1 1v1c0 1-1 2-2 2s-1 .008-1 1.031V20c0 1 0 1 1 1z"></path>
              <path d="M15 21c3 0 7-1 7-8V5c0-1.25-.757-2.017-2-2h-4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2h.75c0 2.25.25 4-2.75 4v3c0 1 0 1 1 1z"></path>
            </svg>

            <!-- Rating -->
            <div class="flex gap-1.5 mb-8">
              {#each Array(testimonial.rating) as _}
                <svg class="w-4 h-4 text-orange-500 fill-current" viewBox="0 0 24 24">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                </svg>
              {/each}
            </div>

            <!-- Testimonial text -->
            <p class="text-gray-700 dark:text-gray-200 mb-10 text-2xl font-medium italic leading-relaxed">
              "{testimonial.text}"
            </p>

            <!-- Author -->
            <div class="flex items-center gap-5">
              <div class="w-14 h-14 rounded-full bg-orange-500/10 border border-orange-500/20 flex items-center justify-center text-orange-500">
                <User class="w-6 h-6" />
              </div>
              <div>
                <h4 class="font-bold text-gray-900 dark:text-white text-lg">
                  {testimonial.author}
                </h4>
                <p class="text-gray-600 dark:text-gray-400 text-sm uppercase tracking-wider">
                  {testimonial.role}
                </p>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  </div>
</section>
