<script>
	// TODO: while functional, this component could use a refactor. Should
	// probably use a store.

	import { ValueType, getValue, updateValue } from '$lib/config-values';
	import ConfigValueInput from './ConfigValueInput.svelte';
	import ConfigKeySelector from './ConfigKeySelector.svelte';
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faCheck, faPencil, faPlus, faTrash, faX } from '@fortawesome/free-solid-svg-icons';

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
	$: fetchConfig(environmentName);

	// Stores all the fetched configuration values for this environment.
	let configuration = [];
	// Stores all the new configuration values that are being added by the user.
	let newValues = [];
	// Used to determine if any configuration value is being edited. Simply maps
	// configvalue names to a boolean indicating their editing status.
	/** @type Map<string, boolean> */
	let editing = new Map();

	const setEditing = (configValue, value) => () => {
		console.log('set edit', configValue.Name, value);
		editing.set(configValue.Name, value);
		editing = editing;
	};

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
		const successful = await setConfigValue(newValues[idx]);
		if (!successful) {
			return;
		}

		// Remove the config value from the list.
		newValues.splice(idx, 1);
		newValues = newValues;
	};

	const saveEdit = (configValue) => async () => {
		const successful = await setConfigValue(configValue);
		if (successful) {
			setEditing(configValue, false)();
		}
	};

	/** @type (configValue: any) => Promise<boolean> */
	const setConfigValue = async (configValue) => {
		const res = await fetch(`/api/v1/config-values/${environmentName}/${configValue.Name}`, {
			method: 'POST',
			body: JSON.stringify(configValue)
		});

		if (!res.ok) return false;
		await fetchConfig(environmentName);
		return true;
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
		// Exclude keys which are already configured for this environment. Note:
		// this includes inherited keys because the proper way to change an
		// inherited key is by "editing" it.
		const directlyConfiguredKeys = configuration.map((cv) => cv.Name);
		// Exclude keys which are already in the newValues list so we don't have
		// duplicates.
		const currentlyConfiguringKeys = newValues.reduce(
			(collector, cv) => (cv.Name === '' ? collector : [cv.Name, ...collector]),
			[]
		);
		return [...directlyConfiguredKeys, ...currentlyConfiguringKeys];
	};

	let configKeys = [];
	const fetchConfigKeys = async () => {
		const res = await fetch('/api/v1/config-keys');
		if (!res.ok) return; // TODO: error handling.
		const data = await res.json();
		configKeys = data;
	};
	fetchConfigKeys();

	let canAddValues = true;
	$: canAddValues = configKeys.length !== configuration.length + newValues.length;

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
		{#each configuration as configValue, i}
			<tr>
				<td>{configValue.Name}</td>
				<td>
					{#if editing.get(configValue.Name)}
						<ConfigValueInput
							valueType={configValue.ValueType}
							value={getValue(configValue)}
							on:updated={(event) => updateValue(configuration[i], event.detail.value)}
						/>
					{:else}
						{getValue(configValue)}
					{/if}
				</td>
				<td>
					{#if !editing.get(configValue.Name)}
						{configValue.InheritedFrom}
					{/if}
				</td>
				<td class="buttons is-centered">
					{#if editing.get(configValue.Name)}
						<button class="button is-success" on:click={saveEdit(configuration[i])}>
							<span class="icon">
								<FontAwesomeIcon icon={faCheck} />
							</span>
						</button>

						<button class="button is-danger" on:click={setEditing(configValue, false)}>
							<span class="icon">
								<FontAwesomeIcon icon={faX} />
							</span>
						</button>
					{:else}
						<button class="button" on:click={setEditing(configValue, true)}>
							<span class="icon">
								<FontAwesomeIcon icon={faPencil} />
							</span>
						</button>
					{/if}
				</td>
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
				<td class="buttons is-centered">
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
			<td class="buttons is-centered">
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
