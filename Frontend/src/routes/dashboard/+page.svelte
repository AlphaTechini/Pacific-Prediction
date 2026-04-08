<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import MarketCard from '$lib/components/MarketCard.svelte';

  let activeFilter = $state<'ALL' | 'CRYPTO' | 'MACRO'>('CRYPTO');

  const markets = [
    {
      id: 'btc-fund-1',
      title: 'BTC Funding Rate Flip',
      subtitle: 'Next 4H Window',
      icon: 'currency_bitcoin',
      badge: 'HIGH VOLATILITY',
      badgeVariant: 'danger' as const,
      yesLabel: 'YES',
      yesOdds: '1.84x',
      noLabel: 'NO',
      noOdds: '2.12x',
      timer: '02:41:09',
      oi: '$12.4M',
      aiInsight: 'Historical funding cycles suggest 72% probability of flip.',
      topBorder: 'border-primary-container'
    },
    {
      id: 'eth-ema-1',
      title: 'ETH/USD 5m EMA',
      subtitle: 'Crossover Event',
      icon: 'equalizer',
      badge: 'TRENDING',
      badgeVariant: 'info' as const,
      yesLabel: 'BULLISH',
      yesOdds: '1.45x',
      noLabel: 'BEARISH',
      noOdds: '2.80x',
      timer: '00:14:52',
      oi: '$4.8M',
      aiInsight: 'Momentum divergence detected on lower timeframes.',
    },
    {
      id: 'sol-vol-1',
      title: 'SOL Volume Spike',
      subtitle: 'Prediction',
      icon: 'bolt',
      badge: 'LOW LIQUIDITY',
      badgeVariant: 'warning' as const,
      yesLabel: 'REACHED',
      yesOdds: '3.10x',
      noLabel: 'FAILED',
      noOdds: '1.18x',
      timer: '12:01:44',
      oi: '$2.1M',
      aiInsight: 'Low liquidity depth suggests higher volatility than expected.',
    }
  ];

  const liveSignals = [
    { asset: 'BTC/USD', price: '$64,281', funding: '-0.002%', fundingColor: 'text-error', oi: '+1.2%', oiColor: 'text-primary-container', vol: '$42.8B', dir: 'trending_up', dirColor: 'text-primary-container' },
    { asset: 'ETH/USD', price: '$3,452', funding: '+0.001%', fundingColor: 'text-primary-container', oi: '-0.5%', oiColor: 'text-error', vol: '$18.2B', dir: 'trending_down', dirColor: 'text-error' },
    { asset: 'SOL/USD', price: '$145.20', funding: '0.000%', fundingColor: 'text-on-surface', oi: '+4.8%', oiColor: 'text-primary-container', vol: '$5.1B', dir: 'trending_up', dirColor: 'text-primary-container' },
  ];
</script>

<svelte:head>
  <title>Dashboard | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<!-- SideNavBar -->
<aside class="fixed left-0 top-16 h-[calc(100vh-64px)] w-64 bg-[#111417] border-r border-[#1d2023] flex-col py-4 hidden lg:flex">
  <div class="px-6 py-4 mb-6">
    <div class="flex items-center gap-3">
      <div class="w-10 h-10 bg-surface-container-highest border border-outline-variant flex items-center justify-center">
        <span class="material-symbols-outlined text-primary-container">shield_person</span>
      </div>
      <div>
        <div class="font-headline text-xs font-bold tracking-widest text-primary uppercase">Alpha Operator</div>
        <div class="text-[10px] text-outline uppercase tracking-tighter">Rank #412</div>
      </div>
    </div>
  </div>
  <nav class="flex-1 flex flex-col gap-1">
    {#each [
      { icon: 'sensors', label: 'Live Feed', active: true },
      { icon: 'bolt', label: 'High Volatility', active: false },
      { icon: 'new_releases', label: 'New Markets', active: false },
      { icon: 'timer', label: 'Ending Soon', active: false },
      { icon: 'star', label: 'Watchlist', active: false }
    ] as item}
      <a
        class="flex items-center gap-4 px-6 py-3 {item.active ? 'bg-[#1d2023] text-[#00F0FF] border-l-4 border-[#00F0FF]' : 'text-[#3b494b] hover:bg-[#191c1f] hover:text-[#dbfcff] transition-all'}"
        href="/dashboard"
      >
        <span class="material-symbols-outlined text-sm">{item.icon}</span>
        <span class="font-label text-xs tracking-widest uppercase">{item.label}</span>
      </a>
    {/each}
  </nav>
  <div class="mt-auto px-6 pt-4 border-t border-outline-variant/20">
    <button class="w-full gradient-primary text-on-primary-fixed py-3 text-[10px] font-bold uppercase tracking-[0.2em] mb-4">Trade Now</button>
    <div class="flex flex-col gap-2">
      <a class="flex items-center gap-2 text-[10px] text-outline uppercase hover:text-primary transition-colors" href="/dashboard">
        <span class="material-symbols-outlined text-xs">help</span> Support
      </a>
      <a class="flex items-center gap-2 text-[10px] text-outline uppercase hover:text-primary transition-colors" href="/dashboard">
        <span class="material-symbols-outlined text-xs">code</span> API Documentation
      </a>
    </div>
  </div>
