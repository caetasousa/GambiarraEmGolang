import { browser } from '$app/environment';

const initialTheme = browser ? (localStorage.getItem('theme') || 'light') : 'light';

const themeState = $state({ value: initialTheme });

if (browser) {
    $effect.root(() => {
        $effect(() => {
            localStorage.setItem('theme', themeState.value);
            if (themeState.value === 'dark') {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }
        });
    });
}

export const theme = {
    get value() { return themeState.value; },
};

export function toggleTheme() {
    themeState.value = themeState.value === 'light' ? 'dark' : 'light';
}
