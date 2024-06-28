<script>
	import { ValueType, getValue, updateValue } from '$lib/config-values';
	import ConfigValueInput from './ConfigValueInput.svelte';
	import ConfigKeySelector from './ConfigKeySelector.svelte';
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faCheck, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';

	export let environmentName;

	/** @type (envName: string) => Promise<void> */
	const fetchConfig = async (envName) => {
		if (envName === '') return;

		const res = await fetch(`/api/v1/config-values/${envName}`);
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

	let configuration = [];
	let newValues = [];

	const addNewConfigValue = () => {
		newValues = [
			...newValues,
			{
				Name: '',
				ValueType: ValueType.FLOAT
			}
		];
	};

	const removeNewConfigValue = (idx) => () => {
		newValues.splice(idx, 1);
		newValues = newValues;
	};

	const saveNewConfigValue = (idx) => async () => {
		console.log(newValues);

		// Remove and retrieve the value from the list.
		const [newConfigValue] = newValues.splice(idx, 1);

		const res = await fetch(`/api/v1/config-values/${environmentName}/${newConfigValue.Name}`, {
			method: 'POST',
			body: JSON.stringify(newConfigValue)
		});

		if (!res.ok) return;
		newValues = newValues;
		fetchConfig(environmentName);
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

	const getExcludedConfigKeys = () => {
		// Exclude keys which are already configured for this environment.
		const directlyConfiguredKeys = configuration.reduce((collector, cv) => {
			return cv.Inherited ? collector : [cv.Name, ...collector];
		}, []);
		// Exclude keys which are already in the newValues list so we don't have
		// duplicates.
		const currentlyConfiguringKeys = newValues.reduce(
			(collector, cv) => (cv.Name === '' ? collector : [cv.Name, ...collector]),
			[]
		);
		return [...directlyConfiguredKeys, ...currentlyConfiguringKeys];
	};

	$: fetchConfig(environmentName);

	let canAddValues = true;
	$: canAddValues = newValues.length !== configuration.filter((cv) => cv.Inherited).length;

	let buttonTitle = '';
	$: {
		buttonTitle = canAddValues
			? 'Add new configuration value'
			: 'All available keys are already configured or being configured!';
	}
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
					<div class="select">
						<!-- Exclude config keys that are already configured directly on this environment. -->
						<ConfigKeySelector
							excludedKeys={getExcludedConfigKeys()}
							preSelectedName={newValue.Name}
							on:updated={updateValueWithNewKey(newValue)}
						/>
					</div>
				</td>
				<td>
					<ConfigValueInput
						valueType={newValue.ValueType}
						value={getValue(newValue)}
						on:updated={(event) => updateValue(newValues[i], event.detail.value)}
					/>
				</td>
				<td></td>
				<td style="text-align: center;">
					<button class="button is-success" on:click={saveNewConfigValue(i)}>
						<span class="icon">
							<FontAwesomeIcon icon={faCheck} />
						</span>
					</button>
					<button class="button is-danger" on:click={removeNewConfigValue(i)}>
						<span class="icon">
							<FontAwesomeIcon icon={faTrash} />
						</span>
					</button>
				</td>
			</tr>
		{/each}

		<tr>
			<td></td>
			<td></td>
			<td></td>
			<td style="text-align: center;">
				<button
					class="button is-success"
					title={buttonTitle}
					disabled={!canAddValues}
					on:click={addNewConfigValue}
				>
					<span class="icon">
						<FontAwesomeIcon icon={faPlus} />
					</span>
				</button>
			</td>
		</tr>
	</tbody>
</table>
