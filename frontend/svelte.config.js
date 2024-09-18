import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import { mdsvex } from 'mdsvex';
import adapter from '@sveltejs/adapter-static';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapter({
			fallback: 'index.html'
		})
	},

	// 	paths: {
	// 		base: '/app'
	// 	},

	extensions: ['.svelte', '.md'],

	preprocess: [
		vitePreprocess({}),
		mdsvex({
			extensions: ['.md']
		})
	]
};

export default config;
