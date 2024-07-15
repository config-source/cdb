<script>
	import EditableConfigRow from './ConfigTable/EditableConfigRow.svelte';
	import NewValues from './ConfigTable/NewValues.svelte';

	export let environmentName;

	// Stores all the fetched configuration values for this environment.
	/** @type any[] */
	let configuration = [];

	/** @type (envName: string) => Promise<void> */
	const fetchConfig = async (envName) => {
		if (envName === '') return;

		const res = await fetch(`/api/v1/config-values/${envName}`);
		if (!res.ok) return;

		/** @type any[] */
		const data = await res.json();
		data.sort((a, b) => {
			if (a.Inherited && !b.Inherited) {
				return 1;
			}

			if (!a.Inherited && b.Inherited) {
				return -1;
			}

			const nameA = a.Name.toUpperCase();
			const nameB = b.Name.toUpperCase();
			if (nameA < nameB) {
				return -1;
			}

			if (nameA > nameB) {
				return 1;
			}

			return 1;
		});
		configuration = data;
	};
	$: fetchConfig(environmentName);
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
				on:updated={() => fetchConfig(environmentName)}
			/>
		{/each}

		<NewValues
			{environmentName}
			existingKeys={configuration.map((cv) => cv.Name)}
			on:updated={() => fetchConfig(environmentName)}
		/>
	</tbody>
</table>