</aside>

<!-- Main Content -->
<main class="lg:ml-64 pt-16 min-h-screen">
  <!-- Summary Strip -->
  <div class="bg-surface-container-low border-b border-outline-variant/10 px-8 py-3 flex flex-wrap items-center gap-8 text-[10px] tracking-[0.15em] font-medium uppercase text-outline">
    <div class="flex items-center gap-2">
      <span class="w-1.5 h-1.5 bg-primary-container rounded-full shadow-[0_0_8px_#00f0ff]"></span>
      <span>Active Markets: <span class="text-on-surface font-bold">142</span></span>
    </div>
    <div class="flex items-center gap-2">
      <span class="material-symbols-outlined text-xs text-error">timer</span>
      <span>Resolving Soon: <span class="text-on-surface font-bold">8</span></span>
    </div>
    <div class="flex items-center gap-2">
      <span class="material-symbols-outlined text-xs text-tertiary-fixed-dim">analytics</span>
      <span>Your Participation: <span class="text-on-surface font-bold">12</span></span>
    </div>
    <div class="flex items-center gap-2 border-l border-outline-variant/20 pl-8">
      <span>Win Rate: <span class="text-primary-container font-bold">68.4%</span></span>
    </div>
    <div class="ml-auto text-primary-fixed-dim/70 flex items-center gap-2">
      <span class="material-symbols-outlined text-sm">wifi_tethering</span>
      NETWORK STATUS: OPTIMAL (14ms)
    </div>
  </div>

  <div class="p-8 grid grid-cols-12 gap-8">
    <!-- Featured Markets -->
    <section class="col-span-12 xl:col-span-8">
      <div class="flex justify-between items-end mb-6">
        <div>
          <h2 class="font-headline text-2xl font-bold tracking-tight text-primary">Featured Live Markets</h2>
          <p class="text-outline text-xs uppercase tracking-widest mt-1">Algorithmic & Volatility Intelligence</p>
        </div>
        <div class="flex gap-2">
          {#each ['ALL', 'CRYPTO', 'MACRO'] as f}
            <button
              onclick={() => activeFilter = f as any}
              class="text-xs px-3 py-1 border transition-colors {activeFilter === f ? 'bg-surface-container-high border-primary-container/50 text-primary-container' : 'bg-surface-container border-outline-variant/30 text-outline hover:text-primary'}"
            >
              {f}
            </button>
          {/each}
        </div>
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        {#each markets as market}
          <MarketCard {market} />
        {/each}
        <!-- Glass Snapshot Card -->
        <div class="bg-[rgba(29,32,35,0.6)] backdrop-blur-xl border border-outline-variant/20 p-5 flex flex-col justify-center items-center text-center">
          <h4 class="font-headline font-bold text-xs uppercase tracking-widest text-primary mb-2">Advanced Liquidity Analysis</h4>
          <p class="text-[10px] text-outline max-w-[200px] mb-4">View real-time depth charts and order flow imbalance for all active prediction pools.</p>
          <button class="text-[10px] font-bold text-primary-container uppercase tracking-[0.2em] hover:text-primary transition-colors flex items-center gap-2">
            Open Terminal <span class="material-symbols-outlined text-sm">arrow_forward</span>
          </button>
        </div>
      </div>

      <!-- Live Signal Table -->
      <div class="mt-8 bg-surface-container-low border border-outline-variant/20 p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="font-headline text-sm font-bold tracking-widest text-primary uppercase">Live Signal Panel</h3>
          <div class="flex items-center gap-2">
            <span class="text-[9px] text-outline">AUTO-UPDATE: ON</span>
            <span class="w-1.5 h-1.5 bg-primary-container rounded-full animate-pulse"></span>
          </div>
        </div>
        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead>
              <tr class="text-[9px] text-outline uppercase tracking-widest border-b border-outline-variant/20">
                <th class="pb-3 font-medium">Asset</th>
                <th class="pb-3 font-medium">Price</th>
                <th class="pb-3 font-medium">1H Funding</th>
                <th class="pb-3 font-medium">OI Change</th>
                <th class="pb-3 font-medium">Vol (24H)</th>
                <th class="pb-3 font-medium">Direction</th>
              </tr>
            </thead>
            <tbody class="font-mono text-xs">
              {#each liveSignals as s}
                <tr class="border-b border-outline-variant/10">
                  <td class="py-3 font-headline font-bold text-on-surface">{s.asset}</td>
                  <td class="py-3">{s.price}</td>
                  <td class="py-3 {s.fundingColor}">{s.funding}</td>
                  <td class="py-3 {s.oiColor}">{s.oi}</td>
                  <td class="py-3">{s.vol}</td>
                  <td class="py-3"><span class="material-symbols-outlined {s.dirColor}">{s.dir}</span></td>
                </tr>
              {/each}
            </tbody>
          </table>
        </div>
      </div>
    </section>

    <!-- Sidebar -->
    <aside class="col-span-12 xl:col-span-4 space-y-8">
      <!-- Active Markets -->
      <div class="bg-surface-container p-6 border-l border-primary-container/20">
        <h3 class="font-headline text-sm font-bold tracking-widest text-primary uppercase mb-6">Your Active Markets</h3>
        <div class="space-y-6">
          {#each [
            { label: 'BTC Dominance Drop', status: 'In Progress', statusColor: 'text-primary-container', pct: 65, barColor: 'bg-primary-container' },
            { label: 'ETH Burn Rate > 1k', status: 'Resolving', statusColor: 'text-error', pct: 92, barColor: 'bg-error' }
          ] as pos}
            <div>
              <div class="flex justify-between text-[10px] uppercase tracking-widest mb-2">
                <span class="text-on-surface">{pos.label}</span>
                <span class="{pos.statusColor}">{pos.status}</span>
              </div>
              <div class="h-1 bg-surface-container-highest w-full overflow-hidden">
                <div class="h-full {pos.barColor}" style="width: {pos.pct}%"></div>
              </div>
            </div>
          {/each}
        </div>
        <div class="mt-8">
          <h4 class="text-[10px] text-outline uppercase tracking-[0.2em] mb-4">Recent Results</h4>
          <div class="space-y-3">
            <div class="flex justify-between items-center bg-surface-container-lowest p-2 border-r-2 border-primary-container">
              <span class="text-[10px] text-on-surface uppercase">USDT Peg stability</span>
              <span class="font-mono text-[10px] text-primary-container">+142.50 USDC</span>
            </div>
            <div class="flex justify-between items-center bg-surface-container-lowest p-2 border-r-2 border-outline-variant">
              <span class="text-[10px] text-on-surface uppercase">Layer 2 Volume Peak</span>
              <span class="font-mono text-[10px] text-outline">-40.00 USDC</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Leaderboard -->
      <div class="bg-surface-container p-6">
        <h3 class="font-headline text-sm font-bold tracking-widest text-primary uppercase mb-6">Leaderboard Preview</h3>
        <div class="space-y-4">
          {#each [
            { rank: '01', initials: 'X7', name: 'OX_REAPER', pct: '+412%' },
            { rank: '02', initials: 'V1', name: 'VOID_WALKER', pct: '+388%' },
            { rank: '03', initials: 'A_', name: 'ALPHA_DOG', pct: '+294%' }
          ] as op}
            <div class="flex items-center gap-3">
              <span class="font-mono text-[10px] text-outline w-4">{op.rank}</span>
              <div class="w-6 h-6 bg-primary-container/20 border border-primary-container/30 flex items-center justify-center text-[8px] font-bold">{op.initials}</div>
              <span class="text-xs text-on-surface flex-1">{op.name}</span>
              <span class="font-mono text-[10px] text-primary-container">{op.pct}</span>
            </div>
          {/each}
        </div>
        <button class="w-full border border-outline-variant/30 mt-6 py-2 text-[9px] font-bold text-outline uppercase tracking-widest hover:border-primary hover:text-primary transition-all">View All Operators</button>
      </div>

      <!-- Trending -->
      <div class="bg-surface-container p-6">
        <h3 class="font-headline text-sm font-bold tracking-widest text-primary uppercase mb-6">Trending Categories</h3>
        <div class="flex flex-wrap gap-2">
          {#each [
            { icon: 'trending_up', label: 'High Yield' },
            { icon: 'schedule', label: 'Short Duration' },
            { icon: 'analytics', label: 'Vol Signals' },
            { icon: 'layers', label: 'L2 Scaling' }
          ] as cat}
            <a class="bg-surface-container-low border border-outline-variant/30 px-3 py-2 text-[10px] text-on-surface uppercase tracking-widest hover:border-primary-container transition-all flex items-center gap-2" href="/dashboard">
              <span class="material-symbols-outlined text-xs text-primary-container">{cat.icon}</span> {cat.label}
            </a>
          {/each}
        </div>
      </div>
    </aside>
  </div>
</main>

<!-- FAB -->
<button class="fixed bottom-8 right-8 gradient-primary text-on-primary-fixed p-4 rounded-lg shadow-[0_0_20px_rgba(0,240,255,0.3)] active:scale-90 transition-transform group">
  <span class="material-symbols-outlined text-2xl group-hover:rotate-12 transition-transform">add_chart</span>
</button>
