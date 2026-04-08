<script lang="ts">
  import TopNavBar from '$lib/components/TopNavBar.svelte';

  type LeaderTab = 'Top Predictors' | 'Top Creators' | 'Best Streaks' | 'Most Active';
  let activeTab = $state<LeaderTab>('Top Predictors');

  const podium = [
    { rank: 'RANK 02', name: '@AlphaQuant_9', winRate: '91.4%', score: '1,842', order: 'order-2 md:order-1', scale: '' },
    { rank: 'RANK 01', name: '@Nexus_Prime', winRate: '94.8%', score: '2,491', order: 'order-1 md:order-2', scale: 'md:scale-105', leader: true },
    { rank: 'RANK 03', name: '@Sovereign_0x', winRate: '89.2%', score: '1,610', order: 'order-3 md:order-3', scale: '' }
  ];

  const tableRows = [
    { rank: '04', initials: 'VK', name: 'Vector_K', winRate: '88.5%', markets: '1,240', streak: '12🔥', pulse: '1,402' },
    { rank: '05', initials: 'DS', name: 'DeepState_Trader', winRate: '87.1%', markets: '952', streak: '8🔥', pulse: '1,298' },
    { rank: '06', initials: 'OR', name: 'Oracle_Redux', winRate: '86.4%', markets: '2,104', streak: '5🔥', pulse: '1,211' },
    { rank: '07', initials: 'CT', name: 'CyberTheory', winRate: '85.9%', markets: '1,012', streak: '14🔥', pulse: '1,180' },
    { rank: '08', initials: 'MK', name: 'Market_Karma', winRate: '85.2%', markets: '540', streak: '2🔥', pulse: '1,095' }
  ];

  const feed = [
    { icon: 'workspace_premium', text: '@AlphaQuant_9 reached a', highlight: '15-day streak', rest: 'on Commodity Markets.', time: '2m ago' },
    { icon: 'add_chart', text: '@Nexus_Prime created a', highlight: '500-participant', rest: 'Geopolitical Event.', time: '18m ago' },
    { icon: 'rocket_launch', text: '@Sovereign_0x advanced to', highlight: 'Tier 5 Operator', rest: 'status.', time: '1h ago' }
  ];
</script>

<svelte:head>
  <title>Intelligence Leaderboard | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Leaderboard" />

