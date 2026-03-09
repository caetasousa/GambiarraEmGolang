export interface ServiceData {
  icon: string;
  title: string;
  description: string;
  features: string[];
}

export interface SpecialistData {
  name: string;
  role: string;
  specialty: string;
  imageUrl: string;
}

export interface TestimonialData {
  text: string;
  author: string;
  role: string;
  rating: number;
}

export interface HeroData {
  badge: string;
  title: string;
  titleAccent: string;
  description: string;
  primaryButtonText: string;
  secondaryButtonText: string;
  stats: {
    value: string;
    label: string;
  }[];
  imageUrl: string;
  floatingCard: {
    title: string;
    subtitle: string;
    text: string;
    rating: number;
  };
}

export interface FooterData {
  description: string;
  sections: {
    title: string;
    links: {
      label: string;
      href: string;
    }[];
  }[];
  contact: {
    address: string;
    city: string;
    phone: string;
    email: string;
  };
  bottomLinks: {
    label: string;
    href: string;
  }[];
  copyright: string;
}

export const landingConfig = {
  hero: {
    badge: "Nova Experiência em Beleza",
    title: "Realce sua",
    titleAccent: "Beleza Natural",
    description: "Descubra uma experiência única de cuidado pessoal. Do design de sobrancelhas a tratamentos capilares exclusivos com excelência e sofisticação.",
    primaryButtonText: "Agendar Visita",
    secondaryButtonText: "Conhecer Serviços",
    stats: [
      { value: "10k+", label: "Clientes Felizes" },
      { value: "98%", label: "Satisfação" },
      { value: "10+", label: "Especialistas" }
    ],
    imageUrl: "https://images.unsplash.com/photo-1560066984-138dadb4c035?auto=format&fit=crop&q=80&w=1000",
    floatingCard: {
      title: "Tratamento VIP",
      subtitle: "Exclusivo para membros",
      text: "O melhor atendimento que já recebi. Me senti uma rainha!",
      rating: 5
    }
  } as HeroData,

  services: [
    {
      icon: "scissors",
      title: "Corte & Styling",
      description: "Cortes personalizados que valorizam suas características únicas com técnicas modernas e tradicionais",
      features: [
        "Análise de Visagismo Detalhada",
        "Produtos Premium de Alta Qualidade",
        "Finalização Profissional Exclusiva"
      ]
    },
    {
      icon: "sparkles",
      title: "Coloração Expert",
      description: "Transforme seu visual com colorações vibrantes e duradouras aplicadas por especialistas certificados",
      features: [
        "Colorimetria Profissional Avançada",
        "Técnicas de Balayage e Ombré Hair",
        "Tratamento de Proteção e Hidratação"
      ]
    },
    {
      icon: "calendar",
      title: "Nail Design",
      description: "Unhas impecáveis com alongamentos em gel, nail art exclusiva e cuidados especiais de manicure",
      features: [
        "Alongamento em Gel de Última Geração",
        "Nail Art Personalizada e Criativa",
        "Biossegurança e Esterilização Total"
      ]
    },
    {
      icon: "star",
      title: "Spa & Estética",
      description: "Relaxe e rejuvenesça com tratamentos faciais e corporais que renovam corpo e mente",
      features: [
        "Tratamentos Faciais Personalizados",
        "Massagens Relaxantes e Terapêuticas",
        "Aromaterapia e Cromoterapia"
      ]
    }
  ] as ServiceData[],

  specialists: [
    {
      name: "Helena Souza",
      role: "Master Stylist",
      specialty: "Expert em Coloração & Corte",
      imageUrl: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?auto=format&fit=crop&q=80&w=200&h=200"
    },
    {
      name: "Ricardo Lima",
      role: "Design de Barba & Corte",
      specialty: "Expert em Visagismo Masculino",
      imageUrl: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?auto=format&fit=crop&q=80&w=200&h=200"
    },
    {
      name: "Bia Santos",
      role: "Nail Artist",
      specialty: "Expert em Alongamento em Gel",
      imageUrl: "https://images.unsplash.com/photo-1580489944761-15a19d654956?auto=format&fit=crop&q=80&w=200&h=200"
    }
  ] as SpecialistData[],

  testimonials: [
    {
      text: "O atendimento da Helena é impecável. Ela entendeu exatamente o que eu queria para o meu cabelo!",
      author: "Ana Clara",
      role: "Cliente Verificada",
      rating: 5
    },
    {
      text: "Espaço moderno e profissionais extremamente qualificados. O Ricardo é o melhor barbeiro da região.",
      author: "Marcos Oliveira",
      role: "Cliente Verificado",
      rating: 5
    },
    {
      text: "As unhas feitas pela Bia duram semanas. O cuidado com a biossegurança me impressionou.",
      author: "Juliana Costa",
      role: "Cliente Verificada",
      rating: 5
    }
  ] as TestimonialData[],

  cta: {
    title: "Sua Melhor Versão Começa Aqui",
    description: "Transformação real feita por quem entende de alma e estilo.",
    buttonText: "Agendar Agora"
  },

  footer: {
    description: "Elevando sua beleza natural com profissionais de elite e produtos de classe mundial. Referência em cuidado e sofisticação.",
    sections: [
      {
        title: "Explorar",
        links: [
          { label: "Início", href: "#home" },
          { label: "Serviços", href: "#services" },
          { label: "Especialistas", href: "#specialists" },
          { label: "Depoimentos", href: "#testimonials" }
        ]
      },
      {
        title: "Nossos Serviços",
        links: [
          { label: "Cabelo & Coloração Expert", href: "#services" },
          { label: "Unhas & Nail Design Artist", href: "#services" },
          { label: "Maquiagem Social & Editorial", href: "#services" },
          { label: "Spa & Estética Avançada", href: "#services" }
        ]
      }
    ],
    contact: {
      address: "Rua das Flores, 123",
      city: "São Paulo, SP",
      phone: "(11) 99999-9999",
      email: "contato@bellasalon.com"
    },
    bottomLinks: [
      { label: "Privacidade", href: "#" },
      { label: "Termos", href: "#" }
    ],
    copyright: "2024 BellaSalon. Crafted with Excellence."
  } as FooterData
};
