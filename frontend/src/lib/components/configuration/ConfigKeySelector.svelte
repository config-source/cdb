<script>
	import { createEventDispatcher } from 'svelte';
	import * as configKeyClient from '$lib/client/config-keys';
	import { isError } from '$lib/client';

	/**
	 * @typedef {Object} Props
	 * @property {(string | number)[]} service
	 * @property {string[]} excludedKeys
	 * @property {string} preSelectedName
	 */

	/** @type {Props} */
	let { excludedKeys, preSelectedName, service } = $props();

	/** @type App.ConfigKey | undefined */
	let selectedKey = $state();

	const dispatch = createEventDispatcher();
	const onUpdate = () => dispatch('updated', { value: selectedKey });

	/** @type Promise<App.ConfigKey[]> */
	let configKeys = $derived(
		configKeyClient.list(...service).then((keys) => {
			if (isError(keys)) throw new Error(keys.Message);

			const allowedKeys = keys.filter((k) => !excludedKeys.includes(k.Name));
			selectedKey = allowedKeys.find((k) => k.Name === preSelectedName) || allowedKeys[0];

			onUpdate();
			return allowedKeys;
		})
	);
</script>

<select bind:value={selectedKey} onchange={onUpdate}>
	{#await configKeys then keys}
		{#each keys as key}
			<option value={key}>
				{key.Name}
			</option>
		{/each}
	{/await}
</select>
