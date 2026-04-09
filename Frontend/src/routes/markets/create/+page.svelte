<script lang="ts">
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  import type {
    MarketCreateContextSymbolResponse,
    MarketValidationModelResponse
  } from '$lib/api-types';
  import { loadCreateMarketContext, submitCreateMarket } from '$lib/create-market-data';
  import { ensureGuestSession } from '$lib/guest-session';
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
    validationModels: [] as MarketValidationModelResponse[]
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

  function settlementSummary(): string {
    const model = selectedModel();
    if (!model) {
      return 'The market rule will appear here once the market type is selected.';
    }

    if (model.market_type === 'candle_direction') {
      return `This market resolves from the ${form.sourceInterval || 'selected'} candle close.`;
    }

    if (model.market_type === 'funding_threshold') {
      if (form.conditionOperator === 'positive') {
        return 'This market settles YES when the funding value is above zero at the funding checkpoint.';
      }

      if (form.conditionOperator === 'negative') {
        return 'This market settles YES when the funding value is below zero at the funding checkpoint.';
      }
    }

    if (needsThreshold() && form.thresholdValue.trim() !== '') {
      return `This market settles YES when the observed value is ${operatorLabel(form.conditionOperator).toLowerCase()} ${form.thresholdValue}.`;
    }

    return `This market uses ${operatorLabel(form.conditionOperator).toLowerCase()} as the settlement rule.`;
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
        title: form.title.trim(),
        symbol: form.symbol,
        market_type: form.marketType,
        condition_operator: form.conditionOperator,
        creator_side: form.creatorSide,
        creator_stake_amount: form.creatorStakeAmount.trim(),
        threshold_value: needsThreshold() ? form.thresholdValue.trim() : '',
        source_type: model.source_type,
        source_interval: model.requires_interval ? form.sourceInterval : '',
        reference_value: '',
        expiry_time: toUtcISOString(form.expiryTime)
      });

      createState.submitStatus = 'idle';
      await goto(`/markets/${createdMarket.id}`);
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
    future.setSeconds(0, 0);

    const timezoneOffset = future.getTimezoneOffset();
    const localTime = new Date(future.getTime() - timezoneOffset * 60 * 1000);

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

