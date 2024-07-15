<script>
	/** @type any */
	export let environment = undefined;

	/** @type string */
	export let environmentName = environment?.Name;

	/** @type number */
	export let size = 5;

	/** @type any */
	let parent;

	/** @type (envName?: string) => Promise<void> */
	const getParent = async (envName) => {
		if (!envName || envName === '') return;

		if (!environment) {
			const selfRes = await fetch(`/api/v1/environments/by-name/${envName}`);
			if (!selfRes.ok) return;

			environment = await selfRes.json();
		}

		const parentID = environment.PromotesToID;
		if (!parentID) return;

		const parentRes = await fetch(`/api/v1/environments/by-id/${parentID}`);
		if (!parentRes.ok) return;

		parent = await parentRes.json();
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
