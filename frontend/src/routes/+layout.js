import { goto } from '$app/navigation';
import { user } from '$lib/stores/user';

export const ssr = false;
export const prerender = false;

export async function load({ url, fetch }) {
	console.log(url.pathname, url.pathname.startsWith('/auth'));
	if (!url.pathname.startsWith('/auth')) {
		const res = await fetch('/api/v1/users/me');
		if (!res.ok) {
			user.set({
				fetched: true,
				data: {}
			});

			return goto('/auth/login');
		}

		const data = await res.json();
		user.set({
			fetched: true,
			data: data
		});
	}
}
