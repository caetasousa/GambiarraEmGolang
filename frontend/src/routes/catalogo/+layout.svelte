<script lang="ts">
    import AuthGuard from "$lib/components/AuthGuard.svelte";
    import RoleGuard from "$lib/components/RoleGuard.svelte";
</script>

<AuthGuard>
    <RoleGuard allowedRoles={['admin']}>
        <slot />
    </RoleGuard>
</AuthGuard>
