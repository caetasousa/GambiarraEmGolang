package teste

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	clienteHttp "meu-servico-agenda/internal/adapters/http/cliente"
	"meu-servico-agenda/internal/adapters/http/cliente/request"
	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// SetupPostClienteRequestGeneric executa request de criação de cliente usando router genérico
func SetupPostClienteRequestGeneric(router *gin.Engine, input request.ClienteRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/clientes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// TestPostPrestador_EmailJaUsadoPorCliente verifica que não é possível
// cadastrar um prestador com um email que já está sendo usado por um cliente
func TestPostPrestador_EmailJaUsadoPorCliente(t *testing.T) {
	// 1. Setup dos repositórios e serviços
	router, prestadorRepo, clienteRepo := SetupPostPrestadorComCliente()
	catalogoResp := CriarCatalogoValido(t, router)

	// 2. Criar um cliente com o email que vamos tentar usar
	emailCompartilhado := "email.compartilhado@teste.com"

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	agendaDiariaRepo := repository.NovoFakeAgendaDiariaRepositorio()
	authService := service.NovoAuthService("test-secret-key", 24, clienteRepo, prestadorRepo)
	clienteService := service.NovoServiceCliente(clienteRepo, prestadorRepo)
	service.NovaPrestadorService(prestadorRepo, catalogoRepo, agendaDiariaRepo, authService, clienteRepo)
	clienteController := clienteHttp.NovoClienteController(clienteService)

	// Adicionar rota de cliente ao router
	router.POST("/api/v1/clientes", clienteController.PostCliente)

	clienteInput := request.ClienteRequest{
		Nome:     "Cliente Existente",
		Email:    emailCompartilhado,
		Telefone: "62999998888",
		Senha:    "senha123",
	}

	// Criar cliente primeiro
	rrCliente := SetupPostClienteRequestGeneric(router, clienteInput)
	require.Equal(t, http.StatusCreated, rrCliente.Code, "Cliente deveria ser criado com sucesso. Resposta: %s", rrCliente.Body.String())

	// 3. Tentar criar prestador com o mesmo email
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Novo Prestador",
		Email:       emailCompartilhado, // Mesmo email do cliente!
		Cpf:         "04423258196",
		Telefone:    "62999677481",
		Senha:       "senha123",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)

	// 4. Deve retornar erro de conflito (409)
	require.Equal(t, http.StatusConflict, rrPrestador.Code, "Deveria retornar 409 Conflict quando email já está em uso por um cliente. Resposta: %s", rrPrestador.Body.String())
	require.Contains(t, rrPrestador.Body.String(), "email", "Resposta deveria mencionar o problema com email")
}

// TestPostCliente_EmailJaUsadoPorPrestador verifica que não é possível
// cadastrar um cliente com um email que já está sendo usado por um prestador
func TestPostCliente_EmailJaUsadoPorPrestador(t *testing.T) {
	// 1. Setup dos repositórios e serviços
	router, prestadorRepo, clienteRepo := SetupPostPrestadorComCliente()
	catalogoResp := CriarCatalogoValido(t, router)

	// 2. Criar um prestador com o email que vamos tentar usar
	emailCompartilhado := "prestador.email@teste.com"

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Prestador Existente",
		Email:       emailCompartilhado,
		Cpf:         "04423258196",
		Telefone:    "62999677481",
		Senha:       "senha123",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	// Criar prestador primeiro
	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrPrestador.Code, "Prestador deveria ser criado com sucesso. Resposta: %s", rrPrestador.Body.String())

	// 3. Configurar o serviço de cliente COM validação de prestador
	catalogoRepo := repository.NovoCatalogoFakeRepo()
	agendaDiariaRepo := repository.NovoFakeAgendaDiariaRepositorio()
	authService := service.NovoAuthService("test-secret-key", 24, clienteRepo, prestadorRepo)
	clienteService := service.NovoServiceCliente(clienteRepo, prestadorRepo)
	service.NovaPrestadorService(prestadorRepo, catalogoRepo, agendaDiariaRepo, authService, clienteRepo)
	clienteController := clienteHttp.NovoClienteController(clienteService)

	// Adicionar rota de cliente ao router
	router.POST("/api/v1/clientes", clienteController.PostCliente)

	clienteInput := request.ClienteRequest{
		Nome:     "Novo Cliente",
		Email:    emailCompartilhado, // Mesmo email do prestador!
		Telefone: "62999998888",
		Senha:    "senha123",
	}

	// 4. Tentar criar cliente com o mesmo email
	rrCliente := SetupPostClienteRequestGeneric(router, clienteInput)

	// 5. Deve retornar erro de conflito (409)
	require.Equal(t, http.StatusConflict, rrCliente.Code, "Deveria retornar 409 Conflict quando email já está em uso por um prestador. Resposta: %s", rrCliente.Body.String())
	require.Contains(t, rrCliente.Body.String(), "email", "Resposta deveria mencionar o problema com email")
}
