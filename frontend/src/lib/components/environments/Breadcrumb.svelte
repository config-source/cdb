<script>
	/** @type string | undefined */
	export let environmentName = undefined;
	/** @type number | undefined */
	export let environmentID = undefined;

	/** @type number */
	export let size = 5;

	/** @type number */
	let parentID;

	/** @type (envID?: number) => Promise<void> */
	const getName = async (envID) => {
		if (!envID) return;

		const res = await fetch(`/api/v1/environments/by-id/${envID}`);
		if (!res.ok) return;

		const env = await res.json();
		environmentName = env.Name;
	};

	/** @type (envName?: string) => Promise<void> */
	const getParentByName = async (envName) => {
		if (!envName || envName === '') return;

		const res = await fetch(`/api/v1/environments/by-name/${envName}`);
		if (!res.ok) return;

		const env = await res.json();
		parentID = env.PromotesToID;
	};

	/** @type (envID?: number) => Promise<void> */
	const getParentByID = async (envName) => {
		if (!envName || envName === 0) return;

		const res = await fetch(`/api/v1/environments/by-name/${envName}`);
		if (!res.ok) return;

		const env = await res.json();
		parentID = env.PromotesToID;
	};

	$: getParentByName(environmentName);
	$: getParentByID(environmentID);
	$: getName(environmentID);
</script>

{#if parentID}
	<svelte:self environmentID={parentID} />
	<span class={`is-size-${size} p-1`}> / </span>
{/if}

{#if environmentName}
	<a href={`/environments/${environmentName}`} class={`is-size-${size} p-1`}>
		{environmentName}
	</a>
{/if}
