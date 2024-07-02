<script>
	import { createEventDispatcher } from 'svelte';

	export let excludedKeys;
	export let preSelectedName;

	let configKeys = [];
	let selectedKey;

	const dispatch = createEventDispatcher();
	const onUpdate = () => dispatch('updated', { value: selectedKey });

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

<select bind:value={selectedKey} on:change={onUpdate}>
	{#each configKeys as key}
		<option value={key}>
			{key.Name}
		</option>
	{/each}
</select>
