<script lang="ts">
  import Button from './Button.svelte';
  import { guestSession } from '$lib/guest-session';

  let { activePage = '' }: { activePage?: string } = $props();

  const navLinks = [
    { label: 'Markets', href: '/dashboard' },
    { label: 'Portfolio', href: '/portfolio' },
    { label: 'Leaderboard', href: '/leaderboard' },
  ];
</script>

<nav class="fixed top-0 w-full z-50 bg-[#111417]/90 backdrop-blur-md border-b border-[#3b494b]/15">
  <div class="flex justify-between items-center h-16 px-6 max-w-7xl mx-auto">
    <div class="text-xl font-bold tracking-tighter text-[#dbfcff] font-headline">Pacifica Pulse</div>
    <div class="hidden md:flex items-center space-x-8">
      {#each navLinks as link}
        <a
          href={link.href}
          class="font-medium transition-colors duration-200 py-1 {activePage === link.label ? 'text-[#00F0FF] border-b border-[#00F0FF]' : 'text-slate-400 hover:text-[#dbfcff]'}"
        >
          {link.label}
        </a>
      {/each}
    </div>
    <div class="flex items-center gap-4">
      <div class="hidden items-center gap-2 text-[10px] uppercase tracking-[0.2em] text-outline sm:flex">
        {#if $guestSession.status === 'ready' && $guestSession.player}
          <span class="h-2 w-2 rounded-full bg-primary-container shadow-[0_0_8px_rgba(0,240,255,0.4)]"></span>
          <span class="text-primary">{$guestSession.player.displayName}</span>
        {:else if $guestSession.status === 'loading'}
          <span class="h-2 w-2 rounded-full bg-tertiary-fixed-dim animate-pulse"></span>
          <span>Starting Guest</span>
        {:else if $guestSession.status === 'error'}
          <span class="h-2 w-2 rounded-full bg-error"></span>
          <span>Guest Unavailable</span>
        {:else}
          <span class="h-2 w-2 rounded-full bg-outline"></span>
          <span>Guest Session</span>
        {/if}
      </div>
      <Button variant="primary" class="px-5 py-2 text-sm scale-95">Launch Terminal</Button>
    </div>
  </div>
</nav>
