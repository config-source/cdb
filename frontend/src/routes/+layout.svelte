<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.scss';
	import { config } from '@fortawesome/fontawesome-svg-core';
	import '@fortawesome/fontawesome-svg-core/styles.css';
	config.autoAddCss = false;

	import { page } from '$app/stores';
	import { goto } from '$app/navigation';

	import { user } from '$lib/stores/user';
	import { getCurrentUser } from '$lib/client/auth';
	/**
	 * @typedef {Object} Props
	 * @property {import('svelte').Snippet} [children]
	 */

	/** @type {Props} */
	let { children } = $props();

	let isLoginPage = $derived($page.url.pathname.startsWith('/auth'));

	(async () => {
		if (!isLoginPage && $user.data?.Email === undefined) {
			const userInfo = await getCurrentUser();
			if (!userInfo.loggedIn) {
				return goto('/auth/login');
			}

			user.set({
				fetched: true,
				data: userInfo.user
			});
		}
	})();

	user.subscribe((data) => {
		if (!$page.url.pathname.startsWith('/auth') && data.fetched && data.data?.Email === undefined) {
			return goto('/auth/login');
		}
	});
</script>

{#if !isLoginPage}
	<Navbar />
{/if}

{@render children?.()}
