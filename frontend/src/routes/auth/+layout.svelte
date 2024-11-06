<script>
	/**
	 * @typedef {Object} Props
	 * @property {import('svelte').Snippet} [children]
	 */

	/** @type {Props} */
	let { children } = $props();
	/** @type number */
	let height = $state();
</script>

<svelte:window bind:innerHeight={height} />

<div
	style={`height: ${height}px`}
	class="is-flex is-flex-direction-row is-align-content-center is-justify-content-center is-align-items-center is-flex-grow-1"
>
	<div>
		{@render children?.()}
	</div>
</div>

<style>
	.is-flex div {
		width: 500px;
	}
</style>
