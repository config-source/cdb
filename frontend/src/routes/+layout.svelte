<script>
	import Navbar from '$lib/components/Navbar.svelte';
	import '../app.scss';
	import { config } from '@fortawesome/fontawesome-svg-core';
	import '@fortawesome/fontawesome-svg-core/styles.css';
	config.autoAddCss = false;
</script>

<Navbar />

<slot />
