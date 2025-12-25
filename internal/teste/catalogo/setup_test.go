package teste

import (
	"bytes"
	"encoding/json"
	"fmt"
	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/catalogo/request_catalogo"
	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func SetupRouterCatalogo() (*gin.Engine, port.CatalogoRepositorio) {
	gin.SetMode(gin.TestMode)

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	cadastroService := service.NovoCatalogoService(catalogoRepo)
	catalogoController := catalogo.NovoCatalogoController(cadastroService)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		apiV1.POST("/catalogos", catalogoController.PostCatalogo)
		apiV1.GET("/catalogos/:id", catalogoController.GetCatalogoPorID)
		apiV1.GET("/catalogos", catalogoController.GetCatalogos)
		apiV1.PUT("/catalogos/:id", catalogoController.Atualizar)
		apiV1.DELETE("/catalogos/:id", catalogoController.Deletar)
	}

	return router, catalogoRepo
}

// ============ POST Helpers ============

func SetupPostCatalogoRequest(router *gin.Engine, input request_catalogo.CatalogoRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalogos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func PostRawJSON(router *gin.Engine, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalogos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func PostJSONFromMap(router *gin.Engine, input map[string]interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	return PostRawJSON(router, body)
}

// ============ GET Helpers ============

func SetupGetCatalogoRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	url := "/api/v1/catalogos/" + id
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func SetupGetCatalogosRequest(router *gin.Engine, page, limit int) *httptest.ResponseRecorder {
	url := "/api/v1/catalogos"
	if page > 0 || limit > 0 {
		url += "?"
		if page > 0 {
			url += "page=" + fmt.Sprintf("%d", page)
		}
		if limit > 0 {
			if page > 0 {
				url += "&"
			}
			url += "limit=" + fmt.Sprintf("%d", limit)
		}
	}
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// ============ PUT Helpers ============

func SetupPutCatalogoRequest(router *gin.Engine, id string, input request_catalogo.CatalogoUpdateRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	url := "/api/v1/catalogos/" + id
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func PutRawJSON(router *gin.Engine, id string, body []byte) *httptest.ResponseRecorder {
	url := "/api/v1/catalogos/" + id
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

func PutJSONFromMap(router *gin.Engine, id string, input map[string]interface{}) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	return PutRawJSON(router, id, body)
}

// ============ DELETE Helpers ============

func SetupDeleteCatalogoRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	url := "/api/v1/catalogos/" + id
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}