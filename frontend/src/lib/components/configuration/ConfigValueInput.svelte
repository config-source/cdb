<script>
	import { createEventDispatcher } from 'svelte';
	import { ValueType } from '$lib/config-values';

	export let valueType;
	export let value;

	const dispatch = createEventDispatcher();
	const onUpdate = ({ target }) => dispatch('updated', { value: target?.value });
</script>

{#if valueType === ValueType.STRING}
	<input class="input" type="text" {value} on:change={onUpdate} />
{:else if valueType === ValueType.BOOLEAN}
	<input class="input" type="checkbox" {value} on:change={onUpdate} />
{:else if valueType === ValueType.INTEGER}
	<input class="input" type="number" {value} step="1" pattern="\d+" on:change={onUpdate} />
{:else if valueType === ValueType.FLOAT}
	<input class="input" type="number" {value} step="0.01" pattern="\d+\.\d+" on:change={onUpdate} />
{/if}
