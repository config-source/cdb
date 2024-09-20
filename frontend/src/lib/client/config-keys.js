/** @type (...serviceNames: string[]) => Promise<App.Response<App.ConfigKey[]>> */
export async function list(...serviceNames) {
	let url = '/api/v1/config-keys';
	if (serviceNames.length > 0) {
		const params = new URLSearchParams();
		for (const serviceID of serviceNames) {
			params.append('service', serviceID.toString());
		}
		url += `?${params.toString()}`;
	}

	const res = await fetch(url);
	return res.json();
}
