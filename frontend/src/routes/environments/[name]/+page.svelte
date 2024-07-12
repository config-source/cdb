<script>
	import ConfigTable from '$lib/components/configuration/ConfigTable.svelte';
	import Heading from '$lib/components/utility/Heading.svelte';
	import Breadcrumbs from '$lib/components/environments/Breadcrumbs.svelte';

	export let data;
</script>

<div class="container mt-6">
	<Breadcrumbs environmentName={data.name} />
	<Heading size={3}>Configuration for {data.name}</Heading>
	<ConfigTable environmentName={data.name} />
</div>
