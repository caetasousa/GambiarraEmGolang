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

	catalogoRepo := repository.NovoCatalogoFakeRepo()
	prestadorRepo := repository.NovoFakePrestadorRepositorio(catalogoRepo)
	clienteRepo := repository.NewFakeClienteRepositorio()
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

func SetupNovoCatalogo(p port.CatalogoRepositorio) (*domain.Catalogo, *[]domain.Catalogo) {
	cat, _ := domain.NovoCatalogo("Manutenção", 60, 20000, "Beleza", "https://tdfuderuzpylkctxbysu.supabase.co/storage/v1/object/public/imagens/b094865b92ed1821.avif")
	p.Salvar(cat)
	catalogos := []domain.Catalogo{*cat}
	return cat, &catalogos
}

func SetupCriaPrestador(p port.PrestadorRepositorio, catalogo []domain.Catalogo) *domain.Prestador {
	pres, _ := domain.NovoPrestador("Eduardo", "04423258196", "caetasousa@gmail.com", "662999687481","https://exemplo.com/img1.jpg", catalogo)
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
		DataHoraInicio: "2030-01-03T10:01:00Z",
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
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
		DataHoraInicio: "2030-01-03T08:00:00Z",
		Notas:          "Estou com preça",
	}

	rr := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusNotFound, rr.Code)
}

func TestPostAgendamento_PrestadorNãoTrabalhaNesseDia(t *testing.T) {
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
		DataHoraInicio: "2030-01-02T08:00:00Z",
		Notas:          "Estou com preça",
	}

	rr := SetupPostAgendamentoRequest(router, input)
	require.Equal(t, http.StatusConflict, rr.Code)
}

func TestPostAgendamento_PeriodoDeTrabalhoDoPrestadorIndisponivel(t *testing.T) {
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

	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T07:59:00Z",
		Notas:          "Estou com preça",
	}

	rr1 := SetupPostAgendamentoRequest(router, input1)
	require.Equal(t, http.StatusConflict, rr1.Code)

	// Pelo tempo medio de serviço ser de 60 minutos um agendamento as 11:01 ja não e mais permitido
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T11:01:00Z",
		Notas:          "Estou com preça",
	}

	rr2 := SetupPostAgendamentoRequest(router, input2)
	require.Equal(t, http.StatusConflict, rr2.Code)
}

func TestPostAgendamento_PrestadorOcupado(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cliente 1
	cliente1 := SetupNovoCliente(clienteRepo)

	// cliente 2
	cliente2, _ := domain.NovoCliente("Maria", "maria@email.com", "62999999999")
	clienteRepo.Salvar(cliente2)

	// catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	// agenda
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria)

	// primeiro agendamento (válido)
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente1.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rr1 := SetupPostAgendamentoRequest(router, input1)
	require.Equal(t, http.StatusCreated, rr1.Code)

	// segundo agendamento (mesmo horário, mesmo prestador)
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente2.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rr2 := SetupPostAgendamentoRequest(router, input2)
	require.Equal(t, http.StatusConflict, rr2.Code)
}

func TestPostAgendamento_ClienteOcupado(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cliente único
	cliente := SetupNovoCliente(clienteRepo)

	// catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// prestador 1
	prestador1 := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)
	agenda1 := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador1.AdicionarAgenda(agenda1)

	// prestador 2
	prestador2, _ := domain.NovoPrestador(
		"Outro Prestador",
		"12345678900",
		"outro@email.com",
		"62988888888",
		"https://exemplo.com/img1.jpg",
		*listaDeCatalogos,
	)
	prestadorRepo.Salvar(prestador2)
	agenda2 := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador2.AdicionarAgenda(agenda2)

	// primeiro agendamento
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador1.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rr1 := SetupPostAgendamentoRequest(router, input1)
	require.Equal(t, http.StatusCreated, rr1.Code)

	// segundo agendamento (mesmo cliente, mesmo horário)
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador2.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
	}
	rr2 := SetupPostAgendamentoRequest(router, input2)
	require.Equal(t, http.StatusConflict, rr2.Code)
}
