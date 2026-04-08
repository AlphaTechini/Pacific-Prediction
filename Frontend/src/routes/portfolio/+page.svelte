<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';

  type Tab = 'ACTIVE' | 'RESOLVED' | 'CREATED MARKETS' | 'ANALYTICS';
  let activeTab = $state<Tab>('ACTIVE');

  const holdings = [
    { id: 'btc-1', title: 'Will BTC Funding turn negative before midnight?', code: 'BTC-FUNDING-VOL-024', side: 'YES', sideColor: 'success', amount: '500.00', odds: '52%', expiry: '04:22:15', status: 'LIVE' },
    { id: 'sol-1', title: 'Will SOL Open Interest exceed $2.5B by UTC rollover?', code: 'SOL-OI-THRESH-99', side: 'NO', sideColor: 'danger', amount: '250.00', odds: '38%', expiry: '01:14:42', status: 'LIVE' }
  ];

  const stats = [
    { label: 'Total Active Positions', value: '12', color: 'text-primary' },
    { label: 'Resolved Markets', value: '148', color: 'text-on-surface' },
    { label: 'Win Rate', value: '64.2%', color: 'text-primary-fixed-dim' },
    { label: 'Current Streak', value: '5 Wins 🔥', color: 'text-tertiary-fixed-dim' },
    { label: 'Total PNL', value: '+$2,450.00', color: 'text-primary-container' },
    { label: 'Participation Score', value: '980', color: 'text-primary' }
  ];
</script>

<svelte:head>
  <title>Portfolio | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Portfolio" />

