<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import { page } from '$app/stores';

  // Snippet: Place Bet Panel (modularized via Svelte snippet)
  let betAmount = $state('');
  let selectedSide = $state<'YES' | 'NO' | null>(null);
</script>

<svelte:head>
  <title>Market Details | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="pt-24 pb-20 px-6 max-w-[1400px] mx-auto grid grid-cols-1 lg:grid-cols-12 gap-6">
  <!-- Left Column -->
  <div class="lg:col-span-8 space-y-6">
    <!-- Market Header -->
    <header class="bg-surface-container-low p-6 rounded-sm border-l-4 border-primary-container">
      <div class="flex flex-wrap justify-between items-start gap-4 mb-4">
        <div class="space-y-1">
          <div class="flex items-center gap-3">
            <span class="bg-primary-container/10 text-primary-container px-2 py-0.5 text-[10px] font-bold tracking-widest uppercase rounded-sm border border-primary-container/20">Active</span>
            <span class="text-outline text-xs font-mono uppercase tracking-widest">Type: Funding-based</span>
          </div>
          <h1 class="text-2xl md:text-3xl font-bold font-headline leading-tight text-on-surface">
            Will BTC Funding Rate turn negative before UTC midnight?
          </h1>
          <p class="text-primary font-mono text-sm tracking-tight opacity-80">ASSET: BTC/USD PERPETUAL</p>
        </div>
        <div class="bg-surface-container p-3 rounded-sm border border-outline-variant/15 flex flex-col items-end">
          <span class="text-[10px] text-outline uppercase tracking-widest font-bold">Time Remaining</span>
          <span class="text-2xl font-mono font-bold text-primary-container">23:59:42</span>
        </div>
      </div>
    </header>

    <!-- Live Signal Data -->
    <section class="grid grid-cols-2 md:grid-cols-4 gap-px bg-outline-variant/15 border border-outline-variant/15 rounded-sm overflow-hidden">
      {#each [
        { label: 'Mark Price', value: '$64,281.40', sub: 'Live from Cluster-A', subColor: 'text-primary-container/70' },
        { label: 'Funding Rate', value: '+0.0210%', valueColor: 'text-tertiary-fixed-dim', sub: 'Live chart →' },
        { label: 'Open Interest', value: '$12.4B', sub: '24h Change: +1.2%' },
        { label: '24h Volume', value: '$42.8B', sub: 'Exchange Aggregate' }
      ] as stat}
        <div class="bg-surface-container-low p-4">
          <p class="text-[10px] text-outline uppercase tracking-widest mb-1">{stat.label}</p>
          <p class="text-xl font-mono font-semibold {stat.valueColor ?? 'text-on-surface'}">{stat.value}</p>
          <p class="text-[10px] {stat.subColor ?? 'text-outline-variant'} font-mono mt-1">{stat.sub}</p>
        </div>
      {/each}
    </section>

    <!-- Settlement Protocol -->
    <section class="bg-surface-container p-6 rounded-sm space-y-4">
      <div class="flex items-center gap-2 border-b border-outline-variant/15 pb-3">
        <span class="material-symbols-outlined text-primary-container text-lg">terminal</span>
        <h3 class="font-headline font-bold uppercase tracking-wider text-sm">Settlement Protocol</h3>
      </div>
      <div class="grid md:grid-cols-3 gap-6">
        <div class="md:col-span-2">
          <p class="text-on-surface-variant text-sm leading-relaxed">
            This market settles <span class="text-primary-container font-bold">YES</span> if the Bybit BTC funding rate is strictly less than <span class="text-on-surface font-mono">0.00%</span> at the target timestamp. Otherwise, it settles <span class="text-error font-bold">NO</span>.
          </p>
        </div>
        <div class="space-y-3">
          <div><p class="text-[10px] text-outline uppercase tracking-widest font-bold">Source</p><p class="text-xs font-mono text-on-surface">Bybit Perpetual Funding</p></div>
          <div><p class="text-[10px] text-outline uppercase tracking-widest font-bold">Threshold</p><p class="text-xs font-mono text-on-surface">&lt; 0.0000%</p></div>
        </div>
      </div>
    </section>

    <!-- AI Insight -->
    <section class="bg-[rgba(50,53,56,0.6)] backdrop-blur-xl p-5 rounded-sm border border-primary-container/10 flex items-start gap-4">
      <div class="p-2 bg-primary-container/10 rounded-sm">
        <span class="material-symbols-outlined text-primary-container">psychology</span>
      </div>
      <div>
        <h4 class="text-xs font-bold uppercase tracking-widest text-primary-container mb-2">Pulse AI Intelligence</h4>
        <p class="text-sm text-on-surface-variant leading-relaxed">
          Volatility index indicates a <span class="text-on-surface font-bold">34% probability</span> of flip. Funding is currently trending flat.
          <span class="block mt-1 text-primary-fixed-dim">Prediction feasibility: High.</span>
        </p>
      </div>
    </section>

    <!-- Activity Feed -->
    <section class="space-y-4">
      <h3 class="font-headline font-bold uppercase tracking-widest text-xs flex items-center gap-2">
        <span class="w-1 h-3 bg-primary-container"></span> Market Activity
      </h3>
      <div class="space-y-px bg-outline-variant/10 rounded-sm overflow-hidden">
        {#each [
          { icon: 'history', text: 'Large position entered', highlight: 'YES', time: '2m ago' },
          { icon: 'add_circle', text: 'Liquidity added to pool', time: '14m ago' },
          { icon: 'rocket_launch', text: 'Signal market created', time: '1h ago' }
        ] as event}
          <div class="bg-surface-container-low p-4 flex justify-between items-center">
            <div class="flex items-center gap-4">
              <span class="material-symbols-outlined text-xs text-outline">{event.icon}</span>
              <span class="text-sm">{event.text} {#if event.highlight}<span class="text-primary-container">{event.highlight}</span>{/if}</span>
            </div>
            <span class="text-xs font-mono text-outline">{event.time}</span>
          </div>
        {/each}
      </div>
    </section>
  </div>

  <!-- Right Column: Trade Terminal -->
  <div class="lg:col-span-4 space-y-6">
    <!-- Place Bet Panel (snippet) -->
    <section class="bg-surface-container-high p-6 rounded-sm border border-outline-variant/20 shadow-[0_0_12px_rgba(0,219,233,0.3)] sticky top-24">
      <h2 class="text-xs font-bold uppercase tracking-widest text-outline mb-6">Trade Terminal</h2>
      <div class="grid grid-cols-2 gap-3 mb-6">
        <button
          onclick={() => selectedSide = 'YES'}
          class="flex flex-col items-center justify-center p-4 border transition-all rounded-sm group {selectedSide === 'YES' ? 'bg-primary-container/20 border-primary-container' : 'bg-surface-container border-primary-container/30 hover:bg-primary-container/10'}"
        >
          <span class="text-primary-container font-headline font-bold text-lg">YES</span>
          <span class="text-xs text-outline font-mono mt-1">50.00%</span>
        </button>
        <button
          onclick={() => selectedSide = 'NO'}
          class="flex flex-col items-center justify-center p-4 border transition-all rounded-sm group {selectedSide === 'NO' ? 'bg-error/20 border-error' : 'bg-surface-container border-error/30 hover:bg-error/10'}"
        >
          <span class="text-error font-headline font-bold text-lg">NO</span>
          <span class="text-xs text-outline font-mono mt-1">50.00%</span>
        </button>
      </div>
      <div class="space-y-4">
        <div class="space-y-2">
          <div class="flex justify-between items-center px-1">
            <label class="text-[10px] text-outline uppercase tracking-widest font-bold">Amount</label>
            <span class="text-[10px] text-outline uppercase font-mono">Bal: 1,240.50 USDC</span>
          </div>
          <div class="relative group">
            <input bind:value={betAmount} class="w-full bg-surface-container-lowest border-none text-xl font-mono p-4 pr-16 focus:ring-0 focus:outline-none placeholder:text-outline-variant" placeholder="0.00" type="number"/>
            <span class="absolute right-4 top-1/2 -translate-y-1/2 text-sm font-bold text-outline">USDC</span>
            <div class="absolute bottom-0 left-0 h-[2px] w-0 bg-primary-container transition-all duration-300 group-focus-within:w-full"></div>
          </div>
        </div>
        <div class="bg-surface-container p-4 rounded-sm space-y-3">
          <div class="flex justify-between items-center text-xs">
            <span class="text-outline">Estimated Payout</span>
            <span class="text-primary-container font-mono font-bold">$0.00</span>
          </div>
          <div class="flex justify-between items-center text-xs">
            <span class="text-outline">Potential Return</span>
            <span class="text-on-surface font-mono">100%</span>
          </div>
          <div class="pt-2 border-t border-outline-variant/15 flex justify-between items-center text-[10px]">
            <span class="text-outline-variant uppercase">Slippage Tolerance</span>
            <span class="text-on-surface font-mono">0.5%</span>
          </div>
        </div>
        <button class="w-full gradient-primary text-on-primary-fixed py-4 font-headline font-bold uppercase tracking-[0.2em] rounded-sm hover:opacity-90 active:scale-[0.99] transition-all">
          Place Trade
        </button>
      </div>
    </section>

    <!-- Related Markets -->
    <section class="space-y-4">
      <h3 class="font-headline font-bold uppercase tracking-widest text-xs">Related Intelligence</h3>
      <div class="space-y-3">
        {#each [
          { title: 'ETH Funding Flip', time: 'Closing in 4h', prob: '62% Probability' },
          { title: 'SOL OI Surge', time: 'Closing in 12h', prob: '18% Probability' }
        ] as rel}
          <a class="block bg-surface-container-low p-4 border-l-2 border-outline-variant hover:border-primary-container transition-all group" href="/markets/eth-funding">
            <p class="text-xs font-bold text-on-surface mb-1 group-hover:text-primary transition-colors">{rel.title}</p>
            <div class="flex justify-between items-center">
              <span class="text-[10px] text-outline uppercase font-mono">{rel.time}</span>
              <span class="text-xs font-mono text-primary-container">{rel.prob}</span>
            </div>
          </a>
        {/each}
      </div>
    </section>
  </div>
</main>
