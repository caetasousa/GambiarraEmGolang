package teste

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/request"
	"meu-servico-agenda/internal/adapters/http/catalogo/response"
	"meu-servico-agenda/internal/adapters/http/prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func SetupPostPrestador() (*gin.Engine, port.PrestadorRepositorio) {
	gin.SetMode(gin.TestMode)

	prestadorRepo := repository.NovoFakePrestadorRepositorio()
	catalogoRepo := repository.NovoCatalogoFakeRepo()
	agendaRepo := repository.NovoFakeAgendaDiariaRepositorio()
	cadastroService := service.NovoCatalogoService(catalogoRepo)

	prestadorService := service.NovaPrestadorService(
		prestadorRepo,
		catalogoRepo,
		agendaRepo,
	)

	prestadorController := prestador.NovoPrestadorController(prestadorService)
	catalogoController := catalogo.NovoCatalogoController(cadastroService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/prestadores", prestadorController.PostPrestador)
		apiV1.GET("/prestadores/:id", prestadorController.GetPrestador)
		apiV1.PUT("/prestadores/:id/agenda", prestadorController.PutAgenda)

		apiV1.POST("/catalogos", catalogoController.PostPrestador)
	}

	return router, prestadorRepo
}

func SetupPostPrestadorRequest(router *gin.Engine, input request_prestador.PrestadorRequest) *httptest.ResponseRecorder {

	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/prestadores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupGetPrestadorRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {

	url := "/api/v1/prestadores/" + id
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupPutAgendaRequest(router *gin.Engine, prestadorID string, input request_prestador.AgendaDiariaRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)

	url := "/api/v1/prestadores/" + prestadorID + "/agenda"
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func TestPostPrestador_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()

	catalogoInput := request.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
	}

	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Email:       "joao@email.com",
		Cpf:         "04423258196",
		Telefone:    "62999677481",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrPrestador.Code)
}

func TestPostPrestador_FalhaCatalogoInexistente(t *testing.T) {
	router, _ := SetupPostPrestador()

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		CatalogoIDs: []string{"catalogo-inexistente"},
	}

	rr := SetupPostPrestadorRequest(router, prestadorInput)

	// O service retorna erro, controller responde com Conflict (409)
	require.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	// Verifica se a mensagem de erro contém a informação do catálogo inexistente
	require.Contains(t, rr.Body.String(), "não existe")
}

func TestGetPrestador_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo válido
	catalogoInput := request.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
	}
	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador usando o ID do catálogo
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	// 3️⃣ Buscar prestador criado
	var resp map[string]interface{}
	_ = json.Unmarshal(rrCreate.Body.Bytes(), &resp)
	id := resp["id"].(string)

	rrGet := SetupGetPrestadorRequest(router, id)
	require.Equal(t, http.StatusOK, rrGet.Code)
}

func TestGetPrestador_UsuarioExistente(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo válido
	catalogoInput := request.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
	}
	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador usando o ID do catálogo
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	ccCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)
	require.Equal(t, http.StatusConflict, ccCreate.Code)
}

func TestGetPrestador_NaoEncontrado(t *testing.T) {
	router, _ := SetupPostPrestador()
	rr := SetupGetPrestadorRequest(router, "id-inexistente")
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func CriarPrestadorValidoParaTeste(t *testing.T) (*gin.Engine, domain.Prestador, port.PrestadorRepositorio) {
	router, repo := SetupPostPrestador()

	// 1️⃣ Criar catálogo
	catalogoInput := request.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
	}
	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrPrestador.Code)

	var prestadorResp domain.Prestador
	err = json.Unmarshal(rrPrestador.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	return router, prestadorResp, repo
}
func TestPutAgenda_Sucesso(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2025-01-03",
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
		Data: "2025-01-03",
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
		Data: "2025-01-03",
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
		Data: "2025-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{
				HoraInicio: "18:00",
				HoraFim:    "08:00",
			},
		},
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	// valida exatamente o erro do domínio
	require.Contains(t, rr.Body.String(), domain.ErrIntervaloHorarioInvalido.Error())
}

func TestPutAgenda_AgendaSemIntervalos(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data:       "2025-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{}, // vazio
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), domain.ErrAgendaSemIntervalos.Error())
}
