<script>
	import { createEventDispatcher } from 'svelte';
	import { selectedEnvTreeNode } from '$lib/stores/selectedEnvTreeNode';

	export let envTree;
	export let depth = 0;
	const thisEnvironment = envTree.Env;

	let active = $selectedEnvTreeNode === thisEnvironment.Name;
	selectedEnvTreeNode.subscribe(
		(name) => (active = name === thisEnvironment.Name)
	);

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
			on:click={onSelect}
		>
			{envTree.Env.Name}
		</button>
	</div>
</div>

{#each envTree.Children as child}
	<svelte:self envTree={child} depth={depth + 1} on:envSelected />
{/each}
