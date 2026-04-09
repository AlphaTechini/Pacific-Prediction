<script lang="ts">
	import { goto } from '$app/navigation';

	import TopNavBar from '$lib/components/TopNavBar.svelte';
	import Button from '$lib/components/Button.svelte';
</script>

<svelte:head>
	<title>Pacifica Pulse | Technical Prediction Markets</title>
</svelte:head>

<TopNavBar />

<!-- Hero Section -->
<main class="mx-auto max-w-7xl px-6 pt-32 pb-20">
	<div class="grid items-center gap-16 lg:grid-cols-2">
		<div class="space-y-8">
			<div
				class="bg-surface-container-high border-outline-variant/20 inline-flex items-center rounded-sm border px-3 py-1"
			>
				<span class="bg-primary-container mr-2 h-2 w-2 animate-pulse rounded-full"></span>
				<span class="text-primary text-[10px] font-bold tracking-[0.1em] uppercase"
					>Live Terminal Active</span
				>
			</div>
			<h1
				class="text-primary font-headline text-6xl leading-none font-bold tracking-tighter md:text-7xl"
			>
				Pacifica Pulse
			</h1>
			<p class="text-on-surface-variant max-w-lg text-2xl leading-relaxed font-light">
				Turn live market structure into conviction. Create and join fast prediction markets around
				price, candles, and funding momentum.
			</p>
			<div class="flex flex-wrap gap-4 pt-4">
				<Button variant="primary" class="px-8 py-4 text-lg" onclick={() => goto('/dashboard')}
					>Launch Terminal</Button
				>
				<Button variant="outline" class="px-8 py-4 text-lg" onclick={() => goto('/markets/create')}
					>Create Market</Button
				>
			</div>
		</div>

		<div class="group relative">
			<div
				class="bg-primary-container/5 group-hover:bg-primary-container/10 absolute -inset-4 rounded-full blur-3xl transition-all duration-700"
			></div>
			<div
				class="bg-surface-container-low border-outline-variant/20 terminal-glow relative overflow-hidden rounded-lg border p-4"
			>
				<div class="border-outline-variant/10 mb-4 flex items-center justify-between border-b pb-2">
					<div class="flex gap-2">
						<div class="bg-error/40 h-3 w-3 rounded-full"></div>
						<div class="bg-tertiary-fixed-dim/40 h-3 w-3 rounded-full"></div>
						<div class="bg-primary-container/40 h-3 w-3 rounded-full"></div>
					</div>
					<div class="text-outline font-mono text-[10px] tracking-widest uppercase">
						System_v4.2.0
					</div>
				</div>
				<div class="grid grid-cols-2 gap-4">
					<div class="space-y-4">
						<div class="bg-surface-container border-primary rounded-sm border-l-2 p-4">
							<div class="text-outline mb-1 text-[10px] uppercase">Funding Signal</div>
							<div class="text-primary-fixed font-mono text-xl tracking-tight">0.0124%</div>
						</div>
						<div class="bg-surface-container border-tertiary-fixed-dim rounded-sm border-l-2 p-4">
							<div class="text-outline mb-1 text-[10px] uppercase">Candle Bias</div>
							<div class="text-tertiary-fixed-dim font-mono text-xl tracking-tight">Bullish 5m</div>
						</div>
					</div>
					<div
						class="bg-surface-container-lowest border-outline-variant/10 flex items-center justify-center rounded-sm border p-4"
					>
						<img
							class="h-full w-full object-cover opacity-50 grayscale transition-all duration-500 hover:grayscale-0"
							alt="Abstract 3D visualization of glowing blue financial data"
							src="https://lh3.googleusercontent.com/aida-public/AB6AXuDBurO6IHXqymGQnO1WsxtgfS979MCFBrS5whziL8_GcOI5G3Hx7RqCVaB0ghUjvBcjNOipVUIKuAZ2-taCe2oI2uTWTN0gnrOsLql95XTrUFTsc1Ynzq57Apkp70q-sbRrlFCJhb451MazjaibWQysI2nIp2qpdYr27qpr3B6dmnXEIPErocIfs-vSjeXkVvo1dRAh-k-vF7J8ijMTxM4mKct6yWyhqzdGfJD456cl6btkJ-mOx9LXo9Rrwh5Zg53LqhnvGMc9Zw-1"
						/>
					</div>
				</div>
				<div
					class="bg-surface-container-highest/30 border-outline-variant/10 mt-4 rounded-sm border p-3"
				>
					<div class="text-primary mb-2 flex items-center justify-between font-mono text-[10px]">
						<span>LIVE MARKET BOARD</span>
						<span>PRICE • CANDLE • FUNDING</span>
					</div>
					<div class="bg-surface-container-lowest h-1 w-full overflow-hidden rounded-full">
						<div class="bg-primary-container h-full w-3/4"></div>
					</div>
				</div>
			</div>
		</div>
	</div>
