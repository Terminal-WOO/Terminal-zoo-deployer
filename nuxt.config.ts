import { defineNuxtConfig } from 'nuxt/config'
import tailwindcss from "@tailwindcss/vite";
import { resolve } from 'pathe';

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  alias: {
    '@': resolve(__dirname) + '/app',
  },
  runtimeConfig: {
    externalApiBase: process.env.NUXT_EXTERNAL_API_BASE || 'https://backend-deployer.clappform.com',//'http://localhost:8080', //'https://backend-deployer.clappform.com/',
    externalApiAuth: process.env.NUXT_EXTERNAL_API_AUTH || '', // e.g. 'Bearer abc'
  },
  compatibilityDate: '2025-07-15',
  devtools: { enabled: false },
  modules: ['@nuxt/ui', '@nuxt/test-utils', '@primevue/nuxt-module', '@pinia/nuxt'],
  primevue: {
    importTheme: { from: '@/themes/theme.js' },
  },
  app: {
    head: {
      link: [
        { rel: 'preconnect', href: 'https://fonts.googleapis.com' },
        { rel: 'preconnect', href: 'https://fonts.gstatic.com', crossorigin: '' },
        {
          rel: 'stylesheet',
          href: 'https://fonts.googleapis.com/css2?family=IBM+Plex+Sans:wght@400;500;600;700&display=swap'
        },
      ]
    }
  },
    nitro: {
      storage: {
      data: { driver: 'fs', base: '.data' } // writes go to ./.data/
    },
    routeRules: {
      "/": {
        headers: {
          "cors": "*",
          "Referrer-Policy": "no-referrer",
          "X-XSS-Protection": "1; mode=block",
          "X-Content-Type-Options": "nosniff",
          "X-Frame-Options": "sameorigin",
          "Content-Security-Policy": `default-src 'self';
            script-src  'self' 'unsafe-inline' 'unsafe-eval';
            style-src   'self' 'unsafe-inline';
            img-src     'self' data: https://clappformimages.blob.core.windows.net;
            font-src    'self';
            connect-src 'self' https://*.clappform.com http://localhost:8080;
            frame-ancestors 'none';
            object-src  'none';
            base-uri    'self';
            form-action 'self';
            worker-src  'self' blob:;
            upgrade-insecure-requests;
            report-uri  https://clappform.com/nlcontact;`
            .replace(/\s*\n\s*/g, " ")
            .trim(),
          "Permissions-Policy":
            'geolocation=(self "https://*.mapbox.com"), camera=(), microphone=()',
        },
      },
    },
  },
  css: [
    'primeflex/primeflex.css',
    'primeicons/primeicons.css',
    '@/assets/css/typography.css'
  ],
  vite: {
    plugins: [
      tailwindcss(),
    ],
    server: {
      watch: {
        ignored: ['**/assets/userApps.json'],
        usePolling: true,
        interval: 200,     // probeer 100–500ms
      }
    }
  },
})
