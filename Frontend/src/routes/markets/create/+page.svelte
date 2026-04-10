<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onMount } from 'svelte';

	import type {
		MarketCreateContextSymbolResponse,
		MarketValidationModelResponse
	} from '$lib/api-types';
	import { loadCreateMarketContext, submitCreateMarket } from '$lib/create-market-data';
	import { ensureGuestSession } from '$lib/guest-session';
	import { formatAmount } from '$lib/number-display';
	import TopNavBar from '$lib/components/TopNavBar.svelte';

	type LoadStatus = 'loading' | 'ready' | 'error';
	type SubmitStatus = 'idle' | 'submitting' | 'error';

	interface FormState {
		title: string;
		symbol: string;
		marketType: string;
		conditionOperator: string;
		sourceInterval: string;
		thresholdValue: string;
		expiryTime: string;
		creatorSide: 'yes' | 'no';
		creatorStakeAmount: string;
	}

	const createState = $state({
		loadStatus: 'loading' as LoadStatus,
		submitStatus: 'idle' as SubmitStatus,
		loadError: null as string | null,
		submitError: null as string | null,
		symbols: [] as MarketCreateContextSymbolResponse[],
		validationModels: [] as MarketValidationModelResponse[],
		priceThresholdCreationBandPercent: ''
	});

	const form = $state<FormState>({
		title: '',
		symbol: '',
		marketType: '',
		conditionOperator: '',
		sourceInterval: '',
		thresholdValue: '',
		expiryTime: defaultExpiryTime(),
		creatorSide: 'yes',
		creatorStakeAmount: ''
	});

	onMount(() => {
		void initializePage();
	});

	async function initializePage(): Promise<void> {
		createState.loadStatus = 'loading';
		createState.loadError = null;

		try {
			await ensureGuestSession();
			const context = await loadCreateMarketContext();

			createState.symbols = context.symbols;
			createState.validationModels = context.validation_models;
			createState.priceThresholdCreationBandPercent =
				context.price_threshold_creation_band_percent;
			createState.loadStatus = 'ready';

			initializeFormDefaults();
		} catch (error) {
			createState.loadStatus = 'error';
			createState.loadError = toErrorMessage(error, 'Unable to load the create-market flow.');
		}
	}

	function initializeFormDefaults(): void {
		const preferredModel = pickDefaultModel();
		const defaultSymbol = createState.symbols[0]?.symbol ?? '';

		form.symbol = defaultSymbol;
		form.marketType = preferredModel?.market_type ?? '';

		syncFormWithModel();
	}

	function pickDefaultModel(): MarketValidationModelResponse | undefined {
		return (
			createState.validationModels.find((model) => model.market_type === 'funding_threshold') ??
			createState.validationModels[0]
		);
	}

	function selectedModel(): MarketValidationModelResponse | undefined {
		return createState.validationModels.find((model) => model.market_type === form.marketType);
	}

	function selectedSymbol(): MarketCreateContextSymbolResponse | undefined {
		return createState.symbols.find((symbol) => symbol.symbol === form.symbol);
	}

	function syncFormWithModel(): void {
		const model = selectedModel();
		if (!model) {
			return;
		}

		if (!model.allowed_operators.includes(form.conditionOperator)) {
			form.conditionOperator = model.allowed_operators[0] ?? '';
		}

		if (model.requires_interval) {
			const intervals = model.allowed_intervals ?? [];
			if (!intervals.includes(form.sourceInterval)) {
				form.sourceInterval = intervals[0] ?? '';
			}
		} else {
			form.sourceInterval = '';
		}

		if (!needsThreshold()) {
			form.thresholdValue = '';
		}
	}

	function needsThreshold(): boolean {
		const model = selectedModel();
		if (!model?.requires_threshold) {
			return false;
		}

		return !['positive', 'negative'].includes(form.conditionOperator);
	}

	function usesAutomaticExpiry(): boolean {
		return selectedModel()?.requires_interval ?? false;
	}

	function isPriceThresholdMarket(): boolean {
		return selectedModel()?.market_type === 'price_threshold';
	}

	function operatorLabel(value: string): string {
		const labels: Record<string, string> = {
			gt: 'Greater Than',
			gte: 'Greater Than Or Equal',
			lt: 'Less Than',
			lte: 'Less Than Or Equal',
			bullish_close: 'Bullish Close',
			bearish_close: 'Bearish Close',
			positive: 'Positive',
			negative: 'Negative'
		};

		return labels[value] ?? value.replaceAll('_', ' ');
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

	function formatNumber(value?: string): string {
		return value && value.trim() !== '' ? value : 'Not available';
	}

	function decimalScale(value?: string): number {
		if (!value) {
			return 0;
		}

		const parts = value.trim().split('.', 2);
		if (parts.length < 2) {
			return 0;
		}

		return parts[1].replace(/0+$/, '').length;
	}

	function parseDecimalNumber(value?: string): number | null {
		if (!value) {
			return null;
		}

		const parsed = Number(value);
		return Number.isFinite(parsed) ? parsed : null;
	}

	function thresholdDisplayScale(): number {
		const scale = Math.max(
			decimalScale(selectedSymbol()?.min_tick),
			decimalScale(selectedSymbol()?.mark_price)
		);
		return Math.min(Math.max(scale, 2), 8);
	}

	function formatThresholdNumber(value: number): string {
		return value.toFixed(thresholdDisplayScale());
	}

	function thresholdGuardrail():
		| {
				referenceValue: string;
				tickSize: string;
				bandPercent: string;
				minThreshold: string;
				maxThreshold: string;
		  }
		| null {
		if (!isPriceThresholdMarket()) {
			return null;
		}

		const symbol = selectedSymbol();
		const referenceValue = parseDecimalNumber(symbol?.mark_price);
		const tickSize = parseDecimalNumber(symbol?.min_tick);
		const bandPercent = parseDecimalNumber(createState.priceThresholdCreationBandPercent);
		if (referenceValue === null || tickSize === null || bandPercent === null || tickSize <= 0) {
			return null;
		}

		const bandRatio = bandPercent / 100;
		const lowerBound = referenceValue * (1 - bandRatio);
		const upperBound = referenceValue * (1 + bandRatio);
		const increasingRule = ['gt', 'gte'].includes(form.conditionOperator);
		const minThreshold = increasingRule ? referenceValue + tickSize : lowerBound;
		const maxThreshold = increasingRule ? upperBound : referenceValue - tickSize;

		if (minThreshold > maxThreshold) {
			return null;
		}

		return {
			referenceValue: formatThresholdNumber(referenceValue),
			tickSize: formatThresholdNumber(tickSize),
			bandPercent: createState.priceThresholdCreationBandPercent,
			minThreshold: formatThresholdNumber(minThreshold),
			maxThreshold: formatThresholdNumber(maxThreshold)
		};
	}

	function currentPriceText(): string {
		return formatNumber(selectedSymbol()?.mark_price);
	}

	function normalizeText(value: string | number | null | undefined): string {
		if (value === null || value === undefined) {
			return '';
		}

		return String(value).trim();
	}

	function settlementSummary(): string {
		const model = selectedModel();
		if (!model) {
			return 'The market rule will appear here once the market type is selected.';
		}

		if (model.market_type === 'candle_direction') {
			return `This market auto-locks to the next ${form.sourceInterval || 'selected'} candle close after you submit it.`;
		}

		if (model.market_type === 'funding_threshold') {
			if (form.conditionOperator === 'positive') {
				return 'This market targets the next funding checkpoint after submission and settles YES when funding is above zero.';
			}

			if (form.conditionOperator === 'negative') {
				return 'This market targets the next funding checkpoint after submission and settles YES when funding is below zero.';
			}
		}

		const thresholdValue = normalizeText(form.thresholdValue);
		const guardrail = thresholdGuardrail();
		if (guardrail) {
			const directionText = ['gt', 'gte'].includes(form.conditionOperator)
				? `${guardrail.referenceValue} and ${guardrail.maxThreshold}`
				: `${guardrail.minThreshold} and ${guardrail.referenceValue}`;
			return `This market settles YES when the observed value is ${operatorLabel(form.conditionOperator).toLowerCase()} the threshold you choose. Price thresholds must stay within ${guardrail.bandPercent}% of the creation reference, on the correct side of spot, so the valid zone is between ${directionText}.`;
		}

		if (needsThreshold() && thresholdValue !== '') {
			return `This market settles YES when the observed value is ${operatorLabel(form.conditionOperator).toLowerCase()} ${thresholdValue}.`;
		}

		return `This market uses ${operatorLabel(form.conditionOperator).toLowerCase()} as the settlement rule.`;
	}

	function timingSummary(): string {
		const model = selectedModel();
		if (!model) {
			return 'Timing is set when the market type is selected.';
		}

		if (model.market_type === 'candle_direction') {
			return `Auto-set to the next ${form.sourceInterval || 'selected'} candle close after submission.`;
		}

		if (model.market_type === 'funding_threshold') {
			return `Auto-set to the next funding checkpoint after submission (${nextFundingCheckpointPreview()}).`;
		}

		return form.expiryTime || 'Pick an expiry';
	}

	function nextFundingCheckpointPreview(): string {
		const currentTime = new Date();
		const nextHour = new Date(
			currentTime.getFullYear(),
			currentTime.getMonth(),
			currentTime.getDate(),
			currentTime.getHours() + 1,
			0,
			0,
			0
		);

		return new Intl.DateTimeFormat(undefined, {
			month: 'short',
			day: 'numeric',
			hour: 'numeric',
			minute: '2-digit'
		}).format(nextHour);
	}

	async function handleSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		createState.submitStatus = 'submitting';
		createState.submitError = null;

		const model = selectedModel();
		if (!model) {
			createState.submitStatus = 'error';
			createState.submitError = 'Pick a supported market type first.';
			return;
		}

		const player = await ensureGuestSession();
		if (!player) {
			createState.submitStatus = 'error';
			createState.submitError = 'Your guest session is not ready yet.';
			return;
		}

		try {
			const createdMarket = await submitCreateMarket({
				title: normalizeText(form.title),
				symbol: form.symbol,
				market_type: form.marketType,
				condition_operator: form.conditionOperator,
				creator_side: form.creatorSide,
				creator_stake_amount: normalizeText(form.creatorStakeAmount),
				threshold_value: needsThreshold() ? normalizeText(form.thresholdValue) : '',
				source_type: model.source_type,
				source_interval: model.requires_interval ? form.sourceInterval : '',
				reference_value: '',
				...(usesAutomaticExpiry() ? {} : { expiry_time: toUtcISOString(form.expiryTime) })
			});

			createState.submitStatus = 'idle';
			await goto(resolve(`/markets/${createdMarket.id}`));
		} catch (error) {
			createState.submitStatus = 'error';
			createState.submitError = toErrorMessage(error, 'Unable to create the market right now.');
		}
	}

	function toUtcISOString(value: string): string {
		return new Date(value).toISOString();
	}

	function defaultExpiryTime(): string {
		const future = new Date(Date.now() + 4 * 60 * 60 * 1000);
		const roundedLocalTime = new Date(
			future.getFullYear(),
			future.getMonth(),
			future.getDate(),
			future.getHours(),
			future.getMinutes(),
			0,
			0
		);

		const timezoneOffset = roundedLocalTime.getTimezoneOffset();
		const localTime = new Date(roundedLocalTime.getTime() - timezoneOffset * 60 * 1000);

		return localTime.toISOString().slice(0, 16);
	}

	function toErrorMessage(error: unknown, fallback: string): string {
		if (error instanceof Error && error.message) {
			return error.message;
		}

		return fallback;
	}