<main class="pt-24 pb-12 px-6 md:px-10">
  <div class="max-w-6xl mx-auto">
    <header class="mb-10">
      <h1 class="font-headline text-4xl md:text-5xl font-extrabold tracking-tighter text-primary mb-2">Create Market</h1>
      <p class="text-outline text-sm tracking-wide max-w-2xl">Define one supported market, pick your side, add your opening stake, and submit the full action once.</p>
    </header>

    {#if createState.loadStatus === 'loading'}
      <section class="bg-surface-container-low p-8 border border-outline-variant/15">
        <p class="text-sm text-outline">Loading available symbols and market rules.</p>
      </section>
    {:else if createState.loadStatus === 'error'}
      <section class="bg-surface-container-low p-8 border border-error/20 space-y-4">
        <p class="text-sm text-error">{createState.loadError}</p>
        <button class="px-5 py-3 gradient-primary text-on-primary-fixed text-xs font-bold uppercase tracking-[0.2em]" onclick={initializePage}>
          Try Again
        </button>
      </section>
    {:else}
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-10">
        <form class="lg:col-span-7 space-y-8" onsubmit={handleSubmit}>
          <section class="bg-surface-container-low p-6 rounded-sm border-l-2 border-primary-container">
            <div class="mb-4 text-xs font-headline uppercase tracking-widest text-outline">Market Question</div>
            <label class="block">
              <span class="sr-only">Market question</span>
              <textarea
                bind:value={form.title}
                class="w-full bg-surface-container-lowest border border-outline-variant/20 text-xl font-headline text-on-surface placeholder:text-outline-variant/30 focus:ring-1 focus:ring-primary-container/30 p-4 rounded-sm resize-none"
                placeholder="Will BTC funding rate turn negative before the next funding checkpoint?"
                rows="3"
                required
              ></textarea>
            </label>
          </section>

          <section class="bg-surface-container-low p-6 rounded-sm">
            <div class="mb-6 text-xs font-headline uppercase tracking-widest text-outline">Market Setup</div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <label class="block">
                <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Symbol</span>
                <select
                  bind:value={form.symbol}
                  class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30"
                >
                  {#each createState.symbols as symbol}
                    <option value={symbol.symbol}>{symbol.symbol}</option>
                  {/each}
                </select>
              </label>

              <label class="block">
                <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Market Type</span>
                <select
                  bind:value={form.marketType}
                  class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30"
                  onchange={syncFormWithModel}
                >
                  {#each createState.validationModels as model}
                    <option value={model.market_type}>{marketTypeLabel(model.market_type)}</option>
                  {/each}
                </select>
              </label>

              <label class="block">
                <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Rule</span>
                <select
                  bind:value={form.conditionOperator}
                  class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30"
                  onchange={syncFormWithModel}
                >
                  {#each selectedModel()?.allowed_operators ?? [] as operator}
                    <option value={operator}>{operatorLabel(operator)}</option>
                  {/each}
                </select>
              </label>

              {#if selectedModel()?.requires_interval}
                <label class="block">
                  <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Interval</span>
                  <select
                    bind:value={form.sourceInterval}
                    class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30"
                  >
                    {#each selectedModel()?.allowed_intervals ?? [] as interval}
                      <option value={interval}>{interval}</option>
                    {/each}
                  </select>
                </label>
              {/if}

              {#if needsThreshold()}
                <label class="block">
                  <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Threshold Value</span>
                  <input
                    bind:value={form.thresholdValue}
                    class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30 outline-none"
                    placeholder="0.00"
                    required={needsThreshold()}
                    type="text"
                  />
                </label>
              {/if}

              <label class="block">
                <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Expiry Time</span>
                <input
                  bind:value={form.expiryTime}
                  class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30 outline-none"
                  required
                  type="datetime-local"
                />
              </label>

              <div class="md:col-span-2 bg-surface-container-lowest p-4 border border-outline-variant/10 rounded-sm">
                <div class="text-[10px] text-outline uppercase tracking-widest mb-3">Settlement Source</div>
                <div class="text-sm text-on-surface font-medium">{sourceTypeLabel(selectedModel()?.source_type ?? '')}</div>
                <div class="mt-2 text-xs text-outline">The backend decides settlement using the supported source for this market type.</div>
              </div>
            </div>
          </section>

          <section class="bg-surface-container-lowest p-4 border border-outline-variant/20 flex gap-4 items-start">
            <span class="material-symbols-outlined text-tertiary-fixed-dim mt-1">info</span>
            <p class="text-xs text-on-surface leading-relaxed">{settlementSummary()}</p>
          </section>

          <section class="bg-surface-container-low p-6 rounded-sm">
            <div class="mb-6 text-xs font-headline uppercase tracking-widest text-outline">Creator Position</div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <div class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Choose Side</div>
                <div class="flex gap-2">
                  {#each ['yes', 'no'] as side}
                    <button
                      class="flex-1 py-3 border text-xs font-bold rounded-sm uppercase transition-colors {form.creatorSide === side ? 'border-primary-container bg-primary-container/10 text-primary-container' : 'border-outline-variant/30 text-outline'}"
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
                <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">Opening Stake</span>
                <input
                  bind:value={form.creatorStakeAmount}
                  class="w-full bg-surface-container-lowest border border-outline-variant/20 text-sm text-on-surface py-3 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30 outline-none"
                  placeholder="500.00"
                  required
                  type="number"
                />
              </label>
            </div>
          </section>

          {#if createState.submitError}
            <div class="bg-surface-container-low border border-error/20 p-4 text-sm text-error">
              {createState.submitError}
            </div>
          {/if}

          <button
            class="w-full py-5 gradient-primary text-on-primary-fixed font-headline font-extrabold text-sm tracking-[0.2em] uppercase rounded-sm hover:shadow-[0_0_20px_rgba(0,240,255,0.4)] transition-all active:scale-95 disabled:opacity-70"
            disabled={createState.submitStatus === 'submitting'}
            type="submit"
          >
            {createState.submitStatus === 'submitting' ? 'Creating Market...' : 'Create Market'}
          </button>
        </form>

        <aside class="lg:col-span-5 space-y-6">
          <div class="sticky top-24">
            <h2 class="font-headline text-xs uppercase tracking-widest text-outline mb-4">Preview</h2>
            <div class="bg-[rgba(29,32,35,0.6)] backdrop-blur-xl p-6 rounded-sm border border-primary-container/20 shadow-[0_0_12px_rgba(0,219,233,0.3)]">
              <div class="flex justify-between items-start mb-6 gap-4">
                <div class="flex items-center gap-2 px-2 py-1 bg-surface-container-highest rounded-sm border border-outline-variant/30">
                  <span class="text-[10px] font-bold tracking-tighter uppercase">{form.symbol || 'Symbol'} / {marketTypeLabel(form.marketType || 'market')}</span>
                </div>
                <div class="text-right">
                  <div class="text-[10px] text-outline uppercase">Creator Side</div>
                  <div class="text-sm font-mono text-primary-container uppercase">{form.creatorSide}</div>
                </div>
              </div>

              <h3 class="text-xl font-headline font-bold text-primary mb-6 leading-tight">
                {form.title || 'Your market question will appear here.'}
              </h3>

              <div class="grid grid-cols-2 gap-4 mb-6">
                <div class="bg-surface-container-lowest p-4 border-l border-primary-container">
                  <span class="text-[10px] text-outline uppercase block mb-1">Source</span>
                  <span class="text-sm font-headline font-bold text-on-surface">{sourceTypeLabel(selectedModel()?.source_type ?? '')}</span>
                </div>
                <div class="bg-surface-container-lowest p-4 border-l border-error/40">
                  <span class="text-[10px] text-outline uppercase block mb-1">Stake</span>
                  <span class="text-sm font-headline font-bold text-on-surface">{form.creatorStakeAmount || '0.00'}</span>
                </div>
              </div>

              <div class="space-y-3 pt-4 border-t border-outline-variant/10 text-xs">
                <div class="flex justify-between gap-4">
                  <span class="text-outline uppercase">Mark Price</span>
                  <span class="text-on-surface font-mono">{formatNumber(selectedSymbol()?.mark_price)}</span>
                </div>
                <div class="flex justify-between gap-4">
                  <span class="text-outline uppercase">Funding Rate</span>
                  <span class="text-on-surface font-mono">{formatNumber(selectedSymbol()?.funding_rate)}</span>
                </div>
                <div class="flex justify-between gap-4">
                  <span class="text-outline uppercase">24H Volume</span>
                  <span class="text-on-surface font-mono">{formatNumber(selectedSymbol()?.volume_24h)}</span>
                </div>
                <div class="flex justify-between gap-4">
                  <span class="text-outline uppercase">Expiry</span>
                  <span class="text-on-surface font-mono">{form.expiryTime || 'Pick an expiry'}</span>
                </div>
              </div>
            </div>

            <div class="mt-6 p-4 bg-surface-container-low border border-outline-variant/15 rounded-sm">
              <div class="text-[10px] uppercase font-bold text-primary tracking-widest mb-2">Important</div>
              <p class="text-xs text-on-surface/80 leading-relaxed">Creating a market also opens your first position. If the backend rejects any part of that action, nothing should be half-created.</p>
            </div>
          </div>
        </aside>
      </div>
    {/if}
  </div>
</main>
