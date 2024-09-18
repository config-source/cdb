import { writable } from 'svelte/store';

/** @type import('svelte/store').Writable<{ fetched: boolean, data?: App.User }> */
export const user = writable({
	fetched: false
});
