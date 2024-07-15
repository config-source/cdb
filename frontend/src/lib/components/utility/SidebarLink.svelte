<script>
	import Heading from './Heading.svelte';

	export let href;
	export let size = 1;
</script>

<li>
	<a {href}>
		<Heading {size}>
			<slot />
		</Heading>
	</a>
</li>
