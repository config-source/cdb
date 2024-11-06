<script>
	import { createEventDispatcher } from 'svelte';
	import { ValueType } from '$lib/config-values';

	let { valueType, value } = $props();

	const dispatch = createEventDispatcher();
	const onUpdate = ({ target }) => {
		let castedValue = target?.value;
		if (valueType === ValueType.BOOLEAN) {
			castedValue = target?.checked;
		} else if (valueType === ValueType.INTEGER) {
			castedValue = parseInt(castedValue, 10);
		} else if (valueType === ValueType.FLOAT) {
			castedValue = parseFloat(castedValue);
		}

		dispatch('updated', { value: castedValue });
	};
</script>

{#if valueType === ValueType.STRING}
	<input class="input" type="text" {value} onchange={onUpdate} />
{:else if valueType === ValueType.BOOLEAN}
	<input type="checkbox" checked={value} onchange={onUpdate} />
{:else if valueType === ValueType.INTEGER}
	<input class="input" type="number" {value} step="1" pattern="\d+" onchange={onUpdate} />
{:else if valueType === ValueType.FLOAT}
	<input class="input" type="number" {value} step="0.01" pattern="\d+\.\d+" onchange={onUpdate} />
{/if}
