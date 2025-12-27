package teste

import (
	"bytes"
	"encoding/json"
	"fmt"
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

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	prestadorRepo := repository.NovoFakePrestadorRepositorio(catalogoRepo)
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
		apiV1.GET("/prestadores", prestadorController.GetPreestadores)
		apiV1.PUT("/prestadores/:id/agenda", prestadorController.PutAgenda)
		apiV1.PUT("/prestadores/:id", prestadorController.UpdatePrestador)

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

func SetupPutPrestadorRequest(router *gin.Engine, id string, input request_prestador.PrestadorUpdateRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)

	url := "/api/v1/prestadores/" + id
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupGetPrestadoresRequest(router *gin.Engine, page, limit int) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/prestadores?page=%d&limit=%d", page, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func gerarCPFValido(seed int) string {
	// CPFs válidos fornecidos pelo usuário
	cpfsValidos := []string{
		"91663575002",
		"10886737087",
		"55964702015",
		"35212899079",
		"42864297094",
		"44187423010",
		"45537518015",
		"23204646033",
		"86306650091",
		"40933461003",
		"33935391080",
		"29466173006",
		"32886059021",
		"77487008002",
		"88992049005",
		"06724785014",
		"90035795042",
		"39308388001",
		"86148883090",
		"26345031054",
		"60314052020",
		"25176594005",
	}

	return cpfsValidos[seed%len(cpfsValidos)]
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

func TestUpdatePrestador_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo
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

	// 2️⃣ Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "Maria Silva",
		Cpf:         "04423258196",
		Email:       "maria@email.com",
		Telefone:    "62999677482",
		ImagemUrl:   "https://exemplo.com/img-original.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)

	// 3️⃣ Atualizar prestador
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "Maria Santos Atualizada",
		Email:       "maria.santos@email.com",
		Telefone:    "62988888888",
		ImagemUrl:   "https://exemplo.com/img-atualizada.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrUpdate := SetupPutPrestadorRequest(router, createResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rrUpdate.Code)

	// 4️⃣ Buscar prestador atualizado para validar
	rrGet := SetupGetPrestadorRequest(router, createResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	var prestadorAtualizado response_prestador.PrestadorResponse
	err = json.Unmarshal(rrGet.Body.Bytes(), &prestadorAtualizado)
	require.NoError(t, err)

	// Validações
	assert.Equal(t, createResp.ID, prestadorAtualizado.ID, "ID não deve mudar")
	assert.Equal(t, "04423258196", prestadorAtualizado.Cpf, "CPF não deve mudar")
	assert.Equal(t, "Maria Santos Atualizada", prestadorAtualizado.Nome)
	assert.Equal(t, "maria.santos@email.com", prestadorAtualizado.Email)
	assert.Equal(t, "62988888888", prestadorAtualizado.Telefone)
	assert.Equal(t, "https://exemplo.com/img-atualizada.jpg", prestadorAtualizado.ImagemUrl)
}

// TestUpdatePrestador_PrestadorNaoEncontrado testa atualização de prestador inexistente
func TestUpdatePrestador_PrestadorNaoEncontrado(t *testing.T) {
	router, _ := SetupPostPrestador()

	// Criar catálogo válido
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

	// Tentar atualizar prestador inexistente
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}

	rr := SetupPutPrestadorRequest(router, "id-inexistente", updateInput)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "prestador não encontrado")
}

// TestUpdatePrestador_CatalogoInexistente testa atualização com catálogo que não existe
func TestUpdatePrestador_CatalogoInexistente(t *testing.T) {
	router, _ := SetupPostPrestador()

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
		Nome:        "João Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)

	// 3️⃣ Tentar atualizar com catálogo inexistente
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva Atualizado",
		Email:       "joao.novo@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img2.jpg",
		CatalogoIDs: []string{"catalogo-inexistente"},
	}

	rr := SetupPutPrestadorRequest(router, createResp.ID, updateInput)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Contains(t, rr.Body.String(), "catálogo")
	require.Contains(t, rr.Body.String(), "não existe")
}

