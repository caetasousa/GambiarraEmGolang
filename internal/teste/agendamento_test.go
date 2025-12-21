package teste

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
	"github.com/stretchr/testify/require"
)

// SetupRouter inicializa router com controllers necessários para testes
func SetupRouterAgendamento() (*gin.Engine, port.PrestadorRepositorio, port.ClienteRepositorio, port.CatalogoRepositorio, port.AgendaDiariaRepositorio) {
	gin.SetMode(gin.TestMode)

	clienteRepo := repository.NewFakeClienteRepositorio()
	prestadorRepo := repository.NovoFakePrestadorRepositorio()
	catalogoRepo := repository.NovoCatalogoFakeRepo()
	agendaDiariaRepo := repository.NovoFakeAgendaDiariaRepositorio()
	agendamentoRepo := repository.NovoFakeAgendamentoRepositorio()

	cadastroCliente := service.NovoServiceCliente(clienteRepo)
	cadastroPrestador := service.NovaPrestadorService(prestadorRepo, catalogoRepo, agendaDiariaRepo)
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
	}

	return router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo
}
func SetupPostAgendamentoRequest(router *gin.Engine, input request_agendamento.AgendamentoRequest) *httptest.ResponseRecorder {
	body, _ := json.Marshal(input)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/agendamentos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}

func SetupNovoCliente(p port.ClienteRepositorio) *domain.Cliente {
	cli, _ := domain.NovoCliente("Eduardo", "caetasousa@gmail.com", "62999697581")
	p.Salvar(cli)
	cliente, _ := p.BuscarPorId(cli.ID)
	return cliente
}

func SetupNovoCatalogo(p port.CatalogoRepositorio) (*domain.Catalogo, *[]domain.Catalogo){
	cat, _ := domain.NovoCatalogo("Manutenção", 60, 20000, "Beleza")
	p.Salvar(cat)
	catalogos := []domain.Catalogo{*cat}
	return cat, &catalogos
}

func SetupCriaPrestador(p port.PrestadorRepositorio, catalogo []domain.Catalogo) *domain.Prestador {
	pres, _ := domain.NovoPrestador("Eduardo", "04423258196", "caetasousa@gmail.com", "662999687481", catalogo)
	p.Salvar(pres)
	return pres
}

func SetupCriaAgendaDiaria(p port.AgendaDiariaRepositorio) *domain.AgendaDiaria {
	horaInicio, _ := time.Parse("15:04", "08:00")
	horaFim, _ := time.Parse("15:04", "12:00")
	data, _ := time.Parse("2006-01-02", "2025-01-03")

	intervalo, _ := domain.NovoIntervaloDiario(horaInicio, horaFim)
	intervalos := []domain.IntervaloDiario{*intervalo}
	agendaDiaria, _ := domain.NovaAgendaDiaria(data, intervalos)
	p.Salvar(agendaDiaria)

	return agendaDiaria
}

func TestPostAgendamento_Sucesso(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	//cadastro do catalogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	//cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	//cria agenda diaria
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	
	//adiciona agenda diaria a prestador
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2025-01-03T08:00:00Z",
		Notas:          "Estou com preça",
	}

	rr := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusCreated, rr.Code)
}

func TestPostAgendamento_ErroPrestadorInexistente(t *testing.T) {
	router, _, clienteRepo, catalogoRepo, _ := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	//cadastro do catalogo
	catalogo, _ := SetupNovoCatalogo(catalogoRepo)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    "prestador-inexistente",
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2025-01-03T08:00:00Z",
		Notas:          "Estou com preça",
	}

	rr := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPostAgendamento_ErroClienteInexistente(t *testing.T) {
	router, prestadorRepo, _, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()
	//cadastro do catalogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	//cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	//cria agenda diaria
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	
	//adiciona agenda diaria a prestador
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      "cliente-inesistente",
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2025-01-03T08:00:00Z",
		Notas:          "Estou com preça",
	}

	rr := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPostAgendamento_ErroCatalogoInexistente(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	//cadastro do catalogo
	_, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	//cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	//cria agenda diaria
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	
	//adiciona agenda diaria a prestador
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     "catalogo-inexistente",
		DataHoraInicio: "2025-01-03T08:00:00Z",
		Notas:          "Estou com preça",
	}

	rr := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusNotFound, rr.Code)
}
