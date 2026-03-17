<script lang="ts">
  import ArrowRight from 'lucide-svelte/icons/arrow-right';
  import { landingConfig } from '$lib/config/landingConfig';

  const serviceImages: Record<string, string> = {
    'scissors': 'https://images.unsplash.com/photo-1560066984-138dadb4c035?auto=format&fit=crop&q=80&w=800',
    'sparkles': 'https://images.unsplash.com/photo-1604654894610-df63bc536371?auto=format&fit=crop&q=80&w=800',
    'calendar': 'https://images.unsplash.com/photo-1604654894610-df63bc536371?auto=format&fit=crop&q=80&w=800',
    'star': 'https://images.unsplash.com/photo-1544161515-4ab6ce6db874?auto=format&fit=crop&q=80&w=800'
  };

  const serviceCategoryLabels: Record<string, string> = {
    'scissors': 'Capilar',
    'sparkles': 'Coloracao',
    'calendar': 'Estetica',
    'star': 'Bem-Estar'
  };

  // Decorative large numbers for each card
  const cardNumbers = ['01', '02', '03'];
</script>

<section id="services" class="py-32 bg-[#FDF2F8] dark:bg-black">
  <div class="container mx-auto px-4">
    <!-- Section header -->
    <div class="mb-16">
      <span class="text-xs font-bold font-montserrat text-orange-500 uppercase tracking-[0.3em]">
        Nossos Servicos
      </span>
      <h2 class="text-5xl md:text-6xl lg:text-7xl font-bold font-cormorant mt-6 text-slate-800 dark:text-white leading-tight tracking-tight">
        Experiencias de <br />
        <span class="text-gradient-warm italic">
          Beleza &amp; Bem-Estar
        </span>
      </h2>
    </div>

    <!-- Services grid -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      {#each landingConfig.services.slice(0, 3) as service, i}
        <a
          href="/listagem"
          class="group relative overflow-hidden rounded-[2rem] h-[600px] cursor-pointer block card-hover-lift"
          style="animation-delay: {i * 100}ms;"
        >
          <img
            src={serviceImages[service.icon] || serviceImages['star']}
            alt={service.title}
            class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
          />
          <!-- Gradient overlay -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>
          <!-- Subtle pink tint on hover -->
          <div class="absolute inset-0 bg-gradient-to-t from-orange-900/20 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>

          <!-- Watermark number in Cormorant italic -->
          <div class="absolute top-6 right-6 text-white/10 font-cormorant italic font-bold text-[120px] leading-none select-none pointer-events-none group-hover:text-white/15 transition-colors duration-300">
            {cardNumbers[i]}
          </div>

          <div class="absolute bottom-0 left-0 right-0 p-8 text-white z-10">
            <!-- Category badge -->
            <span class="inline-block px-3 py-1 rounded-full bg-white/15 backdrop-blur-sm border border-white/20 text-xs font-bold font-montserrat text-orange-300 uppercase tracking-wider mb-3">
              {serviceCategoryLabels[service.icon] || 'Servico'}
            </span>
            <h3 class="text-3xl font-bold font-cormorant mt-1 mb-3">{service.title}</h3>
            <p class="text-sm font-montserrat text-gray-300 mb-5 leading-relaxed">{service.description}</p>
            <div class="w-11 h-11 rounded-full bg-white/10 backdrop-blur-sm border border-white/20 flex items-center justify-center hover:bg-gradient-to-r hover:from-orange-500 hover:to-pink-500 hover:border-transparent transition-all duration-300 group-hover:scale-110">
              <ArrowRight class="w-5 h-5" />
            </div>
          </div>
        </a>
      {/each}
    </div>
  </div>
</section>
