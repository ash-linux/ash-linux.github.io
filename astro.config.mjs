import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://ash-linux.github.io',
  base: '/',
  output: 'static',
  build: { format: 'file' }
});