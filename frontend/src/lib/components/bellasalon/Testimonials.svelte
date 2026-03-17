<script lang="ts">
  import User from 'lucide-svelte/icons/user';
  import ArrowRight from 'lucide-svelte/icons/arrow-right';
  import type { TestimonialData } from '$lib/config/landingConfig';

  let { testimonials }: { testimonials: TestimonialData[] } = $props();

  let currentPage = $state(0);
  const itemsPerPage = 3;

  let totalPages = $derived(Math.ceil(testimonials.length / itemsPerPage));
  let visibleTestimonials = $derived(testimonials.slice(
    currentPage * itemsPerPage,
    (currentPage + 1) * itemsPerPage
  ));

  function prevPage() {
    currentPage = (currentPage - 1 + totalPages) % totalPages;
  }

  function nextPage() {
    currentPage = (currentPage + 1) % totalPages;
  }
</script>

<section id="testimonials" class="py-32 bg-[#FDF2F8] dark:bg-black">
  <div class="container mx-auto px-4">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-24 items-start">
      <!-- Testimonials header — sticky -->
      <div class="md:sticky md:top-32">
        <span class="text-xs font-bold font-montserrat text-orange-500 uppercase tracking-[0.3em]">
          Testemunhos
        </span>
        <h2 class="text-5xl md:text-6xl lg:text-7xl font-bold font-cormorant mt-6 text-slate-800 dark:text-white leading-tight tracking-tight">
          O que dizem nossas <br />
          <span class="text-gradient-warm italic">
            Clientes
          </span>
        </h2>
        <p class="text-gray-500 dark:text-gray-400 mt-8 text-lg font-montserrat leading-relaxed max-w-md">
          A satisfacao de quem confia em nossa arte e o nosso maior premio. Descubra por que somos referencia em Sao Paulo.
        </p>

        <!-- Navigation buttons -->
        <div class="flex items-center gap-4 mt-12">
          <button
            onclick={prevPage}
            class="p-4 rounded-full border border-pink-200 dark:border-white/10 bg-white dark:bg-transparent text-slate-800 dark:text-white hover:bg-gradient-to-r hover:from-orange-500 hover:to-pink-500 hover:text-white hover:border-transparent transition-all duration-250 active:scale-90 cursor-pointer shadow-sm"
            aria-label="Previous"
          >
            <ArrowRight class="w-5 h-5 rotate-180" />
          </button>
          <span class="text-sm font-montserrat text-gray-500 dark:text-gray-400">
            {currentPage + 1} / {totalPages}
          </span>
          <button
            onclick={nextPage}
            class="p-4 rounded-full border border-pink-200 dark:border-white/10 bg-white dark:bg-transparent text-slate-800 dark:text-white hover:bg-gradient-to-r hover:from-orange-500 hover:to-pink-500 hover:text-white hover:border-transparent transition-all duration-250 active:scale-90 cursor-pointer shadow-sm"
            aria-label="Next"
          >
            <ArrowRight class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- Testimonials list -->
      <div class="flex flex-col gap-6">
        {#each visibleTestimonials as testimonial}
          <div
            class="bg-white dark:bg-zinc-900/50 border border-pink-100 dark:border-zinc-800/50 p-10 rounded-3xl relative transition-all duration-300 hover:border-orange-300/60 dark:hover:border-orange-500/30 hover:shadow-xl hover:shadow-pink-100/40 dark:hover:shadow-black/30 hover:-translate-y-1 group shadow-md shadow-pink-50 dark:shadow-none"
          >
            <!-- Large quote watermark -->
            <svg
              class="absolute -top-5 -left-5 text-orange-500/8 w-28 h-28 -rotate-12 group-hover:text-orange-500/15 transition-colors duration-300"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M3 21c3 0 7-1 7-8V5c0-1.25-.756-2.017-2-2H4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2 1 0 1 0 1 1v1c0 1-1 2-2 2s-1 .008-1 1.031V20c0 1 0 1 1 1z"></path>
              <path d="M15 21c3 0 7-1 7-8V5c0-1.25-.757-2.017-2-2h-4c-1.25 0-2 .75-2 1.972V11c0 1.25.75 2 2 2h.75c0 2.25.25 4-2.75 4v3c0 1 0 1 1 1z"></path>
            </svg>

            <!-- Rating stars -->
            <div class="flex gap-1 mb-6">
              {#each Array(testimonial.rating) as _}
                <svg class="w-4 h-4 text-orange-500 fill-current" viewBox="0 0 24 24">
                  <polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon>
                </svg>
              {/each}
            </div>

            <!-- Testimonial text in Cormorant italic -->
            <p class="text-slate-700 dark:text-gray-200 mb-8 text-xl font-cormorant italic leading-relaxed">
              "{testimonial.text}"
            </p>

            <!-- Author -->
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-full bg-gradient-to-br from-orange-500/15 to-pink-500/15 border border-orange-200/60 dark:border-orange-500/20 flex items-center justify-center text-orange-500">
                <User class="w-5 h-5" />
              </div>
              <div>
                <h4 class="font-bold font-cormorant text-slate-800 dark:text-white text-lg">
                  {testimonial.author}
                </h4>
                <p class="text-gray-500 dark:text-gray-400 text-xs font-montserrat uppercase tracking-wider">
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