// TestUpdatePrestador_AtualizarCatalogos testa atualização dos catálogos associados
func TestUpdatePrestador_AtualizarCatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar dois catálogos
	catalogo1Input := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrCatalogo1 := SetupPostCatalogoRequest(router, catalogo1Input)
	require.Equal(t, http.StatusCreated, rrCatalogo1.Code)

	var catalogo1Resp response_catalogo.CatalogoResponse
	err := json.Unmarshal(rrCatalogo1.Body.Bytes(), &catalogo1Resp)
	require.NoError(t, err)

	catalogo2Input := request_catalogo.CatalogoRequest{
		Nome:          "Barba",
		DuracaoPadrao: 20,
		Preco:         2000.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img2.jpg",
	}
	rrCatalogo2 := SetupPostCatalogoRequest(router, catalogo2Input)
	require.Equal(t, http.StatusCreated, rrCatalogo2.Code)

	var catalogo2Resp response_catalogo.CatalogoResponse
	err = json.Unmarshal(rrCatalogo2.Body.Bytes(), &catalogo2Resp)
	require.NoError(t, err)

	// 2️⃣ Criar prestador com primeiro catálogo
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogo1Resp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)

	// 3️⃣ Atualizar prestador com ambos os catálogos
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogo1Resp.ID, catalogo2Resp.ID},
	}
	rrUpdate := SetupPutPrestadorRequest(router, createResp.ID, updateInput)
	require.Equal(t, http.StatusNoContent, rrUpdate.Code)

	// 4️⃣ Buscar e validar que tem 2 catálogos
	rrGet := SetupGetPrestadorRequest(router, createResp.ID)
	require.Equal(t, http.StatusOK, rrGet.Code)

	var prestadorAtualizado response_prestador.PrestadorResponse
	err = json.Unmarshal(rrGet.Body.Bytes(), &prestadorAtualizado)
	require.NoError(t, err)

	assert.Len(t, prestadorAtualizado.Catalogo, 2, "Deve ter 2 catálogos")
}

// TestUpdatePrestador_DadosInvalidos testa validação de dados inválidos
func TestUpdatePrestador_DadosInvalidos(t *testing.T) {
	router, _ := SetupPostPrestador()

	// Criar catálogo válido
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

	// Criar prestador
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)

	// Testes de validação
	testCases := []struct {
		name        string
		input       request_prestador.PrestadorUpdateRequest
		expectedMsg string
	}{
		{
			name: "Nome muito curto",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "Jo", // menos de 3 caracteres
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Nome",
		},
		{
			name: "Email inválido",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "João Silva",
				Email:       "email-invalido", // sem @
				Telefone:    "62999677481",
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Email",
		},
		{
			name: "Telefone muito curto",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "João Silva",
				Email:       "joao@email.com",
				Telefone:    "1234567", // menos de 8 caracteres
				ImagemUrl:   "https://exemplo.com/img1.jpg",
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "Telefone",
		},
		{
			name: "URL inválida",
			input: request_prestador.PrestadorUpdateRequest{
				Nome:        "João Silva",
				Email:       "joao@email.com",
				Telefone:    "62999677481",
				ImagemUrl:   "url-invalida", // não é URL
				CatalogoIDs: []string{catalogoResp.ID},
			},
			expectedMsg: "ImagemUrl",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rr := SetupPutPrestadorRequest(router, createResp.ID, tc.input)
			require.Equal(t, http.StatusBadRequest, rr.Code)
			require.Contains(t, rr.Body.String(), tc.expectedMsg)
		})
	}
}

