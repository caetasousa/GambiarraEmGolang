package main

import (
	"log"
	_ "meu-servico-agenda/docs"

	"meu-servico-agenda/internal/adapters/http/agendamento"
	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/prestador"
	"meu-servico-agenda/internal/infra/database"

	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API de Agendamentos
// @version 1.0
// @description API para gestão de clientes e serviços.
// @host localhost:8080
// @BasePath /api/v1
func main() {

	// 0. Conexão com o banco de dadoss (Infraestrutura)
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 1. Camada de Repositório (Infraestrutura)
	clienteRepo := repository.NovoClientePostgresRepositorio(db)
	prestadorRepo := repository.NewPrestadorPostgresRepository(db)
	catalogoRepo := repository.NovoCatalogoPostgresRepositorio(db)
	agendaDiariaRepo := repository.NovoAgendaDiariaPostgresRepository(db)
	agendamentoRepo := repository.NovoAgendamentoPostgresRepository(db)

	// 2. Camada de Aplicação (Serviços/Casos de Uso)
	cadastroCliente := service.NovoServiceCliente(clienteRepo)
	cadastroPrestador := service.NovaPrestadorService(prestadorRepo, catalogoRepo, agendaDiariaRepo)
	cadastraCatalogo := service.NovoCatalogoService(catalogoRepo)
	cadastraAgendamento := service.NovaAgendamentoService(prestadorRepo, agendamentoRepo, catalogoRepo, clienteRepo)

	// 3. Camada de Adaptador HTTP (Controller)
	clienteController := cliente.NovoClienteController(cadastroCliente)
	prestadorController := prestador.NovoPrestadorController(cadastroPrestador)
	catalogoController := catalogo.NovoCatalogoController(cadastraCatalogo)
	agendamentoController := agendamento.NovoAgendamentoController(cadastraAgendamento)

	// --- 4. Inicialização do Servidor Gin ---
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 5. Define as Rotas
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.GET("/clientes/:id", clienteController.GetCliente)

		apiV1.POST("/prestadores", prestadorController.PostPrestador)
		apiV1.GET("/prestadores/:id", prestadorController.GetPrestador)
		apiV1.PUT("/prestadores/:id/agenda", prestadorController.PutAgenda)

		apiV1.POST("/catalogos", catalogoController.PostCatalogo)
		apiV1.GET("/catalogos/:id", catalogoController.GetCatalogoPorID)
		apiV1.GET("/catalogos", catalogoController.GetCatalogos)
		apiV1.PUT("/catalogos/:id", catalogoController.Atualizar)
		apiV1.DELETE(("/catalogos/:id"), catalogoController.Deletar)

		apiV1.POST("/agendamentos", agendamentoController.PostAgendamento)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Pong"})
	})

	// 6. Inicia o Servidor
	log.Println("Servidor Gin rodando na porta 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}
