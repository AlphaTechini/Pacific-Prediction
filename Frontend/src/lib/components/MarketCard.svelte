<script lang="ts">
  import StatusBadge from './StatusBadge.svelte';

  interface Market {
    id: string;
    title: string;
    subtitle: string;
    icon: string;
    badge: string;
    badgeVariant: 'success' | 'danger' | 'warning' | 'info';
    yesLabel: string;
    yesOdds: string;
    noLabel: string;
    noOdds: string;
    timer: string;
    oi: string;
    aiInsight: string;
    topBorder?: string;
  }

  let { market }: { market: Market } = $props();
</script>

<div class="bg-surface-container-low border-t-2 {market.topBorder ?? 'border-outline-variant'} p-5 relative group hover:bg-surface-container transition-colors">
  <div class="flex justify-between items-start mb-4">
    <div class="flex items-center gap-3">
      <div class="w-8 h-8 bg-surface-container-highest flex items-center justify-center border border-outline-variant/30">
        <span class="material-symbols-outlined text-primary-container text-lg">{market.icon}</span>
      </div>
      <h3 class="font-headline font-bold text-sm tracking-tight leading-tight max-w-[140px]">
        {market.title}<br/>
        <span class="text-[10px] text-outline font-normal uppercase tracking-wider">{market.subtitle}</span>
      </h3>
    </div>
    <StatusBadge label={market.badge} variant={market.badgeVariant} />
  </div>

  <div class="flex gap-2 mb-6">
    <button class="flex-1 bg-surface-container-highest py-3 text-xs font-bold border border-outline-variant/30 hover:border-primary-container transition-all group/btn">
      <span class="block text-[8px] text-outline group-hover/btn:text-primary-container mb-1">{market.yesLabel}</span>
      <span class="font-mono">{market.yesOdds}</span>
    </button>
    <button class="flex-1 bg-surface-container-highest py-3 text-xs font-bold border border-outline-variant/30 hover:border-primary-container transition-all group/btn">
      <span class="block text-[8px] text-outline group-hover/btn:text-primary-container mb-1">{market.noLabel}</span>
      <span class="font-mono">{market.noOdds}</span>
    </button>
  </div>

  <div class="flex justify-between items-center text-[10px] border-t border-outline-variant/10 pt-4">
    <div class="flex items-center gap-1 text-primary-container font-medium">
      <span class="material-symbols-outlined text-xs">timer</span>
      {market.timer}
    </div>
    <div class="text-outline uppercase">OI: {market.oi}</div>
  </div>

  <div class="mt-4 bg-surface-container-lowest p-2 text-[10px] italic text-outline border-l-2 border-primary-container/30">
    <span class="text-primary-container font-bold not-italic mr-1">AI INSIGHT:</span>{market.aiInsight}
  </div>
</div>
