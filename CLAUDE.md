# AnnyGo — Instruções para o Agente

## Stack
- **Backend**: Go + Gin, arquitetura hexagonal, JWT, PostgreSQL, porta 8080
- **Frontend**: SvelteKit 5, TailwindCSS, TypeScript, porta 5173
- **Banco**: PostgreSQL 16 via Docker, migrações com Flyway (`backend/flyway/sql/`)
- **Docs**: Swagger em `http://localhost:8080/swagger/index.html`

## Estrutura do Projeto
```
annygo/
├── backend/              → Código Go (backend)
│   ├── cmd/api/main.go   → Entry point
│   ├── docs/             → Swagger gerado
│   ├── flyway/           → Migrações SQL
│   ├── go.mod / go.sum
│   └── internal/
│       ├── core/domain/          → Entidades: Prestador, Cliente, Catalogo, Agendamento
│       ├── core/application/
│       │   ├── service/          → Regras de negócio
│       │   ├── port/             → Interfaces de repositório
│       │   ├── input/ output/    → DTOs de entrada/saída
│       │   └── mapper/           → Conversão entre camadas
│       ├── adapters/
│       │   ├── http/             → Controllers Gin + request/response structs
│       │   ├── http/middleware/  → auth.go (JWT middleware)
│       │   └── repository/       → Implementações reais e fakes (para testes)
│       └── infra/
│           ├── auth/jwt.go       → Geração/validação JWT
│           └── database/database.go → Conexão PostgreSQL
├── frontend/             → Código SvelteKit (frontend)
│   └── src/
│       ├── lib/
│       │   ├── stores/auth.ts      → Estado de autenticação (persiste no localStorage)
│       │   ├── utils/api.ts        → fetchApi() — sempre usar para requisições autenticadas
│       │   ├── utils/validation.ts → Validações de formulário
│       │   ├── components/         → Componentes reutilizáveis (AuthGuard, Sidebar, etc.)
│       │   └── types/profile.ts    → Tipos TypeScript
│       └── routes/               → File-based routing SvelteKit
├── docker-compose.yaml   → PostgreSQL + pgAdmin + Flyway
└── .env                  → Variáveis de ambiente (nunca commitar)
```

## Autenticação
- `POST /api/v1/login` body: `{ email, senha }` → `{ token, role, user }`
- Roles: `admin`, `cliente`, `prestador`
- Token JWT, validade 24h, header: `Authorization: Bearer <token>`
- Sempre usar `fetchApi()` (de `src/lib/utils/api.ts`) — injeta o token automaticamente
- User ID vem de `$user.id` da store, nunca hardcodar IDs

## Rotas protegidas vs públicas
- **Públicas**: `/`, `/login`, `/listagem`, `/clientes/cadastro`, `/prestadores/cadastro`, `/catalogo`
- **Protegidas**: `/perfil`, `/clientes/agendamento`, `/prestadores/editar`, `/prestadores/agenda`
- Proteção via `+layout.svelte` com `AuthGuard`

## Comandos essenciais
```bash
# Backend (executar sempre dentro de backend/)
cd backend
go run cmd/api/main.go          # Iniciar backend
go test ./...                   # Rodar todos os testes
go build ./...                  # Compilar
swag init -g cmd/api/main.go    # Regenerar Swagger

# Frontend (executar sempre dentro de frontend/)
cd frontend
npm run dev                     # Iniciar frontend
npm run build                   # Build de produção
npm run check                   # Verificar tipos TypeScript

# Docker (executar na raiz do projeto)
docker compose up -d            # Subir PostgreSQL + pgAdmin + Flyway
docker network create local-network  # Criar rede (primeira vez)
```

## Clean Architecture — Regras obrigatórias

O backend segue **Clean Architecture** estritamente. Violar essas regras é proibido:

### Direção das dependências
```
adapters/http (Controllers)
    ↓
core/application/service (Use cases)
    ↓
core/domain (Entidades)
    ↑
core/application/port (Interfaces — repositórios)
    ↑
adapters/repository (Implementações)
```
- Dependências sempre apontam para dentro (em direção ao domínio)
- O domínio não conhece nada das camadas externas
- Services dependem de interfaces (ports), nunca de implementações concretas

### Regras por camada
- **domain/**: apenas entidades e regras de negócio puras. Sem imports de Gin, SQL, ou qualquer framework
- **application/service/**: orquestra casos de uso. Recebe e retorna tipos do domínio. Sem imports HTTP
- **application/port/**: apenas interfaces de repositório. Sem implementações
- **adapters/http/**: converte HTTP → domínio e domínio → HTTP. Sem lógica de negócio
- **adapters/repository/**: acesso ao banco. Implementa as interfaces de `port/`
- **infra/**: configurações técnicas (DB, JWT). Sem regras de negócio

### DTOs e conversão
- Dados externos entram via `input/` (nunca expor entidades do domínio diretamente)
- Dados saem via `output/`
- Conversão feita exclusivamente em `mapper/`

## Padrões de código
- Controllers ficam em `backend/internal/adapters/http/<entidade>/controller.go`
- Novos endpoints: seguir padrão existente (controller → service → repository)
- Testes usam fake repos (`backend/internal/adapters/repository/*_fake_repo.go`)
- Migrações SQL: criar novo arquivo `backend/flyway/sql/V<N>__descricao.sql`
- Frontend: componentes Svelte com TailwindCSS, sem CSS externo

## Fluxo obrigatório após alterações no backend
Sempre que qualquer arquivo em `backend/internal/`, `backend/cmd/` for modificado:
1. Executar `go test ./internal/...` **a partir de `backend/`** ao final
2. Apresentar relatório detalhado contendo:
   - Arquivos alterados (com caminho completo)
   - O que foi adicionado/modificado em cada arquivo
   - Resultado dos testes (quantos passaram, quantos falharam, quais falharam e por quê)
   - Se algum teste falhou, investigar e corrigir antes de encerrar

## Git e GitHub
- **Nunca executar `git push` sem pedir permissão explícita ao usuário antes**
- Commits locais podem ser criados normalmente quando solicitado
- Sempre confirmar com o usuário antes de enviar qualquer coisa ao repositório remoto
- **Commits granulares**: fazer commits pequenos e frequentes **durante** o desenvolvimento, um para cada mudança lógica (ex: nova feature, fix de bug, refactor). Não acumular múltiplas alterações para commitar tudo de uma vez no final

## Relatório de progresso
Antes de cada commit, gerar um relatório resumido de tudo que foi feito desde o último `git pull` (ou desde o início do branch). O relatório deve conter:
1. **Lista de commits** desde o último pull (hash curto, data, mensagem)
2. **Resumo por área** (backend, frontend, infra, docs) com o que foi adicionado/modificado
3. **Alterações pendentes** (não commitadas) — arquivos modificados e novos
4. **Status geral** — visão de alto nível do progresso do projeto

Exibir o relatório ao usuário sempre que solicitado ou antes de commits importantes.

## Variáveis de ambiente (`.env`)
- Nunca editar ou commitar o `.env` diretamente
- Variáveis usadas: DB connection string, JWT secret
