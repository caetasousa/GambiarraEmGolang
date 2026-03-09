<script lang="ts">
  import { onMount } from 'svelte';
  import Navbar from '$lib/components/Navbar.svelte';
  import Footer from '$lib/components/Footer.svelte';
  import Hero from '$lib/components/bellasalon/Hero.svelte';
  import Services from '$lib/components/bellasalon/Services.svelte';
  import Specialists from '$lib/components/bellasalon/Specialists.svelte';
  import Testimonials from '$lib/components/bellasalon/Testimonials.svelte';
  import CTA from '$lib/components/bellasalon/CTA.svelte';
  import { landingConfig } from '$lib/config/landingConfig';

  // Smooth scroll behavior for anchor links
  onMount(() => {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
      anchor.addEventListener('click', (e) => {
        e.preventDefault();
        const href = anchor.getAttribute('href');
        const target = href ? document.querySelector(href) : null;
        if (target) {
          target.scrollIntoView({
            behavior: 'smooth',
            block: 'start'
          });
        }
      });
    });
  });
</script>

<svelte:head>
  <title>BellaSalon - Realce sua Beleza Natural</title>
  <meta name="description" content="Descubra uma experiência única de cuidado pessoal. Do design de sobrancelhas a tratamentos capilares exclusivos com excelência e sofisticação." />
</svelte:head>

<div class="min-h-screen bg-[hsl(var(--bs-background))] text-[hsl(var(--bs-foreground))] font-sans-bs">
  <Navbar />

  <main class="overflow-x-hidden">
    <Hero data={landingConfig.hero} />
    <Services />
    <Specialists specialists={landingConfig.specialists} />
    <Testimonials testimonials={landingConfig.testimonials} />
    <CTA
      title={landingConfig.cta.title}
      description={landingConfig.cta.description}
      buttonText={landingConfig.cta.buttonText}
    />
  </main>

  <Footer data={landingConfig.footer} />
</div>
