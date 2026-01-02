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
// @Description Realiza o agendamento de um serviço para um prestador em uma data e horário específicos.
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param agendamento body request_agendamento.AgendamentoRequest true "Dados do agendamento"
// @Success 201 {object} response_agendamento.AgendamentoResponse "Agendamento criado com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 404 {object} domain.ErrorResponse "Cliente, prestador ou serviço não encontrado"
// @Failure 409 {object} domain.ErrorResponse "Conflito de agenda (dia ou horário indisponível)"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos [post]
func (ag *AgendamentoController) PostAgendamento(c *gin.Context) {
	var input request_agendamento.AgendamentoRequest
	// 1️⃣ Validação estrutural (JSON)
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	// 2️⃣ Chamada da service
	agendamento, err := ag.agendamentoService.CadastraAgendamento(input)
	if err != nil {

		switch {

		// 404 — recurso inexistente
		case errors.Is(err, service.ErrClienteNaoExiste),
			errors.Is(err, service.ErrPrestadorNaoExiste),
			errors.Is(err, service.ErrCatalogoNaoExiste):

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
// @Description Retorna todos os agendamentos de um cliente a partir da data especificada, ordenados por data/hora de início
// @Tags Agendamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do cliente"
// @Param data query string true "Data de início da busca (formato: YYYY-MM-DD)" example(2025-01-03)
// @Success 200 {object} response_agendamento.BuscaClienteDataResponse "Lista de agendamentos encontrados"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos ou formato de data incorreto"
// @Failure 404 {object} domain.ErrorResponse "Cliente não encontrado"
// @Failure 500 {object} domain.ErrorResponse "Erro interno do servidor"
// @Router /agendamentos/cliente/{id} [get]
func (ag *AgendamentoController) GetAgendamentoClienteData(c *gin.Context) {
	id := c.Param("id")

	var input request_agendamento.AgendamentoClienteDataRequest

	// 1️⃣ Validação estrutural (Query params)
	if err := c.ShouldBindQuery(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "dados inválidos",
			"detail": err.Error(),
		})
		return
	}

	// 2️⃣ Conversão para input
	req, err := input.ToAgendamentoClienteDataInput()
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

	response := response_agendamento.ToBuscaClienteDataResponse(agendamentos)
	c.JSON(http.StatusOK, response)
}
