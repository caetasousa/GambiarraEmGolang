package teste

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrestadores_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar 3 prestadores
	for i := 0; i < 3; i++ {
		prestadorInput := request_prestador.PrestadorRequest{
			Nome:        fmt.Sprintf("Prestador %d", i+1),
			Cpf:         gerarCPFValido(i),
			Email:       fmt.Sprintf("prestador%d@email.com", i+1),
			Telefone:    fmt.Sprintf("6299967748%d", i+1),
			ImagemUrl:   fmt.Sprintf("https://exemplo.com/img%d.jpg", i+1),
			CatalogoIDs: []string{catalogoResp.ID},
		}
		rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
		require.Equal(t, http.StatusCreated, rrPrestador.Code)
	}

	// ✅ Listar ativos
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?ativo=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 3)
	assert.Equal(t, 3, resp.Total)
}

func TestGetPrestadores_ComPaginacao(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar 15 prestadores
	cpfsValidos := []string{
		"91663575002", "10886737087", "55964702015", "35212899079", "42864297094",
		"44187423010", "45537518015", "23204646033", "86306650091", "40933461003",
		"33935391080", "29466173006", "32886059021", "77487008002", "88992049005",
	}

	for i := 0; i < 15; i++ {
		prestadorInput := request_prestador.PrestadorRequest{
			Nome:        fmt.Sprintf("Prestador %d", i+1),
			Cpf:         cpfsValidos[i],
			Email:       fmt.Sprintf("prestador%d@email.com", i+1),
			Telefone:    fmt.Sprintf("62999%06d", i+1),
			ImagemUrl:   fmt.Sprintf("https://exemplo.com/img%d.jpg", i+1),
			CatalogoIDs: []string{catalogoResp.ID},
		}
		SetupPostPrestadorRequest(router, prestadorInput)
	}

	// ✅ Página 1 com filtro ativo
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=1&limit=5&ativo=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 5)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 5, resp.Limit)
	assert.Equal(t, 15, resp.Total)

	// ✅ Página 2
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=2&limit=5&ativo=true", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp2 response_prestador.PrestadorListResponse
	json.Unmarshal(rr2.Body.Bytes(), &resp2)

	assert.Len(t, resp2.Data, 5)
	assert.Equal(t, 2, resp2.Page)
	assert.Equal(t, 5, resp2.Limit)
	assert.Equal(t, 15, resp2.Total)

	// ✅ Página 3
	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=3&limit=5&ativo=true", nil)
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	require.Equal(t, http.StatusOK, rr3.Code)

	var resp3 response_prestador.PrestadorListResponse
	json.Unmarshal(rr3.Body.Bytes(), &resp3)

	assert.Len(t, resp3.Data, 5)
	assert.Equal(t, 3, resp3.Page)
	assert.Equal(t, 15, resp3.Total)

	// ✅ Página 4 (vazia)
	req4, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=4&limit=5&ativo=true", nil)
	rr4 := httptest.NewRecorder()
	router.ServeHTTP(rr4, req4)

	require.Equal(t, http.StatusOK, rr4.Code)

	var resp4 response_prestador.PrestadorListResponse
	json.Unmarshal(rr4.Body.Bytes(), &resp4)

	assert.Len(t, resp4.Data, 0)
	assert.Equal(t, 4, resp4.Page)
	assert.Equal(t, 15, resp4.Total)
}

func TestGetPrestadores_ListaVazia(t *testing.T) {
	router, _ := SetupPostPrestador()

	// ✅ Buscar ativos em lista vazia
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?ativo=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 0)
	assert.Equal(t, 0, resp.Total)
}

func TestGetPrestadores_ParametrosInvalidos(t *testing.T) {
	router, _ := SetupPostPrestador()

	testCases := []struct {
		name           string
		url            string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Page negativo",
			url:            "/api/v1/prestadores?page=-1&ativo=true", // ✅ Adicionado ativo
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Page",
		},
		{
			name:           "Limit negativo",
			url:            "/api/v1/prestadores?limit=-5&ativo=true", // ✅ Adicionado ativo
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Limit",
		},
		{
			name:           "Page não numérico",
			url:            "/api/v1/prestadores?page=abc&ativo=true", // ✅ Adicionado ativo
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid syntax",
		},
		{
			name:           "Limit não numérico",
			url:            "/api/v1/prestadores?limit=xyz&ativo=true", // ✅ Adicionado ativo
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid syntax",
		},
		{
			name:           "Sem parâmetro ativo", // ✅ NOVO
			url:            "/api/v1/prestadores",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "ativo",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, tc.url, nil)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			require.Equal(t, tc.expectedStatus, rr.Code)
			if tc.expectedMsg != "" {
				require.Contains(t, rr.Body.String(), tc.expectedMsg)
			}
		})
	}
}

