<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.scss';
	import { config } from '@fortawesome/fontawesome-svg-core';
	import '@fortawesome/fontawesome-svg-core/styles.css';
	config.autoAddCss = false;

	import { page } from '$app/stores';
	import { user } from '$lib/stores/user';
	import { goto } from '$app/navigation';

	let isLoginPage = false;
	$: isLoginPage = $page.url.pathname.startsWith('/auth');

	(async () => {
		if (!isLoginPage && $user.data.Email === undefined) {
			const res = await fetch('/api/v1/users/me', { credentials: 'include' });
			if (res.status === 401) {
				user.set({
					fetched: true,
					data: {}
				});

				return goto('/auth/login');
			} else if (!res.ok) {
				return;
			}

			const data = await res.json();
			user.set({
				fetched: true,
				data: data
			});
		}
	})();

	user.subscribe(({ data }) => {
		if (!$page.url.pathname.startsWith('/auth') && data.fetched && data.Email === undefined) {
			return goto('/auth/login');
		}
	});
</script>

{#if !isLoginPage}
	<Navbar />
{/if}

<slot />
