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
	"github.com/stretchr/testify/require"
)

// SetupPostPrestador configura o ambiente de testes com todas as dependências
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
		apiV1.GET("/prestadores", prestadorController.GetPrestadores)
		apiV1.GET("/prestadores/:id", prestadorController.GetPrestador)
		apiV1.PUT("/prestadores/:id/agenda", prestadorController.PutAgenda)
		apiV1.PUT("/prestadores/:id", prestadorController.UpdatePrestador)
		apiV1.PUT("/prestadores/:id/inativar", prestadorController.InativarPrestador) 
		apiV1.PUT("/prestadores/:id/ativar", prestadorController.AtivarPrestador)
		apiV1.DELETE("/prestadores/:id/agenda", prestadorController.DeleteAgenda)

		apiV1.POST("/catalogos", catalogoController.PostCatalogo)
	}

	return router, prestadorRepo
}

// SetupPostPrestadorRequest executa request de criação de prestador
func SetupPostPrestadorRequest(router *gin.Engine, input request_prestador.PrestadorRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/prestadores", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// SetupGetPrestadorRequest executa request de busca de prestador por ID
func SetupGetPrestadorRequest(router *gin.Engine, id string) *httptest.ResponseRecorder {
	url := "/api/v1/prestadores/" + id
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// SetupPostCatalogoRequest executa request de criação de catálogo
func SetupPostCatalogoRequest(router *gin.Engine, input request_catalogo.CatalogoRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/catalogos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// SetupPutAgendaRequest executa request de criação de agenda
func SetupPutAgendaRequest(router *gin.Engine, prestadorID string, input request_prestador.AgendaDiariaRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	url := "/api/v1/prestadores/" + prestadorID + "/agenda"
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// SetupPutPrestadorRequest executa request de atualização de prestador
func SetupPutPrestadorRequest(router *gin.Engine, id string, input request_prestador.PrestadorUpdateRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)
	url := "/api/v1/prestadores/" + id
	req, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// gerarCPFValido retorna um CPF válido baseado no seed
func gerarCPFValido(seed int) string {
	cpfsValidos := []string{
		"91663575002", "10886737087", "55964702015", "35212899079", "42864297094",
		"44187423010", "45537518015", "23204646033", "86306650091", "40933461003",
		"33935391080", "29466173006", "32886059021", "77487008002", "88992049005",
		"06724785014", "90035795042", "39308388001", "86148883090", "26345031054",
		"60314052020", "25176594005",
	}
	return cpfsValidos[seed%len(cpfsValidos)]
}

// CriarCatalogoValido cria um catálogo válido para testes
func CriarCatalogoValido(t *testing.T, router *gin.Engine) response_catalogo.CatalogoResponse {
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

	return catalogoResp
}

// CriarPrestadorValido cria um prestador válido para testes
func CriarPrestadorValido(t *testing.T, router *gin.Engine, catalogoID string, cpf string) response_prestador.PrestadorResponse {
	prestadorInput := request_prestador.PrestadorRequest{
		Nome:        "João da Silva",
		Cpf:         cpf,
		Email:       "joao@email.com",
		Telefone:    "62999677481",
		ImagemUrl:   "https://exemplo.com/img1.jpg",
		CatalogoIDs: []string{catalogoID},
	}
	rrPrestador := SetupPostPrestadorRequest(router, prestadorInput)
	require.Equal(t, http.StatusCreated, rrPrestador.Code)

	var prestadorResp response_prestador.PrestadorResponse
	err := json.Unmarshal(rrPrestador.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	return prestadorResp
}

// CriarPrestadorValidoParaTeste cria prestador completo com catálogo (mantido para compatibilidade)
func CriarPrestadorValidoParaTeste(t *testing.T) (*gin.Engine, domain.Prestador, port.PrestadorRepositorio) {
	router, repo := SetupPostPrestador()
	catalogoResp := CriarCatalogoValido(t, router)

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
	err := json.Unmarshal(rrPrestador.Body.Bytes(), &prestadorResp)
	require.NoError(t, err)

	return router, prestadorResp, repo
}
