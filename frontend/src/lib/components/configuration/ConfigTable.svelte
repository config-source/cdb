<script>
	import EditableConfigRow from './ConfigTable/EditableConfigRow.svelte';
	import NewValues from './ConfigTable/NewValues.svelte';
	import { fetchConfig } from '$lib/client/config-values';

	// Stores all the fetched configuration values for this environment.
	/** @type any[] */
	export let configuration = [];
	/** @type string */
	export let environmentName;

	/** @type (envName: string) => Promise<void> */
	const updateConfiguration = async (envName) => {
		configuration = await fetchConfig(envName);
	};

	$: updateConfiguration(environmentName);
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
				{environmentName}
				on:updated={() => updateConfiguration(environmentName)}
			/>
		{/each}

		<NewValues
			{environmentName}
			existingKeys={configuration.map((cv) => cv.Name)}
			on:updated={() => updateConfiguration(environmentName)}
		/>
	</tbody>
</table>
