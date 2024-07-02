<script>
	import { createEventDispatcher } from 'svelte';
	import { ValueTypes } from '$lib/config-values';

	export let valueType;
	export let value;

	const dispatch = createEventDispatcher();
	const onUpdate = ({ target }) => dispatch('updated', { value: target?.value });
</script>

{#if valueType === ValueTypes.STRING}
	<input type="text" {value} on:change={onUpdate} />
{:else if valueType === ValueTypes.BOOLEAN}
	<input type="checkbox" {value} on:change={onUpdate} />
{:else if valueType === ValueTypes.INTEGER}
	<input type="number" {value} step="1" pattern="\d+" on:change={onUpdate} />
{:else if valueType === ValueTypes.FLOAT}
	<input type="number" {value} step="0.01" pattern="\d+\.\d+" on:change={onUpdate} />
{/if}
