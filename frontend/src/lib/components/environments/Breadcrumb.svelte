<script>
	import { isError } from '$lib/client';
	import { getByID } from '$lib/client/environments';

	/** @type App.Environment | undefined */
	export let environment = undefined;

	/** @type number */
	export let size = 5;

	/** @type App.Environment | undefined */
	let parent;

	/** @type (parentId?: number) => Promise<void> */
	const getParent = async (parentId) => {
		if (!parentId || parentId === 0) return;

		const parentEnv = await getByID(parentId);
		if (!isError(parentEnv)) {
			parent = parentEnv;
		}
	};

	$: getParent(environment?.PromotesToID);
</script>

{#if parent}
	<svelte:self environment={parent} />
	<span class={`is-size-${size} p-1`}> / </span>
{/if}

{#if environment}
	<a href={`/environments/${environment.ID}`} class={`is-size-${size} p-1`}>
		{environment.Name}
	</a>
{/if}
