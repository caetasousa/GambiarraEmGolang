package teste

import (
	"net/http"
	"testing"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/require"
)

func TestPutAgenda_Sucesso(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
			{HoraInicio: "14:00", HoraFim: "18:00"},
		},
	}

	rrAgenda := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)
}

func TestPutAgenda_PrestadorNaoEncontrado(t *testing.T) {
	router, _, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}

	rr := SetupPutAgendaRequest(router, "id-inexistente", agendaInput)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "prestador não encontrado")
}

func TestPutAgenda_AgendaDuplicada(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}

	rr1 := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rr1.Code)

	rr2 := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusConflict, rr2.Code)
	require.Contains(t, rr2.Body.String(), "agenda")
}

func TestPutAgenda_PrestadorInativo(t *testing.T) {
	router, prestadorResp, repo := CriarPrestadorValidoParaTeste(t)

	prestadorResp.Ativo = false
	err := repo.Salvar(&prestadorResp)
	require.NoError(t, err)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2025-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusConflict, rr.Code)
	require.Contains(t, rr.Body.String(), "inativo")
}

func TestPutAgenda_HorarioInicioMaiorQueFim(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{
				HoraInicio: "18:00",
				HoraFim:    "08:00",
			},
		},
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), domain.ErrIntervaloHorarioInvalido.Error())
}

func TestPutAgenda_AgendaSemIntervalos(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data:       "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{},
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), domain.ErrAgendaSemIntervalos.Error())
}
func TestPutAgenda_DataNoPassado(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2020-01-01", // Data no passado
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "passado")
}

// 3. Múltiplos intervalos válidos
func TestPutAgenda_MultiposIntervalos(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
			{HoraInicio: "14:00", HoraFim: "18:00"},
			{HoraInicio: "19:00", HoraFim: "22:00"},
		},
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rr.Code)
}