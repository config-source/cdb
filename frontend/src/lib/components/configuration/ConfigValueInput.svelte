<script>
	import { ValueTypes } from '$lib/config-values';

	export let valueType;
</script>

{#if valueType === ValueTypes.STRING}
	<input type="text" />
{:else if valueType === ValueTypes.BOOLEAN}
	<input type="checkbox" />
{:else if valueType === ValueTypes.INTEGER}
	<input type="number" step="1" pattern="\d+" />
{:else if valueType === ValueTypes.FLOAT}
	<input type="number" step="0.01" pattern="\d+\.\d+" />
{/if}
