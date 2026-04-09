<script lang="ts">
  import { onMount } from 'svelte';

  import type { MarketResponse, PositionResponse } from '$lib/api-types';
  import Button from '$lib/components/Button.svelte';
  import MarketCard from '$lib/components/MarketCard.svelte';
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import { loadDashboardData } from '$lib/dashboard-data';
  import { ensureGuestSession, guestSession } from '$lib/guest-session';

  type DashboardStatus = 'loading' | 'ready' | 'error';
  type MarketView = 'ACTIVE' | 'RESOLVED';

  interface DashboardState {
    status: DashboardStatus;
    error: string | null;
    activeMarkets: MarketResponse[];
    resolvedMarkets: MarketResponse[];
    positions: PositionResponse[];
    availableBalance: string;
    lockedBalance: string;
  }

  const dashboardState = $state<DashboardState>({
    status: 'loading',
    error: null,
    activeMarkets: [],
    resolvedMarkets: [],
    positions: [],
    availableBalance: '0.00',
    lockedBalance: '0.00'
  });

  let activeView = $state<MarketView>('ACTIVE');

  onMount(() => {
    void refreshDashboard();
  });

  async function refreshDashboard(): Promise<void> {
    dashboardState.status = 'loading';
    dashboardState.error = null;

    const player = await ensureGuestSession();
    if (!player) {
      dashboardState.status = 'error';
      dashboardState.error = 'I could not start your guest session.';
      return;
    }

    try {
      const data = await loadDashboardData();

      dashboardState.status = 'ready';
      dashboardState.activeMarkets = data.activeMarkets;
      dashboardState.resolvedMarkets = data.resolvedMarkets;
      dashboardState.positions = data.positions;
      dashboardState.availableBalance = data.balance.available_balance;
      dashboardState.lockedBalance = data.balance.locked_balance;
    } catch (error) {
      dashboardState.status = 'error';
      dashboardState.error = toErrorMessage(error);
    }
  }

  function toErrorMessage(error: unknown): string {
    if (error instanceof Error && error.message) {
      return error.message;
    }

    return 'Unable to load the dashboard right now.';
  }

  function visibleMarkets(): MarketResponse[] {
    return activeView === 'ACTIVE' ? dashboardState.activeMarkets : dashboardState.resolvedMarkets;
  }

  function openPositions(): PositionResponse[] {
    return dashboardState.positions.filter((position) => position.status === 'open').slice(0, 5);
  }

  function findMarketTitle(marketID: string): string {
    const market = [...dashboardState.activeMarkets, ...dashboardState.resolvedMarkets].find((item) => item.id === marketID);
    return market?.title ?? `Market ${marketID}`;
  }

  function formatPositionStatus(status: string): string {
    return status.replaceAll('_', ' ');
  }

  function formatPositionSide(side: string): string {
    return side.toUpperCase();
  }
</script>

