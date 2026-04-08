<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import Button from '$lib/components/Button.svelte';
</script>

<svelte:head>
  <title>Pacifica Pulse | Technical Prediction Markets</title>
</svelte:head>

<TopNavBar />

<!-- Hero Section -->
<main class="pt-32 pb-20 px-6 max-w-7xl mx-auto">
  <div class="grid lg:grid-cols-2 gap-16 items-center">
    <div class="space-y-8">
      <div class="inline-flex items-center px-3 py-1 bg-surface-container-high border border-outline-variant/20 rounded-sm">
        <span class="w-2 h-2 rounded-full bg-primary-container animate-pulse mr-2"></span>
        <span class="text-[10px] font-bold tracking-[0.1em] uppercase text-primary">Live Terminal Active</span>
      </div>
      <h1 class="text-6xl md:text-7xl font-bold tracking-tighter leading-none text-primary font-headline">
        Pacifica Pulse
      </h1>
      <p class="text-2xl text-on-surface-variant font-light max-w-lg leading-relaxed">
        Trade the Pulse of the Market. A high-fidelity prediction environment for technical traders.
      </p>
      <div class="flex flex-wrap gap-4 pt-4">
        <Button variant="primary" class="px-8 py-4 text-lg">Launch Terminal</Button>
        <Button variant="outline" class="px-8 py-4 text-lg">View Markets</Button>
      </div>
    </div>

    <div class="relative group">
      <div class="absolute -inset-4 bg-primary-container/5 blur-3xl rounded-full group-hover:bg-primary-container/10 transition-all duration-700"></div>
      <div class="relative bg-surface-container-low border border-outline-variant/20 rounded-lg p-4 overflow-hidden terminal-glow">
        <div class="flex items-center justify-between mb-4 pb-2 border-b border-outline-variant/10">
          <div class="flex gap-2">
            <div class="w-3 h-3 rounded-full bg-error/40"></div>
            <div class="w-3 h-3 rounded-full bg-tertiary-fixed-dim/40"></div>
            <div class="w-3 h-3 rounded-full bg-primary-container/40"></div>
          </div>
          <div class="text-[10px] font-mono text-outline uppercase tracking-widest">System_v4.2.0</div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div class="space-y-4">
            <div class="bg-surface-container p-4 rounded-sm border-l-2 border-primary">
              <div class="text-[10px] text-outline uppercase mb-1">BTC Funding Rate</div>
              <div class="text-xl font-mono text-primary-fixed tracking-tight">0.0124%</div>
            </div>
            <div class="bg-surface-container p-4 rounded-sm border-l-2 border-tertiary-fixed-dim">
              <div class="text-[10px] text-outline uppercase mb-1">Open Interest</div>
              <div class="text-xl font-mono text-tertiary-fixed-dim tracking-tight">1.2B USD</div>
            </div>
          </div>
          <div class="bg-surface-container-lowest p-4 rounded-sm flex items-center justify-center border border-outline-variant/10">
            <img
              class="w-full h-full object-cover opacity-50 grayscale hover:grayscale-0 transition-all duration-500"
              alt="Abstract 3D visualization of glowing blue financial data"
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuDBurO6IHXqymGQnO1WsxtgfS979MCFBrS5whziL8_GcOI5G3Hx7RqCVaB0ghUjvBcjNOipVUIKuAZ2-taCe2oI2uTWTN0gnrOsLql95XTrUFTsc1Ynzq57Apkp70q-sbRrlFCJhb451MazjaibWQysI2nIp2qpdYr27qpr3B6dmnXEIPErocIfs-vSjeXkVvo1dRAh-k-vF7J8ijMTxM4mKct6yWyhqzdGfJD456cl6btkJ-mOx9LXo9Rrwh5Zg53LqhnvGMc9Zw-1"
            />
          </div>
        </div>
        <div class="mt-4 p-3 bg-surface-container-highest/30 rounded-sm border border-outline-variant/10">
          <div class="flex justify-between items-center text-[10px] font-mono text-primary mb-2">
            <span>MARKET_PREDICTION_ACTIVE</span>
            <span>98.2% CONFIDENCE</span>
          </div>
          <div class="w-full bg-surface-container-lowest h-1 rounded-full overflow-hidden">
            <div class="bg-primary-container h-full w-3/4"></div>
          </div>
        </div>
      </div>
    </div>
  </div>
