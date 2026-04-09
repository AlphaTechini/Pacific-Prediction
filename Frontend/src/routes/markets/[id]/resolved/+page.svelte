<script lang="ts">
  import { onMount } from 'svelte';
  import { get } from 'svelte/store';
  import { page } from '$app/stores';

  import type { MarketResponse, PositionResponse } from '$lib/api-types';
  import Button from '$lib/components/Button.svelte';
  import TopNavBar from '$lib/components/TopNavBar.svelte';
  import { ensureGuestSession } from '$lib/guest-session';
  import {
    loadMarketAccountData,
    loadMarketDetail,
    loadRelatedMarkets
  } from '$lib/market-detail-data';

  type LoadStatus = 'loading' | 'ready' | 'error';

  const resolvedState = $state({
    loadStatus: 'loading' as LoadStatus,
    loadError: null as string | null,
    market: null as MarketResponse | null,
    positions: [] as PositionResponse[],
    relatedMarkets: [] as MarketResponse[],
    availableBalance: '0.00',
    lockedBalance: '0.00'
  });

  onMount(() => {
    void loadPage();
  });

  async function loadPage(): Promise<void> {
    resolvedState.loadStatus = 'loading';
    resolvedState.loadError = null;

    const marketID = currentMarketID();
    if (!marketID) {
      resolvedState.loadStatus = 'error';
      resolvedState.loadError = 'This market id is missing.';
      return;
    }

    try {
      const [market, relatedMarkets] = await Promise.all([
        loadMarketDetail(marketID),
        loadRelatedMarkets(marketID)
      ]);

      resolvedState.market = market;
      resolvedState.relatedMarkets = relatedMarkets;
      resolvedState.loadStatus = 'ready';
    } catch (error) {
      resolvedState.loadStatus = 'error';
      resolvedState.loadError = toErrorMessage(error, 'Unable to load this resolved market.');
      return;
    }

    try {
      await ensureGuestSession();
      const accountData = await loadMarketAccountData();

      resolvedState.availableBalance = accountData.balance.available_balance;
      resolvedState.lockedBalance = accountData.balance.locked_balance;
      resolvedState.positions = accountData.positions;
    } catch {
      resolvedState.positions = [];
    }
  }

  function currentMarketID(): string {
    return get(page).params.id ?? '';
  }

  function formatDateTime(value?: string): string {
    if (!value) {
      return 'Not available';
    }

    return new Intl.DateTimeFormat(undefined, {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    }).format(new Date(value));
  }

  function marketTypeLabel(value: string): string {
    const labels: Record<string, string> = {
      price_threshold: 'Price Threshold',
      candle_direction: 'Candle Direction',
      funding_threshold: 'Funding Direction Or Threshold'
    };

    return labels[value] ?? value.replaceAll('_', ' ');
  }

  function sourceTypeLabel(value: string): string {
    const labels: Record<string, string> = {
      mark_price: 'Mark Price',
      mark_price_candle: 'Mark Price Candle',
      funding_rate: 'Funding Rate'
    };

    return labels[value] ?? value.replaceAll('_', ' ');
  }

  function operatorLabel(value: string): string {
    const labels: Record<string, string> = {
      gt: 'Greater than',
      gte: 'Greater than or equal',
      lt: 'Less than',
      lte: 'Less than or equal',
      bullish_close: 'Bullish close',
      bearish_close: 'Bearish close',
      positive: 'Positive',
      negative: 'Negative'
    };

    return labels[value] ?? value.replaceAll('_', ' ');
  }

  function formatStatus(value: string): string {
    return value.replaceAll('_', ' ');
  }

  function resultAccent(result?: string): string {
    return result === 'yes' ? 'text-primary-container' : 'text-error';
  }

  function myPositions(): PositionResponse[] {
    const marketID = resolvedState.market?.id;
    if (!marketID) {
      return [];
    }

    return resolvedState.positions.filter((position) => position.market_id === marketID);
  }

  function settlementSummary(): string {
    const market = resolvedState.market;
    if (!market) {
      return '';
    }

    if (market.market_type === 'candle_direction') {
      return `This market resolved from the ${market.source_interval || 'selected'} candle close using a ${operatorLabel(market.condition_operator).toLowerCase()} rule.`;
    }

    if (market.market_type === 'funding_threshold' && ['positive', 'negative'].includes(market.condition_operator)) {
      const direction = market.condition_operator === 'positive' ? 'above zero' : 'below zero';
      return `This market resolved by checking whether the funding value was ${direction} at the funding checkpoint.`;
    }

    if (market.threshold_value) {
      return `This market resolved by checking whether the observed value was ${operatorLabel(market.condition_operator).toLowerCase()} ${market.threshold_value}.`;
    }

    return `This market resolved using a ${operatorLabel(market.condition_operator).toLowerCase()} rule.`;
  }

  function toErrorMessage(error: unknown, fallback: string): string {
    if (error instanceof Error && error.message) {
      return error.message;
    }

    return fallback;
  }
