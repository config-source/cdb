<script>
	/**
	 * @typedef {Object} Props
	 * @property {number} [size]
	 * @property {import('svelte').Snippet} [children]
	 */

	/** @type {Props} */
	let { size = 1, children } = $props();
</script>

<h1 class={`is-size-${size}`}>
	{@render children?.()}
</h1>
