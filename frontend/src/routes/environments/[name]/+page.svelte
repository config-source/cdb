<script>
	import ConfigTable from '$lib/components/configuration/ConfigTable.svelte';
	import Heading from '$lib/components/utility/Heading.svelte';

	export let data;
</script>

<div class="container mt-6">
	<Heading size={3}>Configuration for {data.name}</Heading>
	<ConfigTable environmentName={data.name} />
</div>
