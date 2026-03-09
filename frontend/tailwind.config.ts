import type { Config } from 'tailwindcss';
import forms from '@tailwindcss/forms';
import typography from '@tailwindcss/typography';

export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        "brand": "#f48c25",
        "brand-hover": "#e07b14",
        "primary": "#f48c25",
        "secondary": "#232630",
        "background-light": "#f8f7f5",
        "background-dark": "#221910",
        "surface-dark": "#1F2937",
        "accent-yellow": "#FDBA74",
        "brand-orange": "#f48c25",
        "brand-orange-hover": "#e07b14",
        "surface-light": "#FFFFFF",
        "text-light": "#1F2937",
        "text-dark": "#F3F4F6",
        "border-light": "#E5E7EB",
        "border-dark": "#374151"
      },
      fontFamily: {
        display: "Spline Sans",
        sans: ["Poppins", "sans-serif"],
        body: ["Inter", "sans-serif"]
      },
      borderRadius: {
        DEFAULT: "0.25rem",
        lg: "0.5rem",
        xl: "0.75rem",
        full: "9999px"
      }
    }
  },
  plugins: [forms, typography]
} satisfies Config;
