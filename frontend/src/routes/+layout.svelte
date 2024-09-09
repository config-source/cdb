<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.scss';
	import { config } from '@fortawesome/fontawesome-svg-core';
	import '@fortawesome/fontawesome-svg-core/styles.css';
	config.autoAddCss = false;

	import { page } from '$app/stores';
	import { user } from '$lib/stores/user';
	import { goto } from '$app/navigation';

	user.subscribe(({ data }) => {
		if (!$page.url.pathname.startsWith('/auth') && data.Email === undefined) {
			return goto('/auth/login');
		}
	});
</script>

{#if $user.fetched && $user.data.Email}
	<Navbar />
{/if}

<slot />
