package agendamento

import (
	"errors"
	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/response_agendamento"
	"meu-servico-agenda/internal/core/application/service"
	"meu-servico-agenda/internal/core/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgendamentoController struct {
	agendamentoService *service.AgendamentoService
}

func NovoAgendamentoController(ag *service.AgendamentoService) *AgendamentoController {
	return &AgendamentoController{
		agendamentoService: ag,
	}
}

// @Summary Cria um novo agendamento
// @Description Realiza o agendamento de um serviço para um prestador em uma data e horário específicos. Requer autenticação JWT (Cliente ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param agendamento body request_agendamento.AgendamentoRequest true "Dados do agendamento"
// @Success 201 {object} response_agendamento.AgendamentoResponse "Agendamento criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - apenas Clientes e Admin podem criar agendamentos"
// @Failure 404 {object} domain.ErrorResponse "Cliente, prestador ou serviço não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Conflito de agenda (dia ou horário indisponível)"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos [post]
func (ag *AgendamentoController) PostAgendamento(c *gin.Context) {
	var req request_agendamento.AgendamentoRequest
	// 1️⃣ Validação estrutural (JSON)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	input, err := req.ToAgendamento()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2️⃣ Chamada da service
	agendamento, err := ag.agendamentoService.CadastraAgendamento(input)
	if err != nil {

		switch {

		// 404 — recurso inexistente
		case errors.Is(err, service.ErrClienteNaoEncontrado),
			errors.Is(err, service.ErrPrestadorNaoEncontrado),
			errors.Is(err, service.ErrCatalogoNaoEncontrado):

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		// 409 — conflito de agenda
		case errors.Is(err, service.ErrDiaIndisponivel),
			errors.Is(err, service.ErrHorarioIndisponivel),
			errors.Is(err, service.ErrPrestadorOcupado),
			errors.Is(err, domain.ErrDataEstaNoPassado),
			errors.Is(err, service.ErrAgendamentoDuplo),
			errors.Is(err, service.ErrClienteOcupado):

			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

		// 500 — erro inesperado
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": service.ErrFalhaInfraestrutura.Error(),
			})
		}

		return
	}

	// 3️⃣ Sucesso

	resp := response_agendamento.NovoAgendamentoResponse(agendamento)
	c.JSON(http.StatusCreated, resp)
}

// @Summary Busca agendamentos de um cliente a partir de uma data
// @Description Retorna todos os agendamentos de um cliente a partir da data especificada, ordenados por data/hora de início. Requer autenticação JWT (próprio cliente ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cliente"
// @Param data query string true "Data de início da busca (formato: YYYY-MM-DD)" example(2025-01-03)
// @Success 200 {object} response_agendamento.BuscaDataResponse "Lista de agendamentos encontrados"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - você só pode ver seus próprios agendamentos ou ser Admin"
// @Failure 404 {object} domain.ErrorResponse "Cliente não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/cliente/{id} [get]
func (ag *AgendamentoController) GetAgendamentoClienteData(c *gin.Context) {
	id := c.Param("id")

	var input request_agendamento.AgendamentoDataRequest

	// 1️⃣ Validação estrutural (Query params)
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	// 2️⃣ Conversão para input
	req, err := input.ToAgendamentoDataInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "formato de data inválido",
			"detail": err.Error(),
		})
		return
	}

	// 3️⃣ Chamada da service
	agendamentos, err := ag.agendamentoService.ConsultaAgendamentoClienteData(*req, id)
	if err != nil {
		switch {
		// 404 — cliente não existe
		case errors.Is(err, service.ErrClienteNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		// 500 — erro inesperado
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": service.ErrFalhaInfraestrutura.Error(),
			})
		}
		return
	}

	response := response_agendamento.ToBuscaDataResponse(agendamentos)
	c.JSON(http.StatusOK, response)
}

// @Summary Busca agendamentos de um prestador a partir de uma data
// @Description Retorna todos os agendamentos de um prestador a partir da data especificada, ordenados por data/hora de início. Requer autenticação JWT (próprio prestador ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do prestador"
// @Param data query string true "Data de início da busca (formato: YYYY-MM-DD)" example(2025-01-03)
// @Success 200 {object} response_agendamento.BuscaDataResponse "Lista de agendamentos encontrados"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - você só pode ver seus próprios agendamentos ou ser Admin"
// @Failure 404 {object} domain.ErrorResponse "prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/prestador/{id} [get]
func (ag *AgendamentoController) GetAgendamentoPrestadorData(c *gin.Context) {
	id := c.Param("id")

	var input request_agendamento.AgendamentoDataRequest

	// 1️⃣ Validação estrutural (Query params)
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	// 2️⃣ Conversão para input
	req, err := input.ToAgendamentoDataInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "formato de data inválido",
			"detail": err.Error(),
		})
		return
	}

	// 3️⃣ Chamada da service
	agendamentos, err := ag.agendamentoService.ConsultaAgendamentoPrestadorData(*req, id)
	if err != nil {
		switch {
		// 404 — prestador não existe
		case errors.Is(err, service.ErrPrestadorNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		// 500 — erro inesperado
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": service.ErrFalhaInfraestrutura.Error(),
			})
		}
		return
	}

	response := response_agendamento.ToBuscaDataResponse(agendamentos)
	c.JSON(http.StatusOK, response)
}

