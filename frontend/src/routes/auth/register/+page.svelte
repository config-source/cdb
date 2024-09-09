<script>
	import CredentialForm from '$lib/components/auth/CredentialForm.svelte';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores/user';

	let errorMessage = '';

	/** @type (email: string, password: string) => Promise<void> */
	const onSubmit = async (email, password) => {
		const res = await fetch('/api/v1/register', {
			method: 'POST',
			body: JSON.stringify({ Email: email, Password: password })
		});
		const data = await res.json();
		if (!res.ok) {
			errorMessage = data.message;
		} else {
			user.set({
				fetched: true,
				data: data
			});

			return goto('/');
		}
	};
</script>

<CredentialForm title="Register" {onSubmit} {errorMessage} />
