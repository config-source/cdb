/** @type (envName: string) => Promise<App.Response<App.Environment>> */
export async function getByName(envName) {
	return (await fetch(`/api/v1/environments/by-name/${envName}`)).json();
}

/** @type (id: number) => Promise<App.Response<App.Environment>> */
export async function getByID(id) {
	return (await fetch(`/api/v1/environments/by-id/${id}`)).json();
}
