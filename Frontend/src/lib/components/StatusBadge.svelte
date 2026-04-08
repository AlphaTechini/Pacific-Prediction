<script lang="ts">
  interface Props {
    label: string;
    variant?: 'success' | 'danger' | 'warning' | 'info';
  }

  let { label, variant = 'info' }: Props = $props();

  const styles = {
    success: "bg-primary-container/10 text-primary-container border border-primary-container/20",
    danger: "bg-error/10 text-error border border-error/20",
    warning: "bg-tertiary-fixed-dim/10 text-tertiary-fixed-dim border border-tertiary-fixed-dim/20",
    info: "bg-surface-container-highest text-outline border border-outline-variant/30"
  };
</script>

<span class="text-[9px] px-2 py-0.5 font-bold uppercase tracking-tighter {styles[variant]}">
  {label}
</span>
