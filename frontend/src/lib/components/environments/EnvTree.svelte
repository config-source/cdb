<script>
	import { createEventDispatcher } from 'svelte';
	import { selectedEnvTreeNode } from '$lib/stores/selectedEnvTreeNode';

	export let envTree;
	export let depth = 0;

	let active = false;
	selectedEnvTreeNode.subscribe((id) => (active = id === envTree.env.ID));

	const dispatch = createEventDispatcher();
	const onSelect = () => {
		dispatch('envSelected', {
			...envTree.env
		});
		selectedEnvTreeNode.set(envTree.env.ID);
	};

	if ($selectedEnvTreeNode === 0 && depth === 0) {
		onSelect();
	}
</script>

<div class="level" style:margin-bottom="0">
	<div class="level-left">
		{#if depth > 0}
			<div style:margin-left={`${depth * 2}rem`} class="level-item">┗</div>
		{/if}

		<button
			type="button"
			class={`level-item button ${active ? 'is-primary' : 'is-outline'}`}
			on:click={onSelect}
		>
			{envTree.env.Name}
		</button>
	</div>
</div>

{#each envTree.children as child}
	<svelte:self envTree={child} depth={depth + 1} on:envSelected />
{/each}
