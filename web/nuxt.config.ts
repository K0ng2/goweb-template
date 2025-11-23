import tailwindcss from '@tailwindcss/vite'
import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'

// Read package.json to get version
const packageJson = JSON.parse(readFileSync(resolve(__dirname, 'package.json'), 'utf-8'))

// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
	compatibilityDate: '2025-07-15',
	devtools: { enabled: true },
	ssr: false,
	vite: {
		define: {
			__APP_VERSION__: JSON.stringify(packageJson.version),
		},
		plugins: [
			tailwindcss(),
		],
	},
	app: {
		head: {
			title: 'App', // change this
			meta: [
				{ name: 'description', content: 'App' }, // change this
				{ name: 'viewport', content: 'width=device-width, initial-scale=1' },
				{ name: 'theme-color', content: '#ffffff' },
				{ charset: 'utf-8' },
			],
			link: [
				{ rel: 'icon', type: 'image/icon', href: '/favicon.ico' },
			]
		}
	},
	pwa: {
		manifest: {
			name: 'App', // change this
			short_name: 'App', // change this
			start_url: '/',
			display: 'standalone',
			background_color: '#ffffff',
			orientation: "portrait",
			description: "App", // change this
			scope: '/',
			icons: [
				{
					src: "favicon.ico",
					sizes: "256x256",
					type: 'image/webp'
				}
			]
		},
		workbox: {
			globPatterns: ['**/*.{js,css,html,svg,png,ico}'],
			cleanupOutdatedCaches: true,
			clientsClaim: true,
			maximumFileSizeToCacheInBytes: 5 * 1024 * 1024, // 5MB
		},
		registerType: 'autoUpdate',
	},
	modules: [
		'@vite-pwa/nuxt',
	],
})
