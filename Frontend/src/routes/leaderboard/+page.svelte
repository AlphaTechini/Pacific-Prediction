<script lang="ts">
	import type {
		ActivityLeaderboardEntryResponse,
		CreatorLeaderboardEntryResponse,
		PredictorLeaderboardEntryResponse,
		StreakLeaderboardEntryResponse
	} from '$lib/api-types';
	import TopNavBar from '$lib/components/TopNavBar.svelte';

	import type { PageData } from './$types';

	type LeaderTab = 'top_predictors' | 'top_creators' | 'best_streaks' | 'most_active';

	interface ViewMetric {
		label: string;
		value: string;
	}

	interface ViewEntry {
		rank: number;
		playerID: string;
		displayName: string;
		accentLabel: string;
		accentValue: string;
		metrics: [ViewMetric, ViewMetric, ViewMetric];
	}

	interface SummaryCard {
		label: string;
		value: string;
		tone: 'primary' | 'accent' | 'default';
	}

	let { data }: { data: PageData } = $props();
	let activeTab = $state<LeaderTab>('top_predictors');

	const tabs: Array<{ key: LeaderTab; label: string; description: string }> = [
		{
			key: 'top_predictors',
			label: 'Top Predictors',
			description: 'Realized profit ranked first, then win rate and resolved volume.'
		},
		{
			key: 'top_creators',
			label: 'Top Creators',
			description: 'Created market count ranked first, then reach and participation.'
		},
		{
			key: 'best_streaks',
			label: 'Best Streaks',
			description: 'Current consecutive wins ranked first, then longest streak and win rate.'
		},
		{
			key: 'most_active',
			label: 'Most Active',
			description: 'Total positions ranked first, with market creation as the tie-breaker.'
		}
	];

	const tabColumns: Record<LeaderTab, { accent: string; labels: [string, string, string] }> = {
		top_predictors: {
			accent: 'Net',
			labels: ['Win Rate', 'Resolved', 'Staked']
		},
		top_creators: {
			accent: 'Markets',
			labels: ['Participants', 'Positions', 'Resolved']
		},
		best_streaks: {
			accent: 'Current',
			labels: ['Longest', 'Win Rate', 'Net']
		},
		most_active: {
			accent: 'Positions',
			labels: ['Open', 'Resolved', 'Markets']
		}
	};

	function summaryCards(): SummaryCard[] {
		const overview = data.leaderboard?.overview;
		if (!overview) {
			return [];
		}

		return [
			{
				label: 'Total Predictions',
				value: formatInteger(overview.total_predictions),
				tone: 'primary'
			},
			{
				label: 'Resolved Predictions',
				value: formatInteger(overview.resolved_predictions),
				tone: 'default'
			},
			{
				label: 'Active Predictors',
				value: formatInteger(overview.active_predictors),
				tone: 'accent'
			},
			{
				label: 'Average Win Rate',
				value: formatPercent(overview.average_win_rate),
				tone: 'default'
			}
		];
	}

	function activeEntries(): ViewEntry[] {
		const leaderboard = data.leaderboard;
		if (!leaderboard) {
			return [];
		}

		switch (activeTab) {
			case 'top_predictors':
				return leaderboard.top_predictors.map(mapPredictorEntry);
			case 'top_creators':
				return leaderboard.top_creators.map(mapCreatorEntry);
			case 'best_streaks':
				return leaderboard.best_streaks.map(mapStreakEntry);
			case 'most_active':
				return leaderboard.most_active.map(mapActivityEntry);
		}
	}

	function podiumEntries(): ViewEntry[] {
		return activeEntries().slice(0, 3);
	}

	function tableEntries(): ViewEntry[] {
		return activeEntries().slice(3);
	}

	function activeDescription(): string {
		return tabs.find((tab) => tab.key === activeTab)?.description ?? '';
	}

	function activeColumnLabels(): [string, string, string] {
		return tabColumns[activeTab].labels;
	}

	function activeAccentLabel(): string {
		return tabColumns[activeTab].accent;
	}

	function mapPredictorEntry(entry: PredictorLeaderboardEntryResponse): ViewEntry {
		return {
			rank: entry.rank,
			playerID: entry.player_id,
			displayName: entry.display_name,
			accentLabel: 'Net',
			accentValue: formatSignedAmount(entry.net_profit),
			metrics: [
				{ label: 'Win Rate', value: formatPercent(entry.win_rate) },
				{ label: 'Resolved', value: formatInteger(entry.resolved_positions) },
				{ label: 'Staked', value: formatAmount(entry.total_staked) }
			]
		};
	}

	function mapCreatorEntry(entry: CreatorLeaderboardEntryResponse): ViewEntry {
		return {
			rank: entry.rank,
			playerID: entry.player_id,
			displayName: entry.display_name,
			accentLabel: 'Markets',
			accentValue: formatInteger(entry.created_markets),
			metrics: [
				{ label: 'Participants', value: formatInteger(entry.unique_participants) },
				{ label: 'Positions', value: formatInteger(entry.total_positions) },
				{ label: 'Resolved', value: formatInteger(entry.resolved_markets) }
			]
		};
	}

	function mapStreakEntry(entry: StreakLeaderboardEntryResponse): ViewEntry {
		return {
			rank: entry.rank,
			playerID: entry.player_id,
			displayName: entry.display_name,
			accentLabel: 'Current',
			accentValue: formatInteger(entry.current_win_streak),
			metrics: [
				{ label: 'Longest', value: formatInteger(entry.longest_win_streak) },
				{ label: 'Win Rate', value: formatPercent(entry.win_rate) },
				{ label: 'Net', value: formatSignedAmount(entry.net_profit) }
			]
		};
	}

	function mapActivityEntry(entry: ActivityLeaderboardEntryResponse): ViewEntry {
		return {
			rank: entry.rank,
			playerID: entry.player_id,
			displayName: entry.display_name,
			accentLabel: 'Positions',
			accentValue: formatInteger(entry.total_positions),
			metrics: [
				{ label: 'Open', value: formatInteger(entry.open_positions) },
				{ label: 'Resolved', value: formatInteger(entry.resolved_positions) },
				{ label: 'Markets', value: formatInteger(entry.created_markets) }
			]
		};
	}

	function formatInteger(value: number): string {
		return new Intl.NumberFormat().format(value);
	}

	function formatAmount(value: string): string {
		const numericValue = Number(value);
		if (!Number.isFinite(numericValue)) {
			return value;
		}

		return new Intl.NumberFormat(undefined, {
			minimumFractionDigits: 2,
			maximumFractionDigits: 2
		}).format(numericValue);
	}

	function formatSignedAmount(value: string): string {
		const numericValue = Number(value);
		if (!Number.isFinite(numericValue)) {
			return value;
		}

		const prefix = numericValue > 0 ? '+' : '';
		return `${prefix}${formatAmount(value)}`;
	}

	function formatPercent(value: string): string {
		const numericValue = Number(value);
		if (!Number.isFinite(numericValue)) {
			return `${value}%`;
		}

		return `${numericValue.toFixed(1)}%`;
	}

	function formatGeneratedAt(value: string | undefined): string {
		if (!value) {
			return 'Not available';
		}

		const timestamp = new Date(value);
		if (Number.isNaN(timestamp.getTime())) {
			return value;
		}

		return new Intl.DateTimeFormat(undefined, {
			dateStyle: 'medium',
			timeStyle: 'short'
		}).format(timestamp);
	}

	function initialsFor(name: string): string {
		const cleaned = name.replace(/[^A-Za-z0-9_ ]+/g, ' ').trim();
		const parts = cleaned.split(/[\s_]+/).filter(Boolean);
		if (parts.length === 0) {
			return 'PP';
		}

		if (parts.length === 1) {
			return parts[0].slice(0, 2).toUpperCase();
		}

		return `${parts[0][0] ?? ''}${parts[1][0] ?? ''}`.toUpperCase();
	}
