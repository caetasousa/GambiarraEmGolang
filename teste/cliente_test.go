package teste

import (
	"bytes"
	"encoding/json"
	"fmt"

	Http "meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouterCliente() (*gin.Engine, *repository.FakeClienteRepositorio) {
	gin.SetMode(gin.TestMode)

	clienteRepo := repository.NewFakeClienteRepositorio()
	cadastradorService := service.NovoServiceCliente(clienteRepo)
	clienteController := Http.NovoClienteController(cadastradorService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.GET("/clientes/:id", clienteController.GetCliente)
	}

	return router, clienteRepo
}

func SetupPostClienteRequest(router *gin.Engine, input interface{}) *httptest.ResponseRecorder {
	// 1. Converte o corpo (input) para JSON (aceita structs ou map[string]interface{})
	body, _ := json.Marshal(input)

	// 2. Cria a Requisição HTTP
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// 3. Cria o Response Recorder (rr)
	rr := httptest.NewRecorder()

	// 4. Executa a Requisição no Router
	router.ServeHTTP(rr, req)

	// 5. Retorna o ResponseRecorder com o resultado do teste
	return rr
}

func TestPostCliente_ResultadoEsperado(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Ana",
		Email:    "ana@example.com",
		Telefone: "6299697481",
	}
	rr := SetupPostClienteRequest(router, input)

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

func TestPostCliente_EmailInvalido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "Carlos",
		Email:    "email-invalido", // Email inválido para simular erro de validação
		Telefone: "6299697482",
	}
	rr := SetupPostClienteRequest(router, input)

	// === Validações === //
	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para email inválido")
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")
	// Verifica se a mensagem de erro contém informações sobre o email inválido
	errorMsg, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Equal(t, errorMsg, "Dados inválidos: Key: 'ClienteRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag")
}

func TestPostCliente_NomeRequerido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "",
		Email:    "email@invalido.com",
		Telefone: "6299697482",
	}
	rr := SetupPostClienteRequest(router, input)

	// === Validações === //

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para email inválido")
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")

	errorMsg, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Equal(t, errorMsg, "Dados inválidos: Key: 'ClienteRequest.Nome' Error:Field validation for 'Nome' failed on the 'required' tag")
}

func TestPostCliente_TelefoneReequerido(t *testing.T) {
	router, _ := SetupRouterCliente()

	input := request.ClienteRequest{
		Nome:     "eduardo",
		Email:    "email@invalido.com",
		Telefone: "",
	}
	rr := SetupPostClienteRequest(router, input)

	// === Validações === //

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para email inválido")
	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")

	errorMsg, exists := resp["error"]
	assert.True(t, exists, "Deve conter campo 'error' na resposta")
	assert.Equal(t, errorMsg, "Dados inválidos: Key: 'ClienteRequest.Telefone' Error:Field validation for 'Telefone' failed on the 'required' tag")
}

func SetupGetClienteRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	// 1. Cria a Requisição HTTP
	url := fmt.Sprintf("/api/v1/clientes/%s", id)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	// 2. Cria o Response Recorder (rr)
	rr := httptest.NewRecorder()

	// 3. Executa a Requisição no Router
	// router.ServeHTTP(rr, req) não retorna nada, mas modifica o 'rr'
	router.ServeHTTP(rr, req)

	// 4. Retorna o ResponseRecorder com o resultado do teste
	return rr
}

func TestGetCliente_Sucesso(t *testing.T) {
	// Garante que a rota GET está mapeada no SetupRouterCliente!
	router, _ := SetupRouterCliente()

	// 1. Ação de Setup (Cadastro de um Cliente Real)
	input := request.ClienteRequest{
		Nome:     "Beatriz",
		Email:    "caetaasousa@gmail.com",
		Telefone: "6299977848",
	}
	rrPost := SetupPostClienteRequest(router, input)
	var clienteCriado domain.Cliente
	err := json.Unmarshal(rrPost.Body.Bytes(), &clienteCriado)
	assert.NoError(t, err, "Resposta do POST deve ser um JSON válido")

	// 2. Execução (Busca pelo Cliente Recém-Criado)
	rrGet := SetupGetClienteRequest(router, clienteCriado.ID)
	// === Validações === //
	assert.Equal(t, http.StatusOK, rrGet.Code, "Esperado 200 OK para cliente existente")

	var clienteBuscado domain.Cliente
	err = json.Unmarshal(rrGet.Body.Bytes(), &clienteBuscado)
	assert.NoError(t, err, "Resposta do GET deve ser um JSON válido")
	// 3. Valida os Dados do Cliente Buscado
	assert.Equal(t, clienteCriado.ID, clienteBuscado.ID, "IDs devem ser iguais")
	assert.Equal(t, clienteCriado.Nome, clienteBuscado.Nome, "Nomes devem ser iguais")
	assert.Equal(t, clienteCriado.Email, clienteBuscado.Email, "Emails devem ser iguais")
	assert.Equal(t, clienteCriado.Telefone, clienteBuscado.Telefone, "Telefones devem ser iguais")
}

func TestGetCliente_NaoEncontrado(t *testing.T) {
	// Garante que a rota GET está mapeada no SetupRouterCliente!
	router, _ := SetupRouterCliente()
	// 1. Execução (Busca por um ID que não existe)
	// O ID "id-inexistente" é um ID que o sistema com certeza não gerou.
	rr := SetupGetClienteRequest(router, "id-inexistente")

	// === Validações === //

	// 3. Verifica o Status Code (Deve ser 404)
	assert.Equal(t, http.StatusNotFound, rr.Code, "Esperado 404 Not Found para cliente inexistente")
}

func TestPostCliente_NomeMuitoCurto_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCliente()

	// Nome com menos de 3 caracteres
	input := request.ClienteRequest{
		Nome:     "AB",
		Email:    "ab@example.com",
		Telefone: "62999974848",
	}
	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para nome com < 3 caracteres")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCliente_TelefoneMuitoCurto_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCliente()

	// Telefone com menos de 8 dígitos
	input := request.ClienteRequest{
		Nome:     "João",
		Email:    "joao@example.com",
		Telefone: "1234567",
	}
	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para telefone com < 8 dígitos")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCliente_TelefoneComTipoInvalido_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCliente()

	// Telefone como número em vez de string
	input := map[string]interface{}{
		"nome":     "Carlos",
		"telefone": 6299974848,
		"email":    "carlos@example.com",
	}
	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para tipo inválido de telefone")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCliente_EmailOmitido_DeveRetornar201(t *testing.T) {
	router, _ := SetupRouterCliente()

	// Email é omitempty, então não é obrigatório
	input := request.ClienteRequest{
		Nome:     "Maria",
		Email:    "",
		Telefone: "62999974848",
	}
	rr := SetupPostClienteRequest(router, input)

	assert.Equal(t, http.StatusCreated, rr.Code, "Email omitido deveria ser aceito (201)")
	var resp domain.Cliente
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err, "Resposta deve ser um JSON válido")
	assert.Equal(t, input.Nome, resp.Nome)
	assert.Equal(t, "", resp.Email)
}
