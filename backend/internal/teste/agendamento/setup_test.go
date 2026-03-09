package teste

import (
	"bytes"
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"time"

	"meu-servico-agenda/internal/adapters/http/agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/catalogo"
	"meu-servico-agenda/internal/adapters/http/cliente"

	"meu-servico-agenda/internal/adapters/http/prestador"

	"meu-servico-agenda/internal/adapters/repository"
	"meu-servico-agenda/internal/core/application/port"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

// SetupRouter inicializa router com controllers necessários para testes
func SetupRouterAgendamento() (*gin.Engine, port.PrestadorRepositorio, port.ClienteRepositorio, port.CatalogoRepositorio, port.AgendaDiariaRepositorio) {
	gin.SetMode(gin.TestMode)

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	prestadorRepo := repository.NovoFakePrestadorRepositorio(catalogoRepo)
	clienteRepo := repository.NewFakeClienteRepositorio()
	agendaDiariaRepo := repository.NovoFakeAgendaDiariaRepositorio()
	agendamentoRepo := repository.NovoFakeAgendamentoRepositorio()
	authService := service.NovoAuthService("test-secret-key", 24, clienteRepo, prestadorRepo)
	cadastroCliente := service.NovoServiceCliente(clienteRepo, prestadorRepo)
	cadastroPrestador := service.NovaPrestadorService(prestadorRepo, catalogoRepo, agendaDiariaRepo, authService, clienteRepo)
	cadastraCatalogo := service.NovoCatalogoService(catalogoRepo)
	cadastraAgendamento := service.NovaAgendamentoService(prestadorRepo, agendamentoRepo, catalogoRepo, clienteRepo)

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		clienteController := cliente.NovoClienteController(cadastroCliente)
		prestadorController := prestador.NovoPrestadorController(cadastroPrestador)
		catalogoController := catalogo.NovoCatalogoController(cadastraCatalogo)
		agendamentoController := agendamento.NovoAgendamentoController(cadastraAgendamento)

		apiV1.POST("/clientes", clienteController.PostCliente)
		apiV1.POST("/prestadores", prestadorController.PostPrestador)
		apiV1.PUT("/prestadores/:id/agenda", prestadorController.PutAgenda)
		apiV1.POST("/catalogos", catalogoController.PostCatalogo)
		apiV1.POST("/agendamentos", agendamentoController.PostAgendamento)
		apiV1.GET("/agendamentos/cliente/:id", agendamentoController.GetAgendamentoClienteData)
		apiV1.GET("/agendamentos/cliente/:id/periodo", agendamentoController.GetAgendamentosClientePeriodo)
		apiV1.GET("/agendamentos/prestador/:id", agendamentoController.GetAgendamentoPrestadorData)
		apiV1.GET("/agendamentos/prestador/:id/periodo", agendamentoController.GetAgendamentosPrestadorPeriodo)
		apiV1.PUT("/agendamentos/:id/confirmar", agendamentoController.ConfirmarAgendamento)
		apiV1.PUT("/agendamentos/:id/cancelar", agendamentoController.CancelarAgendamento)
	}

	return router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo
}

// ============ POST Helpers ============

func SetupPostAgendamentoRequest(router *gin.Engine, input request_agendamento.AgendamentoRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/agendamentos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

// ============ GET Helpers ============

func SetupGetAgendamentoClienteDataRequest(router *gin.Engine, clienteID string, data string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/agendamentos/cliente/%s?data=%s", clienteID, data)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupGetAgendamentosClientePeriodoRequest(router *gin.Engine, clienteID, dataInicio, dataFim string, page, limit int) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/agendamentos/cliente/%s/periodo?data_inicio=%s&data_fim=%s&page=%d&limit=%d",
		clienteID, dataInicio, dataFim, page, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupGetAgendamentosPrestadorPeriodoRequest(router *gin.Engine, prestadorID, dataInicio, dataFim string, page, limit int) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/agendamentos/prestador/%s/periodo?data_inicio=%s&data_fim=%s&page=%d&limit=%d",
		prestadorID, dataInicio, dataFim, page, limit)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

// ============ PUT Helpers ============

func SetupConfirmarAgendamentoRequest(router *gin.Engine, agendamentoID string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/agendamentos/%s/confirmar", agendamentoID)
	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupCancelarAgendamentoRequest(router *gin.Engine, agendamentoID string) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/api/v1/agendamentos/%s/cancelar", agendamentoID)
	req, _ := http.NewRequest(http.MethodPut, url, nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

// ============ Domain Setup Helpers ============

func SetupNovoCliente(p port.ClienteRepositorio) *domain.Cliente {
	// Gerar hash da senha para testes
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	cli := domain.NovoCliente(xid.New().String(), "Eduardo", "caetasousa@gmail.com", "62999697581", string(hash))
	p.Salvar(cli)
	cliente, _ := p.BuscarPorId(cli.ID)
	return cliente
}

func SetupNovoCatalogo(p port.CatalogoRepositorio) (*domain.Catalogo, *[]domain.Catalogo) {
	cat, _ := domain.NovoCatalogo(xid.New().String(), "Manutenção", 60, 20000, "Beleza", "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif")
	p.Salvar(cat)
	catalogos := []domain.Catalogo{*cat}
	return cat, &catalogos
}

func SetupCriaPrestador(p port.PrestadorRepositorio, catalogo []domain.Catalogo) *domain.Prestador {
	// Gerar hash da senha para testes
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	pres, err := domain.NovoPrestador(xid.New().String(), "Eduardo", "04423258196", "caetasousa@gmail.com", "662999687481", string(hash), "https://exemplo.com/img1.jpg", catalogo)
	if err != nil {
		panic(err) // Em testes, podemos dar panic se falhar o setup básico
	}
	p.Salvar(pres)
	return pres
}

func SetupCriaAgendaDiaria(p port.AgendaDiariaRepositorio) *domain.AgendaDiaria {
	horaInicio, _ := time.Parse("15:04", "08:00")
	horaFim, _ := time.Parse("15:04", "12:00")
	data, _ := time.Parse("2006-01-02", "2030-01-03")

	intervalo, _ := domain.NovoIntervaloDiario(horaInicio, horaFim)
	intervalos := []domain.IntervaloDiario{*intervalo}
	agendaDiaria, _ := domain.NovaAgendaDiaria(data, intervalos)
	p.Salvar(agendaDiaria, "dasdf")

	return agendaDiaria
}

// Cria múltiplos agendamentos para testes de paginação
func SetupCriaMultiplosAgendamentos(router *gin.Engine, clienteID, prestadorID, catalogoID string, quantidade int) {
	for i := 0; i < quantidade; i++ {
		hora := 8 + i
		input := request_agendamento.AgendamentoRequest{
			ClienteID:      clienteID,
			PrestadorID:    prestadorID,
			CatalogoID:     catalogoID,
			DataHoraInicio: fmt.Sprintf("2030-01-03T%02d:00:00Z", hora),
			Notas:          fmt.Sprintf("Agendamento %d", i+1),
		}
		SetupPostAgendamentoRequest(router, input)
	}
}