</script>

<svelte:head>
	<title>Leaderboard | Pacifica Pulse</title>
</svelte:head>

<TopNavBar activePage="Leaderboard" />

<main class="flex-grow px-6 pt-24 pb-16">
	<div class="mx-auto flex w-full max-w-7xl flex-col gap-10">
		<header class="flex flex-col gap-4 lg:flex-row lg:items-end lg:justify-between">
			<div class="space-y-3">
				<p class="text-primary-container text-[10px] tracking-[0.3em] uppercase">
					Leaderboard Snapshot
				</p>
				<h1 class="font-headline text-on-surface text-4xl font-bold tracking-tight md:text-5xl">
					Leaderboard
				</h1>
				<p class="text-outline max-w-2xl text-sm">
					Live ranking across predictors, creators, streaks, and activity. Every tab is derived from
					market and position history already stored by the backend.
				</p>
			</div>

			<div
				class="border-outline-variant/20 bg-surface-container-low text-outline border px-5 py-4 text-sm"
			>
				<div class="text-primary-container text-[10px] tracking-[0.2em] uppercase">Updated</div>
				<div class="text-on-surface mt-2 font-medium">
					{formatGeneratedAt(data.leaderboard?.generated_at)}
				</div>
			</div>
		</header>

		{#if summaryCards().length > 0}
			<section class="grid grid-cols-1 gap-4 md:grid-cols-2 xl:grid-cols-4">
				{#each summaryCards() as card}
					<div class="border-outline-variant/15 bg-surface-container-low border p-5">
						<div class="text-outline text-[10px] tracking-[0.2em] uppercase">{card.label}</div>
						<div
							class={`font-headline mt-3 text-3xl font-bold ${
								card.tone === 'primary'
									? 'text-primary'
									: card.tone === 'accent'
										? 'text-primary-container'
										: 'text-on-surface'
							}`}
						>
							{card.value}
						</div>
					</div>
				{/each}
			</section>
		{/if}

		<div class="grid grid-cols-1 gap-8 xl:grid-cols-12">
			<section class="xl:col-span-8">
				<div class="border-outline-variant/20 mb-6 flex flex-wrap gap-2 border-b pb-4">
					{#each tabs as tab}
						<button
							class={`px-4 py-2 text-xs font-bold tracking-[0.2em] uppercase transition-colors ${
								activeTab === tab.key
									? 'border-primary-container text-primary border-b-2'
									: 'text-outline hover:text-primary'
							}`}
							onclick={() => (activeTab = tab.key)}
						>
							{tab.label}
						</button>
					{/each}
				</div>

				<p class="text-outline mb-8 max-w-2xl text-sm">{activeDescription()}</p>

				{#if data.error}
					<section class="border-error/20 bg-surface-container-low border p-8">
						<p class="text-error text-sm">{data.error}</p>
					</section>
				{:else if activeEntries().length === 0}
					<section class="border-outline-variant/20 bg-surface-container-low border p-8">
						<p class="text-outline text-sm">
							The leaderboard will populate as soon as player activity lands in the database.
						</p>
					</section>
				{:else}
					<div class="grid grid-cols-1 gap-5 md:grid-cols-3">
						{#each podiumEntries() as entry}
							<div
								class={`bg-surface-container-low border p-6 ${
									entry.rank === 1
										? 'border-primary-container/40 shadow-[0_0_24px_rgba(0,219,233,0.08)]'
										: 'border-outline-variant/15'
								}`}
							>
								<div class="flex items-start justify-between gap-4">
									<div>
										<div class="text-outline text-[10px] tracking-[0.2em] uppercase">
											Rank {entry.rank}
										</div>
										<h2 class="font-headline text-on-surface mt-3 text-2xl font-bold">
											{entry.displayName}
										</h2>
									</div>
									<div
										class="border-outline-variant/20 bg-surface-container text-primary flex h-12 w-12 items-center justify-center border text-sm font-bold"
									>
										{initialsFor(entry.displayName)}
									</div>
								</div>

								<div class="border-outline-variant/10 mt-6 border-t pt-5">
									<div class="text-outline text-[10px] tracking-[0.2em] uppercase">
										{entry.accentLabel}
									</div>
									<div class="font-headline text-primary mt-2 text-3xl font-bold">
										{entry.accentValue}
									</div>
								</div>

								<div class="border-outline-variant/10 mt-6 grid grid-cols-3 gap-3 border-t pt-5">
									{#each entry.metrics as metric}
										<div>
											<div class="text-outline text-[10px] tracking-[0.2em] uppercase">
												{metric.label}
											</div>
											<div class="text-on-surface mt-2 text-sm font-semibold">{metric.value}</div>
										</div>
									{/each}
								</div>
							</div>
						{/each}
					</div>

					<div
						class="border-outline-variant/15 bg-surface-container-low mt-8 overflow-hidden border"
					>
						<table class="w-full border-collapse text-left">
							<thead class="border-outline-variant/20 bg-surface-container border-b">
								<tr>
									<th class="text-outline px-5 py-4 text-[10px] tracking-[0.2em] uppercase">Rank</th
									>
									<th class="text-outline px-5 py-4 text-[10px] tracking-[0.2em] uppercase"
										>Player</th
									>
									<th
										class="text-outline px-5 py-4 text-right text-[10px] tracking-[0.2em] uppercase"
									>
										{activeAccentLabel()}
									</th>
									{#each activeColumnLabels() as label}
										<th
											class="text-outline px-5 py-4 text-right text-[10px] tracking-[0.2em] uppercase"
										>
											{label}
										</th>
									{/each}
								</tr>
							</thead>
							<tbody class="divide-outline-variant/10 divide-y">
								{#each tableEntries() as entry}
									<tr class="hover:bg-surface-container transition-colors">
										<td class="text-outline px-5 py-4 font-mono text-sm">{entry.rank}</td>
										<td class="px-5 py-4">
											<div class="flex items-center gap-3">
												<div
													class="border-outline-variant/20 bg-surface-container text-primary flex h-9 w-9 items-center justify-center border text-xs font-bold"
												>
													{initialsFor(entry.displayName)}
												</div>
												<div class="min-w-0">
													<div class="text-on-surface truncate text-sm font-semibold">
														{entry.displayName}
													</div>
													<div class="text-outline truncate text-[10px] tracking-[0.2em] uppercase">
														{entry.playerID}
													</div>
												</div>
											</div>
										</td>
										<td class="font-headline text-primary px-5 py-4 text-right text-sm font-bold">
											{entry.accentValue}
										</td>
										{#each entry.metrics as metric}
											<td class="text-on-surface px-5 py-4 text-right text-sm">{metric.value}</td>
										{/each}
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</section>

			<aside class="space-y-6 xl:col-span-4">
				<section class="border-outline-variant/15 bg-surface-container-low border p-6">
					<h2 class="font-headline text-primary text-sm font-bold tracking-[0.2em] uppercase">
						Ranking Logic
					</h2>
					<div class="text-outline mt-5 space-y-4 text-sm">
						<p>
							Predictor rank uses settled performance only, so unresolved exposure does not distort
							the board.
						</p>
						<p>
							Creator rank rewards both output and reach by combining market count with
							participation depth.
						</p>
						<p>
							Streak rank follows consecutive settled wins, while activity rank tracks overall usage
							volume.
						</p>
					</div>
				</section>

				<section class="border-outline-variant/15 bg-surface-container-low border p-6">
					<h2 class="font-headline text-primary text-sm font-bold tracking-[0.2em] uppercase">
						Network Snapshot
					</h2>
					{#if data.leaderboard}
						<div class="mt-5 grid grid-cols-2 gap-4">
							<div class="border-outline-variant/10 bg-surface-container border p-4">
								<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Creators</div>
								<div class="font-headline text-on-surface mt-2 text-2xl font-bold">
									{formatInteger(data.leaderboard.overview.active_creators)}
								</div>
							</div>
							<div class="border-outline-variant/10 bg-surface-container border p-4">
								<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Predictors</div>
								<div class="font-headline text-primary mt-2 text-2xl font-bold">
									{formatInteger(data.leaderboard.overview.active_predictors)}
								</div>
							</div>
						</div>
						<div class="border-outline-variant/10 bg-surface-container mt-5 border p-4">
							<div class="text-outline text-[10px] tracking-[0.2em] uppercase">Resolved Share</div>
							<div class="font-headline text-primary-container mt-2 text-2xl font-bold">
								{formatPercent(
									String(
										data.leaderboard.overview.total_predictions === 0
											? 0
											: (data.leaderboard.overview.resolved_predictions /
													data.leaderboard.overview.total_predictions) *
													100
									)
								)}
							</div>
						</div>
					{:else}
						<p class="text-outline mt-5 text-sm">
							Snapshot metrics will appear here once the leaderboard response is available.
						</p>
					{/if}
				</section>
			</aside>
		</div>
	</div>
</main>
