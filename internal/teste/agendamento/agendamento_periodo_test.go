package teste

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/response_agendamento"
	"meu-servico-agenda/internal/core/domain"

	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// TESTES PARA LISTAGEM PAGINADA DE AGENDAMENTOS DO CLIENTE

func TestGetAgendamentosClientePeriodo_Sucesso(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria um agendamento
	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Teste paginação",
	}
	SetupPostAgendamentoRequest(router, input)

	// Busca paginada
	rr := SetupGetAgendamentosClientePeriodoRequest(router, cliente.ID, "2030-01-01", "2030-01-31", 1, 10)

	require.Equal(t, http.StatusOK, rr.Code)

	var response response_agendamento.AgendamentoListPaginadoResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, 1, len(response.Data))
	require.Equal(t, 1, response.Page)
	require.Equal(t, 10, response.Limit)
	require.Equal(t, 1, response.Total)
}

func TestGetAgendamentosClientePeriodo_ClienteNaoEncontrado(t *testing.T) {
	router, _, _, _, _ := SetupRouterAgendamento()

	rr := SetupGetAgendamentosClientePeriodoRequest(router, "cliente-inexistente", "2030-01-01", "2030-01-31", 1, 10)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetAgendamentosClientePeriodo_Paginacao(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)

	// Cria segundo catálogo para permitir múltiplos agendamentos no mesmo dia
	cat1, listaDeCatalogos1 := SetupNovoCatalogo(catalogoRepo)
	cat2, _ := domain.NovoCatalogo(xid.New().String(), "Corte", 30, 5000, "Cabelo", "https://exemplo.com/img2.jpg")
	catalogoRepo.Salvar(cat2)

	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos1)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria 2 agendamentos (categorias diferentes para evitar erro de agendamento duplo)
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     cat1.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	SetupPostAgendamentoRequest(router, input1)

	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     cat2.ID,
		DataHoraInicio: "2030-01-03T09:00:00Z",
	}
	SetupPostAgendamentoRequest(router, input2)

	// Página 1 com limite 1
	rr1 := SetupGetAgendamentosClientePeriodoRequest(router, cliente.ID, "2030-01-01", "2030-01-31", 1, 1)
	require.Equal(t, http.StatusOK, rr1.Code)

	var response1 response_agendamento.AgendamentoListPaginadoResponse
	json.Unmarshal(rr1.Body.Bytes(), &response1)

	require.Equal(t, 1, len(response1.Data))
	require.Equal(t, 1, response1.Page)
	require.Equal(t, 1, response1.Limit)
	require.Equal(t, 2, response1.Total)

	// Página 2 com limite 1
	rr2 := SetupGetAgendamentosClientePeriodoRequest(router, cliente.ID, "2030-01-01", "2030-01-31", 2, 1)
	require.Equal(t, http.StatusOK, rr2.Code)

	var response2 response_agendamento.AgendamentoListPaginadoResponse
	json.Unmarshal(rr2.Body.Bytes(), &response2)

	require.Equal(t, 1, len(response2.Data))
	require.Equal(t, 2, response2.Page)
	require.Equal(t, 1, response2.Limit)
	require.Equal(t, 2, response2.Total)
}

func TestGetAgendamentosClientePeriodo_SemAgendamentos(t *testing.T) {
	router, _, clienteRepo, _, _ := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)

	rr := SetupGetAgendamentosClientePeriodoRequest(router, cliente.ID, "2030-01-01", "2030-01-31", 1, 10)

	require.Equal(t, http.StatusOK, rr.Code)

	var response response_agendamento.AgendamentoListPaginadoResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, 0, len(response.Data))
	require.Equal(t, 0, response.Total)
}

func TestGetAgendamentosClientePeriodo_ValidaDefaultsPaginacao(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
	}
	SetupPostAgendamentoRequest(router, input)

	// Testa sem passar page e limit (deveria usar defaults)
	url := fmt.Sprintf("/api/v1/agendamentos/cliente/%s/periodo?data_inicio=2030-01-01&data_fim=2030-01-31", cliente.ID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var response response_agendamento.AgendamentoListPaginadoResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verifica que encontrou o agendamento (defaults funcionaram)
	require.Equal(t, 1, len(response.Data))
	require.Equal(t, 1, response.Total)
}

// TESTES PARA LISTAGEM PAGINADA DE AGENDAMENTOS DO PRESTADOR

