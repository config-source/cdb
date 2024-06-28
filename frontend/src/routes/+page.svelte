<script>
	import EnvTree from '$lib/components/environments/EnvTree.svelte';
	import ConfigTable from '$lib/components/configuration/ConfigTable.svelte';
	import Heading from '$lib/components/utility/Heading.svelte';
	import { selectedEnvTreeNode } from '$lib/stores/selectedEnvTreeNode';

	let envTrees = [];
	let environmentName = '';

	selectedEnvTreeNode.subscribe((name) => (environmentName = name))

	fetch('/api/v1/environments/tree')
		.then((res) => res.json())
		.then((data) => {
			envTrees = data;
			if (data.length > 0) {
				environmentName = data[0].Env.Name;
			}
		});
</script>

<div class="container mt-6">
	<div class="fixed-grid">
		<div class="grid">
			<div class="cell box">
				<Heading size=3>Environments</Heading>
				{#each envTrees as envTree}
					<EnvTree 
						{envTree} 
					/>
				{/each}
			</div>

			<div class="cell box">
				<Heading size=3>Configuration for {environmentName}</Heading>
				<ConfigTable {environmentName} />
			</div>
		</div>
	</div>
</div>
