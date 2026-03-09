<script lang="ts">
  export let title: string;
  export let description: string;
  export let buttonText: string;
</script>

<section class="py-32 bg-[hsl(var(--bs-primary))] relative overflow-hidden">
  <!-- Pattern background -->
  <div
    class="absolute inset-0 opacity-10 pointer-events-none"
    style="background-image: radial-gradient(circle at 2px 2px, black 1px, transparent 0); background-size: 60px 60px;"
  ></div>

  <!-- Content -->
  <div class="container mx-auto px-4 relative z-10">
    <div class="text-center flex flex-col gap-10 items-center">
      <!-- Title -->
      <h2 class="text-6xl md:text-7xl lg:text-8xl font-black text-black leading-[1.4] tracking-tighter uppercase italic font-serif-bs">
        {#each title.split(' ') as word, i}
          {word}{#if (i + 1) % 2 === 0}<br />{:else}&nbsp;{/if}
        {/each}
      </h2>

      <!-- Description -->
      <p class="text-2xl md:text-3xl font-bold text-black/70 max-w-2xl uppercase tracking-tight">
        {description}
      </p>

      <!-- CTA Button -->
      <a
        href="/listagem"
        class="bg-black text-white px-16 py-5 text-xl font-black uppercase tracking-wider rounded-full shadow-[0_20px_50px_rgba(0,0,0,0.3)] hover:shadow-[0_25px_60px_rgba(0,0,0,0.4)] hover:scale-105 hover:bg-zinc-900 transition-all duration-300 inline-block"
      >
        {buttonText}
      </a>
    </div>
  </div>
</section>