func TestGetAgendamentosPrestadorPeriodo_Sucesso(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Teste prestador",
	}
	SetupPostAgendamentoRequest(router, input)

	rr := SetupGetAgendamentosPrestadorPeriodoRequest(router, prestador.ID, "2030-01-01", "2030-01-31", 1, 10)

	require.Equal(t, http.StatusOK, rr.Code)

	var response response_agendamento.AgendamentoListPaginadoResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, 1, len(response.Data))
	require.Equal(t, 1, response.Page)
	require.Equal(t, 10, response.Limit)
	require.Equal(t, 1, response.Total)
	require.Equal(t, "Teste prestador", response.Data[0].Notas)
}

func TestGetAgendamentosPrestadorPeriodo_PrestadorNaoEncontrado(t *testing.T) {
	router, _, _, _, _ := SetupRouterAgendamento()

	rr := SetupGetAgendamentosPrestadorPeriodoRequest(router, "prestador-inexistente", "2030-01-01", "2030-01-31", 1, 10)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestGetAgendamentosPrestadorPeriodo_MultipleClientes(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// Cria 2 clientes
	cliente1 := SetupNovoCliente(clienteRepo)
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	cliente2 := domain.NovoCliente(xid.New().String(), "Maria", "maria@email.com", "62988888888", string(hash))
	clienteRepo.Salvar(cliente2)

	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria 2 agendamentos com clientes diferentes
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente1.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	SetupPostAgendamentoRequest(router, input1)

	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente2.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T09:30:00Z",
	}
	SetupPostAgendamentoRequest(router, input2)

	rr := SetupGetAgendamentosPrestadorPeriodoRequest(router, prestador.ID, "2030-01-01", "2030-01-31", 1, 10)

	require.Equal(t, http.StatusOK, rr.Code)

	var response response_agendamento.AgendamentoListPaginadoResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	require.Equal(t, 2, len(response.Data))
	require.Equal(t, 2, response.Total)
}

func TestGetAgendamentosPrestadorPeriodo_OrdenacaoPorDataHora(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)

	// Cria dois catálogos diferentes para evitar erro de agendamento duplo
	cat1, listaDeCatalogos1 := SetupNovoCatalogo(catalogoRepo)
	cat2, _ := domain.NovoCatalogo(xid.New().String(), "Corte", 30, 5000, "Cabelo", "https://exemplo.com/img2.jpg")
	catalogoRepo.Salvar(cat2)

	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos1)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria agendamentos fora de ordem com categorias diferentes
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     cat2.ID,
		DataHoraInicio: "2030-01-03T11:00:00Z",
		Notas:          "Segundo",
	}
	SetupPostAgendamentoRequest(router, input2)

	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     cat1.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
		Notas:          "Primeiro",
	}
	SetupPostAgendamentoRequest(router, input1)

	rr := SetupGetAgendamentosPrestadorPeriodoRequest(router, prestador.ID, "2030-01-01", "2030-01-31", 1, 10)

	var response response_agendamento.AgendamentoListPaginadoResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verifica que está ordenado por data_hora_inicio
	require.Equal(t, 2, len(response.Data))
	require.Equal(t, "Primeiro", response.Data[0].Notas)
	require.Equal(t, "Segundo", response.Data[1].Notas)
}

func TestGetAgendamentosPeriodo_DataInvalidaRetorna400(t *testing.T) {
	router, _, clienteRepo, _, _ := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)

	// Data inválida
	url := fmt.Sprintf("/api/v1/agendamentos/cliente/%s/periodo?data_inicio=data-invalida&data_fim=2030-01-31", cliente.ID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAgendamentosPeriodo_ValidaDadosCompletos(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Validação completa",
	}
	SetupPostAgendamentoRequest(router, input)

	rr := SetupGetAgendamentosPrestadorPeriodoRequest(router, prestador.ID, "2030-01-01", "2030-01-31", 1, 10)

	var response response_agendamento.AgendamentoListPaginadoResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	agendamento := response.Data[0]

	// Valida estrutura completa do agendamento
	require.NotEmpty(t, agendamento.ID)

	// Cliente
	require.Equal(t, cliente.ID, agendamento.Cliente.ID)
	require.Equal(t, cliente.Nome, agendamento.Cliente.Nome)
	require.Equal(t, cliente.Email, agendamento.Cliente.Email)

	// Prestador
	require.Equal(t, prestador.ID, agendamento.Prestador.ID)
	require.Equal(t, prestador.Nome, agendamento.Prestador.Nome)
	require.Equal(t, prestador.Cpf.String(), agendamento.Prestador.CPF)

	// Serviço
	require.Equal(t, catalogo.ID, agendamento.Servico.ID)
	require.Equal(t, catalogo.Nome, agendamento.Servico.Nome)
	require.Equal(t, catalogo.Preco, agendamento.Servico.Preco)

	// Dados do agendamento
	require.Equal(t, domain.Pendente, agendamento.Status)
	require.Equal(t, "Validação completa", agendamento.Notas)
}
