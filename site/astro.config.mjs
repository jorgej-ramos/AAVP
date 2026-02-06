import { defineConfig } from 'astro/config';
import rehypeMermaid from 'rehype-mermaid';

export default defineConfig({
  site: 'https://jorgej-ramos.github.io',
  base: '/AAVP/',
  markdown: {
    syntaxHighlight: {
      type: 'shiki',
      excludeLangs: ['mermaid'],
    },
    rehypePlugins: [
      [rehypeMermaid, { strategy: 'img-svg', dark: false }],
    ],
  },
});
