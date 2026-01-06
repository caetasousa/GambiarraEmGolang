package teste

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/response_agendamento"
	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/require"
)

func TestGetAgendamentoClienteData_ValidaDadosCompletos(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	// cadastro do catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	// cria agenda diária
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)

	// adiciona agenda diária ao prestador
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria agendamento
	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Agendamento de teste",
	}

	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	// Parse do agendamento criado para pegar o ID
	var createdAgendamento response_agendamento.AgendamentoResponse
	err := json.Unmarshal(rrPost.Body.Bytes(), &createdAgendamento)
	require.NoError(t, err)

	// Busca os agendamentos do cliente
	rr := SetupGetAgendamentoClienteDataRequest(router, cliente.ID, "2030-01-03")

	// Verifica status code
	require.Equal(t, http.StatusOK, rr.Code)

	// Parse da resposta
	var response response_agendamento.BuscaDataResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verifica que retornou 1 agendamento
	require.Len(t, response.Data, 1)

	agendamento := response.Data[0]

	// Valida ID do agendamento
	require.Equal(t, createdAgendamento.ID, agendamento.ID)
	require.NotEmpty(t, agendamento.ID)

	// ===== Valida dados do CLIENTE =====
	require.NotNil(t, agendamento.Cliente)
	require.Equal(t, cliente.ID, agendamento.Cliente.ID)
	require.Equal(t, cliente.Nome, agendamento.Cliente.Nome)
	require.Equal(t, cliente.Email, agendamento.Cliente.Email)
	require.Equal(t, cliente.Telefone, agendamento.Cliente.Telefone)

	// ===== Valida dados do PRESTADOR =====
	require.NotNil(t, agendamento.Prestador)
	require.Equal(t, prestador.ID, agendamento.Prestador.ID)
	require.Equal(t, prestador.Nome, agendamento.Prestador.Nome)
	require.Equal(t, prestador.Cpf, agendamento.Prestador.CPF)
	require.Equal(t, prestador.Email, agendamento.Prestador.Email)
	require.Equal(t, prestador.Telefone, agendamento.Prestador.Telefone)
	require.True(t, agendamento.Prestador.Ativo)

	// ===== Valida dados do SERVIÇO =====
	require.NotNil(t, agendamento.Servico)
	require.Equal(t, catalogo.ID, agendamento.Servico.ID)
	require.Equal(t, catalogo.Nome, agendamento.Servico.Nome)
	require.Equal(t, catalogo.DuracaoPadrao, agendamento.Servico.Duracao)
	require.Equal(t, catalogo.Preco, agendamento.Servico.Preco)
	require.Equal(t, catalogo.Categoria, agendamento.Servico.Categoria)

	// ===== Valida DATAS =====
	expectedInicio, _ := time.Parse(time.RFC3339, "2030-01-03T10:00:00Z")
	expectedFim := expectedInicio.Add(time.Duration(catalogo.DuracaoPadrao) * time.Minute)

	require.True(t, agendamento.DataInicio.Equal(expectedInicio),
		"Data de início esperada: %v, recebida: %v", expectedInicio, agendamento.DataInicio)
	require.True(t, agendamento.DataFim.Equal(expectedFim),
		"Data de fim esperada: %v, recebida: %v", expectedFim, agendamento.DataFim)

	// ===== Valida STATUS e NOTAS =====
	require.Equal(t, domain.Pendente, agendamento.Status)
	require.Equal(t, "Agendamento de teste", agendamento.Notas)
}

func TestGetAgendamentoPrestadorData_ValidaDadosCompletos(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	// cadastro do catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	// cria agenda diária
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)

	// adiciona agenda diária ao prestador
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria agendamento
	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Agendamento do prestador",
	}

	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	// Parse do agendamento criado para pegar o ID
	var createdAgendamento response_agendamento.AgendamentoResponse
	err := json.Unmarshal(rrPost.Body.Bytes(), &createdAgendamento)
	require.NoError(t, err)

	// Busca os agendamentos do prestador
	url := fmt.Sprintf("/api/v1/agendamentos/prestador/%s?data=%s", prestador.ID, "2030-01-03")
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Verifica status code
	require.Equal(t, http.StatusOK, rr.Code)

	// Parse da resposta
	var response response_agendamento.BuscaDataResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Verifica que retornou 1 agendamento
	require.Len(t, response.Data, 1)

	agendamento := response.Data[0]

	// Valida ID do agendamento
	require.Equal(t, createdAgendamento.ID, agendamento.ID)
	require.NotEmpty(t, agendamento.ID)

	// ===== Valida dados do CLIENTE =====
	require.NotNil(t, agendamento.Cliente)
	require.Equal(t, cliente.ID, agendamento.Cliente.ID)
	require.Equal(t, cliente.Nome, agendamento.Cliente.Nome)
	require.Equal(t, cliente.Email, agendamento.Cliente.Email)
	require.Equal(t, cliente.Telefone, agendamento.Cliente.Telefone)

	// ===== Valida dados do PRESTADOR =====
	require.NotNil(t, agendamento.Prestador)
	require.Equal(t, prestador.ID, agendamento.Prestador.ID)
	require.Equal(t, prestador.Nome, agendamento.Prestador.Nome)
	require.Equal(t, prestador.Cpf, agendamento.Prestador.CPF)
	require.Equal(t, prestador.Email, agendamento.Prestador.Email)
	require.Equal(t, prestador.Telefone, agendamento.Prestador.Telefone)
	require.True(t, agendamento.Prestador.Ativo)

	// ===== Valida dados do SERVIÇO =====
	require.NotNil(t, agendamento.Servico)
	require.Equal(t, catalogo.ID, agendamento.Servico.ID)
	require.Equal(t, catalogo.Nome, agendamento.Servico.Nome)
	require.Equal(t, catalogo.DuracaoPadrao, agendamento.Servico.Duracao)
	require.Equal(t, catalogo.Preco, agendamento.Servico.Preco)
	require.Equal(t, catalogo.Categoria, agendamento.Servico.Categoria)

	// ===== Valida DATAS =====
	expectedInicio, _ := time.Parse(time.RFC3339, "2030-01-03T10:00:00Z")
	expectedFim := expectedInicio.Add(time.Duration(catalogo.DuracaoPadrao) * time.Minute)

	require.True(t, agendamento.DataInicio.Equal(expectedInicio),
		"Data de início esperada: %v, recebida: %v", expectedInicio, agendamento.DataInicio)
	require.True(t, agendamento.DataFim.Equal(expectedFim),
		"Data de fim esperada: %v, recebida: %v", expectedFim, agendamento.DataFim)

	// ===== Valida STATUS e NOTAS =====
	require.Equal(t, domain.Pendente, agendamento.Status)
	require.Equal(t, "Agendamento do prestador", agendamento.Notas)
}
