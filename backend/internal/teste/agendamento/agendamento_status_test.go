package teste

import (
	"encoding/json"
	"net/http"
	"testing"

	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/response_agendamento"
	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/require"
)

// ============ Confirmar ============

func TestConfirmarAgendamento_Sucesso(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	cliente := SetupNovoCliente(clienteRepo)
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	// Cria o agendamento
	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	var criado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrPost.Body.Bytes(), &criado))
	require.Equal(t, domain.Pendente, criado.Status)

	// Confirma
	rrConfirmar := SetupConfirmarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rrConfirmar.Code)

	var confirmado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrConfirmar.Body.Bytes(), &confirmado))
	require.Equal(t, domain.Confirmado, confirmado.Status)
	require.Equal(t, criado.ID, confirmado.ID)
}

func TestConfirmarAgendamento_ErroNaoEncontrado(t *testing.T) {
	router, _, _, _, _ := SetupRouterAgendamento()

	rr := SetupConfirmarAgendamentoRequest(router, "id-inexistente")
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestConfirmarAgendamento_ErroJaCancelado(t *testing.T) {
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	var criado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrPost.Body.Bytes(), &criado))

	// Cancela primeiro
	rrCancelar := SetupCancelarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rrCancelar.Code)

	// Tenta confirmar após cancelado
	rrConfirmar := SetupConfirmarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusConflict, rrConfirmar.Code)
}

// ============ Cancelar ============

func TestCancelarAgendamento_Sucesso(t *testing.T) {
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	var criado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrPost.Body.Bytes(), &criado))
	require.Equal(t, domain.Pendente, criado.Status)

	// Cancela
	rrCancelar := SetupCancelarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rrCancelar.Code)

	var cancelado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrCancelar.Body.Bytes(), &cancelado))
	require.Equal(t, domain.Cancelado, cancelado.Status)
	require.Equal(t, criado.ID, cancelado.ID)
}

func TestCancelarAgendamento_ErroNaoEncontrado(t *testing.T) {
	router, _, _, _, _ := SetupRouterAgendamento()

	rr := SetupCancelarAgendamentoRequest(router, "id-inexistente")
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestCancelarAgendamento_ErroJaCancelado(t *testing.T) {
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	var criado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrPost.Body.Bytes(), &criado))

	// Cancela pela primeira vez
	rr1 := SetupCancelarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rr1.Code)

	// Tenta cancelar novamente
	rr2 := SetupCancelarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusConflict, rr2.Code)
}

func TestCancelarAgendamentoConfirmado_Sucesso(t *testing.T) {
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	var criado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrPost.Body.Bytes(), &criado))

	// Confirma
	rrConfirmar := SetupConfirmarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rrConfirmar.Code)

	// Cancela após confirmado (deve ser permitido)
	rrCancelar := SetupCancelarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rrCancelar.Code)

	var cancelado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrCancelar.Body.Bytes(), &cancelado))
	require.Equal(t, domain.Cancelado, cancelado.Status)
}

func TestHorarioCanceladoFicaDisponivel(t *testing.T) {
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}

	// Cria o agendamento
	rrPost := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost.Code)

	var criado response_agendamento.AgendamentoResponse
	require.NoError(t, json.Unmarshal(rrPost.Body.Bytes(), &criado))

	// Cancela
	rrCancelar := SetupCancelarAgendamentoRequest(router, criado.ID)
	require.Equal(t, http.StatusOK, rrCancelar.Code)

	// Tenta criar novo agendamento no mesmo horário — deve ser permitido
	rrPost2 := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rrPost2.Code, "horário cancelado deve ficar disponível para novo agendamento")
}
