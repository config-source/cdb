<script>
	import EnvTree from './EnvTree.svelte';
	import { createEventDispatcher } from 'svelte';
	import { selectedEnvTreeNode } from '$lib/stores/selectedEnvTreeNode';

	/**
	 * @typedef {Object} Props
	 * @property {any} envTree
	 * @property {number} [depth]
	 */

	/** @type {Props} */
	let { envTree, depth = 0 } = $props();
	const thisEnvironment = envTree.Env;

	let active = $state($selectedEnvTreeNode === thisEnvironment.Name);
	selectedEnvTreeNode.subscribe((name) => (active = name === thisEnvironment.Name));

	const dispatch = createEventDispatcher();
	const onSelect = () => {
		dispatch('envSelected', {
			...envTree.Env
		});
		selectedEnvTreeNode.set(envTree.Env.Name);
	};

	if ($selectedEnvTreeNode === '' && depth === 0) {
		onSelect();
	}
</script>

<div class="level" style:margin-bottom="0.3rem">
	<div class="level-left">
		{#if depth > 0}
			<div style:margin-left={`${depth * 2}rem`} class="level-item">┗</div>
		{/if}

		<button
			type="button"
			class={`level-item button ${active ? 'is-primary' : 'is-outline'}`}
			onclick={onSelect}
		>
			{envTree.Env.Name}
		</button>
	</div>
</div>

{#each envTree.Children as child}
	<EnvTree envTree={child} depth={depth + 1} on:envSelected />
{/each}
