package main

import (
	"log"
	_ "meu-servico-agenda/docs"

	"meu-servico-agenda/internal/adapters/http/agenda"
	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/prestador"

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
	// 1. Camada de Repositório (Infraestrutura)
	clienteRepo := repository.NewFakeClienteRepositorio()
	prestadorRepo := repository.NovoFakePrestadorRepositorio()
	catalogoRepo := repository.NovoCatalogoFakeRepo()
	agendaDiariaRepo := repository.NovoAgendaFakeRepo()

	// 2. Camada de Aplicação (Serviços/Casos de Uso)
	cadastroCliente := service.NovoServiceCliente(clienteRepo)
	cadastroPrestador := service.NovoPrestadorService(prestadorRepo)
	cadastraCatalogo := service.NovoCatalogoService(catalogoRepo)
	cadastraAgendaDiaria := service.NovaServiceAgendaDiaria(agendaDiariaRepo)

	// 3. Camada de Adaptador HTTP (Controller)
	clienteController := cliente.NovoClienteController(cadastroCliente)
	prestadorController := prestador.NovoPrestadorController(cadastroPrestador, catalogoRepo)
	catalogoController := catalogo.NovoCatalogoController(cadastraCatalogo)
	agendaDiariaController := agenda.NovaAgendaDiariaController(cadastraAgendaDiaria)

	// --- 4. Inicialização do Servidor Gin ---
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 5. Define as Rotas
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.GET("/clientes/:id", clienteController.GetCliente)

		apiV1.POST("/prestadores", prestadorController.PostPrestador)

		apiV1.POST("/catalogos", catalogoController.PostPrestador)
		apiV1.GET("/catalogos/:id", catalogoController.GetCatalogoPorID)

		apiV1.POST("/agendas", agendaDiariaController.PostPrestador)
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