<main class="flex-grow pt-24 pb-16 px-6 max-w-7xl mx-auto w-full">
  <header class="mb-12">
    <h1 class="font-headline text-4xl font-bold tracking-tight mb-2 text-on-surface">Intelligence Leaderboard</h1>
    <p class="text-on-surface-variant font-light max-w-2xl">Visualizing high-fidelity market performance. Tracking the most precise operators within the Pacifica Pulse ecosystem.</p>
  </header>

  <!-- Tabs -->
  <div class="flex flex-wrap gap-2 mb-10 border-b border-outline-variant/15 pb-4">
    {#each ['Top Predictors', 'Top Creators', 'Best Streaks', 'Most Active'] as tab}
      <button
        onclick={() => activeTab = tab as LeaderTab}
        class="px-4 py-2 text-sm font-medium tracking-wide transition-colors {activeTab === tab ? 'text-primary border-b-2 border-primary-container' : 'text-slate-400 hover:text-primary'}"
      >
        {tab}
      </button>
    {/each}
  </div>

  <!-- Podium -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
    {#each podium as op}
      <div class="{op.order} {op.scale} bg-[rgba(50,53,56,0.6)] backdrop-blur-xl p-6 {op.leader ? 'p-8 border-t-4 border-primary-container shadow-[0_0_12px_rgba(0,219,233,0.15)]' : 'border-t-2 border-slate-500/30'} flex flex-col items-center text-center relative">
        {#if op.leader}
          <div class="absolute -top-3 left-1/2 -translate-x-1/2 bg-primary-container text-on-primary-fixed px-3 py-1 text-[10px] font-black uppercase tracking-[0.2em] font-headline">Global Leader</div>
        {/if}
        <div class="relative mb-4 {op.leader ? 'mb-6' : ''}">
          <div class="w-20 h-20 {op.leader ? 'w-24 h-24' : ''} rounded-sm bg-surface-container-highest border {op.leader ? 'border-2 border-primary-container/50' : 'border border-outline-variant/30'} flex items-center justify-center">
            <span class="font-headline font-bold text-2xl text-primary">{op.name[1].toUpperCase()}</span>
          </div>
          <span class="absolute {op.leader ? '-bottom-3' : '-top-2 -right-2'} {op.leader ? 'left-1/2 -translate-x-1/2' : ''} bg-{op.leader ? 'primary-container text-on-primary-fixed' : 'surface-container-highest border border-outline-variant/30'} px-2 py-1 text-[10px] font-bold font-headline">{op.rank}</span>
        </div>
        <h3 class="font-headline text-{op.leader ? '2xl' : 'lg'} font-bold">{op.name}</h3>
        <div class="flex items-center gap-1 my-2 bg-primary-container/10 px-2 py-1 rounded-sm">
          <span class="material-symbols-outlined text-[14px] text-primary-container" style="font-variation-settings: 'FILL' 1;">verified</span>
          <span class="text-[10px] font-bold uppercase tracking-widest text-primary-container">Verified Alpha</span>
        </div>
        <div class="mt-4 w-full grid grid-cols-2 gap-4 border-t border-outline-variant/10 pt-4">
          <div><p class="text-[10px] text-slate-400 uppercase tracking-tighter">Win Rate</p><p class="text-xl font-bold font-headline text-primary">{op.winRate}</p></div>
          <div><p class="text-[10px] text-slate-400 uppercase tracking-tighter">Pulse Score</p><p class="text-xl font-bold font-headline">{op.score}</p></div>
        </div>
      </div>
    {/each}
  </div>

  <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
    <!-- Main Table -->
    <div class="lg:col-span-8">
      <div class="bg-surface-container-low overflow-hidden rounded-sm border border-outline-variant/10">
        <table class="w-full text-left border-collapse">
          <thead class="bg-surface-container-high border-b border-outline-variant/20">
            <tr>
              {#each ['Rank', 'Operator', 'Win Rate', 'Markets', 'Streak', 'Pulse'] as col}
                <th class="px-6 py-4 text-[10px] font-black uppercase tracking-widest text-slate-400 font-headline {col === 'Win Rate' || col === 'Markets' || col === 'Pulse' ? 'text-right' : ''} {col === 'Streak' ? 'text-center' : ''}">{col}</th>
              {/each}
            </tr>
          </thead>
          <tbody class="divide-y divide-outline-variant/10">
            {#each tableRows as row}
              <tr class="hover:bg-surface-container-highest/30 transition-colors cursor-pointer group">
                <td class="px-6 py-4 font-mono text-sm text-slate-400 group-hover:text-primary">{row.rank}</td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3">
                    <div class="w-8 h-8 rounded-sm bg-surface-container-highest border border-outline-variant/20 flex items-center justify-center">
                      <span class="text-[10px] font-bold">{row.initials}</span>
                    </div>
                    <span class="text-sm font-medium">{row.name}</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-right text-primary font-mono text-sm">{row.winRate}</td>
                <td class="px-6 py-4 text-right text-slate-300 font-mono text-sm">{row.markets}</td>
                <td class="px-6 py-4 text-center">
                  <span class="px-2 py-0.5 bg-surface-container-highest text-[10px] font-bold rounded-sm border border-outline-variant/20">{row.streak}</span>
                </td>
                <td class="px-6 py-4 text-right font-bold font-headline text-sm">{row.pulse}</td>
              </tr>
            {/each}
          </tbody>
        </table>
        <div class="p-4 bg-surface-container-lowest text-center">
          <button class="text-xs font-bold uppercase tracking-widest text-slate-500 hover:text-primary transition-colors">Load Detailed Registry</button>
        </div>
      </div>
    </div>

    <!-- Sidebar -->
    <div class="lg:col-span-4 flex flex-col gap-8">
      <!-- Intelligence Feed -->
      <section class="bg-surface-container-low p-6 rounded-sm border border-outline-variant/10">
        <h4 class="font-headline text-sm font-bold uppercase tracking-widest mb-6 border-l-2 border-primary-container pl-3">Intelligence Feed</h4>
        <div class="flex flex-col gap-4">
          {#each feed as item}
            <div class="flex gap-4 items-start pb-4 border-b border-outline-variant/5 last:border-0 last:pb-0">
              <span class="material-symbols-outlined text-primary text-lg" style="font-variation-settings: 'FILL' 1;">{item.icon}</span>
              <div>
                <p class="text-xs leading-relaxed">{item.text} <span class="text-primary-container">{item.highlight}</span> {item.rest}</p>
                <span class="text-[10px] text-slate-500 uppercase mt-1 block">{item.time}</span>
              </div>
            </div>
          {/each}
        </div>
      </section>

      <!-- Network Vitals -->
      <section class="bg-surface-container-low p-6 rounded-sm border border-outline-variant/10">
        <h4 class="font-headline text-sm font-bold uppercase tracking-widest mb-6 border-l-2 border-primary-container pl-3">Network Vitals</h4>
        <div class="space-y-6">
          <div>
            <div class="flex justify-between items-end mb-2">
              <span class="text-[10px] text-slate-500 uppercase tracking-tighter">Total Predictions</span>
              <span class="font-mono text-sm font-bold text-primary">1.28M</span>
            </div>
            <div class="h-1 bg-surface-container-highest w-full rounded-full">
              <div class="h-1 bg-primary-container w-[82%] rounded-full shadow-[0_0_8px_rgba(0,240,255,0.4)]"></div>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-surface-container-highest p-4 rounded-sm">
              <p class="text-[10px] text-slate-500 uppercase mb-1">Avg Win Rate</p>
              <p class="text-xl font-bold font-headline">64.2%</p>
            </div>
            <div class="bg-surface-container-highest p-4 rounded-sm">
              <p class="text-[10px] text-slate-500 uppercase mb-1">Active Labs</p>
              <p class="text-xl font-bold font-headline">412</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</main>

<footer class="w-full py-8 mt-auto bg-[#0b0e11] border-t border-[#3b494b]/15">
  <div class="flex flex-col md:flex-row justify-between items-center px-8 max-w-7xl mx-auto font-body text-xs uppercase tracking-widest">
    <span class="text-slate-500 mb-4 md:mb-0">© 2024 Pacifica Pulse Intelligence</span>
    <div class="flex gap-8">
      {#each ['System Status', 'Legal', 'Privacy Policy', 'Terms of Service'] as link}
        <a class="text-slate-500 hover:text-[#00F0FF] transition-colors" href="/">{link}</a>
      {/each}
    </div>
  </div>
</footer>
