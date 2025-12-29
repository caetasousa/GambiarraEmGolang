package teste

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrestadoresPorData_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Criar 3 prestadores
	p1 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(0))
	p2 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(1))
	p3 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(2))

	// Adicionar agenda para p1 e p2 na data futura
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: dataFutura,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "09:00", HoraFim: "12:00"},
		},
	}
	rrAgenda1 := SetupPutAgendaRequest(router, p1.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda1.Code)

	rrAgenda2 := SetupPutAgendaRequest(router, p2.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda2.Code)

	// p3 não tem agenda nessa data

	// Buscar prestadores disponíveis na data
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 2)
	assert.Equal(t, 2, resp.Total)

	// Verificar que são os prestadores corretos
	ids := []string{resp.Data[0].ID, resp.Data[1].ID}
	assert.Contains(t, ids, p1.ID)
	assert.Contains(t, ids, p2.ID)
	assert.NotContains(t, ids, p3.ID) // Não tem agenda
}

func TestGetPrestadoresPorData_ComPaginacao(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Criar 15 prestadores com agenda
	for i := 0; i < 15; i++ {
		prestador := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(i))

		agendaInput := request_prestador.AgendaDiariaRequest{
			Data: dataFutura,
			Intervalos: []request_prestador.IntervaloDiarioRequest{
				{HoraInicio: "09:00", HoraFim: "12:00"},
			},
		}
		rrAgenda := SetupPutAgendaRequest(router, prestador.ID, agendaInput)
		require.Equal(t, http.StatusNoContent, rrAgenda.Code)
	}

	// Página 1
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&page=1&limit=5", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 5)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 5, resp.Limit)
	assert.Equal(t, 15, resp.Total)

	// Página 2
	url2 := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&page=2&limit=5", dataFutura)
	req2, _ := http.NewRequest(http.MethodGet, url2, nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp2 response_prestador.PrestadorListResponse
	json.Unmarshal(rr2.Body.Bytes(), &resp2)

	assert.Len(t, resp2.Data, 5)
	assert.Equal(t, 2, resp2.Page)
	assert.Equal(t, 15, resp2.Total)

	// Página 4 (vazia)
	url4 := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&page=4&limit=5", dataFutura)
	req4, _ := http.NewRequest(http.MethodGet, url4, nil)
	rr4 := httptest.NewRecorder()
	router.ServeHTTP(rr4, req4)

	require.Equal(t, http.StatusOK, rr4.Code)

	var resp4 response_prestador.PrestadorListResponse
	json.Unmarshal(rr4.Body.Bytes(), &resp4)

	assert.Len(t, resp4.Data, 0)
	assert.Equal(t, 4, resp4.Page)
	assert.Equal(t, 15, resp4.Total)
}

func TestGetPrestadoresPorData_ListaVazia(t *testing.T) {
	router, _ := SetupPostPrestador()

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Buscar em data sem nenhum prestador com agenda
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 0)
	assert.Equal(t, 0, resp.Total)
}

func TestGetPrestadoresPorData_DataObrigatoria(t *testing.T) {
	router, _ := SetupPostPrestador()

	// Sem parâmetro data
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores/disponiveis", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "data")
}

func TestGetPrestadoresPorData_FormatoDataInvalido(t *testing.T) {
	router, _ := SetupPostPrestador()

	testCases := []struct {
		name string
		data string
	}{
		{
			name: "Formato DD/MM/YYYY",
			data: "29/12/2024",
		},
		{
			name: "Formato MM-DD-YYYY",
			data: "12-29-2024",
		},
		{
			name: "Formato inválido",
			data: "2024/12/29",
		},
		{
			name: "String não data",
			data: "abc",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", tc.data)
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusBadRequest, rr.Code)
		})
	}
}

func TestGetPrestadoresPorData_DataNoPassado(t *testing.T) {
	router, _ := SetupPostPrestador()

	dataPassada := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataPassada)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "passado")
}

func TestGetPrestadoresPorData_NaoRetornaInativos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Criar 2 prestadores
	p1 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(0))
	p2 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(1))

	// Adicionar agenda para ambos
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: dataFutura,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "09:00", HoraFim: "12:00"},
		},
	}
	SetupPutAgendaRequest(router, p1.ID, agendaInput)
	SetupPutAgendaRequest(router, p2.ID, agendaInput)

	// Inativar p1
	reqInativar, _ := http.NewRequest(http.MethodPut, "/api/v1/prestadores/"+p1.ID+"/inativar", nil)
	rrInativar := httptest.NewRecorder()
	router.ServeHTTP(rrInativar, reqInativar)
	require.Equal(t, http.StatusNoContent, rrInativar.Code)

	// Buscar prestadores disponíveis
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 1, "Deve retornar apenas 1 prestador ativo")
	assert.Equal(t, p2.ID, resp.Data[0].ID)
	assert.True(t, resp.Data[0].Ativo)
}

