<script lang="ts">
  import { onMount } from 'svelte';

  import type { MarketResponse, PositionResponse } from '$lib/api-types';
  import Button from '$lib/components/Button.svelte';
  import { formatAmount } from '$lib/number-display';
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import { ensureGuestSession } from '$lib/guest-session';
  import { loadPortfolioData } from '$lib/portfolio-data';

  type Tab = 'ACTIVE' | 'RESOLVED';
  type LoadStatus = 'loading' | 'ready' | 'error';

  const portfolioState = $state({
    loadStatus: 'loading' as LoadStatus,
    error: null as string | null,
    activeMarkets: [] as MarketResponse[],
    resolvedMarkets: [] as MarketResponse[],
    positions: [] as PositionResponse[],
    availableBalance: '0.00',
    lockedBalance: '0.00'
  });

  let activeTab = $state<Tab>('ACTIVE');

  onMount(() => {
    void loadPage();
  });

  async function loadPage(): Promise<void> {
    portfolioState.loadStatus = 'loading';
    portfolioState.error = null;

    try {
      await ensureGuestSession();
      const data = await loadPortfolioData();

      portfolioState.activeMarkets = data.activeMarkets;
      portfolioState.resolvedMarkets = data.resolvedMarkets;
      portfolioState.positions = data.positions;
      portfolioState.availableBalance = data.availableBalance;
      portfolioState.lockedBalance = data.lockedBalance;
      portfolioState.loadStatus = 'ready';
    } catch (error) {
      portfolioState.loadStatus = 'error';
      portfolioState.error = toErrorMessage(error, 'Unable to load your portfolio right now.');
    }
  }

  function toErrorMessage(error: unknown, fallback: string): string {
    if (error instanceof Error && error.message) {
      return error.message;
    }

    return fallback;
  }

  function activePositions(): PositionResponse[] {
    return portfolioState.positions.filter((position) => position.status === 'open');
  }

  function resolvedPositions(): PositionResponse[] {
    return portfolioState.positions.filter((position) => position.status !== 'open');
  }

  function visiblePositions(): PositionResponse[] {
    return activeTab === 'ACTIVE' ? activePositions() : resolvedPositions();
  }

  function findMarket(marketID: string): MarketResponse | undefined {
    return [...portfolioState.activeMarkets, ...portfolioState.resolvedMarkets].find((market) => market.id === marketID);
  }

  function marketTitle(marketID: string): string {
    return findMarket(marketID)?.title ?? `Market ${marketID}`;
  }

  function marketSymbol(marketID: string): string {
    return findMarket(marketID)?.symbol ?? 'Unknown symbol';
  }

  function marketHref(marketID: string): string {
    const market = findMarket(marketID);
    if (!market) {
      return '/dashboard';
    }

    return market.status === 'resolved' ? `/markets/${marketID}/resolved` : `/markets/${marketID}`;
  }

  function statusLabel(status: string): string {
    return status.replaceAll('_', ' ');
  }

  function statusClass(status: string): string {
    if (status === 'open') {
      return 'text-primary-container';
    }

    if (status === 'won') {
      return 'text-primary';
    }

    if (status === 'lost') {
      return 'text-error';
    }

    return 'text-outline';
  }

  function sideClass(side: string): string {
    return side === 'yes'
      ? 'bg-primary-container/10 text-primary-container border-primary-container/20'
      : 'bg-error/10 text-error border-error/20';
  }

  function formatDateTime(value: string): string {
    return new Intl.DateTimeFormat(undefined, {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    }).format(new Date(value));
  }
</script>

<svelte:head>
  <title>Portfolio | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Portfolio" />

