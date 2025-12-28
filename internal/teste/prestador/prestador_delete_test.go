package teste

import (
	"encoding/json"
	"fmt"
	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteAgenda_Sucesso(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	// Criar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	rrCreate := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrCreate.Code)

	// Deletar agenda
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prestadores/%s/agenda?data=2030-01-03", prestadorResp.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)

	// Verificar que foi deletada
	rrGet := SetupGetPrestadorRequest(router, prestadorResp.ID)
	var prestador response_prestador.PrestadorResponse
	json.Unmarshal(rrGet.Body.Bytes(), &prestador)

	assert.Len(t, prestador.Agenda, 0, "Agenda deve estar vazia")
}

func TestDeleteAgenda_AgendaNaoEncontrada(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	// Tentar deletar agenda que não existe
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prestadores/%s/agenda?data=2030-01-03", prestadorResp.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestDeleteAgenda_PrestadorInativo(t *testing.T) {
	router, prestadorResp, repo := CriarPrestadorValidoParaTeste(t)

	// Criar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	// Inativar prestador
	prestadorResp.Ativo = false
	repo.Salvar(&prestadorResp)

	// Tentar deletar
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/prestadores/%s/agenda?data=2030-01-03", prestadorResp.ID), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)
}
