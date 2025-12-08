package main

import (
	"log"
	Http "meu-servico-agenda/internal/adapters/http"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Camada de Repositório (Infraestrutura)
	clienteRepo := repository.NewFakeClienteRepositorio()

	// 2. Camada de Aplicação (Serviços/Casos de Uso)
	// Tipo do retorno: *services.CadastroDeCliente
	cadastradorService := services.NovoCadastradoDeCliente(clienteRepo)

	// 3. Camada de Adaptador HTTP (Controller)
	// Passamos *services.CadastroDeCliente para uma função que espera *services.CadastroDeCliente.
	// O erro será resolvido após a correção do Controller (Passo 1).
	clienteController := Http.NovoClienteController(cadastradorService)

	// --- 2. Inicialização do Servidor Gin ---
	router := gin.Default()

	// 3. Define as Rotas
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
	}
	
	router.GET("/saudacao", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensagem": "Serviço rodando e pronto para rotas!"})
	})

	// 4. Inicia o Servidor
	log.Println("Servidor Gin rodando na porta 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Erro ao iniciar o servidor: ", err)
	}
}