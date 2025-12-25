package teste

import (
	"encoding/json"
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/core/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtualizarCatalogo_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	assert.Equal(t, http.StatusCreated, rrPost.Code)

	var criado map[string]interface{}
	err := json.Unmarshal(rrPost.Body.Bytes(), &criado)
	assert.NoError(t, err)

	id := criado["id"].(string)

	inputMap := map[string]interface{}{
		"nome":           "Corte Premium",
		"duracao_padrao": 45,
		"preco":          5000,
		"categoria":      "Beleza Premium",
		"image_url":      "https://exemplo.com/img2.jpg",
	}

	body, _ := json.Marshal(inputMap)
	rrPut := PutRawJSON(router, id, body)
	assert.Equal(t, http.StatusNoContent, rrPut.Code)

	rrGet := SetupGetCatalogoRequest(router, id)
	assert.Equal(t, http.StatusOK, rrGet.Code)

	var atualizado map[string]interface{}
	err = json.Unmarshal(rrGet.Body.Bytes(), &atualizado)
	assert.NoError(t, err)

	assert.Equal(t, "Corte Premium", atualizado["nome"])
	assert.Equal(t, float64(45), atualizado["duracao_padrao"])
	assert.Equal(t, float64(5000), atualizado["preco"])
	assert.Equal(t, "Beleza Premium", atualizado["categoria"])
	assert.Equal(t, "https://exemplo.com/img2.jpg", atualizado["image_url"])
}

func TestAtualizarCatalogo_ApenasNome_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Massagem",
		DuracaoPadrao: 60,
		Preco:         8000,
		Categoria:     "Saúde",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	assert.Equal(t, http.StatusCreated, rrPost.Code)

	var criado map[string]interface{}
	err := json.Unmarshal(rrPost.Body.Bytes(), &criado)
	assert.NoError(t, err)

	id := criado["id"].(string)
	duracaoOriginal := criado["duracao_padrao"].(float64)
	precoOriginal := criado["preco"].(float64)

	inputMap := map[string]interface{}{
		"nome":           "Massagem Relaxante",
		"duracao_padrao": 60,
		"preco":          8000,
		"categoria":      "Saúde",
		"image_url":      "https://exemplo.com/img1.jpg",
	}

	body, _ := json.Marshal(inputMap)
	rrPut := PutRawJSON(router, id, body)
	assert.Equal(t, http.StatusNoContent, rrPut.Code)

	rrGet := SetupGetCatalogoRequest(router, id)
	assert.Equal(t, http.StatusOK, rrGet.Code)

	var atualizado map[string]interface{}
	err = json.Unmarshal(rrGet.Body.Bytes(), &atualizado)
	assert.NoError(t, err)

	assert.Equal(t, "Massagem Relaxante", atualizado["nome"])
	assert.Equal(t, duracaoOriginal, atualizado["duracao_padrao"], "Duração não deve mudar")
	assert.Equal(t, precoOriginal, atualizado["preco"], "Preço não deve mudar")
	assert.Equal(t, "Saúde", atualizado["categoria"])
}

func TestAtualizarCatalogo_CatalogoNaoEncontrado(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPut := request_catalogo.CatalogoUpdateRequest{
		Nome:          "Teste",
		DuracaoPadrao: 30,
		Preco:         1000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}

	rr := SetupPutCatalogoRequest(router, "id-inexistente", inputPut)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestAtualizarCatalogo_NomeVazio_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	var criado domain.Catalogo
	json.Unmarshal(rrPost.Body.Bytes(), &criado)

	inputPut := request_catalogo.CatalogoUpdateRequest{
		Nome:          "",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}

	rr := SetupPutCatalogoRequest(router, criado.ID, inputPut)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestAtualizarCatalogo_DuracaoInvalida_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	var criado domain.Catalogo
	json.Unmarshal(rrPost.Body.Bytes(), &criado)

	inputPut := request_catalogo.CatalogoUpdateRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 1,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}

	rr := SetupPutCatalogoRequest(router, criado.ID, inputPut)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestAtualizarCatalogo_PrecoNegativo_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	var criado domain.Catalogo
	json.Unmarshal(rrPost.Body.Bytes(), &criado)

	inputPut := request_catalogo.CatalogoUpdateRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         -100,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}

	rr := SetupPutCatalogoRequest(router, criado.ID, inputPut)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestAtualizarCatalogo_CategoriaMuitoCurta_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	var criado domain.Catalogo
	json.Unmarshal(rrPost.Body.Bytes(), &criado)

	inputPut := request_catalogo.CatalogoUpdateRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "AB",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}

	rr := SetupPutCatalogoRequest(router, criado.ID, inputPut)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestAtualizarCatalogo_JSONInvalido_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Serviço Teste",
		DuracaoPadrao: 30,
		Preco:         2000,
		Categoria:     "Categoria",
		ImagemUrl:     "https://exemplo.com/img.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	var criado domain.Catalogo
	json.Unmarshal(rrPost.Body.Bytes(), &criado)

	jsonInvalido := []byte(`{"nome": "Teste", "duracao_padrao": "abc"}`)
	rr := PutRawJSON(router, criado.ID, jsonInvalido)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}
