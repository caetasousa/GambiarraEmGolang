package teste

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/response_agendamento"
	"meu-servico-agenda/internal/core/domain"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

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
	hash, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	cliente2 := domain.NovoCliente("Maria", "maria@email.com", "62999999999", string(hash))
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
	hash2, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)
	prestador2 := domain.NovoPrestador(
		"Outro Prestador",
		"12345678900",
		"outro@email.com",
		"62988888888",
		string(hash2),
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

func TestPostAgendamento_AgendamentoDuploMesmaCategoriaMesmoDia(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	// cadastro do catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	// cria agenda diária
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)

	// adiciona agenda diária ao prestador
	prestador.AdicionarAgenda(agendaDiaria)

	// primeiro agendamento (válido) - 08:00 às 09:00
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
		Notas:          "Primeiro agendamento",
	}
	rr1 := SetupPostAgendamentoRequest(router, input1)
	require.Equal(t, http.StatusCreated, rr1.Code)

	// segundo agendamento (mesmo dia, mesma categoria, horário diferente) - 10:00 às 11:00
	// Deve falhar porque é a mesma categoria no mesmo dia
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Segundo agendamento - mesma categoria",
	}
	rr2 := SetupPostAgendamentoRequest(router, input2)
	require.Equal(t, http.StatusConflict, rr2.Code)

	// Verificar mensagem de erro específica
	var response map[string]interface{}
	err := json.Unmarshal(rr2.Body.Bytes(), &response)
	require.NoError(t, err)
	require.Contains(t, response["error"], "já existe um agendamento para essa categoria neste dia")
}

func TestPostAgendamento_CategoriasDiferentesMesmoDia_Permitido(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	// cadastro do primeiro catálogo (Manutenção)
	catalogo1, listaDeCatalogos1 := SetupNovoCatalogo(catalogoRepo)

	// cadastro do segundo catálogo (categoria diferente)
	catalogo2, _ := domain.NovoCatalogo(
		"Corte de Cabelo",
		30,
		5000,
		"Beleza",
		"https://exemplo.com/img2.jpg",
	)
	catalogoRepo.Salvar(catalogo2)

	// adiciona o segundo catálogo à lista de catálogos do prestador
	catalogosCompletos := append(*listaDeCatalogos1, *catalogo2)

	// cria prestador com ambos os catálogos
	prestador := SetupCriaPrestador(prestadorRepo, catalogosCompletos)

	// cria agenda diária
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)

	// adiciona agenda diária ao prestador
	prestador.AdicionarAgenda(agendaDiaria)

	// primeiro agendamento - categoria "Manutenção"
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo1.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
		Notas:          "Manutenção",
	}
	rr1 := SetupPostAgendamentoRequest(router, input1)
	require.Equal(t, http.StatusCreated, rr1.Code)

	// segundo agendamento - categoria "Corte de Cabelo" (diferente)
	// Deve ser permitido porque é categoria diferente
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo2.ID,
		DataHoraInicio: "2030-01-03T09:30:00Z",
		Notas:          "Corte de cabelo",
	}
	rr2 := SetupPostAgendamentoRequest(router, input2)
	require.Equal(t, http.StatusCreated, rr2.Code)
}

func TestPostAgendamento_MesmaCategoriaDiasDiferentes_Permitido(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	// cadastro do catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	// cria agenda diária para o dia 03
	agendaDiaria1 := SetupCriaAgendaDiaria(agendaDiariaRepo)
	prestador.AdicionarAgenda(agendaDiaria1)

	// cria agenda diária para o dia 04
	horaInicio, _ := time.Parse("15:04", "08:00")
	horaFim, _ := time.Parse("15:04", "12:00")
	data2, _ := time.Parse("2006-01-02", "2030-01-04")
	intervalo, _ := domain.NovoIntervaloDiario(horaInicio, horaFim)
	intervalos := []domain.IntervaloDiario{*intervalo}
	agendaDiaria2, _ := domain.NovaAgendaDiaria(data2, intervalos)
	agendaDiariaRepo.Salvar(agendaDiaria2, prestador.ID)
	prestador.AdicionarAgenda(agendaDiaria2)

	// primeiro agendamento - dia 03
	input1 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T08:00:00Z",
		Notas:          "Dia 03",
	}
	rr1 := SetupPostAgendamentoRequest(router, input1)
	require.Equal(t, http.StatusCreated, rr1.Code)

	// segundo agendamento - dia 04 (mesma categoria)
	// Deve ser permitido porque é outro dia
	input2 := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-04T08:00:00Z",
		Notas:          "Dia 04",
	}
	rr2 := SetupPostAgendamentoRequest(router, input2)
	require.Equal(t, http.StatusCreated, rr2.Code)
}