<main class="pt-24 pb-12 px-6 max-w-[1600px] mx-auto grid grid-cols-12 gap-6">
  <section class="col-span-12 grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
    <div class="bg-surface-container-low p-5 flex flex-col justify-between">
      <span class="text-[10px] uppercase tracking-[0.2em] text-outline">Available Balance</span>
      <div class="text-3xl font-headline font-bold text-primary tracking-tighter mt-2">{formatAmount(portfolioState.availableBalance)}</div>
    </div>
    <div class="bg-surface-container-low p-5 flex flex-col justify-between">
      <span class="text-[10px] uppercase tracking-[0.2em] text-outline">Locked Balance</span>
      <div class="text-3xl font-headline font-bold text-on-surface tracking-tighter mt-2">{formatAmount(portfolioState.lockedBalance)}</div>
    </div>
    <div class="bg-surface-container-low p-5 flex flex-col justify-between">
      <span class="text-[10px] uppercase tracking-[0.2em] text-outline">Open Positions</span>
      <div class="text-3xl font-headline font-bold text-primary-container tracking-tighter mt-2">{activePositions().length}</div>
    </div>
    <div class="bg-surface-container-low p-5 flex flex-col justify-between">
      <span class="text-[10px] uppercase tracking-[0.2em] text-outline">Resolved Positions</span>
      <div class="text-3xl font-headline font-bold text-primary-fixed-dim tracking-tighter mt-2">{resolvedPositions().length}</div>
    </div>
  </section>

  <div class="col-span-12 flex items-center justify-between border-b border-outline-variant/15 pb-4">
    <div class="flex gap-8">
      {#each ['ACTIVE', 'RESOLVED'] as tab}
        <button
          onclick={() => activeTab = tab as Tab}
          class="text-primary font-headline font-semibold text-sm tracking-widest pb-4 transition-colors {activeTab === tab ? 'border-b-2 border-primary text-primary' : 'text-outline hover:text-on-surface'}"
        >
          {tab}
        </button>
      {/each}
    </div>
    <div class="flex items-center gap-2 text-xs font-label text-outline uppercase tracking-tighter">
      <span class="material-symbols-outlined text-sm">account_balance_wallet</span> Guest Portfolio
    </div>
  </div>

  <section class="col-span-12 lg:col-span-8 space-y-1">
    <div class="bg-surface-container-low/50 px-6 py-3 grid grid-cols-12 text-[10px] uppercase tracking-[0.2em] text-outline font-bold">
      <div class="col-span-4">Market</div>
      <div class="col-span-2 text-center">Side</div>
      <div class="col-span-2 text-right">Stake</div>
      <div class="col-span-2 text-right">Payout</div>
      <div class="col-span-2 text-right">Status</div>
    </div>

    {#if portfolioState.loadStatus === 'loading'}
      <div class="bg-surface-container p-6 text-sm text-outline">Loading your positions.</div>
    {:else if portfolioState.loadStatus === 'error'}
      <div class="bg-surface-container p-6 border border-error/20 space-y-4">
        <p class="text-sm text-error">{portfolioState.error}</p>
        <Button class="px-5 py-3 text-xs uppercase tracking-[0.2em]" onclick={loadPage}>Try Again</Button>
      </div>
    {:else if visiblePositions().length === 0}
      <div class="bg-surface-container p-6 text-sm text-outline">
        {activeTab === 'ACTIVE' ? 'You do not have any open positions yet.' : 'You do not have any resolved positions yet.'}
      </div>
    {:else}
      {#each visiblePositions() as position}
        <a class="bg-surface-container hover:bg-surface-container-high transition-colors px-6 py-5 grid grid-cols-12 items-center group" href={marketHref(position.market_id)}>
          <div class="col-span-4 flex items-center gap-4">
            <div class="w-10 h-10 bg-surface-container-lowest flex items-center justify-center">
              <span class="material-symbols-outlined text-primary-container text-lg">query_stats</span>
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-headline font-medium text-on-surface">{marketTitle(position.market_id)}</span>
              <span class="text-[10px] text-outline uppercase tracking-widest mt-1">{marketSymbol(position.market_id)}</span>
            </div>
          </div>
          <div class="col-span-2 text-center">
            <span class="{sideClass(position.side)} text-[10px] px-2 py-1 font-bold uppercase tracking-widest border">{position.side}</span>
          </div>
          <div class="col-span-2 text-right font-headline font-bold text-on-surface tracking-tight">
            {formatAmount(position.stake_amount)}
          </div>
          <div class="col-span-2 text-right font-headline font-bold text-primary-fixed-dim">
            {formatAmount(position.potential_payout)}
          </div>
          <div class="col-span-2 text-right">
            <div class="flex flex-col items-end gap-1">
              <span class="text-[10px] font-bold uppercase tracking-widest {statusClass(position.status)}">{statusLabel(position.status)}</span>
              <span class="text-[10px] text-outline font-mono">{formatDateTime(position.created_at)}</span>
            </div>
          </div>
        </a>
      {/each}
    {/if}
  </section>

  <aside class="col-span-12 lg:col-span-4 space-y-6">
    <section class="bg-surface-container-low p-6">
      <h2 class="text-[10px] uppercase tracking-[0.2em] text-outline mb-6">Account Snapshot</h2>
      <div class="space-y-6">
        <div class="flex flex-col gap-2">
          <div class="flex justify-between items-center">
            <span class="text-xs font-headline font-medium uppercase tracking-widest">Active Positions</span>
            <span class="text-primary font-bold">{activePositions().length}</span>
          </div>
          <div class="h-1.5 bg-surface-container-lowest w-full overflow-hidden">
            <div class="h-full bg-primary-container" style="width: {Math.min(activePositions().length * 20, 100)}%"></div>
          </div>
        </div>
        <div class="flex flex-col gap-2">
          <div class="flex justify-between items-center">
            <span class="text-xs font-headline font-medium uppercase tracking-widest">Resolved Positions</span>
            <span class="text-primary font-bold">{resolvedPositions().length}</span>
          </div>
          <div class="h-1.5 bg-surface-container-lowest w-full overflow-hidden">
            <div class="h-full bg-surface-variant" style="width: {Math.min(resolvedPositions().length * 20, 100)}%"></div>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4 pt-6">
        <div class="bg-surface-container-highest p-4 rounded-sm">
          <p class="text-[10px] text-outline uppercase mb-1">Markets Created</p>
          <p class="text-xl font-bold font-headline">{portfolioState.activeMarkets.filter((market) => market.created_by_player_id === portfolioState.positions[0]?.player_id).length + portfolioState.resolvedMarkets.filter((market) => market.created_by_player_id === portfolioState.positions[0]?.player_id).length}</p>
        </div>
        <div class="bg-surface-container-highest p-4 rounded-sm">
          <p class="text-[10px] text-outline uppercase mb-1">Tracked Markets</p>
          <p class="text-xl font-bold font-headline">{portfolioState.positions.length}</p>
        </div>
      </div>
    </section>

    <section class="bg-surface-container-low p-6">
      <h2 class="text-[10px] uppercase tracking-[0.2em] text-outline mb-6">Quick Links</h2>
      <div class="space-y-3">
        <a class="block w-full py-3 bg-surface-container-high hover:bg-surface-bright text-[10px] font-bold uppercase tracking-[0.2em] transition-all border border-outline-variant/20 text-center" href="/dashboard">
          Browse Markets
        </a>
        <a class="block w-full py-3 bg-surface-container-high hover:bg-surface-bright text-[10px] font-bold uppercase tracking-[0.2em] transition-all border border-outline-variant/20 text-center" href="/markets/create">
          Create Market
        </a>
      </div>
    </section>
  </aside>
</main>
