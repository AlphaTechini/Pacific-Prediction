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

		if (
			market.market_type === 'funding_threshold' &&
			['positive', 'negative'].includes(market.condition_operator)
		) {
			const direction = market.condition_operator === 'positive' ? 'above zero' : 'below zero';
			return `This market resolved by checking whether the funding value was ${direction} at the funding checkpoint.`;
		}

		if (market.threshold_value) {
			return `This market resolved by checking whether the observed value was ${operatorLabel(market.condition_operator).toLowerCase()} ${market.threshold_value}.`;
		}

		return `This market resolved using a ${operatorLabel(market.condition_operator).toLowerCase()} rule.`;
	}

	function timingLabel(): string {
		const market = resolvedState.market;
		if (!market) {
			return 'Expiry';
		}

		if (market.market_type === 'funding_threshold') {
			return 'Funding Checkpoint';
		}

		if (market.market_type === 'candle_direction') {
			return 'Candle Close';
		}

		return 'Expiry';
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

<main class="mx-auto max-w-7xl px-6 pt-24 pb-20">
	{#if resolvedState.loadStatus === 'loading'}
		<section class="bg-surface-container-low border-outline-variant/20 border p-8">
			<p class="text-outline text-sm">Loading resolved market details.</p>
		</section>
	{:else if resolvedState.loadStatus === 'error'}
		<section class="bg-surface-container-low border-error/20 space-y-4 border p-8">
			<p class="text-error text-sm">{resolvedState.loadError}</p>
			<Button class="px-5 py-3 text-xs tracking-[0.2em] uppercase" onclick={loadPage}
				>Try Again</Button
			>
		</section>
	{:else if resolvedState.market}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
			<div class="space-y-6 lg:col-span-8">
				<section class="bg-surface-container border-primary-container border-l-4 p-8">
					<div class="flex flex-col items-start justify-between gap-6 md:flex-row md:items-center">
						<div class="space-y-2">
							<div class="flex flex-wrap items-center gap-3">
								<span
									class="bg-surface-container-highest text-primary-container border-outline-variant/20 border px-2 py-0.5 text-[10px] font-bold tracking-widest uppercase"
								>
									{formatStatus(resolvedState.market.status)}
								</span>
								<span class="text-outline font-mono text-xs uppercase"
									>{marketTypeLabel(resolvedState.market.market_type)}</span
								>
							</div>
							<h1
								class="font-headline text-primary text-3xl leading-none font-bold tracking-tight md:text-4xl"
							>
								{resolvedState.market.title}
							</h1>
							<p class="text-outline font-mono text-sm">{resolvedState.market.symbol}</p>
						</div>
						<div
							class="bg-surface-container-highest border-outline-variant/20 flex min-w-[160px] flex-col items-center justify-center border p-6"
						>
							<span class="text-outline mb-1 text-xs tracking-widest uppercase">Final Result</span>
							<span
								class="font-headline text-5xl font-extrabold {resultAccent(
									resolvedState.market.result
								)}"
							>
								{resolvedState.market.result?.toUpperCase() || 'N/A'}
							</span>
						</div>
					</div>
				</section>

				<section class="bg-surface-container-low border-outline-variant/15 border p-6">
					<div class="mb-4 flex items-center gap-2">
						<span class="material-symbols-outlined text-primary-container">verified</span>
						<h2 class="font-headline text-primary text-sm font-bold tracking-widest uppercase">
							Settlement Summary
						</h2>
					</div>
					<p class="text-on-surface-variant text-sm leading-relaxed">{settlementSummary()}</p>
					<div class="mt-6 grid grid-cols-1 gap-4 md:grid-cols-3">
						<div class="bg-surface-container border-outline-variant/10 border p-4">
							<div class="text-outline text-[10px] tracking-[0.2em] uppercase">
								Settlement Value
							</div>
							<div class="text-on-surface mt-2 font-mono text-lg">
								{resolvedState.market.settlement_value || 'Not available'}
							</div>
						</div>
						<div class="bg-surface-container border-outline-variant/10 border p-4">
							<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Resolved At</div>
							<div class="text-on-surface mt-2 font-mono text-lg">
								{formatDateTime(resolvedState.market.resolved_at)}
							</div>
						</div>
						<div class="bg-surface-container border-outline-variant/10 border p-4">
							<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Source</div>
							<div class="font-headline text-on-surface mt-2 text-lg">
								{sourceTypeLabel(resolvedState.market.source_type)}
							</div>
						</div>
					</div>
				</section>

				<section class="grid grid-cols-1 gap-6 md:grid-cols-2">
					<div class="bg-surface-container-low border-outline-variant/15 border p-6">
						<h2 class="font-headline text-primary mb-4 text-sm font-bold tracking-widest uppercase">
							Resolution Inputs
						</h2>
						<div class="space-y-4 text-sm">
							<div class="flex justify-between gap-4">
								<span class="text-outline text-[10px] tracking-[0.2em] uppercase">Rule</span>
								<span class="text-on-surface text-right font-mono"
									>{operatorLabel(resolvedState.market.condition_operator)}</span
								>
							</div>
							<div class="flex justify-between gap-4">
								<span class="text-outline text-[10px] tracking-[0.2em] uppercase">Threshold</span>
								<span class="text-on-surface text-right font-mono"
									>{resolvedState.market.threshold_value || 'Not required'}</span
								>
							</div>
							<div class="flex justify-between gap-4">
								<span class="text-outline text-[10px] tracking-[0.2em] uppercase">Interval</span>
								<span class="text-on-surface text-right font-mono"
									>{resolvedState.market.source_interval || 'Not required'}</span
								>
							</div>
							<div class="flex justify-between gap-4">
								<span class="text-outline text-[10px] tracking-[0.2em] uppercase"
									>{timingLabel()}</span
								>
								<span class="text-on-surface text-right font-mono"
									>{formatDateTime(resolvedState.market.expiry_time)}</span
								>
							</div>
						</div>
					</div>

					<div class="bg-surface-container-low border-outline-variant/15 border p-6">
						<h2 class="font-headline text-primary mb-4 text-sm font-bold tracking-widest uppercase">
							Resolution Notes
						</h2>
						<div class="space-y-4 text-sm">
							<div>
								<div class="text-outline mb-1 text-[10px] tracking-[0.2em] uppercase">Reason</div>
								<div class="text-on-surface">
									{resolvedState.market.resolution_reason || 'No additional reason was provided.'}
								</div>
							</div>
							<div>
								<div class="text-outline mb-1 text-[10px] tracking-[0.2em] uppercase">
									Created At
								</div>
								<div class="text-on-surface font-mono">
									{formatDateTime(resolvedState.market.created_at)}
								</div>
							</div>
							<div>
								<div class="text-outline mb-1 text-[10px] tracking-[0.2em] uppercase">
									Market Id
								</div>
								<div class="text-on-surface font-mono break-all">{resolvedState.market.id}</div>
							</div>
						</div>
					</div>
				</section>

				<section class="bg-surface-container-low border-outline-variant/15 border p-6">
					<div class="mb-4 flex items-center gap-2">
						<span class="material-symbols-outlined text-primary-container"
							>account_balance_wallet</span
						>
						<h2 class="font-headline text-primary text-sm font-bold tracking-widest uppercase">
							Your Outcome On This Market
						</h2>
					</div>

					{#if myPositions().length > 0}
						<div class="space-y-3">
							{#each myPositions() as position}
								<div class="bg-surface-container border-outline-variant/10 border p-4">
									<div class="flex items-start justify-between gap-4">
										<div>
											<div class="text-outline text-[10px] tracking-[0.2em] uppercase">
												Position Status
											</div>
											<div class="font-headline text-on-surface mt-1 text-sm font-bold">
												{formatStatus(position.status)}
											</div>
										</div>
										<span
											class="px-2 py-1 text-[10px] font-bold tracking-widest uppercase {position.side ===
											'yes'
												? 'bg-primary-container/10 text-primary-container border-primary-container/20 border'
												: 'bg-error/10 text-error border-error/20 border'}"
										>
											{position.side.toUpperCase()}
										</span>
									</div>
									<div class="mt-4 grid grid-cols-2 gap-4 text-xs">
										<div>
											<div class="text-outline tracking-[0.2em] uppercase">Stake</div>
											<div class="text-on-surface mt-1 font-mono">{position.stake_amount}</div>
										</div>
										<div>
											<div class="text-outline tracking-[0.2em] uppercase">Potential Payout</div>
											<div class="text-on-surface mt-1 font-mono">{position.potential_payout}</div>
										</div>
									</div>
								</div>
							{/each}
						</div>
					{:else}
						<p class="text-outline text-sm">You did not place a tracked position on this market.</p>
					{/if}
				</section>
			</div>

			<aside class="space-y-6 lg:col-span-4">
				<section class="bg-surface-container border-outline-variant/15 border p-6">
					<h2
						class="font-headline border-outline-variant/10 mb-6 border-b pb-2 text-xs font-bold tracking-widest uppercase"
					>
						Account Snapshot
					</h2>
					<div class="grid grid-cols-2 gap-4">
						<div class="bg-surface-container-highest p-4">
							<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Available</div>
							<div class="font-headline text-primary mt-2 text-xl font-bold">
								{resolvedState.availableBalance}
							</div>
						</div>
						<div class="bg-surface-container-highest p-4">
							<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Locked</div>
							<div class="font-headline text-on-surface mt-2 text-xl font-bold">
								{resolvedState.lockedBalance}
							</div>
						</div>
					</div>
					<div class="mt-6 space-y-3">
						<a
							class="gradient-primary text-on-primary-fixed block py-4 text-center text-xs font-bold tracking-widest"
							href="/dashboard"
						>
							Back To Markets
						</a>
						<a
							class="border-outline-variant/30 text-primary hover:bg-surface-container-high block border py-4 text-center text-xs font-bold tracking-widest transition-colors"
							href="/portfolio"
						>
							View Portfolio
						</a>
					</div>
				</section>

				<section class="space-y-4">
					<h2 class="font-headline text-xs font-bold tracking-widest uppercase">
						Other Active Markets
					</h2>
					<div class="space-y-3">
						{#if resolvedState.relatedMarkets.length > 0}
							{#each resolvedState.relatedMarkets as relatedMarket}
								<a
									class="bg-surface-container-low border-outline-variant hover:border-primary-container group block border-l-2 p-4 transition-all"
									href={`/markets/${relatedMarket.id}`}
								>
									<p
										class="text-on-surface group-hover:text-primary mb-1 text-xs font-bold transition-colors"
									>
										{relatedMarket.title}
									</p>
									<div class="flex items-center justify-between gap-4">
										<span class="text-outline font-mono text-[10px] uppercase"
											>{relatedMarket.symbol}</span
										>
										<span class="text-primary-container font-mono text-[10px] uppercase"
											>{formatStatus(relatedMarket.status)}</span
										>
									</div>
								</a>
							{/each}
						{:else}
							<div
								class="bg-surface-container-low border-outline-variant/15 text-outline border p-4 text-sm"
							>
								No other active markets are available right now.
							</div>
						{/if}
					</div>
				</section>
			</aside>
		</div>
	{/if}
</main>