func TestPostAgendamento_ValidaResponse(t *testing.T) {
	router, prestadorRepo, clienteRepo, catalogoRepo, agendaDiariaRepo := SetupRouterAgendamento()

	// cadastro do cliente
	cliente := SetupNovoCliente(clienteRepo)

	// cadastro do catálogo
	catalogo, listaDeCatalogos := SetupNovoCatalogo(catalogoRepo)

	// cria prestador
	prestador := SetupCriaPrestador(prestadorRepo, *listaDeCatalogos)

	// cria agenda diária
	agendaDiaria := SetupCriaAgendaDiaria(agendaDiariaRepo)

	// adiciona agenda diária ao prestador
	prestador.AdicionarAgenda(agendaDiaria)

	input := request_agendamento.AgendamentoRequest{
		ClienteID:      cliente.ID,
		PrestadorID:    prestador.ID,
		CatalogoID:     catalogo.ID,
		DataHoraInicio: "2030-01-03T10:00:00Z",
		Notas:          "Estou com pressa",
	}

	rr := SetupPostAgendamentoRequest(router, input)

	// Verifica status code
	require.Equal(t, http.StatusCreated, rr.Code)

	// Parse da resposta
	var response response_agendamento.AgendamentoResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Valida que foi gerado um ID
	require.NotEmpty(t, response.ID)

	// Valida dados do cliente
	require.Equal(t, cliente.ID, response.Cliente.ID)
	require.Equal(t, cliente.Nome, response.Cliente.Nome)
	require.Equal(t, cliente.Email, response.Cliente.Email)
	require.Equal(t, cliente.Telefone, response.Cliente.Telefone)

	// Valida dados do prestador
	require.Equal(t, prestador.ID, response.Prestador.ID)
	require.Equal(t, prestador.Nome, response.Prestador.Nome)
	require.Equal(t, prestador.Cpf, response.Prestador.CPF)
	require.Equal(t, prestador.Email, response.Prestador.Email)
	require.Equal(t, prestador.Telefone, response.Prestador.Telefone)
	require.True(t, response.Prestador.Ativo)

	// Valida dados do serviço
	require.Equal(t, catalogo.ID, response.Servico.ID)
	require.Equal(t, catalogo.Nome, response.Servico.Nome)
	require.Equal(t, catalogo.DuracaoPadrao, response.Servico.Duracao)
	require.Equal(t, catalogo.Preco, response.Servico.Preco)
	require.Equal(t, catalogo.Categoria, response.Servico.Categoria)

	// Valida datas
	expectedInicio, _ := time.Parse(time.RFC3339, "2030-01-03T10:00:00Z")
	expectedFim := expectedInicio.Add(time.Duration(catalogo.DuracaoPadrao) * time.Minute)

	require.True(t, response.DataInicio.Equal(expectedInicio),
		"Data de início esperada: %v, recebida: %v", expectedInicio, response.DataInicio)
	require.True(t, response.DataFim.Equal(expectedFim),
		"Data de fim esperada: %v, recebida: %v", expectedFim, response.DataFim)

	// Valida status
	require.Equal(t, domain.Pendente, response.Status)

	// Valida notas
	require.Equal(t, "Estou com pressa", response.Notas)
}
