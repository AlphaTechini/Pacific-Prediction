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

	function normalizeText(value: string | number | null | undefined): string {
		if (value === null || value === undefined) {
			return '';
		}

		return String(value).trim();
	}

	function settlementRule(): string {
		const market = detailState.market;
		if (!market) {
			return '';
		}

		if (market.market_type === 'candle_direction') {
			return `This market resolves from the ${market.source_interval || 'selected'} candle close using a ${operatorLabel(market.condition_operator).toLowerCase()} rule.`;
		}

		if (
			market.market_type === 'funding_threshold' &&
			['positive', 'negative'].includes(market.condition_operator)
		) {
			const direction = market.condition_operator === 'positive' ? 'above zero' : 'below zero';
			return `This market settles YES when the funding value is ${direction} at the funding checkpoint.`;
		}

		if (market.threshold_value) {
			return `This market settles YES when the observed value is ${operatorLabel(market.condition_operator).toLowerCase()} ${market.threshold_value}.`;
		}

		return `This market uses ${operatorLabel(market.condition_operator).toLowerCase()} as its settlement rule.`;
	}

	function timingLabel(): string {
		const market = detailState.market;
		if (!market) {
			return 'Expiry Time';
		}

		if (market.market_type === 'funding_threshold') {
			return 'Funding Checkpoint';
		}

		if (market.market_type === 'candle_direction') {
			return 'Candle Close';
		}

		return 'Expiry Time';
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
				stake_amount: normalizeText(stakeAmount)
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

<main class="mx-auto max-w-[1400px] px-6 pt-24 pb-20">
	{#if detailState.loadStatus === 'loading'}
		<section class="bg-surface-container-low border-outline-variant/20 border p-8">
			<p class="text-outline text-sm">Loading this market.</p>
		</section>
	{:else if detailState.loadStatus === 'error'}
		<section class="bg-surface-container-low border-error/20 space-y-4 border p-8">
			<p class="text-error text-sm">{detailState.loadError}</p>
			<Button class="px-5 py-3 text-xs tracking-[0.2em] uppercase" onclick={loadPage}
				>Try Again</Button
			>
		</section>
	{:else if detailState.market}
		<div class="grid grid-cols-1 gap-6 lg:grid-cols-12">
			<div class="space-y-6 lg:col-span-8">
				<header class="bg-surface-container-low border-primary-container rounded-sm border-l-4 p-6">
					<div class="mb-4 flex flex-wrap items-start justify-between gap-4">
						<div class="space-y-2">
							<div class="flex flex-wrap items-center gap-3">
								<span
									class="rounded-sm border px-2 py-0.5 text-[10px] font-bold tracking-widest uppercase {statusVariant(
										detailState.market.status
									)}"
								>
									{formatStatus(detailState.market.status)}
								</span>
								<span class="text-outline font-mono text-xs tracking-widest uppercase">
									Type: {marketTypeLabel(detailState.market.market_type)}
								</span>
							</div>
							<h1
								class="font-headline text-on-surface text-2xl leading-tight font-bold md:text-3xl"
							>
								{detailState.market.title}
							</h1>
							<p class="text-primary font-mono text-sm tracking-tight opacity-80">
								{detailState.market.symbol}
							</p>
						</div>
						<div
							class="bg-surface-container border-outline-variant/15 min-w-[220px] rounded-sm border p-4"
						>
							<div class="text-outline text-[10px] font-bold tracking-widest uppercase">
								{timingLabel()}
							</div>
							<div class="text-primary-container mt-2 font-mono text-base font-bold">
								{formatDateTime(detailState.market.expiry_time)}
							</div>
							{#if detailState.market.resolved_at}
								<div class="text-outline mt-3 text-[10px] font-bold tracking-widest uppercase">
									Resolved At
								</div>
								<div class="text-on-surface mt-1 font-mono text-sm">
									{formatDateTime(detailState.market.resolved_at)}
								</div>
							{/if}
						</div>
					</div>
				</header>

				<section class="grid grid-cols-2 gap-4 md:grid-cols-4">
					<div class="bg-surface-container-low border-outline-variant/15 border p-4">
						<p class="text-outline mb-2 text-[10px] tracking-widest uppercase">Source</p>
						<p class="font-headline text-on-surface text-sm font-bold">
							{sourceTypeLabel(detailState.market.source_type)}
						</p>
					</div>
					<div class="bg-surface-container-low border-outline-variant/15 border p-4">
						<p class="text-outline mb-2 text-[10px] tracking-widest uppercase">Rule</p>
						<p class="font-headline text-on-surface text-sm font-bold">
							{operatorLabel(detailState.market.condition_operator)}
						</p>
					</div>
					<div class="bg-surface-container-low border-outline-variant/15 border p-4">
						<p class="text-outline mb-2 text-[10px] tracking-widest uppercase">Threshold</p>
						<p class="text-on-surface font-mono text-sm font-bold">
							{detailState.market.threshold_value || 'Not required'}
						</p>
					</div>
					<div class="bg-surface-container-low border-outline-variant/15 border p-4">
						<p class="text-outline mb-2 text-[10px] tracking-widest uppercase">Interval</p>
						<p class="text-on-surface font-mono text-sm font-bold">
							{detailState.market.source_interval || 'Not required'}
						</p>
					</div>
				</section>

				<section class="bg-surface-container space-y-4 rounded-sm p-6">
					<div class="border-outline-variant/15 flex items-center gap-2 border-b pb-3">
						<span class="material-symbols-outlined text-primary-container text-lg">rule</span>
						<h2 class="font-headline text-sm font-bold tracking-wider uppercase">
							Settlement Rule
						</h2>
					</div>
					<p class="text-on-surface-variant text-sm leading-relaxed">{settlementRule()}</p>
					{#if detailState.market.result}
						<div class="grid grid-cols-1 gap-4 pt-2 md:grid-cols-3">
							<div>
								<p class="text-outline text-[10px] font-bold tracking-widest uppercase">Result</p>
								<p class="text-on-surface mt-1 font-mono text-sm">
									{detailState.market.result.toUpperCase()}
								</p>
							</div>
							<div>
								<p class="text-outline text-[10px] font-bold tracking-widest uppercase">
									Settlement Value
								</p>
								<p class="text-on-surface mt-1 font-mono text-sm">
									{detailState.market.settlement_value || 'Not available'}
								</p>
							</div>
							<div>
								<p class="text-outline text-[10px] font-bold tracking-widest uppercase">
									Resolution Reason
								</p>
								<p class="text-on-surface mt-1 text-sm">
									{detailState.market.resolution_reason || 'Not available'}
								</p>
							</div>
						</div>
					{/if}
				</section>

				<section class="bg-surface-container-low border-outline-variant/15 rounded-sm border p-6">
					<div class="mb-6 flex items-center gap-2">
						<span class="material-symbols-outlined text-primary-container"
							>account_balance_wallet</span
						>
						<h2 class="font-headline text-primary text-sm font-bold tracking-widest uppercase">
							Your Positions On This Market
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
						<p class="text-outline text-sm">You have not placed a position on this market yet.</p>
					{/if}
				</section>
			</div>

			<div class="space-y-6 lg:col-span-4">
				<section
					class="bg-surface-container-high border-outline-variant/20 sticky top-24 rounded-sm border p-6 shadow-[0_0_12px_rgba(0,219,233,0.3)]"
				>
					<h2 class="text-outline mb-6 text-xs font-bold tracking-widest uppercase">
						{tradingClosed() ? 'Market Status' : 'Place Position'}
					</h2>

					{#if tradingClosed()}
						<div class="space-y-4">
							<p class="text-on-surface-variant text-sm leading-relaxed">
								Trading is closed for this market because it is currently {formatStatus(
									detailState.market.status
								)}.
							</p>
							{#if detailState.market.result}
								<div class="bg-surface-container border-outline-variant/15 rounded-sm border p-4">
									<div class="text-outline text-[10px] tracking-[0.2em] uppercase">
										Final Result
									</div>
									<div class="font-headline text-primary-container mt-2 text-2xl font-bold">
										{detailState.market.result.toUpperCase()}
									</div>
								</div>
							{/if}
						</div>
					{:else}
						<form class="space-y-4" onsubmit={handleSubmit}>
							<div class="grid grid-cols-2 gap-3">
								{#each ['yes', 'no'] as side}
									<button
										class="flex flex-col items-center justify-center rounded-sm border p-4 transition-all {selectedSide ===
										side
											? side === 'yes'
												? 'bg-primary-container/20 border-primary-container'
												: 'bg-error/20 border-error'
											: side === 'yes'
												? 'bg-surface-container border-primary-container/30 hover:bg-primary-container/10'
												: 'bg-surface-container border-error/30 hover:bg-error/10'}"
										onclick={() => (selectedSide = side as 'yes' | 'no')}
										type="button"
									>
										<span
											class="{side === 'yes'
												? 'text-primary-container'
												: 'text-error'} font-headline text-lg font-bold uppercase">{side}</span
										>
									</button>
								{/each}
							</div>

							<label class="block space-y-2">
								<div class="flex items-center justify-between px-1">
									<span class="text-outline text-[10px] font-bold tracking-widest uppercase"
										>Stake Amount</span
									>
									<span class="text-outline font-mono text-[10px] uppercase"
										>Available: {detailState.availableBalance}</span
									>
								</div>
								<div class="group relative">
									<input
										bind:value={stakeAmount}
										class="bg-surface-container-lowest border-outline-variant/20 focus:ring-primary-container/30 placeholder:text-outline-variant w-full rounded-sm border p-4 pr-16 font-mono text-xl focus:ring-1 focus:outline-none"
										placeholder="0.00"
										required
										step="0.00000001"
										type="number"
									/>
									<span
										class="text-outline absolute top-1/2 right-4 -translate-y-1/2 text-sm font-bold"
										>USDC</span
									>
								</div>
							</label>

							<div class="bg-surface-container space-y-3 rounded-sm p-4">
								<div class="flex items-center justify-between text-xs">
									<span class="text-outline">Selected Side</span>
									<span
										class="{selectedSide === 'yes'
											? 'text-primary-container'
											: 'text-error'} font-mono font-bold uppercase">{selectedSide}</span
									>
								</div>
								<div class="flex items-center justify-between text-xs">
									<span class="text-outline">Locked Balance</span>
									<span class="text-on-surface font-mono">{detailState.lockedBalance}</span>
								</div>
							</div>

							{#if detailState.submitError}
								<div class="bg-surface-container-low border-error/20 text-error border p-4 text-sm">
									{detailState.submitError}
								</div>
							{/if}

							{#if detailState.successMessage}
								<div
									class="bg-surface-container-low border-primary-container/20 text-primary border p-4 text-sm"
								>
									{detailState.successMessage}
								</div>
							{/if}

							<button
								class="gradient-primary text-on-primary-fixed font-headline w-full rounded-sm py-4 font-bold tracking-[0.2em] uppercase transition-all hover:opacity-90 active:scale-[0.99] disabled:opacity-70"
								disabled={detailState.submitStatus === 'submitting'}
								type="submit"
							>
								{detailState.submitStatus === 'submitting'
									? 'Placing Position...'
									: 'Place Position'}
							</button>
						</form>
					{/if}
				</section>

				<section class="space-y-4">
					<h2 class="font-headline text-xs font-bold tracking-widest uppercase">
						Other Active Markets
					</h2>
					<div class="space-y-3">
						{#if detailState.relatedMarkets.length > 0}
							{#each detailState.relatedMarkets as relatedMarket}
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
			</div>
		</div>
	{/if}
</main>
