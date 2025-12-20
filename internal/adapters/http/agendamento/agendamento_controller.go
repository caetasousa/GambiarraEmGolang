package agendamento

import (
	"errors"
	requesta_gendamento "meu-servico-agenda/internal/adapters/http/agendamento/request_agendamento"
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

func (ag *AgendamentoController) PostAgendamento(c *gin.Context) {
	var input requesta_gendamento.AgendamentoRequest
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
				"error": "erro interno do servidor",
			})
		}

		return
	}

	// 3️⃣ Sucesso
	c.JSON(http.StatusCreated, agendamento)
}
