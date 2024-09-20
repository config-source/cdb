/** @type (id: number) => Promise<App.Response<App.Environment>> */
export async function getByID(id) {
	return (await fetch(`/api/v1/environments/by-id/${id}`)).json();
}
