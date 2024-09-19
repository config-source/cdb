/** @type (environmentId: string, configValue: any) => Promise<boolean> */
export async function setConfigValue(environmentId, configValue) {
	const res = await fetch(`/api/v1/config-values/${environmentId}/${configValue.Name}`, {
		method: 'POST',
		body: JSON.stringify(configValue)
	});

	if (!res.ok) return false;
	return true;
}

/** @type (envId: number) => Promise<any[]> */
export async function fetchConfig(envId) {
	if (envId === 0) return [];

	const res = await fetch(`/api/v1/config-values/${envId}`);
	if (!res.ok) return [];

	/** @type any[] */
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

	return data;
}