// TestUpdatePrestador_RemoverTodosCatalogos testa remoção de todos os catálogos (se permitido)
func TestUpdatePrestador_SemCatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()

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
		Nome:        "João Silva",
		Cpf:         "04423258196",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoResp.ID},
	}
	rrCreate := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrCreate.Code)

	var createResp response_prestador.PrestadorResponse
	err = json.Unmarshal(rrCreate.Body.Bytes(), &createResp)
	require.NoError(t, err)

	// 3️⃣ Tentar atualizar sem catálogos (array vazio)
	updateInput := request_prestador.PrestadorUpdateRequest{
		Nome:        "João Silva",
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{}, // vazio
	}

	rr := SetupPutPrestadorRequest(router, createResp.ID, updateInput)

	// Dependendo da regra de negócio, pode ser 400 ou aceitar
	// Ajuste conforme sua regra
	require.Equal(t, http.StatusNoContent, rr.Code)
}

func TestGetPrestadores_Sucesso(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo
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

	// 2️⃣ Criar 3 prestadores com CPFs válidos
	cpfsValidos := []string{
		"91663575002",
		"10886737087",
		"55964702015",
	}

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

	// 3️⃣ Listar prestadores
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 3, "Deve retornar 3 prestadores")
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 3, resp.Total)

	assert.NotEmpty(t, resp.Data[0].ID)
	assert.NotEmpty(t, resp.Data[0].Nome)
	assert.NotEmpty(t, resp.Data[0].Email)
	assert.NotEmpty(t, resp.Data[0].Cpf)
	assert.NotEmpty(t, resp.Data[0].Telefone)
	assert.True(t, resp.Data[0].Ativo)
	assert.NotEmpty(t, resp.Data[0].Catalogo)
}

func TestGetPrestadores_ComPaginacao(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo
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

	// 2️⃣ CPFs válidos (fornecidos pelo usuário)
	cpfsValidos := []string{
		"91663575002", "10886737087", "55964702015", "35212899079", "42864297094",
		"44187423010", "45537518015", "23204646033", "86306650091", "40933461003",
		"33935391080", "29466173006", "32886059021", "77487008002", "88992049005",
	}

	// 3️⃣ Criar 15 prestadores
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
			t.Fatalf("Erro ao criar prestador %d com CPF %s: Status=%d, Body=%s",
				i+1, prestadorInput.Cpf, rrPrestador.Code, rrPrestador.Body.String())
		}
	}

	// 4️⃣ Página 1 com 5 itens
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=1&limit=5", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 5, "Deve retornar 5 prestadores")
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 5, resp.Limit)
	assert.Equal(t, 15, resp.Total)

	// 5️⃣ Página 2 com 5 itens
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=2&limit=5", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp2 response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr2.Body.Bytes(), &resp2)
	require.NoError(t, err)

	assert.Len(t, resp2.Data, 5, "Deve retornar 5 prestadores")
	assert.Equal(t, 2, resp2.Page)
	assert.Equal(t, 5, resp2.Limit)
	assert.Equal(t, 15, resp2.Total)

	// 6️⃣ Página 3 com 5 itens
	req3, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=3&limit=5", nil)
	rr3 := httptest.NewRecorder()
	router.ServeHTTP(rr3, req3)

	require.Equal(t, http.StatusOK, rr3.Code)

	var resp3 response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr3.Body.Bytes(), &resp3)
	require.NoError(t, err)

	assert.Len(t, resp3.Data, 5, "Deve retornar 5 prestadores")
	assert.Equal(t, 3, resp3.Page)
	assert.Equal(t, 15, resp3.Total)

	// 7️⃣ Página 4 (vazia)
	req4, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=4&limit=5", nil)
	rr4 := httptest.NewRecorder()
	router.ServeHTTP(rr4, req4)

	require.Equal(t, http.StatusOK, rr4.Code)

	var resp4 response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr4.Body.Bytes(), &resp4)
	require.NoError(t, err)

	assert.Len(t, resp4.Data, 0, "Deve retornar 0 prestadores")
	assert.Equal(t, 4, resp4.Page)
	assert.Equal(t, 15, resp4.Total)
}

