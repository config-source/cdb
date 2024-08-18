<script>
	/** @type number */
	let height;
</script>

<svelte:window bind:innerHeight={height} />

<div
	style={`height: ${height}px`}
	class="is-flex is-flex-direction-row is-align-content-center is-justify-content-center is-align-items-center is-flex-grow-1"
>
	<div>
		<slot />
	</div>
</div>

<style>
	.is-flex div {
		width: 500px;
	}
</style>
