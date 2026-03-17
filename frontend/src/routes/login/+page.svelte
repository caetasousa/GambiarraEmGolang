<script lang="ts">
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { goto } from "$app/navigation";
    import { setAuth, type AuthUser, type UserRole } from "$lib/stores/auth.svelte";

    let email = $state("");
    let senha = $state("");
    let errorMsg = $state("");
    let loading = $state(false);

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
                cliente: "/clientes/meus-agendamentos",
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

<div class="font-montserrat bg-white dark:bg-zinc-950 text-slate-800 dark:text-gray-100 antialiased h-screen flex overflow-hidden">
    <!-- Left Side: Branding & Gradient -->
    <div class="hidden md:flex md:w-1/2 lg:w-3/5 relative bg-gradient-to-br from-orange-400 via-orange-500 to-pink-500 overflow-hidden">
        <!-- Decorative pattern -->
        <div
            class="absolute inset-0 opacity-[0.07] pointer-events-none"
            style="background-image: radial-gradient(circle at 2px 2px, white 1.5px, transparent 0); background-size: 40px 40px;"
        ></div>

        <!-- Decorative blobs -->
        <div class="absolute top-16 right-16 w-64 h-64 bg-white/10 rounded-full blur-[80px] pointer-events-none"></div>
        <div class="absolute bottom-24 left-8 w-48 h-48 bg-pink-700/20 rounded-full blur-[60px] pointer-events-none"></div>

        <!-- Brand Logo -->
        <div class="absolute top-10 left-10 z-20 flex items-center gap-3">
            <div class="w-9 h-9 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center text-white font-bold text-base border border-white/30">
                B
            </div>
            <span class="font-bold text-xl tracking-wide text-white font-cormorant">BellaVita</span>
        </div>

        <!-- Bottom branding text -->
        <div class="relative z-10 w-full flex flex-col justify-end p-12 lg:p-16">
            <div class="mb-8">
                <div class="inline-flex items-center justify-center p-3 bg-white/15 backdrop-blur-md rounded-2xl mb-6 border border-white/25 shadow-lg">
                    <span class="material-icons text-white text-4xl">spa</span>
                </div>
                <h1 class="text-4xl lg:text-5xl font-bold font-cormorant text-white mb-4 leading-tight italic">
                    Realce sua<br />beleza natural.
                </h1>
                <p class="text-lg text-white/80 max-w-lg font-montserrat font-light leading-relaxed">
                    Agende seus serviços favoritos, gerencie seus horários e
                    aproveite uma experiência exclusiva no BellaVita Salon.
                </p>
            </div>
            <!-- Progress dots -->
            <div class="flex space-x-2">
                <div class="w-8 h-1 bg-white rounded-full"></div>
                <div class="w-2 h-1 bg-white/30 rounded-full"></div>
                <div class="w-2 h-1 bg-white/30 rounded-full"></div>
            </div>
        </div>
    </div>

    <!-- Right Side: Login Form -->
    <div class="w-full md:w-1/2 lg:w-2/5 flex flex-col justify-center bg-white dark:bg-zinc-950 relative overflow-y-auto">
        <div class="absolute top-6 right-6 z-20">
            <ThemeToggle />
        </div>

        <div class="w-full max-w-md mx-auto px-8 py-12">
            <!-- Mobile Brand -->
            <div class="md:hidden flex items-center mb-10 justify-center gap-2">
                <div class="w-8 h-8 bg-gradient-to-br from-orange-500 to-pink-500 rounded-xl flex items-center justify-center text-white font-bold text-sm">B</div>
                <span class="font-bold text-xl tracking-wide text-slate-800 dark:text-white font-cormorant">BellaVita</span>
            </div>

            <div class="mb-8 animate-fade-in-up-1">
                <h2 class="text-4xl font-bold font-cormorant text-slate-800 dark:text-white mb-2">
                    Bem-vindo de volta
                </h2>
                <p class="text-gray-500 dark:text-gray-400 text-sm font-montserrat">
                    Insira seus dados para continuar.
                </p>
            </div>

            <form class="space-y-5 animate-fade-in-up-2" onsubmit={(e) => { e.preventDefault(); handleLogin(); }}>
                {#if errorMsg}
                    <div class="rounded-xl bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 px-4 py-3 text-sm font-montserrat text-red-700 dark:text-red-400">
                        {errorMsg}
                    </div>
                {/if}

                <!-- Email field -->
                <div>
                    <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300 mb-2" for="email">
                        E-mail
                    </label>
                    <div class="relative group">
                        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                            <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">email</span>
                        </div>
                        <input
                            bind:value={email}
                            class="block w-full pl-11 pr-4 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl leading-5 bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                            id="email"
                            placeholder="seu@email.com"
                            type="email"
                            required
                        />
                    </div>
                </div>

                <!-- Password field -->
                <div>
                    <div class="flex items-center justify-between mb-2">
                        <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300" for="password">
                            Senha
                        </label>
                        <a class="text-xs font-medium font-montserrat text-orange-500 hover:text-orange-600 dark:hover:text-orange-400 transition-colors cursor-pointer" href="/">
                            Esqueceu a senha?
                        </a>
                    </div>
                    <div class="relative group">
                        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                            <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">lock</span>
                        </div>
                        <input
                            bind:value={senha}
                            class="block w-full pl-11 pr-4 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl leading-5 bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                            id="password"
                            placeholder="••••••••"
                            type="password"
                            required
                        />
                    </div>
                </div>

                <!-- Submit button -->
                <div class="pt-2">
                    <button
                        class="w-full flex justify-center items-center gap-2 py-3.5 px-4 rounded-xl text-sm font-semibold font-montserrat text-white bg-gradient-to-r from-orange-500 to-orange-600 shadow-lg shadow-orange-500/30 hover:shadow-orange-500/40 hover:scale-[1.02] hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-400 transition-all duration-200 active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
                        type="submit"
                        disabled={loading}
                    >
                        {#if loading}
                            <span class="material-icons text-sm animate-spin">refresh</span>
                            Entrando...
                        {:else}
                            Acessar Conta
                            <span class="material-icons text-sm">arrow_forward</span>
                        {/if}
                    </button>
                </div>
            </form>

            <!-- Divider & signup link -->
            <div class="mt-8 animate-fade-in-up-3">
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <div class="w-full border-t border-pink-100 dark:border-zinc-800"></div>
                    </div>
                    <div class="relative flex justify-center text-sm">
                        <span class="px-3 bg-white dark:bg-zinc-950 text-gray-400 font-medium font-montserrat">
                            Novo por aqui?
                        </span>
                    </div>
                </div>
                <div class="mt-5">
                    <a
                        href="/clientes/cadastro"
                        class="w-full flex justify-center py-3 px-4 border border-pink-200 dark:border-zinc-700 rounded-xl text-sm font-medium font-montserrat text-slate-700 dark:text-gray-200 bg-white dark:bg-zinc-900 hover:bg-pink-50/60 dark:hover:bg-zinc-800 hover:border-orange-300 focus:outline-none transition-all duration-200 cursor-pointer"
                    >
                        Criar uma conta gratuita
                    </a>
                </div>
            </div>

            <footer class="mt-10 text-center text-xs text-gray-400 dark:text-gray-500 font-montserrat">
                <p class="mb-3">© 2024 BellaVita Salon. Todos os direitos reservados.</p>
                <div class="space-x-4">
                    <a class="hover:text-orange-500 transition-colors" href="/">Termos de Uso</a>
                    <a class="hover:text-orange-500 transition-colors" href="/">Política de Privacidade</a>
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
        background: #fbcfe8;
        border-radius: 4px;
    }
    :global(.dark) ::-webkit-scrollbar-thumb {
        background: #4b5563;
    }
</style>
