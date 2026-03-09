import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';

export type UserRole = 'admin' | 'cliente' | 'prestador';

export interface AuthUser {
    id: string;
    nome: string;
    email: string;
    telefone: string;
    role: UserRole;
    cpf?: string;
    image_url?: string;
}

interface AuthState {
    token: string | null;
    user: AuthUser | null;
}

function loadFromStorage(): AuthState {
    if (!browser) return { token: null, user: null };
    try {
        const token = localStorage.getItem('auth_token');
        const userRaw = localStorage.getItem('auth_user');
        const user = userRaw ? JSON.parse(userRaw) : null;
        return { token, user };
    } catch {
        return { token: null, user: null };
    }
}

const initial = loadFromStorage();
const authStore = writable<AuthState>(initial);

if (browser) {
    authStore.subscribe((state) => {
        if (state.token) {
            localStorage.setItem('auth_token', state.token);
        } else {
            localStorage.removeItem('auth_token');
        }
        if (state.user) {
            localStorage.setItem('auth_user', JSON.stringify(state.user));
        } else {
            localStorage.removeItem('auth_user');
        }
    });
}

export const user = derived(authStore, ($auth) => $auth.user);
export const isAuthenticated = derived(authStore, ($auth) => !!$auth.token && !!$auth.user);

export function getToken(): string | null {
    let token: string | null = null;
    authStore.subscribe((s) => (token = s.token))();
    return token;
}

export function setAuth(token: string, userData: AuthUser) {
    authStore.set({ token, user: userData });
}

export function logout() {
    authStore.set({ token: null, user: null });
}
