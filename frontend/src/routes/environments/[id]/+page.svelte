<script>
	import ConfigTable from '$lib/components/configuration/ConfigTable.svelte';
	import Heading from '$lib/components/utility/Heading.svelte';
	import Breadcrumbs from '$lib/components/environments/Breadcrumbs.svelte';
	import { isError } from '$lib/client';
	import { getByID } from '$lib/client/environments';

	/**
	 * @typedef Props
	 * @property {{id: string}} data
	 */

	/** @type {Props} */
	let { data } = $props();

	/** @type number */
	let environmentId = $derived(parseInt(data.id, 10));

	let env = $derived.by(async () => {
		const res = await getByID(environmentId);
		if (isError(res)) throw new Error(res.Message);
		return res;
	});
</script>

<div class="container mt-6">
	{#await env then environment}
		<Breadcrumbs {environment} />
		<Heading size={3}>Configuration for {environment.Name}</Heading>
	{/await}

	<ConfigTable {environmentId} />
</div>
