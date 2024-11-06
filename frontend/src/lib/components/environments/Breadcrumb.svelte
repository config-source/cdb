<script>
	import Breadcrumb from './Breadcrumb.svelte';

	import { isError } from '$lib/client';
	import { getByID } from '$lib/client/environments';

	/**
	 * @typedef {Object} Props
	 * @property {App.Environment} environment
	 * @property {number} size
	 */

	/** @type {Props} */
	let { environment, size = 5 } = $props();

	/** @type Promise<App.Environment | undefined> */
	let promotesTo = $derived.by(async () => {
		const parentID = environment.PromotesToID;
		if (parentID === undefined) return undefined;

		const res = await getByID(parentID);
		if (isError(res)) return undefined;
		return res;
	});
</script>

{#await promotesTo then parent}
	{#if parent}
		<Breadcrumb environment={parent} {size} />
		<span class={`is-size-${size} p-1`}> / </span>
	{/if}
{/await}

<a href={`/environments/${environment.ID}`} class={`is-size-${size} p-1`}>
	{environment.Name}
</a>
