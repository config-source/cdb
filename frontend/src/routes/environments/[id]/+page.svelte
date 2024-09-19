<script>
	import ConfigTable from '$lib/components/configuration/ConfigTable.svelte';
	import Heading from '$lib/components/utility/Heading.svelte';
	import Breadcrumbs from '$lib/components/environments/Breadcrumbs.svelte';
	import { isError } from '$lib/client';
	import { getByID } from '$lib/client/environments';

	export let data;

	/** @type number */
	let environmentId = 0;
	$: environmentId = parseInt(data.id, 10);

	/** @type App.Environment */
	let environment;

	$: {
		getByID(environmentId).then((e) => {
			if (isError(e)) {
				// TODO: handle error
			} else {
				environment = e;
			}
		});
	}
</script>

<div class="container mt-6">
	<Breadcrumbs {environment} />
	<Heading size={3}>Configuration for {environment?.Name}</Heading>
	<ConfigTable {environmentId} />
</div>
