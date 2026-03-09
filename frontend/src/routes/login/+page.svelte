<script lang="ts">
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { goto } from "$app/navigation";
    import { setAuth, type AuthUser, type UserRole } from "$lib/stores/auth";

    let email = "";
    let senha = "";
    let errorMsg = "";
    let loading = false;

    async function handleLogin() {
        errorMsg = "";
        if (!email || !senha) {
            errorMsg = "Preencha e-mail e senha.";
            return;
        }
        loading = true;
        try {
            const res = await fetch("/api/v1/login", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, senha }),
            });

            if (res.status === 401) {
                errorMsg = "E-mail ou senha incorretos.";
                return;
            }
            if (!res.ok) {
                errorMsg = "Erro ao fazer login. Tente novamente.";
                return;
            }

            const data = await res.json();
            const userData: AuthUser = {
                id: data.user.id,
                nome: data.user.nome,
                email: data.user.email,
                telefone: data.user.telefone ?? "",
                role: data.role as UserRole,
                cpf: data.user.cpf,
                image_url: data.user.image_url,
            };
            setAuth(data.token, userData);

            const redirectMap: Record<UserRole, string> = {
                admin: "/prestadores/editar",
                cliente: "/clientes/agendamento",
                prestador: "/perfil",
            };
            goto(redirectMap[userData.role] ?? "/");
        } catch {
            errorMsg = "Erro de conexão. Verifique se o servidor está ativo.";
        } finally {
            loading = false;
        }
    }
</script>