</script>

<svelte:head>
	<title>Create Market | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Markets" />

<main class="px-6 pt-24 pb-12 md:px-10">
	<div class="mx-auto max-w-6xl">
		<header class="mb-10">
			<h1
				class="font-headline text-primary mb-2 text-4xl font-extrabold tracking-tighter md:text-5xl"
			>
				Create Market
			</h1>
			<p class="text-outline max-w-2xl text-sm tracking-wide">
				Define one supported market, pick your side, add your opening stake, and submit the full
				action once.
			</p>
		</header>

		{#if createState.loadStatus === 'loading'}
			<section class="bg-surface-container-low border-outline-variant/15 border p-8">
				<p class="text-outline text-sm">Loading available symbols and market rules.</p>
			</section>
		{:else if createState.loadStatus === 'error'}
			<section class="bg-surface-container-low border-error/20 space-y-4 border p-8">
				<p class="text-error text-sm">{createState.loadError}</p>
				<button
					class="gradient-primary text-on-primary-fixed px-5 py-3 text-xs font-bold tracking-[0.2em] uppercase"
					onclick={initializePage}
				>
					Try Again
				</button>
			</section>
		{:else}
			<div class="grid grid-cols-1 gap-10 lg:grid-cols-12">
				<form class="space-y-8 lg:col-span-7" onsubmit={handleSubmit}>
					<section
						class="bg-surface-container-low border-primary-container rounded-sm border-l-2 p-6"
					>
						<div class="font-headline text-outline mb-4 text-xs tracking-widest uppercase">
							Market Question
						</div>
						<label class="block">
							<span class="sr-only">Market question</span>
							<textarea
								bind:value={form.title}
								class="bg-surface-container-lowest border-outline-variant/20 font-headline text-on-surface placeholder:text-outline-variant/30 focus:ring-primary-container/30 w-full resize-none rounded-sm border p-4 text-xl focus:ring-1"
								placeholder="Will BTC funding rate turn negative before the next funding checkpoint?"
								rows="3"
								required
							></textarea>
						</label>
					</section>

					<section class="bg-surface-container-low rounded-sm p-6">
						<div class="font-headline text-outline mb-6 text-xs tracking-widest uppercase">
							Market Setup
						</div>

						<div class="mb-6 grid grid-cols-1 gap-4 md:grid-cols-3">
							<div
								class="bg-surface-container-lowest border-outline-variant/10 rounded-sm border p-4"
							>
								<div class="text-outline text-[10px] tracking-widest uppercase">Current Price</div>
								<div class="text-primary mt-2 font-mono text-2xl font-bold">
									{currentPriceText()}
								</div>
								<div class="text-outline mt-1 text-[10px] uppercase">
									{form.symbol || 'Pick a symbol'}
								</div>
							</div>

							<div
								class="bg-surface-container-lowest border-outline-variant/10 rounded-sm border p-4"
							>
								<div class="text-outline text-[10px] tracking-widest uppercase">Funding Rate</div>
								<div class="text-on-surface mt-2 font-mono text-sm font-bold">
									{formatNumber(selectedSymbol()?.funding_rate)}
								</div>
							</div>

							<div
								class="bg-surface-container-lowest border-outline-variant/10 rounded-sm border p-4"
							>
								<div class="text-outline text-[10px] tracking-widest uppercase">24H Volume</div>
								<div class="text-on-surface mt-2 font-mono text-sm font-bold">
									{formatNumber(selectedSymbol()?.volume_24h)}
								</div>
							</div>
						</div>

						<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
							<label class="block">
								<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
									>Symbol</span
								>
								<select
									bind:value={form.symbol}
									class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm focus:ring-1"
								>
									{#each createState.symbols as symbol (symbol.symbol)}
										<option value={symbol.symbol}>{symbol.symbol}</option>
									{/each}
								</select>
							</label>

							<label class="block">
								<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
									>Market Type</span
								>
								<select
									bind:value={form.marketType}
									class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm focus:ring-1"
									onchange={syncFormWithModel}
								>
									{#each createState.validationModels as model (model.market_type)}
										<option value={model.market_type}>{marketTypeLabel(model.market_type)}</option>
									{/each}
								</select>
							</label>

							<label class="block">
								<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
									>Rule</span
								>
								<select
									bind:value={form.conditionOperator}
									class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm focus:ring-1"
									onchange={syncFormWithModel}
								>
									{#each selectedModel()?.allowed_operators ?? [] as operator (operator)}
										<option value={operator}>{operatorLabel(operator)}</option>
									{/each}
								</select>
							</label>

							{#if selectedModel()?.requires_interval}
								<label class="block">
									<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
										>Interval</span
									>
									<select
										bind:value={form.sourceInterval}
										class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm focus:ring-1"
									>
										{#each selectedModel()?.allowed_intervals ?? [] as interval (interval)}
											<option value={interval}>{interval}</option>
										{/each}
									</select>
								</label>
							{/if}

							{#if needsThreshold()}
								<label class="block">
									<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
										>Threshold Value</span
									>
									<input
										bind:value={form.thresholdValue}
										class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm outline-none focus:ring-1"
										placeholder="0.00"
										required={needsThreshold()}
										type="text"
									/>
									{#if thresholdGuardrail()}
										<div
											class="bg-surface-container mt-3 space-y-2 rounded-sm border border-primary-container/15 p-3 text-xs"
										>
											<div class="flex items-center justify-between gap-4">
												<span class="text-outline uppercase">Creation Reference</span>
												<span class="text-primary font-mono"
													>{thresholdGuardrail()?.referenceValue}</span
												>
											</div>
											<div class="flex items-center justify-between gap-4">
												<span class="text-outline uppercase">Allowed Threshold Range</span>
												<span class="text-on-surface font-mono"
													>{thresholdGuardrail()?.minThreshold} to
													{thresholdGuardrail()?.maxThreshold}</span
												>
											</div>
											<div class="flex items-center justify-between gap-4">
												<span class="text-outline uppercase">Band / Minimum Step</span>
												<span class="text-on-surface font-mono"
													>{thresholdGuardrail()?.bandPercent}% / {thresholdGuardrail()?.tickSize}</span
												>
											</div>
										</div>
									{/if}
								</label>
							{/if}

							{#if !usesAutomaticExpiry()}
								<label class="block">
									<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
										>Expiry Time</span
									>
									<div class="text-outline mb-2 text-xs leading-relaxed">
										Use the picker or click into the field and use your arrow keys to adjust the
										date and time.
									</div>
									<input
										bind:value={form.expiryTime}
										class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm outline-none focus:ring-1"
										required
										step="60"
										type="datetime-local"
									/>
								</label>
							{:else}
								<div
									class="bg-surface-container-lowest border-outline-variant/10 rounded-sm border p-4"
								>
									<div class="text-outline mb-2 text-[10px] tracking-widest uppercase">Timing</div>
									<div class="text-on-surface text-sm font-medium">{timingSummary()}</div>
									<div class="text-outline mt-2 text-xs">
										The backend sets the final expiry when you submit so interval markets stay
										aligned with live settlement boundaries.
									</div>
								</div>
							{/if}

							<div
								class="bg-surface-container-lowest border-outline-variant/10 rounded-sm border p-4 md:col-span-2"
							>
								<div class="text-outline mb-3 text-[10px] tracking-widest uppercase">
									Settlement Source
								</div>
								<div class="text-on-surface text-sm font-medium">
									{sourceTypeLabel(selectedModel()?.source_type ?? '')}
								</div>
								<div class="text-outline mt-2 text-xs">
									The backend decides settlement using the supported source for this market type.
								</div>
							</div>
						</div>
					</section>

					<section
						class="bg-surface-container-lowest border-outline-variant/20 flex items-start gap-4 border p-4"
					>
						<span class="material-symbols-outlined text-tertiary-fixed-dim mt-1">info</span>
						<p class="text-on-surface text-xs leading-relaxed">{settlementSummary()}</p>
					</section>

					<section class="bg-surface-container-low rounded-sm p-6">
						<div class="font-headline text-outline mb-6 text-xs tracking-widest uppercase">
							Creator Position
						</div>
						<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
							<div>
								<div class="text-outline mb-2 block text-[10px] tracking-widest uppercase">
									Choose Side
								</div>
								<div class="flex gap-2">
									{#each ['yes', 'no'] as side (side)}
										<button
											class="flex-1 rounded-sm border py-3 text-xs font-bold uppercase transition-colors {form.creatorSide ===
											side
												? 'border-primary-container bg-primary-container/10 text-primary-container'
												: 'border-outline-variant/30 text-outline'}"
											onclick={(event) => {
												event.preventDefault();
												form.creatorSide = side as 'yes' | 'no';
											}}
											type="button"
										>
											{side}
										</button>
									{/each}
								</div>
							</div>

							<label class="block">
								<span class="text-outline mb-2 block text-[10px] tracking-widest uppercase"
									>Opening Stake</span
								>
								<input
									bind:value={form.creatorStakeAmount}
									class="bg-surface-container-lowest border-outline-variant/20 text-on-surface focus:ring-primary-container/30 w-full rounded-sm border px-4 py-3 text-sm outline-none focus:ring-1"
									min="1"
									placeholder="500"
									required
									step="1"
									type="number"
								/>
							</label>
						</div>
					</section>

					{#if createState.submitError}
						<div class="bg-surface-container-low border-error/20 text-error border p-4 text-sm">
							{createState.submitError}
						</div>
					{/if}

					<button
						class="gradient-primary text-on-primary-fixed font-headline w-full rounded-sm py-5 text-sm font-extrabold tracking-[0.2em] uppercase transition-all hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] active:scale-95 disabled:opacity-70"
						disabled={createState.submitStatus === 'submitting'}
						type="submit"
					>
						{createState.submitStatus === 'submitting' ? 'Creating Market...' : 'Create Market'}
					</button>
				</form>

				<aside class="space-y-6 lg:col-span-5">
					<div class="sticky top-24">
						<h2 class="font-headline text-outline mb-4 text-xs tracking-widest uppercase">
							Preview
						</h2>
						<div
							class="border-primary-container/20 rounded-sm border bg-[rgba(29,32,35,0.6)] p-6 shadow-[0_0_12px_rgba(0,219,233,0.3)] backdrop-blur-xl"
						>
							<div class="mb-6 flex items-start justify-between gap-4">
								<div
									class="bg-surface-container-highest border-outline-variant/30 flex items-center gap-2 rounded-sm border px-2 py-1"
								>
									<span class="text-[10px] font-bold tracking-tighter uppercase"
										>{form.symbol || 'Symbol'} / {marketTypeLabel(
											form.marketType || 'market'
										)}</span
									>
								</div>
								<div class="text-right">
									<div class="text-outline text-[10px] uppercase">Creator Side</div>
									<div class="text-primary-container font-mono text-sm uppercase">
										{form.creatorSide}
									</div>
								</div>
							</div>

							<h3 class="font-headline text-primary mb-6 text-xl leading-tight font-bold">
								{form.title || 'Your market question will appear here.'}
							</h3>

							<div class="mb-6 grid grid-cols-2 gap-4">
								<div class="bg-surface-container-lowest border-primary-container border-l p-4">
									<span class="text-outline mb-1 block text-[10px] uppercase">Source</span>
									<span class="font-headline text-on-surface text-sm font-bold"
										>{sourceTypeLabel(selectedModel()?.source_type ?? '')}</span
									>
								</div>
								<div class="bg-surface-container-lowest border-error/40 border-l p-4">
									<span class="text-outline mb-1 block text-[10px] uppercase">Stake</span>
									<span class="font-headline text-on-surface text-sm font-bold"
										>{formatAmount(form.creatorStakeAmount || '0')}</span
									>
								</div>
							</div>

							<div
								class="bg-surface-container-lowest border-outline-variant/10 mb-6 rounded-sm border p-4"
							>
								<div class="flex items-center justify-between gap-4">
									<span class="text-outline text-[10px] tracking-widest uppercase"
										>Current Price</span
									>
									<span class="text-primary font-mono text-lg font-bold">{currentPriceText()}</span>
								</div>
							</div>

							<div class="border-outline-variant/10 space-y-3 border-t pt-4 text-xs">
								<div class="flex justify-between gap-4">
									<span class="text-outline uppercase">Mark Price</span>
									<span class="text-on-surface font-mono"
										>{formatNumber(selectedSymbol()?.mark_price)}</span
									>
								</div>
								<div class="flex justify-between gap-4">
									<span class="text-outline uppercase">Funding Rate</span>
									<span class="text-on-surface font-mono"
										>{formatNumber(selectedSymbol()?.funding_rate)}</span
									>
								</div>
								<div class="flex justify-between gap-4">
									<span class="text-outline uppercase">24H Volume</span>
									<span class="text-on-surface font-mono"
										>{formatNumber(selectedSymbol()?.volume_24h)}</span
									>
								</div>
								{#if thresholdGuardrail()}
									<div class="flex justify-between gap-4">
										<span class="text-outline uppercase">Threshold Band</span>
										<span class="text-on-surface font-mono"
											>{thresholdGuardrail()?.minThreshold} to
											{thresholdGuardrail()?.maxThreshold}</span
										>
									</div>
								{/if}
								<div class="flex justify-between gap-4">
									<span class="text-outline uppercase"
										>{usesAutomaticExpiry() ? 'Timing' : 'Expiry'}</span
									>
									<span class="text-on-surface font-mono">{timingSummary()}</span>
								</div>
							</div>
						</div>

						<div
							class="bg-surface-container-low border-outline-variant/15 mt-6 rounded-sm border p-4"
						>
							<div class="text-primary mb-2 text-[10px] font-bold tracking-widest uppercase">
								Important
							</div>
							<p class="text-on-surface/80 text-xs leading-relaxed">
								Creating a market also opens your first position. If the backend rejects any part of
								that action, nothing should be half-created.
							</p>
						</div>
					</div>
				</aside>
			</div>
		{/if}
	</div>
</main>
