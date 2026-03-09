<script lang="ts">
    import ThemeToggle from "$lib/components/ThemeToggle.svelte";
    import { goto } from "$app/navigation";

    let nome = "";
    let email = "";
    let telefone = "";
    let senha = "";
    let confirmarSenha = "";
    let termos = false;
    let errorMsg = "";
    let successMsg = "";
    let loading = false;

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

<div
    class="font-body bg-background-light dark:bg-background-dark text-text-light dark:text-text-dark antialiased h-screen flex overflow-hidden transition-colors duration-200"
>
    <!-- Left Side: Branding & Image (Premium Auth Style) -->
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
            alt="Experiência de Beleza"
            class="absolute inset-0 w-full h-full object-cover opacity-60 mix-blend-overlay"
            src="https://images.unsplash.com/photo-1562322140-8baeececf3df?auto=format&fit=crop&q=80&w=1200"
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
                    Comece sua jornada<br />de autocuidado.
                </h1>
                <p class="text-lg text-gray-300 max-w-lg font-light">
                    Cadastre-se para agendar serviços, acompanhar seu histórico
                    e receber ofertas exclusivas do BellaVita Salon.
                </p>
            </div>
            <div class="flex space-x-2">
                <div
                    class="w-2 h-1 bg-gray-600 rounded-full hover:bg-gray-400 transition-colors"
                ></div>
                <div
                    class="w-8 h-1 bg-primary rounded-full transition-all duration-300"
                ></div>
                <div
                    class="w-2 h-1 bg-gray-600 rounded-full hover:bg-gray-400 transition-colors"
                ></div>
            </div>
        </div>
    </div>

    <!-- Right Side: Registration Form -->
    <div
        class="w-full md:w-1/2 lg:w-2/5 flex flex-col bg-surface-light dark:bg-surface-dark relative overflow-y-auto"
    >
        <div class="absolute top-6 right-6 z-20">
            <ThemeToggle />
        </div>

        <div
            class="w-full max-w-md mx-auto px-8 py-16 md:py-12 flex-1 flex flex-col justify-center"
        >
            <!-- Mobile Brand -->
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
                    Crie sua conta
                </h2>
                <p class="text-gray-500 dark:text-gray-400 text-sm">
                    Preencha seus dados para começar a aproveitar.
                </p>
            </div>

            <form class="space-y-5" on:submit={handleCadastro}>
                {#if errorMsg}
                    <div
                        class="rounded-lg bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 px-4 py-3 text-sm text-red-700 dark:text-red-400"
                    >
                        {errorMsg}
                    </div>
                {/if}
                {#if successMsg}
                    <div
                        class="rounded-lg bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 px-4 py-3 text-sm text-green-700 dark:text-green-400"
                    >
                        {successMsg}
                    </div>
                {/if}

                <div>
                    <label
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5"
                        for="name">Nome Completo</label
                    >
                    <div class="relative rounded-lg shadow-sm group">
                        <div
                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                        >
                            <span
                                class="material-icons text-gray-400 text-[20px] group-focus-within:text-blue-500 transition-colors"
                                >badge</span
                            >
                        </div>
                        <input
                            bind:value={nome}
                            class="block w-full pl-10 pr-3 py-3 border border-border-light dark:border-border-dark rounded-lg leading-5 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                            id="name"
                            placeholder="Ex: Maria Silva"
                            required
                            type="text"
                        />
                    </div>
                </div>

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
                                >mail</span
                            >
                        </div>
                        <input
                            bind:value={email}
                            class="block w-full pl-10 pr-3 py-3 border border-border-light dark:border-border-dark rounded-lg leading-5 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                            id="email"
                            placeholder="seu@email.com"
                            required
                            type="email"
                        />
                    </div>
                </div>

                <div>
                    <label
                        class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5"
                        for="phone">Telefone / WhatsApp</label
                    >
                    <div class="relative rounded-lg shadow-sm group">
                        <div
                            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                        >
                            <span
                                class="material-icons text-gray-400 text-[20px] group-focus-within:text-blue-500 transition-colors"
                                >smartphone</span
                            >
                        </div>
                        <input
                            bind:value={telefone}
                            class="block w-full pl-10 pr-3 py-3 border border-border-light dark:border-border-dark rounded-lg leading-5 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                            id="phone"
                            placeholder="(00) 00000-0000"
                            required
                            type="tel"
                        />
                    </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5"
                            for="password">Senha</label
                        >
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
                                required
                                type="password"
                            />
                        </div>
                    </div>
                    <div>
                        <label
                            class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1.5"
                            for="confirm_password">Confirmar</label
                        >
                        <div class="relative rounded-lg shadow-sm group">
                            <div
                                class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
                            >
                                <span
                                    class="material-icons text-gray-400 text-[20px] group-focus-within:text-blue-500 transition-colors"
                                    >lock_reset</span
                                >
                            </div>
                            <input
                                bind:value={confirmarSenha}
                                class="block w-full pl-10 pr-3 py-3 border border-border-light dark:border-border-dark rounded-lg leading-5 bg-gray-50 dark:bg-gray-800 text-gray-900 dark:text-gray-100 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 sm:text-sm"
                                id="confirm_password"
                                placeholder="••••••••"
                                required
                                type="password"
                            />
                        </div>
                    </div>
                </div>

                <div class="flex items-center pt-2">
                    <input
                        bind:checked={termos}
                        class="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded cursor-pointer"
                        id="terms"
                        required
                        type="checkbox"
                    />
                    <label
                        class="ml-2 block text-sm text-gray-600 dark:text-gray-400 cursor-pointer"
                        for="terms"
                    >
                        Eu li e aceito os <a
                            class="text-primary hover:underline font-medium"
                            href="/">Termos de Uso</a
                        >.
                    </label>
                </div>

                <div class="pt-2">
                    <button
                        class="w-full flex justify-center items-center py-3 px-4 border border-transparent rounded-lg shadow-md text-sm font-semibold text-white bg-primary hover:bg-primary-hover focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary transition-all duration-200 transform active:scale-[0.98] disabled:opacity-60 disabled:cursor-not-allowed"
                        type="submit"
                        disabled={loading}
                    >
                        {#if loading}
                            <span class="material-icons text-sm mr-2 animate-spin">refresh</span>
                            Criando conta...
                        {:else}
                            Criar minha conta
                            <span class="material-icons text-sm ml-2">check_circle</span>
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
                            >Já possui uma conta?</span
                        >
                    </div>
                </div>
                <div class="mt-6">
                    <a
                        href="/login"
                        class="w-full flex justify-center py-3 px-4 border border-border-light dark:border-border-dark rounded-lg shadow-sm text-sm font-medium text-gray-700 dark:text-gray-200 bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-all duration-200"
                    >
                        Fazer login agora
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
