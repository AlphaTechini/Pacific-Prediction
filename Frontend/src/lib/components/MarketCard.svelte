<script lang="ts">
  import StatusBadge from './StatusBadge.svelte';
  import type { MarketResponse } from '$lib/api-types';

  let { market }: { market: MarketResponse } = $props();

  function formatStatus(status: string): string {
    return status.replaceAll('_', ' ');
  }

  function getBadgeVariant(status: string): 'success' | 'danger' | 'warning' | 'info' {
    if (status === 'active') {
      return 'success';
    }

    if (status === 'resolving') {
      return 'warning';
    }

    if (status === 'resolved') {
      return 'info';
    }

    return 'danger';
  }

  function formatMarketType(value: string): string {
    return value.replaceAll('_', ' ');
  }

  function formatSourceType(value: string): string {
    return value.replaceAll('_', ' ');
  }

  function formatDateTime(value: string): string {
    return new Intl.DateTimeFormat(undefined, {
      month: 'short',
      day: 'numeric',
      hour: 'numeric',
      minute: '2-digit'
    }).format(new Date(value));
  }

  function badgeLabel(): string {
    return formatStatus(market.status);
  }

  function badgeVariant(): 'success' | 'danger' | 'warning' | 'info' {
    return getBadgeVariant(market.status);
  }

  function marketHref(): string {
    return market.status === 'resolved' ? `/markets/${market.id}/resolved` : `/markets/${market.id}`;
  }

  function detailValue(): string {
    return market.threshold_value || market.reference_value || 'Rule-based';
  }

  function footerLabel(): string {
    return market.status === 'resolved' ? 'Resolved' : 'Expires';
  }

  function footerValue(): string {
    return formatDateTime(market.status === 'resolved' ? market.resolved_at ?? market.expiry_time : market.expiry_time);
  }

  function resultValue(): string {
    return market.result ? market.result.toUpperCase() : 'Pending';
  }
</script>

<a class="block bg-surface-container-low border-t-2 border-outline-variant p-5 relative group hover:bg-surface-container transition-colors" href={marketHref()}>
  <div class="flex justify-between items-start mb-4">
    <div class="flex items-center gap-3">
      <div class="w-8 h-8 bg-surface-container-highest flex items-center justify-center border border-outline-variant/30">
        <span class="material-symbols-outlined text-primary-container text-lg">query_stats</span>
      </div>
      <h3 class="font-headline font-bold text-sm tracking-tight leading-tight max-w-[140px]">
        {market.title}<br/>
        <span class="text-[10px] text-outline font-normal uppercase tracking-wider">{market.symbol}</span>
      </h3>
    </div>
    <StatusBadge label={badgeLabel()} variant={badgeVariant()} />
  </div>

  <div class="grid grid-cols-2 gap-2 mb-6">
    <div class="bg-surface-container-highest py-3 px-3 text-xs border border-outline-variant/30">
      <span class="block text-[8px] text-outline mb-1 uppercase tracking-widest">Market Type</span>
      <span class="font-mono uppercase">{formatMarketType(market.market_type)}</span>
    </div>
    <div class="bg-surface-container-highest py-3 px-3 text-xs border border-outline-variant/30">
      <span class="block text-[8px] text-outline mb-1 uppercase tracking-widest">Source</span>
      <span class="font-mono uppercase">{formatSourceType(market.source_type)}</span>
    </div>
  </div>

  <div class="flex justify-between items-center text-[10px] border-t border-outline-variant/10 pt-4">
    <div class="flex items-center gap-1 text-primary-container font-medium">
      <span class="material-symbols-outlined text-xs">{market.status === 'resolved' ? 'verified' : 'timer'}</span>
      {footerLabel()}: {footerValue()}
    </div>
    <div class="text-outline uppercase">{market.condition_operator.replaceAll('_', ' ')}</div>
  </div>

  <div class="mt-4 bg-surface-container-lowest p-2 text-[10px] italic text-outline border-l-2 border-primary-container/30">
    <span class="text-primary-container font-bold not-italic mr-1">{market.status === 'resolved' ? 'RESULT:' : 'RULE:'}</span>
    {#if market.status === 'resolved'}
      {resultValue()}{#if market.settlement_value} at {market.settlement_value}{/if}
    {:else}
      {detailValue()}
    {/if}
  </div>
</a>
