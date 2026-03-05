package teste

import (
	"encoding/json"
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/core/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostCatalogo_Sucesso(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	input := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
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

	input := request_catalogo.CatalogoRequest{
		Nome:          "",
		DuracaoPadrao: 30,
		Preco:         2500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
	}

	rr := SetupPostCatalogoRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 Bad Request quando nome está vazio")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestPostCatalogo_DuracaoInvalida_DeveRetornar400(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	input := map[string]interface{}{
		"nome":           "Teste",
		"duracao_padrao": "abc",
		"preco":          1200.0,
		"categoria":      "Serviço",
		"ImagemUrl":      "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
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

	input := map[string]interface{}{
		"nome":           "Teste",
		"duracao_padrao": 45,
		"preco":          "sem-numero",
		"categoria":      "Serviço",
		"ImagemUrl":      "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
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

	input := request_catalogo.CatalogoRequest{
		Nome:          "AB",
		DuracaoPadrao: 30,
		Preco:         2500.0,
		Categoria:     "Beleza",
		ImagemUrl:     "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
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

	input := request_catalogo.CatalogoRequest{
		Nome:          "Corte Premium",
		DuracaoPadrao: 45,
		Preco:         5000.0,
		Categoria:     "AB",
		ImagemUrl:     "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif",
	}
	rr := SetupPostCatalogoRequest(router, input)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Esperado 400 para categoria com < 3 caracteres")
	var resp map[string]string
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}