<div
    class="font-body bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <div class="hidden md:flex md:w-1/2 lg:w-3/5 relative bg-gray-900">
        <!-- Brand Logo on Desktop -->
        <div class="absolute top-12 left-12 z-20 flex items-center">
            <span class="material-icons text-primary text-3xl mr-2">spa</span>
            <span
                class="font-bold text-2xl tracking-tight text-white font-display"
                >BellaVita</span
            >
        </div>

        <img
            alt="Interior do Salão"
            class="absolute inset-0 w-full h-full object-cover opacity-60 mix-blend-overlay"
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuDhlQ2Nki0e3bAK5mzrTEtFcf15J6aqqWFX_lxXgqcG-0SWM9ODly4KEJbNeq798MTrB8n_vNiSAlfU_Z4MW360wNrGbEsWpnOAbpyrof-xl_uhZEjt87EDQ4QpboPhvqCK8wPZdNrAsTpuxa2OlYFt0V6ZGweu-3QUV4iIADJdeS29wuYs2yGdjxLHDEkmmH8Yi6gZGB9RNUO7U3FhWNlRl27MeBhXYaUMmh6buEMqq7BazskeYraYgNjuG6ILXsyXo3jzIAtmbqQ"
        />
        <div
            class="absolute inset-0 bg-gradient-to-t from-background-dark via-background-dark/50 to-transparent"
        ></div>
        <div
            class="relative z-10 w-full flex flex-col justify-end p-12 lg:p-16"
        >
            <div class="mb-8">
                <div
                    class="inline-flex items-center justify-center p-3 bg-primary/20 backdrop-blur-md rounded-xl mb-6 border border-primary/30 shadow-lg"
                >
                    <span class="material-icons text-primary text-4xl">spa</span
                    >
                </div>
                <h1
                    class="text-4xl lg:text-5xl font-bold text-white font-display mb-4 tracking-tight leading-tight"
                >
                    Realce sua<br />beleza natural.
                </h1>
                <p class="text-lg text-gray-300 max-w-lg font-light">
                    Agende seus serviços favoritos, gerencie seus horários e
                    aproveite uma experiência exclusiva no BellaVita Salon.
                </p>
            </div>
            <div class="flex space-x-2">
                <div
                    class="w-8 h-1 bg-primary rounded-full transition-all duration-300"
                ></div>
                <div
                    class="w-2 h-1 bg-gray-600 rounded-full hover:bg-gray-400 transition-colors"
                ></div>
                <div
                    class="w-2 h-1 bg-gray-600 rounded-full hover:bg-gray-400 transition-colors"
                ></div>
            </div>
        </div>
    </div>
    <div
        class="w-full md:w-1/2 lg:w-2/5 flex flex-col justify-center bg-surface-light dark:bg-surface-dark relative overflow-y-auto"
    >
        <div class="absolute top-6 right-6 z-20">
            <ThemeToggle />
        </div>
        <div class="w-full max-w-md mx-auto px-8 py-12">
            <div class="md:hidden flex items-center mb-10 justify-center">
                <span class="material-icons text-primary text-3xl mr-2"
                    >spa</span
                >
                <span
                    class="font-bold text-2xl tracking-tight text-gray-900 dark:text-white font-display"
                    >BellaVita</span
                >
            </div>
            <div class="mb-8">
                <h2
                    class="text-3xl font-bold text-gray-900 dark:text-white font-display mb-2"
                >
                    Login
                </h2>
                <p class="text-gray-500 dark:text-gray-400 text-sm">
                    Bem-vindo de volta! Insira seus dados para continuar.
                </p>
            </div>
            <form class="space-y-5" on:submit|preventDefault={handleLogin}>
                {#if errorMsg}
                    <div
                        class="rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 px-4 py-3 text-sm text-red-700 dark:text-red-400"
                    >
                        {errorMsg}
                    </div>
                {/if}
                <div>
                    <label
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5"
                        for="email">E-mail</label
                    >
                    <div class="relative rounded-lg shadow-sm group">
                        <div
                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                        >
                            <span
                                class="material-icons text-gray-400 text-[20px] group-focus-within:text-blue-500 transition-colors"
                                >email</span
                            >
                        </div>
                        <input
                            bind:value={email}
                            class="block w-full pl-10 pr-3 py-3 border border-border-light dark:border-border-dark rounded-lg leading-5 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                            id="email"
                            placeholder="seu@email.com"
                            type="email"
                            required
                        />
                    </div>
                </div>
                <div>
                    <div class="flex items-center justify-between mb-1.5">
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                            for="password">Senha</label
                        >
                        <a
                            class="text-xs font-medium text-blue-600 hover:text-blue-500 dark:text-blue-400 dark:hover:text-blue-300 transition-colors"
                            href="/">Esqueceu a senha?</a
                        >
                    </div>
                    <div class="relative rounded-lg shadow-sm group">
                        <div
                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                        >
                            <span
                                class="material-icons text-gray-400 text-[20px] group-focus-within:text-blue-500 transition-colors"
                                >lock</span
                            >
                        </div>
                        <input
                            bind:value={senha}
                            class="block w-full pl-10 pr-3 py-3 border border-border-light dark:border-border-dark rounded-lg leading-5 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                            id="password"
                            placeholder="••••••••"
                            type="password"
                            required
                        />
                    </div>
                </div>
                <div class="pt-2">
                    <button
                        class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-md text-sm font-semibold text-white bg-primary hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-all duration-200 transform active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed"
                        type="submit"
                        disabled={loading}
                    >
                        {#if loading}
                            <span class="material-icons text-sm mr-2 animate-spin">refresh</span>
                            Entrando...
                        {:else}
                            Acessar Conta
                            <span class="material-icons text-sm ml-2">arrow_forward</span>
                        {/if}
                    </button>
                </div>
            </form>
            <div class="mt-8">
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <div
                            class="w-full border-t border-border-light dark:border-border-dark"
                        ></div>
                    </div>
                    <div class="relative flex justify-center text-sm">
                        <span
                            class="px-3 bg-surface-light dark:bg-surface-dark text-gray-500 dark:text-gray-400 font-medium"
                            >Novo por aqui?</span
                        >
                    </div>
                </div>
                <div class="mt-6">
                    <a
                        href="/clientes/cadastro"
                        class="w-full flex justify-center py-3 px-4 border border-border-light dark:border-border-dark rounded-lg shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-all duration-200"
                    >
                        Criar uma conta gratuita
                    </a>
                </div>
            </div>
            <footer
                class="mt-12 text-center text-xs text-gray-500 dark:text-gray-500"
            >
                <p class="mb-3">
                    © 2024 BellaVita Salon. Todos os direitos reservados.
                </p>
                <div class="space-x-4">
                    <a
                        class="hover:text-primary transition-colors underline decoration-transparent hover:decoration-primary underline-offset-2"
                        href="/">Termos de Uso</a
                    >
                    <a
                        class="hover:text-primary transition-colors underline decoration-transparent hover:decoration-primary underline-offset-2"
                        href="/">Política de Privacidade</a
                    >
                </div>
            </footer>
        </div>
    </div>
</div>

<style>
    ::-webkit-scrollbar {
        width: 8px;
        height: 8px;
    }
    ::-webkit-scrollbar-track {
        background: transparent;
    }
    ::-webkit-scrollbar-thumb {
        background: #d1d5db;
        border-radius: 4px;
    }
    :global(.dark) ::-webkit-scrollbar-thumb {
        background: #4b5563;
    }
</style>
