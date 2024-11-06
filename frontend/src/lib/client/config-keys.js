/** @type (...serviceNamesOrIDs: (string | number)[]) => Promise<App.Response<App.ConfigKey[]>> */
export async function list(...serviceNamesOrIDs) {
	let url = '/api/v1/config-keys';
	if (serviceNamesOrIDs.length > 0) {
		const params = new URLSearchParams();
		for (const serviceID of serviceNamesOrIDs) {
			params.append('service', serviceID.toString());
		}
		url += `?${params.toString()}`;
	}

	const res = await fetch(url);
	return res.json();
}
