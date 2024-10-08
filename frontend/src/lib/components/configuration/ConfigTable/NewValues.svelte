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

	/** @type string[] */
	export let existingKeys;
	/** @type number */
	export let environmentId;

	/** @type App.Environment | undefined */
	let environment;

	/** @type (envName: number) => Promise<void> */
	const fetchEnv = async (envId) => {
		const env = await envClient.getByID(envId);
		if (isError(env)) return; // TODO: error handling
		environment = env;
	};
	$: fetchEnv(environmentId);

	// Stores all the new configuration values that are being added by the user.
	/** @type any[] */
	let newValues = [];

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
		const configValue = newValues[idx];
		const successful = await setConfigValue(environmentName, configValue);
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

	/** @type App.ConfigKey[] */
	let configKeys = [];
	/** @type (serviceID: number) => Promise<void> */
	const fetchConfigKeys = async (serviceID) => {
		if (serviceID === 0) return;
		const keys = await configKeyClient.list(serviceID);
		if (isError(keys)) return; // TODO: error handling.
		configKeys = keys;
	};
	$: fetchConfigKeys(environment?.ServiceID ?? 0);

	let canAddValues = true;
	$: canAddValues = configKeys.length !== existingKeys.length + newValues.length;

	let buttonTitle = '';
	$: {
		buttonTitle = canAddValues
			? 'Add new configuration value'
			: 'All available keys are already configured or being configured!';
	}
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
