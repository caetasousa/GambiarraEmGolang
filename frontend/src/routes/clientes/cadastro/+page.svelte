<script lang="ts">
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { goto } from "$app/navigation";

    let nome = $state("");
    let email = $state("");
    let telefone = $state("");
    let senha = $state("");
    let confirmarSenha = $state("");
    let termos = $state(false);
    let errorMsg = $state("");
    let successMsg = $state("");
    let loading = $state(false);

    async function handleCadastro(e: Event) {
        e.preventDefault();
        errorMsg = "";
        successMsg = "";

        if (!nome || !email || !telefone || !senha || !confirmarSenha) {
            errorMsg = "Preencha todos os campos obrigatórios.";
            return;
        }
        if (senha !== confirmarSenha) {
            errorMsg = "As senhas não coincidem.";
            return;
        }
        if (!termos) {
            errorMsg = "Você precisa aceitar os Termos de Uso.";
            return;
        }

        loading = true;
        try {
            const res = await fetch("/api/v1/clientes", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ nome, email, telefone, senha }),
            });

            if (res.status === 409) {
                errorMsg = "Este e-mail já está cadastrado.";
                return;
            }
            if (!res.ok) {
                const body = await res.json().catch(() => ({}));
                errorMsg = body?.message ?? "Erro ao criar conta. Tente novamente.";
                return;
            }

            successMsg = "Conta criada com sucesso! Redirecionando para o login...";
            setTimeout(() => goto("/login"), 2000);
        } catch {
            errorMsg = "Erro de conexão. Verifique se o servidor está ativo.";
        } finally {
            loading = false;
        }
    }
</script>

