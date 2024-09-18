// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
	namespace App {
		interface Error {
			Message: string;
		}

		type Response<T> = T | Error;

		interface User {
			Email: string;
		}

		interface Service {
			ID: number;
			Name: string;
			CreatedAt: string;
		}

		interface ConfigKey {
			ID: number;
			Name: string;
			CreatedAt: string;
			ServiceID: number;
			Service: string;
		}

		interface Environment {
			ID: number;
			Name: string;
			PromotesToID?: number;
			ServiceID: number;
			Service: string;
			CreatedAt: string;
			Sensitive: boolean;
		}

		interface CurrentUserInfo {
			loggedIn: boolean;
			user?: User;
		}
	}
}

export {};
