package teste

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/response_catalogo"
	"meu-servico-agenda/internal/adapters/http/prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/request_prestador"
	"meu-servico-agenda/internal/adapters/http/prestador/response_prestador"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

		apiV1.POST("/catalogos", catalogoController.PostCatalogo)
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

func SetupPostCatalogoRequest(router *gin.Engine, input request_catalogo.CatalogoRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalogos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
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

	catalogoInput := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}

	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response_catalogo.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Email:       "joao@email.com",
		Cpf:         "04423258196",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
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
		ImagemUrl:   "https://exemplo.com/img1.jpg",
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
	catalogoInput := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img-catalogo.jpg",
	}
	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response_catalogo.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador usando o ID do catálogo
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria Silva",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		ImagemUrl:   "https://exemplo.com/img-prestador.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	// 3️⃣ Extrair ID do prestador criado
	var createResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)
	require.NotEmpty(t, createResp.ID)

	// 4️⃣ Buscar prestador criado
	rrGet := SetupGetPrestadorRequest(router, createResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	// 5️⃣ Validar todos os dados retornados
	var prestadorResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrGet.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	// Validar dados do prestador
	assert.Equal(t, createResp.ID, prestadorResp.ID, "ID do prestador deve ser igual")
	assert.Equal(t, "Maria Silva", prestadorResp.Nome, "Nome do prestador deve ser igual")
	assert.Equal(t, "04423258196", prestadorResp.Cpf, "CPF do prestador deve ser igual")
	assert.Equal(t, "maria@email.com", prestadorResp.Email, "Email do prestador deve ser igual")
	assert.Equal(t, "62999677482", prestadorResp.Telefone, "Telefone do prestador deve ser igual")
	assert.True(t, prestadorResp.Ativo, "Prestador deve estar ativo")
	assert.Equal(t, "https://exemplo.com/img-prestador.jpg", prestadorResp.ImagemUrl, "URL da imagem do prestador deve ser igual")
	assert.NotEmpty(t, prestadorResp.ImagemUrl, "URL da imagem do prestador não deve estar vazia")

	// Validar catálogos
	require.NotEmpty(t, prestadorResp.Catalogo, "Prestador deve ter pelo menos um catálogo")
	assert.Len(t, prestadorResp.Catalogo, 1, "Prestador deve ter exatamente 1 catálogo")

	catalogoRetornado := prestadorResp.Catalogo[0]
	assert.Equal(t, catalogoResp.ID, catalogoRetornado.ID, "ID do catálogo deve ser igual")
	assert.Equal(t, "Corte de Cabelo", catalogoRetornado.Nome, "Nome do catálogo deve ser igual")
	assert.Equal(t, 30, catalogoRetornado.DuracaoPadrao, "Duração padrão do catálogo deve ser igual")
	assert.Equal(t, 3500, catalogoRetornado.Preco, "Preço do catálogo deve ser igual")
	assert.Equal(t, "Beleza", catalogoRetornado.Categoria, "Categoria do catálogo deve ser igual")
	assert.Equal(t, "https://exemplo.com/img-catalogo.jpg", catalogoRetornado.ImagemUrl, "URL da imagem do catálogo deve ser igual")
	assert.NotEmpty(t, catalogoRetornado.ImagemUrl, "URL da imagem do catálogo não deve estar vazia")

	// Validar que agenda está vazia (já que não foi criada nenhuma)
	assert.NotNil(t, prestadorResp.Agenda, "Agenda não deve ser nil")
	assert.Empty(t, prestadorResp.Agenda, "Agenda deve estar vazia")
}

func TestGetPrestador_UsuarioExistente(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo válido
	catalogoInput := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response_catalogo.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador usando o ID do catálogo
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
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
	catalogoInput := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrCatalogo := SetupPostCatalogoRequest(router, catalogoInput)
	require.Equal(t, http.StatusCreated, rrCatalogo.Code)

	var catalogoResp response_catalogo.CatalogoResponse
	err := json.Unmarshal(rrCatalogo.Body.Bytes(), &catalogoResp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
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

	// valida exatamente o erro do domínio
	require.Contains(t, rr.Body.String(), domain.ErrIntervaloHorarioInvalido.Error())
}

func TestPutAgenda_AgendaSemIntervalos(t *testing.T) {
	router, prestadorResp, _ := CriarPrestadorValidoParaTeste(t)

	agendaInput := request_prestador.AgendaDiariaRequest{
		Data:       "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{}, // vazio
	}

	rr := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), domain.ErrAgendaSemIntervalos.Error())
}
