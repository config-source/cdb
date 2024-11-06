<script>
	import { ValueType, getValue, initialiseValue, updateValue } from '$lib/config-values';
	import { setConfigValue } from '$lib/client/config-values';
	import ConfigValueInput from '../ConfigValueInput.svelte';
	import ConfigKeySelector from '../ConfigKeySelector.svelte';
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faCheck, faPlus, faTrash } from '@fortawesome/free-solid-svg-icons';
	import { createEventDispatcher } from 'svelte';
	import * as configKeyClient from '$lib/client/config-keys';
	import { isError } from '$lib/client';
	import * as envClient from '$lib/client/environments';

	/**
	 * @typedef {Object} Props
	 * @property {any[]} existingKeys
	 * @property {number} environmentId
	 */

	/** @type {Props} */
	let { existingKeys, environmentId } = $props();

	/** @type Promise<App.Environment> */
	let environment = $derived.by(async () => {
		const env = await envClient.getByID(environmentId);
		if (isError(env)) throw new Error(env.Message); // TODO: error handling
		return env;
	});

	// Stores all the new configuration values that are being added by the user.
	/** @type any[] */
	let newValues = $state([]);

	const dispatch = createEventDispatcher();

	const addNewConfigValue = () => {
		newValues = [
			...newValues,
			{
				Name: '',
				ValueType: ValueType.FLOAT
			}
		];
	};

	/** @type (idx: number) => () => void */
	const removeNewConfigValue = (idx) => () => {
		newValues.splice(idx, 1);
		newValues = newValues;
	};

	/** @type (idx: number) => () => void */
	const saveNewConfigValue = (idx) => async () => {
		const { Name } = await environment;
		const configValue = newValues[idx];
		const successful = await setConfigValue(Name, configValue);
		if (successful) {
			// Remove the config value from the list.
			removeNewConfigValue(idx)();
			dispatch('updated', { configValue });
		}
	};

	/** @type (configValue: any) => (event: any) => void */
	const updateValueWithNewKey =
		(configValue) =>
		({ detail }) => {
			const newKey = detail.value;
			if (!newKey) return;

			configValue.Name = newKey.Name;
			configValue.ValueType = newKey.ValueType;
			configValue.ConfigKeyID = newKey.ID;
			initialiseValue(configValue);

			dispatch('newKeySelected', { configValue, newKey });

			newValues = newValues;
		};

	const getExcludedConfigKeys = () => {
		// Exclude keys which are already in the newValues list so we don't have
		// duplicates.
		const currentlyConfiguringKeys = newValues.reduce(
			(collector, cv) => (cv.Name === '' ? collector : [cv.Name, ...collector]),
			[]
		);
		// Exclude keys which are already configured for this environment. Note:
		// this includes inherited keys because the proper way to change an
		// inherited key is by "editing" it.
		return [...existingKeys, ...currentlyConfiguringKeys];
	};

	/** @type Promise<App.ConfigKey[]> */
	let configKeys = $derived.by(async () => {
		const { ServiceID } = await environment;
		if (ServiceID === 0) return [];

		const keys = await configKeyClient.list(ServiceID);
		if (isError(keys)) throw new Error(keys.Message); // TODO: error handling.
		return keys;
	});

	let canAddValues = $derived.by(
		async () => (await configKeys).length !== existingKeys.length + newValues.length
	);

	let buttonTitle = $derived.by(async () =>
		(await canAddValues)
			? 'Add new configuration value'
			: 'All available keys are already configured or being configured!'
	);
</script>

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
			<button class="button is-success" onclick={saveNewConfigValue(i)}>
				<span class="icon">
					<FontAwesomeIcon icon={faCheck} />
				</span>
			</button>
			<button class="button is-danger" onclick={removeNewConfigValue(i)}>
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
		{#await buttonTitle then title}
			<button
				{title}
				class="button is-success"
				disabled={!canAddValues}
				onclick={addNewConfigValue}
			>
				<span class="icon">
					<FontAwesomeIcon icon={faPlus} />
				</span>
			</button>
		{/await}
	</td>
</tr>
