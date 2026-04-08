<script lang="ts">
  interface Props {
    id?: string;
    type?: 'text' | 'number' | 'search';
    placeholder?: string;
    value?: string | number;
    class?: string;
    label?: string;
  }

  let { id, type = 'text', placeholder = '', value = $bindable(''), class: className = '', label }: Props = $props();
</script>

<div class="relative group {className}">
  {#if label}
    <span class="text-[10px] text-outline uppercase tracking-widest mb-2 block">{label}</span>
  {/if}
  <input
    {id}
    {type}
    {placeholder}
    bind:value
    class="w-full bg-surface-container-lowest border-none text-sm text-on-surface py-2.5 px-4 rounded-sm focus:ring-1 focus:ring-primary-container/30 outline-none placeholder:text-outline-variant/40"
  />
  <!-- Terminal cursor underline on focus -->
  <div class="absolute bottom-0 left-0 w-0 h-[2px] bg-primary-container transition-all duration-500 group-focus-within:w-full"></div>
</div>
