<script lang="ts">
    let {
        stats = { hoje: 8, pendentes: 3, confirmados: 12, receita: 450 },
    }: {
        stats?: { hoje: number; pendentes: number; confirmados: number; receita: number };
    } = $props();
</script>

<div class="animate-fade-in space-y-6">
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
            <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 mr-4">
                <span class="material-symbols-outlined">calendar_today</span>
            </div>
            <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Hoje</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{stats.hoje}</p>
            </div>
        </div>
        <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
            <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-yellow-100 dark:bg-yellow-900/30 text-yellow-600 dark:text-yellow-400 mr-4">
                <span class="material-symbols-outlined">pending</span>
            </div>
            <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Pendentes</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{stats.pendentes}</p>
            </div>
        </div>
        <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
            <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-green-100 dark:bg-green-900/30 text-green-600 dark:text-green-400 mr-4">
                <span class="material-symbols-outlined">check_circle</span>
            </div>
            <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Confirmados</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">{stats.confirmados}</p>
            </div>
        </div>
        <div class="bg-[hsl(var(--bs-card))] p-4 rounded-lg shadow-sm border border-border-light dark:border-border-dark flex items-center">
            <div class="w-12 h-12 rounded-full flex items-center justify-center flex-shrink-0 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 mr-4">
                <span class="material-symbols-outlined">attach_money</span>
            </div>
            <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Receita Est. (Hoje)</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white">R$ {stats.receita}</p>
            </div>
        </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 space-y-4">
            <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
                    <span class="material-symbols-outlined text-primary mr-2">schedule</span>
                    Atendimentos de Hoje
                </h3>
                <button class="text-sm text-primary font-medium hover:underline">Ver Todos</button>
            </div>

            <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border-l-4 border-l-warning border-y border-r border-gray-200 dark:border-gray-700 p-4 transition-all hover:shadow-md">
                <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center">
                    <div class="flex items-start space-x-4 mb-4 sm:mb-0">
                        <div class="bg-warning/10 text-warning p-2 rounded-full">
                            <span class="material-symbols-outlined">hourglass_top</span>
                        </div>
                        <div>
                            <div class="flex items-center">
                                <span class="text-lg font-bold text-gray-900 dark:text-white">09:00</span>
                                <span class="mx-2 text-gray-400">•</span>
                                <span class="text-sm text-gray-600 dark:text-gray-300 font-medium">Corte de Cabelo + Hidratação</span>
                            </div>
                            <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Cliente: <span class="font-medium text-gray-800 dark:text-gray-200">Maria Oliveira</span></p>
                            <div class="flex items-center mt-2 text-xs text-gray-400">
                                <span class="material-symbols-outlined text-[16px] mr-1">timer</span>
                                <span>60 min</span>
                                <span class="mx-2">•</span>
                                <span class="text-warning font-medium">Aguardando Confirmação</span>
                            </div>
                        </div>
                    </div>
                    <div class="flex space-x-2 w-full sm:w-auto">
                        <button class="flex-1 sm:flex-none px-3 py-1.5 border border-gray-300 dark:border-gray-600 text-gray-600 dark:text-gray-300 rounded text-sm hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">Cancelar</button>
                        <button class="flex-1 sm:flex-none px-3 py-1.5 bg-primary hover:bg-primary-hover text-white rounded text-sm shadow-sm transition-colors">Confirmar</button>
                    </div>
                </div>
            </div>
        </div>

        <div class="space-y-6">
            <div class="bg-[hsl(var(--bs-card))] rounded-lg shadow-sm border border-border-light dark:border-border-dark p-6">
                <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">Solicitações Recentes</h3>
                <div class="space-y-4">
                    <div class="flex items-start space-x-3 pb-4 border-b border-gray-100 dark:border-gray-800">
                        <div class="h-10 w-10 rounded-full bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center text-blue-600 dark:text-blue-400 shrink-0">
                            <span class="material-symbols-outlined text-sm">person</span>
                        </div>
                        <div>
                            <p class="text-sm font-medium text-gray-900 dark:text-white">Novo cliente se cadastrou</p>
                            <p class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Beatriz M. agendou para amanhã às 14h</p>
                            <p class="text-xs text-gray-400 mt-1">Há 15 min</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<style>
    .animate-fade-in { animation: fadeIn 0.3s ease-in-out; }
    @keyframes fadeIn {
        from { opacity: 0; transform: translateY(10px); }
        to { opacity: 1; transform: translateY(0); }
    }
</style>
