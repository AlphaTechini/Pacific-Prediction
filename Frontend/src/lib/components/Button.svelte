<script lang="ts">
	interface Props {
		onclick?: () => void;
		children?: any;
		class?: string;
		variant?: 'primary' | 'outline' | 'ghost';
	}

	let { onclick, children, class: className = '', variant = 'primary' }: Props = $props();

	const base = 'font-bold rounded-sm transition-all active:scale-95';
	const variants = {
		primary: 'bg-[#00bfd1] text-white hover:bg-[#00cddd]',
		outline: 'border border-outline-variant/30 text-primary hover:bg-surface-container-high',
		ghost: 'text-slate-400 hover:text-[#dbfcff]'
	};
</script>

<button {onclick} class="{base} {variants[variant]} {className}">
	{@render children?.()}
</button>
