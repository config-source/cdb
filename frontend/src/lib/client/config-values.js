/** @type (environmentName: string, configValue: any) => Promise<boolean> */
export const setConfigValue = async (environmentName, configValue) => {
	const res = await fetch(`/api/v1/config-values/${environmentName}/${configValue.Name}`, {
		method: 'POST',
		body: JSON.stringify(configValue)
	});

	if (!res.ok) return false;
	return true;
};