// TestGetPrestadores_ListaVazia testa listagem quando não há prestadores
func TestGetPrestadores_ListaVazia(t *testing.T) {
	router, _ := SetupPostPrestador()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Len(t, resp.Data, 0, "Deve retornar lista vazia")
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 10, resp.Limit)
	assert.Equal(t, 0, resp.Total)
}

// TestGetPrestadores_ParametrosInvalidos testa validação de parâmetros
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
			expectedMsg:    "Page", // ✅ Vai aparecer na mensagem de validação
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

			require.Equal(t, tc.expectedStatus, rr.Code, "Response: %s", rr.Body.String())
			if tc.expectedMsg != "" {
				require.Contains(t, rr.Body.String(), tc.expectedMsg)
			}
		})
	}
}

// TestGetPrestadores_LimiteMaximo testa limite máximo de 100
func TestGetPrestadores_LimiteMaximo(t *testing.T) {
	router, _ := SetupPostPrestador()

	// Criar catálogo
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

	// ✅ Teste 1: Limit maior que 100 é rejeitado no binding
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?limit=150", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code, "Limit > 100 deve retornar 400")
	require.Contains(t, rr.Body.String(), "Limit", "Deve mencionar Limit na mensagem de erro")

	// ✅ Teste 2: Limit exatamente 100 funciona
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?limit=100", nil)
	rr2 := httptest.NewRecorder()
	router.ServeHTTP(rr2, req2)

	require.Equal(t, http.StatusOK, rr2.Code)

	var resp response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr2.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 100, resp.Limit, "Limit deve ser 100")
	assert.Equal(t, 15, resp.Total)
	assert.Len(t, resp.Data, 15, "Deve retornar todos os 15 prestadores")
}

// TestGetPrestadores_ComAgendasECatalogos testa que retorna prestadores com dados relacionados
func TestGetPrestadores_ComAgendasECatalogos(t *testing.T) {
	router, _ := SetupPostPrestador()

	// 1️⃣ Criar catálogo
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

	// 2️⃣ Criar prestador
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
	err = json.Unmarshal(rrPrestador.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	// 3️⃣ Adicionar agenda
	agendaInput := request_prestador.AgendaDiariaRequest{
		Data: "2030-01-03",
		Intervalos: []request_prestador.IntervaloDiarioRequest{
			{HoraInicio: "08:00", HoraFim: "12:00"},
		},
	}
	rrAgenda := SetupPutAgendaRequest(router, prestadorResp.ID, agendaInput)
	require.Equal(t, http.StatusNoContent, rrAgenda.Code)

	// 4️⃣ Listar prestadores
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	// Validações
	require.Len(t, resp.Data, 1)

	prestador := resp.Data[0]
	assert.NotEmpty(t, prestador.Catalogo, "Deve ter catálogo")
	assert.Len(t, prestador.Catalogo, 1)
	assert.Equal(t, "Corte de Cabelo", prestador.Catalogo[0].Nome)

	assert.NotEmpty(t, prestador.Agenda, "Deve ter agenda")
	assert.Len(t, prestador.Agenda, 1)
	assert.Equal(t, "2030-01-03", prestador.Agenda[0].Data)
	assert.Len(t, prestador.Agenda[0].Intervalos, 1)
}

// TestGetPrestadores_PageZeroAjustaParaUm testa que page=0 é ajustado para 1
func TestGetPrestadores_PageZeroAjustaParaUm(t *testing.T) {
	router, _ := SetupPostPrestador()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=0&limit=10", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 1, resp.Page, "Page deve ser ajustado para 1")
}

// TestGetPrestadores_LimitZeroAjustaParaDez testa que limit=0 é ajustado para 10
func TestGetPrestadores_LimitZeroAjustaParaDez(t *testing.T) {
	router, _ := SetupPostPrestador()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/prestadores?page=1&limit=0", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var resp response_prestador.PrestadorListResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	require.NoError(t, err)

	assert.Equal(t, 10, resp.Limit, "Limit deve ser ajustado para 10")
}
