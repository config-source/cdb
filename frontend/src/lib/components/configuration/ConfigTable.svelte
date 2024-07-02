<script>
	import { ValueTypes, getValue } from '$lib/config-values';
	import ConfigValueInput from './ConfigValueInput.svelte';
	import ConfigKeySelector from './ConfigKeySelector.svelte';

	export let environmentName;

	let configuration = [];
	let newValues = [];

	const addNewConfigValue = () => {
		newValues = [
			...newValues,
			{
				Name: '',
				ValueType: ValueTypes.FLOAT
			}
		];
	};

	const removeNewConfigValue = (idx) => () => {
		newValues.splice(idx, 1);
		newValues = newValues;
	};

	const updateValueWithNewKey =
		(configValue) =>
		({ detail }) => {
			const newKey = detail.value;
			if (!newKey) return;

			configValue.Name = newKey.Name;
			configValue.ValueType = newKey.ValueType;
			configValue.ConfigKeyID = newKey.ID;

			for (const key of Object.keys(configValue).filter((k) => k.endsWith('Value'))) {
				configValue[key] = undefined;
			}

			newValues = [...newValues];
		};

	const fetchConfig = async (name) => {
		if (name === '') return;

		const res = await fetch(`/api/v1/config-values/${name}`);
		if (!res.ok) return;

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

	let excludedConfigKeys = [];
	$: {
		// Exclude keys which are already configured for this environment.
		const directlyConfiguredKeys = configuration.filter((cv) => !cv.Inherited).map((cv) => cv.Name);
		// Exclude keys which are already in the newValues list so we don't have
		// duplicates.
		const currentlyConfiguringKeys = newValues.filter((cv) => cv.Name);
		excludedConfigKeys = [...directlyConfiguredKeys, ...currentlyConfiguringKeys];
		console.log(excludedConfigKeys);
	}

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
			<tr>
				<td>{configValue.Name}</td>
				<td>{getValue(configValue)}</td>
				<td>
					{configValue.InheritedFrom}
				</td>
				<td></td>
			</tr>
		{/each}

		{#each newValues as newValue, i}
			<tr>
				<td>
					<!-- Exclude config keys that are already configured directly on this environment. -->
					<ConfigKeySelector
						excludedKeys={excludedConfigKeys}
						preSelectedName={newValue.Name}
						on:updated={updateValueWithNewKey(newValue)}
					/>
				</td>
				<td>
					<ConfigValueInput
						valueType={newValue.ValueType}
						value={getValue(newValue)}
						on:updated={(rawVal) => updateValue(newValue, rawVal)}
					/>
				</td>
				<td></td>
				<td>
					<button class="button is-error" on:click={removeNewConfigValue(i)}>
						<span class="icon">
							<i class="fab fa-trash"></i>
						</span>
					</button>
				</td>
			</tr>
		{/each}

		<tr>
			<td></td>
			<td></td>
			<td></td>
			<td>
				<button class="button is-success" on:click={addNewConfigValue}>
					<span class="icon">
						<i class="fab fa-plus"></i>
					</span>
				</button>
			</td>
		</tr>
	</tbody>
</table>
