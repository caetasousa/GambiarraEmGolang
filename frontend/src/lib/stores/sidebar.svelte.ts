const sidebarState = $state({ open: false });

export const sidebarOpen = {
    get value() { return sidebarState.open; }
};

export function openSidebar() {
    sidebarState.open = true;
}

export function closeSidebar() {
    sidebarState.open = false;
}
