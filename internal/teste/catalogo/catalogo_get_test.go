package teste

import (
	"encoding/json"
	"fmt"
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/core/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============ Testes de GET individual ============

func TestGetCatalogo_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	input := request_catalogo.CatalogoRequest{
		Nome:          "Massagem",
		DuracaoPadrao: 60,
		Preco:         5000.0,
		Categoria:     "Saúde",
		ImagemUrl:     "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
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

// ============ Testes de GET listagem (paginação) ============

func TestGetCatalogos_ListagemVazia_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	rr := SetupGetCatalogosRequest(router, 1, 10)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].([]interface{})
	assert.Equal(t, 0, len(data), "Lista deve estar vazia")
	assert.Equal(t, float64(0), resp["total"], "Total deve ser 0")
}

func TestGetCatalogos_ComDados_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	catalogos := []request_catalogo.CatalogoRequest{
		{
			Nome:          "Corte de Cabelo",
			DuracaoPadrao: 30,
			Preco:         3500.0,
			Categoria:     "Beleza",
			ImagemUrl:     "https://exemplo.com/img1.jpg",
		},
		{
			Nome:          "Massagem",
			DuracaoPadrao: 60,
			Preco:         8000.0,
			Categoria:     "Saúde",
			ImagemUrl:     "https://exemplo.com/img2.jpg",
		},
		{
			Nome:          "Manicure",
			DuracaoPadrao: 45,
			Preco:         4500.0,
			Categoria:     "Beleza",
			ImagemUrl:     "https://exemplo.com/img3.jpg",
		},
	}

	for _, cat := range catalogos {
		rr := SetupPostCatalogoRequest(router, cat)
		assert.Equal(t, http.StatusCreated, rr.Code)
	}

	rr := SetupGetCatalogosRequest(router, 1, 10)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].([]interface{})
	assert.Equal(t, 3, len(data), "Deve retornar 3 catálogos")
	assert.Equal(t, float64(3), resp["total"], "Total deve ser 3")
	assert.Equal(t, float64(1), resp["page"], "Página deve ser 1")
	assert.Equal(t, float64(10), resp["limit"], "Limit deve ser 10")
}

func TestGetCatalogos_Paginacao_PrimeiraPagina(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	for i := 1; i <= 5; i++ {
		cat := request_catalogo.CatalogoRequest{
			Nome:          "Serviço " + fmt.Sprintf("%d", i),
			DuracaoPadrao: 30,
			Preco:         1000 * i,
			Categoria:     "Categoria",
			ImagemUrl:     "https://exemplo.com/img.jpg",
		}
		rr := SetupPostCatalogoRequest(router, cat)
		assert.Equal(t, http.StatusCreated, rr.Code)
	}

	rr := SetupGetCatalogosRequest(router, 1, 2)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].([]interface{})
	assert.Equal(t, 2, len(data), "Deve retornar 2 catálogos")
	assert.Equal(t, float64(5), resp["total"], "Total deve ser 5")
	assert.Equal(t, float64(1), resp["page"], "Página deve ser 1")
	assert.Equal(t, float64(2), resp["limit"], "Limit deve ser 2")
}

func TestGetCatalogos_Paginacao_SegundaPagina(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	for i := 1; i <= 5; i++ {
		cat := request_catalogo.CatalogoRequest{
			Nome:          "Serviço " + fmt.Sprintf("%d", i),
			DuracaoPadrao: 30,
			Preco:         1000 * i,
			Categoria:     "Categoria",
			ImagemUrl:     "https://exemplo.com/img.jpg",
		}
		rr := SetupPostCatalogoRequest(router, cat)
		assert.Equal(t, http.StatusCreated, rr.Code)
	}

	rr := SetupGetCatalogosRequest(router, 2, 2)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].([]interface{})
	assert.Equal(t, 2, len(data), "Deve retornar 2 catálogos")
	assert.Equal(t, float64(5), resp["total"], "Total deve ser 5")
	assert.Equal(t, float64(2), resp["page"], "Página deve ser 2")
	assert.Equal(t, float64(2), resp["limit"], "Limit deve ser 2")
}

func TestGetCatalogos_Paginacao_UltimaPaginaParcial(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	for i := 1; i <= 5; i++ {
		cat := request_catalogo.CatalogoRequest{
			Nome:          "Serviço " + fmt.Sprintf("%d", i),
			DuracaoPadrao: 30,
			Preco:         1000 * i,
			Categoria:     "Categoria",
			ImagemUrl:     "https://exemplo.com/img.jpg",
		}
		rr := SetupPostCatalogoRequest(router, cat)
		assert.Equal(t, http.StatusCreated, rr.Code)
	}

	rr := SetupGetCatalogosRequest(router, 3, 2)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].([]interface{})
	assert.Equal(t, 1, len(data), "Deve retornar 1 catálogo (último)")
	assert.Equal(t, float64(5), resp["total"], "Total deve ser 5")
	assert.Equal(t, float64(3), resp["page"], "Página deve ser 3")
}

func TestGetCatalogos_Paginacao_PaginaAlemDoTotal(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	for i := 1; i <= 3; i++ {
		cat := request_catalogo.CatalogoRequest{
			Nome:          "Serviço " + fmt.Sprintf("%d", i),
			DuracaoPadrao: 30,
			Preco:         1000.0,
			Categoria:     "Categoria",
			ImagemUrl:     "https://exemplo.com/img.jpg",
		}
		rr := SetupPostCatalogoRequest(router, cat)
		assert.Equal(t, http.StatusCreated, rr.Code)
	}

	rr := SetupGetCatalogosRequest(router, 10, 10)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	data := resp["data"].([]interface{})
	assert.Equal(t, 0, len(data), "Deve retornar lista vazia")
	assert.Equal(t, float64(3), resp["total"], "Total deve ser 3")
}

func TestGetCatalogos_SemParametros_UsaValoresPadrao(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	cat := request_catalogo.CatalogoRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         1000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}
	rr := SetupPostCatalogoRequest(router, cat)
	assert.Equal(t, http.StatusCreated, rr.Code)

	rr = SetupGetCatalogosRequest(router, 1, 10)
	assert.Equal(t, http.StatusOK, rr.Code)

	var resp map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)

	assert.Equal(t, float64(1), resp["page"], "Page padrão deve ser 1")
	assert.Equal(t, float64(10), resp["limit"], "Limit padrão deve ser 10")
}
