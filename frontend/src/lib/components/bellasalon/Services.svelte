<script lang="ts">
  import ArrowRight from 'lucide-svelte/icons/arrow-right';
</script>

<section id="services" class="py-32 bg-gray-50 dark:bg-black">
  <div class="container mx-auto px-4">
    <!-- Section header -->
    <div class="mb-16">
      <span class="text-xs font-bold text-orange-500 uppercase tracking-[0.3em]">
        Nossos Serviços
      </span>
      <h2 class="text-5xl md:text-6xl lg:text-7xl font-bold mt-6 text-gray-900 dark:text-white leading-tight tracking-tight font-sans-bs">
        Experiências de <br />
        <span class="text-orange-500 underline decoration-orange-500/20 underline-offset-8">
          Beleza & Bem-Estar
        </span>
      </h2>
    </div>

    <!-- Services grid com imagens -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <!-- Manicure & Pedicure -->
      <a href="/listagem" class="group relative overflow-hidden rounded-[2rem] h-[600px] cursor-pointer block">
        <img
          src="https://images.unsplash.com/photo-1604654894610-df63bc536371?auto=format&fit=crop&q=80&w=800"
          alt="Manicure & Pedicure"
          class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>

        <div class="absolute bottom-0 left-0 right-0 p-8 text-white z-10">
          <span class="text-xs font-bold text-orange-500 uppercase tracking-wider">Estética</span>
          <h3 class="text-3xl font-bold mt-2 mb-4">Manicure & Pedicure</h3>
          <p class="text-sm text-gray-300 mb-4">A partir de R$ 45,00</p>
          <div class="w-12 h-12 rounded-full bg-white/10 backdrop-blur-sm border border-white/20 flex items-center justify-center hover:bg-orange-500 transition-all duration-300 group-hover:scale-110">
            <ArrowRight class="w-5 h-5" />
          </div>
        </div>
      </a>

      <!-- Corte & Coloração -->
      <a href="/listagem" class="group relative overflow-hidden rounded-[2rem] h-[600px] cursor-pointer block">
        <img
          src="https://images.unsplash.com/photo-1560066984-138dadb4c035?auto=format&fit=crop&q=80&w=800"
          alt="Corte & Coloração"
          class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>

        <div class="absolute bottom-0 left-0 right-0 p-8 text-white z-10">
          <span class="text-xs font-bold text-orange-500 uppercase tracking-wider">Capilar</span>
          <h3 class="text-3xl font-bold mt-2 mb-2">Corte & Coloração</h3>
          <p class="text-sm text-gray-300 mb-4">Mechas criativas, morena iluminada e cortes de precisão com L'Oréal e Wella.</p>
          <p class="text-sm text-gray-300 mb-4">A partir de R$ 120,00</p>
          <div class="w-12 h-12 rounded-full bg-white/10 backdrop-blur-sm border border-white/20 flex items-center justify-center hover:bg-orange-500 transition-all duration-300 group-hover:scale-110">
            <ArrowRight class="w-5 h-5" />
          </div>
        </div>
      </a>

      <!-- Spa & Massagens -->
      <a href="/listagem" class="group relative overflow-hidden rounded-[2rem] h-[600px] cursor-pointer block">
        <img
          src="https://images.unsplash.com/photo-1544161515-4ab6ce6db874?auto=format&fit=crop&q=80&w=800"
          alt="Spa & Massagens"
          class="absolute inset-0 w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/90 via-black/40 to-transparent"></div>

        <div class="absolute bottom-0 left-0 right-0 p-8 text-white z-10">
          <span class="text-xs font-bold text-orange-500 uppercase tracking-wider">Bem-Estar</span>
          <h3 class="text-3xl font-bold mt-2 mb-4">Spa & Massagens</h3>
          <p class="text-sm text-gray-300 mb-4">A partir de R$ 90,00</p>
          <div class="w-12 h-12 rounded-full bg-white/10 backdrop-blur-sm border border-white/20 flex items-center justify-center hover:bg-orange-500 transition-all duration-300 group-hover:scale-110">
            <ArrowRight class="w-5 h-5" />
          </div>
        </div>
      </a>
    </div>
  </div>
</section>
