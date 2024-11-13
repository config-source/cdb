<script>
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faPencil, faX, faCheck } from '@fortawesome/free-solid-svg-icons';
	import NewValues from './ConfigTable/NewValues.svelte';
	import { fetchConfig, setConfigValues } from '$lib/client/config-values';
	import ConfigValueInput from './ConfigValueInput.svelte';
	import { getValue, updateValue } from '$lib/config-values';

	/**
	 * @typedef Props
	 * @property {number} environmentId
	 */

	/** @type {Props} */
	let { environmentId } = $props();

	let editing = $state(false);
	let query = $state('');

	/** @type (cv: any) => boolean */
	const matchesFilter = (cv) => {
		return cv.Name.toLowerCase().includes(query);
	};

	// Stores all the fetched configuration values for this environment.
	/** @type any[] */
	let configuration = $state([]);

	const saveEdit = async () => {
		configuration = await setConfigValues(environmentId, configuration);
		editing = false;
	};

	$effect(() => {
		fetchConfig(environmentId).then((cfg) => (configuration = cfg));
	});
</script>

<div class="is-flex is-justify-content-right w-full">
	<input
		class="is-input is-flex-grow-1"
		type="text"
		placeholder="Filter configuration by key"
		bind:value={query}
	/>

	<div class="buttons ml-2">
		{#if editing}
			<button class="button is-success" onclick={saveEdit}>
				<span class="icon">
					<FontAwesomeIcon icon={faCheck} />
				</span>

				<span style="margin-left: 2px;"> Save </span>
			</button>

			<button class="button is-danger" onclick={() => (editing = false)}>
				<span class="icon">
					<FontAwesomeIcon icon={faX} />
				</span>

				<span style="margin-left: 2px;"> Cancel </span>
			</button>
		{:else}
			<button class="button" onclick={() => (editing = true)}>
				<span class="icon">
					<FontAwesomeIcon icon={faPencil} />
				</span>

				<span style="margin-left: 2px;"> Edit </span>
			</button>
		{/if}
	</div>
</div>

<table class="table is-fullwidth is-hoverable">
	<thead>
		<tr>
			<th>Key</th>
			<th>Value</th>
			<th>
				{#if !editing}
					Inherited From
				{/if}
			</th>
		</tr>
	</thead>

	<tbody>
		{#each configuration as configValue}
			{#if matchesFilter(configValue)}
				<tr class={editing ? 'is-editing' : ''}>
					<td>{configValue.Name}</td>
					<td>
						{#if editing}
							<ConfigValueInput
								valueType={configValue.ValueType}
								value={getValue(configValue)}
								on:updated={(event) => updateValue(configValue, event.detail.value)}
							/>
						{:else}
							{getValue(configValue)}
						{/if}
					</td>
					<td>
						{#if !editing}
							{configValue.InheritedFrom}
						{/if}
					</td>
				</tr>
			{/if}
		{/each}

		{#if editing}
			<NewValues {environmentId} existingKeys={configuration.map((i) => i.Name)} />
		{/if}
	</tbody>
</table>

<style>
	.is-editing {
		height: 57px;
	}
</style>