<div class="font-montserrat bg-white dark:bg-zinc-950 text-slate-800 dark:text-gray-100 antialiased h-screen flex overflow-hidden">
    <!-- Left Side: Branding & Gradient -->
    <div class="hidden md:flex md:w-1/2 lg:w-3/5 relative bg-gradient-to-br from-pink-500 via-orange-400 to-orange-500 overflow-hidden">
        <!-- Decorative pattern -->
        <div
            class="absolute inset-0 opacity-[0.07] pointer-events-none"
            style="background-image: radial-gradient(circle at 2px 2px, white 1.5px, transparent 0); background-size: 40px 40px;"
        ></div>

        <!-- Blobs -->
        <div class="absolute top-10 left-10 w-72 h-72 bg-white/10 rounded-full blur-[90px] pointer-events-none"></div>
        <div class="absolute bottom-16 right-8 w-56 h-56 bg-orange-700/15 rounded-full blur-[70px] pointer-events-none"></div>

        <!-- Brand Logo -->
        <div class="absolute top-10 left-10 z-20 flex items-center gap-3">
            <div class="w-9 h-9 bg-white/20 backdrop-blur-sm rounded-xl flex items-center justify-center text-white font-bold text-base border border-white/30">
                B
            </div>
            <span class="font-bold text-xl tracking-wide text-white font-cormorant">BellaVita</span>
        </div>

        <!-- Salon image overlay -->
        <img
            alt="Experiência de Beleza"
            class="absolute inset-0 w-full h-full object-cover opacity-25 mix-blend-overlay"
            src="https://images.unsplash.com/photo-1562322140-8baeececf3df?auto=format&fit=crop&q=80&w=1200"
        />

        <!-- Bottom branding text -->
        <div class="relative z-10 w-full flex flex-col justify-end p-12 lg:p-16">
            <div class="mb-8">
                <div class="inline-flex items-center justify-center p-3 bg-white/15 backdrop-blur-md rounded-2xl mb-6 border border-white/25 shadow-lg">
                    <span class="material-icons text-white text-4xl">spa</span>
                </div>
                <h1 class="text-4xl lg:text-5xl font-bold font-cormorant text-white mb-4 leading-tight italic">
                    Comece sua jornada<br />de autocuidado.
                </h1>
                <p class="text-lg text-white/80 max-w-lg font-montserrat font-light leading-relaxed">
                    Cadastre-se para agendar serviços, acompanhar seu histórico
                    e receber ofertas exclusivas do BellaVita Salon.
                </p>
            </div>
            <!-- Progress dots -->
            <div class="flex space-x-2">
                <div class="w-2 h-1 bg-white/30 rounded-full"></div>
                <div class="w-8 h-1 bg-white rounded-full"></div>
                <div class="w-2 h-1 bg-white/30 rounded-full"></div>
            </div>
        </div>
    </div>

    <!-- Right Side: Registration Form -->
    <div class="w-full md:w-1/2 lg:w-2/5 flex flex-col bg-white dark:bg-zinc-950 relative overflow-y-auto">
        <div class="absolute top-6 right-6 z-20">
            <ThemeToggle />
        </div>

        <div class="w-full max-w-md mx-auto px-8 py-14 md:py-10 flex-1 flex flex-col justify-center">
            <!-- Mobile Brand -->
            <div class="md:hidden flex items-center mb-10 justify-center gap-2">
                <div class="w-8 h-8 bg-gradient-to-br from-orange-500 to-pink-500 rounded-xl flex items-center justify-center text-white font-bold text-sm">B</div>
                <span class="font-bold text-xl tracking-wide text-slate-800 dark:text-white font-cormorant">BellaVita</span>
            </div>

            <div class="mb-7 animate-fade-in-up-1">
                <h2 class="text-4xl font-bold font-cormorant text-slate-800 dark:text-white mb-2">
                    Crie sua conta
                </h2>
                <p class="text-gray-500 dark:text-gray-400 text-sm font-montserrat">
                    Preencha seus dados para começar a aproveitar.
                </p>
            </div>

            <form class="space-y-4 animate-fade-in-up-2" onsubmit={handleCadastro}>
                {#if errorMsg}
                    <div class="rounded-xl bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 px-4 py-3 text-sm font-montserrat text-red-700 dark:text-red-400">
                        {errorMsg}
                    </div>
                {/if}
                {#if successMsg}
                    <div class="rounded-xl bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 px-4 py-3 text-sm font-montserrat text-green-700 dark:text-green-400">
                        {successMsg}
                    </div>
                {/if}

                <!-- Nome -->
                <div>
                    <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300 mb-1.5" for="name">
                        Nome Completo
                    </label>
                    <div class="relative group">
                        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                            <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">badge</span>
                        </div>
                        <input
                            bind:value={nome}
                            class="block w-full pl-11 pr-4 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                            id="name"
                            placeholder="Ex: Maria Silva"
                            required
                            type="text"
                        />
                    </div>
                </div>

                <!-- Email -->
                <div>
                    <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300 mb-1.5" for="email">
                        E-mail
                    </label>
                    <div class="relative group">
                        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                            <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">mail</span>
                        </div>
                        <input
                            bind:value={email}
                            class="block w-full pl-11 pr-4 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                            id="email"
                            placeholder="seu@email.com"
                            required
                            type="email"
                        />
                    </div>
                </div>

                <!-- Telefone -->
                <div>
                    <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300 mb-1.5" for="phone">
                        Telefone / WhatsApp
                    </label>
                    <div class="relative group">
                        <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                            <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">smartphone</span>
                        </div>
                        <input
                            bind:value={telefone}
                            class="block w-full pl-11 pr-4 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                            id="phone"
                            placeholder="(00) 00000-0000"
                            required
                            type="tel"
                        />
                    </div>
                </div>

                <!-- Senha + Confirmar (grid) -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300 mb-1.5" for="password">
                            Senha
                        </label>
                        <div class="relative group">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                                <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">lock</span>
                            </div>
                            <input
                                bind:value={senha}
                                class="block w-full pl-11 pr-3 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                                id="password"
                                placeholder="••••••••"
                                required
                                type="password"
                            />
                        </div>
                    </div>
                    <div>
                        <label class="block text-sm font-medium font-montserrat text-slate-700 dark:text-gray-300 mb-1.5" for="confirm_password">
                            Confirmar
                        </label>
                        <div class="relative group">
                            <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none">
                                <span class="material-icons text-gray-400 text-[20px] group-focus-within:text-orange-500 transition-colors duration-200">lock_reset</span>
                            </div>
                            <input
                                bind:value={confirmarSenha}
                                class="block w-full pl-11 pr-3 py-3 border border-pink-200 dark:border-zinc-700 rounded-xl bg-pink-50/40 dark:bg-zinc-900 text-slate-800 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-orange-400/40 focus:border-orange-400 transition-all duration-200 text-sm font-montserrat"
                                id="confirm_password"
                                placeholder="••••••••"
                                required
                                type="password"
                            />
                        </div>
                    </div>
                </div>

                <!-- Termos -->
                <div class="flex items-center pt-1">
                    <input
                        bind:checked={termos}
                        class="h-4 w-4 text-orange-500 focus:ring-orange-400 border-pink-300 rounded cursor-pointer accent-orange-500"
                        id="terms"
                        required
                        type="checkbox"
                    />
                    <label class="ml-2.5 block text-sm font-montserrat text-gray-600 dark:text-gray-400 cursor-pointer" for="terms">
                        Eu li e aceito os <a class="text-orange-500 hover:underline font-medium" href="/">Termos de Uso</a>.
                    </label>
                </div>

                <!-- Submit -->
                <div class="pt-1">
                    <button
                        class="w-full flex justify-center items-center gap-2 py-3.5 px-4 rounded-xl text-sm font-semibold font-montserrat text-white bg-gradient-to-r from-orange-500 to-orange-600 shadow-lg shadow-orange-500/30 hover:shadow-orange-500/40 hover:scale-[1.02] hover:-translate-y-0.5 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-400 transition-all duration-200 active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed cursor-pointer"
                        type="submit"
                        disabled={loading}
                    >
                        {#if loading}
                            <span class="material-icons text-sm animate-spin">refresh</span>
                            Criando conta...
                        {:else}
                            Criar minha conta
                            <span class="material-icons text-sm">check_circle</span>
                        {/if}
                    </button>
                </div>
            </form>

            <!-- Divider & login link -->
            <div class="mt-7 animate-fade-in-up-3">
                <div class="relative">
                    <div class="absolute inset-0 flex items-center">
                        <div class="w-full border-t border-pink-100 dark:border-zinc-800"></div>
                    </div>
                    <div class="relative flex justify-center text-sm">
                        <span class="px-3 bg-white dark:bg-zinc-950 text-gray-400 font-medium font-montserrat">
                            Já possui uma conta?
                        </span>
                    </div>
                </div>
                <div class="mt-5">
                    <a
                        href="/login"
                        class="w-full flex justify-center py-3 px-4 border border-pink-200 dark:border-zinc-700 rounded-xl text-sm font-medium font-montserrat text-slate-700 dark:text-gray-200 bg-white dark:bg-zinc-900 hover:bg-pink-50/60 dark:hover:bg-zinc-800 hover:border-orange-300 focus:outline-none transition-all duration-200 cursor-pointer"
                    >
                        Fazer login agora
                    </a>
                </div>
            </div>

            <footer class="mt-8 text-center text-xs text-gray-400 dark:text-gray-500 font-montserrat">
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
