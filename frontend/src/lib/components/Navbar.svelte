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
	let services = ['All'];

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
				<a class="navbar-link"> Service: {$selectedService} </a>

				<div class="navbar-dropdown">
					{#each services as service}
						<a class="navbar-item" on:click={() => selectedService.set(service)}>
							{service}
						</a>
					{/each}
				</div>
			</div>

			<div class="navbar-item has-dropdown is-hoverable">
				<a class="navbar-link"> {$user.data?.Email} </a>

				<div class="navbar-dropdown">
					<a class="navbar-item" on:click={doLogout}> Log Out </a>
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
</style>
