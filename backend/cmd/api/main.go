package main

import (
	"log"
	"os"

	_ "meu-servico-agenda/docs"

	"meu-servico-agenda/internal/adapters/http/agendamento"
	"meu-servico-agenda/internal/adapters/http/auth"
	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/middleware"
	"meu-servico-agenda/internal/adapters/http/prestador"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"
	"meu-servico-agenda/internal/infra/database"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Agendamentos
// @version 1.0
// @description API para gestão de clientes e serviços com autenticação JWT.
// @description
// @description Credenciais do Administrador (para testes):
// @description Email: admin@admin.com
// @description Senha: admin123
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Digite "Bearer " seguido do token JWT recebido no login
func main() {

	// 0. Configurações JWT
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my-super-secret-key-change-in-production" // Default para desenvolvimento
		log.Println("AVISO: Usando JWT_SECRET padrão. Configure JWT_SECRET em produção!")
	}
	jwtExpirationHrs := 24 // Token expira em 24 horas

	// 1. Conexão com o banco de dados (Infraestrutura)
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. Camada de Repositório (Infraestrutura)
	clienteRepo := repository.NovoClientePostgresRepositorio(db)
	prestadorRepo := repository.NewPrestadorPostgresRepository(db)
	catalogoRepo := repository.NovoCatalogoPostgresRepositorio(db)
	agendaDiariaRepo := repository.NovoAgendaDiariaPostgresRepository(db)
	agendamentoRepo := repository.NovoAgendamentoPostgresRepository(db)

	// 3. Camada de Aplicação (Serviços/Casos de Uso)
	authService := service.NovoAuthService(jwtSecret, jwtExpirationHrs, clienteRepo, prestadorRepo)
	cadastroCliente := service.NovoServiceCliente(clienteRepo, prestadorRepo)
	cadastroPrestador := service.NovaPrestadorService(prestadorRepo, catalogoRepo, agendaDiariaRepo, authService, clienteRepo)
	cadastraCatalogo := service.NovoCatalogoService(catalogoRepo)
	cadastraAgendamento := service.NovaAgendamentoService(prestadorRepo, agendamentoRepo, catalogoRepo, clienteRepo)

	// 4. Camada de Adaptador HTTP (Controller)
	authController := auth.NovoAuthController(authService)
	clienteController := cliente.NovoClienteController(cadastroCliente)
	prestadorController := prestador.NovoPrestadorController(cadastroPrestador)
	catalogoController := catalogo.NovoCatalogoController(cadastraCatalogo)
	agendamentoController := agendamento.NovoAgendamentoController(cadastraAgendamento)

	// 5. Inicialização do Servidor Gin
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 6. Define as Rotas PÚBLICAS (sem autenticação)
	public := router.Group("/api/v1")
	{
		// Autenticação
		public.POST("/login", authController.Login)

		// Cadastro de clientes (qualquer um pode se cadastrar)
		public.POST("/clientes", clienteController.PostCliente)

		// Cadastro de prestadores (qualquer um pode se cadastrar)
		public.POST("/prestadores", prestadorController.PostPrestador)


		// Catálogos - Visualização pública
		public.GET("/catalogos", catalogoController.GetCatalogos)
		public.GET("/catalogos/:id", catalogoController.GetCatalogoPorID)
	}

	// 7. Define as Rotas PROTEGIDAS (requer autenticação JWT)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(jwtSecret)) // Aplica middleware de autenticação
	{
		// ========== CLIENTES ==========
		// GET /clientes/:id - Cliente próprio OU Admin
		protected.GET("/clientes/:id",
			middleware.RequireOwnerOrAdmin(),
			clienteController.GetCliente,
		)

		// ========== PRESTADORES ==========
		// GET /prestadores - Todos podem ver (Cliente, Prestador, Admin)
		protected.GET("/prestadores", prestadorController.GetPrestadores)

		// GET /prestadores/disponiveis - Todos podem ver
		protected.GET("/prestadores/disponiveis", prestadorController.GetPrestadoresPorData)

		// GET /prestadores/:id - Todos podem ver
		protected.GET("/prestadores/:id", prestadorController.GetPrestador)

		// PUT /prestadores/:id - Prestador próprio OU Admin
		protected.PUT("/prestadores/:id",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			prestadorController.UpdatePrestador,
		)

		// PUT /prestadores/:id/agenda - Prestador próprio OU Admin
		protected.PUT("/prestadores/:id/agenda",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			prestadorController.PutAgenda,
		)

		// DELETE /prestadores/:id/agenda - Prestador próprio OU Admin
		protected.DELETE("/prestadores/:id/agenda",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			prestadorController.DeleteAgenda,
		)

		// PUT /prestadores/:id/inativar - APENAS Admin
		protected.PUT("/prestadores/:id/inativar",
			middleware.RequireRole(domain.RoleAdmin),
			prestadorController.InativarPrestador,
		)

		// PUT /prestadores/:id/ativar - APENAS Admin
		protected.PUT("/prestadores/:id/ativar",
			middleware.RequireRole(domain.RoleAdmin),
			prestadorController.AtivarPrestador,
		)

		// ========== CATÁLOGOS ==========
		// POST /catalogos - APENAS Admin
		protected.POST("/catalogos",
			middleware.RequireRole(domain.RoleAdmin),
			catalogoController.PostCatalogo,
		)

		// PUT /catalogos/:id - APENAS Admin
		protected.PUT("/catalogos/:id",
			middleware.RequireRole(domain.RoleAdmin),
			catalogoController.Atualizar,
		)

		// DELETE /catalogos/:id - APENAS Admin
		protected.DELETE("/catalogos/:id",
			middleware.RequireRole(domain.RoleAdmin),
			catalogoController.Deletar,
		)

		// ========== AGENDAMENTOS ==========
		// POST /agendamentos - Cliente OU Admin
		protected.POST("/agendamentos",
			middleware.RequireRole(domain.RoleCliente, domain.RoleAdmin),
			agendamentoController.PostAgendamento,
		)

		// GET /agendamentos/cliente/:id - Cliente próprio OU Admin
		protected.GET("/agendamentos/cliente/:id",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			agendamentoController.GetAgendamentoClienteData,
		)

		// GET /agendamentos/cliente/:id/periodo - Cliente próprio OU Admin
		protected.GET("/agendamentos/cliente/:id/periodo",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			agendamentoController.GetAgendamentosClientePeriodo,
		)

		// GET /agendamentos/prestador/:id - Prestador próprio OU Admin
		protected.GET("/agendamentos/prestador/:id",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			agendamentoController.GetAgendamentoPrestadorData,
		)

		// GET /agendamentos/prestador/:id/periodo - Prestador próprio OU Admin
		protected.GET("/agendamentos/prestador/:id/periodo",
			middleware.RequireOwnerOrRole(domain.RoleAdmin),
			agendamentoController.GetAgendamentosPrestadorPeriodo,
		)

		// PUT /agendamentos/:id/confirmar - Prestador OU Admin
		protected.PUT("/agendamentos/:id/confirmar",
			middleware.RequireRole(domain.RolePrestador, domain.RoleAdmin),
			agendamentoController.ConfirmarAgendamento,
		)

		// PUT /agendamentos/:id/cancelar - Prestador, Cliente OU Admin
		protected.PUT("/agendamentos/:id/cancelar",
			middleware.RequireRole(domain.RolePrestador, domain.RoleCliente, domain.RoleAdmin),
			agendamentoController.CancelarAgendamento,
		)
	}

	router.GET("/health", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy"})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// 6. Inicia o Servidor
	log.Println("Servidor Gin rodando na porta 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
