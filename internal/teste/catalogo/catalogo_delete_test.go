package teste

import (
	"encoding/json"
	"fmt"
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeletarCatalogo_Sucesso(t *testing.T) {
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

	rrDelete := SetupDeleteCatalogoRequest(router, id)
	assert.Equal(t, http.StatusNoContent, rrDelete.Code)

	rrGet := SetupGetCatalogoRequest(router, id)
	assert.Equal(t, http.StatusNotFound, rrGet.Code)

	var resp map[string]string
	err = json.Unmarshal(rrGet.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestDeletarCatalogo_NaoEncontrado(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	rr := SetupDeleteCatalogoRequest(router, "id-inexistente")
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var resp map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	_, exists := resp["error"]
	assert.True(t, exists, "Resposta deve conter campo 'error'")
}

func TestDeletarCatalogo_VerificaListagemAposDelete(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	ids := []string{}
	for i := 1; i <= 3; i++ {
		cat := request_catalogo.CatalogoRequest{
			Nome:          "Serviço " + fmt.Sprintf("%d", i),
			DuracaoPadrao: 30,
			Preco:         1000 * i,
			Categoria:     "Categoria",
			ImagemUrl:     "https://exemplo.com/img.jpg",
		}
		rr := SetupPostCatalogoRequest(router, cat)
		assert.Equal(t, http.StatusCreated, rr.Code)

		var criado map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &criado)
		ids = append(ids, criado["id"].(string))
	}

	rrList := SetupGetCatalogosRequest(router, 1, 10)
	var respList map[string]interface{}
	json.Unmarshal(rrList.Body.Bytes(), &respList)
	assert.Equal(t, float64(3), respList["total"])

	rrDelete := SetupDeleteCatalogoRequest(router, ids[1])
	assert.Equal(t, http.StatusNoContent, rrDelete.Code)

	rrListAfter := SetupGetCatalogosRequest(router, 1, 10)
	var respListAfter map[string]interface{}
	json.Unmarshal(rrListAfter.Body.Bytes(), &respListAfter)
	assert.Equal(t, float64(2), respListAfter["total"])

	data := respListAfter["data"].([]interface{})
	assert.Equal(t, 2, len(data))
}

func TestDeletarCatalogo_DeleteMultiplos(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	ids := []string{}
	for i := 1; i <= 5; i++ {
		cat := request_catalogo.CatalogoRequest{
			Nome:          "Serviço " + fmt.Sprintf("%d", i),
			DuracaoPadrao: 30,
			Preco:         1000,
			Categoria:     "Categoria",
			ImagemUrl:     "https://exemplo.com/img.jpg",
		}
		rr := SetupPostCatalogoRequest(router, cat)
		var criado map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &criado)
		ids = append(ids, criado["id"].(string))
	}

	for _, id := range ids {
		rr := SetupDeleteCatalogoRequest(router, id)
		assert.Equal(t, http.StatusNoContent, rr.Code)
	}

	rrList := SetupGetCatalogosRequest(router, 1, 10)
	var respList map[string]interface{}
	json.Unmarshal(rrList.Body.Bytes(), &respList)
	assert.Equal(t, float64(0), respList["total"])

	data := respList["data"].([]interface{})
	assert.Equal(t, 0, len(data))
}

func TestDeletarCatalogo_TentaDeletarDuasVezes(t *testing.T) {
	router, _ := SetupRouterCatalogo()

	inputPost := request_catalogo.CatalogoRequest{
		Nome:          "Corte de Cabelo",
		DuracaoPadrao: 30,
		Preco:         3500,
		Categoria:     "Beleza",
		ImagemUrl:     "https://exemplo.com/img1.jpg",
	}
	rrPost := SetupPostCatalogoRequest(router, inputPost)
	var criado map[string]interface{}
	json.Unmarshal(rrPost.Body.Bytes(), &criado)
	id := criado["id"].(string)

	rrDelete1 := SetupDeleteCatalogoRequest(router, id)
	assert.Equal(t, http.StatusNoContent, rrDelete1.Code)

	rrDelete2 := SetupDeleteCatalogoRequest(router, id)
	assert.Equal(t, http.StatusNotFound, rrDelete2.Code)
}
