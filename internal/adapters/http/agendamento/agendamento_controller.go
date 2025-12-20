package agendamento

import (
	"errors"
	"meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
	"meu-servico-agenda/internal/adapters/http/agendamento/response_agendamento"
	"meu-servico-agenda/internal/core/application/service"
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

		// 400 — erro de input
		case errors.Is(err, service.ErrDataHoraInvalida),
			errors.Is(err, service.ErrClienteInvalido),
			errors.Is(err, service.ErrPrestadorInvalido),
			errors.Is(err, service.ErrCatalogoInvalido):

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		// 404 — recurso inexistente
		case errors.Is(err, service.ErrClienteNaoExiste),
			errors.Is(err, service.ErrPrestadorNaoExiste),
			errors.Is(err, service.ErrCatalogoNaoExiste):

			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

		// 409 — conflito de agenda
		case errors.Is(err, service.ErrDiaIndisponivel),
			errors.Is(err, service.ErrHorarioIndisponivel),
			errors.Is(err, service.ErrPrestadorOcupado),
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
