<script>
	import { isError } from '$lib/client';
	import { getByID } from '$lib/client/environments';

	/** @type App.Environment | undefined */
	export let environment = undefined;

	export let environmentId = 0;

	/** @type string */
	let environmentName = '';
	$: environmentName = environment?.Name ?? '';

	/** @type number */
	export let size = 5;

	/** @type App.Environment | undefined */
	let parent;

	/** @type (envId?: number) => Promise<void> */
	const getParent = async (envId) => {
		if (!envId || envId === 0) return;

		if (!environment) {
			const self = await getByID(envId);
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

	$: getParent(environmentId);
</script>

{#if parent}
	<svelte:self environment={parent} />
	<span class={`is-size-${size} p-1`}> / </span>
{/if}

{#if environmentName}
	<a href={`/environments/${environmentId}`} class={`is-size-${size} p-1`}>
		{environmentName}
	</a>
{/if}
