/** @type (data: any) => data is App.Error */
export function isError(data) {
	return data.Message !== undefined;
}
