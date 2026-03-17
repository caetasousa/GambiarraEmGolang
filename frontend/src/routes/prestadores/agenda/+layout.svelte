<script lang="ts">
    import AuthGuard from "$lib/components/AuthGuard.svelte";
    import RoleGuard from "$lib/components/RoleGuard.svelte";
    let { children } = $props();
</script>

<AuthGuard>
    <RoleGuard allowedRoles={['prestador', 'admin']}>{@render children()}</RoleGuard>
</AuthGuard>