</main>

<!-- Value Prop Section -->
<section class="bg-surface-container-low py-24 px-6 border-y border-outline-variant/10">
  <div class="max-w-7xl mx-auto">
    <div class="mb-16">
      <h2 class="text-4xl font-bold tracking-tight mb-4 text-primary font-headline">Prediction Markets for Technical Signals</h2>
      <p class="text-on-surface-variant max-w-2xl text-lg">Harness the power of high-fidelity data signals. Our markets aren't based on opinion, they're built on the undeniable truth of on-chain and perpetual trading metrics.</p>
    </div>
    <div class="grid md:grid-cols-3 gap-1">
      <div class="bg-surface-container p-8 border border-outline-variant/15 hover:bg-surface-container-high transition-colors">
        <span class="material-symbols-outlined text-primary mb-6 scale-125" style="font-variation-settings: 'FILL' 1;">trending_up</span>
        <h3 class="text-xl font-bold mb-3 text-on-surface font-headline">Price Action</h3>
        <p class="text-sm text-on-secondary-container leading-relaxed">Predict volatility breakouts, liquidation cascades, and structural shifts based on sub-second candle data.</p>
      </div>
      <div class="bg-surface-container p-8 border border-outline-variant/15 hover:bg-surface-container-high transition-colors">
        <span class="material-symbols-outlined text-primary mb-6 scale-125" style="font-variation-settings: 'FILL' 1;">account_balance_wallet</span>
        <h3 class="text-xl font-bold mb-3 text-on-surface font-headline">Funding Rates</h3>
        <p class="text-sm text-on-secondary-container leading-relaxed">Trade on the imbalance of the market. Forecast funding flips and perpetual premium shifts before the hourly update.</p>
      </div>
      <div class="bg-surface-container p-8 border border-outline-variant/15 hover:bg-surface-container-high transition-colors">
        <span class="material-symbols-outlined text-primary mb-6 scale-125" style="font-variation-settings: 'FILL' 1;">bar_chart</span>
        <h3 class="text-xl font-bold mb-3 text-on-surface font-headline">Volume Surges</h3>
        <p class="text-sm text-on-secondary-container leading-relaxed">Identify organic vs. inorganic volume growth. Set markets on relative volume profiles and whale absorption events.</p>
      </div>
    </div>
  </div>
</section>

<!-- How It Works -->
<section class="py-24 px-6 max-w-7xl mx-auto">
  <h2 class="text-center text-3xl font-bold mb-16 tracking-tight font-headline">The Signal Lifecycle</h2>
  <div class="grid grid-cols-1 md:grid-cols-4 gap-8 relative">
    <div class="hidden md:block absolute top-12 left-0 w-full h-[1px] bg-gradient-to-r from-transparent via-outline-variant/30 to-transparent -z-10"></div>
    {#each [
      { icon: 'fort', title: 'Signal Forge', desc: 'Define your parameters and market triggers using our logic engine.' },
      { icon: 'rocket_launch', title: 'Deployment', desc: 'Join existing pools or initiate a new signal market with instant liquidity.' },
      { icon: 'monitoring', title: 'Live Monitoring', desc: 'Track real-time resolution odds as market data streams in.' },
      { icon: 'verified', title: 'Automated Settlement', desc: 'Trustless execution via smart oracles immediately upon event closure.' }
    ] as step}
      <div class="flex flex-col items-center text-center space-y-4">
        <div class="w-12 h-12 rounded-sm bg-surface-container-highest flex items-center justify-center border border-primary/20 text-primary">
          <span class="material-symbols-outlined">{step.icon}</span>
        </div>
        <div>
          <h4 class="font-bold text-on-surface font-headline">{step.title}</h4>
          <p class="text-xs text-on-surface-variant px-4 mt-2">{step.desc}</p>
        </div>
      </div>
    {/each}
  </div>
