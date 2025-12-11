package main

import (
	"log"
	_ "meu-servico-agenda/docs"

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

	// 2. Camada de Aplicação (Serviços/Casos de Uso)
	cadastradoService := service.NovoServiceCliente(clienteRepo)
	cadastroService := service.NovoPrestadorService(prestadorRepo)

	// 3. Camada de Adaptador HTTP (Controller)
	clienteController := cliente.NovoClienteController(cadastradoService)
	prestadorController := prestador.NovoPrestadorController(cadastroService)

	// --- 4. Inicialização do Servidor Gin ---
	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 5. Define as Rotas
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.GET("/clientes/:id", clienteController.GetCliente)

		apiV1.POST("/prestadores", prestadorController.PostPrestador)
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
