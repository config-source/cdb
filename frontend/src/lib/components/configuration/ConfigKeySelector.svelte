<script>
	import { createEventDispatcher } from 'svelte';

	/**
	 * @typedef {Object} Props
	 * @property {App.ConfigKey[]} excludedKeys
	 * @property {string} preSelectedName
	 */

	/** @type {Props} */
	let { excludedKeys, preSelectedName } = $props();

	/** @type App.ConfigKey[] */
	let configKeys = $state([]);
	/** @type App.ConfigKey | undefined */
	let selectedKey = $state();

	const dispatch = createEventDispatcher();
	const onUpdate = () => dispatch('updated', { value: selectedKey });

	// TODO: should be scoped to the current service
	fetch('/api/v1/config-keys')
		.then((r) => r.json())
		.then((keys) => {
			configKeys = keys.filter((k) => !excludedKeys.includes(k.Name));
			const matchingKey = configKeys.find((k) => k.Name === preSelectedName);
			if (matchingKey) {
				selectedKey = matchingKey;
			} else {
				selectedKey = configKeys[0];
			}

			onUpdate();
		});
</script>

<select bind:value={selectedKey} onchange={onUpdate}>
	{#each configKeys as key}
		<option value={key}>
			{key.Name}
		</option>
	{/each}
</select>