</section>

<!-- Use Cases -->
<section class="py-24 px-6 bg-surface-container-lowest border-y border-outline-variant/10 overflow-hidden">
  <div class="max-w-7xl mx-auto grid md:grid-cols-2 gap-12 items-center">
    <div class="space-y-12">
      <div class="p-8 bg-surface-container-low border-l-4 border-primary">
        <h3 class="text-2xl font-bold mb-4 font-headline">Sentiment Testing for Alpha</h3>
        <p class="text-on-secondary-container leading-relaxed">Instead of guessing where the herd is moving, Pacifica Pulse lets you quantify conviction. Use prediction markets to hedge against technical bias or confirm your trade setup with financial consensus.</p>
      </div>
      <div class="p-8 bg-surface-container-low border-l-4 border-primary">
        <h3 class="text-2xl font-bold mb-4 font-headline">Community Driven Conviction</h3>
        <p class="text-on-secondary-container leading-relaxed">Collaborate with high-frequency traders to build niche markets. From obscure altcoin funding rates to specific CME gap closures, if there is a signal, there is a market.</p>
      </div>
    </div>
    <div class="relative">
      <div class="aspect-square bg-surface-container rounded-sm border border-outline-variant/20 overflow-hidden">
        <img
          class="w-full h-full object-cover mix-blend-overlay opacity-60"
          alt="Sophisticated dark financial interface showing complex charts"
          src="https://lh3.googleusercontent.com/aida-public/AB6AXuAiNd1eiclUlS4Gzn2IdfQ-NITHlpO0nf3_bFiAodl2Cll_tJ2Xc9NetI_wl_fQ_dtaas1rMl5hU_JqBfoNpMVpJ4zeiw4iL5bRBGoWL-RjXBSWP1rY7-5qJbuIiwV4R9Lf0QjJy9XKPk5I0c-TExZp6-INZg4VXl8p6NuAV6p323EWU05MHt_GUwPgeYHFmf527uGEhBLxmhS-wUQ_im1mXoCHgMqC1ZstpCJ5ds_dW8aXUEBwi0htH6uoSvOHukWcSDXEmUJmT-Tz"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-surface-container-lowest via-transparent to-transparent"></div>
        <div class="absolute bottom-8 left-8 right-8">
          <div class="flex items-center gap-2 mb-2">
            <span class="px-2 py-0.5 bg-primary/10 text-primary text-[10px] font-bold tracking-tighter">PREMIUM ANALYTICS</span>
          </div>
          <div class="text-2xl font-bold font-headline text-white">Aggregated conviction flows for 24h volatility.</div>
        </div>
      </div>
    </div>
  </div>
</section>

