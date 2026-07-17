import { defineConfig } from 'astro/config';

export default defineConfig({
  site: 'https://ash.sh',
  output: 'static',
  build: { format: 'file' }
});