func TestGetPrestadoresPorData_ComAgendasECatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Criar prestador
	prestador := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(0))

	// Adicionar agenda com múltiplos intervalos
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: dataFutura,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
			{HoraInicio: "14:00", HoraFim: "18:00"},
		},
	}
	rrAgenda := SetupPutAgendaRequest(router, prestador.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)

	// Buscar prestadores disponíveis
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	require.Len(t, resp.Data, 1)

	p := resp.Data[0]
	
	// Verificar catálogos
	assert.NotEmpty(t, p.Catalogo)
	assert.Len(t, p.Catalogo, 1)
	assert.Equal(t, "Corte de Cabelo", p.Catalogo[0].Nome)

	// Verificar agenda
	assert.NotEmpty(t, p.Agenda)
	assert.Len(t, p.Agenda, 1)
	assert.Equal(t, dataFutura, p.Agenda[0].Data)
	
	// Verificar intervalos
	assert.Len(t, p.Agenda[0].Intervalos, 2)
	assert.Equal(t, "08:00:00", p.Agenda[0].Intervalos[0].HoraInicio)
	assert.Equal(t, "12:00:00", p.Agenda[0].Intervalos[0].HoraFim)
	assert.Equal(t, "14:00:00", p.Agenda[0].Intervalos[1].HoraInicio)
	assert.Equal(t, "18:00:00", p.Agenda[0].Intervalos[1].HoraFim)
}

func TestGetPrestadoresPorData_SemDuplicacaoDeIntervalos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Criar prestador
	prestador := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(0))

	// Adicionar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: dataFutura,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
			{HoraInicio: "14:00", HoraFim: "18:00"},
		},
	}
	SetupPutAgendaRequest(router, prestador.ID, agendaInput)

	// Buscar prestadores disponíveis
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	require.Len(t, resp.Data, 1)
	require.Len(t, resp.Data[0].Agenda, 1)

	// Validação crítica: NÃO deve ter intervalos duplicados
	intervalos := resp.Data[0].Agenda[0].Intervalos
	assert.Len(t, intervalos, 2, "Deve ter exatamente 2 intervalos únicos")

	// Validar IDs únicos
	idsIntervalos := make(map[string]bool)
	for _, intervalo := range intervalos {
		assert.False(t, idsIntervalos[intervalo.ID], "ID do intervalo %s está duplicado", intervalo.ID)
		idsIntervalos[intervalo.ID] = true
	}
}

func TestGetPrestadoresPorData_DataHoje(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	dataHoje := time.Now().Format("2006-01-02")

	// Criar prestador com agenda hoje
	prestador := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(0))

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: dataHoje,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "09:00", HoraFim: "12:00"},
		},
	}
	rrAgenda := SetupPutAgendaRequest(router, prestador.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)

	// Buscar prestadores disponíveis hoje
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", dataHoje)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.GreaterOrEqual(t, resp.Total, 1)
	assert.NotEmpty(t, resp.Data)
}

func TestGetPrestadoresPorData_PageZeroAjustaParaUm(t *testing.T) {
	router, _ := SetupPostPrestador()

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&page=0&limit=10", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, 1, resp.Page)
}

func TestGetPrestadoresPorData_LimitZeroAjustaParaDez(t *testing.T) {
	router, _ := SetupPostPrestador()

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&page=1&limit=0", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, 10, resp.Limit)
}

func TestGetPrestadoresPorData_LimiteMaximo(t *testing.T) {
	router, _ := SetupPostPrestador()

	dataFutura := time.Now().AddDate(0, 0, 5).Format("2006-01-02")

	// Limite maior que 100 é rejeitado
	url := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&limit=150", dataFutura)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "Limit")

	// Limite exatamente 100 funciona
	url2 := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s&limit=100", dataFutura)
	req2, _ := http.NewRequest(http.MethodGet, url2, nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr2.Body.Bytes(), &resp)

	assert.Equal(t, 100, resp.Limit)
}

func TestGetPrestadoresPorData_DiferentesDatas(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	data1 := time.Now().AddDate(0, 0, 5).Format("2006-01-02")
	data2 := time.Now().AddDate(0, 0, 10).Format("2006-01-02")

	// Criar prestadores
	p1 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(0))
	p2 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(1))
	p3 := CriarPrestadorValido(t, router, catalogoResp.ID, gerarCPFValido(2))

	// p1 e p2 na data1
	agenda1 := request_prestador.AgendaDiariaRequest{
		Data: data1,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "09:00", HoraFim: "12:00"},
		},
	}
	SetupPutAgendaRequest(router, p1.ID, agenda1)
	SetupPutAgendaRequest(router, p2.ID, agenda1)

	// p3 na data2
	agenda2 := request_prestador.AgendaDiariaRequest{
		Data: data2,
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "09:00", HoraFim: "12:00"},
		},
	}
	SetupPutAgendaRequest(router, p3.ID, agenda2)

	// Buscar na data1
	url1 := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", data1)
	req1, _ := http.NewRequest(http.MethodGet, url1, nil)
	rr1 := httptest.NewRecorder()
	router.ServeHTTP(rr1, req1)

	var resp1 response_prestador.PrestadorListResponse
	json.Unmarshal(rr1.Body.Bytes(), &resp1)

	assert.Len(t, resp1.Data, 2)
	assert.Equal(t, 2, resp1.Total)

	// Buscar na data2
	url2 := fmt.Sprintf("/api/v1/prestadores/disponiveis?data=%s", data2)
	req2, _ := http.NewRequest(http.MethodGet, url2, nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	var resp2 response_prestador.PrestadorListResponse
	json.Unmarshal(rr2.Body.Bytes(), &resp2)

	assert.Len(t, resp2.Data, 1)
	assert.Equal(t, 1, resp2.Total)
	assert.Equal(t, p3.ID, resp2.Data[0].ID)
}