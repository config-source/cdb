<script>
	import { getValue, updateValue } from '$lib/config-values';
	import { setConfigValue } from '$lib/client/config-values';
	import ConfigValueInput from '../ConfigValueInput.svelte';
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faPencil, faX, faCheck } from '@fortawesome/free-solid-svg-icons';
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

<tr>
	<td>{configValue.Name}</td>
	<td>
		{#if editing}
			<ConfigValueInput
				valueType={configValue.ValueType}
				value={getValue(configValue)}
				on:updated={(event) => updateValue(configValue, event.detail.value)}
			/>
		{:else}
			{getValue(configValue)}
		{/if}
	</td>
	<td>
		{#if !editing}
			{configValue.InheritedFrom}
		{/if}
	</td>
	<td class="buttons is-centered">
		{#if editing}
			<button class="button is-success" onclick={saveEdit(environmentId, configValue)}>
				<span class="icon">
					<FontAwesomeIcon icon={faCheck} />
				</span>
			</button>

			<button class="button is-danger" onclick={() => (editing = false)}>
				<span class="icon">
					<FontAwesomeIcon icon={faX} />
				</span>
			</button>
		{:else}
			<button class="button" onclick={() => (editing = true)}>
				<span class="icon">
					<FontAwesomeIcon icon={faPencil} />
				</span>
			</button>
		{/if}
	</td>
</tr>
