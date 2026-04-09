<script lang="ts">
	import { goto } from '$app/navigation';

	import Button from './Button.svelte';
	import { guestSession } from '$lib/guest-session';

	let { activePage = '' }: { activePage?: string } = $props();

	const navLinks = [
		{ label: 'Markets', href: '/dashboard' },
		{ label: 'Portfolio', href: '/portfolio' },
		{ label: 'Leaderboard', href: '/leaderboard' }
	];
</script>

<nav class="fixed top-0 z-50 w-full border-b border-[#3b494b]/15 bg-[#111417]/90 backdrop-blur-md">
	<div class="mx-auto flex h-16 max-w-7xl items-center justify-between px-6">
		<div class="font-headline text-xl font-bold tracking-tighter text-[#dbfcff]">
			Pacifica Pulse
		</div>
		<div class="hidden items-center space-x-8 md:flex">
			{#each navLinks as link}
				<a
					href={link.href}
					class="py-1 font-medium transition-colors duration-200 {activePage === link.label
						? 'border-b border-[#00F0FF] text-[#00F0FF]'
						: 'text-slate-400 hover:text-[#dbfcff]'}"
				>
					{link.label}
				</a>
			{/each}
		</div>
		<div class="flex items-center gap-4">
			<div
				class="text-outline hidden items-center gap-2 text-[10px] tracking-[0.2em] uppercase sm:flex"
			>
				{#if $guestSession.status === 'ready' && $guestSession.player}
					<span
						class="bg-primary-container h-2 w-2 rounded-full shadow-[0_0_8px_rgba(0,240,255,0.4)]"
					></span>
					<span class="text-primary">{$guestSession.player.displayName}</span>
				{:else if $guestSession.status === 'loading'}
					<span class="bg-tertiary-fixed-dim h-2 w-2 animate-pulse rounded-full"></span>
					<span>Starting Guest</span>
				{:else if $guestSession.status === 'error'}
					<span class="bg-error h-2 w-2 rounded-full"></span>
					<span>Guest Unavailable</span>
				{:else}
					<span class="bg-outline h-2 w-2 rounded-full"></span>
					<span>Guest Session</span>
				{/if}
			</div>
			<Button
				variant="primary"
				class="scale-95 px-5 py-2 text-sm"
				onclick={() => goto('/dashboard')}>Launch Terminal</Button
			>
		</div>
	</div>
</nav>
