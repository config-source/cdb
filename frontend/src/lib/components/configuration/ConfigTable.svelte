<script>
	import EditableConfigRow from './ConfigTable/EditableConfigRow.svelte';
	import NewValues from './ConfigTable/NewValues.svelte';
	import { fetchConfig } from '$lib/client/config-values';

	// Stores all the fetched configuration values for this environment.
	/** @type any[] */
	export let configuration = [];
	/** @type number */
	export let environmentId;

	/** @type (envName: number) => Promise<void> */
	const updateConfiguration = async (envId) => {
		configuration = await fetchConfig(envId);
	};

	$: updateConfiguration(environmentId);
</script>

<table class="table is-fullwidth is-hoverable">
	<thead>
		<th>Key</th>
		<th>Value</th>
		<th>Inherited From</th>
		<th></th>
	</thead>
	<tbody>
		{#each configuration as configValue}
			<EditableConfigRow
				{configValue}
				{environmentId}
				on:updated={() => updateConfiguration(environmentId)}
			/>
		{/each}

		<NewValues
			{environmentId}
			existingKeys={configuration.map((cv) => cv.Name)}
			on:updated={() => updateConfiguration(environmentId)}
		/>
	</tbody>
</table>
