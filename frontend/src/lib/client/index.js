/** @type (data: any) => data is App.ApiError */
export function isError(data) {
	return data.Message !== undefined;
}
