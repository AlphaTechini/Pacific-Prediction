<script lang="ts">
	import { resolve } from '$app/paths';
	import StatusBadge from './StatusBadge.svelte';
	import type { MarketResponse } from '$lib/api-types';
	import { formatMarketTitleDisplay } from '$lib/market-title';

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

	function detailValue(): string {
		return market.threshold_value || market.reference_value || 'Rule-based';
	}

	function footerLabel(): string {
		if (market.status === 'resolved') {
			return 'Resolved';
		}

		if (market.market_type === 'funding_threshold') {
			return 'Funding Checkpoint';
		}

		if (market.market_type === 'candle_direction') {
			return 'Candle Close';
		}

		return 'Expires';
	}

	function footerValue(): string {
		return formatDateTime(
			market.status === 'resolved' ? (market.resolved_at ?? market.expiry_time) : market.expiry_time
		);
	}

	function resultValue(): string {
		return market.result ? market.result.toUpperCase() : 'Pending';
	}
</script>

<a
	class="bg-surface-container-low border-outline-variant group hover:bg-surface-container relative block border-t-2 p-5 transition-colors"
	href={resolve(market.status === 'resolved' ? `/markets/${market.id}/resolved` : `/markets/${market.id}`)}
>
	<div class="mb-4 flex items-start justify-between">
		<div class="flex items-center gap-3">
			<div
				class="bg-surface-container-highest border-outline-variant/30 flex h-8 w-8 items-center justify-center border"
			>
				<span class="material-symbols-outlined text-primary-container text-lg">query_stats</span>
			</div>
			<h3 class="font-headline max-w-35 text-sm leading-tight font-bold tracking-tight">
				{formatMarketTitleDisplay(market.title)}<br />
				<span class="text-outline text-[10px] font-normal tracking-wider uppercase"
					>{market.symbol}</span
				>
			</h3>
		</div>
		<StatusBadge label={badgeLabel()} variant={badgeVariant()} />
	</div>

	<div class="mb-6 grid grid-cols-2 gap-2">
		<div class="bg-surface-container-highest border-outline-variant/30 border px-3 py-3 text-xs">
			<span class="text-outline mb-1 block text-[8px] tracking-widest uppercase">Market Type</span>
			<span class="font-mono uppercase">{formatMarketType(market.market_type)}</span>
		</div>
		<div class="bg-surface-container-highest border-outline-variant/30 border px-3 py-3 text-xs">
			<span class="text-outline mb-1 block text-[8px] tracking-widest uppercase">Source</span>
			<span class="font-mono uppercase">{formatSourceType(market.source_type)}</span>
		</div>
	</div>

	<div
		class="border-outline-variant/10 flex items-center justify-between border-t pt-4 text-[10px]"
	>
		<div class="text-primary-container flex items-center gap-1 font-medium">
			<span class="material-symbols-outlined text-xs"
				>{market.status === 'resolved' ? 'verified' : 'timer'}</span
			>
			{footerLabel()}: {footerValue()}
		</div>
		<div class="text-outline uppercase">{market.condition_operator.replaceAll('_', ' ')}</div>
	</div>

	<div
		class="bg-surface-container-lowest text-outline border-primary-container/30 mt-4 border-l-2 p-2 text-[10px] italic"
	>
		<span class="text-primary-container mr-1 font-bold not-italic"
			>{market.status === 'resolved' ? 'RESULT:' : 'RULE:'}</span
		>
		{#if market.status === 'resolved'}
			{resultValue()}{#if market.settlement_value}
				&nbsp;at {market.settlement_value}{/if}
		{:else}
			{detailValue()}
		{/if}
	</div>
</a>