<svelte:head>
  <title>Dashboard | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="pt-24 pb-12 px-6 md:px-10 lg:px-12 max-w-[1600px] mx-auto">
  <section class="grid grid-cols-1 xl:grid-cols-12 gap-8">
    <div class="xl:col-span-8 space-y-8">
      <header class="space-y-4">
        <div class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
          <div>
            <h1 class="font-headline text-4xl md:text-5xl font-bold tracking-tight text-primary">Live Markets</h1>
            <p class="text-sm text-outline max-w-2xl">This dashboard now reflects the real backend state: active markets, resolved markets, and your current position activity.</p>
          </div>
          <div class="flex gap-2">
            {#each ['ACTIVE', 'RESOLVED'] as view}
              <button
                class="px-4 py-2 text-xs font-bold border uppercase tracking-[0.2em] transition-colors {activeView === view ? 'bg-surface-container-high border-primary-container/50 text-primary-container' : 'bg-surface-container border-outline-variant/30 text-outline hover:text-primary'}"
                onclick={() => activeView = view as MarketView}
              >
                {view}
              </button>
            {/each}
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="bg-surface-container-low p-5 border border-outline-variant/15">
            <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Active Markets</div>
            <div class="mt-2 text-3xl font-headline font-bold text-primary">{dashboardState.activeMarkets.length}</div>
          </div>
          <div class="bg-surface-container-low p-5 border border-outline-variant/15">
            <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Resolved Markets</div>
            <div class="mt-2 text-3xl font-headline font-bold text-on-surface">{dashboardState.resolvedMarkets.length}</div>
          </div>
          <div class="bg-surface-container-low p-5 border border-outline-variant/15">
            <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Open Positions</div>
            <div class="mt-2 text-3xl font-headline font-bold text-primary-container">{openPositions().length}</div>
          </div>
          <div class="bg-surface-container-low p-5 border border-outline-variant/15">
            <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Available Balance</div>
            <div class="mt-2 text-3xl font-headline font-bold text-on-surface">{dashboardState.availableBalance}</div>
          </div>
        </div>
      </header>

      {#if dashboardState.status === 'loading'}
        <section class="bg-surface-container-low border border-outline-variant/20 p-8">
          <p class="text-sm text-outline">Loading markets and your account state.</p>
        </section>
      {:else if dashboardState.status === 'error'}
        <section class="bg-surface-container-low border border-error/20 p-8 space-y-4">
          <p class="text-sm text-error">{dashboardState.error}</p>
          <Button class="px-5 py-3 text-xs uppercase tracking-[0.2em]" onclick={refreshDashboard}>Try Again</Button>
        </section>
      {:else}
        <section class="grid grid-cols-1 md:grid-cols-2 gap-4">
          {#if visibleMarkets().length > 0}
            {#each visibleMarkets() as market}
              <MarketCard {market} />
            {/each}
          {:else}
            <div class="md:col-span-2 bg-surface-container-low border border-outline-variant/20 p-8">
              <p class="text-sm text-outline">There are no {activeView.toLowerCase()} markets to show yet.</p>
            </div>
          {/if}
        </section>
      {/if}
    </div>

    <aside class="xl:col-span-4 space-y-6">
      <section class="bg-surface-container-low p-6 border border-outline-variant/15">
        <div class="flex items-center justify-between mb-4">
          <h2 class="font-headline text-sm font-bold uppercase tracking-widest text-primary">Your Session</h2>
          <button class="text-[10px] uppercase tracking-[0.2em] text-outline hover:text-primary transition-colors" onclick={refreshDashboard}>
            Refresh
          </button>
        </div>
        <div class="space-y-4">
          <div>
            <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Guest Name</div>
            <div class="mt-2 text-lg font-headline font-bold text-on-surface">
              {$guestSession.player?.displayName ?? 'Starting guest session'}
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-surface-container p-4">
              <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Available</div>
              <div class="mt-2 text-xl font-headline font-bold text-primary">{dashboardState.availableBalance}</div>
            </div>
            <div class="bg-surface-container p-4">
              <div class="text-[10px] uppercase tracking-[0.2em] text-outline">Locked</div>
              <div class="mt-2 text-xl font-headline font-bold text-on-surface">{dashboardState.lockedBalance}</div>
            </div>
          </div>
          <div class="flex flex-col gap-3 pt-2">
            <a class="gradient-primary text-on-primary-fixed px-4 py-3 text-[10px] font-bold uppercase tracking-[0.2em] text-center" href="/markets/create">
              Create Market
            </a>
            <a class="border border-outline-variant/30 px-4 py-3 text-[10px] font-bold uppercase tracking-[0.2em] text-center text-outline hover:text-primary hover:border-primary-container/40 transition-colors" href="/portfolio">
              View Portfolio
            </a>
          </div>
        </div>
      </section>

      <section class="bg-surface-container-low p-6 border border-outline-variant/15">
        <h2 class="font-headline text-sm font-bold uppercase tracking-widest text-primary mb-6">Open Positions</h2>
        {#if dashboardState.status === 'ready' && openPositions().length > 0}
          <div class="space-y-3">
            {#each openPositions() as position}
              <a class="block bg-surface-container p-4 border border-outline-variant/10 hover:border-primary-container/30 transition-colors" href={`/markets/${position.market_id}`}>
                <div class="flex items-start justify-between gap-4">
                  <div>
                    <div class="text-sm font-headline font-bold text-on-surface">{findMarketTitle(position.market_id)}</div>
                    <div class="mt-1 text-[10px] uppercase tracking-[0.2em] text-outline">{formatPositionStatus(position.status)}</div>
                  </div>
                  <span class="px-2 py-1 text-[10px] font-bold uppercase tracking-widest {position.side === 'yes' ? 'bg-primary-container/10 text-primary-container border border-primary-container/20' : 'bg-error/10 text-error border border-error/20'}">
                    {formatPositionSide(position.side)}
                  </span>
                </div>
                <div class="mt-4 flex justify-between text-[10px] uppercase tracking-[0.2em] text-outline">
                  <span>Stake {position.stake_amount}</span>
                  <span>Payout {position.potential_payout}</span>
                </div>
              </a>
            {/each}
          </div>
        {:else if dashboardState.status === 'ready'}
          <p class="text-sm text-outline">You do not have any open positions yet.</p>
        {:else}
          <p class="text-sm text-outline">Position activity will appear here after the dashboard loads.</p>
        {/if}
      </section>
    </aside>
  </section>
</main>
