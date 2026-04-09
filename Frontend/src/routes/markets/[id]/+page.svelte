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
    loadRelatedMarkets,
    submitPosition
  } from '$lib/market-detail-data';

  type LoadStatus = 'loading' | 'ready' | 'error';
  type SubmitStatus = 'idle' | 'submitting' | 'error';

  const detailState = $state({
    loadStatus: 'loading' as LoadStatus,
    submitStatus: 'idle' as SubmitStatus,
    loadError: null as string | null,
    submitError: null as string | null,
    successMessage: null as string | null,
    market: null as MarketResponse | null,
    relatedMarkets: [] as MarketResponse[],
    positions: [] as PositionResponse[],
    availableBalance: '0.00',
    lockedBalance: '0.00'
  });

  let selectedSide = $state<'yes' | 'no'>('yes');
  let stakeAmount = $state('');

  onMount(() => {
    void loadPage();
  });

  async function loadPage(): Promise<void> {
    detailState.loadStatus = 'loading';
    detailState.loadError = null;
    detailState.successMessage = null;

    const marketID = currentMarketID();
    if (!marketID) {
      detailState.loadStatus = 'error';
      detailState.loadError = 'This market id is missing.';
      return;
    }

    try {
      const [market, relatedMarkets] = await Promise.all([
        loadMarketDetail(marketID),
        loadRelatedMarkets(marketID)
      ]);

      detailState.market = market;
      detailState.relatedMarkets = relatedMarkets;
      detailState.loadStatus = 'ready';
    } catch (error) {
      detailState.loadStatus = 'error';
      detailState.loadError = toErrorMessage(error, 'Unable to load this market.');
      return;
    }

    try {
      await ensureGuestSession();
      const accountData = await loadMarketAccountData();

      detailState.availableBalance = accountData.balance.available_balance;
      detailState.lockedBalance = accountData.balance.locked_balance;
      detailState.positions = accountData.positions;
    } catch (error) {
      detailState.submitError = toErrorMessage(error, 'Your guest session is not ready yet.');
    }
  }

  function currentMarketID(): string {
    return get(page).params.id ?? '';
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

  function statusVariant(value: string): string {
    if (value === 'active') {
      return 'bg-primary-container/10 text-primary-container border-primary-container/20';
    }

    if (value === 'resolving') {
      return 'bg-tertiary-fixed-dim/10 text-tertiary-fixed-dim border-tertiary-fixed-dim/20';
    }

    if (value === 'resolved') {
      return 'bg-surface-container-highest text-outline border-outline-variant/30';
    }

    return 'bg-error/10 text-error border-error/20';
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

  function settlementRule(): string {
    const market = detailState.market;
    if (!market) {
      return '';
    }

    if (market.market_type === 'candle_direction') {
      return `This market resolves from the ${market.source_interval || 'selected'} candle close using a ${operatorLabel(market.condition_operator).toLowerCase()} rule.`;
    }

    if (market.market_type === 'funding_threshold' && ['positive', 'negative'].includes(market.condition_operator)) {
      const direction = market.condition_operator === 'positive' ? 'above zero' : 'below zero';
      return `This market settles YES when the funding value is ${direction} at the funding checkpoint.`;
    }

    if (market.threshold_value) {
      return `This market settles YES when the observed value is ${operatorLabel(market.condition_operator).toLowerCase()} ${market.threshold_value}.`;
    }

    return `This market uses ${operatorLabel(market.condition_operator).toLowerCase()} as its settlement rule.`;
  }

  function myPositions(): PositionResponse[] {
    const marketID = detailState.market?.id;
    if (!marketID) {
      return [];
    }

    return detailState.positions.filter((position) => position.market_id === marketID);
  }

  function tradingClosed(): boolean {
    return detailState.market?.status !== 'active';
  }

  async function handleSubmit(event: SubmitEvent): Promise<void> {
    event.preventDefault();

    const market = detailState.market;
    if (!market) {
      return;
    }

    if (tradingClosed()) {
      detailState.submitStatus = 'error';
      detailState.submitError = 'Trading is closed for this market.';
      return;
    }

    detailState.submitStatus = 'submitting';
    detailState.submitError = null;
    detailState.successMessage = null;

    const player = await ensureGuestSession();
    if (!player) {
      detailState.submitStatus = 'error';
      detailState.submitError = 'Your guest session is not ready yet.';
      return;
    }

    try {
      await submitPosition(market.id, {
        side: selectedSide,
        stake_amount: stakeAmount.trim()
      });

      stakeAmount = '';
      detailState.submitStatus = 'idle';
      detailState.successMessage = 'Your position was placed successfully.';

      const accountData = await loadMarketAccountData();
      detailState.availableBalance = accountData.balance.available_balance;
      detailState.lockedBalance = accountData.balance.locked_balance;
      detailState.positions = accountData.positions;
    } catch (error) {
      detailState.submitStatus = 'error';
      detailState.submitError = toErrorMessage(error, 'Unable to place that position right now.');
    }
  }

  function toErrorMessage(error: unknown, fallback: string): string {
    if (error instanceof Error && error.message) {
      return error.message;
    }

    return fallback;
  }
</script>

<svelte:head>
  <title>Market Details | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="pt-24 pb-20 px-6 max-w-[1400px] mx-auto">
  {#if detailState.loadStatus === 'loading'}
    <section class="bg-surface-container-low p-8 border border-outline-variant/20">
      <p class="text-sm text-outline">Loading this market.</p>
    </section>
  {:else if detailState.loadStatus === 'error'}
    <section class="bg-surface-container-low p-8 border border-error/20 space-y-4">
      <p class="text-sm text-error">{detailState.loadError}</p>
      <Button class="px-5 py-3 text-xs uppercase tracking-[0.2em]" onclick={loadPage}>Try Again</Button>
    </section>
  {:else if detailState.market}
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <div class="lg:col-span-8 space-y-6">
        <header class="bg-surface-container-low p-6 rounded-sm border-l-4 border-primary-container">
          <div class="flex flex-wrap justify-between items-start gap-4 mb-4">
            <div class="space-y-2">
              <div class="flex items-center gap-3 flex-wrap">
                <span class="px-2 py-0.5 text-[10px] font-bold tracking-widest uppercase rounded-sm border {statusVariant(detailState.market.status)}">
                  {formatStatus(detailState.market.status)}
                </span>
                <span class="text-outline text-xs font-mono uppercase tracking-widest">
                  Type: {marketTypeLabel(detailState.market.market_type)}
                </span>
              </div>
              <h1 class="text-2xl md:text-3xl font-bold font-headline leading-tight text-on-surface">
                {detailState.market.title}
              </h1>
              <p class="text-primary font-mono text-sm tracking-tight opacity-80">{detailState.market.symbol}</p>
            </div>
            <div class="bg-surface-container p-4 rounded-sm border border-outline-variant/15 min-w-[220px]">
              <div class="text-[10px] text-outline uppercase tracking-widest font-bold">Expiry Time</div>
              <div class="mt-2 text-base font-mono font-bold text-primary-container">{formatDateTime(detailState.market.expiry_time)}</div>
              {#if detailState.market.resolved_at}
                <div class="mt-3 text-[10px] text-outline uppercase tracking-widest font-bold">Resolved At</div>
                <div class="mt-1 text-sm font-mono text-on-surface">{formatDateTime(detailState.market.resolved_at)}</div>
              {/if}
            </div>
          </div>
        </header>

        <section class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="bg-surface-container-low p-4 border border-outline-variant/15">
            <p class="text-[10px] text-outline uppercase tracking-widest mb-2">Source</p>
            <p class="text-sm font-headline font-bold text-on-surface">{sourceTypeLabel(detailState.market.source_type)}</p>
          </div>
          <div class="bg-surface-container-low p-4 border border-outline-variant/15">
            <p class="text-[10px] text-outline uppercase tracking-widest mb-2">Rule</p>
            <p class="text-sm font-headline font-bold text-on-surface">{operatorLabel(detailState.market.condition_operator)}</p>
          </div>
          <div class="bg-surface-container-low p-4 border border-outline-variant/15">
            <p class="text-[10px] text-outline uppercase tracking-widest mb-2">Threshold</p>
            <p class="text-sm font-mono font-bold text-on-surface">{detailState.market.threshold_value || 'Not required'}</p>
          </div>
          <div class="bg-surface-container-low p-4 border border-outline-variant/15">
            <p class="text-[10px] text-outline uppercase tracking-widest mb-2">Interval</p>
            <p class="text-sm font-mono font-bold text-on-surface">{detailState.market.source_interval || 'Not required'}</p>
          </div>
        </section>

        <section class="bg-surface-container p-6 rounded-sm space-y-4">
          <div class="flex items-center gap-2 border-b border-outline-variant/15 pb-3">
            <span class="material-symbols-outlined text-primary-container text-lg">rule</span>
            <h2 class="font-headline font-bold uppercase tracking-wider text-sm">Settlement Rule</h2>
          </div>
          <p class="text-on-surface-variant text-sm leading-relaxed">{settlementRule()}</p>
          {#if detailState.market.result}
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 pt-2">
              <div>
                <p class="text-[10px] text-outline uppercase tracking-widest font-bold">Result</p>
                <p class="mt-1 text-sm font-mono text-on-surface">{detailState.market.result.toUpperCase()}</p>
              </div>
              <div>
                <p class="text-[10px] text-outline uppercase tracking-widest font-bold">Settlement Value</p>
                <p class="mt-1 text-sm font-mono text-on-surface">{detailState.market.settlement_value || 'Not available'}</p>
              </div>
              <div>
                <p class="text-[10px] text-outline uppercase tracking-widest font-bold">Resolution Reason</p>
                <p class="mt-1 text-sm text-on-surface">{detailState.market.resolution_reason || 'Not available'}</p>
              </div>
            </div>
          {/if}
        </section>

        <section class="bg-surface-container-low p-6 rounded-sm border border-outline-variant/15">
          <div class="flex items-center gap-2 mb-6">
            <span class="material-symbols-outlined text-primary-container">account_balance_wallet</span>
            <h2 class="font-headline text-sm font-bold uppercase tracking-widest text-primary">Your Positions On This Market</h2>
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
            <p class="text-sm text-outline">You have not placed a position on this market yet.</p>
          {/if}
        </section>
      </div>

      <div class="lg:col-span-4 space-y-6">
        <section class="bg-surface-container-high p-6 rounded-sm border border-outline-variant/20 shadow-[0_0_12px_rgba(0,219,233,0.3)] sticky top-24">
          <h2 class="text-xs font-bold uppercase tracking-widest text-outline mb-6">
            {tradingClosed() ? 'Market Status' : 'Place Position'}
          </h2>

          {#if tradingClosed()}
            <div class="space-y-4">
              <p class="text-sm text-on-surface-variant leading-relaxed">Trading is closed for this market because it is currently {formatStatus(detailState.market.status)}.</p>
              {#if detailState.market.result}
                <div class="bg-surface-container p-4 rounded-sm border border-outline-variant/15">
                  <div class="text-[10px] text-outline uppercase tracking-[0.2em]">Final Result</div>
                  <div class="mt-2 text-2xl font-headline font-bold text-primary-container">{detailState.market.result.toUpperCase()}</div>
                </div>
              {/if}
            </div>
          {:else}
            <form class="space-y-4" onsubmit={handleSubmit}>
              <div class="grid grid-cols-2 gap-3">
                {#each ['yes', 'no'] as side}
                  <button
                    class="flex flex-col items-center justify-center p-4 border transition-all rounded-sm {selectedSide === side ? side === 'yes' ? 'bg-primary-container/20 border-primary-container' : 'bg-error/20 border-error' : side === 'yes' ? 'bg-surface-container border-primary-container/30 hover:bg-primary-container/10' : 'bg-surface-container border-error/30 hover:bg-error/10'}"
                    onclick={() => selectedSide = side as 'yes' | 'no'}
                    type="button"
                  >
                    <span class="{side === 'yes' ? 'text-primary-container' : 'text-error'} font-headline font-bold text-lg uppercase">{side}</span>
                  </button>
                {/each}
              </div>

              <label class="block space-y-2">
                <div class="flex justify-between items-center px-1">
                  <span class="text-[10px] text-outline uppercase tracking-widest font-bold">Stake Amount</span>
                  <span class="text-[10px] text-outline uppercase font-mono">Available: {detailState.availableBalance}</span>
                </div>
                <div class="relative group">
                  <input
                    bind:value={stakeAmount}
                    class="w-full bg-surface-container-lowest border border-outline-variant/20 text-xl font-mono p-4 pr-16 rounded-sm focus:ring-1 focus:ring-primary-container/30 focus:outline-none placeholder:text-outline-variant"
                    placeholder="0.00"
                    required
                    step="0.00000001"
                    type="number"
                  />
                  <span class="absolute right-4 top-1/2 -translate-y-1/2 text-sm font-bold text-outline">USDC</span>
                </div>
              </label>

              <div class="bg-surface-container p-4 rounded-sm space-y-3">
                <div class="flex justify-between items-center text-xs">
                  <span class="text-outline">Selected Side</span>
                  <span class="{selectedSide === 'yes' ? 'text-primary-container' : 'text-error'} font-mono font-bold uppercase">{selectedSide}</span>
                </div>
                <div class="flex justify-between items-center text-xs">
                  <span class="text-outline">Locked Balance</span>
                  <span class="text-on-surface font-mono">{detailState.lockedBalance}</span>
                </div>
              </div>

              {#if detailState.submitError}
                <div class="bg-surface-container-low border border-error/20 p-4 text-sm text-error">
                  {detailState.submitError}
                </div>
              {/if}

              {#if detailState.successMessage}
                <div class="bg-surface-container-low border border-primary-container/20 p-4 text-sm text-primary">
                  {detailState.successMessage}
                </div>
              {/if}

              <button
                class="w-full gradient-primary text-on-primary-fixed py-4 font-headline font-bold uppercase tracking-[0.2em] rounded-sm hover:opacity-90 active:scale-[0.99] transition-all disabled:opacity-70"
                disabled={detailState.submitStatus === 'submitting'}
                type="submit"
              >
                {detailState.submitStatus === 'submitting' ? 'Placing Position...' : 'Place Position'}
              </button>
            </form>
          {/if}
        </section>

        <section class="space-y-4">
          <h2 class="font-headline font-bold uppercase tracking-widest text-xs">Other Active Markets</h2>
          <div class="space-y-3">
            {#if detailState.relatedMarkets.length > 0}
              {#each detailState.relatedMarkets as relatedMarket}
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
      </div>
    </div>
  {/if}
</main>