</main>

<!-- Value Prop Section -->
<section class="bg-surface-container-low border-outline-variant/10 border-y px-6 py-24">
	<div class="mx-auto max-w-7xl">
		<div class="mb-16">
			<h2 class="text-primary font-headline mb-4 text-4xl font-bold tracking-tight">
				Prediction Markets for Technical Signals
			</h2>
			<p class="text-on-surface-variant max-w-2xl text-lg">
				Build conviction around the metrics traders already watch. Pacifica Pulse turns live price
				structure, candle direction, and funding pressure into clear yes-or-no markets.
			</p>
		</div>
		<div class="grid gap-1 md:grid-cols-3">
			<div
				class="bg-surface-container border-outline-variant/15 hover:bg-surface-container-high border p-8 transition-colors"
			>
				<span
					class="material-symbols-outlined text-primary mb-6 scale-125"
					style="font-variation-settings: 'FILL' 1;">trending_up</span
				>
				<h3 class="text-on-surface font-headline mb-3 text-xl font-bold">Price Action</h3>
				<p class="text-on-secondary-container text-sm leading-relaxed">
					Frame fast markets around threshold breaks and expiry targets without leaving the terminal
					flow.
				</p>
			</div>
			<div
				class="bg-surface-container border-outline-variant/15 hover:bg-surface-container-high border p-8 transition-colors"
			>
				<span
					class="material-symbols-outlined text-primary mb-6 scale-125"
					style="font-variation-settings: 'FILL' 1;">account_balance_wallet</span
				>
				<h3 class="text-on-surface font-headline mb-3 text-xl font-bold">Funding Rates</h3>
				<p class="text-on-secondary-container text-sm leading-relaxed">
					Forecast funding flips and pressure shifts using the same settlement source the backend
					resolves against.
				</p>
			</div>
			<div
				class="bg-surface-container border-outline-variant/15 hover:bg-surface-container-high border p-8 transition-colors"
			>
				<span
					class="material-symbols-outlined text-primary mb-6 scale-125"
					style="font-variation-settings: 'FILL' 1;">bar_chart</span
				>
				<h3 class="text-on-surface font-headline mb-3 text-xl font-bold">Candle Direction</h3>
				<p class="text-on-secondary-container text-sm leading-relaxed">
					Challenge the next close with bullish and bearish candle markets anchored to explicit
					mark-price intervals.
				</p>
			</div>
		</div>
	</div>
</section>

