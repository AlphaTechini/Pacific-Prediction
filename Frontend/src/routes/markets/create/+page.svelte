<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import Input from '$lib/components/Input.svelte';

  let selectedAsset = $state<'BTC' | 'ETH' | 'SOL'>('BTC');
  let selectedExpiry = $state<'1H' | '4H' | '24H' | 'CUSTOM'>('24H');
  let visibility = $state<'PUBLIC' | 'PRIVATE'>('PUBLIC');
  let question = $state('');
  let stake = $state('');
</script>

<svelte:head>
  <title>Create Market | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="pt-24 pb-10 px-6 md:px-10 lg:ml-64">
  <div class="max-w-6xl mx-auto">
    <header class="mb-10">
      <h1 class="font-headline text-4xl md:text-5xl font-extrabold tracking-tighter text-primary mb-2">Signal Forge</h1>
      <p class="text-outline text-sm tracking-wide max-w-xl">Architect a new market intelligence vector. Define the parameters, set the conditions, and deploy the terminal.</p>
    </header>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-10">
      <div class="lg:col-span-7 space-y-8">
        <section class="bg-surface-container-low p-6 rounded-sm border-l-2 border-primary-container">
          <label class="block font-headline text-xs uppercase tracking-widest text-outline mb-4">Market Question</label>
          <div class="relative group">
            <textarea bind:value={question} class="w-full bg-surface-container-lowest border-none text-xl font-headline text-on-surface placeholder:text-outline-variant/30 focus:ring-0 p-4 rounded-sm resize-none" placeholder="Will BTC Funding Rate turn negative before UTC midnight?" rows="2"></textarea>
            <div class="absolute bottom-0 left-0 w-0 h-0.5 bg-primary-container transition-all duration-500 group-focus-within:w-full"></div>
          </div>
        </section>

        <section class="bg-surface-container-low p-6 rounded-sm">
          <label class="block font-headline text-xs uppercase tracking-widest text-outline mb-6">Market Configuration</label>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Primary Asset</span>
              <div class="flex gap-2">
                {#each ['BTC', 'ETH', 'SOL'] as asset}
                  <button onclick={() => selectedAsset = asset as any} class="flex-1 py-2 border text-xs font-bold rounded-sm transition-colors {selectedAsset === asset ? 'border-primary-container bg-primary-container/10 text-primary-container' : 'border-outline-variant/30 text-outline'}">
                    {asset}
                  </button>
                {/each}
              </div>
            </div>
            <div>
              <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Market Type</span>
              <select class="w-full bg-surface-container-lowest border-none text-xs font-bold text-on-surface py-2.5 rounded-sm focus:ring-1 focus:ring-primary-container/30">
                <option>Price-based</option><option>Candle-based</option><option selected>Funding-based</option><option>Open-interest-based</option>
              </select>
            </div>
            <div><Input label="Threshold Value" type="text" placeholder="0.00%" /></div>
            <div>
              <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Expiry Duration</span>
              <div class="flex gap-1">
                {#each ['1H', '4H', '24H', 'CUSTOM'] as exp}
                  <button onclick={() => selectedExpiry = exp as any} class="px-3 py-2 text-[10px] font-bold border transition-colors {selectedExpiry === exp ? 'border-primary-container bg-primary-container/10 text-primary-container' : 'bg-surface-container-lowest border-outline-variant/20'}">{exp}</button>
                {/each}
              </div>
            </div>
            <div class="md:col-span-2">
              <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Settlement Source</span>
              <div class="flex items-center bg-surface-container-lowest p-3 rounded-sm border border-outline-variant/10">
                <span class="material-symbols-outlined text-primary-container text-lg mr-3">account_balance</span>
                <span class="text-xs font-medium">Bybit Perpetual Funding Rate (8H Interval)</span>
              </div>
            </div>
          </div>
        </section>

        <section class="bg-surface-container-lowest p-4 border border-outline-variant/20 flex gap-4 items-start">
          <span class="material-symbols-outlined text-tertiary-fixed-dim mt-1">info</span>
          <p class="text-xs text-on-surface leading-relaxed">
            <span class="text-tertiary-fixed-dim font-bold">Settlement Logic:</span>
            This market settles <span class="text-primary-container font-bold">YES</span> if funding is strictly less than 0.00% at the target timestamp. Otherwise settles <span class="text-error font-bold">NO</span>.
          </p>
        </section>

        <section class="bg-surface-container-low p-6 rounded-sm">
          <label class="block font-headline text-xs uppercase tracking-widest text-outline mb-6">Participation Parameters</label>
          <div class="flex flex-col md:flex-row gap-8">
            <div class="flex-1 relative">
              <Input label="Initial Stake (USDT)" type="number" placeholder="500.00" bind:value={stake} />
              <span class="absolute right-3 top-9 text-[10px] text-outline font-bold">USDT</span>
            </div>
            <div class="flex-1">
              <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Visibility</span>
              <div class="flex bg-surface-container-lowest p-1 rounded-sm">
                {#each ['PUBLIC', 'PRIVATE'] as v}
                  <button onclick={() => visibility = v as any} class="flex-1 py-2 text-[10px] font-bold rounded-sm transition-colors {visibility === v ? 'bg-surface-container-high text-primary-container' : 'text-outline'}">{v}</button>
                {/each}
              </div>
            </div>
          </div>
        </section>

        <button class="w-full py-5 gradient-primary text-on-primary-fixed font-headline font-extrabold text-sm tracking-[0.2em] uppercase rounded-sm hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] transition-all active:scale-95">
          Deploy Signal Market
        </button>
      </div>

      <div class="lg:col-span-5 space-y-6">
        <div class="sticky top-24">
          <h2 class="font-headline text-xs uppercase tracking-widest text-outline mb-4">Signal Preview</h2>
          <div class="bg-[rgba(29,32,35,0.6)] backdrop-blur-xl p-6 rounded-sm border border-primary-container/20 shadow-[0_0_12px_rgba(0,219,233,0.3)]">
            <div class="flex justify-between items-start mb-6">
              <div class="flex items-center gap-2 px-2 py-1 bg-surface-container-highest rounded-sm border border-outline-variant/30">
                <span class="text-[10px] font-bold tracking-tighter">{selectedAsset} / FUNDING</span>
              </div>
              <div class="flex flex-col items-end">
                <span class="text-[10px] text-outline uppercase">Time to Expiry</span>
                <span class="text-sm font-mono text-primary-container">23:59:42</span>
              </div>
            </div>
            <h3 class="text-xl font-headline font-bold text-primary mb-8 leading-tight">
              {question || 'Will BTC Funding Rate turn negative before UTC midnight?'}
            </h3>
            <div class="grid grid-cols-2 gap-4 mb-8">
              <div class="bg-surface-container-lowest p-4 border-l border-primary-container">
                <span class="text-[10px] text-outline uppercase block mb-1">Yes Odds</span>
                <span class="text-2xl font-headline font-bold text-on-surface">50.00%</span>
              </div>
              <div class="bg-surface-container-lowest p-4 border-l border-error/40">
                <span class="text-[10px] text-outline uppercase block mb-1">No Odds</span>
                <span class="text-2xl font-headline font-bold text-on-surface">50.00%</span>
              </div>
            </div>
            <div class="flex items-center justify-between pt-4 border-t border-outline-variant/10">
              <span class="text-[10px] text-outline">0 TRADERS</span>
              <span class="text-[10px] font-bold text-primary-container tracking-widest">LIVE SIGNAL</span>
            </div>
          </div>
          <div class="mt-6 p-4 bg-primary-container/5 border border-primary-container/10 rounded-sm relative overflow-hidden">
            <div class="absolute top-0 left-0 w-1 h-full bg-primary-container"></div>
            <div class="flex items-center gap-2 mb-2">
              <span class="material-symbols-outlined text-primary-container text-sm" style="font-variation-settings: 'FILL' 1;">auto_awesome</span>
              <span class="text-[10px] uppercase font-bold text-primary-container tracking-widest">AI Intelligence</span>
            </div>
            <p class="text-[11px] text-on-surface/80 leading-relaxed italic">"Volatility index indicates a 34% probability of flip. Prediction feasibility: High."</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</main>
