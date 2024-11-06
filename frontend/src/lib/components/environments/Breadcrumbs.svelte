<script>
	import Breadcrumb from './Breadcrumb.svelte';

	let { environment, size = undefined } = $props();
</script>

<div class="is-flex-direction-row is-flex">
	<Breadcrumb {size} {environment} />
</div>
