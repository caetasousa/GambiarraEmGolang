package teste

import (
	"bytes"
	"encoding/json"

	Http "meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/services"
	"meu-servico-agenda/internal/core/domain"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouterCliente() (*gin.Engine, *repository.FakeClienteRepositorio, *Http.ClienteController) {
	gin.SetMode(gin.TestMode)

	// 1. Camada de Repositório (Infraestrutura)
	clienteRepo := repository.NewFakeClienteRepositorio()

	// 2. Camada de Aplicação (Serviços/Casos de Uso)
	cadastradorService := services.NovoCadastradoDeCliente(clienteRepo)

	// 3. Camada de Adaptador HTTP (Controller)
	clienteController := Http.NovoClienteController(cadastradorService)

	// 🔥 Router REAL
	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
	}

	return router, clienteRepo, clienteController
}

func TestPostCliente_ResultadoEsperado(t *testing.T) {
	router, _, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Ana",
		Email:    "ana@example.com",
		Telefone: "6299697481",
	}

	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// 	// === Validações === //

	assert.NotEqual(t, http.StatusBadRequest, rr.Code, "Não deveria retornar 400, JSON é válido")
	assert.NotEqual(t, http.StatusInternalServerError, rr.Code, "Serviço real não deveria causar panic ou erro interno")
	assert.Equal(t, http.StatusCreated, rr.Code, "Esperado que o serviço real retorne 201 Created")

	var resp domain.Cliente
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")

	// Campos obrigatórios
	assert.Equal(t, input.Nome, resp.Nome)
	assert.Equal(t, input.Telefone, resp.Telefone)

	// O serviço real deve gerar ID
	assert.NotZero(t, resp.ID, "O serviço real deve gerar ID")
}