<!-- How It Works -->
<section class="mx-auto max-w-7xl px-6 py-24">
	<h2 class="font-headline mb-16 text-center text-3xl font-bold tracking-tight">The Market Loop</h2>
	<div class="relative grid grid-cols-1 gap-8 md:grid-cols-4">
		<div
			class="via-outline-variant/30 absolute top-12 left-0 -z-10 hidden h-[1px] w-full bg-gradient-to-r from-transparent to-transparent md:block"
		></div>
		{#each [{ icon: 'fort', title: 'Define The Market', desc: 'Choose a symbol, expiry, and rule built around price, candle, or funding behavior.' }, { icon: 'rocket_launch', title: 'Commit A Side', desc: 'Every new market opens with a real creator position instead of passive metadata.' }, { icon: 'monitoring', title: 'Track Live State', desc: 'Watch active markets, your balance, and open positions update through the app.' }, { icon: 'verified', title: 'Settle From Source', desc: 'The backend resolves against explicit Pacifica-derived data and applies virtual payouts.' }] as step}
			<div class="flex flex-col items-center space-y-4 text-center">
				<div
					class="bg-surface-container-highest border-primary/20 text-primary flex h-12 w-12 items-center justify-center rounded-sm border"
				>
					<span class="material-symbols-outlined">{step.icon}</span>
				</div>
				<div>
					<h4 class="text-on-surface font-headline font-bold">{step.title}</h4>
					<p class="text-on-surface-variant mt-2 px-4 text-xs">{step.desc}</p>
				</div>
			</div>
		{/each}
	</div>
</section>

<!-- Use Cases -->
<section
	class="bg-surface-container-lowest border-outline-variant/10 overflow-hidden border-y px-6 py-24"
>
	<div class="mx-auto grid max-w-7xl items-center gap-12 md:grid-cols-2">
		<div class="space-y-12">
			<div class="bg-surface-container-low border-primary border-l-4 p-8">
				<h3 class="font-headline mb-4 text-2xl font-bold">Fast Conviction Checks</h3>
				<p class="text-on-secondary-container leading-relaxed">
					Pressure-test your thesis in minutes. Use short-duration markets to challenge a level, a
					candle close, or a funding turn before the move is obvious.
				</p>
			</div>
			<div class="bg-surface-container-low border-primary border-l-4 p-8">
				<h3 class="font-headline mb-4 text-2xl font-bold">Creator-Led Markets</h3>
				<p class="text-on-secondary-container leading-relaxed">
					Every market starts with skin in the game. Creators define the setup, choose a side, and
					open the first position in one flow.
				</p>
			</div>
		</div>
		<div class="relative">
			<div
				class="bg-surface-container border-outline-variant/20 aspect-square overflow-hidden rounded-sm border"
			>
				<img
					class="h-full w-full object-cover opacity-60 mix-blend-overlay"
					alt="Sophisticated dark financial interface showing complex charts"
					src="https://lh3.googleusercontent.com/aida-public/AB6AXuAiNd1eiclUlS4Gzn2IdfQ-NITHlpO0nf3_bFiAodl2Cll_tJ2Xc9NetI_wl_fQ_dtaas1rMl5hU_JqBfoNpMVpJ4zeiw4iL5bRBGoWL-RjXBSWP1rY7-5qJbuIiwV4R9Lf0QjJy9XKPk5I0c-TExZp6-INZg4VXl8p6NuAV6p323EWU05MHt_GUwPgeYHFmf527uGEhBLxmhS-wUQ_im1mXoCHgMqC1ZstpCJ5ds_dW8aXUEBwi0htH6uoSvOHukWcSDXEmUJmT-Tz"
				/>
				<div
					class="from-surface-container-lowest absolute inset-0 bg-gradient-to-t via-transparent to-transparent"
				></div>
				<div class="absolute right-8 bottom-8 left-8">
					<div class="mb-2 flex items-center gap-2">
						<span
							class="bg-primary/10 text-primary px-2 py-0.5 text-[10px] font-bold tracking-tighter"
							>LIVE GAME FLOW</span
						>
					</div>
					<div class="font-headline text-2xl font-bold text-white">
						Markets, positions, balances, and results tied together in one loop.
					</div>
				</div>
			</div>
		</div>
	</div>
</section>

<!-- Live Surface -->
<section class="mx-auto max-w-7xl px-6 py-24">
	<div class="flex flex-col items-start gap-12 md:flex-row">
		<div class="sticky top-24 md:w-1/3">
			<div class="from-primary/30 rounded-sm bg-gradient-to-br to-transparent p-1">
				<div class="bg-surface-container-high rounded-sm p-6">
					<div class="mb-6 flex items-center gap-3">
						<span class="material-symbols-outlined text-primary">hub</span>
						<span class="text-on-surface text-sm font-bold tracking-widest uppercase"
							>Operator Surface</span
						>
					</div>
					<p class="text-on-surface-variant mb-6 text-sm leading-relaxed">
						Step into the live product loop: build a market, place a side, watch settlement, and
						climb the leaderboard from actual activity.
					</p>
					<div class="space-y-4">
						{#each [{ label: 'MARKET TYPES', value: '3 LIVE', color: 'text-primary' }, { label: 'SETTLEMENT PATH', value: 'SOURCE-DRIVEN', color: '' }, { label: 'PLAYER MODE', value: 'GUEST SESSION', color: '' }] as stat}
							<div class="flex justify-between font-mono text-[10px]">
								<span class="text-outline">{stat.label}</span>
								<span class={stat.color || 'text-on-surface'}>{stat.value}</span>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>
		<div class="grid grid-cols-1 gap-6 md:w-2/3">
			{#each [{ dot: 'bg-error', label: 'Live Price Setup', title: 'Challenge a level before expiry', desc: 'Create a threshold market around the exact price boundary you care about.', pct: 'YES / NO', pool: 'CREATOR OPENING STAKE' }, { dot: 'bg-primary', label: 'Live Funding Setup', title: 'Track the next funding turn', desc: 'Frame the next funding epoch with a clean sign or threshold call.', pct: 'SETTLED FROM SOURCE', pool: 'BACKEND VERIFIED' }] as card}
				<div
					class="group bg-surface-container-low border-outline-variant/15 hover:border-primary/40 flex cursor-crosshair flex-col items-center justify-between gap-6 border p-6 transition-all md:flex-row"
				>
					<div class="flex-1">
						<div class="mb-2 flex items-center gap-2">
							<span class="h-2 w-2 rounded-full {card.dot}"></span>
							<span
								class="font-mono text-[10px] {card.dot === 'bg-error'
									? 'text-error'
									: 'text-primary'} uppercase">{card.label}</span
							>
						</div>
						<h4 class="font-headline text-xl font-bold">{card.title}</h4>
						<p class="text-on-secondary-container mt-1 text-sm">{card.desc}</p>
					</div>
					<div class="flex flex-col items-end">
						<div class="text-primary font-mono text-3xl font-bold">{card.pct}</div>
						<div class="text-outline mt-1 text-[10px] uppercase">Pool: {card.pool}</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</section>

<!-- Footer -->
<footer class="border-t border-[#3b494b]/20 bg-[#111417]">
	<div
		class="mx-auto grid max-w-7xl grid-cols-1 gap-8 px-6 py-12 text-sm tracking-wide md:grid-cols-4"
	>
		<div class="col-span-1">
			<div class="font-headline mb-6 text-lg font-bold text-[#dbfcff]">Pacifica Pulse</div>
			<div
				class="bg-surface-container-low border-primary/10 inline-flex items-center rounded-sm border px-3 py-1"
			>
				<span class="bg-primary mr-2 h-2 w-2 rounded-full"></span>
				<span class="text-primary text-[10px] font-bold tracking-tighter uppercase"
					>System Status: Optimal</span
				>
			</div>
		</div>
		<div class="space-y-4">
			<h5 class="font-headline text-[10px] font-bold tracking-widest text-[#00F0FF] uppercase">
				Platform
			</h5>
			<ul class="space-y-2">
				<li>
					<a class="text-slate-500 underline transition-all hover:text-[#00F0FF]" href="/dashboard"
						>Terminal</a
					>
				</li>
				<li>
					<a
						class="text-slate-500 underline transition-all hover:text-[#00F0FF]"
						href="/markets/create">Create Market</a
					>
				</li>
			</ul>
		</div>
		<div class="space-y-4">
			<h5 class="font-headline text-[10px] font-bold tracking-widest text-[#00F0FF] uppercase">
				Legal
			</h5>
			<ul class="space-y-2">
				<li>
					<a class="text-slate-500 underline transition-all hover:text-[#00F0FF]" href="/"
						>Risk Disclosure</a
					>
				</li>
				<li>
					<a class="text-slate-500 underline transition-all hover:text-[#00F0FF]" href="/"
						>Privacy Policy</a
					>
				</li>
			</ul>
		</div>
		<div class="space-y-4">
			<h5 class="font-headline text-[10px] font-bold tracking-widest text-[#00F0FF] uppercase">
				Updates
			</h5>
			<ul class="space-y-2">
				<li>
					<a
						class="text-slate-500 underline transition-all hover:text-[#00F0FF]"
						href="/leaderboard">Leaderboard</a
					>
				</li>
			</ul>
			<div class="pt-4 text-xs leading-relaxed text-slate-600">
				© 2024 Pacifica Pulse Intelligence. All rights reserved.
			</div>
		</div>
	</div>
</footer>
