<script>
	import CredentialForm from '$lib/components/auth/CredentialForm.svelte';
	import { isError } from '$lib/client';
	import { login } from '$lib/client/auth';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores/user';

	let errorMessage = $state('');

	/** @type (email: string, password: string) => Promise<void> */
	const onSubmit = async (email, password) => {
		const result = await login(email, password);
		if (isError(result)) {
			errorMessage = result.message;
		} else {
			user.set({
				fetched: true,
				data: result
			});

			return goto('/');
		}
	};
</script>

<CredentialForm title="Login" {onSubmit} {errorMessage} />
