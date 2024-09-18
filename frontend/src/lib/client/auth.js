/** @type () => Promise<App.CurrentUserInfo> */
export async function getCurrentUser() {
	const res = await fetch('/api/v1/users/me', { credentials: 'include' });
	if (!res.ok) {
		return {
			loggedIn: false
		};
	}

	return {
		loggedIn: true,
		user: await res.json()
	};
}

/** @type (email: string, password: string) => Promise<App.Response<App.User>> */
export async function login(email, password) {
	const res = await fetch('/api/v1/auth/login', {
		method: 'POST',
		body: JSON.stringify({ Email: email, Password: password })
	});

	return res.json();
}

/** @type (email: string, password: string) => Promise<App.Response<App.User>> */
export async function register(email, password) {
	const res = await fetch('/api/v1/register', {
		method: 'POST',
		body: JSON.stringify({ Email: email, Password: password })
	});

	return res.json();
}

/** @type () => Promise<App.Response<boolean>> */
export async function logout() {
	const res = await fetch('/api/v1/auth/logout', {
		method: 'DELETE'
	});
	if (!res.ok) {
		return res.json();
	}

	return true;
}
