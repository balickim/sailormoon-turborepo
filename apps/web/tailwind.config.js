/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/renderer/index.html', './src/renderer/src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        background: 'var(--color-background)',
        text1: 'var(--ev-c-text-1)',
        text2: 'var(--ev-c-text-2)',
        gray2: 'var(--ev-c-gray-2)'
      }
    }
  },
  plugins: []
}
