package teste

import (
	"bytes"
	"encoding/json"
	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/request"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupRouterCatalogo() (*gin.Engine, *repository.CatalogoFakeRepo) {
	gin.SetMode(gin.TestMode)

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	cadastroService := service.NovoCatalogoService(catalogoRepo)
	catalogoController := catalogo.NovoCatalogoController(cadastroService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/catalogos", catalogoController.PostPrestador)
		apiV1.GET("/catalogos/:id", catalogoController.GetCatalogoPorID)
	}

	return router, catalogoRepo
}

func SetupPostCatalogoRequest(router *gin.Engine, input request.CatalogoRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalogos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func SetupGetCatalogoRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	url := "/api/v1/catalogos/" + id
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// PostRawJSON envia um payload JSON bruto para o endpoint de criação de catálogo
func PostRawJSON(router *gin.Engine, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalogos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// PostJSONFromMap converte um map em JSON e envia para o endpoint
func PostJSONFromMap(router *gin.Engine, input map[string]interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	return PostRawJSON(router, body)
}

func TestPostCatalogo_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	input := request.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
	}

	rr := SetupPostCatalogoRequest(router, input)

	assert.NotEqual(t, http.StatusBadRequest, rr.Code)
	assert.NotEqual(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp domain.Catalogo
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, input.Nome, resp.Nome)
	assert.Equal(t, input.Categoria, resp.Categoria)
	assert.Equal(t, input.Preco, resp.Preco)
	assert.NotZero(t, resp.ID)
}

func TestPostCatalogo_NomeVazio_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	// Nome vazio — validação de entrada deveria falhar
	input := request.CatalogoRequest{
		Nome:          "",
		DuracaoPadrao: 30,
		Preco:         2500.0,
		Categoria:     "Beleza",
	}

	rr := SetupPostCatalogoRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request quando nome está vazio")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	print(exists)
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCatalogo_DuracaoInvalida_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	// Corpo com tipo inválido para duracao_padrao (string em vez de int)
	input := map[string]interface{}{
		"nome":           "Teste",
		"duracao_padrao": "abc",
		"preco":          1200.0,
		"categoria":      "Serviço",
	}
	rr := PostJSONFromMap(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para tipo inválido de duracao_padrao")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCatalogo_PrecoInvalido_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	// Corpo com tipo inválido para preco (string em vez de number)
	input := map[string]interface{}{
		"nome":           "Teste",
		"duracao_padrao": 45,
		"preco":          "sem-numero",
		"categoria":      "Serviço",
	}
	rr := PostJSONFromMap(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request para tipo inválido de preco")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCatalogo_NomeMuitoCurto_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	// Nome com menos de 3 caracteres
	input := request.CatalogoRequest{
		Nome:          "AB",
		DuracaoPadrao: 30,
		Preco:         2500.0,
		Categoria:     "Beleza",
	}
	rr := SetupPostCatalogoRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para nome com < 3 caracteres")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCatalogo_CategoriaMuitoCurta_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	// Categoria com menos de 3 caracteres
	input := request.CatalogoRequest{
		Nome:          "Corte Premium",
		DuracaoPadrao: 45,
		Preco:         5000.0,
		Categoria:     "AB",
	}
	rr := SetupPostCatalogoRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para categoria com < 3 caracteres")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestGetCatalogo_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	input := request.CatalogoRequest{
		Nome:          "Massagem",
		DuracaoPadrao: 60,
		Preco:         5000.0,
		Categoria:     "Saúde",
	}

	rrPost := SetupPostCatalogoRequest(router, input)
	var criado domain.Catalogo
	err := json.Unmarshal(rrPost.Body.Bytes(), &criado)
	assert.NoError(t, err)

	rrGet := SetupGetCatalogoRequest(router, criado.ID)
	assert.Equal(t, http.StatusOK, rrGet.Code)

	var buscado domain.Catalogo
	err = json.Unmarshal(rrGet.Body.Bytes(), &buscado)
	assert.NoError(t, err)
	assert.Equal(t, criado.ID, buscado.ID)
}

func TestGetCatalogo_NaoEncontrado(t *testing.T) {
	router, _ := SetupRouterCatalogo()
	rr := SetupGetCatalogoRequest(router, "id-inexistente")
	assert.Equal(t, http.StatusNotFound, rr.Code)
}
