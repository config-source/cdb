<script>
	import 'github-markdown-css/github-markdown.css';
	import Heading from '$lib/components/utility/Heading.svelte';
	import SidebarLinks from '$lib/components/utility/SidebarLinks.svelte';
	import SidebarLink from '$lib/components/utility/SidebarLink.svelte';

	const sidebarHeadingSize = 3;
	const sidebarChildSize = sidebarHeadingSize + 1;
</script>

<div class="is-flex is-flex-direction-row is-flex-wrap-nowrap">
	<nav class="sidebar">
		<Heading size={sidebarHeadingSize}>Configuration</Heading>
		<SidebarLinks>
			<SidebarLink size={sidebarChildSize} href="/docs/configuration/keys">Keys</SidebarLink>
			<SidebarLink size={sidebarChildSize} href="/docs/configuration/values">Values</SidebarLink>
		</SidebarLinks>
	</nav>

	<article class="markdown-body mt-4" style="width: 980px;">
		<slot />
	</article>
</div>

<style>
	:global(html),
	:global(body) {
		width: 100%;
		height: 100%;
	}

	.markdown-body {
		box-sizing: border-box;
		min-width: 200px;
		max-width: 980px;
		margin: 0 auto;
		padding: 45px;
	}

	@media (max-width: 767px) {
		.markdown-body {
			padding: 15px;
		}
	}

	.sidebar {
		height: 100%;
		max-width: fit-content;
		padding-top: 1rem;
		padding-left: 3rem;
		padding-right: 2.5rem;
		border-right: 1px solid lightgrey;
		position: fixed;
	}

	div {
		min-width: 100%;
		height: 100%;
	}
</style>
