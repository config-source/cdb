<script>
	export let environmentName;

	let configuration = [];

	const fetchConfig = async (name) => {
		if (name === '') return;

		const res = await fetch(`/api/v1/config-values/${name}`);
		if (!res.ok) return;

		const data = await res.json();
		data.sort((a, b) => {
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

	const getValue = (configValue) => {
		const valueKeys = ['StrValue', 'IntValue', 'BoolValue', 'FloatValue'];
		for (const key of valueKeys) {
			const value = configValue[key];
			if (value) {
				return value;
			}
		}

		return 'UNRECOGNISED';
	};

	$: fetchConfig(environmentName);
</script>

<table class="table is-fullwidth is-hoverable">
	<thead>
		<th>Key</th>
		<th>Value</th>
		<th>Inherited</th>
	</thead>
	{#each configuration as configValue}
		<tr>
			<td>{configValue.Name}</td>
			<td>{getValue(configValue)}</td>
			<td>{configValue.Inherited}</td>
		</tr>
	{/each}
</table>
