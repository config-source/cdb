/** @type () => Promise<App.Response<App.Service[]>> */
export async function list() {
	const res = await fetch('/api/v1/services');
	return res.json();
}
