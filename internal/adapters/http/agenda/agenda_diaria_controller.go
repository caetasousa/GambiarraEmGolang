package agenda

import (
	"meu-servico-agenda/internal/adapters/http/agenda/request"
	"meu-servico-agenda/internal/adapters/http/agenda/response"
	"meu-servico-agenda/internal/core/application/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AgendaDiariaController struct {
	criarAgendaDiariaService *service.ServiceAgendaDiaria
}

func NovaAgendaDiariaController(criarAgendaDiariaService *service.ServiceAgendaDiaria) *AgendaDiariaController {
	return &AgendaDiariaController{criarAgendaDiariaService: criarAgendaDiariaService}
}

// @Summary Cadastra uma nova agenda diária
// @Description Registra uma agenda diária contendo data e intervalos de horário disponíveis.
// @Tags Agendas
// @Accept json
// @Produce json
// @Param agenda body request.AgendaDiariaRequest true "Dados da Agenda Diária"
// @Success 201 {object} response.AgendaDiariaResponse "Agenda criada com sucesso"
// @Failure 400 {object} domain.ErrorResponse "Dados inválidos (erro de validação)"
// @Failure 409 {object} domain.ErrorResponse "Conflito ao cadastrar agenda"
// @Failure 500 {object} domain.ErrorResponse "Erro interno ao criar agenda"
// @Router /agendas [post]
func (agd *AgendaDiariaController) PostPrestador(c *gin.Context) {
	var input request.AgendaDiariaRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos: " + err.Error()})
		return
	}

	agenda, err := input.ToAgendaDiaria()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	agendaSalva, err := agd.criarAgendaDiariaService.Cadastra(agenda)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	response := response.FromAgendaDiaria(agendaSalva)

	c.JSON(http.StatusCreated, response)

}
