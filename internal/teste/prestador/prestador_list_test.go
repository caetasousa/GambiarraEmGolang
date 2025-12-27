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

	// Criar 3 prestadores com CPFs válidos
	cpfsValidos := []string{"91663575002", "10886737087", "55964702015"}

	for i := 0; i < 3; i++ {
		prestadorInput := request_prestador.PrestadorRequest{
			Nome:        fmt.Sprintf("Prestador %d", i+1),
			Cpf:         cpfsValidos[i],
			Email:       fmt.Sprintf("prestador%d@email.com", i+1),
			Telefone:    fmt.Sprintf("6299967748%d", i+1),
			ImagemUrl:   fmt.Sprintf("https://exemplo.com/img%d.jpg", i+1),
			CatalogoIDs: []string{catalogoResp.ID},
		}
		rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
		require.Equal(t, http.StatusCreated, rrPrestador.Code)
	}

	// Listar prestadores
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 3, "Deve retornar 3 prestadores")
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 3, resp.Total)

	assert.NotEmpty(t, resp.Data[0].ID)
	assert.NotEmpty(t, resp.Data[0].Nome)
	assert.NotEmpty(t, resp.Data[0].Catalogo)
}

func TestGetPrestadores_ComPaginacao(t *testing.T) {
	router, _ := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

	// CPFs válidos
	cpfsValidos := []string{
		"91663575002", "10886737087", "55964702015", "35212899079", "42864297094",
		"44187423010", "45537518015", "23204646033", "86306650091", "40933461003",
		"33935391080", "29466173006", "32886059021", "77487008002", "88992049005",
	}

	// Criar 15 prestadores
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
		if rrPrestador.Code != http.StatusCreated {
			t.Fatalf("Erro ao criar prestador %d: %s", i+1, rrPrestador.Body.String())
		}
	}

	// Página 1 com 5 itens
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=1&limit=5", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 5)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 5, resp.Limit)
	assert.Equal(t, 15, resp.Total)

	// Página 2 com 5 itens
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=2&limit=5", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp2 response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	require.NoError(t, err)

	assert.Len(t, resp2.Data, 5)
	assert.Equal(t, 2, resp2.Page)

	// Página 4 (vazia)
	req4, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=4&limit=5", nil)
	rr4 := httptest.NewRecorder()
	router.ServeHTTP(rr4, req4)

	require.Equal(t, http.StatusOK, rr4.Code)

	var resp4 response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr4.Body.Bytes(), &resp4)
	require.NoError(t, err)

	assert.Len(t, resp4.Data, 0)
	assert.Equal(t, 4, resp4.Page)
	assert.Equal(t, 15, resp4.Total)
}

func TestGetPrestadores_ListaVazia(t *testing.T) {
	router, _ := SetupPostPrestador()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 0)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)
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
			url:            "/api/v1/prestadores?page=-1",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Page",
		},
		{
			name:           "Limit negativo",
			url:            "/api/v1/prestadores?limit=-5",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Limit",
		},
		{
			name:           "Page não numérico",
			url:            "/api/v1/prestadores?page=abc",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid syntax",
		},
		{
			name:           "Limit não numérico",
			url:            "/api/v1/prestadores?limit=xyz",
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "invalid syntax",
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
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?limit=150", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "Limit")

	// Teste 2: Limit exatamente 100 funciona
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?limit=100", nil)
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

	// Listar prestadores
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores", nil)
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

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=0&limit=10", nil)
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

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=1&limit=0", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 10, resp.Limit)
}
