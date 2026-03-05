package teste

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	clienteHttp "meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"

	"github.com/gin-gonic/gin"
)

// ============ Setup Router ============

// SetupRouterCliente configura o ambiente de testes para cliente
func SetupRouterCliente() (*gin.Engine, port.ClienteRepositorio) {
	gin.SetMode(gin.TestMode)

	clienteRepo := repository.NewFakeClienteRepositorio()
	service.NovoAuthService("test-secret-key", 24, clienteRepo, nil)
	clienteService := service.NovoServiceCliente(clienteRepo, nil)
	clienteController := clienteHttp.NovoClienteController(clienteService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.GET("/clientes/:id", clienteController.GetCliente)
	}

	return router, clienteRepo
}

// SetupRouterClienteComPrestador configura o ambiente de testes com validação de email de prestador
func SetupRouterClienteComPrestador() (*gin.Engine, port.ClienteRepositorio, port.PrestadorRepositorio) {
	gin.SetMode(gin.TestMode)

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	clienteRepo := repository.NewFakeClienteRepositorio()
	prestadorRepo := repository.NovoFakePrestadorRepositorio(catalogoRepo)
	
	service.NovoAuthService("test-secret-key", 24, clienteRepo, prestadorRepo)
	clienteService := service.NovoServiceCliente(clienteRepo, prestadorRepo)
	clienteController := clienteHttp.NovoClienteController(clienteService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.GET("/clientes/:id", clienteController.GetCliente)
	}

	return router, clienteRepo, prestadorRepo
}

// ============ POST Helpers ============

// SetupPostClienteRequest executa request de criação de cliente
func SetupPostClienteRequest(router *gin.Engine, input request.ClienteRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// PostRawJSON executa request POST com JSON bruto
func PostRawJSON(router *gin.Engine, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// PostJSONFromMap executa request POST com map
func PostJSONFromMap(router *gin.Engine, input map[string]interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	return PostRawJSON(router, body)
}

// ============ GET Helpers ============

// SetupGetClienteRequest executa request de busca de cliente por ID
func SetupGetClienteRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/clientes/%s", id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