<main class="pt-24 pb-12 px-6 max-w-[1600px] mx-auto grid grid-cols-12 gap-6">
  <!-- Stats Bento -->
  <section class="col-span-12 grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4 mb-4">
    {#each stats as stat}
      <div class="bg-surface-container-low p-5 flex flex-col justify-between">
        <span class="text-[10px] uppercase tracking-[0.2em] text-outline">{stat.label}</span>
        <div class="text-3xl font-headline font-bold {stat.color} tracking-tighter mt-2">{stat.value}</div>
      </div>
    {/each}
  </section>

  <!-- Tab Navigation -->
  <div class="col-span-12 flex items-center justify-between border-b border-outline-variant/15 pb-4">
    <div class="flex gap-8">
      {#each ['ACTIVE', 'RESOLVED', 'CREATED MARKETS', 'ANALYTICS'] as tab}
        <button
          onclick={() => activeTab = tab as Tab}
          class="text-primary font-headline font-semibold text-sm tracking-widest pb-4 transition-colors {activeTab === tab ? 'border-b-2 border-primary text-primary' : 'text-outline hover:text-on-surface'}"
        >
          {tab}
        </button>
      {/each}
    </div>
    <div class="flex items-center gap-2 text-xs font-label text-outline uppercase tracking-tighter">
      <span class="material-symbols-outlined text-sm">filter_list</span> Filter By Asset
    </div>
  </div>

  <!-- Holdings Table -->
  <section class="col-span-12 lg:col-span-8 space-y-1">
    <div class="bg-surface-container-low/50 px-6 py-3 grid grid-cols-12 text-[10px] uppercase tracking-[0.2em] text-outline font-bold">
      <div class="col-span-5">Market Description</div>
      <div class="col-span-1 text-center">Side</div>
      <div class="col-span-2 text-right">Amount</div>
      <div class="col-span-1 text-right">Odds</div>
      <div class="col-span-2 text-right">Expiry</div>
      <div class="col-span-1 text-right">Status</div>
    </div>

    {#each holdings as h}
      <div class="bg-surface-container hover:bg-surface-container-high transition-colors px-6 py-5 grid grid-cols-12 items-center group">
        <div class="col-span-5 flex items-center gap-4">
          <div class="w-10 h-10 bg-surface-container-lowest flex items-center justify-center">
            <span class="material-symbols-outlined text-primary-container text-lg">currency_bitcoin</span>
          </div>
          <div class="flex flex-col">
            <span class="text-sm font-headline font-medium text-on-surface">{h.title}</span>
            <span class="text-[10px] text-outline uppercase tracking-widest mt-1">{h.code}</span>
          </div>
        </div>
        <div class="col-span-1 text-center">
          <span class="{h.sideColor === 'success' ? 'bg-primary-container/10 text-primary-container border-primary-container/20' : 'bg-error/10 text-error border-error/20'} text-[10px] px-2 py-1 font-bold uppercase tracking-widest border">{h.side}</span>
        </div>
        <div class="col-span-2 text-right font-headline font-bold text-on-surface tracking-tight">{h.amount} <span class="text-[10px] text-outline">USDC</span></div>
        <div class="col-span-1 text-right font-headline font-bold text-primary-fixed-dim">{h.odds}</div>
        <div class="col-span-2 text-right font-headline font-bold text-on-surface tracking-widest font-mono">{h.expiry}</div>
        <div class="col-span-1 text-right">
          <div class="flex items-center justify-end gap-2">
            <span class="w-1.5 h-1.5 bg-primary-container rounded-full animate-pulse"></span>
            <span class="text-[10px] font-bold text-primary-container tracking-widest uppercase">{h.status}</span>
          </div>
        </div>
      </div>
    {/each}

    <!-- Summary Viz -->
    <div class="grid grid-cols-2 gap-1 mt-6">
      <div class="bg-surface-container p-6 relative overflow-hidden h-48">
        <div class="relative z-10">
          <h4 class="text-[10px] uppercase tracking-[0.2em] text-outline mb-6">Asset Concentration</h4>
          <div class="space-y-4">
            {#each [{ label: 'BITCOIN (BTC)', pct: 72 }, { label: 'ETHEREUM (ETH)', pct: 58 }] as asset}
              <div>
                <div class="flex justify-between text-[10px] font-bold mb-1">
                  <span class="text-on-surface uppercase">{asset.label}</span>
                  <span class="text-primary-container">{asset.pct}%</span>
                </div>
                <div class="h-1 bg-surface-container-lowest w-full">
                  <div class="h-full bg-primary-container" style="width: {asset.pct}%"></div>
                </div>
              </div>
            {/each}
          </div>
        </div>
        <div class="absolute -right-12 -bottom-12 opacity-5 scale-150">
          <span class="material-symbols-outlined text-[160px] text-primary">analytics</span>
        </div>
      </div>
      <div class="bg-surface-container p-6 flex flex-col justify-between">
        <div>
          <h4 class="text-[10px] uppercase tracking-[0.2em] text-outline mb-2">Confidence vs Outcome</h4>
          <p class="text-xs text-outline leading-relaxed">High confidence signals (weighted &gt;80%) maintain a historical win rate of 78%.</p>
        </div>
        <div class="flex items-baseline gap-2">
          <span class="text-4xl font-headline font-bold text-primary tracking-tighter">78%</span>
          <span class="text-[10px] font-bold text-primary-container uppercase tracking-widest">ACCURACY</span>
        </div>
      </div>
    </div>
  </section>

  <!-- Right Panels -->
  <aside class="col-span-12 lg:col-span-4 space-y-6">
    <!-- AI Insights -->
    <div class="bg-[rgba(29,32,35,0.6)] backdrop-blur-xl p-6 border-l-2 border-primary-container">
      <div class="flex items-center gap-2 mb-4">
        <span class="material-symbols-outlined text-primary text-xl">psychology</span>
        <h3 class="text-sm font-headline font-bold uppercase tracking-widest text-primary">Predictive Insights</h3>
      </div>
      <div class="bg-surface-container-lowest/40 p-4 mb-4">
        <p class="text-sm leading-relaxed text-on-surface/90 font-label italic">
          "Your performance is strongest in funding-rate volatility markets. Highest accuracy on 'NO' outcomes for SOL open interest signals."
        </p>
      </div>
      <div class="flex flex-col gap-3">
        <div class="flex items-center justify-between text-xs">
          <span class="text-outline uppercase tracking-tighter font-medium">Top Alpha Category</span>
          <span class="text-primary-container font-bold">Funding Rates</span>
        </div>
        <div class="flex items-center justify-between text-xs">
          <span class="text-outline uppercase tracking-tighter font-medium">Strategic Bias</span>
          <span class="text-primary-container font-bold">Short-OI SOL</span>
        </div>
      </div>
    </div>

    <!-- Performance Breakdown -->
    <div class="bg-surface-container-low p-6">
      <h3 class="text-[10px] uppercase tracking-[0.2em] text-outline mb-6">Performance Matrix</h3>
      <div class="space-y-6">
        {#each [
          { label: 'Funding Markets', pct: 80, value: '80% Win Rate', color: 'bg-primary-container' },
          { label: 'Price Prediction', pct: 42, value: '42% Win Rate', color: 'bg-surface-variant' }
        ] as perf}
          <div class="flex flex-col gap-2">
            <div class="flex justify-between items-center">
              <span class="text-xs font-headline font-medium uppercase tracking-widest">{perf.label}</span>
              <span class="text-primary font-bold">{perf.value}</span>
            </div>
            <div class="h-1.5 bg-surface-container-lowest w-full overflow-hidden">
              <div class="h-full {perf.color}" style="width: {perf.pct}%"></div>
            </div>
          </div>
        {/each}
      </div>
      <button class="w-full mt-8 py-3 bg-surface-container-high hover:bg-surface-bright text-[10px] font-bold uppercase tracking-[0.2em] transition-all border border-outline-variant/20">
        Export Performance Audit
      </button>
    </div>
  </aside>
</main>
