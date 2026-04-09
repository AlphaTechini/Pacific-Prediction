<script lang="ts">
	import '../app.css';
	import { onMount } from 'svelte';

	import { ensureGuestSession } from '$lib/guest-session';

	let { children } = $props();

	onMount(() => {
		void ensureGuestSession();
	});
</script>

<svelte:head>
  <link href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;600;700&family=Inter:wght@300;400;500;600&display=swap" rel="stylesheet" />
  <link href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&display=swap" rel="stylesheet" />
</svelte:head>

<!-- We apply the base background and text color to a wrapping div, or via global css which we just configured as a default layout -->
<div class="bg-surface text-on-surface selection:bg-primary-container/30 min-h-screen font-body text-base">
	{@render children()}
</div>
