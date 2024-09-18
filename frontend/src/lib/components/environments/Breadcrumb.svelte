<script>
	import { isError } from '$lib/client';
	import { getByID, getByName } from '$lib/client/environments';

	/** @type App.Environment | undefined */
	export let environment = undefined;

	/** @type string */
	export let environmentName = environment?.Name ?? '';

	/** @type number */
	export let size = 5;

	/** @type App.Environment | undefined */
	let parent;

	/** @type (envName?: string) => Promise<void> */
	const getParent = async (envName) => {
		if (!envName || envName === '') return;

		if (!environment) {
			const self = await getByName(envName);
			if (!isError(self)) {
				environment = self;
			}
		}

		const parentID = environment?.PromotesToID;
		if (!parentID) return;

		const parentEnv = await getByID(parentID);
		if (!isError(parentEnv)) {
			parent = parentEnv;
		}
	};

	$: getParent(environmentName);
</script>

{#if parent}
	<svelte:self environment={parent} />
	<span class={`is-size-${size} p-1`}> / </span>
{/if}

{#if environmentName}
	<a href={`/environments/${environmentName}`} class={`is-size-${size} p-1`}>
		{environmentName}
	</a>
{/if}
