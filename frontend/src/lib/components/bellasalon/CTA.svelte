<script lang="ts">
  let { title, description, buttonText }: { title: string; description: string; buttonText: string } = $props();
</script>

<section class="py-32 relative overflow-hidden">
  <!-- Gradient background -->
  <div class="absolute inset-0 bg-gradient-to-br from-orange-500 via-orange-400 to-pink-500"></div>

  <!-- Dot pattern -->
  <div
    class="absolute inset-0 opacity-[0.08] pointer-events-none"
    style="background-image: radial-gradient(circle at 2px 2px, black 1.5px, transparent 0); background-size: 48px 48px;"
  ></div>

  <!-- Soft blob accents -->
  <div class="absolute top-0 right-0 w-96 h-96 bg-white/10 rounded-full blur-[80px] pointer-events-none"></div>
  <div class="absolute bottom-0 left-0 w-80 h-80 bg-pink-600/20 rounded-full blur-[60px] pointer-events-none"></div>

  <!-- Content -->
  <div class="container mx-auto px-4 relative z-10">
    <div class="text-center flex flex-col gap-10 items-center">
      <!-- Title in Cormorant -->
      <h2 class="text-6xl md:text-7xl lg:text-8xl font-bold font-cormorant text-white leading-[1.1] tracking-tight italic">
        {#each title.split(' ') as word, i}
          {word}{#if (i + 1) % 2 === 0}<br />{:else}&nbsp;{/if}
        {/each}
      </h2>

      <!-- Description -->
      <p class="text-xl md:text-2xl font-montserrat font-medium text-white/80 max-w-2xl">
        {description}
      </p>

      <!-- CTA Button -->
      <a
        href="/listagem"
        class="bg-white text-orange-600 px-14 py-5 text-base font-bold font-montserrat uppercase tracking-wider rounded-full shadow-[0_20px_50px_rgba(0,0,0,0.2)] hover:shadow-[0_25px_60px_rgba(0,0,0,0.3)] hover:scale-105 hover:-translate-y-1 hover:bg-orange-50 transition-all duration-250 inline-block cursor-pointer"
      >
        {buttonText}
      </a>
    </div>
  </div>
</section>
