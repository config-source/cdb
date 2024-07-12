<script>
	import Breadcrumb from './Breadcrumb.svelte';

	export let environmentName;
	export let size = undefined;
</script>

<div class="is-flex-direction-row is-flex">
	<Breadcrumb {size} {environmentName} />
</div>
