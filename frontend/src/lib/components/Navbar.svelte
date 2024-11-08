<script>
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faGithub } from '@fortawesome/free-brands-svg-icons';
	import { user } from '$lib/stores/user';
	import * as serviceClient from '$lib/client/services';
	import { selectedService } from '$lib/stores/selectedService';
	import { goto } from '$app/navigation';
	import { logout } from '$lib/client/auth';
	import { isError } from '$lib/client';

	/** @type string[] */
	let services = $state(['All']);

	const fetchServices = async () => {
		const data = await serviceClient.list();
		if (isError(data)) return;
		services = ['All', ...data.map((s) => s.Name)];
	};

	fetchServices();
	user.subscribe(() => fetchServices());

	const doLogout = async () => {
		// TODO: handle error scenario
		if (!isError(await logout())) {
			user.set({
				fetched: false,
				data: undefined
			});

			return goto('/auth/login');
		}
	};
</script>

<nav class="navbar" aria-label="main navigation">
	<div class="navbar-brand">
		<a class="navbar-item" href="/"> CDB </a>
	</div>

	<div id="navbarMenu" class="navbar-menu">
		<div class="navbar-start">
			<a class="navbar-item" href="/docs"> Documentation </a>
			<a class="navbar-item" href="https://github.com/config-source/cdb">
				<span class="icon" style="margin-right: 0.2rem">
					<FontAwesomeIcon icon={faGithub} />
				</span>
				Source Code
			</a>
		</div>

		<div class="navbar-end">
			<div class="navbar-item has-dropdown is-hoverable">
				<button class="navbar-link button is-white"> Service: {$selectedService} </button>

				<div class="navbar-dropdown">
					{#each services as service}
						<button
							class="navbar-item button is-white"
							onclick={() => selectedService.set(service)}
						>
							{service}
						</button>
					{/each}
				</div>
			</div>

			<div class="navbar-item has-dropdown is-hoverable">
				<button class="navbar-link button is-white"> {$user.data?.Email} </button>

				<div class="navbar-dropdown">
					<button class="navbar-item button is-white" onclick={doLogout}> Log Out </button>
				</div>
			</div>
		</div>
	</div>
</nav>

<!-- 
Occupies space so that other components don't go under the navbar until scrolling
happens. 
-->
<div class="navbar-spacer"></div>

<style>
	.navbar {
		width: 100%;
		position: fixed;
		border-bottom: 1px solid lightgrey;
	}

	.navbar-spacer {
		height: 53px;
	}

	/* 
	* Bulma doesn't like that we use button tags instead of a tags, but a11y doesn't
	* like having a tags with no href. So this just patches in the hover behaviour of
	* a tags for our buttons in the navbar dropdown.
	*/
	.navbar-dropdown .navbar-item.button {
		display: block;
		width: 100%;
		text-align: left;
	}
	.navbar-dropdown .navbar-item.button:hover {
		--bulma-navbar-item-background-l-delta: var(--bulma-navbar-item-hover-background-l-delta);
		--bulma-navbar-item-background-a: 1;
	}
</style>
