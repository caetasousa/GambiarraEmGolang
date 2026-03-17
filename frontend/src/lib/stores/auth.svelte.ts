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

class AuthStore {
    token = $state<string | null>(null);
    user = $state<AuthUser | null>(null);

    isAuthenticated = $derived(!!this.token && !!this.user);

    constructor() {
        if (browser) {
            try {
                this.token = localStorage.getItem('auth_token');
                const userRaw = localStorage.getItem('auth_user');
                this.user = userRaw ? JSON.parse(userRaw) : null;
            } catch {
                this.token = null;
                this.user = null;
            }

            $effect.root(() => {
                $effect(() => {
                    if (this.token) {
                        localStorage.setItem('auth_token', this.token);
                    } else {
                        localStorage.removeItem('auth_token');
                    }
                    if (this.user) {
                        localStorage.setItem('auth_user', JSON.stringify(this.user));
                    } else {
                        localStorage.removeItem('auth_user');
                    }
                });
            });
        }
    }

    setAuth(token: string, userData: AuthUser) {
        this.token = token;
        this.user = userData;
    }

    logout() {
        this.token = null;
        this.user = null;
    }
}

export const auth = new AuthStore();

// Compat: exporta helpers usados em toda a codebase
export function getToken(): string | null {
    return auth.token;
}

export function setAuth(token: string, userData: AuthUser) {
    auth.setAuth(token, userData);
}

export function logout() {
    auth.logout();
}
