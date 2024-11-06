<script>
	import EditableConfigRow from './ConfigTable/EditableConfigRow.svelte';
	import NewValues from './ConfigTable/NewValues.svelte';
	import { fetchConfig } from '$lib/client/config-values';

	/**
	 * @typedef Props
	 * @property {number} environmentId
	 */

	/** @type {Props} */
	let { environmentId } = $props();

	// Stores all the fetched configuration values for this environment.
	/** @type Promise<any[]> */
	let configuration = $derived(fetchConfig(environmentId));
</script>

<table class="table is-fullwidth is-hoverable">
	<thead>
		<tr>
			<th>Key</th>
			<th>Value</th>
			<th>Inherited From</th>
			<th></th>
		</tr>
	</thead>

	<tbody>
		{#await configuration then items}
			{#each items as configValue}
				<EditableConfigRow
					{configValue}
					{environmentId}
					on:updated={() => (environmentId = environmentId)}
				/>
			{/each}

			<NewValues
				{environmentId}
				existingKeys={items.map((i) => i.Name)}
				on:updated={() => (environmentId = environmentId)}
			/>
		{/await}
	</tbody>
</table>
