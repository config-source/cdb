<script>
	import CredentialForm from '$lib/components/auth/CredentialForm.svelte';
	import { isError } from '$lib/client';
	import { register } from '$lib/client/auth';
	import { goto } from '$app/navigation';
	import { user } from '$lib/stores/user';

	let errorMessage = $state('');

	/** @type (email: string, password: string) => Promise<void> */
	const onSubmit = async (email, password) => {
		const result = await register(email, password);
		if (isError(result)) {
			errorMessage = result.Message;
		} else {
			user.set({
				fetched: true,
				data: result
			});

			return goto('/');
		}
	};
</script>

<CredentialForm title="Register" {onSubmit} {errorMessage} />
