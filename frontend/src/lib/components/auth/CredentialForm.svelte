<script>
	import { goto } from '$app/navigation';
	import Heading from '$lib/components/utility/Heading.svelte';

	/**
	 * @typedef {Object} Props
	 * @property {string} title
	 * @property {string} [errorMessage]
	 * @property {(email: string, password: string) => unknown} [onSubmit]
	 */

	/** @type {Props} */
	let { title, errorMessage = '', onSubmit } = $props();

	let email = $state('');
	let password = $state('');
	let isLogin = $derived(title === 'Login');

	/** @type (event: Event) => void */
	const handleSubmit = (event) => {
		event.preventDefault();
		onSubmit && onSubmit(email, password);
	};
</script>

<div class="has-text-centered">
	<Heading size={4}>{title}</Heading>
</div>

<form onsubmit={handleSubmit}>
	{#if errorMessage != ''}
		<div class="field">
			{errorMessage}
		</div>
	{/if}

	<div class="field">
		<label class="label" for="email"><b>Email</b></label>
		<div class="control">
			<input
				bind:value={email}
				class="input"
				type="email"
				placeholder="Enter Email"
				name="email"
				required
			/>
		</div>
	</div>

	<div class="field">
		<label for="password"><b>Password</b></label>
		<div class="control">
			<input
				bind:value={password}
				class="input"
				type="password"
				placeholder="Enter Password"
				name="password"
				required
			/>
		</div>
	</div>

	<div class="field is-grouped">
		<div class="control">
			<button class="button is-link" type="submit">
				{isLogin ? 'Login' : 'Register'}
			</button>
		</div>

		<div class="control">
			<button
				class="button is-link"
				onclick={() => goto(isLogin ? '/auth/register' : '/auth/login')}
			>
				{isLogin ? 'Need an account?' : 'Already have an account?'}
			</button>
		</div>
	</div>
</form>
