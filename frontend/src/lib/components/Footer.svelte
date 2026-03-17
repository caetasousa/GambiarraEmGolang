<script lang="ts">
  import type { FooterData } from '$lib/config/landingConfig';

  let { data }: { data: FooterData } = $props();
</script>

<footer class="bg-gradient-to-b from-slate-800 to-slate-900 border-t border-slate-700/50 dark:border-white/5 pt-28 pb-12">
  <div class="container mx-auto px-4">
    <div class="grid grid-cols-1 md:grid-cols-[2fr_1fr_1.5fr_1.5fr] gap-16 mb-24">
      <!-- Brand -->
      <div class="flex flex-col gap-8">
        <a href="#home" class="flex items-center gap-3 cursor-pointer group">
          <div class="w-10 h-10 bg-gradient-to-br from-orange-500 to-pink-500 rounded-xl flex items-center justify-center text-white font-bold text-xl shadow-[0_10px_15px_-3px_rgba(249,115,22,0.35)] transition-transform duration-300 group-hover:scale-105">
            B
          </div>
          <span class="text-2xl font-bold text-white tracking-wide font-cormorant">
            BellaVita
          </span>
        </a>
        <p class="text-gray-400 text-lg leading-relaxed">
          {data.description}
        </p>
      </div>

      <!-- Sections -->
      {#each data.sections as section}
        <div>
          <h4 class="font-bold text-white mb-8 uppercase tracking-wider text-sm">
            {section.title}
          </h4>
          <ul class="flex flex-col gap-5">
            {#each section.links as link}
              <li>
                <a
                  href={link.href}
                  class="text-gray-400 text-sm font-medium hover:text-orange-500 transition-colors duration-300"
                >
                  {link.label}
                </a>
              </li>
            {/each}
          </ul>
        </div>
      {/each}

      <!-- Contact -->
      <div>
        <h4 class="font-bold text-white mb-8 uppercase tracking-wider text-sm">
          Visite-nos
        </h4>
        <ul class="flex flex-col gap-5">
          <li class="flex items-center gap-3">
            <span class="w-1.5 h-1.5 rounded-full bg-orange-500"></span>
            <span class="text-gray-400 text-sm font-medium">
              {data.contact.address}
            </span>
          </li>
          <li class="flex items-center gap-3">
            <span class="w-1.5 h-1.5 rounded-full bg-orange-500"></span>
            <span class="text-gray-400 text-sm font-medium">
              {data.contact.city}
            </span>
          </li>
          <li class="flex items-center gap-3">
            <span class="w-1.5 h-1.5 rounded-full bg-orange-500"></span>
            <span class="text-gray-400 text-sm font-medium">
              {data.contact.phone}
            </span>
          </li>
          <li class="flex items-center gap-3">
            <span class="w-1.5 h-1.5 rounded-full bg-orange-500"></span>
            <span class="text-gray-400 text-sm font-medium">
              {data.contact.email}
            </span>
          </li>
        </ul>
      </div>
    </div>

    <!-- Bottom -->
    <div class="border-t border-white/5 pt-12 flex flex-col md:flex-row justify-between items-center gap-6 text-xs text-gray-500 font-bold uppercase tracking-wider">
      <p>&copy; {data.copyright}</p>
      <div class="flex gap-10">
        {#each data.bottomLinks as link}
          <a
            href={link.href}
            class="text-gray-500 underline decoration-orange-500/30 underline-offset-1 hover:text-gray-300 transition-colors duration-300"
          >
            {link.label}
          </a>
        {/each}
      </div>
    </div>
  </div>
</footer>
