/** @type (...serviceIDs: number[]) => Promise<App.Response<App.ConfigKey[]>> */
export async function list(...serviceIDs) {
	let url = '/api/v1/config-keys';
	if (serviceIDs.length > 0) {
		const params = new URLSearchParams();
		for (const serviceID of serviceIDs) {
			params.append('service', serviceID.toString());
		}
		url += params.toString();
	}

	const res = await fetch(url);
	return res.json();
}