<!-- Intelligence Layer -->
<section class="py-24 px-6 max-w-7xl mx-auto">
  <div class="flex flex-col md:flex-row gap-12 items-start">
    <div class="md:w-1/3 sticky top-24">
      <div class="p-1 bg-gradient-to-br from-primary/30 to-transparent rounded-sm">
        <div class="bg-surface-container-high p-6 rounded-sm">
          <div class="flex items-center gap-3 mb-6">
            <span class="material-symbols-outlined text-primary">psychology</span>
            <span class="text-sm font-bold uppercase tracking-widest text-on-surface">AI Pulse Insight</span>
          </div>
          <p class="text-sm text-on-surface-variant leading-relaxed mb-6">Our neural engine processes over 10TB of technical data daily to provide a baseline for your predictions.</p>
          <div class="space-y-4">
            {#each [
              { label: 'MARK PRICE', value: '$64,281.04', color: '' },
              { label: 'OPEN INTEREST', value: '$842.1M', color: '' },
              { label: 'VOLATILITY (IV)', value: '42.8%', color: 'text-primary' }
            ] as stat}
              <div class="flex justify-between text-[10px] font-mono">
                <span class="text-outline">{stat.label}</span>
                <span class="{stat.color || 'text-on-surface'}">{stat.value}</span>
              </div>
            {/each}
          </div>
        </div>
      </div>
    </div>
    <div class="md:w-2/3 grid grid-cols-1 gap-6">
      {#each [
        { dot: 'bg-error', label: 'High Volatility Alert', title: 'ETH/USD 5m EMA Crossover', desc: 'Market predicting a bullish flip within the next 45 minutes.', pct: '64% YES', pool: '$142,000' },
        { dot: 'bg-primary', label: 'Stable Signal', title: 'BTC Funding Rate Flip', desc: 'Will BTC funding rate turn negative before UTC midnight?', pct: '22% YES', pool: '$89,500' }
      ] as card}
        <div class="group bg-surface-container-low border border-outline-variant/15 p-6 flex flex-col md:flex-row justify-between items-center gap-6 hover:border-primary/40 transition-all cursor-crosshair">
          <div class="flex-1">
            <div class="flex items-center gap-2 mb-2">
              <span class="w-2 h-2 rounded-full {card.dot}"></span>
              <span class="text-[10px] font-mono {card.dot === 'bg-error' ? 'text-error' : 'text-primary'} uppercase">{card.label}</span>
            </div>
            <h4 class="text-xl font-bold font-headline">{card.title}</h4>
            <p class="text-sm text-on-secondary-container mt-1">{card.desc}</p>
          </div>
          <div class="flex flex-col items-end">
            <div class="text-3xl font-mono font-bold text-primary">{card.pct}</div>
            <div class="text-[10px] text-outline mt-1 uppercase">Pool: {card.pool}</div>
          </div>
        </div>
      {/each}
    </div>
  </div>
</section>

<!-- Footer -->
<footer class="bg-[#111417] border-t border-[#3b494b]/20">
  <div class="grid grid-cols-1 md:grid-cols-4 gap-8 py-12 px-6 max-w-7xl mx-auto text-sm tracking-wide">
    <div class="col-span-1">
      <div class="text-lg font-bold text-[#dbfcff] mb-6 font-headline">Pacifica Pulse</div>
      <div class="inline-flex items-center px-3 py-1 bg-surface-container-low border border-primary/10 rounded-sm">
        <span class="w-2 h-2 rounded-full bg-primary mr-2"></span>
        <span class="text-[10px] font-bold text-primary tracking-tighter uppercase">System Status: Optimal</span>
      </div>
    </div>
    <div class="space-y-4">
      <h5 class="text-[#00F0FF] font-bold uppercase text-[10px] tracking-widest font-headline">Platform</h5>
      <ul class="space-y-2">
        <li><a class="text-slate-500 hover:text-[#00F0FF] underline transition-all" href="/dashboard">Terminal</a></li>
        <li><a class="text-slate-500 hover:text-[#00F0FF] underline transition-all" href="/">API Documentation</a></li>
      </ul>
    </div>
    <div class="space-y-4">
      <h5 class="text-[#00F0FF] font-bold uppercase text-[10px] tracking-widest font-headline">Legal</h5>
      <ul class="space-y-2">
        <li><a class="text-slate-500 hover:text-[#00F0FF] underline transition-all" href="/">Risk Disclosure</a></li>
        <li><a class="text-slate-500 hover:text-[#00F0FF] underline transition-all" href="/">Privacy Policy</a></li>
      </ul>
    </div>
    <div class="space-y-4">
      <h5 class="text-[#00F0FF] font-bold uppercase text-[10px] tracking-widest font-headline">Updates</h5>
      <ul class="space-y-2">
        <li><a class="text-slate-500 hover:text-[#00F0FF] underline transition-all" href="/">System Status</a></li>
      </ul>
      <div class="pt-4 text-slate-600 text-xs leading-relaxed">
        © 2024 Pacifica Pulse Intelligence. All rights reserved.
      </div>
    </div>
  </div>
</footer>
