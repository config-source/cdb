import { writable } from 'svelte/store';

export const user = writable({
	fetched: false,
	data: {}
});