</script>

<svelte:head>
  <title>Market Resolved | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="pt-24 pb-20 px-6 max-w-7xl mx-auto">
  {#if resolvedState.loadStatus === 'loading'}
    <section class="bg-surface-container-low p-8 border border-outline-variant/20">
      <p class="text-sm text-outline">Loading resolved market details.</p>
    </section>
  {:else if resolvedState.loadStatus === 'error'}
    <section class="bg-surface-container-low p-8 border border-error/20 space-y-4">
      <p class="text-sm text-error">{resolvedState.loadError}</p>
      <Button class="px-5 py-3 text-xs uppercase tracking-[0.2em]" onclick={loadPage}>Try Again</Button>
    </section>
  {:else if resolvedState.market}
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <div class="lg:col-span-8 space-y-6">
        <section class="bg-surface-container p-8 border-l-4 border-primary-container">
          <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6">
            <div class="space-y-2">
              <div class="flex items-center gap-3 flex-wrap">
                <span class="bg-surface-container-highest text-primary-container text-[10px] font-bold px-2 py-0.5 tracking-widest uppercase border border-outline-variant/20">
                  {formatStatus(resolvedState.market.status)}
                </span>
                <span class="text-outline text-xs font-mono uppercase">{marketTypeLabel(resolvedState.market.market_type)}</span>
              </div>
              <h1 class="text-3xl md:text-4xl font-headline font-bold text-primary tracking-tight leading-none">
                {resolvedState.market.title}
              </h1>
              <p class="text-outline text-sm font-mono">{resolvedState.market.symbol}</p>
            </div>
            <div class="flex flex-col items-center justify-center bg-surface-container-highest p-6 min-w-[160px] border border-outline-variant/20">
              <span class="text-xs text-outline uppercase tracking-widest mb-1">Final Result</span>
              <span class="text-5xl font-headline font-extrabold {resultAccent(resolvedState.market.result)}">
                {resolvedState.market.result?.toUpperCase() || 'N/A'}
              </span>
            </div>
          </div>
        </section>

        <section class="bg-surface-container-low p-6 border border-outline-variant/15">
          <div class="flex items-center gap-2 mb-4">
            <span class="material-symbols-outlined text-primary-container">verified</span>
            <h2 class="text-sm font-headline font-bold uppercase tracking-widest text-primary">Settlement Summary</h2>
          </div>
          <p class="text-sm text-on-surface-variant leading-relaxed">{settlementSummary()}</p>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mt-6">
            <div class="bg-surface-container p-4 border border-outline-variant/10">
              <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Settlement Value</div>
              <div class="mt-2 text-lg font-mono text-on-surface">{resolvedState.market.settlement_value || 'Not available'}</div>
            </div>
            <div class="bg-surface-container p-4 border border-outline-variant/10">
              <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Resolved At</div>
              <div class="mt-2 text-lg font-mono text-on-surface">{formatDateTime(resolvedState.market.resolved_at)}</div>
            </div>
            <div class="bg-surface-container p-4 border border-outline-variant/10">
              <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Source</div>
              <div class="mt-2 text-lg font-headline text-on-surface">{sourceTypeLabel(resolvedState.market.source_type)}</div>
            </div>
          </div>
        </section>

        <section class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="bg-surface-container-low p-6 border border-outline-variant/15">
            <h2 class="text-sm font-headline font-bold uppercase tracking-widest mb-4 text-primary">Resolution Inputs</h2>
            <div class="space-y-4 text-sm">
              <div class="flex justify-between gap-4">
                <span class="text-outline uppercase tracking-[0.2em] text-[10px]">Rule</span>
                <span class="text-on-surface font-mono text-right">{operatorLabel(resolvedState.market.condition_operator)}</span>
              </div>
              <div class="flex justify-between gap-4">
                <span class="text-outline uppercase tracking-[0.2em] text-[10px]">Threshold</span>
                <span class="text-on-surface font-mono text-right">{resolvedState.market.threshold_value || 'Not required'}</span>
              </div>
              <div class="flex justify-between gap-4">
                <span class="text-outline uppercase tracking-[0.2em] text-[10px]">Interval</span>
                <span class="text-on-surface font-mono text-right">{resolvedState.market.source_interval || 'Not required'}</span>
              </div>
              <div class="flex justify-between gap-4">
                <span class="text-outline uppercase tracking-[0.2em] text-[10px]">Expiry</span>
                <span class="text-on-surface font-mono text-right">{formatDateTime(resolvedState.market.expiry_time)}</span>
              </div>
            </div>
          </div>

          <div class="bg-surface-container-low p-6 border border-outline-variant/15">
            <h2 class="text-sm font-headline font-bold uppercase tracking-widest mb-4 text-primary">Resolution Notes</h2>
            <div class="space-y-4 text-sm">
              <div>
                <div class="text-outline uppercase tracking-[0.2em] text-[10px] mb-1">Reason</div>
                <div class="text-on-surface">{resolvedState.market.resolution_reason || 'No additional reason was provided.'}</div>
              </div>
              <div>
                <div class="text-outline uppercase tracking-[0.2em] text-[10px] mb-1">Created At</div>
                <div class="text-on-surface font-mono">{formatDateTime(resolvedState.market.created_at)}</div>
              </div>
              <div>
                <div class="text-outline uppercase tracking-[0.2em] text-[10px] mb-1">Market Id</div>
                <div class="text-on-surface font-mono break-all">{resolvedState.market.id}</div>
              </div>
            </div>
          </div>
        </section>

        <section class="bg-surface-container-low p-6 border border-outline-variant/15">
          <div class="flex items-center gap-2 mb-4">
            <span class="material-symbols-outlined text-primary-container">account_balance_wallet</span>
            <h2 class="text-sm font-headline font-bold uppercase tracking-widest text-primary">Your Outcome On This Market</h2>
          </div>

          {#if myPositions().length > 0}
            <div class="space-y-3">
              {#each myPositions() as position}
                <div class="bg-surface-container p-4 border border-outline-variant/10">
                  <div class="flex items-start justify-between gap-4">
                    <div>
                      <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Position Status</div>
                      <div class="mt-1 text-sm font-headline font-bold text-on-surface">{formatStatus(position.status)}</div>
                    </div>
                    <span class="px-2 py-1 text-[10px] font-bold uppercase tracking-widest {position.side === 'yes' ? 'bg-primary-container/10 text-primary-container border border-primary-container/20' : 'bg-error/10 text-error border border-error/20'}">
                      {position.side.toUpperCase()}
                    </span>
                  </div>
                  <div class="grid grid-cols-2 gap-4 mt-4 text-xs">
                    <div>
                      <div class="text-outline uppercase tracking-[0.2em]">Stake</div>
                      <div class="mt-1 font-mono text-on-surface">{position.stake_amount}</div>
                    </div>
                    <div>
                      <div class="text-outline uppercase tracking-[0.2em]">Potential Payout</div>
                      <div class="mt-1 font-mono text-on-surface">{position.potential_payout}</div>
                    </div>
                  </div>
                </div>
              {/each}
            </div>
          {:else}
            <p class="text-sm text-outline">You did not place a tracked position on this market.</p>
          {/if}
        </section>
      </div>

      <aside class="lg:col-span-4 space-y-6">
        <section class="bg-surface-container p-6 border border-outline-variant/15">
          <h2 class="text-xs font-headline font-bold uppercase tracking-widest mb-6 border-b border-outline-variant/10 pb-2">Account Snapshot</h2>
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-surface-container-highest p-4">
              <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Available</div>
              <div class="mt-2 text-xl font-headline font-bold text-primary">{resolvedState.availableBalance}</div>
            </div>
            <div class="bg-surface-container-highest p-4">
              <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Locked</div>
              <div class="mt-2 text-xl font-headline font-bold text-on-surface">{resolvedState.lockedBalance}</div>
            </div>
          </div>
          <div class="mt-6 space-y-3">
            <a class="block gradient-primary text-on-primary-fixed py-4 text-xs tracking-widest font-bold text-center" href="/dashboard">
              Back To Markets
            </a>
            <a class="block border border-outline-variant/30 py-4 text-xs tracking-widest font-bold text-center text-primary hover:bg-surface-container-high transition-colors" href="/portfolio">
              View Portfolio
            </a>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="font-headline font-bold uppercase tracking-widest text-xs">Other Active Markets</h2>
          <div class="space-y-3">
            {#if resolvedState.relatedMarkets.length > 0}
              {#each resolvedState.relatedMarkets as relatedMarket}
                <a class="block bg-surface-container-low p-4 border-l-2 border-outline-variant hover:border-primary-container transition-all group" href={`/markets/${relatedMarket.id}`}>
                  <p class="text-xs font-bold text-on-surface mb-1 group-hover:text-primary transition-colors">{relatedMarket.title}</p>
                  <div class="flex justify-between items-center gap-4">
                    <span class="text-[10px] text-outline uppercase font-mono">{relatedMarket.symbol}</span>
                    <span class="text-[10px] text-primary-container uppercase font-mono">{formatStatus(relatedMarket.status)}</span>
                  </div>
                </a>
              {/each}
            {:else}
              <div class="bg-surface-container-low p-4 border border-outline-variant/15 text-sm text-outline">
                No other active markets are available right now.
              </div>
            {/if}
          </div>
        </section>
      </aside>
    </div>
  {/if}
</main>
