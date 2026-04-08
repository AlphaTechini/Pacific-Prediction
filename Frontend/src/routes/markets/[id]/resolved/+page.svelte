<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import Button from '$lib/components/Button.svelte';
</script>

<svelte:head>
  <title>Market Resolved | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="pt-24 pb-20 px-6 max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-6">
  <!-- Left Column -->
  <div class="lg:col-span-8 space-y-6">
    <!-- Market Header -->
    <section class="bg-surface-container p-8 border-l-4 border-primary-container">
      <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
        <div class="space-y-2">
          <div class="flex items-center gap-3">
            <span class="bg-surface-container-highest text-primary-container text-[10px] font-bold px-2 py-0.5 tracking-widest uppercase">RESOLVED</span>
            <span class="text-outline text-xs font-mono">ID: BTC-FUND-241024</span>
          </div>
          <h1 class="text-3xl md:text-4xl font-headline font-bold text-primary tracking-tight leading-none">
            Will BTC Funding Rate turn negative before UTC midnight?
          </h1>
        </div>
        <div class="flex flex-col items-center justify-center bg-surface-container-highest p-6 min-w-[140px] border border-outline-variant/20">
          <span class="text-xs text-outline uppercase tracking-widest font-label mb-1">Final Result</span>
          <span class="text-5xl font-headline font-extrabold text-primary-container">NO</span>
        </div>
      </div>
    </section>

    <!-- Winner Card (active glow) -->
    <section class="bg-[rgba(29,32,35,0.6)] backdrop-blur-xl p-8 border border-primary-container/20 relative overflow-hidden">
      <div class="absolute top-0 right-0 p-4 opacity-10">
        <span class="material-symbols-outlined text-9xl text-primary-container" style="font-variation-settings: 'FILL' 1;">verified</span>
      </div>
      <div class="relative z-10">
        <div class="flex items-center gap-2 text-primary-container mb-4">
          <span class="material-symbols-outlined">account_balance_wallet</span>
          <span class="text-xs font-bold uppercase tracking-widest">Settlement Complete</span>
        </div>
        <h2 class="text-4xl font-headline font-bold text-on-surface mb-2 uppercase">STAKE RECOVERED + PROFIT</h2>
        <div class="flex items-baseline gap-4 mt-6">
          <div class="flex flex-col">
            <span class="text-xs text-outline uppercase tracking-tighter">Initial Stake</span>
            <span class="text-xl font-mono text-secondary">500.00 USDC</span>
          </div>
          <span class="material-symbols-outlined text-outline">arrow_forward</span>
          <div class="flex flex-col">
            <span class="text-xs text-primary uppercase tracking-tighter font-bold">Total Return</span>
            <span class="text-4xl font-mono font-bold text-primary-container">962.50 USDC</span>
          </div>
        </div>
        <div class="mt-8">
          <button class="bg-surface-container-highest hover:bg-surface-bright text-primary text-xs font-bold uppercase tracking-widest px-6 py-3 transition-all">
            Claim Funds to Wallet
          </button>
        </div>
      </div>
    </section>

    <!-- Oracle Audit -->
    <section class="grid grid-cols-1 md:grid-cols-2 gap-px bg-outline-variant/15 border border-outline-variant/15">
      <div class="bg-surface-container-low p-6">
        <div class="flex items-center gap-2 mb-6">
          <span class="material-symbols-outlined text-primary text-sm">fact_check</span>
          <h3 class="text-sm font-headline font-bold uppercase tracking-widest">Oracle Audit</h3>
        </div>
        <div class="space-y-6">
          <div>
            <p class="text-[10px] text-outline uppercase font-label mb-1">Final Observed Value</p>
            <p class="text-2xl font-mono text-on-surface">+0.0210%</p>
          </div>
          <div>
            <p class="text-[10px] text-outline uppercase font-label mb-1">Target Threshold</p>
            <p class="text-lg font-mono text-error">&lt; 0.0000%</p>
          </div>
        </div>
      </div>
      <div class="bg-surface-container-low p-6">
        <div class="flex items-center gap-2 mb-6">
          <span class="material-symbols-outlined text-primary text-sm">database</span>
          <h3 class="text-sm font-headline font-bold uppercase tracking-widest">Data Source</h3>
        </div>
        <div class="space-y-6">
          <div>
            <p class="text-[10px] text-outline uppercase font-label mb-1">Settlement Source</p>
            <p class="text-lg font-body text-on-surface">Bybit Perpetual Funding Rate</p>
          </div>
          <div>
            <p class="text-[10px] text-outline uppercase font-label mb-1">Resolution Time</p>
            <p class="text-lg font-mono text-on-surface">Oct 24, 2024 - 00:00:00 UTC</p>
          </div>
        </div>
      </div>
    </section>
  </div>

  <!-- Right Column -->
  <aside class="lg:col-span-4 space-y-6">
    <!-- Market Intelligence -->
    <div class="bg-surface-container p-6">
      <h3 class="text-xs font-headline font-bold uppercase tracking-widest mb-6 border-b border-outline-variant/10 pb-2">Market Intelligence</h3>
      <div class="space-y-6">
        <div>
          <div class="flex justify-between items-end mb-2">
            <span class="text-xs font-bold text-on-surface">DISTRIBUTION</span>
            <span class="text-[10px] font-mono text-outline">1.2M USDC VOLUME</span>
          </div>
          <div class="flex h-2 w-full bg-surface-container-highest overflow-hidden">
            <div class="h-full bg-outline-variant" style="width: 38%"></div>
            <div class="h-full bg-primary-container" style="width: 62%"></div>
          </div>
          <div class="flex justify-between mt-2 font-mono text-[11px]">
            <span class="text-outline">YES (38%)</span>
            <span class="text-primary-container">NO (62%)</span>
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4 pt-4 border-t border-outline-variant/10">
          <div><p class="text-[10px] text-outline uppercase">Payout Ratio</p><p class="text-xl font-mono text-primary">1.925x</p></div>
          <div><p class="text-[10px] text-outline uppercase">Participants</p><p class="text-xl font-mono text-on-surface">4,812</p></div>
        </div>
      </div>
    </div>

    <!-- Event Lifecycle -->
    <div class="bg-surface-container p-6">
      <h3 class="text-xs font-headline font-bold uppercase tracking-widest mb-6 border-b border-outline-variant/10 pb-2">Event Lifecycle</h3>
      <div class="space-y-6">
        {#each [
          { label: 'Market Settled', time: 'OCT 24, 00:05 UTC', active: true },
          { label: 'Data Verified', time: 'OCT 24, 00:02 UTC', active: true },
          { label: 'Expiry Reached', time: 'OCT 24, 00:00 UTC', active: false },
          { label: 'Trading Active', time: 'OCT 23, 00:00 UTC', active: false },
          { label: 'Created', time: 'OCT 22, 14:12 UTC', active: false }
        ] as event}
          <div class="flex gap-4">
            <div class="flex flex-col items-center">
              <div class="w-3 h-3 rounded-full {event.active ? 'bg-primary-container shadow-[0_0_8px_rgba(0,240,255,0.4)]' : 'bg-primary-container/50'}"></div>
              <div class="w-px h-full {event.active ? 'bg-primary-container/30' : 'bg-outline-variant/20'}"></div>
            </div>
            <div>
              <p class="text-xs font-bold {event.active ? 'text-on-surface' : 'text-on-surface-variant'} uppercase">{event.label}</p>
              <p class="text-[10px] text-outline">{event.time}</p>
            </div>
          </div>
        {/each}
      </div>
    </div>

    <!-- CTAs -->
    <div class="space-y-3">
      <Button variant="primary" class="w-full py-4 text-xs tracking-widest flex items-center justify-center gap-2">
        <span class="material-symbols-outlined text-sm">trending_up</span> Join Similar Markets
      </Button>
      <Button variant="outline" class="w-full py-4 text-xs tracking-widest flex items-center justify-center gap-2">
        <span class="material-symbols-outlined text-sm">pie_chart</span> View Portfolio
      </Button>
    </div>
  </aside>
</main>
