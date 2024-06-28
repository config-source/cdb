<script>
	import Heading from '$lib/components/utility/Heading.svelte';

	let environments = [];
	let filteredEnvironments = [];

	const filterEnvironments = (evt) => {
		const name = evt.target.value;
		if (name === '') {
			filteredEnvironments = environments;
			return;
		}

		filteredEnvironments = environments.filter((env) =>
			env.Name.toLowerCase().includes(name.toLowerCase())
		);
	};

	const getNameFromId = (id) => {
		const env = environments.find((env) => env.ID === id);
		if (env) {
			return env.Name;
		}

		return '';
	};

	fetch('/api/v1/environments')
		.then((res) => res.json())
		.then((data) => {
			environments = data;
			filteredEnvironments = data;
		});
</script>

<div class="container mt-6">
	<Heading size={3}>Environments</Heading>
	<div class="columns">
		<div class="column is-full p-3">
			<input
				class="is-input column is-full"
				type="text"
				placeholder="Filter environments by name"
				on:keyup={filterEnvironments}
			/>
		</div>
	</div>
	<table class="table is-fullwidth is-hoverable">
		<thead>
			<th>Name</th>
			<th>Promotes To</th>
		</thead>
		<tbody>
			{#each filteredEnvironments as env}
				<tr>
					<td>
						<a href="/environments/{env.Name}">
							{env.Name}
						</a>
					</td>
					<td>{getNameFromId(env.PromotesToID)}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
