<script>
	import { FontAwesomeIcon } from '@fortawesome/svelte-fontawesome';
	import { faGithub } from '@fortawesome/free-brands-svg-icons';
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

		<div class="navbar-end"></div>
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