// @Summary Lista agendamentos de um prestador por período com paginação
// @Description Retorna agendamentos de um prestador em um período específico com suporte a paginação. Requer autenticação JWT (próprio prestador ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do prestador"
// @Param data_inicio query string true "Data de início do período (formato: YYYY-MM-DD)" example(2025-01-01)
// @Param data_fim query string true "Data de fim do período (formato: YYYY-MM-DD)" example(2025-01-31)
// @Param page query int false "Número da página (padrão: 1)"
// @Param limit query int false "Itens por página (padrão: 10, máx: 100)"
// @Success 200 {object} response_agendamento.AgendamentoListPaginadoResponse "Lista paginada de agendamentos"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - você só pode ver seus próprios agendamentos ou ser Admin"
// @Failure 404 {object} domain.ErrorResponse "Prestador não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/prestador/{id}/periodo [get]
func (ag *AgendamentoController) GetAgendamentosPrestadorPeriodo(c *gin.Context) {
	id := c.Param("id")

	var req request_agendamento.AgendamentoPrestadorPeriodoRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	input, err := req.ToInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "formato de data inválido",
			"detail": err.Error(),
		})
		return
	}

	agendamentos, total, err := ag.agendamentoService.ListarAgendamentosPrestadorPorPeriodo(id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPrestadorNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": service.ErrFalhaInfraestrutura.Error(),
			})
		}
		return
	}

	response := response_agendamento.ToAgendamentoListPaginadoResponse(agendamentos, input.Page, input.Limit, total)
	c.JSON(http.StatusOK, response)
}

// @Summary Confirma um agendamento
// @Description Atualiza o status de um agendamento para Confirmado. Requer autenticação JWT (Prestador ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do agendamento"
// @Success 200 {object} response_agendamento.AgendamentoResponse "Agendamento confirmado com sucesso"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - apenas Prestador e Admin podem confirmar agendamentos"
// @Failure 404 {object} domain.ErrorResponse "Agendamento não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Agendamento já foi cancelado ou concluído"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/{id}/confirmar [put]
func (ag *AgendamentoController) ConfirmarAgendamento(c *gin.Context) {
	id := c.Param("id")

	agendamento, err := ag.agendamentoService.ConfirmarAgendamento(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAgendamentoNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrAgendamentoJaCancelado),
			errors.Is(err, service.ErrAgendamentoJaConcluido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrFalhaInfraestrutura.Error()})
		}
		return
	}

	resp := response_agendamento.NovoAgendamentoResponse(agendamento)
	c.JSON(http.StatusOK, resp)
}

// @Summary Cancela um agendamento
// @Description Atualiza o status de um agendamento para Cancelado. Requer autenticação JWT (Prestador, Cliente ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do agendamento"
// @Success 200 {object} response_agendamento.AgendamentoResponse "Agendamento cancelado com sucesso"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - apenas Prestador, Cliente e Admin podem cancelar agendamentos"
// @Failure 404 {object} domain.ErrorResponse "Agendamento não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Agendamento já foi cancelado ou concluído"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/{id}/cancelar [put]
func (ag *AgendamentoController) CancelarAgendamento(c *gin.Context) {
	id := c.Param("id")

	agendamento, err := ag.agendamentoService.CancelarAgendamento(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAgendamentoNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrAgendamentoJaCancelado),
			errors.Is(err, service.ErrAgendamentoJaConcluido):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": service.ErrFalhaInfraestrutura.Error()})
		}
		return
	}

	resp := response_agendamento.NovoAgendamentoResponse(agendamento)
	c.JSON(http.StatusOK, resp)
}

// @Summary Lista agendamentos de um cliente por período com paginação
// @Description Retorna agendamentos de um cliente em um período específico com suporte a paginação. Requer autenticação JWT (próprio cliente ou Admin).
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID do cliente"
// @Param data_inicio query string true "Data de início do período (formato: YYYY-MM-DD)" example(2025-01-01)
// @Param data_fim query string true "Data de fim do período (formato: YYYY-MM-DD)" example(2025-01-31)
// @Param page query int false "Número da página (padrão: 1)"
// @Param limit query int false "Itens por página (padrão: 10, máx: 100)"
// @Success 200 {object} response_agendamento.AgendamentoListPaginadoResponse "Lista paginada de agendamentos"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 401 {object} domain.ErrorResponse "Token não fornecido ou inválido"
// @Failure 403 {object} domain.ErrorResponse "Acesso negado - você só pode ver seus próprios agendamentos ou ser Admin"
// @Failure 404 {object} domain.ErrorResponse "Cliente não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/cliente/{id}/periodo [get]
func (ag *AgendamentoController) GetAgendamentosClientePeriodo(c *gin.Context) {
	id := c.Param("id")

	var req request_agendamento.AgendamentoClientePeriodoRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	input, err := req.ToInput()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "formato de data inválido",
			"detail": err.Error(),
		})
		return
	}

	agendamentos, total, err := ag.agendamentoService.ListarAgendamentosClientePorPeriodo(id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrClienteNaoEncontrado):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": service.ErrFalhaInfraestrutura.Error(),
			})
		}
		return
	}

	response := response_agendamento.ToAgendamentoListPaginadoResponse(agendamentos, input.Page, input.Limit, total)
	c.JSON(http.StatusOK, response)
}
