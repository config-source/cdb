<script>
	import { getValue, updateValue } from '$lib/config-values';
	import { setConfigValue } from '$lib/client/config-values';
	import ConfigValueInput from '../ConfigValueInput.svelte';
	import { createEventDispatcher } from 'svelte';

	/**
	 * @typedef {Object} Props
	 * @property {any} configValue
	 * @property {any} environmentId
	 * @property {boolean} [editing]
	 */

	/** @type {Props} */
	let { configValue, environmentId, editing = $bindable(false) } = $props();

	const dispatch = createEventDispatcher();
	const saveEdit = (envId, configValue) => async () => {
		if (await setConfigValue(envId, configValue)) {
			editing = false;
			dispatch('updated', { value: configValue });
		}
	};
</script>
