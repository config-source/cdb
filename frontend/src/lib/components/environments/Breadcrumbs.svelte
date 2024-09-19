<script>
	import Breadcrumb from './Breadcrumb.svelte';

	export let environment;
	export let size = undefined;
</script>

<div class="is-flex-direction-row is-flex">
	<Breadcrumb {size} {environment} />
</div>
