<script>
	import Heading from '$lib/components/utility/Heading.svelte';
	import { selectedService } from '$lib/stores/selectedService';

	/** @type App.Environment[] */
	let environments = $state([]);
	let query = $state('');

	/** @type App.Environment[] */
	let filteredEnvironments = $derived.by(() => {
		let results =
			query === ''
				? environments
				: environments
						.filter(
							(env) =>
								env.Name.toLowerCase().includes(query.toLowerCase()) ||
								env.Service.toLowerCase().includes(query.toLowerCase())
						)
						// TODO: prioritise things that match the name over things that match the
						// service if both match the query.
						.sort((a, b) => a.Name.localeCompare(b.Name));

		if ($selectedService !== 'All') {
			return results.filter((env) => env.Service === $selectedService);
		}

		return results;
	});

	/** @type (id?: number) => string */
	const getNameFromId = (id) => {
		const env = environments.find((env) => env.ID === id);
		return env?.Name ?? '';
	};

	fetch('/api/v1/environments')
		.then((res) => res.json())
		.then((data) => {
			environments = data;
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
				bind:value={query}
			/>
		</div>
	</div>
	<table class="table is-fullwidth is-hoverable">
		<thead>
			<tr>
				<th>Service</th>
				<th>Name</th>
				<th>Promotes To</th>
			</tr>
		</thead>

		<tbody>
			{#each filteredEnvironments as env}
				<tr>
					<td>
						{env.Service}
					</td>
					<td>
						<a href="/environments/{env.ID}">
							{env.Name}
						</a>
					</td>
					<td>{getNameFromId(env.PromotesToID)}</td>
				</tr>
			{/each}
		</tbody>
	</table>
</div>