func TestGetPrestadores_LimiteMaximo(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar 15 prestadores
	cpfsValidos := []string{
		"91663575002", "10886737087", "55964702015", "35212899079", "42864297094",
		"44187423010", "45537518015", "23204646033", "86306650091", "40933461003",
		"33935391080", "29466173006", "32886059021", "77487008002", "88992049005",
	}

	for i := 0; i < 15; i++ {
		prestadorInput := request_prestador.PrestadorRequest{
			Nome:        fmt.Sprintf("Prestador %d", i+1),
			Cpf:         cpfsValidos[i],
			Email:       fmt.Sprintf("prestador%d@email.com", i+1),
			Telefone:    fmt.Sprintf("62999%06d", i+1),
			ImagemUrl:   fmt.Sprintf("https://exemplo.com/img%d.jpg", i+1),
			CatalogoIDs: []string{catalogoResp.ID},
		}
		rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
		require.Equal(t, http.StatusCreated, rrPrestador.Code)
	}

	// Teste 1: Limit maior que 100 é rejeitado
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?limit=150&ativo=true", nil) // ✅ Adicionado ativo
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "Limit")

	// Teste 2: Limit exatamente 100 funciona
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?limit=100&ativo=true", nil) // ✅ Adicionado ativo
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr2.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 100, resp.Limit)
	assert.Equal(t, 15, resp.Total)
	assert.Len(t, resp.Data, 15)
}

func TestGetPrestadores_ComAgendasECatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João Silva",
		Cpf:         "91663575002",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrPrestador.Code)

	var prestadorResp response_prestador.PrestadorResponse
	err := json.Unmarshal(rrPrestador.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	// Adicionar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	rrAgenda := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)

	// ✅ Listar prestadores ativos
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?ativo=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Len(t, resp.Data, 1)

	prestador := resp.Data[0]
	assert.NotEmpty(t, prestador.Catalogo)
	assert.Len(t, prestador.Catalogo, 1)
	assert.Equal(t, "Corte de Cabelo", prestador.Catalogo[0].Nome)

	assert.NotEmpty(t, prestador.Agenda)
	assert.Len(t, prestador.Agenda, 1)
	assert.Equal(t, "2030-01-03", prestador.Agenda[0].Data)
	assert.Len(t, prestador.Agenda[0].Intervalos, 1)
}

func TestGetPrestadores_PageZeroAjustaParaUm(t *testing.T) {
	router, _ := SetupPostPrestador()

	// ✅ Adicionado ativo
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=0&limit=10&ativo=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 1, resp.Page)
}

func TestGetPrestadores_LimitZeroAjustaParaDez(t *testing.T) {
	router, _ := SetupPostPrestador()

	// ✅ Adicionado ativo
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=1&limit=0&ativo=true", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 10, resp.Limit)
}

// ✅ NOVOS TESTES - Filtro de status

func TestGetPrestadores_FiltroAtivos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar 2 prestadores ativos
	CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")
	CriarPrestadorValido(t, router, catalogoResp.ID, "10886737087")

	// Criar 1 prestador e inativar
	prestador3 := CriarPrestadorValido(t, router, catalogoResp.ID, "55964702015")
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/prestadores/"+prestador3.ID+"/inativar", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Buscar só ativos
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?ativo=true", nil)
	rrGet := httptest.NewRecorder()
	router.ServeHTTP(rrGet, reqGet)

	require.Equal(t, http.StatusOK, rrGet.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rrGet.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 2, "Deve retornar apenas 2 ativos")
	assert.Equal(t, 2, resp.Total)

	// Verificar que todos são ativos
	for _, p := range resp.Data {
		assert.True(t, p.Ativo)
	}
}

func TestGetPrestadores_FiltroInativos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar 3 prestadores
	p1 := CriarPrestadorValido(t, router, catalogoResp.ID, "04423258196")
	p2 := CriarPrestadorValido(t, router, catalogoResp.ID, "10886737087")
	CriarPrestadorValido(t, router, catalogoResp.ID, "55964702015")

	// Inativar 2 deles
	for _, id := range []string{p1.ID, p2.ID} {
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/prestadores/"+id+"/inativar", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		require.Equal(t, http.StatusNoContent, rr.Code)
	}

	// Buscar só inativos
	reqGet, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?ativo=false", nil)
	rrGet := httptest.NewRecorder()
	router.ServeHTTP(rrGet, reqGet)

	require.Equal(t, http.StatusOK, rrGet.Code)

	var resp response_prestador.PrestadorListResponse
	json.Unmarshal(rrGet.Body.Bytes(), &resp)

	assert.Len(t, resp.Data, 2, "Deve retornar apenas 2 inativos")
	assert.Equal(t, 2, resp.Total)

	// Verificar que todos são inativos
	for _, p := range resp.Data {
		assert.False(t, p.Ativo)
	}
}

func TestGetPrestador_SemDuplicacaoDeIntervalos(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// Criar prestador
	prestadorResp := CriarPrestadorValido(t, router, catalogoResp.ID, "91663575002")

	// Adicionar agenda com 2 intervalos
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
			{HoraInicio: "14:00", HoraFim: "18:00"},
		},
	}
	rrAgenda := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)

	// Buscar prestador por ID
	rrGet := SetupGetPrestadorRequest(router, prestadorResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	var prestador response_prestador.PrestadorResponse
	json.Unmarshal(rrGet.Body.Bytes(), &prestador)

	require.Len(t, prestador.Agenda, 1, "Deve ter 1 agenda")
	
	// ✅ Validação crítica: NÃO deve ter intervalos duplicados
	assert.Len(t, prestador.Agenda[0].Intervalos, 2, "Deve ter exatamente 2 intervalos únicos")

	// ✅ Validar IDs únicos
	idsIntervalos := make(map[string]bool)
	for _, intervalo := range prestador.Agenda[0].Intervalos {
		assert.False(t, idsIntervalos[intervalo.ID], "ID do intervalo %s está duplicado", intervalo.ID)
		idsIntervalos[intervalo.ID] = true
	}
}