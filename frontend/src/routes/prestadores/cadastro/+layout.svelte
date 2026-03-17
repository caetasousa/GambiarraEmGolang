<script lang="ts">
    import AuthGuard from "$lib/components/AuthGuard.svelte";
    import RoleGuard from "$lib/components/RoleGuard.svelte";
    let { children } = $props();
</script>

<AuthGuard>
    <RoleGuard allowedRoles={['admin']}>{@render children()}</RoleGuard>
</AuthGuard>
