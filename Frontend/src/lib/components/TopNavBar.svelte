<script lang="ts">
  import Button from './Button.svelte';
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
      <Button variant="ghost" class="px-4 py-2 text-sm">Log In</Button>
      <Button variant="primary" class="px-5 py-2 text-sm scale-95">Launch Terminal</Button>
    </div>
  </div>
</nav>
