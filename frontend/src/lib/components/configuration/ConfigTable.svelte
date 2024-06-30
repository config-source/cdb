<script>
	import { ValueTypes } from '$lib/config-values';
	import ConfigValueInput from './ConfigValueInput.svelte';

	export let environmentName;

	let configuration = [];
	let newValues = [];

	const addNewConfigValue = () => {
		newValues = [
			...newValues,
			{
				Name: '',
				ValueType: ValueTypes.FLOAT
			}
		];
	};

	const fetchConfig = async (name) => {
		if (name === '') return;

		const res = await fetch(`/api/v1/config-values/${name}`);
		if (!res.ok) return;

		const data = await res.json();
		data.sort((a, b) => {
			if (a.Inherited && !b.Inherited) {
				return 1;
			}

			if (!a.Inherited && b.Inherited) {
				return -1;
			}

			const nameA = a.Name.toUpperCase();
			const nameB = b.Name.toUpperCase();
			if (nameA < nameB) {
				return -1;
			}

			if (nameA > nameB) {
				return 1;
			}

			return 1;
		});
		configuration = data;
	};

	const getValue = (configValue) => {
		const valueKeys = ['StrValue', 'IntValue', 'BoolValue', 'FloatValue'];
		for (const key of valueKeys) {
			const value = configValue[key];
			if (value) {
				return value;
			}
		}

		return 'UNRECOGNISED';
	};

	$: fetchConfig(environmentName);
</script>

<table class="table is-fullwidth is-hoverable">
	<thead>
		<th>Key</th>
		<th>Value</th>
		<th>Inherited From</th>
		<th></th>
	</thead>
	<tbody>
		{#each configuration as configValue}
			<tr>
				<td>{configValue.Name}</td>
				<td>{getValue(configValue)}</td>
				<td>
					{configValue.InheritedFrom}
				</td>
				<td></td>
			</tr>
		{/each}

		{#each newValues as newValue}
			<tr>
				<td>
					<!-- <ConfigKeySelector -->
					<!-- 	on:keySelected={(key) => { -->
					<!-- 		newValue.Name = key.Name; -->
					<!-- 		newValue.ValueType = key.ValueType; -->
					<!-- 	}} -->
					<!-- /> -->
				</td>
				<td>
					<ConfigValueInput
						valueType={newValue.ValueType}
						on:change={(newValue) => {
							switch (newValue.ValueType) {
								case ValueTypes.STRING:
									newValue.StrValue = newValue;
									break;
								case ValueTypes.INTEGER:
									newValue.IntValue = newValue;
									break;
								case ValueTypes.FLOAT:
									newValue.FloatValue = newValue;
									break;
								case ValueTypes.BOOLEAN:
									newValue.FloatValue = newValue;
									break;
								default:
									throw new Error('Somehow reached unreachable code!');
							}
						}}
					/>
				</td>
				<td></td>
				<td>DELETE ME</td>
			</tr>
		{/each}

		<tr>
			<td></td>
			<td></td>
			<td></td>
			<td>
				<button class="button is-success" on:click={addNewConfigValue}>
					<span class="icon">
						<i class="fab fa-plus"></i>
					</span>
				</button>
			</td>
		</tr>
	</tbody>
</table>
